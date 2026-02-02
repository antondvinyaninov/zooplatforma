# Auth Service - Сервис авторизации

**Порт:** 7000  
**Назначение:** Единая точка авторизации для всей платформы (SSO)

---

## Описание

Auth Service - это центральный сервис авторизации, который обеспечивает:
- Регистрацию и вход пользователей
- Генерацию и проверку JWT токенов
- Управление сессиями
- Сброс и изменение паролей
- Подтверждение email

Все остальные сервисы (Main, Admin, Clinic и т.д.) используют Auth Service для проверки авторизации.

---

## База данных

`auth/database/auth.db`

### Таблицы

```sql
-- Пользователи
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    last_name TEXT,
    role TEXT DEFAULT 'user',
    avatar TEXT,
    email_verified BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Сессии
CREATE TABLE sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token TEXT UNIQUE NOT NULL,
    expires_at DATETIME NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Refresh токены
CREATE TABLE refresh_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token TEXT UNIQUE NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Сброс пароля
CREATE TABLE password_resets (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token TEXT UNIQUE NOT NULL,
    expires_at DATETIME NOT NULL,
    used BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Подтверждение email
CREATE TABLE email_verifications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    token TEXT UNIQUE NOT NULL,
    expires_at DATETIME NOT NULL,
    verified BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

---

## API Endpoints

### Регистрация и вход

#### POST /api/auth/register
Регистрация нового пользователя

**Request:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "name": "Иван",
  "last_name": "Иванов"
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "Иван",
      "last_name": "Иванов",
      "role": "user"
    },
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

#### POST /api/auth/login
Вход в систему

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
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "Иван",
      "role": "user"
    },
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

#### POST /api/auth/logout
Выход из системы

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "message": "Logged out successfully"
}
```

### Проверка токена

#### GET /api/auth/me
Получить текущего пользователя (используется для SSO)

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "Иван",
      "last_name": "Иванов",
      "role": "user",
      "avatar": "/uploads/users/1/avatar.jpg"
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

#### GET /api/auth/verify
Проверить токен (для других сервисов)

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "Иван",
      "role": "user"
    }
  }
}
```

### Обновление токена

#### POST /api/auth/refresh
Обновить access token используя refresh token

**Request:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response:**
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### Управление паролем

#### POST /api/auth/forgot-password
Запросить сброс пароля

**Request:**
```json
{
  "email": "user@example.com"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password reset email sent"
}
```

#### POST /api/auth/reset-password
Сбросить пароль

**Request:**
```json
{
  "token": "reset_token_from_email",
  "new_password": "newpassword123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password reset successfully"
}
```

#### POST /api/auth/change-password
Изменить пароль (требует авторизации)

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "current_password": "oldpassword123",
  "new_password": "newpassword123"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

### Email верификация

#### POST /api/auth/send-verification
Отправить письмо подтверждения

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "message": "Verification email sent"
}
```

#### POST /api/auth/verify-email
Подтвердить email

**Request:**
```json
{
  "token": "verification_token_from_email"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Email verified successfully"
}
```

---

## Использование в других сервисах

### Проверка авторизации

Все сервисы должны проверять токены через Auth Service:

```go
// clinic/backend/middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "clinic/backend/pkg/clients"
)

var authClient = clients.NewAuthClient("http://localhost:7000")

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", 401)
            return
        }
        
        // Проверить через Auth API
        user, err := authClient.VerifyToken(token)
        if err != nil {
            http.Error(w, "Invalid token", 401)
            return
        }
        
        // Добавить в контекст
        ctx := context.WithValue(r.Context(), "userID", user.ID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## Конфигурация

### .env
```bash
# Порт сервера
PORT=7000

# База данных
DATABASE_URL=./database/auth.db

# JWT секрет
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=24h
REFRESH_TOKEN_EXPIRATION=168h

# Email (для сброса пароля и верификации)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASSWORD=your-password
SMTP_FROM=noreply@zooplatform.ru

# CORS
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:4000,http://localhost:4100
```

---

## Запуск

### Development
```bash
cd auth/backend
air
```

### Production
```bash
cd auth/backend
go build -o auth-server
./auth-server
```

### Docker
```bash
docker build -t auth-service -f infrastructure/docker/auth/Dockerfile.backend .
docker run -p 7000:7000 auth-service
```

---

## Swagger документация

После запуска сервера:
```
http://localhost:7000/swagger/index.html
```

Обновить документацию:
```bash
cd auth/backend
swag init
```

---

## Тестирование

```bash
# Регистрация
curl -X POST http://localhost:7000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123","name":"Test"}'

# Вход
curl -X POST http://localhost:7000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'

# Проверка токена
curl http://localhost:7000/api/auth/me \
  -H "Authorization: Bearer <token>"
```

---

## Зависимости

```bash
go get github.com/golang-jwt/jwt/v5
go get github.com/joho/godotenv
go get golang.org/x/crypto/bcrypt
go get github.com/mattn/go-sqlite3
go get github.com/swaggo/swag/cmd/swag
go get github.com/swaggo/http-swagger
```

---

## Безопасность

- ✅ Пароли хешируются с bcrypt
- ✅ JWT токены с expiration
- ✅ Refresh токены для обновления
- ✅ Rate limiting на login/register
- ✅ CORS защита
- ✅ SQL injection защита (prepared statements)
- ✅ XSS защита (sanitization)

---

**Дата обновления:** 17 января 2025
