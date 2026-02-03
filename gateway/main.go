package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузить .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found")
	}

	// Инициализировать JWT Secret
	InitJWT()

	// Инициализировать БД для авторизации
	if err := InitAuthDB(); err != nil {
		log.Fatal("❌ Failed to connect to auth database:", err)
	}
	defer authDB.Close()

	// Инициализировать сервисы
	services := InitServices()

	// Создать роутер
	r := mux.NewRouter()

	// Middleware
	r.Use(LoggingMiddleware)
	r.Use(CORSMiddleware)
	r.Use(RateLimitMiddleware)

	// Health check
	r.HandleFunc("/health", HealthCheckHandler(services)).Methods("GET")
	r.HandleFunc("/api/health", HealthCheckHandler(services)).Methods("GET")

	// Auth endpoints (встроенные в Gateway)
	r.HandleFunc("/api/auth/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/auth/login", LoginHandler).Methods("POST")
	r.HandleFunc("/api/auth/logout", LogoutHandler).Methods("POST")
	r.HandleFunc("/api/auth/me", GetMeHandler).Methods("GET")

	// Main Backend - публичные endpoints (без авторизации)
	r.HandleFunc("/api/posts", ProxyHandler(services.Main)).Methods("GET")
	r.HandleFunc("/api/users/{id:[0-9]+}", ProxyHandler(services.Main)).Methods("GET")
	r.HandleFunc("/api/organizations/all", ProxyHandler(services.Main)).Methods("GET")
	r.HandleFunc("/api/species", ProxyHandler(services.Main)).Methods("GET")
	r.HandleFunc("/api/breeds", ProxyHandler(services.Main)).Methods("GET")

	// Main Backend - все остальные endpoints (требуют авторизации)
	mainProtected := r.NewRoute().Subrouter()
	mainProtected.Use(AuthMiddleware)
	mainProtected.PathPrefix("/api").HandlerFunc(ProxyHandler(services.Main))

	// PetBase Backend
	petbaseRouter := r.PathPrefix("/api/petbase").Subrouter()
	petbaseRouter.PathPrefix("").HandlerFunc(ProxyHandler(services.PetBase))

	// Clinic Backend
	clinicRouter := r.PathPrefix("/api/clinic").Subrouter()
	clinicRouter.Use(AuthMiddleware)
	clinicRouter.PathPrefix("").HandlerFunc(ProxyHandler(services.Clinic))

	// Owner Backend
	ownerRouter := r.PathPrefix("/api/owner").Subrouter()
	ownerRouter.Use(AuthMiddleware)
	ownerRouter.PathPrefix("").HandlerFunc(ProxyHandler(services.Owner))

	// Shelter Backend
	shelterRouter := r.PathPrefix("/api/shelter").Subrouter()
	shelterRouter.Use(AuthMiddleware)
	shelterRouter.PathPrefix("").HandlerFunc(ProxyHandler(services.Shelter))

	// Volunteer Backend
	volunteerRouter := r.PathPrefix("/api/volunteer").Subrouter()
	volunteerRouter.Use(AuthMiddleware)
	volunteerRouter.PathPrefix("").HandlerFunc(ProxyHandler(services.Volunteer))

	// Admin Backend
	adminRouter := r.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(AuthMiddleware)
	adminRouter.Use(AdminOnlyMiddleware)
	adminRouter.PathPrefix("").HandlerFunc(ProxyHandler(services.Admin))

	// Статические файлы (uploads)
	uploadsDir := os.Getenv("UPLOAD_PATH")
	if uploadsDir == "" {
		uploadsDir = "/app/uploads"
	}
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	// Frontend (все остальные запросы проксируем на Next.js)
	// Next.js сам обработает роутинг страниц
	r.PathPrefix("/").HandlerFunc(ProxyHandler(&Service{
		Name:    "Main Frontend",
		URL:     "http://localhost:3000",
		Timeout: 30,
	}))

	// Запустить сервер
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "80"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("🚀 API Gateway started on port %s", port)
	log.Printf("📋 Services:")
	log.Printf("   - Auth Service: %s", services.Auth.URL)
	log.Printf("   - Main Backend: %s", services.Main.URL)
	log.Printf("   - PetBase Backend: %s", services.PetBase.URL)
	log.Printf("   - Clinic Backend: %s", services.Clinic.URL)
	log.Printf("   - Owner Backend: %s", services.Owner.URL)
	log.Printf("   - Shelter Backend: %s", services.Shelter.URL)
	log.Printf("   - Volunteer Backend: %s", services.Volunteer.URL)
	log.Printf("   - Admin Backend: %s", services.Admin.URL)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
