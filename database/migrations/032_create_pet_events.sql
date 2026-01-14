-- Миграция 032: Таблица pet_events для истории событий питомцев
-- Дата: 2025-01-14
-- Описание: Полная история жизни питомца от рождения до смерти

-- Таблица событий питомца
CREATE TABLE IF NOT EXISTS pet_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER NOT NULL,
    event_type TEXT NOT NULL, -- registration, ownership_change, sterilization, vaccination, treatment, lost, found, death, shelter_intake, adoption
    event_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Кто создал событие
    created_by_user_id INTEGER, -- ID пользователя (владелец, волонтёр)
    created_by_clinic_id INTEGER, -- ID клиники (для медицинских событий)
    created_by_organization_id INTEGER, -- ID организации (для событий приюта)
    
    -- Детали события
    title TEXT, -- Краткое название события
    description TEXT, -- Подробное описание
    details TEXT, -- JSON с дополнительными данными
    
    -- Медицинские данные (для vaccination, treatment, sterilization)
    vaccine_name TEXT, -- Название вакцины
    vaccine_batch TEXT, -- Серия вакцины
    medication_name TEXT, -- Название препарата
    dosage TEXT, -- Дозировка
    next_date DATETIME, -- Дата следующей процедуры
    
    -- Смена владельца (для ownership_change)
    previous_owner_id INTEGER, -- Предыдущий владелец
    new_owner_id INTEGER, -- Новый владелец
    transfer_reason TEXT, -- Причина передачи
    
    -- Потеря/находка (для lost, found)
    location TEXT, -- Место потери/находки
    circumstances TEXT, -- Обстоятельства
    contact_phone TEXT, -- Контактный телефон
    contact_name TEXT, -- Контактное лицо
    
    -- Смерть (для death)
    death_reason TEXT, -- Причина смерти (natural, euthanasia, accident, disease)
    death_confirmed_by_clinic_id INTEGER, -- Клиника, подтвердившая смерть
    
    -- Приют (для shelter_intake, adoption)
    shelter_id INTEGER, -- ID приюта
    adoption_contract TEXT, -- Номер договора пристройства
    
    -- Метаданные
    is_verified BOOLEAN DEFAULT 0, -- Подтверждено ли событие (клиникой/приютом)
    verified_by_user_id INTEGER, -- Кто подтвердил
    verified_at DATETIME, -- Когда подтверждено
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE,
    FOREIGN KEY (created_by_user_id) REFERENCES users(id),
    FOREIGN KEY (created_by_clinic_id) REFERENCES organizations(id),
    FOREIGN KEY (created_by_organization_id) REFERENCES organizations(id),
    FOREIGN KEY (previous_owner_id) REFERENCES users(id),
    FOREIGN KEY (new_owner_id) REFERENCES users(id),
    FOREIGN KEY (death_confirmed_by_clinic_id) REFERENCES organizations(id),
    FOREIGN KEY (shelter_id) REFERENCES organizations(id),
    FOREIGN KEY (verified_by_user_id) REFERENCES users(id)
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_pet_events_pet_id ON pet_events(pet_id);
CREATE INDEX IF NOT EXISTS idx_pet_events_event_type ON pet_events(event_type);
CREATE INDEX IF NOT EXISTS idx_pet_events_event_date ON pet_events(event_date);
CREATE INDEX IF NOT EXISTS idx_pet_events_created_by_user ON pet_events(created_by_user_id);
CREATE INDEX IF NOT EXISTS idx_pet_events_created_by_clinic ON pet_events(created_by_clinic_id);
CREATE INDEX IF NOT EXISTS idx_pet_events_created_by_org ON pet_events(created_by_organization_id);

-- Составной индекс для получения истории питомца
CREATE INDEX IF NOT EXISTS idx_pet_events_pet_date ON pet_events(pet_id, event_date DESC);

-- Индекс для медицинских событий
CREATE INDEX IF NOT EXISTS idx_pet_events_medical ON pet_events(pet_id, event_type) 
WHERE event_type IN ('vaccination', 'treatment', 'sterilization');

-- Индекс для событий приюта
CREATE INDEX IF NOT EXISTS idx_pet_events_shelter ON pet_events(shelter_id, event_date DESC) 
WHERE shelter_id IS NOT NULL;

-- Обновление статистики
ANALYZE;
