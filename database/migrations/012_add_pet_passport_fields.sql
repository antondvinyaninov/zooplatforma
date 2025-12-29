-- Миграция 012: Добавление полей ветеринарного паспорта
-- Дата: 2025-12-29
-- Описание: Добавляем поля для полноценного ветеринарного паспорта питомца

-- Добавляем поля для основной информации
ALTER TABLE pets ADD COLUMN distinctive_marks TEXT DEFAULT ''; -- Особые приметы
ALTER TABLE pets ADD COLUMN tattoo_number TEXT DEFAULT ''; -- Номер клейма (татуировки)

-- Добавляем поля владельца
ALTER TABLE pets ADD COLUMN owner_name TEXT DEFAULT ''; -- ФИО владельца
ALTER TABLE pets ADD COLUMN owner_address TEXT DEFAULT ''; -- Адрес владельца
ALTER TABLE pets ADD COLUMN owner_phone TEXT DEFAULT ''; -- Телефон владельца
ALTER TABLE pets ADD COLUMN owner_email TEXT DEFAULT ''; -- Email владельца

-- Добавляем дополнительные медицинские поля
ALTER TABLE pets ADD COLUMN blood_type TEXT DEFAULT ''; -- Группа крови
ALTER TABLE pets ADD COLUMN allergies TEXT DEFAULT ''; -- Аллергии
ALTER TABLE pets ADD COLUMN chronic_diseases TEXT DEFAULT ''; -- Хронические заболевания
ALTER TABLE pets ADD COLUMN current_medications TEXT DEFAULT ''; -- Текущие лекарства

-- Добавляем поля для документов
ALTER TABLE pets ADD COLUMN pedigree_number TEXT DEFAULT ''; -- Номер родословной
ALTER TABLE pets ADD COLUMN registration_org TEXT DEFAULT ''; -- Организация регистрации (РКФ и т.д.)

-- Примечание: Прививки, обработки, операции и болезни будут храниться в отдельных таблицах
-- для возможности вести полную историю (множественные записи)
