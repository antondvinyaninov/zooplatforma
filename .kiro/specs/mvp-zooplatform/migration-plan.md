# План миграции PetBase → PetID

**Дата:** 28 декабря 2025  
**Статус:** Планирование

---

## 🎯 ЦЕЛЬ

Расширить микросервис PetBase (справочник пород) до PetID (единый реестр животных) с сохранением всех существующих данных.

---

## 📋 ЧТО МЕНЯЕТСЯ

### Переименование
- `petbase/` → `petid/`
- Все упоминания "PetBase" → "PetID" в коде и документации
- URL остаётся: `localhost:4100` и `localhost:8100`

### Расширение функциональности
**Было:** Справочник видов и пород животных  
**Станет:** Справочник + Единый реестр всех животных с историей событий

### Новые таблицы
- `pet_registry` - реестр животных
- `pet_events` - история событий
- `pet_chips` - чипы и метки

### Удаление из Main
- Таблица `pets` будет удалена
- Данные мигрируют в `pet_registry` (PetID)
- API endpoints `/api/pets/*` будут перенаправлять на PetID API

---

## 🗂️ ЭТАПЫ МИГРАЦИИ

### ЭТАП 1: Подготовка (1 день)

#### 1.1. Резервное копирование
```bash
# Создать резервную копию базы данных
DATE=$(date +%Y%m%d_%H%M%S)
cp database/data.db database/backups/data_before_petid_migration_${DATE}.db

# Создать SQL dump
sqlite3 database/data.db .dump > database/backups/backup_before_petid_${DATE}.sql
```

#### 1.2. Создать ветку в Git
```bash
git checkout -b feature/petbase-to-petid-migration
git add -A
git commit -m "Checkpoint before PetBase → PetID migration"
```

#### 1.3. Документация текущего состояния
- [ ] Список всех файлов, которые будут изменены
- [ ] Список всех API endpoints, которые будут изменены
- [ ] Список всех компонентов frontend, использующих pets API

---

### ЭТАП 2: Переименование директории (1 день)

#### 2.1. Переименовать директорию
```bash
mv petbase petid
```

#### 2.2. Обновить файлы в petid/

**petid/backend/main.go:**
```go
// Изменить сообщения
fmt.Printf("🐾 PetID API starting on port %s\n", port)

// Изменить welcome message
fmt.Fprintf(w, `{"message": "PetID API", "version": "2.0.0"}`)
```

**petid/backend/go.mod:**
```go
module petid

go 1.25
```

**petid/README.md:**
- Заменить все "ЗооБаза" → "PetID"
- Заменить все "PetBase" → "PetID"
- Обновить описание: "Единый реестр животных + справочник пород"

#### 2.3. Обновить ссылки в других файлах

**README.md:**
- `petbase/` → `petid/`
- "ЗооБаза" → "PetID"
- Обновить описание микросервиса

**docs/ROADMAP.md:**
- Уже обновлён ✅

**docs/STRUCTURE.md:**
- Обновить структуру проекта
- Обновить таблицу портов

**run script:**
```bash
# Изменить пути
cd petid/backend
```

#### 2.4. Тестирование после переименования
```bash
# Запустить PetID backend
cd petid/backend
go mod tidy
go run main.go

# Проверить endpoints
curl http://localhost:8100/
curl http://localhost:8100/api/health
curl http://localhost:8100/api/species
curl http://localhost:8100/api/breeds
```

---

### ЭТАП 3: Создание новых таблиц (1 день)

#### 3.1. Создать миграцию
**Файл:** `database/migrations/008_create_petid_tables.sql`

```sql
-- ============================================
-- Миграция 008: Создание таблиц PetID
-- Дата: 2025-12-28
-- Описание: Создание реестра животных, событий и чипов
-- ============================================

-- Таблица: pet_registry
CREATE TABLE IF NOT EXISTS pet_registry (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  species_id INTEGER NOT NULL,
  breed_id INTEGER,
  gender TEXT CHECK(gender IN ('male', 'female', 'unknown')),
  birth_date DATE,
  color TEXT,
  special_marks TEXT,
  photos TEXT,  -- JSON array
  city TEXT,
  status TEXT NOT NULL DEFAULT 'unknown' CHECK(status IN ('shelter', 'home', 'street', 'deceased', 'unknown')),
  responsible_id INTEGER,
  responsible_type TEXT CHECK(responsible_type IN ('owner', 'shelter', 'volunteer', 'clinic')),
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (species_id) REFERENCES species(id),
  FOREIGN KEY (breed_id) REFERENCES breeds(id)
);

CREATE INDEX idx_pet_registry_responsible ON pet_registry(responsible_id, responsible_type);
CREATE INDEX idx_pet_registry_status ON pet_registry(status);
CREATE INDEX idx_pet_registry_city ON pet_registry(city);

-- Таблица: pet_events
CREATE TABLE IF NOT EXISTS pet_events (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  pet_id INTEGER NOT NULL,
  event_type TEXT NOT NULL CHECK(event_type IN (
    'registration',
    'ownership_change',
    'sterilization',
    'vaccination',
    'lost',
    'found',
    'death'
  )),
  event_date DATE NOT NULL,
  description TEXT,
  confirmed_by_id INTEGER,
  confirmed_by_type TEXT CHECK(confirmed_by_type IN ('clinic', 'shelter', 'owner', 'volunteer')),
  confirmation_required BOOLEAN DEFAULT 0,
  confirmed_at DATETIME,
  metadata TEXT,  -- JSON
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (pet_id) REFERENCES pet_registry(id) ON DELETE CASCADE
);

CREATE INDEX idx_pet_events_pet ON pet_events(pet_id);
CREATE INDEX idx_pet_events_type ON pet_events(event_type);
CREATE INDEX idx_pet_events_date ON pet_events(event_date);

-- Таблица: pet_chips
CREATE TABLE IF NOT EXISTS pet_chips (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  pet_id INTEGER NOT NULL,
  chip_type TEXT NOT NULL CHECK(chip_type IN ('chip', 'tattoo', 'tag')),
  chip_number TEXT NOT NULL UNIQUE,
  registered_by_id INTEGER NOT NULL,
  registered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  is_active BOOLEAN DEFAULT 1,
  FOREIGN KEY (pet_id) REFERENCES pet_registry(id) ON DELETE CASCADE
);

CREATE INDEX idx_pet_chips_number ON pet_chips(chip_number);
CREATE INDEX idx_pet_chips_pet ON pet_chips(pet_id);
```

#### 3.2. Применить миграцию
```bash
sqlite3 database/data.db < database/migrations/008_create_petid_tables.sql

# Проверить создание таблиц
sqlite3 database/data.db "SELECT name FROM sqlite_master WHERE type='table' AND name LIKE 'pet_%';"
```

---

### ЭТАП 4: Миграция данных pets → pet_registry (1 день)

#### 4.1. Создать миграцию данных
**Файл:** `database/migrations/009_migrate_pets_to_registry.sql`

```sql
-- ============================================
-- Миграция 009: Перенос данных pets → pet_registry
-- Дата: 2025-12-28
-- Описание: Миграция существующих питомцев в новый реестр
-- ============================================

-- Шаг 1: Проверить количество записей
SELECT COUNT(*) as total_pets FROM pets;

-- Шаг 2: Мигрировать данные
INSERT INTO pet_registry (
  name,
  species_id,
  breed_id,
  photos,
  responsible_id,
  responsible_type,
  status,
  created_at,
  updated_at
)
SELECT
  p.name,
  COALESCE(
    (SELECT id FROM species WHERE LOWER(name) = LOWER(p.species) LIMIT 1),
    (SELECT id FROM species WHERE name = 'Другое' LIMIT 1)
  ) as species_id,
  NULL as breed_id,
  CASE
    WHEN p.photo IS NOT NULL AND p.photo != '' THEN json_array(p.photo)
    ELSE json_array()
  END as photos,
  p.user_id as responsible_id,
  'owner' as responsible_type,
  'home' as status,
  p.created_at,
  p.created_at as updated_at
FROM pets p;

-- Шаг 3: Создать события регистрации
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
  DATE(pr.created_at) as event_date,
  'Миграция из старой системы' as description,
  pr.responsible_id as confirmed_by_id,
  'owner' as confirmed_by_type,
  pr.created_at as confirmed_at,
  pr.created_at
FROM pet_registry pr
WHERE pr.responsible_type = 'owner';

-- Шаг 4: Проверить результаты
SELECT
  (SELECT COUNT(*) FROM pets) as old_pets_count,
  (SELECT COUNT(*) FROM pet_registry) as new_registry_count,
  (SELECT COUNT(*) FROM pet_events) as events_count;
```

#### 4.2. Применить миграцию
```bash
# Применить миграцию
sqlite3 database/data.db < database/migrations/009_migrate_pets_to_registry.sql

# Проверить данные
sqlite3 database/data.db "SELECT id, name, species_id, responsible_id FROM pet_registry LIMIT 5;"
sqlite3 database/data.db "SELECT id, pet_id, event_type FROM pet_events LIMIT 5;"
```

#### 4.3. Создать таблицу соответствия (временно)
```sql
-- Для обратной совместимости на время миграции
CREATE TABLE IF NOT EXISTS pets_migration_map (
  old_pet_id INTEGER PRIMARY KEY,
  new_pet_id INTEGER NOT NULL,
  FOREIGN KEY (new_pet_id) REFERENCES pet_registry(id)
);

INSERT INTO pets_migration_map (old_pet_id, new_pet_id)
SELECT p.id as old_pet_id, pr.id as new_pet_id
FROM pets p
JOIN pet_registry pr ON pr.name = p.name AND pr.responsible_id = p.user_id;
```

---

### ЭТАП 5: Создание API для PetID (2 дня)

#### 5.1. Создать модели
**Файл:** `petid/backend/models/pet_registry.go`

```go
package models

type PetRegistry struct {
    ID              int      `json:"id"`
    Name            string   `json:"name"`
    SpeciesID       int      `json:"species_id"`
    BreedID         *int     `json:"breed_id"`
    Gender          string   `json:"gender"`
    BirthDate       *string  `json:"birth_date"`
    Color           string   `json:"color"`
    SpecialMarks    string   `json:"special_marks"`
    Photos          string   `json:"photos"` // JSON array
    City            string   `json:"city"`
    Status          string   `json:"status"`
    ResponsibleID   *int     `json:"responsible_id"`
    ResponsibleType string   `json:"responsible_type"`
    CreatedAt       string   `json:"created_at"`
    UpdatedAt       string   `json:"updated_at"`
}

type PetEvent struct {
    ID                   int     `json:"id"`
    PetID                int     `json:"pet_id"`
    EventType            string  `json:"event_type"`
    EventDate            string  `json:"event_date"`
    Description          string  `json:"description"`
    ConfirmedByID        *int    `json:"confirmed_by_id"`
    ConfirmedByType      string  `json:"confirmed_by_type"`
    ConfirmationRequired bool    `json:"confirmation_required"`
    ConfirmedAt          *string `json:"confirmed_at"`
    Metadata             string  `json:"metadata"` // JSON
    CreatedAt            string  `json:"created_at"`
}

type PetChip struct {
    ID             int    `json:"id"`
    PetID          int    `json:"pet_id"`
    ChipType       string `json:"chip_type"`
    ChipNumber     string `json:"chip_number"`
    RegisteredByID int    `json:"registered_by_id"`
    RegisteredAt   string `json:"registered_at"`
    IsActive       bool   `json:"is_active"`
}
```

#### 5.2. Создать handlers
**Файл:** `petid/backend/handlers/registry.go`

```go
package handlers

// GET /api/petid/registry/:id
func GetPetRegistryHandler(w http.ResponseWriter, r *http.Request) {
    // Реализация
}

// POST /api/petid/registry
func CreatePetRegistryHandler(w http.ResponseWriter, r *http.Request) {
    // Проверка прав (clinic, shelter, volunteer)
    // Создание записи
}

// PUT /api/petid/registry/:id
func UpdatePetRegistryHandler(w http.ResponseWriter, r *http.Request) {
    // Проверка прав (responsible, clinic, shelter)
    // Обновление записи
}

// GET /api/petid/registry/search
func SearchPetRegistryHandler(w http.ResponseWriter, r *http.Request) {
    // Поиск по имени, чипу, городу
}
```

**Файл:** `petid/backend/handlers/events.go`

```go
package handlers

// GET /api/petid/events/:pet_id
func GetPetEventsHandler(w http.ResponseWriter, r *http.Request) {
    // Получить историю событий
}

// POST /api/petid/events
func CreatePetEventHandler(w http.ResponseWriter, r *http.Request) {
    // Создать событие
}

// PUT /api/petid/events/:id/confirm
func ConfirmPetEventHandler(w http.ResponseWriter, r *http.Request) {
    // Подтвердить событие (для клиник)
}
```

**Файл:** `petid/backend/handlers/chips.go`

```go
package handlers

// POST /api/petid/chips
func RegisterChipHandler(w http.ResponseWriter, r *http.Request) {
    // Зарегистрировать чип (только клиники)
}

// GET /api/petid/chips/search/:number
func SearchChipHandler(w http.ResponseWriter, r *http.Request) {
    // Поиск по номеру чипа
}
```

#### 5.3. Обновить main.go
```go
// Добавить новые routes
http.HandleFunc("/api/petid/registry", enableCORS(handlers.PetRegistryHandler))
http.HandleFunc("/api/petid/registry/", enableCORS(handlers.PetRegistryDetailHandler))
http.HandleFunc("/api/petid/events", enableCORS(handlers.PetEventsHandler))
http.HandleFunc("/api/petid/events/", enableCORS(handlers.PetEventDetailHandler))
http.HandleFunc("/api/petid/chips", enableCORS(handlers.ChipsHandler))
http.HandleFunc("/api/petid/chips/search/", enableCORS(handlers.SearchChipHandler))
```

---

### ЭТАП 6: Обновление Main API (1 день)

#### 6.1. Создать proxy endpoints в Main
**Файл:** `main/backend/handlers/petid_proxy.go`

```go
package handlers

import (
    "io"
    "net/http"
)

const petidAPIURL = "http://localhost:8100"

// Proxy для /api/pets/* → PetID API
func PetIDProxyHandler(w http.ResponseWriter, r *http.Request) {
    // Перенаправить запрос на PetID API
    targetURL := petidAPIURL + "/api/petid" + r.URL.Path[len("/api/pets"):]
    
    proxyReq, err := http.NewRequest(r.Method, targetURL, r.Body)
    if err != nil {
        http.Error(w, "Proxy error", http.StatusInternalServerError)
        return
    }
    
    // Копировать headers
    for key, values := range r.Header {
        for _, value := range values {
            proxyReq.Header.Add(key, value)
        }
    }
    
    client := &http.Client{}
    resp, err := client.Do(proxyReq)
    if err != nil {
        http.Error(w, "Proxy error", http.StatusBadGateway)
        return
    }
    defer resp.Body.Close()
    
    // Копировать response
    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }
    w.WriteHeader(resp.StatusCode)
    io.Copy(w, resp.Body)
}
```

#### 6.2. Обновить routes в main.go
```go
// Старые endpoints теперь проксируют на PetID
http.HandleFunc("/api/pets/user/", enableCORS(authMiddleware(handlers.PetIDProxyHandler)))
http.HandleFunc("/api/pets", enableCORS(authMiddleware(handlers.PetIDProxyHandler)))
http.HandleFunc("/api/pets/", enableCORS(authMiddleware(handlers.PetIDProxyHandler)))
```

---

### ЭТАП 7: Обновление Frontend (2 дня)

#### 7.1. Обновить API клиент
**Файл:** `main/frontend/lib/api.ts`

```typescript
// Обновить endpoints
export const getPets = async (userId: number) => {
  const response = await fetch(`${API_URL}/api/petid/registry?responsible_id=${userId}`, {
    credentials: 'include',
  });
  return response.json();
};

export const createPet = async (petData: CreatePetRequest) => {
  const response = await fetch(`${API_URL}/api/petid/registry`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify(petData),
  });
  return response.json();
};
```

#### 7.2. Создать страницу карточки PetID
**Файл:** `petid/frontend/app/petid/[id]/page.tsx`

```typescript
'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';

export default function PetIDPage() {
  const params = useParams();
  const [activeTab, setActiveTab] = useState<'info' | 'events' | 'posts'>('info');
  const [pet, setPet] = useState(null);
  
  // Вкладки: О питомце, История событий, Посты о питомце
  
  return (
    <div>
      {/* UI карточки PetID */}
    </div>
  );
}
```

---

### ЭТАП 8: Удаление старой таблицы pets (1 день)

#### 8.1. Проверка перед удалением
```bash
# Убедиться, что все данные мигрированы
sqlite3 database/data.db "
SELECT
  (SELECT COUNT(*) FROM pets) as old_count,
  (SELECT COUNT(*) FROM pet_registry WHERE responsible_type = 'owner') as new_count;
"

# Убедиться, что все endpoints работают
curl http://localhost:8000/api/pets/user/1
curl http://localhost:8100/api/petid/registry/1
```

#### 8.2. Создать финальную миграцию
**Файл:** `database/migrations/010_drop_pets_table.sql`

```sql
-- ============================================
-- Миграция 010: Удаление старой таблицы pets
-- Дата: 2025-12-28
-- Описание: Удаление таблицы pets после успешной миграции
-- ============================================

-- Проверка перед удалением
SELECT 'Pets count:' as check_name, COUNT(*) as count FROM pets
UNION ALL
SELECT 'Registry count:', COUNT(*) FROM pet_registry
UNION ALL
SELECT 'Events count:', COUNT(*) FROM pet_events;

-- Удалить таблицу соответствия
DROP TABLE IF EXISTS pets_migration_map;

-- Удалить старую таблицу
DROP TABLE IF EXISTS pets;

-- Проверка после удаления
SELECT name FROM sqlite_master WHERE type='table' AND name='pets';
```

#### 8.3. Применить миграцию
```bash
# Последняя резервная копия перед удалением
cp database/data.db database/backups/data_before_drop_pets_$(date +%Y%m%d_%H%M%S).db

# Удалить таблицу
sqlite3 database/data.db < database/migrations/010_drop_pets_table.sql

# Проверить
sqlite3 database/data.db "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;"
```

#### 8.4. Удалить старые файлы
```bash
# Удалить старые handlers и models
rm main/backend/handlers/pets.go
rm main/backend/models/pet.go

# Обновить imports в main.go
```

---

### ЭТАП 9: Тестирование (2 дня)

#### 9.1. Функциональное тестирование
- [ ] Создание карточки PetID
- [ ] Просмотр карточки PetID
- [ ] Редактирование карточки
- [ ] Создание события
- [ ] Просмотр истории событий
- [ ] Регистрация чипа
- [ ] Поиск по чипу
- [ ] Прикрепление PetID к посту
- [ ] Отображение постов в карточке PetID

#### 9.2. Тестирование прав доступа
- [ ] Обычный пользователь может создать своего питомца
- [ ] Волонтёр может создать карточку (через модерацию)
- [ ] Клиника может создать карточку и регистрировать чипы
- [ ] Приют может создать карточку
- [ ] Только ответственный может редактировать карточку
- [ ] Только клиника может подтверждать медицинские события

#### 9.3. Тестирование интеграции
- [ ] Main → PetID API работает
- [ ] Посты с прикреплёнными PetID отображаются
- [ ] Профиль показывает питомцев из PetID
- [ ] SSO работает между Main и PetID

---

### ЭТАП 10: Документация и релиз (1 день)

#### 10.1. Обновить документацию
- [x] README.md - обновлён ✅
- [x] docs/ROADMAP.md - обновлён ✅
- [ ] docs/STRUCTURE.md - обновить
- [ ] petid/README.md - обновить с новой функциональностью
- [ ] Создать docs/PETID_API.md - документация API

#### 10.2. Создать CHANGELOG
```markdown
## [0.4.0] - 2025-12-XX

### Добавлено
- 🐾 **PetID - Единый реестр животных**
  - Реестр животных (pet_registry)
  - История событий (pet_events)
  - Чипы и метки (pet_chips)
  - API для работы с реестром
  - Карточка PetID с 3 вкладками

### Изменено
- Переименован микросервис PetBase → PetID
- Расширена функциональность: справочник + реестр

### Удалено
- Таблица pets из Main (мигрирована в pet_registry)
- Старые handlers pets.go и models pet.go

### Миграция
- Все существующие питомцы автоматически мигрированы в PetID
- Созданы события регистрации для всех животных
```

#### 10.3. Git commit и tag
```bash
git add -A
git commit -m "Release v0.4.0: PetID - Единый реестр животных

Новое:
- PetID микросервис (расширение PetBase)
- Реестр животных (pet_registry)
- История событий (pet_events)
- Чипы и метки (pet_chips)
- API для работы с реестром
- Карточка PetID с 3 вкладками

Изменено:
- PetBase → PetID (переименование)
- Миграция данных pets → pet_registry

Удалено:
- Таблица pets из Main
- Старые handlers и models для pets"

git tag -a v0.4.0 -m "Release v0.4.0: PetID - Единый реестр животных"
```

---

## ✅ CHECKLIST МИГРАЦИИ

### Подготовка
- [ ] Резервная копия базы данных создана
- [ ] Git ветка создана
- [ ] Документация текущего состояния

### Переименование
- [ ] Директория petbase → petid
- [ ] Обновлены все файлы в petid/
- [ ] Обновлены ссылки в README.md
- [ ] Обновлены ссылки в docs/
- [ ] Обновлён run script
- [ ] Тестирование после переименования

### База данных
- [ ] Миграция 008: создание таблиц PetID
- [ ] Миграция 009: перенос данных pets → pet_registry
- [ ] Миграция 010: удаление таблицы pets
- [ ] Все данные сохранены
- [ ] События регистрации созданы

### Backend
- [ ] Модели созданы (pet_registry.go, pet_event.go, pet_chip.go)
- [ ] Handlers созданы (registry.go, events.go, chips.go)
- [ ] Routes добавлены в main.go
- [ ] Proxy endpoints в Main
- [ ] Middleware проверки прав

### Frontend
- [ ] API клиент обновлён
- [ ] Страница карточки PetID создана
- [ ] Компоненты для вкладок
- [ ] Интеграция с постами

### Тестирование
- [ ] Функциональное тестирование
- [ ] Тестирование прав доступа
- [ ] Тестирование интеграции
- [ ] Нагрузочное тестирование

### Документация
- [ ] README.md обновлён
- [ ] docs/ROADMAP.md обновлён
- [ ] docs/STRUCTURE.md обновлён
- [ ] petid/README.md обновлён
- [ ] docs/PETID_API.md создан
- [ ] CHANGELOG.md обновлён

### Релиз
- [ ] Git commit создан
- [ ] Git tag v0.4.0 создан
- [ ] Все файлы закоммичены

---

## 🚨 ОТКАТ (если что-то пошло не так)

### Откат базы данных
```bash
# Восстановить из резервной копии
cp database/backups/data_before_petid_migration_YYYYMMDD_HHMMSS.db database/data.db
```

### Откат кода
```bash
# Вернуться к предыдущему коммиту
git checkout main
git branch -D feature/petbase-to-petid-migration
```

---

## 📊 ОЖИДАЕМЫЕ РЕЗУЛЬТАТЫ

### После миграции:
- ✅ Все существующие питомцы сохранены
- ✅ PetID микросервис работает
- ✅ Старые API endpoints работают через proxy
- ✅ Новые API endpoints доступны
- ✅ Frontend отображает данные из PetID
- ✅ Права доступа работают корректно
- ✅ Документация актуальна

### Метрики:
- 0 потерянных записей
- 100% успешных миграций
- <2s время ответа API
- 0 критичных багов

---

**Последнее обновление:** 28 декабря 2025  
**Статус:** Готов к выполнению  
**Ответственный:** Команда разработки
