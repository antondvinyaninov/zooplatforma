# Admin Service - Административная панель

**Порт Backend:** 9000  
**Порт Frontend:** 4000  
**Назначение:** Модерация и управление платформой

---

## Описание

Admin Service - это панель администратора для:
- Модерации пользователей и контента
- Управления организациями
- Просмотра статистики
- Рассмотрения жалоб
- Управления системными настройками

**Доступ:** Только для пользователей с ролью `superadmin`

---

## База данных

`admin/database/admin.db`

### Таблицы

```sql
-- Логи администраторов
CREATE TABLE admin_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    admin_id INTEGER NOT NULL,
    action TEXT NOT NULL,
    target_type TEXT,
    target_id INTEGER,
    details TEXT,
    ip_address TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Очередь модерации
CREATE TABLE moderation_queue (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    item_type TEXT NOT NULL,
    item_id INTEGER NOT NULL,
    reason TEXT,
    status TEXT DEFAULT 'pending',
    assigned_to INTEGER,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Забаненные пользователи
CREATE TABLE banned_users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    banned_by INTEGER NOT NULL,
    reason TEXT NOT NULL,
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Забаненные IP
CREATE TABLE banned_ips (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ip_address TEXT NOT NULL,
    reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Настройки системы
CREATE TABLE system_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT UNIQUE NOT NULL,
    value TEXT,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

## API Endpoints

### Пользователи

#### GET /api/admin/users
Получить список пользователей

**Headers:**
```
Authorization: Bearer <token>
```

**Query параметры:**
- `page` - номер страницы
- `limit` - количество на странице
- `search` - поиск по имени/email
- `role` - фильтр по роли
- `banned` - только забаненные

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "email": "user@example.com",
      "name": "Иван Иванов",
      "role": "user",
      "verified": true,
      "banned": false,
      "created_at": "2025-01-01T10:00:00Z",
      "last_login": "2025-01-17T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 150
  }
}
```

#### GET /api/admin/users/:id
Получить информацию о пользователе

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
    "email": "user@example.com",
    "name": "Иван",
    "last_name": "Иванов",
    "role": "user",
    "avatar": "/uploads/users/1/avatar.jpg",
    "verified": true,
    "banned": false,
    "created_at": "2025-01-01T10:00:00Z",
    "stats": {
      "posts_count": 50,
      "friends_count": 25,
      "pets_count": 3
    }
  }
}
```

#### PUT /api/admin/users/:id
Обновить пользователя

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "role": "moderator",
  "verified": true
}
```

**Response:**
```json
{
  "success": true,
  "message": "User updated successfully"
}
```

#### DELETE /api/admin/users/:id
Забанить пользователя

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "reason": "Нарушение правил",
  "expires_at": "2025-02-17T10:00:00Z"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User banned successfully"
}
```

#### POST /api/admin/users/:id/unban
Разбанить пользователя

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "message": "User unbanned successfully"
}
```

### Посты

#### GET /api/admin/posts
Получить список постов

**Headers:**
```
Authorization: Bearer <token>
```

**Query параметры:**
- `page` - номер страницы
- `limit` - количество на странице
- `author_id` - фильтр по автору
- `status` - фильтр по статусу
- `reported` - только с жалобами

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "author_id": 5,
      "author_name": "Иван Иванов",
      "content": "Текст поста...",
      "status": "published",
      "reports_count": 2,
      "created_at": "2025-01-17T10:00:00Z"
    }
  ]
}
```

#### GET /api/admin/posts/:id
Получить информацию о посте

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
    "author_id": 5,
    "author_name": "Иван Иванов",
    "content": "Текст поста...",
    "attachments": [],
    "likes_count": 10,
    "comments_count": 5,
    "reports": [
      {
        "id": 1,
        "reporter_id": 10,
        "reason": "spam",
        "description": "Рекламный спам"
      }
    ],
    "created_at": "2025-01-17T10:00:00Z"
  }
}
```

#### DELETE /api/admin/posts/:id
Удалить пост

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "reason": "Нарушение правил"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Post deleted successfully"
}
```

#### POST /api/admin/posts/:id/hide
Скрыть пост

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "message": "Post hidden successfully"
}
```

#### POST /api/admin/posts/:id/show
Показать пост

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "message": "Post shown successfully"
}
```

### Организации

#### GET /api/admin/organizations
Получить список организаций

**Headers:**
```
Authorization: Bearer <token>
```

**Query параметры:**
- `type` - фильтр по типу
- `verified` - фильтр по верификации
- `city` - фильтр по городу

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "name": "Ветклиника Айболит",
      "type": "clinic",
      "city": "Москва",
      "verified": true,
      "members_count": 5,
      "created_at": "2025-01-01T10:00:00Z"
    }
  ]
}
```

#### PUT /api/admin/organizations/:id
Верифицировать/обновить организацию

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "verified": true,
  "notes": "Проверено, документы в порядке"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Organization updated successfully"
}
```

#### DELETE /api/admin/organizations/:id
Удалить организацию

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "reason": "Фейковая организация"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Organization deleted successfully"
}
```

#### GET /api/admin/organizations/stats
Статистика организаций

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total": 150,
    "by_type": {
      "clinic": 80,
      "shelter": 50,
      "store": 20
    },
    "verified": 100,
    "pending_verification": 50
  }
}
```

### Модерация

#### GET /api/admin/moderation/reports
Получить список жалоб

**Headers:**
```
Authorization: Bearer <token>
```

**Query параметры:**
- `status` - фильтр по статусу (pending, resolved, rejected)
- `target_type` - фильтр по типу (post, user, comment)

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "reporter_id": 10,
      "reporter_name": "Петр Петров",
      "target_type": "post",
      "target_id": 123,
      "reason": "spam",
      "description": "Рекламный спам",
      "status": "pending",
      "created_at": "2025-01-17T10:00:00Z"
    }
  ]
}
```

#### GET /api/admin/moderation/reports/:id
Получить информацию о жалобе

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
    "reporter": {
      "id": 10,
      "name": "Петр Петров"
    },
    "target_type": "post",
    "target_id": 123,
    "target_content": "Текст поста...",
    "reason": "spam",
    "description": "Рекламный спам",
    "status": "pending",
    "created_at": "2025-01-17T10:00:00Z"
  }
}
```

#### PUT /api/admin/moderation/reports/:id
Рассмотреть жалобу

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "status": "resolved",
  "action": "post_deleted",
  "notes": "Пост удален за нарушение правил"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Report resolved successfully"
}
```

#### GET /api/admin/moderation/queue
Получить очередь модерации

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
      "item_type": "post",
      "item_id": 123,
      "reason": "Автоматическая проверка",
      "status": "pending",
      "created_at": "2025-01-17T10:00:00Z"
    }
  ]
}
```

#### GET /api/admin/moderation/stats
Статистика модерации

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_reports": 150,
    "pending_reports": 20,
    "resolved_reports": 120,
    "rejected_reports": 10,
    "reports_by_type": {
      "spam": 50,
      "inappropriate": 30,
      "harassment": 20
    }
  }
}
```

### Статистика

#### GET /api/admin/stats/overview
Общая статистика платформы

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "users": {
      "total": 10000,
      "active_today": 500,
      "new_this_week": 100
    },
    "posts": {
      "total": 50000,
      "today": 200
    },
    "pets": {
      "total": 15000
    },
    "organizations": {
      "total": 150,
      "verified": 100
    }
  }
}
```

#### GET /api/admin/stats/users
Статистика пользователей

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total": 10000,
    "by_role": {
      "user": 9900,
      "moderator": 50,
      "superadmin": 5
    },
    "verified": 5000,
    "banned": 100,
    "registrations_by_month": [
      {"month": "2025-01", "count": 500}
    ]
  }
}
```

### Логи

#### GET /api/admin/logs
Получить системные логи

**Headers:**
```
Authorization: Bearer <token>
```

**Query параметры:**
- `page` - номер страницы
- `limit` - количество на странице
- `action` - фильтр по действию
- `admin_id` - фильтр по администратору

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "admin_id": 1,
      "admin_name": "Супер Админ",
      "action": "user_banned",
      "target_type": "user",
      "target_id": 123,
      "details": "Забанен за спам",
      "ip_address": "192.168.1.1",
      "created_at": "2025-01-17T10:00:00Z"
    }
  ]
}
```

#### GET /api/admin/logs/user/:id
Получить логи пользователя

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
      "action": "login",
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0...",
      "created_at": "2025-01-17T10:00:00Z"
    }
  ]
}
```

### Настройки

#### GET /api/admin/settings
Получить настройки системы

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "maintenance_mode": false,
    "registration_enabled": true,
    "max_upload_size": 10485760,
    "allowed_file_types": ["jpg", "png", "gif"]
  }
}
```

#### PUT /api/admin/settings
Обновить настройки

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "maintenance_mode": true,
  "registration_enabled": false
}
```

**Response:**
```json
{
  "success": true,
  "message": "Settings updated successfully"
}
```

---

## Зависимости от других сервисов

- **Auth Service** - проверка прав superadmin
- **Main Service** - модерация постов, пользователей
- **PetBase Service** - статистика питомцев

---

## Конфигурация

### .env
```bash
PORT=9000
DATABASE_URL=./database/admin.db
AUTH_SERVICE_URL=http://localhost:7000
MAIN_SERVICE_URL=http://localhost:8000
PETBASE_SERVICE_URL=http://localhost:8100
ALLOWED_ORIGINS=http://localhost:4000
```

---

## Запуск

```bash
cd admin/backend
air
```

Frontend:
```bash
cd admin/frontend
npm run dev
```

---

## Swagger документация

```
http://localhost:9000/swagger/index.html
```

---

**Дата обновления:** 17 января 2025
