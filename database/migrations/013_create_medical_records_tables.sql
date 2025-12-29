-- Миграция 013: Создание таблиц для медицинских записей
-- Дата: 2025-12-29
-- Описание: Таблицы для прививок, обработок от паразитов, операций и болезней

-- Таблица прививок
CREATE TABLE IF NOT EXISTS pet_vaccinations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER NOT NULL,
    vaccination_date DATE NOT NULL,
    vaccine_name TEXT NOT NULL, -- Название вакцины (Нобивак, Эурикан и т.д.)
    vaccine_series TEXT DEFAULT '', -- Серия вакцины
    disease TEXT NOT NULL, -- От какой болезни (бешенство, чума и т.д.)
    next_vaccination_date DATE, -- Дата следующей прививки
    veterinarian_name TEXT DEFAULT '', -- ФИО ветеринара
    clinic_name TEXT DEFAULT '', -- Название клиники
    clinic_stamp TEXT DEFAULT '', -- Печать клиники (путь к изображению)
    notes TEXT DEFAULT '', -- Примечания
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE
);

-- Таблица обработок от паразитов
CREATE TABLE IF NOT EXISTS pet_parasite_treatments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER NOT NULL,
    treatment_date DATE NOT NULL,
    treatment_type TEXT NOT NULL, -- Тип: 'external' (внешние), 'internal' (внутренние), 'both' (комплексная)
    drug_name TEXT NOT NULL, -- Название препарата
    drug_series TEXT DEFAULT '', -- Серия препарата
    dosage TEXT DEFAULT '', -- Дозировка
    next_treatment_date DATE, -- Дата следующей обработки
    veterinarian_name TEXT DEFAULT '',
    notes TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE
);

-- Таблица операций и процедур
CREATE TABLE IF NOT EXISTS pet_surgeries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER NOT NULL,
    surgery_date DATE NOT NULL,
    surgery_type TEXT NOT NULL, -- Тип операции/процедуры
    description TEXT NOT NULL, -- Описание
    veterinarian_name TEXT NOT NULL, -- ФИО ветеринара
    clinic_name TEXT NOT NULL, -- Название клиники
    anesthesia_type TEXT DEFAULT '', -- Тип анестезии
    complications TEXT DEFAULT '', -- Осложнения
    outcome TEXT DEFAULT '', -- Исход
    notes TEXT DEFAULT '', -- Примечания
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE
);

-- Таблица болезней и лечения
CREATE TABLE IF NOT EXISTS pet_medical_history (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER NOT NULL,
    diagnosis_date DATE NOT NULL,
    diagnosis TEXT NOT NULL, -- Диагноз
    symptoms TEXT DEFAULT '', -- Симптомы
    treatment TEXT DEFAULT '', -- Лечение
    medications TEXT DEFAULT '', -- Назначенные лекарства
    recovery_date DATE, -- Дата выздоровления
    veterinarian_name TEXT DEFAULT '',
    clinic_name TEXT DEFAULT '',
    notes TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_vaccinations_pet_id ON pet_vaccinations(pet_id);
CREATE INDEX IF NOT EXISTS idx_vaccinations_date ON pet_vaccinations(vaccination_date);
CREATE INDEX IF NOT EXISTS idx_parasite_treatments_pet_id ON pet_parasite_treatments(pet_id);
CREATE INDEX IF NOT EXISTS idx_parasite_treatments_date ON pet_parasite_treatments(treatment_date);
CREATE INDEX IF NOT EXISTS idx_surgeries_pet_id ON pet_surgeries(pet_id);
CREATE INDEX IF NOT EXISTS idx_surgeries_date ON pet_surgeries(surgery_date);
CREATE INDEX IF NOT EXISTS idx_medical_history_pet_id ON pet_medical_history(pet_id);
CREATE INDEX IF NOT EXISTS idx_medical_history_date ON pet_medical_history(diagnosis_date);
