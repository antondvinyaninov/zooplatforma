---
inclusion: always
---

# API Design Rules - Правила проектирования API

## ⚠️ КРИТИЧЕСКОЕ ПРАВИЛО: Единый стиль API

**Все микросервисы следуют одинаковым правилам проектирования API!**

## REST Conventions

### Именование endpoints

**Формат:** `/api/<resource>/<action>`

**Правила:**
- Используй существительные во множественном числе
- Используй kebab-case для составных слов
- Избегай глаголов в URL (используй HTTP методы)

**✅ Правильно:**
```
GET    /api/users              # Список пользователей
GET    /api/users/123          # Конкретный пользователь
POST   /api/users              # Создать пользователя
PUT    /api/users/123          # Обновить пользователя
DELETE /api/users/123          # Удалить пользователя

GET    /api/organizations      # Список организаций
GET    /api/organizations/5    # Конкретная организация
POST   /api/organizations/5/members  # Добавить участника

GET    /api/pet-cards          # Список карточек (kebab-case)
```

**❌ Неправильно:**
```
GET    /api/getUsers           # Глагол в URL
GET    /api/user               # Единственное число
POST   /api/createUser         # Глагол в URL
GET    /api/petCards           # camelCase вместо kebab-case
```

### HTTP методы

| Метод | Назначение | Идемпотентность | Body |
|-------|-----------|-----------------|------|
| GET | Получить данные | Да | Нет |
| POST | Создать ресурс | Нет | Да |
| PUT | Обновить ресурс (полностью) | Да | Да |
| PATCH | Обновить ресурс (частично) | Нет | Да |
| DELETE | Удалить ресурс | Да | Нет |

**Примеры:**

```go
// GET - получение данных
GET /api/posts?status=published&limit=10

// POST - создание
POST /api/posts
Body: { "content": "Hello", "status": "draft" }

// PUT - полное обновление
PUT /api/posts/123
Body: { "content": "Updated", "status": "published", "tags": ["news"] }

// PATCH - частичное обновление
PATCH /api/posts/123
Body: { "status": "published" }

// DELETE - удаление
DELETE /api/posts/123
```

## Стандарты ответов

### Успешный ответ

**Формат:**
```json
{
  "success": true,
  "data": { /* данные */ },
  "message": "Optional success message"
}
```

**Примеры:**

```json
// Список ресурсов
{
  "success": true,
  "data": [
    { "id": 1, "name": "Item 1" },
    { "id": 2, "name": "Item 2" }
  ]
}

// Один ресурс
{
  "success": true,
  "data": {
    "id": 1,
    "name": "Item 1",
    "created_at": "2024-01-16T10:00:00Z"
  }
}

// Создание ресурса
{
  "success": true,
  "data": {
    "id": 123,
    "name": "New Item"
  },
  "message": "Item created successfully"
}

// Операция без данных
{
  "success": true,
  "message": "Item deleted successfully"
}
```

### Ответ с ошибкой

**Формат:**
```json
{
  "success": false,
  "error": "Error message",
  "details": { /* опциональные детали */ }
}
```

**Примеры:**

```json
// Простая ошибка
{
  "success": false,
  "error": "User not found"
}

// Ошибка валидации
{
  "success": false,
  "error": "Validation failed",
  "details": {
    "email": "Invalid email format",
    "password": "Password must be at least 8 characters"
  }
}

// Ошибка авторизации
{
  "success": false,
  "error": "Unauthorized"
}
```

### HTTP Status Codes

**Используй правильные коды:**

| Код | Название | Когда использовать |
|-----|----------|-------------------|
| 200 | OK | Успешный GET, PUT, PATCH, DELETE |
| 201 | Created | Успешный POST (создание) |
| 204 | No Content | Успешное удаление без тела ответа |
| 400 | Bad Request | Невалидные данные от клиента |
| 401 | Unauthorized | Не авторизован (нет токена) |
| 403 | Forbidden | Авторизован, но нет прав |
| 404 | Not Found | Ресурс не найден |
| 409 | Conflict | Конфликт (например, email уже существует) |
| 500 | Internal Server Error | Ошибка сервера |

**Примеры в Go:**

```go
// 200 OK
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(response)

// 201 Created
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(response)

// 400 Bad Request
http.Error(w, `{"success":false,"error":"Invalid input"}`, http.StatusBadRequest)

// 401 Unauthorized
http.Error(w, `{"success":false,"error":"Unauthorized"}`, http.StatusUnauthorized)

// 404 Not Found
http.Error(w, `{"success":false,"error":"User not found"}`, http.StatusNotFound)

// 500 Internal Server Error
http.Error(w, `{"success":false,"error":"Internal server error"}`, http.StatusInternalServerError)
```

## Pagination

**Для списков используй pagination:**

**Query параметры:**
- `page` - номер страницы (начиная с 1)
- `limit` - количество элементов на странице (по умолчанию 20, максимум 100)

**Пример запроса:**
```
GET /api/posts?page=2&limit=10
```

**Формат ответа:**
```json
{
  "success": true,
  "data": [
    { "id": 11, "title": "Post 11" },
    { "id": 12, "title": "Post 12" }
  ],
  "pagination": {
    "page": 2,
    "limit": 10,
    "total": 45,
    "total_pages": 5,
    "has_next": true,
    "has_prev": true
  }
}
```

**Реализация в Go:**

```go
func GetPosts(w http.ResponseWriter, r *http.Request) {
    // Получить параметры
    pageStr := r.URL.Query().Get("page")
    limitStr := r.URL.Query().Get("limit")
    
    page := 1
    if pageStr != "" {
        page, _ = strconv.Atoi(pageStr)
        if page < 1 {
            page = 1
        }
    }
    
    limit := 20
    if limitStr != "" {
        limit, _ = strconv.Atoi(limitStr)
        if limit < 1 || limit > 100 {
            limit = 20
        }
    }
    
    offset := (page - 1) * limit
    
    // Получить общее количество
    var total int
    db.QueryRow("SELECT COUNT(*) FROM posts").Scan(&total)
    
    // Получить данные
    rows, err := db.Query("SELECT * FROM posts LIMIT ? OFFSET ?", limit, offset)
    // ... обработка rows
    
    // Ответ
    response := map[string]interface{}{
        "success": true,
        "data": posts,
        "pagination": map[string]interface{}{
            "page": page,
            "limit": limit,
            "total": total,
            "total_pages": (total + limit - 1) / limit,
            "has_next": offset + limit < total,
            "has_prev": page > 1,
        },
    }
    
    json.NewEncoder(w).Encode(response)
}
```

## Filtering & Sorting

### Фильтрация

**Query параметры для фильтров:**

```
GET /api/posts?status=published&author_id=5
GET /api/organizations?type=clinic&region=moscow
GET /api/users?role=admin&verified=true
```

**Реализация:**

```go
func GetPosts(w http.ResponseWriter, r *http.Request) {
    query := "SELECT * FROM posts WHERE 1=1"
    args := []interface{}{}
    
    // Фильтр по статусу
    if status := r.URL.Query().Get("status"); status != "" {
        query += " AND status = ?"
        args = append(args, status)
    }
    
    // Фильтр по автору
    if authorID := r.URL.Query().Get("author_id"); authorID != "" {
        query += " AND author_id = ?"
        args = append(args, authorID)
    }
    
    rows, err := db.Query(query, args...)
    // ...
}
```

### Сортировка

**Query параметры:**
- `sort` - поле для сортировки
- `order` - направление (asc/desc)

```
GET /api/posts?sort=created_at&order=desc
GET /api/users?sort=name&order=asc
```

**Реализация:**

```go
func GetPosts(w http.ResponseWriter, r *http.Request) {
    query := "SELECT * FROM posts WHERE 1=1"
    
    // Сортировка
    sortField := r.URL.Query().Get("sort")
    order := r.URL.Query().Get("order")
    
    // Whitelist полей для безопасности
    allowedSorts := map[string]bool{
        "id": true,
        "created_at": true,
        "updated_at": true,
        "title": true,
    }
    
    if allowedSorts[sortField] {
        query += " ORDER BY " + sortField
        if order == "desc" {
            query += " DESC"
        } else {
            query += " ASC"
        }
    } else {
        query += " ORDER BY created_at DESC" // По умолчанию
    }
    
    rows, err := db.Query(query)
    // ...
}
```

## Поиск

**Query параметр:** `search` или `q`

```
GET /api/users?search=anton
GET /api/organizations?q=ветклиника
```

**Реализация:**

```go
func SearchUsers(w http.ResponseWriter, r *http.Request) {
    searchQuery := r.URL.Query().Get("search")
    
    if searchQuery == "" {
        http.Error(w, `{"success":false,"error":"Search query required"}`, http.StatusBadRequest)
        return
    }
    
    query := `
        SELECT id, name, email, avatar 
        FROM users 
        WHERE name LIKE ? OR email LIKE ?
        LIMIT 10
    `
    
    searchPattern := "%" + searchQuery + "%"
    rows, err := db.Query(query, searchPattern, searchPattern)
    // ...
}
```

## Вложенные ресурсы

**Для связанных ресурсов:**

```
GET    /api/organizations/5/members       # Участники организации 5
POST   /api/organizations/5/members       # Добавить участника
DELETE /api/organizations/5/members/10    # Удалить участника 10

GET    /api/posts/123/comments            # Комментарии к посту 123
POST   /api/posts/123/comments            # Добавить комментарий
```

**Реализация:**

```go
// GET /api/organizations/{id}/members
func GetMembers(w http.ResponseWriter, r *http.Request) {
    // Получить organization_id из URL
    vars := mux.Vars(r) // если используешь gorilla/mux
    orgID := vars["id"]
    
    // Или из контекста (если через middleware)
    orgID, ok := r.Context().Value(middleware.OrgIDKey).(int)
    
    rows, err := db.Query(`
        SELECT om.id, om.user_id, om.role, u.name, u.email, u.avatar
        FROM organization_members om
        JOIN users u ON om.user_id = u.id
        WHERE om.organization_id = ?
    `, orgID)
    // ...
}
```

## Batch операции

**Для массовых операций:**

```
POST /api/posts/batch-delete
Body: { "ids": [1, 2, 3, 4, 5] }

POST /api/users/batch-update
Body: { 
  "ids": [10, 20, 30],
  "updates": { "status": "active" }
}
```

## Версионирование API

**Пока не используем, но в будущем:**

```
GET /api/v1/users
GET /api/v2/users
```

**Или через заголовок:**
```
GET /api/users
Header: Accept: application/vnd.api+json; version=1
```

## Специальные endpoints

### Health Check

**Обязательно для каждого микросервиса:**

```
GET /api/health
```

**Ответ:**
```json
{
  "success": true,
  "service": "clinic-backend",
  "version": "0.1.0",
  "status": "healthy",
  "timestamp": "2024-01-16T10:00:00Z"
}
```

### Статистика

```
GET /api/stats/overview
GET /api/stats/users
GET /api/stats/posts
```

## Примеры из проекта

### Main Backend

```go
// Посты
GET    /api/posts              # Список постов
GET    /api/posts/:id          # Конкретный пост
POST   /api/posts              # Создать пост
PUT    /api/posts/:id          # Обновить пост
DELETE /api/posts/:id          # Удалить пост
GET    /api/posts/drafts       # Черновики

// Лайки
POST   /api/posts/:id/like     # Лайк/анлайк
GET    /api/posts/:id/like     # Статус лайка

// Комментарии
GET    /api/comments?post_id=123
POST   /api/comments
DELETE /api/comments/:id
```

### Clinic Backend

```go
// Клиники
GET    /api/clinics            # Мои клиники
GET    /api/clinics/:id        # Конкретная клиника

// Участники
GET    /api/members            # Участники клиники (с X-Clinic-ID)
POST   /api/members            # Добавить участника
PUT    /api/members            # Обновить участника
DELETE /api/members            # Удалить участника

// Пациенты
GET    /api/patients           # Пациенты клиники
```

### PetBase Backend

```go
// Виды
GET    /api/species
POST   /api/species
GET    /api/species/:id
PUT    /api/species/:id
DELETE /api/species/:id

// Породы
GET    /api/breeds
GET    /api/breeds/species/:id  # Породы конкретного вида
POST   /api/breeds
PUT    /api/breeds/:id
DELETE /api/breeds/:id

// Карточки
GET    /api/cards
GET    /api/cards/breed/:id    # Карточки конкретной породы
POST   /api/cards
PUT    /api/cards/:id
DELETE /api/cards/:id
```

## Checklist для нового endpoint

- [ ] Используется правильный HTTP метод
- [ ] URL следует REST conventions (существительные, множественное число)
- [ ] Ответ в стандартном формате `{ success, data, error }`
- [ ] Правильный HTTP status code
- [ ] Добавлена валидация входных данных
- [ ] Добавлена авторизация (если нужно)
- [ ] Добавлено логирование
- [ ] Обработаны ошибки
- [ ] Добавлена документация в `docs/API_ENDPOINTS.md`
- [ ] Обновлен CHANGELOG.md

## Best Practices

### 1. Валидация данных

```go
func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, `{"success":false,"error":"Invalid JSON"}`, http.StatusBadRequest)
        return
    }
    
    // Валидация
    if user.Email == "" {
        http.Error(w, `{"success":false,"error":"Email is required"}`, http.StatusBadRequest)
        return
    }
    
    if !isValidEmail(user.Email) {
        http.Error(w, `{"success":false,"error":"Invalid email format"}`, http.StatusBadRequest)
        return
    }
    
    // ...
}
```

### 2. Логирование

```go
func GetUser(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    log.Printf("📋 GetUser: user_id=%s", userID)
    
    // ...
    
    if err != nil {
        log.Printf("❌ GetUser error: %v", err)
        http.Error(w, `{"success":false,"error":"Internal error"}`, http.StatusInternalServerError)
        return
    }
    
    log.Printf("✅ GetUser success: user_id=%s", userID)
}
```

### 3. Обработка ошибок

```go
func GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := db.GetUserByID(userID)
    
    if err == sql.ErrNoRows {
        // 404 - не найден
        http.Error(w, `{"success":false,"error":"User not found"}`, http.StatusNotFound)
        return
    }
    
    if err != nil {
        // 500 - ошибка сервера
        log.Printf("Database error: %v", err)
        http.Error(w, `{"success":false,"error":"Internal server error"}`, http.StatusInternalServerError)
        return
    }
    
    // 200 - успех
    response := map[string]interface{}{
        "success": true,
        "data": user,
    }
    json.NewEncoder(w).Encode(response)
}
```

---

**Помни:** Единый стиль API делает проект понятным и предсказуемым!
