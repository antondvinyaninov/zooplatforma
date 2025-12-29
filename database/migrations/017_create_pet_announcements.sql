-- Миграция 017: Создание таблицы объявлений о питомцах
-- Дата: 2024-12-29
-- Описание: Универсальная система карточек для разных типов объявлений

-- ============================================
-- 1. Таблица объявлений о питомцах
-- ============================================

CREATE TABLE IF NOT EXISTS pet_announcements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    
    -- Связь с питомцем из ЗооБазы
    pet_id INTEGER NOT NULL,
    
    -- Тип объявления
    type TEXT NOT NULL, -- 'looking_for_home', 'found', 'lost', 'fundraising'
    
    -- Заголовок и описание
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    
    -- Автор объявления (кто создал карточку)
    author_id INTEGER NOT NULL,
    
    -- Контактное лицо (может отличаться от автора)
    contact_person_id INTEGER, -- ID пользователя (куратор, нашедший, владелец, организатор)
    contact_person_name TEXT, -- Если контакт не зарегистрирован в системе
    contact_person_phone TEXT,
    
    -- Локация события
    location_city TEXT,
    location_address TEXT,
    location_coordinates TEXT, -- JSON: {lat, lng}
    
    -- Дата и время события
    event_date DATETIME, -- Когда потерян/найден
    event_time TEXT, -- Примерное время (например, "около 18:00")
    
    -- Детали для "Потерян"
    lost_last_seen_location TEXT, -- Где видели последний раз
    lost_distinctive_features TEXT, -- Особые приметы для поиска
    lost_reward_amount INTEGER, -- Вознаграждение в рублях
    
    -- Детали для "Найден"
    found_current_location TEXT, -- Где сейчас находится
    found_condition TEXT, -- Состояние животного при находке
    
    -- Детали для "Сбор средств"
    fundraising_goal_amount INTEGER, -- Целевая сумма в рублях
    fundraising_current_amount INTEGER DEFAULT 0, -- Собрано
    fundraising_purpose TEXT, -- На что собираем (лечение, операция и т.д.)
    fundraising_deadline DATE, -- Срок сбора
    fundraising_bank_details TEXT, -- Реквизиты для перевода
    
    -- Статус объявления
    status TEXT DEFAULT 'active', -- 'active', 'closed', 'archived'
    status_reason TEXT, -- Причина закрытия (найден, пристроен и т.д.)
    
    -- Видимость
    is_published BOOLEAN DEFAULT TRUE,
    views_count INTEGER DEFAULT 0,
    
    -- Метаданные
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    closed_at DATETIME,
    
    FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (contact_person_id) REFERENCES users(id) ON DELETE SET NULL
);

-- ============================================
-- 2. Индексы
-- ============================================

CREATE INDEX IF NOT EXISTS idx_pet_announcements_pet_id ON pet_announcements(pet_id);
CREATE INDEX IF NOT EXISTS idx_pet_announcements_author_id ON pet_announcements(author_id);
CREATE INDEX IF NOT EXISTS idx_pet_announcements_type ON pet_announcements(type);
CREATE INDEX IF NOT EXISTS idx_pet_announcements_status ON pet_announcements(status);
CREATE INDEX IF NOT EXISTS idx_pet_announcements_location_city ON pet_announcements(location_city);
CREATE INDEX IF NOT EXISTS idx_pet_announcements_created_at ON pet_announcements(created_at);

-- ============================================
-- 3. Таблица публикаций (обновлений) к объявлению
-- ============================================

CREATE TABLE IF NOT EXISTS announcement_posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    announcement_id INTEGER NOT NULL,
    author_id INTEGER NOT NULL,
    
    -- Тип публикации
    post_type TEXT NOT NULL, -- 'update', 'photo', 'observation', 'search', 'reward', 'donation'
    
    -- Содержимое
    content TEXT NOT NULL,
    media_urls TEXT, -- JSON массив URL фото/видео
    
    -- Для донатов
    donation_amount INTEGER, -- Сумма пожертвования
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (announcement_id) REFERENCES pet_announcements(id) ON DELETE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_announcement_posts_announcement_id ON announcement_posts(announcement_id);
CREATE INDEX IF NOT EXISTS idx_announcement_posts_post_type ON announcement_posts(post_type);
CREATE INDEX IF NOT EXISTS idx_announcement_posts_created_at ON announcement_posts(created_at);

-- ============================================
-- 4. Таблица пожертвований для сборов
-- ============================================

CREATE TABLE IF NOT EXISTS announcement_donations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    announcement_id INTEGER NOT NULL,
    donor_id INTEGER, -- NULL если анонимный донор
    donor_name TEXT, -- Имя для отображения (может быть "Аноним")
    
    amount INTEGER NOT NULL, -- Сумма в рублях
    message TEXT, -- Сообщение от донора
    
    is_anonymous BOOLEAN DEFAULT FALSE,
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (announcement_id) REFERENCES pet_announcements(id) ON DELETE CASCADE,
    FOREIGN KEY (donor_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_announcement_donations_announcement_id ON announcement_donations(announcement_id);
CREATE INDEX IF NOT EXISTS idx_announcement_donations_created_at ON announcement_donations(created_at);

-- ============================================
-- 5. Триггеры
-- ============================================

-- Обновление updated_at при изменении объявления
CREATE TRIGGER IF NOT EXISTS update_pet_announcements_timestamp 
AFTER UPDATE ON pet_announcements
FOR EACH ROW
BEGIN
    UPDATE pet_announcements SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

-- Обновление суммы сбора при добавлении пожертвования
CREATE TRIGGER IF NOT EXISTS update_fundraising_amount
AFTER INSERT ON announcement_donations
FOR EACH ROW
WHEN (SELECT type FROM pet_announcements WHERE id = NEW.announcement_id) = 'fundraising'
BEGIN
    UPDATE pet_announcements 
    SET fundraising_current_amount = fundraising_current_amount + NEW.amount
    WHERE id = NEW.announcement_id;
END;

-- Автоматическое закрытие сбора при достижении цели
CREATE TRIGGER IF NOT EXISTS close_fundraising_on_goal
AFTER UPDATE OF fundraising_current_amount ON pet_announcements
FOR EACH ROW
WHEN NEW.type = 'fundraising' 
    AND NEW.fundraising_current_amount >= NEW.fundraising_goal_amount
    AND NEW.status = 'active'
BEGIN
    UPDATE pet_announcements 
    SET status = 'closed', 
        status_reason = 'Цель сбора достигнута',
        closed_at = CURRENT_TIMESTAMP
    WHERE id = NEW.id;
END;

-- ============================================
-- 6. Проверка миграции
-- ============================================

-- Проверяем структуру таблиц
PRAGMA table_info(pet_announcements);
PRAGMA table_info(announcement_posts);
PRAGMA table_info(announcement_donations);

-- Проверяем индексы
SELECT name FROM sqlite_master WHERE type='index' AND tbl_name='pet_announcements';
