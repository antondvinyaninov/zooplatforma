package main

import (
	"admin/handlers"
	"admin/middleware"
	"database"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Для admin frontend и других admin сервисов
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"http://localhost:4000", // Admin панель
			"http://localhost:4100", // PetBase (ЗооБаза)
			"https://sadmin.zooplatforma.ru",
		}

		for _, allowed := range allowedOrigins {
			if origin == allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	// Load .env file (try current directory first, then parent)
	if err := godotenv.Load(); err != nil {
		// Try loading from current directory
		if err := godotenv.Load(".env"); err != nil {
			log.Println("Warning: .env file not found, using default values")
		}
	}

	// Debug: проверяем JWT_SECRET
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("⚠️  JWT_SECRET not set, using default")
	} else {
		log.Printf("✅ JWT_SECRET loaded: %s...\n", secret[:10])
	}

	// Проверяем что секрет не default
	if secret == "default-secret-key" || secret == "" {
		log.Fatal("❌ JWT_SECRET must be set in .env file!")
	}

	// Initialize database (используем общую БД)
	if err := database.InitDBWithPath("../../database/data.db"); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Create admin tables
	if err := createAdminTables(); err != nil {
		log.Fatal("Failed to create admin tables:", err)
	}

	// Public routes
	http.HandleFunc("/", enableCORS(handleRoot))
	http.HandleFunc("/api/admin/health", enableCORS(handleHealth))

	// Auth routes
	http.HandleFunc("/api/admin/auth/login", enableCORS(handlers.AdminLoginHandler))
	http.HandleFunc("/api/admin/auth/logout", enableCORS(handlers.AdminLogoutHandler))
	http.HandleFunc("/api/admin/auth/me", enableCORS(handlers.AdminMeHandler))

	// Protected admin routes
	http.HandleFunc("/api/admin/users", enableCORS(middleware.SuperAdminMiddleware(handlers.AdminUsersHandler)))
	http.HandleFunc("/api/admin/users/", enableCORS(middleware.SuperAdminMiddleware(handlers.AdminUserHandler)))
	http.HandleFunc("/api/admin/posts", enableCORS(middleware.SuperAdminMiddleware(handlers.AdminPostsHandler)))
	http.HandleFunc("/api/admin/posts/", enableCORS(middleware.SuperAdminMiddleware(handlers.AdminPostHandler)))
	http.HandleFunc("/api/admin/stats/overview", enableCORS(middleware.SuperAdminMiddleware(handlers.AdminStatsOverviewHandler)))
	http.HandleFunc("/api/admin/logs", enableCORS(middleware.SuperAdminMiddleware(handlers.AdminLogsHandler)))

	port := ":9000"
	fmt.Printf("🔐 Admin Panel API starting on port %s\n", port)
	fmt.Println("📊 Dashboard: http://localhost:4000")
	log.Fatal(http.ListenAndServe(port, nil))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message": "ЗооПлатформа Admin API", "version": "0.0.1"}`)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status": "ok", "service": "admin-api"}`)
}

func createAdminTables() error {
	// Создаем таблицу администраторов
	_, err := database.DB.Exec(`
		CREATE TABLE IF NOT EXISTS admins (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL REFERENCES users(id),
			role TEXT NOT NULL DEFAULT 'moderator',
			permissions TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			created_by INTEGER REFERENCES admins(id)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create admins table: %w", err)
	}

	// Создаем таблицу логов администраторов
	_, err = database.DB.Exec(`
		CREATE TABLE IF NOT EXISTS admin_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			admin_id INTEGER NOT NULL REFERENCES admins(id),
			action TEXT NOT NULL,
			target_type TEXT,
			target_id INTEGER,
			details TEXT,
			ip_address TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create admin_logs table: %w", err)
	}

	// Создаем общую таблицу системных логов
	_, err = database.DB.Exec(`
		CREATE TABLE IF NOT EXISTS system_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			level TEXT NOT NULL DEFAULT 'info',
			category TEXT NOT NULL,
			action TEXT NOT NULL,
			user_id INTEGER REFERENCES users(id),
			target_type TEXT,
			target_id INTEGER,
			message TEXT,
			details TEXT,
			ip_address TEXT,
			user_agent TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create system_logs table: %w", err)
	}

	// Создаем индексы
	database.DB.Exec(`CREATE INDEX IF NOT EXISTS idx_admin_logs_admin_id ON admin_logs(admin_id)`)
	database.DB.Exec(`CREATE INDEX IF NOT EXISTS idx_admin_logs_created_at ON admin_logs(created_at)`)
	database.DB.Exec(`CREATE INDEX IF NOT EXISTS idx_admin_logs_action ON admin_logs(action)`)
	database.DB.Exec(`CREATE INDEX IF NOT EXISTS idx_system_logs_level ON system_logs(level)`)
	database.DB.Exec(`CREATE INDEX IF NOT EXISTS idx_system_logs_category ON system_logs(category)`)
	database.DB.Exec(`CREATE INDEX IF NOT EXISTS idx_system_logs_user_id ON system_logs(user_id)`)
	database.DB.Exec(`CREATE INDEX IF NOT EXISTS idx_system_logs_created_at ON system_logs(created_at)`)

	fmt.Println("✅ Admin tables created successfully")
	return nil
}
