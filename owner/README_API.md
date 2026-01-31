# Owner Service - Кабинет владельца питомца

**Порт Backend:** 8400  
**Порт Frontend:** 6100  
**Назначение:** Управление своими питомцами

---

## Описание

Owner Service - это упрощенный кабинет для владельцев питомцев:
- Управление своими питомцами
- Просмотр медицинских записей
- Просмотр вакцинаций
- Загрузка фотографий

**Примечание:** Не имеет своей БД, все операции через PetBase Service

---

## API Endpoints

### Питомцы

#### GET /api/pets
Получить моих питомцев

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
      "name": "Барсик",
      "species": "cat",
      "breed": "Британская короткошерстная",
      "birth_date": "2020-05-15",
      "gender": "male",
      "photo": "/uploads/pets/1.jpg",
      "created_at": "2025-01-01T10:00:00Z"
    }
  ]
}
```

#### GET /api/pets/:id
Получить информацию о питомце

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
    "name": "Барсик",
    "species": "cat",
    "breed": "Британская короткошерстная",
    "birth_date": "2020-05-15",
    "gender": "male",
    "color": "Серый",
    "size": "medium",
    "photo": "/uploads/pets/1.jpg",
    "chip_number": "123456789",
    "tag_number": "TAG001",
    "sterilization_date": "2021-01-10",
    "story": "Добрый и ласковый кот",
    "created_at": "2025-01-01T10:00:00Z"
  }
}
```

#### POST /api/pets
Добавить питомца

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
  "birth_date": "2020-05-15",
  "gender": "male",
  "color": "Серый",
  "size": "medium"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Барсик",
    "species": "cat"
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
  "color": "Темно-серый",
  "story": "Обновленная история"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Барсик Обновленный"
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

#### POST /api/pets/:id/photo
Загрузить фото питомца

**Headers:**
```
Authorization: Bearer <token>
```

**Request (multipart/form-data):**
```
photo: [binary]
```

**Response:**
```json
{
  "success": true,
  "data": {
    "photo_url": "/uploads/pets/1.jpg"
  }
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
      "date": "2025-01-15",
      "clinic_name": "Ветклиника Айболит",
      "diagnosis": "Здоров",
      "treatment": "Вакцинация",
      "veterinarian": "Иванов И.И.",
      "notes": "Плановый осмотр"
    }
  ]
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
      "vaccine_name": "Бешенство",
      "date": "2025-01-10",
      "next_date": "2026-01-10",
      "clinic_name": "Ветклиника Айболит",
      "veterinarian": "Иванов И.И."
    }
  ]
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
    "name": "Иван",
    "last_name": "Иванов",
    "email": "ivan@example.com",
    "avatar": "/uploads/users/5/avatar.jpg",
    "phone": "+7 (999) 123-45-67"
  }
}
```

### Справочники

#### GET /api/breeds
Получить список пород

**Query параметры:**
- `species_id` - фильтр по виду

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Британская короткошерстная",
      "species_id": 2
    }
  ]
}
```

---

## Зависимости от других сервисов

- **Auth Service** - авторизация
- **PetBase Service** - все операции с питомцами

---

## Конфигурация

### .env
```bash
PORT=8400
AUTH_SERVICE_URL=http://localhost:7000
PETBASE_SERVICE_URL=http://localhost:8100
ALLOWED_ORIGINS=http://localhost:6100
```

---

## Запуск

```bash
cd owner/backend
air
```

Frontend:
```bash
cd owner/frontend
npm run dev
```

---

## Swagger документация

```
http://localhost:8400/swagger/index.html
```

---

**Дата обновления:** 17 января 2025
