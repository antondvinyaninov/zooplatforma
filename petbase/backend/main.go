package main

import (
	"database"
	"fmt"
	"log"
	"net/http"
	"os"
	"petbase/handlers"
	"strings"

	"github.com/joho/godotenv"
	"github.com/zooplatforma/pkg/middleware"
)

// enableCORS - для http.HandlerFunc
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Убрали verbose логирование для уменьшения шума в консоли

		origin := r.Header.Get("Origin")

		// Получаем разрешённые origins из переменной окружения
		allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
		allowedOrigins := map[string]bool{
			"http://localhost:3000": true, // Основной сайт (development)
			"http://localhost:4100": true, // ЗооБаза Frontend (development)
			"http://localhost:6100": true, // Кабинет владельца (development)
			"http://localhost:6200": true, // Кабинет волонтёра (development)
			"http://localhost:6300": true, // Кабинет клиники (development)
		}

		// Добавляем origins из .env
		if allowedOriginsEnv != "" {
			for _, o := range strings.Split(allowedOriginsEnv, ",") {
				allowedOrigins[strings.TrimSpace(o)] = true
			}
		}

		// Проверяем, разрешён ли origin
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Printf("✅ Origin allowed: %s", origin)
		} else if origin == "" {
			// Если origin не указан, используем дефолтный
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4100")
			// Убрали verbose логирование
		} else {
			// Origin не разрешён
			log.Printf("❌ Blocked request from unauthorized origin: %s", origin)
			http.Error(w, "Forbidden origin", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-Clinic-ID")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			log.Printf("✅ OPTIONS request handled")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Убрали verbose логирование
		next(w, r)
	}
}

// enableCORSMiddleware - для http.Handler (используется с middleware chain)
func enableCORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("🌐 CORS Middleware: %s %s from origin: %s", r.Method, r.URL.Path, r.Header.Get("Origin"))

		origin := r.Header.Get("Origin")

		// Получаем разрешённые origins из переменной окружения
		allowedOriginsEnv := os.Getenv("ALLOWED_ORIGINS")
		allowedOrigins := map[string]bool{
			"http://localhost:3000": true, // Основной сайт (development)
			"http://localhost:4100": true, // ЗооБаза Frontend (development)
			"http://localhost:6100": true, // Кабинет владельца (development)
			"http://localhost:6200": true, // Кабинет волонтёра (development)
			"http://localhost:6300": true, // Кабинет клиники (development)
		}

		// Добавляем origins из .env
		if allowedOriginsEnv != "" {
			for _, o := range strings.Split(allowedOriginsEnv, ",") {
				allowedOrigins[strings.TrimSpace(o)] = true
			}
		}

		// Проверяем, разрешён ли origin
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			log.Printf("✅ Origin allowed: %s", origin)
		} else if origin == "" {
			// Если origin не указан, используем дефолтный
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4100")
			// Убрали verbose логирование
		} else {
			// Origin не разрешён
			log.Printf("❌ Blocked request from unauthorized origin: %s", origin)
			http.Error(w, "Forbidden origin", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-Clinic-ID")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			log.Printf("✅ OPTIONS request handled")
			w.WriteHeader(http.StatusOK)
			return
		}

		log.Printf("➡️ Passing to next middleware/handler: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
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
	// Публичные endpoints (без аутентификации)
	http.HandleFunc("/api/health", enableCORS(handleHealth))
	http.HandleFunc("/api/species", enableCORS(handlers.SpeciesHandler))
	http.HandleFunc("/api/species/", enableCORS(handlers.SpeciesDetailHandler))
	http.HandleFunc("/api/breeds", enableCORS(handlers.BreedsHandler))
	http.HandleFunc("/api/breeds/", enableCORS(handlers.BreedDetailHandler))
	http.HandleFunc("/api/breeds/species/", enableCORS(handlers.BreedsBySpeciesHandler))
	http.HandleFunc("/api/cards", enableCORS(handlers.CardsHandler))
	http.HandleFunc("/api/cards/breed/", enableCORS(handlers.CardsByBreedHandler))
	http.HandleFunc("/api/cards/", enableCORS(handlers.CardDetailHandler))
	http.HandleFunc("/api/catalog", enableCORS(handlers.CatalogHandler)) // Публичный каталог для главного сайта

	// Защищённые endpoints (требуют аутентификации)
	// Pets routes - реальные питомцы пользователей
	http.HandleFunc("/api/pets/search", enableCORS(handlers.SearchPetsHandler))                                                // Поиск питомцев (без auth для клиник)
	http.Handle("/api/pets/user/", middleware.OptionalAuthMiddleware(http.HandlerFunc(enableCORS(handlers.PetDetailHandler)))) // Получение питомцев пользователя (опциональная auth)

	// ВАЖНО: Используем общий middleware из pkg, который работает с Auth Service (7100)
	// Сначала CORS, потом Auth
	// POST/PUT/DELETE требуют авторизацию, GET - опциональная
	http.Handle("/api/pets", enableCORSMiddleware(middleware.AuthMiddleware(http.HandlerFunc(handlers.PetsHandler))))               // POST - создание (требует auth)
	http.Handle("/api/pets/", enableCORSMiddleware(middleware.OptionalAuthMiddleware(http.HandlerFunc(handlers.PetDetailHandler)))) // GET - просмотр (опциональная auth)

	// PetID Events routes - история событий питомцев
	http.HandleFunc("/api/pet-events", enableCORS(handlers.CreatePetEventSimpleHandler)) // Создание события (для клиник)
	http.HandleFunc("/api/petid/", enableCORS(handlePetIDRoutes))

	// Static files - раздача загруженных файлов с CORS
	// Путь от petbase/backend/ к корневой папке uploads/
	fileServer := http.FileServer(http.Dir("../../uploads"))
	http.Handle("/uploads/", enableCORSMiddleware(http.StripPrefix("/uploads/", fileServer)))
	log.Println("📁 Static files: /uploads/ -> ../../uploads (from petbase/backend/)")

	// Root route закомментирован - не нужен, мешает статическим файлам
	// http.HandleFunc("/", enableCORS(handleRoot))

	port := ":8100"
	fmt.Printf("🐾 ЗооБаза API starting on port %s\n", port)
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
	return "localhost:3000, localhost:4100"
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	// Только для корневого пути
	// НЕ обрабатываем /uploads/ - это статические файлы
	if r.URL.Path != "/" {
		// Если это не корень и не начинается с /uploads/, возвращаем 404
		if !strings.HasPrefix(r.URL.Path, "/uploads/") {
			http.NotFound(w, r)
			return
		}
		// Для /uploads/ пропускаем дальше (обработает FileServer)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "ЗооБаза API", "version": "1.0.0"}`)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok", "service": "petbase"}`)
}

// handlePetIDRoutes обрабатывает роуты для PetID событий
func handlePetIDRoutes(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// /api/petid/:id/events - история событий
	if strings.Contains(path, "/events") {
		switch r.Method {
		case http.MethodGet:
			handlers.GetPetEventsHandler(w, r)
		case http.MethodPost:
			// Используем AuthMiddleware для POST запросов
			middleware.AuthMiddleware(http.HandlerFunc(handlers.CreatePetEventHandler)).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// /api/petid/:id/medical - медицинская история
	if strings.Contains(path, "/medical") {
		switch r.Method {
		case http.MethodGet:
			handlers.GetPetMedicalHistoryHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	http.NotFound(w, r)
}

func createTables() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS species (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			name_en TEXT NOT NULL,
			description TEXT,
			icon TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS breeds (
			id SERIAL PRIMARY KEY,
			species_id INTEGER NOT NULL,
			name TEXT NOT NULL,
			name_en TEXT,
			description TEXT,
			origin TEXT,
			size TEXT,
			weight_min REAL,
			weight_max REAL,
			lifespan_min INTEGER,
			lifespan_max INTEGER,
			temperament TEXT,
			care_level TEXT,
			photo TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (species_id) REFERENCES species(id)
		)`,
		`CREATE TABLE IF NOT EXISTS pet_cards (
			id SERIAL PRIMARY KEY,
			breed_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT,
			characteristics TEXT,
			care_tips TEXT,
			health_info TEXT,
			nutrition TEXT,
			photos TEXT,
			is_published BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (breed_id) REFERENCES breeds(id)
		)`,
	}

	for _, query := range queries {
		if _, err := database.DB.Exec(query); err != nil {
			return err
		}
	}

	log.Println("✅ ЗооБаза tables created successfully")
	return nil
}
