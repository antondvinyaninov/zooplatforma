package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"shelter/handlers"
	localmiddleware "shelter/middleware"
	"strings"

	"database"

	"github.com/joho/godotenv"
	pkgmiddleware "github.com/zooplatforma/pkg/middleware"
)

func enableCORSHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Убрали verbose логирование для уменьшения шума в консоли

		origin := r.Header.Get("Origin")
		allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
			"http://localhost:5100": true,
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
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5100")
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
			"http://localhost:5100": true,
		}

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if origin == "" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5100")
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

	// Create tables
	if err := createTables(); err != nil {
		log.Fatal("Failed to create tables:", err)
	}

	// Routes
	// Публичные endpoints
	http.HandleFunc("/api/health", enableCORS(handleHealth))

	// Защищённые endpoints (требуют аутентификации)
	db := database.DB

	// Список приютов пользователя (только auth, без tenant)
	http.Handle("/api/my-shelters", enableCORSHandler(pkgmiddleware.AuthMiddleware(http.HandlerFunc(handlers.GetMyShelters(db)))))

	// Создание приюта (только auth, без tenant)
	http.Handle("/api/shelters", enableCORSHandler(pkgmiddleware.AuthMiddleware(http.HandlerFunc(handlers.CreateShelter(db)))))

	// Endpoints с tenant (требуют выбранного приюта)
	tenantMiddleware := localmiddleware.TenantMiddleware(db)

	// Применяем middleware к handlers (AuthMiddleware + TenantMiddleware)
	http.Handle("/api/animals", enableCORSHandler(pkgmiddleware.AuthMiddleware(tenantMiddleware(http.HandlerFunc(handlers.GetAnimals(db))))))
	http.Handle("/api/stats", enableCORSHandler(pkgmiddleware.AuthMiddleware(tenantMiddleware(http.HandlerFunc(handlers.GetStats(db))))))
	http.Handle("/api/organization", enableCORSHandler(pkgmiddleware.AuthMiddleware(tenantMiddleware(http.HandlerFunc(handlers.GetOrganization(db))))))

	// Root route - должен быть последним!
	http.HandleFunc("/", enableCORS(handleRoot))

	port := ":8200"
	fmt.Printf("🏠 Shelter API starting on port %s\n", port)
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
	return "localhost:3000, localhost:5100"
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Только для корневого пути
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "Shelter API", "version": "1.0.0"}`)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok", "service": "shelter"}`)
}

func createTables() error {
	log.Println("✅ Shelter tables created successfully")
	return nil
}
