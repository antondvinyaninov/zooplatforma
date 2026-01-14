-- Миграция 027: Создание системы ролей пользователей
-- Дата: 2025-01-14
-- Описание: Таблица для управления ролями пользователей в системе

-- Таблица ролей пользователей
CREATE TABLE IF NOT EXISTS user_roles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    role TEXT NOT NULL CHECK(role IN ('user', 'volunteer', 'shelter_admin', 'clinic_admin', 'moderator', 'superadmin')),
    granted_by INTEGER, -- ID пользователя, который назначил роль
    granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME, -- Опционально: дата истечения роли
    is_active BOOLEAN DEFAULT 1,
    notes TEXT, -- Примечания о назначении роли
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (granted_by) REFERENCES users(id) ON DELETE SET NULL,
    
    -- Уникальность: один пользователь не может иметь дублирующиеся активные роли
    UNIQUE(user_id, role, is_active)
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role ON user_roles(role);
CREATE INDEX IF NOT EXISTS idx_user_roles_active ON user_roles(is_active);

-- Добавляем роль 'user' всем существующим пользователям
INSERT INTO user_roles (user_id, role, granted_at, is_active)
SELECT id, 'user', created_at, 1
FROM users
WHERE id NOT IN (SELECT user_id FROM user_roles WHERE role = 'user');

-- Комментарии к ролям:
-- user: обычный пользователь (владелец питомца)
-- volunteer: волонтёр/зоопомощник
-- shelter_admin: администратор приюта
-- clinic_admin: администратор ветклиники
-- moderator: модератор контента
-- superadmin: суперадминистратор платформы
