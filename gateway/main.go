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

	// ✅ КРИТИЧНО: Middleware в правильном порядке
	// 1. CORS - ПЕРВЫМ! Обрабатывает OPTIONS и устанавливает заголовки
	r.Use(CORSMiddleware)
	// 2. Logging - логирует все запросы
	r.Use(LoggingMiddleware)
	// 3. Rate Limiting - ограничивает количество запросов
	r.Use(RateLimitMiddleware)

	// ✅ КРИТИЧНО: Глобальный обработчик OPTIONS для всех маршрутов
	// Это нужно чтобы preflight запросы не получали 405 Method Not Allowed
	r.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS заголовки уже установлены в CORSMiddleware
		// Просто возвращаем 200 OK
		w.WriteHeader(http.StatusOK)
	})

	// Health check
	r.HandleFunc("/health", HealthCheckHandler(services)).Methods("GET")
	r.HandleFunc("/api/health", HealthCheckHandler(services)).Methods("GET")

	// Auth endpoints (встроенные в Gateway)
	r.HandleFunc("/api/auth/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/api/auth/login", LoginHandler).Methods("POST")
	r.HandleFunc("/api/auth/logout", LogoutHandler).Methods("POST")
	r.HandleFunc("/api/auth/me", GetMeHandler).Methods("GET")

	// PetBase Backend (публичный) - регистрируем ПЕРВЫМ
	petbaseRouter := r.PathPrefix("/api/petbase").Subrouter()
	petbaseRouter.PathPrefix("").HandlerFunc(ProxyHandler(services.PetBase))

	// Clinic Backend (защищенный)
	clinicRouter := r.PathPrefix("/api/clinic").Subrouter()
	clinicRouter.Use(AuthMiddleware)
	clinicRouter.PathPrefix("").HandlerFunc(ProxyHandler(services.Clinic))

	// Owner Backend (защищенный)
	ownerRouter := r.PathPrefix("/api/owner").Subrouter()
	ownerRouter.Use(AuthMiddleware)
	ownerRouter.PathPrefix("").HandlerFunc(ProxyHandler(services.Owner))

	// Shelter Backend (защищенный)
	shelterRouter := r.PathPrefix("/api/shelter").Subrouter()
	shelterRouter.Use(AuthMiddleware)
	shelterRouter.PathPrefix("").HandlerFunc(ProxyHandler(services.Shelter))

	// Volunteer Backend (защищенный)
	volunteerRouter := r.PathPrefix("/api/volunteer").Subrouter()
	volunteerRouter.Use(AuthMiddleware)
	volunteerRouter.PathPrefix("").HandlerFunc(ProxyHandler(services.Volunteer))

	// Admin Backend (защищенный + только админы)
	adminRouter := r.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(AuthMiddleware)
	adminRouter.Use(AdminOnlyMiddleware)
	adminRouter.PathPrefix("").HandlerFunc(ProxyHandler(services.Admin))

	// Main Backend - публичные endpoints (без авторизации)
	// Регистрируем ПОСЛЕ специфичных маршрутов (/api/petbase, /api/clinic и т.д.)
	r.HandleFunc("/api/posts", ProxyHandler(services.Main)).Methods("GET")
	r.HandleFunc("/api/users/{id:[0-9]+}", ProxyHandler(services.Main)).Methods("GET")
	r.HandleFunc("/api/organizations/all", ProxyHandler(services.Main)).Methods("GET")
	r.HandleFunc("/api/species", ProxyHandler(services.Main)).Methods("GET")
	r.HandleFunc("/api/breeds", ProxyHandler(services.Main)).Methods("GET")

	// Main Backend - защищенные endpoints (требуют авторизации)
	// Используем конкретные маршруты вместо PathPrefix("/api")
	mainProtected := r.NewRoute().Subrouter()
	mainProtected.Use(AuthMiddleware)
	mainProtected.HandleFunc("/api/profile", ProxyHandler(services.Main)).Methods("GET", "PUT")
	mainProtected.HandleFunc("/api/posts", ProxyHandler(services.Main)).Methods("POST")
	mainProtected.HandleFunc("/api/posts/{id:[0-9]+}", ProxyHandler(services.Main)).Methods("PUT", "DELETE")
	mainProtected.HandleFunc("/api/comments", ProxyHandler(services.Main)).Methods("POST", "PUT", "DELETE")
	mainProtected.HandleFunc("/api/comments/{id:[0-9]+}", ProxyHandler(services.Main)).Methods("PUT", "DELETE")
	mainProtected.HandleFunc("/api/likes", ProxyHandler(services.Main)).Methods("POST", "DELETE")
	mainProtected.HandleFunc("/api/follows", ProxyHandler(services.Main)).Methods("POST", "DELETE")
	mainProtected.HandleFunc("/api/notifications", ProxyHandler(services.Main)).Methods("GET", "PUT")
	mainProtected.HandleFunc("/api/messages", ProxyHandler(services.Main)).Methods("GET", "POST")
	mainProtected.HandleFunc("/api/chats", ProxyHandler(services.Main)).Methods("GET", "POST")
	mainProtected.HandleFunc("/api/chats/{id:[0-9]+}", ProxyHandler(services.Main)).Methods("GET", "PUT", "DELETE")
	mainProtected.HandleFunc("/api/chats/{id:[0-9]+}/messages", ProxyHandler(services.Main)).Methods("GET", "POST")

	// Статические файлы (uploads)
	uploadsDir := os.Getenv("UPLOAD_PATH")
	if uploadsDir == "" {
		uploadsDir = "/app/uploads"
	}
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadsDir))))

	// ⚠️ ВАЖНО: Frontend проксирование - ПОСЛЕДНИЙ маршрут!
	// Проксирует все запросы (кроме /api/*, /uploads/*, /health) на Main Service
	// Main Service внутри контейнера имеет nginx который направляет:
	//   - /api/* → Backend (localhost:8000)
	//   - /* → Frontend (localhost:3000)
	r.PathPrefix("/").HandlerFunc(ProxyHandler(services.Main))

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
