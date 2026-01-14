-- Миграция 030: Создание системы жалоб и модерации
-- Дата: 2025-01-14
-- Описание: Таблица для жалоб на контент и пользователей

-- ============================================
-- Таблица жалоб (reports)
-- ============================================

CREATE TABLE IF NOT EXISTS reports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    
    -- Кто пожаловался
    reporter_id INTEGER NOT NULL,
    
    -- На что жалоба
    target_type TEXT NOT NULL CHECK(target_type IN ('post', 'comment', 'user', 'organization', 'pet')),
    target_id INTEGER NOT NULL,
    
    -- Причина жалобы
    reason TEXT NOT NULL CHECK(reason IN (
        'spam',           -- Спам
        'harassment',     -- Оскорбления/домогательства
        'violence',       -- Насилие
        'hate_speech',    -- Разжигание ненависти
        'misinformation', -- Дезинформация
        'inappropriate',  -- Неприемлемый контент
        'copyright',      -- Нарушение авторских прав
        'animal_abuse',   -- Жестокое обращение с животными
        'fraud',          -- Мошенничество
        'other'           -- Другое
    )),
    
    -- Детали жалобы
    description TEXT,
    
    -- Статус модерации
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'reviewing', 'resolved', 'rejected')),
    
    -- Модератор и решение
    moderator_id INTEGER,
    moderator_action TEXT CHECK(moderator_action IN ('no_action', 'warning', 'content_removed', 'user_banned', 'user_suspended')),
    moderator_comment TEXT,
    reviewed_at DATETIME,
    
    -- Метаданные
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    -- Внешние ключи
    FOREIGN KEY (reporter_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (moderator_id) REFERENCES users(id) ON DELETE SET NULL
);

-- ============================================
-- Индексы для быстрого поиска
-- ============================================

-- Индекс для поиска жалоб по статусу (очередь модерации)
CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status, created_at DESC);

-- Индекс для поиска жалоб на конкретный объект
CREATE INDEX IF NOT EXISTS idx_reports_target ON reports(target_type, target_id);

-- Индекс для поиска жалоб от пользователя
CREATE INDEX IF NOT EXISTS idx_reports_reporter ON reports(reporter_id, created_at DESC);

-- Индекс для поиска жалоб модератора
CREATE INDEX IF NOT EXISTS idx_reports_moderator ON reports(moderator_id, reviewed_at DESC);

-- Индекс для поиска по причине
CREATE INDEX IF NOT EXISTS idx_reports_reason ON reports(reason, status);

-- ============================================
-- Таблица действий модератора (moderation_log)
-- ============================================

CREATE TABLE IF NOT EXISTS moderation_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    
    -- Модератор
    moderator_id INTEGER NOT NULL,
    moderator_email TEXT NOT NULL,
    
    -- Действие
    action_type TEXT NOT NULL CHECK(action_type IN (
        'report_reviewed',    -- Рассмотрена жалоба
        'content_removed',    -- Удален контент
        'user_warned',        -- Предупреждение пользователю
        'user_suspended',     -- Временная блокировка
        'user_banned',        -- Постоянная блокировка
        'user_unbanned',      -- Разблокировка
        'comment_hidden',     -- Скрыт комментарий
        'post_hidden'         -- Скрыт пост
    )),
    
    -- Цель действия
    target_type TEXT NOT NULL,
    target_id INTEGER NOT NULL,
    target_name TEXT,
    
    -- Связь с жалобой (если есть)
    report_id INTEGER,
    
    -- Детали
    reason TEXT,
    details TEXT,
    
    -- Метаданные
    ip_address TEXT,
    user_agent TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (moderator_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (report_id) REFERENCES reports(id) ON DELETE SET NULL
);

-- ============================================
-- Индексы для логов модерации
-- ============================================

CREATE INDEX IF NOT EXISTS idx_moderation_log_moderator ON moderation_log(moderator_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_moderation_log_action ON moderation_log(action_type, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_moderation_log_target ON moderation_log(target_type, target_id);
CREATE INDEX IF NOT EXISTS idx_moderation_log_report ON moderation_log(report_id);

-- ============================================
-- Таблица блокировок пользователей (user_bans)
-- ============================================

CREATE TABLE IF NOT EXISTS user_bans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    
    -- Заблокированный пользователь
    user_id INTEGER NOT NULL,
    
    -- Кто заблокировал
    banned_by INTEGER NOT NULL,
    
    -- Тип блокировки
    ban_type TEXT NOT NULL CHECK(ban_type IN ('warning', 'temporary', 'permanent')),
    
    -- Причина
    reason TEXT NOT NULL,
    details TEXT,
    
    -- Срок блокировки
    banned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME, -- NULL для permanent
    
    -- Статус
    is_active BOOLEAN DEFAULT 1,
    
    -- Связь с жалобой
    report_id INTEGER,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (banned_by) REFERENCES users(id) ON DELETE SET NULL,
    FOREIGN KEY (report_id) REFERENCES reports(id) ON DELETE SET NULL
);

-- ============================================
-- Индексы для блокировок
-- ============================================

CREATE INDEX IF NOT EXISTS idx_user_bans_user ON user_bans(user_id, is_active);
CREATE INDEX IF NOT EXISTS idx_user_bans_active ON user_bans(is_active, expires_at);
CREATE INDEX IF NOT EXISTS idx_user_bans_banned_by ON user_bans(banned_by, banned_at DESC);

-- ============================================
-- Комментарии
-- ============================================

-- Система жалоб позволяет:
-- 1. Пользователям жаловаться на неприемлемый контент
-- 2. Модераторам просматривать очередь жалоб
-- 3. Принимать решения и логировать действия
-- 4. Блокировать пользователей временно или навсегда
-- 5. Отслеживать статистику модерации

-- Типы контента для жалоб:
-- - post: жалоба на пост
-- - comment: жалоба на комментарий
-- - user: жалоба на пользователя
-- - organization: жалоба на организацию
-- - pet: жалоба на карточку питомца

-- Причины жалоб:
-- - spam: спам, реклама
-- - harassment: оскорбления, домогательства
-- - violence: насилие, жестокость
-- - hate_speech: разжигание ненависти
-- - misinformation: ложная информация
-- - inappropriate: неприемлемый контент
-- - copyright: нарушение авторских прав
-- - animal_abuse: жестокое обращение с животными
-- - fraud: мошенничество
-- - other: другое

-- Действия модератора:
-- - no_action: жалоба необоснованна
-- - warning: предупреждение пользователю
-- - content_removed: контент удален
-- - user_banned: пользователь заблокирован навсегда
-- - user_suspended: пользователь заблокирован временно
