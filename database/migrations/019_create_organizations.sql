-- Миграция 019: Создание таблицы организаций
-- Дата: 2024-12-29
-- Описание: Создание профилей организаций (приюты, ветклиники, зоомагазины, фонды)

-- ============================================
-- 1. Создание таблицы organizations
-- ============================================

CREATE TABLE IF NOT EXISTS organizations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    
    -- Основная информация
    name TEXT NOT NULL,                      -- Название организации
    short_name TEXT,                         -- Краткое название
    legal_form TEXT,                         -- ОПФ (ООО, ИП, АНО, Фонд и т.д.)
    type TEXT DEFAULT 'other',               -- Тип: shelter, vet_clinic, pet_shop, foundation, kennel, other
    
    -- Юридическая информация
    inn TEXT UNIQUE,                         -- ИНН (10 или 12 цифр)
    ogrn TEXT,                               -- ОГРН (13 или 15 цифр)
    kpp TEXT,                                -- КПП (9 цифр)
    registration_date DATE,                  -- Дата регистрации
    
    -- Контактная информация
    email TEXT,                              -- Email
    phone TEXT,                              -- Телефон
    website TEXT,                            -- Сайт
    
    -- Адрес
    address_full TEXT,                       -- Полный адрес
    address_postal_code TEXT,                -- Индекс
    address_region TEXT,                     -- Регион
    address_city TEXT,                       -- Город
    address_street TEXT,                     -- Улица
    address_house TEXT,                      -- Дом
    address_office TEXT,                     -- Офис/помещение
    
    -- Координаты
    geo_lat REAL,                            -- Широта
    geo_lon REAL,                            -- Долгота
    
    -- Описание
    description TEXT,                        -- Описание деятельности
    bio TEXT,                                -- Краткая информация (для карточки)
    
    -- Медиа
    logo TEXT,                               -- Логотип организации
    cover_photo TEXT,                        -- Обложка профиля
    
    -- Руководство
    director_name TEXT,                      -- ФИО руководителя
    director_position TEXT,                  -- Должность руководителя
    
    -- Владелец профиля (пользователь, который создал)
    owner_user_id INTEGER NOT NULL,         -- ID пользователя-владельца
    
    -- Настройки приватности
    profile_visibility TEXT DEFAULT 'public', -- public, private
    show_phone TEXT DEFAULT 'everyone',      -- everyone, nobody
    show_email TEXT DEFAULT 'everyone',      -- everyone, nobody
    allow_messages TEXT DEFAULT 'everyone',  -- everyone, nobody
    
    -- Статус
    is_verified BOOLEAN DEFAULT 0,           -- Верифицирована ли организация
    is_active BOOLEAN DEFAULT 1,             -- Активна ли организация
    status TEXT DEFAULT 'active',            -- active, inactive, blocked
    
    -- Метаданные
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (owner_user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ============================================
-- 2. Индексы
-- ============================================

CREATE INDEX IF NOT EXISTS idx_organizations_inn ON organizations(inn);
CREATE INDEX IF NOT EXISTS idx_organizations_ogrn ON organizations(ogrn);
CREATE INDEX IF NOT EXISTS idx_organizations_type ON organizations(type);
CREATE INDEX IF NOT EXISTS idx_organizations_owner ON organizations(owner_user_id);
CREATE INDEX IF NOT EXISTS idx_organizations_city ON organizations(address_city);
CREATE INDEX IF NOT EXISTS idx_organizations_status ON organizations(status);
CREATE INDEX IF NOT EXISTS idx_organizations_verified ON organizations(is_verified);

-- ============================================
-- 3. Таблица сотрудников организации
-- ============================================

CREATE TABLE IF NOT EXISTS organization_members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    organization_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role TEXT DEFAULT 'member',              -- owner, admin, moderator, member
    position TEXT,                           -- Должность
    can_post BOOLEAN DEFAULT 0,              -- Может создавать посты от имени организации
    can_edit BOOLEAN DEFAULT 0,              -- Может редактировать профиль
    can_manage_members BOOLEAN DEFAULT 0,    -- Может управлять участниками
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(organization_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_org_members_org ON organization_members(organization_id);
CREATE INDEX IF NOT EXISTS idx_org_members_user ON organization_members(user_id);
CREATE INDEX IF NOT EXISTS idx_org_members_role ON organization_members(role);

-- ============================================
-- 4. Таблица питомцев организации
-- ============================================

CREATE TABLE IF NOT EXISTS organization_pets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    organization_id INTEGER NOT NULL,
    pet_id INTEGER NOT NULL,
    added_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    notes TEXT,                              -- Заметки о питомце
    
    FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
    FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE,
    UNIQUE(organization_id, pet_id)
);

CREATE INDEX IF NOT EXISTS idx_org_pets_org ON organization_pets(organization_id);
CREATE INDEX IF NOT EXISTS idx_org_pets_pet ON organization_pets(pet_id);

-- ============================================
-- 5. Проверка миграции
-- ============================================

PRAGMA table_info(organizations);
PRAGMA table_info(organization_members);
PRAGMA table_info(organization_pets);
