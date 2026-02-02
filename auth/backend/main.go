package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	// Загрузить .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found")
	}

	// Инициализировать JWT Secret
	InitJWTSecret()

	// Подключиться к БД Auth Service
	var err error
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		// Production: используем Volume
		if os.Getenv("ENVIRONMENT") == "production" {
			dbPath = "/app/auth-data/auth.db"
		} else {
			// Development: локальный файл
			dbPath = "./auth.db"
		}
	}

	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}
	defer db.Close()

	// Инициализировать БД
	if err := initDatabase(); err != nil {
		log.Fatal("❌ Failed to initialize database:", err)
	}

	// Создать роутер
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/api/health", healthHandler).Methods("GET")

	// Auth endpoints
	r.HandleFunc("/api/auth/register", registerHandler).Methods("POST")
	r.HandleFunc("/api/auth/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/auth/logout", logoutHandler).Methods("POST")
	r.HandleFunc("/api/auth/refresh", refreshHandler).Methods("POST")
	r.HandleFunc("/api/auth/me", getMeHandler).Methods("GET")
	r.HandleFunc("/api/auth/verify", verifyTokenHandler).Methods("GET")

	// Password management
	r.HandleFunc("/api/auth/forgot-password", forgotPasswordHandler).Methods("POST")
	r.HandleFunc("/api/auth/reset-password", resetPasswordHandler).Methods("POST")
	r.HandleFunc("/api/auth/change-password", changePasswordHandler).Methods("POST")

	// Email verification
	r.HandleFunc("/api/auth/send-verification", sendVerificationHandler).Methods("POST")
	r.HandleFunc("/api/auth/verify-email", verifyEmailHandler).Methods("POST")

	// User profile endpoints
	r.HandleFunc("/api/users/{id}", getUserProfileHandler).Methods("GET")
	r.HandleFunc("/api/users/{id}", updateUserProfileHandler).Methods("PUT")
	r.HandleFunc("/api/users/{id}/avatar", uploadAvatarHandler).Methods("POST")
	r.HandleFunc("/api/users/{id}/avatar", deleteAvatarHandler).Methods("DELETE")

	// CORS middleware
	handler := enableCORS(r)

	// Запустить сервер
	// Auth Service ВСЕГДА использует порт 7100, игнорируя PORT переменную
	// (которую может установить EasyPanel на 80)
	port := "7100"

	log.Printf("🚀 Auth Service started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

func initDatabase() error {
	// Таблица user_roles для ролей пользователей
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS user_roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			role TEXT NOT NULL,
			granted_by INTEGER,
			granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME,
			is_active BOOLEAN DEFAULT 1,
			notes TEXT,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(user_id, role)
		)
	`)
	if err != nil {
		return err
	}

	// Создать индексы
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_user_roles_role ON user_roles(role)`)
	if err != nil {
		return err
	}

	// Создать таблицу sessions
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT UNIQUE NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return err
	}

	// Создать таблицу refresh_tokens
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT UNIQUE NOT NULL,
			expires_at DATETIME NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return err
	}

	// Создать таблицу password_resets
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS password_resets (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT UNIQUE NOT NULL,
			expires_at DATETIME NOT NULL,
			used BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)
	if err != nil {
		return err
	}

	// Создать таблицу email_verifications
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS email_verifications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			token TEXT UNIQUE NOT NULL,
			expires_at DATETIME NOT NULL,
			verified BOOLEAN DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)
	`)

	log.Println("✅ Database initialized (using shared database/data.db)")
	return err
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := []string{
			"http://localhost:3000",                                  // Main Frontend (dev)
			"http://localhost:4000",                                  // Admin Frontend (dev)
			"http://localhost:4100",                                  // PetBase Frontend (dev)
			"http://localhost:5100",                                  // Shelter Frontend (dev)
			"http://localhost:6100",                                  // Owner Frontend (dev)
			"http://localhost:6200",                                  // Volunteer Frontend (dev)
			"http://localhost:6300",                                  // Clinic Frontend (dev)
			"https://my-projects-zooplatforma.crv1ic.easypanel.host", // Main Frontend (prod)
			"https://my-projects-admin.crv1ic.easypanel.host",        // Admin Frontend (prod)
			"https://my-projects-petbase.crv1ic.easypanel.host",      // PetBase Frontend (prod)
			"https://my-projects-shelter.crv1ic.easypanel.host",      // Shelter Frontend (prod)
			"https://my-projects-owner.crv1ic.easypanel.host",        // Owner Frontend (prod)
			"https://my-projects-volunteer.crv1ic.easypanel.host",    // Volunteer Frontend (prod)
			"https://my-projects-clinic.crv1ic.easypanel.host",       // Clinic Frontend (prod)
		}

		origin := r.Header.Get("Origin")
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"success":true,"service":"auth-backend","version":"1.0.0","status":"healthy"}`))
}
