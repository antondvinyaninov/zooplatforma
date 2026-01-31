package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"owner/handlers"
	"strings"

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
			"http://localhost:6100": true,
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
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:6100")
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
			"http://localhost:6100": true,
		}

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if origin == "" {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:6100")
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

	// ✅ Auth Service URL (middleware читает из AUTH_SERVICE_URL в .env)
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "http://localhost:7100"
	}
	log.Printf("🔐 Auth Service URL: %s", authServiceURL)

	// Routes
	// Публичные endpoints
	http.HandleFunc("/api/health", enableCORS(handleHealth))

	// Защищённые endpoints (требуют аутентификации)
	db := database.DB

	// Owner endpoints - прямые запросы к базе данных
	http.Handle("/api/pets", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handlers.GetMyPets(db)(w, r)
		case "POST":
			handlers.CreatePet(db)(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}))))
	http.Handle("/api/pets/", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🔀 Route /api/pets/ - Method: %s, Path: %s", r.Method, r.URL.Path)
		// Проверяем, это запрос на загрузку фото или получение питомца
		if strings.HasSuffix(r.URL.Path, "/photo") {
			switch r.Method {
			case "POST":
				log.Printf("📸 Routing to UploadPetPhoto handler")
				handlers.UploadPetPhoto(db)(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		} else {
			switch r.Method {
			case "GET":
				log.Printf("🐾 Routing to GetPet handler")
				handlers.GetPet(db)(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}
	}))))
	http.Handle("/api/profile", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.GetProfile(db)))))
	http.Handle("/api/breeds", enableCORSHandler(middleware.AuthMiddleware(http.HandlerFunc(handlers.GetBreeds(db)))))

	// Static files - раздача загруженных файлов из общей папки
	http.Handle("/uploads/", enableCORSHandler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("../../uploads")))))
	log.Println("📁 Static files: /uploads/ -> ../../uploads")

	// Root route - должен быть последним!
	http.HandleFunc("/", enableCORS(handleRoot))

	port := ":8400"
	fmt.Printf("👤 Owner API starting on port %s\n", port)
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
	return "localhost:3000, localhost:6100"
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Только для корневого пути
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "Owner API", "version": "1.0.0"}`)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok", "service": "owner"}`)
}

func createTables() error {
	db := database.DB

	// Create treatments table
	query := `
	CREATE TABLE IF NOT EXISTS treatments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		pet_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		date DATETIME NOT NULL,
		medication TEXT NOT NULL,
		dosage TEXT NOT NULL,
		next_date DATETIME,
		notes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_treatments_pet_id ON treatments(pet_id);
	CREATE INDEX IF NOT EXISTS idx_treatments_user_id ON treatments(user_id);
	CREATE INDEX IF NOT EXISTS idx_treatments_next_date ON treatments(next_date);
	CREATE INDEX IF NOT EXISTS idx_treatments_date ON treatments(date);
	`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create treatments table: %w", err)
	}

	// Create trigger for updated_at
	triggerQuery := `
	CREATE TRIGGER IF NOT EXISTS update_treatments_timestamp 
	AFTER UPDATE ON treatments
	FOR EACH ROW
	BEGIN
		UPDATE treatments SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
	END;
	`

	if _, err := db.Exec(triggerQuery); err != nil {
		return fmt.Errorf("failed to create treatments trigger: %w", err)
	}

	log.Println("✅ Owner tables created successfully")
	return nil
}
