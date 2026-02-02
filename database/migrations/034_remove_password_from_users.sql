-- Удаление колонки password из таблицы users
-- Auth Service теперь единственный владелец паролей (хранятся в auth/backend/auth.db)

-- Создать новую таблицу без password
CREATE TABLE users_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    bio TEXT DEFAULT '',
    phone TEXT DEFAULT '',
    location TEXT DEFAULT '',
    avatar TEXT DEFAULT '',
    cover_photo TEXT DEFAULT '',
    updated_at DATETIME,
    last_name TEXT DEFAULT '',
    profile_visibility TEXT DEFAULT 'public',
    show_phone TEXT DEFAULT 'nobody',
    show_email TEXT DEFAULT 'nobody',
    allow_messages TEXT DEFAULT 'everyone',
    show_online TEXT DEFAULT 'yes',
    verified BOOLEAN DEFAULT 0,
    verified_at DATETIME,
    verified_by INTEGER
);

-- Копировать данные (без password)
INSERT INTO users_new 
SELECT 
    id, name, email, created_at, bio, phone, location, avatar, cover_photo,
    updated_at, last_name, profile_visibility, show_phone, show_email,
    allow_messages, show_online, verified, verified_at, verified_by
FROM users;

-- Удалить старую таблицу
DROP TABLE users;

-- Переименовать новую
ALTER TABLE users_new RENAME TO users;

-- Пересоздать индексы
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_verified ON users(verified);
