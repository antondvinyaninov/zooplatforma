package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"volunteer/handlers"

	"database"

	"github.com/joho/godotenv"
	"github.com/zooplatforma/pkg/middleware"
)

func enableCORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Убрали verbose логирование для уменьшения шума в консоли

		origin := r.Header.Get("Origin")
		allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
			"http://localhost:6200": true,
		}

		if allowedOriginsEnv != "" {
			for _, o := range strings.Split(allowedOriginsEnv, ",") {
				allowedOrigins[strings.TrimSpace(o)] = true
			}
		}

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Printf("✅ Origin allowed: %s", origin)
		} else if origin == "" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:6200")
			// Убрали verbose логирование
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			log.Printf("✅ OPTIONS request handled")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Убрали verbose логирование
		next.ServeHTTP(w, r)
	})
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
			"http://localhost:6200": true,
		}

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if origin == "" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:6200")
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Routes
	// Публичные endpoints
	http.HandleFunc("/api/health", enableCORS(handleHealth))

	// Защищённые endpoints (требуют аутентификации)
	db := database.DB

	// Volunteer endpoints
	http.Handle("/api/my-tasks", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.GetMyTasks(db)))))
	http.Handle("/api/my-pets", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.GetMyPets(db)))))
	http.Handle("/api/profile", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.GetProfile(db)))))
	http.Handle("/api/take-custody", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.TakeCustody(db)))))
	http.Handle("/api/release-custody", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.ReleaseCustody(db)))))
	http.Handle("/api/pets", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.CreatePet(db)))))

	// Root route - должен быть последним!
	http.HandleFunc("/", enableCORS(handleRoot))

	port := ":8500"
	fmt.Printf("🤝 Volunteer API starting on port %s\n", port)
	fmt.Printf("🔒 JWT Authentication: %s\n", getAuthStatus())
	fmt.Printf("🌐 CORS: %s\n", getAllowedOrigins())
	log.Fatal(http.ListenAndServe(port, nil))
}

func getAuthStatus() string {
	if os.Getenv("JWT_SECRET") != "" {
		return "Enabled (production mode)"
	}
	return "Enabled (development mode - using default secret)"
}

func getAllowedOrigins() string {
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins != "" {
		return origins
	}
	return "localhost:3000, localhost:6200"
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Только для корневого пути
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "Volunteer API", "version": "1.0.0"}`)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok", "service": "volunteer"}`)
}
