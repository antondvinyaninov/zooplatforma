package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	// Получаем connection string из переменной окружения
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Fallback для локальной разработки
		dbURL = "postgres://postgres_zp:7da0905cd3349f58f368@localhost:5432/bd_zp?sslmode=disable"
	}

	return InitDBWithURL(dbURL)
}

func InitDBWithPath(path string) error {
	// Для обратной совместимости - если передан путь, используем SQLite
	// В продакшене используется InitDB() или InitDBWithURL()
	log.Printf("⚠️ Warning: InitDBWithPath is deprecated. Use InitDB() or InitDBWithURL() instead")
	return InitDB()
}

func InitDBWithURL(dbURL string) error {
	var err error

	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Настройка пула соединений для PostgreSQL
	DB.SetMaxOpenConns(25)                 // Максимум 25 одновременных соединений
	DB.SetMaxIdleConns(5)                  // Держать 5 idle соединений
	DB.SetConnMaxLifetime(5 * time.Minute) // Переиспользовать соединения каждые 5 минут

	// Проверяем подключение
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ PostgreSQL database connected successfully")
	return createTables()
}

func createTables() error {
	// PostgreSQL использует SERIAL для auto-increment
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS user_media (
		id SERIAL PRIMARY KEY,
		user_id INTEGER NOT NULL,
		file_name TEXT NOT NULL,
		original_name TEXT NOT NULL,
		file_path TEXT NOT NULL,
		file_size INTEGER NOT NULL,
		mime_type TEXT NOT NULL,
		media_type TEXT NOT NULL,
		width INTEGER,
		height INTEGER,
		duration INTEGER,
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_user_media_user_id ON user_media(user_id);
	CREATE INDEX IF NOT EXISTS idx_user_media_type ON user_media(media_type);
	CREATE INDEX IF NOT EXISTS idx_user_media_uploaded ON user_media(uploaded_at);
	`
	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("✅ Tables created successfully")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

// ExecWithRetry выполняет SQL запрос
// PostgreSQL не требует retry логики как SQLite, но оставляем для совместимости
func ExecWithRetry(query string, args ...interface{}) (sql.Result, error) {
	result, err := DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("exec error: %w", err)
	}
	return result, nil
}

// QueryWithRetry выполняет SELECT запрос
// PostgreSQL не требует retry логики как SQLite, но оставляем для совместимости
func QueryWithRetry(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	return rows, nil
}

// QueryRowWithRetry выполняет SELECT запрос для одной строки
func QueryRowWithRetry(query string, args ...interface{}) *sql.Row {
	return DB.QueryRow(query, args...)
}
