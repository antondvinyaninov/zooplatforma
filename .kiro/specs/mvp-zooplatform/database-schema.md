# MVP ЗооПлатформа - Схема базы данных

**Версия:** 1.0.0  
**Дата:** 28 декабря 2025

---

## 📊 ОБЗОР

Проект использует SQLite на старте с планом миграции на PostgreSQL в будущем.

### Распределение таблиц по микросервисам:

- **Main** - посты, комментарии, лайки, подписки, друзья, чаты
- **Admin** - пользователи, роли, модерация, логи
- **PetID** - реестр животных, события, чипы, виды, породы
- **Shelter** - организации, привязки волонтёров (планируется)

---

## 🐾 PETID МИКРОСЕРВИС

### Таблица: pet_registry
**Назначение:** Единый реестр всех животных (цифровой паспорт)

```sql
CREATE TABLE pet_registry (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,                    -- Имя животного
  species_id INTEGER NOT NULL,           -- Вид (связь с species)
  breed_id INTEGER,                      -- Порода (связь с breeds)
  gender TEXT CHECK(gender IN ('male', 'female', 'unknown')),
  birth_date DATE,                       -- Дата рождения
  
  -- РОДОСЛОВНАЯ
  parent_mother_id INTEGER,              -- Мать (связь с pet_registry)
  parent_father_id INTEGER,              -- Отец (связь с pet_registry)
  breeder_id INTEGER,                    -- Заводчик (user_id или organization_id)
  breeder_type TEXT CHECK(breeder_type IN ('user', 'organization', 'shelter')),
  
  color TEXT,                            -- Окрас
  special_marks TEXT,                    -- Особые приметы
  photos TEXT,                           -- JSON array фотографий
  city TEXT,                             -- Город
  status TEXT NOT NULL DEFAULT 'unknown' CHECK(status IN ('shelter', 'home', 'street', 'deceased', 'unknown')),
  
  -- СМЕРТЬ
  death_date DATE,                       -- Дата смерти
  death_reason TEXT CHECK(death_reason IN ('natural', 'euthanasia', 'accident', 'disease', 'unknown')),
  death_details TEXT,                    -- Детали смерти
  death_confirmed_by_clinic_id INTEGER,  -- Клиника, подтвердившая смерть
  
  responsible_id INTEGER,                -- ID ответственного (user_id)
  responsible_type TEXT CHECK(responsible_type IN ('owner', 'shelter', 'volunteer', 'clinic')),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  
  FOREIGN KEY (species_id) REFERENCES species(id),
  FOREIGN KEY (breed_id) REFERENCES breeds(id),
  FOREIGN KEY (parent_mother_id) REFERENCES pet_registry(id),
  FOREIGN KEY (parent_father_id) REFERENCES pet_registry(id)
);

CREATE INDEX idx_pet_registry_responsible ON pet_registry(responsible_id, responsible_type);
CREATE INDEX idx_pet_registry_status ON pet_registry(status);
CREATE INDEX idx_pet_registry_city ON pet_registry(city);
CREATE INDEX idx_pet_registry_breeder ON pet_registry(breeder_id, breeder_type);
CREATE INDEX idx_pet_registry_parents ON pet_registry(parent_mother_id, parent_father_id);
CREATE INDEX idx_pet_registry_death ON pet_registry(death_date, death_reason);
```

### Таблица: pet_events
**Назначение:** История событий жизни животного

```sql
CREATE TABLE pet_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  pet_id INTEGER NOT NULL,               -- Связь с pet_registry
  event_type TEXT NOT NULL CHECK(event_type IN (
    'registration',
    'ownership_change',
    'sterilization',
    'vaccination',
    'lost',
    'found',
    'death'
  )),
  event_date DATE NOT NULL,              -- Дата события
  description TEXT,                      -- Описание события
  confirmed_by_id INTEGER,               -- Кто подтвердил (user_id)
  confirmed_by_type TEXT CHECK(confirmed_by_type IN ('clinic', 'shelter', 'owner', 'volunteer')),
  confirmation_required BOOLEAN DEFAULT 0, -- Требуется ли подтверждение
  confirmed_at DATETIME,                 -- Дата подтверждения
  metadata TEXT,                         -- JSON с дополнительными данными
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (pet_id) REFERENCES pet_registry(id) ON DELETE CASCADE
);

CREATE INDEX idx_pet_events_pet ON pet_events(pet_id);
CREATE INDEX idx_pet_events_type ON pet_events(event_type);
CREATE INDEX idx_pet_events_date ON pet_events(event_date);
```

**Примеры metadata для разных типов событий:**

```json
// ownership_change
{
  "previous_owner_id": 123,
  "new_owner_id": 456,
  "confirmed_by_previous": true,
  "confirmed_by_new": true
}

// vaccination
{
  "vaccine_name": "Нобивак DHPPi",
  "vaccine_batch": "A12345",
  "next_vaccination_date": "2026-12-28"
}

// death
{
  "cause": "natural",
  "details": "Старость"
}
```

### Таблица: pet_chips
**Назначение:** Чипы, клейма, бирки для идентификации

```sql
CREATE TABLE pet_chips (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  pet_id INTEGER NOT NULL,               -- Связь с pet_registry
  chip_type TEXT NOT NULL CHECK(chip_type IN ('chip', 'tattoo', 'tag')),
  chip_number TEXT NOT NULL UNIQUE,      -- Номер чипа/клейма/бирки
  registered_by_id INTEGER NOT NULL,     -- Кто зарегистрировал (user_id)
  registered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  is_active BOOLEAN DEFAULT 1,           -- Активен ли (может быть деактивирован при потере)
  FOREIGN KEY (pet_id) REFERENCES pet_registry(id) ON DELETE CASCADE
);

CREATE INDEX idx_pet_chips_number ON pet_chips(chip_number);
CREATE INDEX idx_pet_chips_pet ON pet_chips(pet_id);
```

### Таблица: species (уже существует)
**Назначение:** Виды животных

```sql
CREATE TABLE species (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,                    -- Название (Собака, Кошка)
  name_en TEXT NOT NULL,                 -- English name
  description TEXT,                      -- Описание
  icon TEXT,                             -- Иконка/эмодзи
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### Таблица: breeds (уже существует)
**Назначение:** Породы животных

```sql
CREATE TABLE breeds (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  species_id INTEGER NOT NULL,           -- Связь с видом
  name TEXT NOT NULL,                    -- Название породы
  name_en TEXT,                          -- English name
  description TEXT,                      -- Описание
  origin TEXT,                           -- Страна происхождения
  size TEXT,                             -- Размер (small, medium, large)
  weight_min REAL,                       -- Минимальный вес (кг)
  weight_max REAL,                       -- Максимальный вес (кг)
  lifespan_min INTEGER,                  -- Мин. продолжительность жизни (лет)
  lifespan_max INTEGER,                  -- Макс. продолжительность жизни (лет)
  temperament TEXT,                      -- Темперамент
  care_level TEXT,                       -- Уровень ухода (easy, medium, hard)
  photo TEXT,                            -- URL фото
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (species_id) REFERENCES species(id)
);
```

### Таблица: pet_cards (уже существует)
**Назначение:** Подробные карточки пород

```sql
CREATE TABLE pet_cards (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  breed_id INTEGER NOT NULL,             -- Связь с породой
  title TEXT NOT NULL,                   -- Заголовок карточки
  description TEXT,                      -- Полное описание
  characteristics TEXT,                  -- Характеристики (JSON)
  care_tips TEXT,                        -- Советы по уходу
  health_info TEXT,                      -- Информация о здоровье
  nutrition TEXT,                        -- Питание
  photos TEXT,                           -- Фотографии (JSON array)
  is_published BOOLEAN DEFAULT 0,        -- Опубликована ли
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (breed_id) REFERENCES breeds(id)
);
```

---

### Таблица: pet_medical_records (новая)
**Назначение:** Медицинские записи о визитах к ветеринару

```sql
CREATE TABLE pet_medical_records (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  pet_id INTEGER NOT NULL,               -- Связь с pet_registry
  clinic_id INTEGER NOT NULL,            -- Клиника (organization_id)
  vet_doctor_id INTEGER,                 -- Врач (user_id, опционально)
  visit_date DATE NOT NULL,              -- Дата визита
  visit_type TEXT NOT NULL CHECK(visit_type IN ('checkup', 'vaccination', 'treatment', 'surgery', 'emergency')),
  diagnosis TEXT,                        -- Диагноз
  symptoms TEXT,                         -- Симптомы (JSON array)
  treatment_plan TEXT,                   -- Схема лечения
  prescriptions TEXT,                    -- Назначения (JSON array)
  notes TEXT,                            -- Заметки врача
  next_visit_date DATE,                  -- Дата следующего визита
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (pet_id) REFERENCES pet_registry(id) ON DELETE CASCADE
);

CREATE INDEX idx_pet_medical_pet ON pet_medical_records(pet_id);
CREATE INDEX idx_pet_medical_clinic ON pet_medical_records(clinic_id);
CREATE INDEX idx_pet_medical_date ON pet_medical_records(visit_date);
```

### Таблица: pet_lab_tests (новая)
**Назначение:** Лабораторные анализы и исследования

```sql
CREATE TABLE pet_lab_tests (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  medical_record_id INTEGER NOT NULL,    -- Связь с pet_medical_records
  test_type TEXT NOT NULL CHECK(test_type IN ('blood', 'urine', 'feces', 'xray', 'ultrasound', 'other')),
  test_name TEXT NOT NULL,               -- Название анализа
  test_date DATE NOT NULL,               -- Дата проведения
  results TEXT,                          -- Результаты (JSON)
  files TEXT,                            -- Файлы (JSON array путей)
  interpretation TEXT,                   -- Расшифровка
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (medical_record_id) REFERENCES pet_medical_records(id) ON DELETE CASCADE
);

CREATE INDEX idx_pet_lab_tests_record ON pet_lab_tests(medical_record_id);
CREATE INDEX idx_pet_lab_tests_type ON pet_lab_tests(test_type);
```

### Таблица: pet_medications (новая)
**Назначение:** Назначенные препараты

```sql
CREATE TABLE pet_medications (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  medical_record_id INTEGER NOT NULL,    -- Связь с pet_medical_records
  medication_name TEXT NOT NULL,         -- Название препарата
  dosage TEXT NOT NULL,                  -- Дозировка
  frequency TEXT NOT NULL,               -- Частота приёма
  duration_days INTEGER,                 -- Длительность курса (дней)
  start_date DATE NOT NULL,              -- Дата начала
  end_date DATE,                         -- Дата окончания
  notes TEXT,                            -- Примечания
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (medical_record_id) REFERENCES pet_medical_records(id) ON DELETE CASCADE
);

CREATE INDEX idx_pet_medications_record ON pet_medications(medical_record_id);
CREATE INDEX idx_pet_medications_dates ON pet_medications(start_date, end_date);
```

### Таблица: pet_allergies (новая)
**Назначение:** Аллергии животного

```sql
CREATE TABLE pet_allergies (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  pet_id INTEGER NOT NULL,               -- Связь с pet_registry
  allergen TEXT NOT NULL,                -- Аллерген
  severity TEXT NOT NULL CHECK(severity IN ('mild', 'moderate', 'severe')),
  reaction TEXT,                         -- Реакция
  diagnosed_date DATE,                   -- Дата диагностики
  diagnosed_by_clinic_id INTEGER,        -- Клиника
  notes TEXT,                            -- Примечания
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (pet_id) REFERENCES pet_registry(id) ON DELETE CASCADE
);

CREATE INDEX idx_pet_allergies_pet ON pet_allergies(pet_id);
CREATE INDEX idx_pet_allergies_severity ON pet_allergies(severity);
```

### Таблица: clinic_schedule (новая)
**Назначение:** Расписание работы клиники

```sql
CREATE TABLE clinic_schedule (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  clinic_id INTEGER NOT NULL,            -- Организация (клиника)
  doctor_id INTEGER,                     -- Врач (опционально)
  day_of_week INTEGER NOT NULL CHECK(day_of_week BETWEEN 0 AND 6), -- 0=Пн, 6=Вс
  start_time TIME NOT NULL,              -- Время начала работы
  end_time TIME NOT NULL,                -- Время окончания работы
  slot_duration INTEGER DEFAULT 30,      -- Длительность слота (минуты)
  is_active BOOLEAN DEFAULT 1,           -- Активно ли расписание
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_clinic_schedule_clinic ON clinic_schedule(clinic_id);
CREATE INDEX idx_clinic_schedule_doctor ON clinic_schedule(doctor_id);
CREATE INDEX idx_clinic_schedule_day ON clinic_schedule(day_of_week);
```

### Таблица: clinic_appointments (новая)
**Назначение:** Записи на приём к ветеринару

```sql
CREATE TABLE clinic_appointments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  clinic_id INTEGER NOT NULL,            -- Клиника
  doctor_id INTEGER,                     -- Врач (опционально)
  pet_id INTEGER NOT NULL,               -- Животное
  owner_id INTEGER NOT NULL,             -- Владелец
  appointment_date DATE NOT NULL,        -- Дата приёма
  appointment_time TIME NOT NULL,        -- Время приёма
  duration INTEGER DEFAULT 30,           -- Длительность (минуты)
  visit_type TEXT NOT NULL CHECK(visit_type IN ('checkup', 'vaccination', 'treatment', 'emergency')),
  reason TEXT,                           -- Причина визита
  status TEXT NOT NULL DEFAULT 'scheduled' CHECK(status IN ('scheduled', 'confirmed', 'completed', 'cancelled', 'no_show')),
  notes TEXT,                            -- Заметки клиники
  reminder_sent BOOLEAN DEFAULT 0,       -- Отправлено ли напоминание
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (pet_id) REFERENCES pet_registry(id),
  FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE INDEX idx_clinic_appointments_clinic ON clinic_appointments(clinic_id);
CREATE INDEX idx_clinic_appointments_doctor ON clinic_appointments(doctor_id);
CREATE INDEX idx_clinic_appointments_pet ON clinic_appointments(pet_id);
CREATE INDEX idx_clinic_appointments_owner ON clinic_appointments(owner_id);
CREATE INDEX idx_clinic_appointments_date ON clinic_appointments(appointment_date, appointment_time);
CREATE INDEX idx_clinic_appointments_status ON clinic_appointments(status);
```

### Таблица: clinic_time_off (новая)
**Назначение:** Выходные дни и отпуска клиники/врачей

```sql
CREATE TABLE clinic_time_off (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  clinic_id INTEGER NOT NULL,            -- Клиника
  doctor_id INTEGER,                     -- Врач (если null - вся клиника)
  start_date DATE NOT NULL,              -- Дата начала
  end_date DATE NOT NULL,                -- Дата окончания
  reason TEXT CHECK(reason IN ('vacation', 'sick', 'training', 'other')),
  notes TEXT,                            -- Примечания
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_clinic_time_off_clinic ON clinic_time_off(clinic_id);
CREATE INDEX idx_clinic_time_off_doctor ON clinic_time_off(doctor_id);
CREATE INDEX idx_clinic_time_off_dates ON clinic_time_off(start_date, end_date);
```

---

## 📱 MAIN МИКРОСЕРВИС

### Таблица: post_pets (новая)
**Назначение:** Связь постов с животными из PetID

```sql
CREATE TABLE post_pets (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  post_id INTEGER NOT NULL,              -- Связь с posts
  pet_id INTEGER NOT NULL,               -- Связь с pet_registry (PetID)
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

CREATE INDEX idx_post_pets_post ON post_pets(post_id);
CREATE INDEX idx_post_pets_pet ON post_pets(pet_id);
```

### Таблица: posts (обновить)
**Назначение:** Посты пользователей

```sql
-- Добавить новое поле к существующей таблице
ALTER TABLE posts ADD COLUMN pet_tag TEXT CHECK(pet_tag IN ('looking_for_home', 'lost', 'found', 'new_home', NULL));

CREATE INDEX idx_posts_pet_tag ON posts(pet_tag);
```

**Значения pet_tag:**
- `looking_for_home` - "ищет дом"
- `lost` - "потеряшка"
- `found` - "нашёлся"
- `new_home` - "новоселье"
- `NULL` - обычный пост

### Таблица: subscriptions (новая)
**Назначение:** Подписки пользователей друг на друга

```sql
CREATE TABLE subscriptions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  follower_id INTEGER NOT NULL,          -- Кто подписался
  following_id INTEGER NOT NULL,         -- На кого подписались
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE(follower_id, following_id)
);

CREATE INDEX idx_subscriptions_follower ON subscriptions(follower_id);
CREATE INDEX idx_subscriptions_following ON subscriptions(following_id);
```

### Таблица: chats (новая)
**Назначение:** Чаты между пользователями

```sql
CREATE TABLE chats (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  type TEXT NOT NULL DEFAULT 'direct' CHECK(type IN ('direct', 'group', 'pet', 'city')),
  name TEXT,                             -- Название (для групповых)
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_chats_type ON chats(type);
```

**Типы чатов:**
- `direct` - личный чат 1-1 (MVP)
- `group` - групповой чат (v2.0)
- `pet` - чат по питомцу (v2.0)
- `city` - чат города (v2.0)

### Таблица: chat_members (новая)
**Назначение:** Участники чатов

```sql
CREATE TABLE chat_members (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  chat_id INTEGER NOT NULL,              -- Связь с chats
  user_id INTEGER NOT NULL,              -- Участник
  joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  last_read_at DATETIME,                 -- Последнее прочтение
  FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  UNIQUE(chat_id, user_id)
);

CREATE INDEX idx_chat_members_chat ON chat_members(chat_id);
CREATE INDEX idx_chat_members_user ON chat_members(user_id);
```

### Таблица: messages (новая)
**Назначение:** Сообщения в чатах

```sql
CREATE TABLE messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  chat_id INTEGER NOT NULL,              -- Связь с chats
  sender_id INTEGER NOT NULL,            -- Отправитель
  content TEXT NOT NULL,                 -- Текст сообщения
  is_read BOOLEAN DEFAULT 0,             -- Прочитано ли
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
  FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_messages_chat ON messages(chat_id);
CREATE INDEX idx_messages_sender ON messages(sender_id);
CREATE INDEX idx_messages_created ON messages(created_at);
```

### Таблица: notifications (новая)
**Назначение:** Уведомления пользователей

```sql
CREATE TABLE notifications (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,              -- Кому уведомление
  type TEXT NOT NULL CHECK(type IN ('comment', 'like', 'pet_event', 'message', 'subscription', 'friend_request')),
  title TEXT NOT NULL,                   -- Заголовок
  content TEXT,                          -- Содержание
  link TEXT,                             -- Ссылка на объект
  is_read BOOLEAN DEFAULT 0,             -- Прочитано ли
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_read ON notifications(is_read);
CREATE INDEX idx_notifications_created ON notifications(created_at);
```

### Таблица: pets (удалить)
**Назначение:** Старая таблица питомцев - будет удалена

```sql
-- Миграция данных в pet_registry (PetID)
-- Затем удаление таблицы
DROP TABLE pets;
```

---

## 👥 ADMIN МИКРОСЕРВИС

### Таблица: user_roles (новая)
**Назначение:** Роли пользователей

```sql
CREATE TABLE user_roles (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,              -- Пользователь
  role TEXT NOT NULL CHECK(role IN (
    'user',
    'volunteer',
    'shelter_admin',
    'clinic_admin',
    'moderator',
    'superadmin'
  )),
  organization_id INTEGER,               -- Привязка к организации (для shelter_admin, clinic_admin)
  granted_by_id INTEGER,                 -- Кто назначил роль
  granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_id, role, organization_id)
);

CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_user_roles_role ON user_roles(role);
CREATE INDEX idx_user_roles_org ON user_roles(organization_id);
```

**Роли:**
- `user` - обычный пользователь (по умолчанию)
- `volunteer` - волонтёр (куратор животных)
- `shelter_admin` - администратор приюта
- `clinic_admin` - администратор ветклиники
- `moderator` - модератор контента
- `superadmin` - суперадминистратор

### Таблица: users (обновить)
**Назначение:** Пользователи системы

```sql
-- Добавить новое поле к существующей таблице
ALTER TABLE users ADD COLUMN is_verified BOOLEAN DEFAULT 0;

CREATE INDEX idx_users_verified ON users(is_verified);
```

---

## 🏥 SHELTER МИКРОСЕРВИС

### Таблица: organizations
**Назначение:** Приюты и ветклиники

```sql
CREATE TABLE organizations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  type TEXT NOT NULL CHECK(type IN ('shelter', 'clinic')),
  name TEXT NOT NULL,                    -- Название
  inn TEXT,                              -- ИНН
  ogrn TEXT,                             -- ОГРН
  address TEXT,                          -- Адрес
  city TEXT,                             -- Город
  phone TEXT,                            -- Телефон
  email TEXT,                            -- Email
  website TEXT,                          -- Сайт
  description TEXT,                      -- Описание
  logo TEXT,                             -- Логотип
  is_verified BOOLEAN DEFAULT 0,         -- Проверена ли
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_organizations_type ON organizations(type);
CREATE INDEX idx_organizations_city ON organizations(city);
```

### Таблица: organization_members
**Назначение:** Сотрудники организаций

```sql
CREATE TABLE organization_members (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  organization_id INTEGER NOT NULL,      -- Организация
  user_id INTEGER NOT NULL,              -- Пользователь
  role TEXT NOT NULL CHECK(role IN ('admin', 'member', 'volunteer')),
  joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (organization_id) REFERENCES organizations(id) ON DELETE CASCADE,
  UNIQUE(organization_id, user_id)
);

CREATE INDEX idx_org_members_org ON organization_members(organization_id);
CREATE INDEX idx_org_members_user ON organization_members(user_id);
```

### Таблица: shelter_animals (новая)
**Назначение:** Животные в приюте

```sql
CREATE TABLE shelter_animals (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  shelter_id INTEGER NOT NULL,           -- Приют
  pet_id INTEGER NOT NULL,               -- Животное (pet_registry)
  intake_date DATE NOT NULL,             -- Дата поступления
  intake_reason TEXT CHECK(intake_reason IN ('stray', 'owner_surrender', 'rescue', 'osvv')),
  intake_condition TEXT CHECK(intake_condition IN ('healthy', 'sick', 'injured', 'pregnant')),
  location_in_shelter TEXT,              -- Расположение в приюте
  assigned_volunteer_id INTEGER,         -- Назначенный волонтёр
  adoption_date DATE,                    -- Дата пристройства
  adoption_to_user_id INTEGER,           -- Кому пристроен
  status TEXT NOT NULL DEFAULT 'in_shelter' CHECK(status IN ('in_shelter', 'adopted', 'returned', 'deceased', 'osvv_returned')),
  notes TEXT,                            -- Заметки
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (shelter_id) REFERENCES organizations(id),
  FOREIGN KEY (pet_id) REFERENCES pet_registry(id)
);

CREATE INDEX idx_shelter_animals_shelter ON shelter_animals(shelter_id);
CREATE INDEX idx_shelter_animals_pet ON shelter_animals(pet_id);
CREATE INDEX idx_shelter_animals_status ON shelter_animals(status);
CREATE INDEX idx_shelter_animals_volunteer ON shelter_animals(assigned_volunteer_id);
```

### Таблица: shelter_osvv_program (новая)
**Назначение:** Программа ОСВВ (Отлов-Стерилизация-Вакцинация-Возврат)

```sql
CREATE TABLE shelter_osvv_program (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  shelter_id INTEGER NOT NULL,           -- Приют
  pet_id INTEGER NOT NULL,               -- Животное
  request_id INTEGER,                    -- Связь с заявкой (если есть)
  capture_date DATE NOT NULL,            -- Дата отлова
  capture_location TEXT NOT NULL,        -- Место отлова (адрес)
  capture_by TEXT,                       -- Кто отловил
  sterilization_date DATE,               -- Дата стерилизации
  sterilization_clinic_id INTEGER,       -- Клиника
  vaccination_date DATE,                 -- Дата вакцинации
  vaccination_details TEXT,              -- Детали вакцинации (JSON)
  return_date DATE,                      -- Дата возврата
  return_location TEXT,                  -- Место возврата
  ear_mark TEXT,                         -- Метка на ухе
  status TEXT NOT NULL DEFAULT 'captured' CHECK(status IN ('captured', 'in_treatment', 'ready_to_return', 'returned')),
  photos_before TEXT,                    -- Фото до (JSON)
  photos_after TEXT,                     -- Фото после (JSON)
  notes TEXT,                            -- Заметки
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (shelter_id) REFERENCES organizations(id),
  FOREIGN KEY (pet_id) REFERENCES pet_registry(id),
  FOREIGN KEY (request_id) REFERENCES shelter_osvv_requests(id)
);

CREATE INDEX idx_shelter_osvv_shelter ON shelter_osvv_program(shelter_id);
CREATE INDEX idx_shelter_osvv_pet ON shelter_osvv_program(pet_id);
CREATE INDEX idx_shelter_osvv_request ON shelter_osvv_program(request_id);
CREATE INDEX idx_shelter_osvv_status ON shelter_osvv_program(status);
CREATE INDEX idx_shelter_osvv_dates ON shelter_osvv_program(capture_date, return_date);
```

### Таблица: shelter_osvv_requests (новая)
**Назначение:** Заявки от жителей на ОСВВ

```sql
CREATE TABLE shelter_osvv_requests (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  shelter_id INTEGER,                    -- Приют (назначается при принятии)
  requester_user_id INTEGER,             -- Заявитель (если зарегистрирован)
  requester_name TEXT,                   -- Имя заявителя (если не зарегистрирован)
  requester_phone TEXT NOT NULL,         -- Телефон заявителя
  requester_email TEXT,                  -- Email заявителя
  location_address TEXT NOT NULL,        -- Адрес (где животное)
  location_lat REAL,                     -- Широта (для карты)
  location_lng REAL,                     -- Долгота (для карты)
  animal_description TEXT NOT NULL,      -- Описание животного
  animal_count INTEGER DEFAULT 1,        -- Количество животных
  urgency TEXT DEFAULT 'normal' CHECK(urgency IN ('low', 'normal', 'high', 'emergency')),
  photos TEXT,                           -- Фото (JSON array)
  additional_info TEXT,                  -- Дополнительная информация
  status TEXT NOT NULL DEFAULT 'new' CHECK(status IN ('new', 'accepted', 'in_progress', 'completed', 'rejected', 'duplicate')),
  assigned_to_user_id INTEGER,           -- Назначен ответственный
  rejection_reason TEXT,                 -- Причина отклонения
  completed_osvv_id INTEGER,             -- Связь с выполненным ОСВВ
  internal_notes TEXT,                   -- Внутренние заметки приюта
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (shelter_id) REFERENCES organizations(id),
  FOREIGN KEY (requester_user_id) REFERENCES users(id),
  FOREIGN KEY (completed_osvv_id) REFERENCES shelter_osvv_program(id)
);

CREATE INDEX idx_shelter_osvv_req_shelter ON shelter_osvv_requests(shelter_id);
CREATE INDEX idx_shelter_osvv_req_status ON shelter_osvv_requests(status);
CREATE INDEX idx_shelter_osvv_req_urgency ON shelter_osvv_requests(urgency);
CREATE INDEX idx_shelter_osvv_req_location ON shelter_osvv_requests(location_lat, location_lng);
CREATE INDEX idx_shelter_osvv_req_created ON shelter_osvv_requests(created_at);
```

### Таблица: shelter_osvv_request_updates (новая)
**Назначение:** История обновлений заявки

```sql
CREATE TABLE shelter_osvv_request_updates (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  request_id INTEGER NOT NULL,           -- Заявка
  user_id INTEGER,                       -- Кто обновил
  old_status TEXT,                       -- Старый статус
  new_status TEXT,                       -- Новый статус
  comment TEXT,                          -- Комментарий
  is_public BOOLEAN DEFAULT 0,           -- Видно ли заявителю
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (request_id) REFERENCES shelter_osvv_requests(id) ON DELETE CASCADE
);

CREATE INDEX idx_shelter_osvv_updates_request ON shelter_osvv_request_updates(request_id);
CREATE INDEX idx_shelter_osvv_updates_created ON shelter_osvv_request_updates(created_at);
```

### Таблица: shelter_visits (новая)
**Назначение:** Записи на посещение приюта

```sql
CREATE TABLE shelter_visits (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  shelter_id INTEGER NOT NULL,           -- Приют
  visitor_user_id INTEGER NOT NULL,      -- Посетитель
  visit_date DATE NOT NULL,              -- Дата визита
  visit_time TIME NOT NULL,              -- Время
  visit_type TEXT NOT NULL CHECK(visit_type IN ('adoption_viewing', 'volunteer_work', 'event', 'tour')),
  interested_pets TEXT,                  -- Животные (JSON array pet_ids)
  status TEXT NOT NULL DEFAULT 'scheduled' CHECK(status IN ('scheduled', 'confirmed', 'completed', 'cancelled', 'no_show')),
  notes TEXT,                            -- Заметки приюта
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (shelter_id) REFERENCES organizations(id),
  FOREIGN KEY (visitor_user_id) REFERENCES users(id)
);

CREATE INDEX idx_shelter_visits_shelter ON shelter_visits(shelter_id);
CREATE INDEX idx_shelter_visits_visitor ON shelter_visits(visitor_user_id);
CREATE INDEX idx_shelter_visits_date ON shelter_visits(visit_date, visit_time);
CREATE INDEX idx_shelter_visits_status ON shelter_visits(status);
```

### Таблица: shelter_events (новая)
**Назначение:** Мероприятия приюта

```sql
CREATE TABLE shelter_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  shelter_id INTEGER NOT NULL,           -- Приют
  event_type TEXT NOT NULL CHECK(event_type IN ('adoption_day', 'volunteer_meeting', 'fundraiser', 'educational', 'other')),
  title TEXT NOT NULL,                   -- Название
  description TEXT,                      -- Описание
  event_date DATE NOT NULL,              -- Дата
  start_time TIME NOT NULL,              -- Время начала
  end_time TIME,                         -- Время окончания
  location TEXT,                         -- Место проведения
  max_participants INTEGER,              -- Макс. участников
  current_participants INTEGER DEFAULT 0, -- Текущее кол-во
  is_public BOOLEAN DEFAULT 1,           -- Публичное ли
  registration_required BOOLEAN DEFAULT 0, -- Требуется ли регистрация
  status TEXT NOT NULL DEFAULT 'planned' CHECK(status IN ('planned', 'active', 'completed', 'cancelled')),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (shelter_id) REFERENCES organizations(id)
);

CREATE INDEX idx_shelter_events_shelter ON shelter_events(shelter_id);
CREATE INDEX idx_shelter_events_date ON shelter_events(event_date);
CREATE INDEX idx_shelter_events_status ON shelter_events(status);
CREATE INDEX idx_shelter_events_public ON shelter_events(is_public);
```

### Таблица: shelter_event_registrations (новая)
**Назначение:** Регистрации на мероприятия

```sql
CREATE TABLE shelter_event_registrations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  event_id INTEGER NOT NULL,             -- Мероприятие
  user_id INTEGER NOT NULL,              -- Пользователь
  registration_date DATETIME DEFAULT CURRENT_TIMESTAMP,
  status TEXT NOT NULL DEFAULT 'registered' CHECK(status IN ('registered', 'confirmed', 'attended', 'cancelled')),
  notes TEXT,                            -- Заметки
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (event_id) REFERENCES shelter_events(id) ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES users(id),
  UNIQUE(event_id, user_id)
);

CREATE INDEX idx_shelter_event_reg_event ON shelter_event_registrations(event_id);
CREATE INDEX idx_shelter_event_reg_user ON shelter_event_registrations(user_id);
CREATE INDEX idx_shelter_event_reg_status ON shelter_event_registrations(status);
```

### Таблица: volunteer_assignments
**Назначение:** Привязка волонтёров к животным

```sql
CREATE TABLE volunteer_assignments (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  volunteer_id INTEGER NOT NULL,         -- Волонтёр (user_id)
  pet_id INTEGER NOT NULL,               -- Животное (pet_registry.id)
  organization_id INTEGER,               -- Организация (если через приют)
  assigned_by_id INTEGER,                -- Кто назначил
  assigned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  notes TEXT,                            -- Заметки
  UNIQUE(volunteer_id, pet_id)
);

CREATE INDEX idx_volunteer_assignments_volunteer ON volunteer_assignments(volunteer_id);
CREATE INDEX idx_volunteer_assignments_pet ON volunteer_assignments(pet_id);
```

### Таблица: volunteer_tasks
**Назначение:** Задачи и чек-листы для волонтёров

```sql
CREATE TABLE volunteer_tasks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  volunteer_id INTEGER NOT NULL,         -- Волонтёр
  pet_id INTEGER,                        -- Животное (опционально)
  task_type TEXT CHECK(task_type IN ('medical', 'feeding', 'walking', 'post', 'fundraising', 'transport', 'other')),
  title TEXT NOT NULL,                   -- Название задачи
  description TEXT,                      -- Описание
  priority TEXT DEFAULT 'normal' CHECK(priority IN ('low', 'normal', 'high', 'urgent')),
  due_date DATE,                         -- Срок выполнения
  is_completed BOOLEAN DEFAULT 0,        -- Выполнена ли
  completed_at DATETIME,                 -- Дата выполнения
  assigned_by_id INTEGER,                -- Кто назначил
  notes TEXT,                            -- Заметки
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_volunteer_tasks_volunteer ON volunteer_tasks(volunteer_id);
CREATE INDEX idx_volunteer_tasks_pet ON volunteer_tasks(pet_id);
CREATE INDEX idx_volunteer_tasks_completed ON volunteer_tasks(is_completed);
CREATE INDEX idx_volunteer_tasks_priority ON volunteer_tasks(priority);
CREATE INDEX idx_volunteer_tasks_due ON volunteer_tasks(due_date);
```

### Таблица: volunteer_fosters (новая)
**Назначение:** Передержки волонтёров

```sql
CREATE TABLE volunteer_fosters (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  volunteer_id INTEGER NOT NULL,         -- Волонтёр-передержка
  pet_id INTEGER,                        -- Животное (если конкретное)
  start_date DATE NOT NULL,              -- Дата начала
  end_date DATE,                         -- Дата окончания (null если активна)
  foster_type TEXT CHECK(foster_type IN ('temporary', 'long_term', 'medical')),
  location_address TEXT,                 -- Адрес передержки
  location_lat REAL,                     -- Широта
  location_lng REAL,                     -- Долгота
  max_capacity INTEGER DEFAULT 1,        -- Макс. животных одновременно
  current_count INTEGER DEFAULT 0,       -- Текущее количество
  conditions TEXT,                       -- Условия
  photos TEXT,                           -- Фото условий (JSON)
  status TEXT DEFAULT 'active' CHECK(status IN ('active', 'paused', 'closed')),
  notes TEXT,                            -- Заметки
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (volunteer_id) REFERENCES users(id),
  FOREIGN KEY (pet_id) REFERENCES pet_registry(id)
);

CREATE INDEX idx_volunteer_fosters_volunteer ON volunteer_fosters(volunteer_id);
CREATE INDEX idx_volunteer_fosters_pet ON volunteer_fosters(pet_id);
CREATE INDEX idx_volunteer_fosters_status ON volunteer_fosters(status);
CREATE INDEX idx_volunteer_fosters_location ON volunteer_fosters(location_lat, location_lng);
```

### Таблица: fundraisers (новая)
**Назначение:** Сборы средств

```sql
CREATE TABLE fundraisers (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  creator_id INTEGER NOT NULL,           -- Кто создал (волонтёр)
  pet_id INTEGER,                        -- Животное (опционально)
  title TEXT NOT NULL,                   -- Название
  description TEXT NOT NULL,             -- Описание
  goal_amount DECIMAL(10,2) NOT NULL,    -- Целевая сумма
  current_amount DECIMAL(10,2) DEFAULT 0, -- Собрано
  currency TEXT DEFAULT 'RUB',           -- Валюта
  purpose TEXT CHECK(purpose IN ('medical', 'food', 'shelter', 'transport', 'other')),
  deadline DATE,                         -- Срок
  status TEXT DEFAULT 'active' CHECK(status IN ('active', 'completed', 'cancelled')),
  beneficiary_type TEXT CHECK(beneficiary_type IN ('volunteer', 'shelter', 'clinic')),
  beneficiary_id INTEGER,                -- ID получателя
  payment_details TEXT,                  -- Реквизиты (JSON)
  is_verified BOOLEAN DEFAULT 0,         -- Проверен ли
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (creator_id) REFERENCES users(id),
  FOREIGN KEY (pet_id) REFERENCES pet_registry(id)
);

CREATE INDEX idx_fundraisers_creator ON fundraisers(creator_id);
CREATE INDEX idx_fundraisers_pet ON fundraisers(pet_id);
CREATE INDEX idx_fundraisers_status ON fundraisers(status);
CREATE INDEX idx_fundraisers_deadline ON fundraisers(deadline);
```

### Таблица: fundraiser_donations (новая)
**Назначение:** Пожертвования

```sql
CREATE TABLE fundraiser_donations (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  fundraiser_id INTEGER NOT NULL,        -- Сбор
  donor_user_id INTEGER,                 -- Жертвователь (если зарегистрирован)
  donor_name TEXT,                       -- Имя (если анонимный)
  amount DECIMAL(10,2) NOT NULL,         -- Сумма
  currency TEXT DEFAULT 'RUB',           -- Валюта
  payment_method TEXT,                   -- Способ (card/sbp/yoomoney/cash)
  payment_id TEXT,                       -- ID платежа
  is_anonymous BOOLEAN DEFAULT 0,        -- Анонимный ли
  message TEXT,                          -- Сообщение
  status TEXT DEFAULT 'completed' CHECK(status IN ('pending', 'completed', 'failed', 'refunded')),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (fundraiser_id) REFERENCES fundraisers(id) ON DELETE CASCADE,
  FOREIGN KEY (donor_user_id) REFERENCES users(id)
);

CREATE INDEX idx_fundraiser_donations_fundraiser ON fundraiser_donations(fundraiser_id);
CREATE INDEX idx_fundraiser_donations_donor ON fundraiser_donations(donor_user_id);
CREATE INDEX idx_fundraiser_donations_status ON fundraiser_donations(status);
```

### Таблица: fundraiser_expenses (новая)
**Назначение:** Расходы из сборов

```sql
CREATE TABLE fundraiser_expenses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  fundraiser_id INTEGER NOT NULL,        -- Сбор
  amount DECIMAL(10,2) NOT NULL,         -- Сумма
  purpose TEXT NOT NULL,                 -- Назначение
  description TEXT,                      -- Описание
  receipt_photo TEXT,                    -- Фото чека
  date DATE NOT NULL,                    -- Дата расхода
  added_by_id INTEGER NOT NULL,          -- Кто добавил
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (fundraiser_id) REFERENCES fundraisers(id) ON DELETE CASCADE,
  FOREIGN KEY (added_by_id) REFERENCES users(id)
);

CREATE INDEX idx_fundraiser_expenses_fundraiser ON fundraiser_expenses(fundraiser_id);
CREATE INDEX idx_fundraiser_expenses_date ON fundraiser_expenses(date);
```

### Таблица: volunteer_activity_log (новая)
**Назначение:** Лог активности волонтёра

```sql
CREATE TABLE volunteer_activity_log (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  volunteer_id INTEGER NOT NULL,         -- Волонтёр
  activity_type TEXT NOT NULL CHECK(activity_type IN ('post_created', 'pet_updated', 'task_completed', 'donation_received', 'event_attended', 'pet_adopted')),
  pet_id INTEGER,                        -- Животное (если применимо)
  description TEXT,                      -- Описание
  points INTEGER DEFAULT 0,              -- Баллы активности
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (volunteer_id) REFERENCES users(id),
  FOREIGN KEY (pet_id) REFERENCES pet_registry(id)
);

CREATE INDEX idx_volunteer_activity_volunteer ON volunteer_activity_log(volunteer_id);
CREATE INDEX idx_volunteer_activity_type ON volunteer_activity_log(activity_type);
CREATE INDEX idx_volunteer_activity_created ON volunteer_activity_log(created_at);
```

---

## 🔄 МИГРАЦИЯ ДАННЫХ

### Миграция pets → pet_registry

```sql
-- Шаг 1: Создать записи в pet_registry из существующих pets
INSERT INTO pet_registry (
  name,
  species_id,
  breed_id,
  photos,
  responsible_id,
  responsible_type,
  status,
  created_at
)
SELECT
  p.name,
  (SELECT id FROM species WHERE name = p.species LIMIT 1) as species_id,
  NULL as breed_id,  -- породы не было в старой таблице
  json_array(p.photo) as photos,
  p.user_id as responsible_id,
  'owner' as responsible_type,
  'home' as status,
  p.created_at
FROM pets p;

-- Шаг 2: Создать события регистрации для всех животных
INSERT INTO pet_events (
  pet_id,
  event_type,
  event_date,
  description,
  confirmed_by_id,
  confirmed_by_type,
  confirmed_at,
  created_at
)
SELECT
  pr.id as pet_id,
  'registration' as event_type,
  pr.created_at as event_date,
  'Миграция из старой системы' as description,
  pr.responsible_id as confirmed_by_id,
  'owner' as confirmed_by_type,
  pr.created_at as confirmed_at,
  pr.created_at
FROM pet_registry pr;

-- Шаг 3: Удалить старую таблицу
DROP TABLE pets;
```

---

## 📝 ПРИМЕЧАНИЯ

### Индексы
Все таблицы имеют индексы на часто используемые поля для оптимизации запросов.

### Каскадное удаление
Используется `ON DELETE CASCADE` для автоматического удаления связанных записей.

### JSON поля
Поля с JSON данными (photos, metadata, characteristics) хранятся как TEXT в SQLite.

### Миграция на PostgreSQL
При переходе на PostgreSQL:
- JSON поля станут типом JSONB
- INTEGER → BIGSERIAL для ID
- DATETIME → TIMESTAMP WITH TIME ZONE
- Добавить полнотекстовый поиск (tsvector)

### Резервное копирование
Перед любой миграцией создавать резервную копию:
```bash
cp database/data.db database/backups/data_$(date +%Y%m%d_%H%M%S).db
```

---

**Последнее обновление:** 28 декабря 2025  
**Версия схемы:** 1.0.0
