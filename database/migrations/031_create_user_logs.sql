-- Таблица для логирования действий пользователей
CREATE TABLE IF NOT EXISTS user_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    action_type TEXT NOT NULL,
    action_details TEXT,
    ip_address TEXT,
    user_agent TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_user_logs_user_id ON user_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_user_logs_created_at ON user_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_user_logs_action_type ON user_logs(action_type);

-- Типы действий:
-- 'register' - регистрация
-- 'login' - вход в систему
-- 'logout' - выход из системы
-- 'profile_update' - обновление профиля
-- 'post_create' - создание поста
-- 'post_delete' - удаление поста
-- 'comment_create' - создание комментария
-- 'like' - лайк
-- 'friend_request' - запрос в друзья
-- 'message_send' - отправка сообщения
-- 'pet_create' - создание питомца
-- 'pet_update' - обновление питомца
-- 'organization_create' - создание организации
-- 'organization_join' - вступление в организацию
-- 'error' - ошибка в системе
