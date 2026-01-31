# Volunteer Service - Кабинет зоопомощника

**Порт Backend:** 8500  
**Порт Frontend:** 6200  
**Назначение:** Управление подопечными животными и задачами волонтера

---

## Описание

Volunteer Service - это кабинет для зоопомощников (волонтеров):
- Управление подопечными животными (curator_id)
- Создание и отслеживание задач
- Отчеты о проделанной работе
- Учет времени волонтерской деятельности
- Статистика и аналитика

**База данных:** `volunteer/database/volunteer.db`

---

## API Endpoints

### Подопечные животные

#### GET /api/my-pets
Получить моих подопечных питомцев (где curator_id = user_id)

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
      "id": 15,
      "name": "Шарик",
      "species": "dog",
      "breed": "Дворняга",
      "birth_date": "2022-03-10",
      "gender": "male",
      "photo": "/uploads/pets/15.jpg",
      "status": "shelter",
      "curator_id": 5,
      "location": "Приют 'Добрые руки'",
      "created_at": "2025-01-01T10:00:00Z"
    }
  ]
}
```

#### GET /api/pets/:id
Получить информацию о подопечном питомце

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 15,
    "name": "Шарик",
    "species": "dog",
    "breed": "Дворняга",
    "birth_date": "2022-03-10",
    "gender": "male",
    "color": "Рыжий",
    "size": "medium",
    "photo": "/uploads/pets/15.jpg",
    "status": "shelter",
    "curator_id": 5,
    "curator_name": "Мария Волонтер",
    "location": "Приют 'Добрые руки'",
    "story": "Найден на улице, дружелюбный",
    "health_status": "Здоров",
    "chip_number": "123456789",
    "tag_number": "TAG015",
    "created_at": "2025-01-01T10:00:00Z"
  }
}
```

#### PUT /api/pets/:id
Обновить информацию о подопечном питомце

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "story": "Обновленная история",
  "health_status": "Требуется вакцинация",
  "location": "Новое местоположение"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 15,
    "name": "Шарик",
    "story": "Обновленная история"
  },
  "message": "Pet updated successfully"
}
```

### Задачи

#### GET /api/my-tasks
Получить мои задачи

**Headers:**
```
Authorization: Bearer <token>
```

**Query параметры:**
- `status` - фильтр по статусу (pending, in_progress, completed)
- `pet_id` - фильтр по питомцу
- `date_from` - дата начала периода
- `date_to` - дата окончания периода

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "Выгулять Шарика",
      "description": "Прогулка в парке 30 минут",
      "pet_id": 15,
      "pet_name": "Шарик",
      "status": "pending",
      "priority": "high",
      "due_date": "2025-01-20",
      "created_at": "2025-01-15T10:00:00Z",
      "completed_at": null
    },
    {
      "id": 2,
      "title": "Отвезти на прививку",
      "description": "Вакцинация от бешенства",
      "pet_id": 15,
      "pet_name": "Шарик",
      "status": "completed",
      "priority": "high",
      "due_date": "2025-01-18",
      "created_at": "2025-01-15T10:00:00Z",
      "completed_at": "2025-01-18T14:30:00Z"
    }
  ]
}
```

#### GET /api/tasks/:id
Получить информацию о задаче

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
    "title": "Выгулять Шарика",
    "description": "Прогулка в парке 30 минут",
    "pet_id": 15,
    "pet_name": "Шарик",
    "pet_photo": "/uploads/pets/15.jpg",
    "volunteer_id": 5,
    "volunteer_name": "Мария Волонтер",
    "status": "pending",
    "priority": "high",
    "due_date": "2025-01-20",
    "notes": "Взять поводок и лакомства",
    "created_at": "2025-01-15T10:00:00Z",
    "updated_at": "2025-01-15T10:00:00Z",
    "completed_at": null
  }
}
```

#### POST /api/tasks
Создать новую задачу

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "title": "Выгулять Шарика",
  "description": "Прогулка в парке 30 минут",
  "pet_id": 15,
  "priority": "high",
  "due_date": "2025-01-20",
  "notes": "Взять поводок и лакомства"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Выгулять Шарика",
    "status": "pending"
  },
  "message": "Task created successfully"
}
```

#### PUT /api/tasks/:id
Обновить задачу

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "title": "Выгулять Шарика (обновлено)",
  "status": "in_progress",
  "notes": "Начал выполнение"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Выгулять Шарика (обновлено)",
    "status": "in_progress"
  },
  "message": "Task updated successfully"
}
```

#### POST /api/tasks/:id/complete
Завершить задачу

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "notes": "Задача выполнена успешно"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "status": "completed",
    "completed_at": "2025-01-20T15:30:00Z"
  },
  "message": "Task completed successfully"
}
```

#### DELETE /api/tasks/:id
Удалить задачу

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "message": "Task deleted successfully"
}
```

### Отчеты

#### GET /api/reports
Получить мои отчеты

**Headers:**
```
Authorization: Bearer <token>
```

**Query параметры:**
- `date_from` - дата начала периода
- `date_to` - дата окончания периода
- `pet_id` - фильтр по питомцу

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "title": "Отчет за неделю",
      "description": "Выполнено 5 задач, выгуляно 3 собаки",
      "date": "2025-01-20",
      "volunteer_id": 5,
      "volunteer_name": "Мария Волонтер",
      "tasks_completed": 5,
      "hours_worked": 12,
      "pets_helped": [15, 16, 17],
      "created_at": "2025-01-20T18:00:00Z"
    }
  ]
}
```

#### GET /api/reports/:id
Получить информацию об отчете

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
    "title": "Отчет за неделю",
    "description": "Выполнено 5 задач, выгуляно 3 собаки",
    "date": "2025-01-20",
    "volunteer_id": 5,
    "volunteer_name": "Мария Волонтер",
    "tasks_completed": 5,
    "hours_worked": 12,
    "pets_helped": [
      {
        "id": 15,
        "name": "Шарик",
        "photo": "/uploads/pets/15.jpg"
      }
    ],
    "activities": [
      "Выгул собак - 6 часов",
      "Кормление - 2 часа",
      "Уборка вольеров - 4 часа"
    ],
    "notes": "Все животные в хорошем состоянии",
    "created_at": "2025-01-20T18:00:00Z"
  }
}
```

#### POST /api/reports
Создать отчет

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "title": "Отчет за неделю",
  "description": "Выполнено 5 задач, выгуляно 3 собаки",
  "date": "2025-01-20",
  "tasks_completed": 5,
  "hours_worked": 12,
  "pets_helped": [15, 16, 17],
  "activities": [
    "Выгул собак - 6 часов",
    "Кормление - 2 часа",
    "Уборка вольеров - 4 часа"
  ],
  "notes": "Все животные в хорошем состоянии"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Отчет за неделю",
    "date": "2025-01-20"
  },
  "message": "Report created successfully"
}
```

### Учет времени

#### GET /api/hours
Получить учет часов

**Headers:**
```
Authorization: Bearer <token>
```

**Query параметры:**
- `date_from` - дата начала периода
- `date_to` - дата окончания периода

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "date": "2025-01-20",
      "hours": 4,
      "activity": "Выгул собак",
      "description": "Выгулял Шарика и Бобика",
      "volunteer_id": 5,
      "created_at": "2025-01-20T18:00:00Z"
    }
  ]
}
```

#### POST /api/hours
Добавить часы работы

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "date": "2025-01-20",
  "hours": 4,
  "activity": "Выгул собак",
  "description": "Выгулял Шарика и Бобика"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "date": "2025-01-20",
    "hours": 4
  },
  "message": "Hours logged successfully"
}
```

#### GET /api/hours/stats
Получить статистику часов

**Headers:**
```
Authorization: Bearer <token>
```

**Query параметры:**
- `period` - период (week, month, year)

**Response:**
```json
{
  "success": true,
  "data": {
    "total_hours": 48,
    "this_week": 12,
    "this_month": 48,
    "by_activity": {
      "Выгул собак": 24,
      "Кормление": 12,
      "Уборка вольеров": 12
    },
    "by_date": [
      {
        "date": "2025-01-20",
        "hours": 4
      }
    ]
  }
}
```

### Профиль

#### GET /api/profile
Получить мой профиль

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 5,
    "name": "Мария",
    "last_name": "Волонтер",
    "email": "maria@example.com",
    "avatar": "/uploads/users/5/avatar.jpg",
    "phone": "+7 (999) 123-45-67",
    "role": "volunteer",
    "volunteer_since": "2024-06-01",
    "total_hours": 120,
    "pets_helped": 15,
    "tasks_completed": 45
  }
}
```

### Статистика

#### GET /api/stats/overview
Получить общую статистику

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "my_pets_count": 3,
    "active_tasks": 5,
    "completed_tasks_month": 12,
    "hours_this_month": 48,
    "reports_this_month": 4
  }
}
```

---

## Зависимости от других сервисов

- **Auth Service** - авторизация
- **PetBase Service** - данные подопечных питомцев (чтение и обновление curator_id)

---

## Конфигурация

### .env
```bash
PORT=8500
AUTH_SERVICE_URL=http://localhost:7000
PETBASE_SERVICE_URL=http://localhost:8100
ALLOWED_ORIGINS=http://localhost:6200
DATABASE_PATH=./database/volunteer.db
```

---

## База данных

### Таблицы

#### volunteer_tasks
```sql
CREATE TABLE volunteer_tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    pet_id INTEGER,
    volunteer_id INTEGER NOT NULL,
    status TEXT DEFAULT 'pending', -- pending, in_progress, completed
    priority TEXT DEFAULT 'medium', -- low, medium, high
    due_date DATE,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    FOREIGN KEY (volunteer_id) REFERENCES users(id),
    FOREIGN KEY (pet_id) REFERENCES pets(id)
);
```

#### volunteer_reports
```sql
CREATE TABLE volunteer_reports (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT,
    date DATE NOT NULL,
    volunteer_id INTEGER NOT NULL,
    tasks_completed INTEGER DEFAULT 0,
    hours_worked INTEGER DEFAULT 0,
    pets_helped TEXT, -- JSON array of pet IDs
    activities TEXT, -- JSON array of activities
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (volunteer_id) REFERENCES users(id)
);
```

#### volunteer_hours
```sql
CREATE TABLE volunteer_hours (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    volunteer_id INTEGER NOT NULL,
    date DATE NOT NULL,
    hours INTEGER NOT NULL,
    activity TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (volunteer_id) REFERENCES users(id)
);
```

---

## Запуск

Backend:
```bash
cd volunteer/backend
air
```

Frontend:
```bash
cd volunteer/frontend
npm run dev
```

---

## Swagger документация

```
http://localhost:8500/swagger/index.html
```

---

## Примеры использования

### Получить подопечных питомцев
```bash
curl -X GET http://localhost:8500/api/my-pets \
  -H "Authorization: Bearer <token>"
```

### Создать задачу
```bash
curl -X POST http://localhost:8500/api/tasks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Выгулять Шарика",
    "description": "Прогулка в парке 30 минут",
    "pet_id": 15,
    "priority": "high",
    "due_date": "2025-01-20"
  }'
```

### Завершить задачу
```bash
curl -X POST http://localhost:8500/api/tasks/1/complete \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "notes": "Задача выполнена успешно"
  }'
```

### Создать отчет
```bash
curl -X POST http://localhost:8500/api/reports \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Отчет за неделю",
    "description": "Выполнено 5 задач",
    "date": "2025-01-20",
    "tasks_completed": 5,
    "hours_worked": 12,
    "pets_helped": [15, 16, 17]
  }'
```

### Добавить часы работы
```bash
curl -X POST http://localhost:8500/api/hours \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2025-01-20",
    "hours": 4,
    "activity": "Выгул собак",
    "description": "Выгулял Шарика и Бобика"
  }'
```

---

**Дата обновления:** 29 января 2025
