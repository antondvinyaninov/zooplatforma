# Shelter Service - Кабинет приюта

**Порт Backend:** 8200  
**Порт Frontend:** 5100  
**Назначение:** Управление приютом для животных

---

## Описание

Shelter Service - это кабинет для приютов:
- Управление приютом и командой
- Животные в приюте
- Заявки на пристройство
- Волонтеры
- Пожертвования
- Расходы приюта

**Multi-tenancy:** Заголовок `X-Shelter-ID` для изоляции данных между приютами

---

## База данных

`shelter/database/shelter.db`

### Таблицы

```sql
-- Животные в приюте
CREATE TABLE shelter_animals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    shelter_id INTEGER NOT NULL,
    pet_id INTEGER NOT NULL,
    arrival_date DATE NOT NULL,
    arrival_reason TEXT,
    status TEXT DEFAULT 'available',
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(shelter_id, pet_id)
);

-- Заявки на пристройство
CREATE TABLE shelter_adoptions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    shelter_id INTEGER NOT NULL,
    pet_id INTEGER NOT NULL,
    applicant_id INTEGER NOT NULL,
    status TEXT DEFAULT 'pending',
    application_text TEXT,
    interview_date DATE,
    decision_notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Волонтеры приюта
CREATE TABLE shelter_volunteers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    shelter_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    role TEXT DEFAULT 'volunteer',
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(shelter_id, user_id)
);

-- Пожертвования
CREATE TABLE shelter_donations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    shelter_id INTEGER NOT NULL,
    donor_id INTEGER,
    donor_name TEXT,
    amount DECIMAL(10,2) NOT NULL,
    type TEXT DEFAULT 'money',
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Расходы приюта
CREATE TABLE shelter_expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    shelter_id INTEGER NOT NULL,
    category TEXT NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    description TEXT,
    date DATE NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## API Endpoints

### Приюты

#### GET /api/my-shelters
Получить список приютов пользователя

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
      "name": "Приют Добрые руки",
      "role": "owner",
      "position": "Директор"
    }
  ]
}
```

#### POST /api/shelters
Создать приют

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "name": "Приют Добрые руки",
  "type": "shelter",
  "description": "Приют для бездомных животных",
  "address_city": "Москва",
  "address_street": "ул. Ленина, 20",
  "phone": "+7 (495) 123-45-67",
  "email": "info@shelter.ru"
}
```

#### GET /api/organization
Получить текущий приют

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Приют Добрые руки",
    "description": "Приют для бездомных животных",
    "logo": "/uploads/organizations/1/logo.jpg",
    "address_city": "Москва",
    "phone": "+7 (495) 123-45-67",
    "animals_count": 50,
    "volunteers_count": 10
  }
}
```

#### PUT /api/organization
Обновить приют

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Request:**
```json
{
  "name": "Приют Добрые руки 2.0",
  "description": "Обновленное описание"
}
```

### Животные

#### GET /api/my-animals
Получить животных приюта

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Query параметры:**
- `status` - фильтр по статусу (available, adopted, reserved)
- `species` - фильтр по виду
- `search` - поиск по имени

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "pet_id": 10,
      "pet_name": "Шарик",
      "pet_species": "dog",
      "pet_photo": "/uploads/pets/10.jpg",
      "arrival_date": "2025-01-10",
      "arrival_reason": "Найден на улице",
      "status": "available",
      "days_in_shelter": 7
    }
  ]
}
```

#### GET /api/animals/:id
Получить информацию о животном

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "pet": {
      "id": 10,
      "name": "Шарик",
      "species": "dog",
      "breed": "Дворняга",
      "birth_date": "2023-05-15",
      "photo": "/uploads/pets/10.jpg",
      "story": "Найден на улице"
    },
    "arrival_date": "2025-01-10",
    "arrival_reason": "Найден на улице",
    "status": "available",
    "notes": "Дружелюбный, любит детей",
    "days_in_shelter": 7
  }
}
```

#### POST /api/animals
Добавить животное в приют

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Request:**
```json
{
  "pet_id": 10,
  "arrival_date": "2025-01-10",
  "arrival_reason": "Найден на улице",
  "notes": "Дружелюбный"
}
```

#### PUT /api/animals/:id
Обновить информацию о животном

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Request:**
```json
{
  "status": "adopted",
  "notes": "Пристроен в хорошие руки"
}
```

#### DELETE /api/animals/:id
Удалить животное из приюта

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

### Заявки на пристройство

#### GET /api/adoptions
Получить заявки на пристройство

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Query параметры:**
- `status` - фильтр по статусу (pending, approved, rejected)
- `pet_id` - фильтр по животному

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "pet_id": 10,
      "pet_name": "Шарик",
      "pet_photo": "/uploads/pets/10.jpg",
      "applicant_id": 5,
      "applicant_name": "Иван Иванов",
      "status": "pending",
      "created_at": "2025-01-17T10:00:00Z"
    }
  ]
}
```

#### GET /api/adoptions/:id
Получить информацию о заявке

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "pet": {
      "id": 10,
      "name": "Шарик",
      "photo": "/uploads/pets/10.jpg"
    },
    "applicant": {
      "id": 5,
      "name": "Иван Иванов",
      "phone": "+7 (999) 123-45-67",
      "email": "ivan@example.com"
    },
    "status": "pending",
    "application_text": "Хочу взять Шарика...",
    "interview_date": "2025-01-20",
    "created_at": "2025-01-17T10:00:00Z"
  }
}
```

#### POST /api/adoptions
Создать заявку

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Request:**
```json
{
  "pet_id": 10,
  "application_text": "Хочу взять Шарика в хорошие руки..."
}
```

#### PUT /api/adoptions/:id
Обновить заявку

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Request:**
```json
{
  "status": "approved",
  "interview_date": "2025-01-20",
  "decision_notes": "Одобрено после собеседования"
}
```

#### POST /api/adoptions/:id/approve
Одобрить заявку

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Request:**
```json
{
  "decision_notes": "Одобрено"
}
```

#### POST /api/adoptions/:id/reject
Отклонить заявку

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Request:**
```json
{
  "decision_notes": "Не подходит"
}
```

### Волонтеры

#### GET /api/volunteers
Получить волонтеров приюта

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "user_id": 10,
      "name": "Петр Петров",
      "avatar": "/uploads/users/10/avatar.jpg",
      "role": "volunteer",
      "joined_at": "2025-01-01T10:00:00Z"
    }
  ]
}
```

#### POST /api/volunteers
Добавить волонтера

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Request:**
```json
{
  "user_id": 10,
  "role": "volunteer"
}
```

#### DELETE /api/volunteers/:id
Удалить волонтера

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

### Пожертвования

#### GET /api/donations
Получить пожертвования

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Query параметры:**
- `type` - фильтр по типу (money, food, supplies)
- `date_from` - с даты
- `date_to` - по дату

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "donor_id": 5,
      "donor_name": "Иван Иванов",
      "amount": 5000,
      "type": "money",
      "description": "На корм",
      "created_at": "2025-01-17T10:00:00Z"
    }
  ]
}
```

#### POST /api/donations
Добавить пожертвование

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Request:**
```json
{
  "donor_id": 5,
  "donor_name": "Иван Иванов",
  "amount": 5000,
  "type": "money",
  "description": "На корм"
}
```

#### GET /api/donations/stats
Статистика пожертвований

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_amount": 150000,
    "this_month": 25000,
    "donors_count": 50,
    "by_type": {
      "money": 100000,
      "food": 30000,
      "supplies": 20000
    }
  }
}
```

### Расходы

#### GET /api/expenses
Получить расходы приюта

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Query параметры:**
- `category` - фильтр по категории
- `date_from` - с даты
- `date_to` - по дату

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "category": "food",
      "amount": 10000,
      "description": "Корм для собак",
      "date": "2025-01-15",
      "created_at": "2025-01-15T10:00:00Z"
    }
  ]
}
```

#### POST /api/expenses
Добавить расход

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Request:**
```json
{
  "category": "food",
  "amount": 10000,
  "description": "Корм для собак",
  "date": "2025-01-15"
}
```

#### PUT /api/expenses/:id
Обновить расход

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Request:**
```json
{
  "amount": 12000,
  "description": "Корм для собак (обновлено)"
}
```

#### DELETE /api/expenses/:id
Удалить расход

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

#### GET /api/expenses/stats
Статистика расходов

**Headers:**
```
Authorization: Bearer <token>
X-Shelter-ID: 1
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_amount": 120000,
    "this_month": 20000,
    "by_category": {
      "food": 60000,
      "medical": 40000,
      "supplies": 20000
    }
  }
}
```

---

## Зависимости от других сервисов

- **Auth Service** - авторизация
- **PetBase Service** - данные животных
- **Main Service** - данные организации, участники

---

## Конфигурация

### .env
```bash
PORT=8200
DATABASE_URL=./database/shelter.db
AUTH_SERVICE_URL=http://localhost:7000
PETBASE_SERVICE_URL=http://localhost:8100
MAIN_SERVICE_URL=http://localhost:8000
ALLOWED_ORIGINS=http://localhost:5100
```

---

## Запуск

```bash
cd shelter/backend
air
```

Frontend:
```bash
cd shelter/frontend
npm run dev
```

---

## Swagger документация

```
http://localhost:8200/swagger/index.html
```

---

**Дата обновления:** 17 января 2025
