# Clinic Service - Кабинет ветеринарной клиники

**Порт Backend:** 8600  
**Порт Frontend:** 6300  
**Назначение:** Управление ветеринарной клиникой

---

## Описание

Clinic Service - это кабинет для ветеринарных клиник:
- Управление клиникой и командой
- Записи на прием
- Пациенты клиники
- Расписание работы
- Услуги и прайс-лист
- Счета и оплата

**Multi-tenancy:** Заголовок `X-Clinic-ID` для изоляции данных между клиниками

---

## База данных

`clinic/database/clinic.db`

### Таблицы

```sql
-- Записи на прием
CREATE TABLE clinic_appointments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    clinic_id INTEGER NOT NULL,
    pet_id INTEGER NOT NULL,
    owner_id INTEGER NOT NULL,
    date DATE NOT NULL,
    time TIME NOT NULL,
    service_id INTEGER,
    veterinarian_id INTEGER,
    status TEXT DEFAULT 'scheduled',
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Пациенты клиники
CREATE TABLE clinic_patients (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    clinic_id INTEGER NOT NULL,
    pet_id INTEGER NOT NULL,
    first_visit DATE,
    last_visit DATE,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(clinic_id, pet_id)
);

-- Расписание клиники
CREATE TABLE clinic_schedules (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    clinic_id INTEGER NOT NULL,
    day_of_week INTEGER NOT NULL,
    open_time TIME NOT NULL,
    close_time TIME NOT NULL,
    is_working BOOLEAN DEFAULT 1
);

-- Услуги клиники
CREATE TABLE clinic_services (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    clinic_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    price DECIMAL(10,2),
    duration INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Счета
CREATE TABLE clinic_invoices (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    clinic_id INTEGER NOT NULL,
    appointment_id INTEGER,
    pet_id INTEGER NOT NULL,
    owner_id INTEGER NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    paid BOOLEAN DEFAULT 0,
    paid_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Сотрудники клиники
CREATE TABLE clinic_staff (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    clinic_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    position TEXT NOT NULL,
    specialization TEXT,
    schedule TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## API Endpoints

### Клиники

#### GET /api/my-clinics
Получить список клиник пользователя

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
      "name": "Ветклиника Айболит",
      "role": "owner",
      "position": "Главный врач"
    }
  ]
}
```

#### POST /api/clinics
Создать клинику

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "name": "Ветклиника Айболит",
  "type": "clinic",
  "description": "Современная ветеринарная клиника",
  "address_city": "Москва",
  "address_street": "ул. Ленина, 10",
  "phone": "+7 (495) 123-45-67",
  "email": "info@aibolit.ru",
  "website": "https://aibolit.ru"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Ветклиника Айболит"
  }
}
```

#### GET /api/organization
Получить текущую клинику

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Ветклиника Айболит",
    "description": "Современная ветеринарная клиника",
    "logo": "/uploads/organizations/1/logo.jpg",
    "address_city": "Москва",
    "address_street": "ул. Ленина, 10",
    "phone": "+7 (495) 123-45-67",
    "email": "info@aibolit.ru",
    "website": "https://aibolit.ru",
    "verified": true
  }
}
```

#### PUT /api/organization
Обновить клинику

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Request:**
```json
{
  "name": "Ветклиника Айболит Плюс",
  "description": "Обновленное описание",
  "phone": "+7 (495) 999-88-77"
}
```

### Команда

#### GET /api/members
Получить участников клиники

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "user_id": 5,
      "name": "Иван Иванов",
      "avatar": "/uploads/users/5/avatar.jpg",
      "role": "owner",
      "position": "Главный врач",
      "joined_at": "2025-01-01T10:00:00Z"
    }
  ]
}
```

#### POST /api/members
Добавить участника

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Request:**
```json
{
  "user_id": 10,
  "role": "admin",
  "position": "Ветеринар"
}
```

#### PUT /api/members
Обновить участника

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Request:**
```json
{
  "user_id": 10,
  "role": "moderator",
  "position": "Старший ветеринар"
}
```

#### DELETE /api/members
Удалить участника

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Request:**
```json
{
  "user_id": 10
}
```

### Пациенты

#### GET /api/patients
Получить пациентов клиники

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Query параметры:**
- `search` - поиск по имени питомца/владельца
- `page` - номер страницы
- `limit` - количество на странице

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "pet_id": 10,
      "pet_name": "Барсик",
      "pet_species": "cat",
      "pet_photo": "/uploads/pets/10.jpg",
      "owner_id": 5,
      "owner_name": "Иван Иванов",
      "first_visit": "2025-01-10",
      "last_visit": "2025-01-17",
      "visits_count": 5
    }
  ]
}
```

#### GET /api/patients/:id
Получить информацию о пациенте

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "pet": {
      "id": 10,
      "name": "Барсик",
      "species": "cat",
      "breed": "Британская",
      "birth_date": "2020-05-15",
      "photo": "/uploads/pets/10.jpg"
    },
    "owner": {
      "id": 5,
      "name": "Иван Иванов",
      "phone": "+7 (999) 123-45-67"
    },
    "first_visit": "2025-01-10",
    "last_visit": "2025-01-17",
    "visits_count": 5,
    "notes": "Аллергия на курицу"
  }
}
```

#### POST /api/patients
Добавить пациента

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Request:**
```json
{
  "pet_id": 10,
  "notes": "Аллергия на курицу"
}
```

### Записи на прием

#### GET /api/appointments
Получить записи клиники

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Query параметры:**
- `date` - фильтр по дате
- `status` - фильтр по статусу (scheduled, completed, cancelled)
- `veterinarian_id` - фильтр по ветеринару

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "pet_id": 10,
      "pet_name": "Барсик",
      "owner_id": 5,
      "owner_name": "Иван Иванов",
      "date": "2025-01-20",
      "time": "10:00",
      "service": "Осмотр",
      "veterinarian": "Петров П.П.",
      "status": "scheduled"
    }
  ]
}
```

#### GET /api/appointments/:id
Получить информацию о записи

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "pet": {
      "id": 10,
      "name": "Барсик",
      "photo": "/uploads/pets/10.jpg"
    },
    "owner": {
      "id": 5,
      "name": "Иван Иванов",
      "phone": "+7 (999) 123-45-67"
    },
    "date": "2025-01-20",
    "time": "10:00",
    "service": {
      "id": 1,
      "name": "Осмотр",
      "price": 1000
    },
    "veterinarian": {
      "id": 3,
      "name": "Петров П.П."
    },
    "status": "scheduled",
    "notes": "Плановый осмотр"
  }
}
```

#### POST /api/appointments
Создать запись

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Request:**
```json
{
  "pet_id": 10,
  "owner_id": 5,
  "date": "2025-01-20",
  "time": "10:00",
  "service_id": 1,
  "veterinarian_id": 3,
  "notes": "Плановый осмотр"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "date": "2025-01-20",
    "time": "10:00"
  }
}
```

#### PUT /api/appointments/:id
Обновить запись

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Request:**
```json
{
  "date": "2025-01-21",
  "time": "11:00",
  "status": "completed"
}
```

#### DELETE /api/appointments/:id
Отменить запись

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Response:**
```json
{
  "success": true,
  "message": "Appointment cancelled"
}
```

#### GET /api/appointments/today
Получить записи на сегодня

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "time": "10:00",
      "pet_name": "Барсик",
      "owner_name": "Иван Иванов",
      "service": "Осмотр"
    }
  ]
}
```

### Расписание

#### GET /api/schedules
Получить расписание клиники

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "day_of_week": 1,
      "day_name": "Понедельник",
      "open_time": "09:00",
      "close_time": "18:00",
      "is_working": true
    }
  ]
}
```

#### POST /api/schedules
Создать расписание

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Request:**
```json
{
  "day_of_week": 1,
  "open_time": "09:00",
  "close_time": "18:00",
  "is_working": true
}
```

#### PUT /api/schedules/:id
Обновить расписание

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Request:**
```json
{
  "open_time": "08:00",
  "close_time": "20:00"
}
```

### Услуги

#### GET /api/services
Получить услуги клиники

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Осмотр",
      "description": "Первичный осмотр",
      "price": 1000,
      "duration": 30
    }
  ]
}
```

#### POST /api/services
Добавить услугу

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Request:**
```json
{
  "name": "Вакцинация",
  "description": "Вакцинация от бешенства",
  "price": 1500,
  "duration": 15
}
```

#### PUT /api/services/:id
Обновить услугу

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Request:**
```json
{
  "price": 1200
}
```

#### DELETE /api/services/:id
Удалить услугу

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

### Счета

#### GET /api/invoices
Получить счета клиники

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Query параметры:**
- `paid` - фильтр по оплате (true/false)
- `date_from` - с даты
- `date_to` - по дату

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "appointment_id": 1,
      "pet_name": "Барсик",
      "owner_name": "Иван Иванов",
      "total_amount": 1500,
      "paid": true,
      "paid_at": "2025-01-17T10:00:00Z",
      "created_at": "2025-01-17T09:00:00Z"
    }
  ]
}
```

#### POST /api/invoices
Создать счет

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Request:**
```json
{
  "appointment_id": 1,
  "pet_id": 10,
  "owner_id": 5,
  "items": [
    {
      "service_id": 1,
      "quantity": 1,
      "price": 1000
    }
  ],
  "total_amount": 1000
}
```

#### PUT /api/invoices/:id
Обновить счет (отметить оплаченным)

**Headers:**
```
Authorization: Bearer <token>
X-Clinic-ID: 1
```

**Request:**
```json
{
  "paid": true
}
```

---

## Зависимости от других сервисов

- **Auth Service** - авторизация
- **PetBase Service** - данные питомцев, медицинские записи
- **Main Service** - данные организации, участники

---

## Конфигурация

### .env
```bash
PORT=8600
DATABASE_URL=./database/clinic.db
AUTH_SERVICE_URL=http://localhost:7000
PETBASE_SERVICE_URL=http://localhost:8100
MAIN_SERVICE_URL=http://localhost:8000
ALLOWED_ORIGINS=http://localhost:6300
```

---

## Запуск

```bash
cd clinic/backend
air
```

Frontend:
```bash
cd clinic/frontend
npm run dev
```

---

## Swagger документация

```
http://localhost:8600/swagger/index.html
```

---

**Дата обновления:** 17 января 2025
