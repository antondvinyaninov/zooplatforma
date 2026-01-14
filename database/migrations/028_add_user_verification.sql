-- Миграция 028: Добавление верификации пользователей
-- Дата: 2025-01-14
-- Описание: Поле verified для отметки проверенных пользователей

-- Добавляем поле verified в таблицу users
ALTER TABLE users ADD COLUMN verified BOOLEAN DEFAULT 0;

-- Добавляем поле verified_at для даты верификации
ALTER TABLE users ADD COLUMN verified_at DATETIME;

-- Добавляем поле verified_by для ID модератора, который верифицировал
ALTER TABLE users ADD COLUMN verified_by INTEGER;

-- Индекс для быстрого поиска верифицированных пользователей
CREATE INDEX IF NOT EXISTS idx_users_verified ON users(verified);

-- Комментарий:
-- verified = 1 означает, что пользователь прошёл проверку модератором
-- Проверенные пользователи получают бейдж "✓ Проверенный" на профиле
