# SSO Архитектура ЗооПлатформы

Single Sign-On система для всех микросервисов платформы.

---

## 🎯 Цель

Единая система авторизации для всех сервисов:
- Основное приложение (`zooplatforma.ru`)
- Админ-панель (`sadmin.zooplatforma.ru`)
- Аналитика (`analytics.zooplatforma.ru`)
- Будущие микросервисы

---

## 🏗️ Архитектура

### Вариант 1: Shared Cookie SSO (текущая реализация)

```
┌─────────────────────────────────────────────────────────────┐
│              Main Backend (Auth Service)                     │
│                   backend:8080                               │
│                                                              │
│  POST /api/auth/login    → Выдает JWT в cookie              │
│  POST /api/auth/logout   → Удаляет cookie                   │
│  GET  /api/auth/me       → Проверяет токен                  │
│  GET  /api/auth/verify   → Публичная проверка токена        │
│                                                              │
│  Cookie: auth_token                                          │
│  Domain: .zooplatforma.ru (доступен всем поддоменам)        │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ Shared JWT Token
            ┌───────────────┼───────────────┐
            │               │               │
            ▼               ▼               ▼
    ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
    │   Frontend   │ │  Admin Panel │ │  Analytics   │
    │   :3000      │ │    :3001     │ │    :3002     │
    │              │ │              │ │              │
    │ Читает       │ │ Читает       │ │ Читает       │
    │ auth_token   │ │ auth_token   │ │ auth_token   │
    └──────────────┘ └──────────────┘ └──────────────┘
            │               │               │
            ▼               ▼               ▼
    ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
    │   Backend    │ │ Admin Backend│ │Analytics API │
    │   :8080      │ │    :8081     │ │    :8082     │
    │              │ │              │ │              │
    │ Проверяет    │ │ Проверяет    │ │ Проверяет    │
    │ токен        │ │ токен +      │ │ токен        │
    │              │ │ права админа │ │              │
    └──────────────┘ └──────────────┘ └──────────────┘
```

---

## 🔑 JWT Token Structure

```json
{
  "user_id": 1,
  "email": "user@example.com",
  "roles": ["user", "superadmin"],
  "permissions": ["read", "write", "admin"],
  "exp": 1234567890,
  "iat": 1234567890
}
```

**Важно:**
- `roles` - массив ролей пользователя
- `permissions` - детальные права
- Токен подписан общим `JWT_SECRET`

---

## 🔐 Процесс авторизации

### 1. Вход пользователя

```
User → Frontend (zooplatforma.ru)
  ↓
POST /api/auth/login
  ↓
Backend проверяет credentials
  ↓
Генерирует JWT токен
  ↓
Set-Cookie: auth_token=<JWT>
  Domain: .zooplatforma.ru
  HttpOnly: true
  Secure: true (в production)
  SameSite: Lax
  ↓
Frontend получает токен в cookie
```

### 2. Доступ к админ-панели

```
User → Admin Frontend (sadmin.zooplatforma.ru)
  ↓
Читает cookie auth_token (доступен благодаря .zooplatforma.ru)
  ↓
GET /api/admin/auth/me
  ↓
Admin Backend проверяет:
  1. Токен валиден?
  2. Пользователь существует?
  3. Есть роль superadmin?
  ↓
Если все ОК → доступ разрешен
Если нет → редирект на zooplatforma.ru/auth
```

### 3. Проверка токена в микросервисах

Каждый микросервис:
1. Читает cookie `auth_token`
2. Парсит JWT
3. Проверяет подпись (общий `JWT_SECRET`)
4. Проверяет срок действия
5. Проверяет права (roles/permissions)

---

## 📡 API для SSO

### Main Backend (Auth Service)

#### POST /api/auth/login
Вход пользователя.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "name": "User Name",
    "email": "user@example.com",
    "roles": ["user"]
  }
}
```

**Cookie:**
```
Set-Cookie: auth_token=<JWT>; Domain=.zooplatforma.ru; HttpOnly; Secure; SameSite=Lax; Max-Age=604800
```

#### GET /api/auth/verify
Публичная проверка токена (для микросервисов).

**Request:**
```
Cookie: auth_token=<JWT>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user_id": 1,
    "email": "user@example.com",
    "roles": ["user", "superadmin"],
    "valid": true
  }
}
```

#### POST /api/auth/logout
Выход из всех сервисов.

**Response:**
```
Set-Cookie: auth_token=; Domain=.zooplatforma.ru; Max-Age=0
```

---

## 🛡️ Безопасность

### 1. JWT Secret
- Общий для всех сервисов
- Хранится в `.env` каждого сервиса
- Минимум 32 символа
- Регулярная ротация (раз в 3 месяца)

### 2. Cookie настройки
```javascript
{
  Domain: '.zooplatforma.ru',  // Доступен всем поддоменам
  HttpOnly: true,               // Защита от XSS
  Secure: true,                 // Только HTTPS (production)
  SameSite: 'Lax',             // Защита от CSRF
  MaxAge: 604800               // 7 дней
}
```

### 3. Проверка прав
Каждый микросервис проверяет:
- Токен валиден
- Срок не истек
- Пользователь существует
- Есть нужные права (roles/permissions)

### 4. Логирование
- Все входы логируются
- Все проверки токенов логируются
- Подозрительная активность → алерты

---

## 🔧 Реализация

### 1. Обновить Main Backend

**Добавить роли в JWT:**

```go
// backend/middleware/auth.go
type Claims struct {
    UserID      int      `json:"user_id"`
    Email       string   `json:"email"`
    Roles       []string `json:"roles"`        // NEW
    Permissions []string `json:"permissions"`  // NEW
    jwt.RegisteredClaims
}

func GenerateToken(userID int, email string, roles []string) (string, error) {
    claims := &Claims{
        UserID:      userID,
        Email:       email,
        Roles:       roles,
        Permissions: getPermissions(roles),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    // ...
}
```

**Добавить endpoint для проверки:**

```go
// backend/handlers/auth.go
func VerifyTokenHandler(w http.ResponseWriter, r *http.Request) {
    cookie, err := r.Cookie("auth_token")
    if err != nil {
        sendError(w, "No token", http.StatusUnauthorized)
        return
    }

    token, err := middleware.ParseToken(cookie.Value)
    if err != nil {
        sendError(w, "Invalid token", http.StatusUnauthorized)
        return
    }

    sendSuccess(w, map[string]interface{}{
        "user_id": token.UserID,
        "email":   token.Email,
        "roles":   token.Roles,
        "valid":   true,
    })
}
```

**Обновить cookie domain:**

```go
// backend/handlers/auth.go
http.SetCookie(w, &http.Cookie{
    Name:     "auth_token",
    Value:    token,
    Path:     "/",
    Domain:   ".zooplatforma.ru",  // Доступен всем поддоменам
    HttpOnly: true,
    Secure:   true,  // В production
    SameSite: http.SameSiteLaxMode,
    MaxAge:   604800, // 7 дней
})
```

### 2. Обновить Admin Backend

**Использовать общий токен:**

```go
// admin/backend/middleware/admin.go
func SuperAdminMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Читаем общий токен
        cookie, err := r.Cookie("auth_token")  // Не admin_token!
        if err != nil {
            sendError(w, "Не авторизован", http.StatusUnauthorized)
            return
        }

        // Парсим токен
        token, err := ParseToken(cookie.Value)
        if err != nil {
            sendError(w, "Неверный токен", http.StatusUnauthorized)
            return
        }

        // Проверяем роль superadmin
        if !contains(token.Roles, "superadmin") {
            sendError(w, "Требуются права суперадмина", http.StatusForbidden)
            return
        }

        // Добавляем в контекст
        ctx := context.WithValue(r.Context(), "userID", token.UserID)
        ctx = context.WithValue(ctx, "roles", token.Roles)

        next(w, r.WithContext(ctx))
    }
}
```

**Убрать отдельную авторизацию:**

```go
// admin/backend/handlers/auth.go
// Удалить AdminLoginHandler - используем основной
// Оставить только проверку прав

func AdminMeHandler(w http.ResponseWriter, r *http.Request) {
    // Просто проверяем токен и права
    cookie, err := r.Cookie("auth_token")
    // ...
}
```

### 3. Обновить таблицу users

**Добавить поле roles:**

```sql
ALTER TABLE users ADD COLUMN roles TEXT DEFAULT '["user"]';
```

**Или использовать таблицу admins:**

```sql
-- Оставляем как есть
-- admins.role = 'superadmin' → добавляем в JWT roles
```

---

## 🚀 Миграция

### Шаг 1: Обновить Main Backend
- Добавить roles в JWT
- Обновить cookie domain на `.zooplatforma.ru`
- Добавить endpoint `/api/auth/verify`

### Шаг 2: Обновить Admin Backend
- Использовать общий токен `auth_token`
- Убрать отдельную авторизацию
- Проверять роли из JWT

### Шаг 3: Обновить Frontend
- Админ-панель редиректит на основной сайт для входа
- После входа возвращается обратно

### Шаг 4: Тестирование
- Вход на основном сайте
- Переход на админ-панель без повторного входа
- Выход из одного сервиса → выход из всех

---

## 📊 Преимущества

✅ **Единая авторизация** - один раз вошел, доступ ко всем сервисам  
✅ **Безопасность** - HttpOnly cookies, общий secret  
✅ **Простота** - не нужен отдельный Auth Service  
✅ **Масштабируемость** - легко добавлять новые микросервисы  
✅ **Аудит** - все действия логируются  

---

## 🔮 Будущее (v2.0)

### OAuth2 / OpenID Connect
Для более сложных сценариев:
- Вход через соцсети
- API для сторонних приложений
- Refresh tokens
- Scope и permissions

### Отдельный Auth Service
Если микросервисов станет много:
- Выделить auth в отдельный сервис
- Использовать Redis для сессий
- Добавить rate limiting
- Мониторинг и алерты

---

**Версия:** 1.0  
**Дата:** 12 декабря 2025  
**Статус:** Проектирование
