-- Миграция 020: Добавление полей для каталога питомцев
-- Дата: 2024-12-30
-- Описание: Добавление полей для работы каталога и интеграции с организациями
-- ВАЖНО: Существующие данные сохраняются, добавляются только новые поля

-- ============================================
-- 1. Добавление новых полей для каталога
-- ============================================

-- Регион для фильтрации в каталоге
ALTER TABLE pets ADD COLUMN region TEXT DEFAULT '';

-- Срочность (SOS-животные показываются первыми в каталоге)
ALTER TABLE pets ADD COLUMN urgent BOOLEAN DEFAULT 0;

-- Контактное лицо (для связи по животному)
ALTER TABLE pets ADD COLUMN contact_name TEXT DEFAULT '';

-- Связь с организацией (приют, клиника, фонд)
ALTER TABLE pets ADD COLUMN organization_id INTEGER DEFAULT NULL;

-- ============================================
-- 2. Индексы для быстрого поиска в каталоге
-- ============================================

-- Индекс для фильтрации по региону
CREATE INDEX IF NOT EXISTS idx_pets_region ON pets(region);

-- Индекс для срочных животных
CREATE INDEX IF NOT EXISTS idx_pets_urgent ON pets(urgent);

-- Индекс для поиска по организации
CREATE INDEX IF NOT EXISTS idx_pets_organization_id ON pets(organization_id);

-- Составной индекс для каталога (статус + город + срочность)
CREATE INDEX IF NOT EXISTS idx_pets_catalog ON pets(status, card_location_city, urgent);

-- ============================================
-- 3. Комментарии к полям
-- ============================================

-- region: Регион РФ (например, "Удмуртская Республика", "Московская область")
--         Используется для фильтрации в каталоге по региону

-- urgent: Срочность (0 или 1)
--         1 = SOS-животное (нужна срочная помощь, показывается первым в каталоге)
--         0 = обычное животное

-- contact_name: Имя контактного лица для связи по животному
--               Дополняет card_contact_phone (телефон уже есть)

-- organization_id: ID организации из таблицы organizations
--                  NULL = частное лицо (владелец/зоопомощник)
--                  NOT NULL = животное принадлежит организации (приют/клиника)

-- ============================================
-- 4. Проверка миграции
-- ============================================

-- Проверяем, что новые поля добавлены
SELECT 
    name, 
    type, 
    dflt_value 
FROM pragma_table_info('pets') 
WHERE name IN ('region', 'urgent', 'contact_name', 'organization_id');

-- Проверяем индексы
SELECT name FROM sqlite_master 
WHERE type='index' 
AND tbl_name='pets' 
AND name LIKE 'idx_pets_%';

-- Подсчитываем животных (должно остаться столько же)
SELECT COUNT(*) as total_pets FROM pets;
