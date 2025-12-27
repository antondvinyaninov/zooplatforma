package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	return InitDBWithPath("../../database/data.db")
}

func InitDBWithPath(path string) error {
	var err error
	DB, err = sql.Open("sqlite3", path)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Database connected successfully")
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
