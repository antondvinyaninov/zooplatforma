-- Добавление полей профиля в таблицу users

-- Добавить bio (описание профиля)
ALTER TABLE users ADD COLUMN bio TEXT DEFAULT '';

-- Добавить phone (телефон)
ALTER TABLE users ADD COLUMN phone TEXT DEFAULT '';

-- Добавить date_of_birth (дата рождения)
ALTER TABLE users ADD COLUMN date_of_birth DATE;

-- Добавить gender (пол)
ALTER TABLE users ADD COLUMN gender TEXT DEFAULT '';

-- Добавить city (город)
ALTER TABLE users ADD COLUMN city TEXT DEFAULT '';

-- Добавить country (страна)
ALTER TABLE users ADD COLUMN country TEXT DEFAULT '';
