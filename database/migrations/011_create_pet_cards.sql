-- Миграция 011: Расширение таблицы pets для полноценных карточек
-- Дата: 2025-12-29
-- Описание: Добавление полей для полной информации о питомце
-- ВАЖНО: Таблица pets уже существует, добавляем только новые поля

-- ============================================
-- 1. Добавление новых полей в существующую таблицу pets
-- ============================================

-- Основная информация
ALTER TABLE pets ADD COLUMN breed TEXT;
ALTER TABLE pets ADD COLUMN gender TEXT; -- 'male', 'female', 'unknown'
ALTER TABLE pets ADD COLUMN birth_date DATE;

-- Внешний вид
ALTER TABLE pets ADD COLUMN color TEXT;
ALTER TABLE pets ADD COLUMN size TEXT; -- 'small', 'medium', 'large'
ALTER TABLE pets ADD COLUMN weight REAL; -- в кг

-- Идентификация
ALTER TABLE pets ADD COLUMN chip_number TEXT;
ALTER TABLE pets ADD COLUMN passport_number TEXT;

-- Медицинская информация
ALTER TABLE pets ADD COLUMN is_sterilized BOOLEAN DEFAULT FALSE;
ALTER TABLE pets ADD COLUMN is_vaccinated BOOLEAN DEFAULT FALSE;
ALTER TABLE pets ADD COLUMN health_notes TEXT;

-- Характер и особенности
ALTER TABLE pets ADD COLUMN character_traits TEXT; -- JSON массив
ALTER TABLE pets ADD COLUMN special_needs TEXT;

-- Статус
ALTER TABLE pets ADD COLUMN status TEXT DEFAULT 'home'; -- 'home', 'looking_for_home', 'lost', 'found', 'deceased'
ALTER TABLE pets ADD COLUMN status_updated_at DATETIME;

-- Медиа
ALTER TABLE pets ADD COLUMN photos TEXT; -- JSON массив дополнительных фото

-- История
ALTER TABLE pets ADD COLUMN story TEXT;

-- Метаданные
ALTER TABLE pets ADD COLUMN updated_at DATETIME DEFAULT CURRENT_TIMESTAMP;

-- ============================================
-- 2. Индексы для новых полей
-- ============================================

CREATE INDEX IF NOT EXISTS idx_pets_status ON pets(status);
CREATE INDEX IF NOT EXISTS idx_pets_chip_number ON pets(chip_number);
CREATE INDEX IF NOT EXISTS idx_pets_breed ON pets(breed);

-- ============================================
-- 3. Таблица для связи постов и питомцев
-- ============================================

CREATE TABLE IF NOT EXISTS post_pets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    pet_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE,
    UNIQUE(post_id, pet_id)
);

CREATE INDEX IF NOT EXISTS idx_post_pets_post_id ON post_pets(post_id);
CREATE INDEX IF NOT EXISTS idx_post_pets_pet_id ON post_pets(pet_id);

-- ============================================
-- 4. Триггер для обновления updated_at
-- ============================================

CREATE TRIGGER IF NOT EXISTS update_pets_timestamp 
AFTER UPDATE ON pets
FOR EACH ROW
BEGIN
    UPDATE pets SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

-- ============================================
-- 5. Проверка миграции
-- ============================================

-- Проверяем структуру таблицы
PRAGMA table_info(pets);

-- Проверяем индексы
SELECT name FROM sqlite_master WHERE type='index' AND tbl_name='pets';
