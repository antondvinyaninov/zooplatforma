-- Миграция 018: Добавление типа карточки питомца для постов
-- Дата: 2024-12-29
-- Описание: Добавляем поле card_type для определения типа отображения карточки в постах

-- ============================================
-- 1. Добавление поля card_type в таблицу pets
-- ============================================

-- Тип карточки определяет как она отображается в посте
-- 'default' - обычная карточка (просто информация)
-- 'looking_for_home' - ищет дом (кнопка "Взять домой")
-- 'found' - найден (кнопка "Я владелец")
-- 'lost' - потерян (кнопка "Я видел")
-- 'fundraising' - сбор средств (кнопка "Помочь", прогресс-бар)

ALTER TABLE pets ADD COLUMN card_type TEXT DEFAULT 'default';

-- Дополнительные поля для карточек
ALTER TABLE pets ADD COLUMN card_description TEXT; -- Описание для карточки (отдельно от story)
ALTER TABLE pets ADD COLUMN card_contact_phone TEXT; -- Контактный телефон
ALTER TABLE pets ADD COLUMN card_location_city TEXT; -- Город
ALTER TABLE pets ADD COLUMN card_location_address TEXT; -- Адрес

-- Поля для "Потерян"
ALTER TABLE pets ADD COLUMN card_lost_date DATE; -- Когда потерялся
ALTER TABLE pets ADD COLUMN card_lost_location TEXT; -- Где потерялся
ALTER TABLE pets ADD COLUMN card_reward_amount INTEGER; -- Вознаграждение

-- Поля для "Сбор средств"
ALTER TABLE pets ADD COLUMN card_fundraising_goal INTEGER; -- Целевая сумма
ALTER TABLE pets ADD COLUMN card_fundraising_current INTEGER DEFAULT 0; -- Собрано
ALTER TABLE pets ADD COLUMN card_fundraising_purpose TEXT; -- Цель сбора

-- ============================================
-- 2. Индексы
-- ============================================

CREATE INDEX IF NOT EXISTS idx_pets_card_type ON pets(card_type);

-- ============================================
-- 3. Проверка миграции
-- ============================================

PRAGMA table_info(pets);
