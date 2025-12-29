-- Миграция 016: Добавление даты стерилизации
-- Дата: 2024-12-29

-- Добавляем поле sterilization_date (дата стерилизации/кастрации)
ALTER TABLE pets ADD COLUMN sterilization_date DATE DEFAULT NULL;

-- Комментарий:
-- sterilization_date: Дата проведения операции стерилизации/кастрации
-- NULL если операция не проводилась или дата неизвестна
