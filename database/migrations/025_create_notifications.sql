-- Migration: Create notifications table
-- Description: Система уведомлений для комментариев, лайков и запросов в друзья
-- Date: 2025-01-12

-- Таблица уведомлений
CREATE TABLE IF NOT EXISTS notifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,                    -- Кому уведомление
    type TEXT NOT NULL,                          -- Тип: comment, like, friend_request
    actor_id INTEGER NOT NULL,                   -- Кто совершил действие
    entity_type TEXT,                            -- Тип сущности: post, comment, friendship
    entity_id INTEGER,                           -- ID сущности
    message TEXT NOT NULL,                       -- Текст уведомления
    is_read BOOLEAN DEFAULT 0,                   -- Прочитано ли
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_notifications_user_unread ON notifications(user_id, is_read);

-- Комментарий к таблице
-- Типы уведомлений:
-- - comment: новый комментарий к посту
-- - like: новый лайк поста
-- - friend_request: новый запрос в друзья
-- - friend_accepted: запрос в друзья принят
