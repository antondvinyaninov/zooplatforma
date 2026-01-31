-- Миграция: перенос авторизации в Auth Service
-- Auth Service теперь единственный владелец паролей и ролей

-- 1. Удалить колонку password из users (SQLite не поддерживает DROP COLUMN, нужно пересоздать таблицу)
-- Сначала создадим новую таблицу без password
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
    verified_by INTEGER,
    date_of_birth DATE,
    gender TEXT DEFAULT '',
    city TEXT DEFAULT '',
    country TEXT DEFAULT ''
);

-- Копируем данные (без password)
INSERT INTO users_new 
SELECT 
    id, name, email, created_at, bio, phone, location, avatar, cover_photo,
    updated_at, last_name, profile_visibility, show_phone, show_email,
    allow_messages, show_online, verified, verified_at, verified_by,
    NULL as date_of_birth, '' as gender, '' as city, '' as country
FROM users;

-- Удаляем старую таблицу
DROP TABLE users;

-- Переименовываем новую
ALTER TABLE users_new RENAME TO users;

-- Пересоздаем индексы
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 2. Удалить таблицу user_roles из общей базы (она теперь в Auth Service)
DROP TABLE IF EXISTS user_roles;

-- Готово! Теперь:
-- - Пароли хранятся только в Auth Service (auth/backend/auth.db или общая БД через Auth Service)
-- - Роли хранятся только в Auth Service
-- - Main Backend работает только с публичными данными пользователей
