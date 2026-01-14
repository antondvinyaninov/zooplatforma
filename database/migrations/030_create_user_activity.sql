-- Таблица для отслеживания активности пользователей
CREATE TABLE IF NOT EXISTS user_activity (
    user_id INTEGER PRIMARY KEY,
    last_seen DATETIME NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Индекс для быстрого поиска онлайн пользователей
CREATE INDEX IF NOT EXISTS idx_user_activity_last_seen ON user_activity(last_seen);

-- Пользователь считается онлайн, если был активен в последние 5 минут
