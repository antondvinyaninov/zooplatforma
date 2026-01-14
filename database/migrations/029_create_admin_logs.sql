-- Миграция: Создание таблицы логов администраторов
-- Дата: 2025-01-14

CREATE TABLE IF NOT EXISTS admin_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    admin_id INTEGER NOT NULL,
    admin_email TEXT NOT NULL,
    action_type TEXT NOT NULL, -- verify_user, unverify_user, grant_role, revoke_role, delete_post, delete_user, etc.
    target_type TEXT NOT NULL, -- user, post, role, organization, etc.
    target_id INTEGER NOT NULL,
    target_name TEXT,
    details TEXT, -- JSON с дополнительными деталями
    ip_address TEXT,
    user_agent TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (admin_id) REFERENCES users(id)
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_admin_logs_admin_id ON admin_logs(admin_id);
CREATE INDEX IF NOT EXISTS idx_admin_logs_action_type ON admin_logs(action_type);
CREATE INDEX IF NOT EXISTS idx_admin_logs_target_type ON admin_logs(target_type);
CREATE INDEX IF NOT EXISTS idx_admin_logs_target_id ON admin_logs(target_id);
CREATE INDEX IF NOT EXISTS idx_admin_logs_created_at ON admin_logs(created_at);
