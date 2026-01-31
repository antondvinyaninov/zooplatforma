package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	return InitDBWithPath("../../database/data.db")
}

func InitDBWithPath(path string) error {
	var err error

	// Добавляем параметры для лучшей конкурентности
	// WAL mode - позволяет читать во время записи
	// busy_timeout - ждет до 5 секунд если БД заблокирована
	// journal_mode=WAL - Write-Ahead Logging для лучшей производительности
	// cache=shared - разделяемый кэш между соединениями
	dsn := path + "?_busy_timeout=5000&_journal_mode=WAL&cache=shared"

	DB, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}

	// Настройка пула соединений
	DB.SetMaxOpenConns(25)   // Максимум 25 одновременных соединений
	DB.SetMaxIdleConns(5)    // Держать 5 idle соединений
	DB.SetConnMaxLifetime(0) // Соединения живут бесконечно

	if err = DB.Ping(); err != nil {
		return err
	}

	// Включаем WAL mode явно
	_, err = DB.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		log.Printf("⚠️ Warning: Could not enable WAL mode: %v", err)
	}

	// Устанавливаем busy timeout
	_, err = DB.Exec("PRAGMA busy_timeout=5000;")
	if err != nil {
		log.Printf("⚠️ Warning: Could not set busy timeout: %v", err)
	}

	log.Println("✅ Database connected successfully (WAL mode enabled)")
	return createTables()
}

func createTables() error {
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
		uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_user_media_user_id ON user_media(user_id);
	CREATE INDEX IF NOT EXISTS idx_user_media_type ON user_media(media_type);
	CREATE INDEX IF NOT EXISTS idx_user_media_uploaded ON user_media(uploaded_at);
	`
	_, err := DB.Exec(query)
	if err != nil {
		return err
	}

	log.Println("Tables created successfully")
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

// ExecWithRetry выполняет SQL запрос с повторными попытками при блокировке БД
func ExecWithRetry(query string, args ...interface{}) (sql.Result, error) {
	maxRetries := 3
	var result sql.Result
	var err error

	for i := 0; i < maxRetries; i++ {
		result, err = DB.Exec(query, args...)
		if err == nil {
			return result, nil
		}

		// Проверяем, является ли ошибка блокировкой БД
		if err.Error() == "database is locked" {
			if i < maxRetries-1 {
				// Ждем перед повторной попыткой (экспоненциальная задержка)
				waitTime := time.Duration(100*(i+1)) * time.Millisecond
				log.Printf("⚠️ Database locked, retrying in %v (attempt %d/%d)", waitTime, i+1, maxRetries)
				time.Sleep(waitTime)
				continue
			}
		}

		// Если это не блокировка или исчерпаны попытки, возвращаем ошибку
		return nil, err
	}

	return result, err
}

// QueryWithRetry выполняет SELECT запрос с повторными попытками при блокировке БД
func QueryWithRetry(query string, args ...interface{}) (*sql.Rows, error) {
	maxRetries := 3
	var rows *sql.Rows
	var err error

	for i := 0; i < maxRetries; i++ {
		rows, err = DB.Query(query, args...)
		if err == nil {
			return rows, nil
		}

		// Проверяем, является ли ошибка блокировкой БД
		if err.Error() == "database is locked" {
			if i < maxRetries-1 {
				waitTime := time.Duration(100*(i+1)) * time.Millisecond
				log.Printf("⚠️ Database locked, retrying in %v (attempt %d/%d)", waitTime, i+1, maxRetries)
				time.Sleep(waitTime)
				continue
			}
		}

		return nil, err
	}

	return rows, err
}

// QueryRowWithRetry выполняет SELECT запрос для одной строки с повторными попытками
func QueryRowWithRetry(query string, args ...interface{}) *sql.Row {
	// QueryRow не возвращает ошибку сразу, поэтому просто возвращаем результат
	// Ошибка будет обработана при Scan()
	return DB.QueryRow(query, args...)
}
