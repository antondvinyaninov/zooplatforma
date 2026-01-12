-- Миграция 023: Добавление полей для каталога питомцев
-- Дата: 2025-01-12
-- Описание: Добавление полей для публичного каталога (ЭТАП 0 MVP)
-- Связано с: docs/MVP.md - ЭТАП 0 - ПОДГОТОВКА PETID

-- ============================================
-- 1. Добавление полей для каталога
-- ============================================

-- Местоположение для фильтрации
ALTER TABLE pets ADD COLUMN city TEXT DEFAULT '';
ALTER TABLE pets ADD COLUMN region TEXT DEFAULT '';

-- Срочность (SOS-животные показываются первыми)
ALTER TABLE pets ADD COLUMN urgent BOOLEAN DEFAULT FALSE;

-- Контакты для связи (для каталога)
ALTER TABLE pets ADD COLUMN contact_phone TEXT DEFAULT '';
ALTER TABLE pets ADD COLUMN contact_name TEXT DEFAULT '';

-- ============================================
-- 2. Индексы для быстрого поиска
-- ============================================

CREATE INDEX IF NOT EXISTS idx_pets_city ON pets(city);
CREATE INDEX IF NOT EXISTS idx_pets_region ON pets(region);
CREATE INDEX IF NOT EXISTS idx_pets_urgent ON pets(urgent);

-- ============================================
-- 3. Комментарии к полям
-- ============================================

-- city: Город для фильтрации в каталоге (например: "Ижевск", "Москва")
-- region: Регион для фильтрации (например: "Удмуртская Респ", "Москва")
-- urgent: Срочный случай (SOS) - показывается первым в каталоге
-- contact_phone: Телефон для связи (отображается в каталоге)
-- contact_name: Контактное лицо (отображается в каталоге)

-- ============================================
-- 4. Проверка миграции
-- ============================================

-- Проверяем, что поля добавлены
SELECT COUNT(*) as total_columns FROM pragma_table_info('pets');

-- Проверяем индексы
SELECT name FROM sqlite_master 
WHERE type='index' 
AND tbl_name='pets' 
AND name IN ('idx_pets_city', 'idx_pets_region', 'idx_pets_urgent');
