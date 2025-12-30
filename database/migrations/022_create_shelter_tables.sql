-- Migration 022: Create Shelter tables with multi-tenancy support
-- Каждая таблица имеет organization_id для изоляции данных между приютами

-- Таблица животных приюта
CREATE TABLE IF NOT EXISTS shelter_animals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    organization_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    species TEXT NOT NULL,
    breed TEXT DEFAULT '',
    age_years INTEGER DEFAULT 0,
    age_months INTEGER DEFAULT 0,
    gender TEXT DEFAULT '',
    color TEXT DEFAULT '',
    size TEXT DEFAULT '',
    weight REAL DEFAULT 0,
    chip_number TEXT DEFAULT '',
    status TEXT DEFAULT 'in_shelter',
    arrival_date DATE NOT NULL,
    arrival_reason TEXT DEFAULT '',
    health_status TEXT DEFAULT '',
    vaccinated BOOLEAN DEFAULT 0,
    sterilized BOOLEAN DEFAULT 0,
    character_traits TEXT DEFAULT '',
    special_needs TEXT DEFAULT '',
    photo TEXT DEFAULT '',
    photos TEXT DEFAULT '',
    story TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE
);

-- Индекс для быстрого поиска по приюту
CREATE INDEX IF NOT EXISTS idx_shelter_animals_org ON shelter_animals(organization_id);
CREATE INDEX IF NOT EXISTS idx_shelter_animals_status ON shelter_animals(organization_id, status);

-- Таблица волонтеров приюта
CREATE TABLE IF NOT EXISTS shelter_volunteers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    organization_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role TEXT DEFAULT 'volunteer',
    status TEXT DEFAULT 'active',
    joined_date DATE NOT NULL,
    hours_contributed INTEGER DEFAULT 0,
    specialization TEXT DEFAULT '',
    notes TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(organization_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_shelter_volunteers_org ON shelter_volunteers(organization_id);
CREATE INDEX IF NOT EXISTS idx_shelter_volunteers_user ON shelter_volunteers(user_id);

-- Таблица заявок на пристройство
CREATE TABLE IF NOT EXISTS adoption_requests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    organization_id INTEGER NOT NULL,
    animal_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status TEXT DEFAULT 'pending',
    request_date DATE NOT NULL,
    decision_date DATE DEFAULT NULL,
    decision_by INTEGER DEFAULT NULL,
    applicant_name TEXT NOT NULL,
    applicant_phone TEXT NOT NULL,
    applicant_email TEXT DEFAULT '',
    applicant_address TEXT DEFAULT '',
    living_conditions TEXT DEFAULT '',
    has_experience BOOLEAN DEFAULT 0,
    has_other_pets BOOLEAN DEFAULT 0,
    other_pets_info TEXT DEFAULT '',
    motivation TEXT DEFAULT '',
    notes TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY (animal_id) REFERENCES shelter_animals(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (decision_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_adoption_requests_org ON adoption_requests(organization_id);
CREATE INDEX IF NOT EXISTS idx_adoption_requests_animal ON adoption_requests(animal_id);
CREATE INDEX IF NOT EXISTS idx_adoption_requests_status ON adoption_requests(organization_id, status);

-- Таблица событий приюта (для истории и логирования)
CREATE TABLE IF NOT EXISTS shelter_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    organization_id INTEGER NOT NULL,
    event_type TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id INTEGER NOT NULL,
    user_id INTEGER DEFAULT NULL,
    description TEXT NOT NULL,
    metadata TEXT DEFAULT '',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_shelter_events_org ON shelter_events(organization_id);
CREATE INDEX IF NOT EXISTS idx_shelter_events_entity ON shelter_events(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_shelter_events_date ON shelter_events(created_at);
