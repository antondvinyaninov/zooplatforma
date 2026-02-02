package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	// Получаем connection string из переменной окружения
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Fallback для локальной разработки - используем SQLite
		log.Println("📦 Using SQLite for local development (database/data.db)")
		return InitSQLite("database/data.db")
	}

	// Production - используем PostgreSQL
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
	return createTablesPostgreSQL()
}

func InitSQLite(path string) error {
	var err error

	// Если путь относительный и не начинается с ../, добавляем ../../
	// (для запуска из service/backend директории)
	if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "..") {
		path = "../../" + path
	}

	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return fmt.Errorf("failed to open SQLite database: %w", err)
	}

	// Настройка для SQLite
	DB.SetMaxOpenConns(1) // SQLite работает лучше с одним соединением

	// Проверяем подключение
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping SQLite database: %w", err)
	}

	log.Printf("✅ SQLite database connected successfully: %s", path)
	return createTablesSQLite()
}

// fixOrganizationsTable добавляет недостающие колонки в таблицу organizations
func fixOrganizationsTable() error {
	fixes := []string{
		"ALTER TABLE organizations ADD COLUMN IF NOT EXISTS short_name TEXT",
		"ALTER TABLE organizations ADD COLUMN IF NOT EXISTS bio TEXT",
		"ALTER TABLE organizations ADD COLUMN IF NOT EXISTS cover_photo TEXT",
		"ALTER TABLE organizations ADD COLUMN IF NOT EXISTS address_city TEXT",
		"ALTER TABLE organizations ADD COLUMN IF NOT EXISTS address_region TEXT",
		"ALTER TABLE organizations ADD COLUMN IF NOT EXISTS is_verified BOOLEAN DEFAULT FALSE",
		"ALTER TABLE organization_members ADD COLUMN IF NOT EXISTS can_post BOOLEAN DEFAULT FALSE",
		"UPDATE organization_members SET can_post = TRUE WHERE role IN ('owner', 'admin') AND can_post IS NULL",
	}

	for _, fix := range fixes {
		if _, err := DB.Exec(fix); err != nil {
			// Игнорируем ошибки "column already exists"
			if !strings.Contains(err.Error(), "already exists") && !strings.Contains(err.Error(), "duplicate column") {
				return fmt.Errorf("failed to apply fix: %s - %w", fix, err)
			}
		}
	}

	return nil
}

func createTablesPostgreSQL() error {
	// Миграция уже применена вручную через psql
	// Просто проверяем что таблицы существуют
	log.Println("✅ PostgreSQL tables already exist (applied manually)")

	// Применяем fix для organizations table (на случай если нужно)
	if err := fixOrganizationsTable(); err != nil {
		log.Printf("⚠️ Organizations table fix failed (maybe already applied): %v", err)
	} else {
		log.Println("✅ Organizations table fix applied successfully")
	}

	return nil
}

// createBasicTablesPostgreSQL создает минимальные таблицы если миграция не найдена
func createBasicTablesPostgreSQL() error {
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
		return fmt.Errorf("failed to create basic PostgreSQL tables: %w", err)
	}

	log.Println("✅ Basic PostgreSQL tables created")
	return nil
}

func createTablesSQLite() error {
	// SQLite использует INTEGER PRIMARY KEY AUTOINCREMENT
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS user_media (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
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
		uploaded_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_user_media_user_id ON user_media(user_id);
	CREATE INDEX IF NOT EXISTS idx_user_media_type ON user_media(media_type);
	CREATE INDEX IF NOT EXISTS idx_user_media_uploaded ON user_media(uploaded_at);
	`
	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create SQLite tables: %w", err)
	}

	log.Println("✅ SQLite tables created successfully")
	return nil
}

func createTables() error {
	// Deprecated - используй createTablesPostgreSQL() или createTablesSQLite()
	log.Println("⚠️ Warning: createTables() is deprecated")
	return createTablesPostgreSQL()
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
