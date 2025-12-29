-- Миграция 010: Добавление системы аналитики активности пользователей
-- Дата: 2025-12-28
-- Описание: Добавляет таблицы для отслеживания сессий, действий и статистики пользователей
-- Необходимо для: NPS опросов, retention анализа, engagement метрик

-- ============================================
-- 1. Обновление таблицы users
-- ============================================

-- Добавляем поля для базовой аналитики
ALTER TABLE users ADD COLUMN registered_at DATETIME DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE users ADD COLUMN last_login_at DATETIME;
ALTER TABLE users ADD COLUMN login_count INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN total_sessions INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN analytics_consent BOOLEAN DEFAULT TRUE;

-- Заполняем registered_at из created_at для существующих пользователей
UPDATE users SET registered_at = created_at WHERE registered_at IS NULL;

-- ============================================
-- 2. Таблица user_sessions - сессии пользователей
-- ============================================

CREATE TABLE IF NOT EXISTS user_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    ended_at DATETIME,
    duration_seconds INTEGER, -- вычисляется при завершении сессии
    pages_viewed INTEGER DEFAULT 0,
    actions_count INTEGER DEFAULT 0,
    device_type TEXT, -- 'desktop', 'mobile', 'tablet'
    browser TEXT,
    ip_address TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_started_at ON user_sessions(started_at);
CREATE INDEX idx_user_sessions_ended_at ON user_sessions(ended_at);

-- ============================================
-- 3. Таблица user_activity_log - лог действий
-- ============================================

CREATE TABLE IF NOT EXISTS user_activity_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_id INTEGER,
    action_type TEXT NOT NULL, -- 'post_create', 'comment_add', 'like', 'profile_view', etc.
    entity_type TEXT, -- 'post', 'user', 'pet', 'comment'
    entity_id INTEGER,
    metadata TEXT, -- JSON с дополнительными данными
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES user_sessions(id) ON DELETE SET NULL
);

CREATE INDEX idx_activity_user_id ON user_activity_log(user_id);
CREATE INDEX idx_activity_session_id ON user_activity_log(session_id);
CREATE INDEX idx_activity_action_type ON user_activity_log(action_type);
CREATE INDEX idx_activity_created_at ON user_activity_log(created_at);
CREATE INDEX idx_activity_entity ON user_activity_log(entity_type, entity_id);

-- ============================================
-- 4. Таблица user_stats - агрегированная статистика
-- ============================================

CREATE TABLE IF NOT EXISTS user_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    date DATE NOT NULL,
    sessions_count INTEGER DEFAULT 0,
    total_time_seconds INTEGER DEFAULT 0,
    posts_created INTEGER DEFAULT 0,
    comments_added INTEGER DEFAULT 0,
    likes_given INTEGER DEFAULT 0,
    profile_views INTEGER DEFAULT 0,
    messages_sent INTEGER DEFAULT 0,
    pets_added INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, date)
);

CREATE INDEX idx_user_stats_user_id ON user_stats(user_id);
CREATE INDEX idx_user_stats_date ON user_stats(date);

-- ============================================
-- Комментарии к типам действий (action_type)
-- ============================================

-- Контент:
-- post_create, post_edit, post_delete, post_view
-- comment_add, comment_delete
-- like_add, like_remove

-- Профиль:
-- profile_view, profile_edit
-- avatar_upload, cover_upload

-- Социальные связи:
-- friend_request_send, friend_request_accept, friend_request_reject
-- subscription_add, subscription_remove

-- Питомцы:
-- pet_add, pet_edit, pet_delete, pet_view

-- Коммуникации:
-- message_send, message_read, chat_open

-- Поиск:
-- search_query, search_result_click

-- Авторизация:
-- login, logout, register

-- ============================================
-- Проверка миграции
-- ============================================

-- Проверяем, что таблицы созданы
SELECT 'user_sessions table created' as status 
WHERE EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name='user_sessions');

SELECT 'user_activity_log table created' as status 
WHERE EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name='user_activity_log');

SELECT 'user_stats table created' as status 
WHERE EXISTS (SELECT 1 FROM sqlite_master WHERE type='table' AND name='user_stats');

-- Проверяем, что поля добавлены в users
PRAGMA table_info(users);
