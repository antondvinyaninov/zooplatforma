# PetBase Service - База данных животных

**Порт:** 8100  
**Назначение:** Центральная база данных всех животных платформы

---

## Описание

PetBase Service - это главный сервис для работы с данными о животных:
- Справочники (виды, породы, карточки)
- Все питомцы платформы
- Медицинские записи
- Вакцинации
- Чипы и метки
- Родословная
- История событий

**ВАЖНО:** Все остальные сервисы (Main, Clinic, Shelter, Owner) используют PetBase для работы с животными!

---

## База данных

`petbase/database/petbase.db`

### Таблицы

```sql
-- Виды животных
CREATE TABLE species (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    icon TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Породы
CREATE TABLE breeds (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    species_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    size TEXT,
    temperament TEXT,
    lifespan TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (species_id) REFERENCES species(id)
);

-- Информационные карточки о породах
CREATE TABLE pet_cards (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    breed_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT,
    image TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (breed_id) REFERENCES breeds(id)
);

-- ВСЕ ПИТОМЦЫ (главная таблица!)
CREATE TABLE pets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    organization_id INTEGER,
    name TEXT NOT NULL,
    species TEXT NOT NULL,
    breed TEXT,
    gender TEXT,
    birth_date DATE,
    color TEXT,
    size TEXT,
    photo TEXT,
    status TEXT DEFAULT 'home',
    city TEXT,
    region TEXT,
    urgent BOOLEAN DEFAULT 0,
    story TEXT,
    curator_id INTEGER,
    chip_number TEXT,
    tag_number TEXT,
    sterilization_date DATE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Медицинские записи
CREATE TABLE pet_medical_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER NOT NULL,
    clinic_id INTEGER,
    date DATE NOT NULL,
    diagnosis TEXT,
    treatment TEXT,
    notes TEXT,
    veterinarian TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pet_id) REFERENCES pets(id)
);

-- Вакцинации
CREATE TABLE pet_vaccinations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER NOT NULL,
    vaccine_name TEXT NOT NULL,
    date DATE NOT NULL,
    next_date DATE,
    clinic_id INTEGER,
    veterinarian TEXT,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pet_id) REFERENCES pets(id)
);

-- Чипы
CREATE TABLE pet_chips (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER NOT NULL,
    chip_number TEXT UNIQUE NOT NULL,
    implant_date DATE,
    clinic_id INTEGER,
    location TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pet_id) REFERENCES pets(id)
);

-- История событий
CREATE TABLE pet_events (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER NOT NULL,
    event_type TEXT NOT NULL,
    event_date DATE NOT NULL,
    description TEXT,
    clinic_id INTEGER,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pet_id) REFERENCES pets(id)
);

-- Родословная
CREATE TABLE pet_pedigree (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pet_id INTEGER NOT NULL,
    parent_id INTEGER NOT NULL,
    parent_type TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pet_id) REFERENCES pets(id),
    FOREIGN KEY (parent_id) REFERENCES pets(id)
);
```

---

## API Endpoints

### Справочники (публичные)

#### GET /api/species
Получить все виды животных

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Собака",
      "description": "Домашнее животное семейства псовых",
      "icon": "🐕"
    },
    {
      "id": 2,
      "name": "Кошка",
      "description": "Домашнее животное семейства кошачьих",
      "icon": "🐈"
    }
  ]
}
```

#### GET /api/breeds
Получить все породы

**Query параметры:**
- `species_id` - фильтр по виду

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "species_id": 1,
      "name": "Лабрадор",
      "description": "Дружелюбная порода",
      "size": "large",
      "temperament": "Дружелюбный, активный",
      "lifespan": "10-12 лет"
    }
  ]
}
```

#### GET /api/breeds/species/:id
Получить породы конкретного вида

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Лабрадор"
    },
    {
      "id": 2,
      "name": "Немецкая овчарка"
    }
  ]
}
```

### Питомцы (главная таблица)

#### GET /api/pets
Получить всех питомцев (с фильтрами)

**Headers:**
```
Authorization: Bearer <token>
```

**Query параметры:**
- `user_id` - питомцы пользователя
- `organization_id` - питомцы организации
- `status` - фильтр по статусу (home, looking_for_home, lost, found)
- `species` - фильтр по виду
- `city` - фильтр по городу
- `urgent` - только срочные

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 5,
      "name": "Барсик",
      "species": "cat",
      "breed": "Британская короткошерстная",
      "gender": "male",
      "birth_date": "2020-05-15",
      "color": "Серый",
      "photo": "/uploads/pets/1.jpg",
      "status": "home",
      "chip_number": "123456789",
      "created_at": "2025-01-17T10:00:00Z"
    }
  ]
}
```

#### GET /api/pets/:id
Получить конкретного питомца

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 5,
    "name": "Барсик",
    "species": "cat",
    "breed": "Британская короткошерстная",
    "gender": "male",
    "birth_date": "2020-05-15",
    "color": "Серый",
    "size": "medium",
    "photo": "/uploads/pets/1.jpg",
    "status": "home",
    "city": "Москва",
    "region": "Московская область",
    "story": "Добрый и ласковый кот",
    "chip_number": "123456789",
    "tag_number": "TAG001",
    "sterilization_date": "2021-01-10",
    "created_at": "2025-01-17T10:00:00Z"
  }
}
```

#### POST /api/pets
Создать питомца

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "name": "Барсик",
  "species": "cat",
  "breed": "Британская короткошерстная",
  "gender": "male",
  "birth_date": "2020-05-15",
  "color": "Серый",
  "size": "medium",
  "status": "home"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "user_id": 5,
    "name": "Барсик",
    "species": "cat",
    "breed": "Британская короткошерстная"
  }
}
```

#### PUT /api/pets/:id
Обновить питомца

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "name": "Барсик Обновленный",
  "color": "Темно-серый"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Барсик Обновленный",
    "color": "Темно-серый"
  }
}
```

#### DELETE /api/pets/:id
Удалить питомца

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "message": "Pet deleted successfully"
}
```

### Медицинские записи

#### GET /api/pets/:id/medical
Получить медкарту питомца

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "pet_id": 1,
      "clinic_id": 5,
      "date": "2025-01-15",
      "diagnosis": "Здоров",
      "treatment": "Вакцинация",
      "notes": "Плановый осмотр",
      "veterinarian": "Иванов И.И.",
      "created_at": "2025-01-15T10:00:00Z"
    }
  ]
}
```

#### POST /api/pets/:id/medical
Добавить медицинскую запись

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "clinic_id": 5,
  "date": "2025-01-15",
  "diagnosis": "Здоров",
  "treatment": "Вакцинация",
  "notes": "Плановый осмотр",
  "veterinarian": "Иванов И.И."
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "pet_id": 1,
    "clinic_id": 5,
    "date": "2025-01-15"
  }
}
```

### Вакцинации

#### GET /api/pets/:id/vaccinations
Получить вакцинации питомца

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "pet_id": 1,
      "vaccine_name": "Бешенство",
      "date": "2025-01-10",
      "next_date": "2026-01-10",
      "clinic_id": 5,
      "veterinarian": "Иванов И.И."
    }
  ]
}
```

#### POST /api/pets/:id/vaccinations
Добавить вакцинацию

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "vaccine_name": "Бешенство",
  "date": "2025-01-10",
  "next_date": "2026-01-10",
  "clinic_id": 5,
  "veterinarian": "Иванов И.И.",
  "notes": "Первая вакцинация"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "pet_id": 1,
    "vaccine_name": "Бешенство",
    "date": "2025-01-10"
  }
}
```

### Чипы

#### GET /api/chips/:number
Найти питомца по номеру чипа

**Response:**
```json
{
  "success": true,
  "data": {
    "pet": {
      "id": 1,
      "name": "Барсик",
      "species": "cat",
      "photo": "/uploads/pets/1.jpg"
    },
    "chip": {
      "chip_number": "123456789",
      "implant_date": "2020-06-01",
      "clinic_id": 5
    }
  }
}
```

#### POST /api/chips
Зарегистрировать чип

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "pet_id": 1,
  "chip_number": "123456789",
  "implant_date": "2020-06-01",
  "clinic_id": 5,
  "location": "Между лопаток"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "pet_id": 1,
    "chip_number": "123456789"
  }
}
```

### История событий

#### GET /api/pets/:id/events
Получить историю питомца

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "pet_id": 1,
      "event_type": "registration",
      "event_date": "2025-01-17",
      "description": "Регистрация в клинике",
      "clinic_id": 5,
      "notes": "Первичная регистрация"
    },
    {
      "id": 2,
      "pet_id": 1,
      "event_type": "vaccination",
      "event_date": "2025-01-10",
      "description": "Вакцинация от бешенства"
    }
  ]
}
```

#### POST /api/pet-events
Добавить событие

**Request:**
```json
{
  "pet_id": 1,
  "event_type": "medical_visit",
  "event_date": "2025-01-17",
  "description": "Плановый осмотр",
  "clinic_id": 5,
  "notes": "Все в порядке"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "pet_id": 1,
    "event_type": "medical_visit"
  }
}
```

### Родословная

#### GET /api/pets/:id/parents
Получить родителей питомца

**Response:**
```json
{
  "success": true,
  "data": {
    "mother": {
      "id": 10,
      "name": "Мурка",
      "photo": "/uploads/pets/10.jpg"
    },
    "father": {
      "id": 11,
      "name": "Барон",
      "photo": "/uploads/pets/11.jpg"
    }
  }
}
```

#### POST /api/pets/:id/parents
Указать родителей

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "mother_id": 10,
  "father_id": 11
}
```

**Response:**
```json
{
  "success": true,
  "message": "Parents added successfully"
}
```

---

## Использование в других сервисах

### Main Service
```go
// Получить питомцев для поста
petbaseClient := clients.NewPetBaseClient("http://localhost:8100")
pets, err := petbaseClient.GetPets(petIDs, token)
```

### Clinic Service
```go
// Добавить медицинскую запись
petbaseClient := clients.NewPetBaseClient("http://localhost:8100")
err := petbaseClient.AddMedicalRecord(petID, record, token)
```

### Owner Service
```go
// Получить моих питомцев
petbaseClient := clients.NewPetBaseClient("http://localhost:8100")
pets, err := petbaseClient.GetUserPets(userID, token)
```

---

## Конфигурация

### .env
```bash
PORT=8100
DATABASE_URL=./database/petbase.db
AUTH_SERVICE_URL=http://localhost:7000
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:4100
```

---

## Запуск

```bash
cd petbase/backend
air
```

---

## Swagger документация

```
http://localhost:8100/swagger/index.html
```

---

**Дата обновления:** 17 января 2025
