# Middleware Package

Универсальный middleware для авторизации через Auth Service для всех микросервисов.

## Установка

```bash
go get github.com/zooplatforma/pkg/middleware
```

## Использование

### 1. Обязательная авторизация

```go
package main

import (
    "net/http"
    "github.com/zooplatforma/pkg/middleware"
)

func main() {
    // Защищенный endpoint
    http.Handle("/api/protected", middleware.AuthMiddleware(http.HandlerFunc(protectedHandler)))
    
    http.ListenAndServe(":8000", nil)
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
    // Получить user_id из контекста
    userID, ok := middleware.GetUserID(r)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // Получить email
    email, _ := middleware.GetUserEmail(r)
    
    // Получить роль
    role, _ := middleware.GetUserRole(r)
    
    // Использовать данные пользователя
    log.Printf("User %d (%s) with role %s accessed endpoint", userID, email, role)
}
```

### 2. Опциональная авторизация

```go
// Endpoint доступен всем, но если есть токен - добавляет данные в контекст
http.Handle("/api/public", middleware.OptionalAuthMiddleware(http.HandlerFunc(publicHandler)))

func publicHandler(w http.ResponseWriter, r *http.Request) {
    userID, ok := middleware.GetUserID(r)
    if ok {
        // Пользователь авторизован
        log.Printf("Authorized user %d", userID)
    } else {
        // Анонимный пользователь
        log.Printf("Anonymous user")
    }
}
```

### 3. Проверка роли superadmin

```go
// Только для суперадминов
http.Handle("/api/admin", 
    middleware.AuthMiddleware(
        middleware.RequireSuperAdmin(
            http.HandlerFunc(adminHandler),
        ),
    ),
)

func adminHandler(w http.ResponseWriter, r *http.Request) {
    // Только суперадмины попадут сюда
    userID, _ := middleware.GetUserID(r)
    log.Printf("Superadmin %d accessed admin panel", userID)
}
```

### 4. Прямая проверка токена

```go
import "github.com/zooplatforma/pkg/middleware"

func someFunction() {
    token := "eyJhbGciOiJIUzI1NiIs..."
    
    user, err := middleware.VerifyTokenViaAuthService(token)
    if err != nil {
        log.Printf("Invalid token: %v", err)
        return
    }
    
    log.Printf("User: %d - %s (%s)", user.ID, user.Name, user.Email)
}
```

### 5. Логин через Auth Service

```go
import "github.com/zooplatforma/pkg/middleware"

func loginHandler(w http.ResponseWriter, r *http.Request) {
    email := r.FormValue("email")
    password := r.FormValue("password")
    
    token, user, err := middleware.LoginViaAuthService(email, password)
    if err != nil {
        http.Error(w, "Login failed", http.StatusUnauthorized)
        return
    }
    
    // Установить cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "auth_token",
        Value:    token,
        Path:     "/",
        HttpOnly: true,
        MaxAge:   86400 * 7, // 7 days
    })
    
    // Вернуть данные пользователя
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": map[string]interface{}{
            "token": token,
            "user":  user,
        },
    })
}
```

## Переменные окружения

```bash
# URL Auth Service (по умолчанию http://localhost:7100)
AUTH_SERVICE_URL=http://localhost:7100
```

## Context Keys

Middleware добавляет следующие данные в контекст запроса:

- `UserIDKey` - ID пользователя (int)
- `UserEmailKey` - Email пользователя (string)
- `UserRoleKey` - Роль пользователя (string)

## Helper функции

### GetUserID(r *http.Request) (int, bool)
Получает user_id из контекста запроса.

### GetUserEmail(r *http.Request) (string, bool)
Получает email из контекста запроса.

### GetUserRole(r *http.Request) (string, bool)
Получает роль из контекста запроса.

### VerifyTokenViaAuthService(token string) (*User, error)
Проверяет токен через Auth Service и возвращает данные пользователя.

### LoginViaAuthService(email, password string) (string, *User, error)
Выполняет вход через Auth Service и возвращает токен и данные пользователя.

## Как это работает

1. Middleware получает токен из `Authorization` header или cookie `auth_token`
2. Отправляет запрос к Auth Service `/api/auth/verify` с токеном
3. Auth Service проверяет токен и возвращает данные пользователя
4. Middleware добавляет данные в контекст запроса
5. Handler получает доступ к данным через helper функции

## Преимущества

- ✅ Централизованная авторизация через Auth Service
- ✅ Единый JWT secret для всех сервисов
- ✅ Автоматическая проверка токенов
- ✅ Простая интеграция в любой микросервис
- ✅ Поддержка ролей (user, superadmin)
- ✅ Опциональная авторизация для публичных endpoints

## Миграция существующих сервисов

### До (старый middleware):

```go
// main/backend/middleware/auth.go
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Локальная проверка JWT
        token, err := jwt.Parse(...)
        // ...
    }
}
```

### После (новый middleware):

```go
import "github.com/zooplatforma/pkg/middleware"

// Просто используем готовый middleware
http.Handle("/api/protected", middleware.AuthMiddleware(http.HandlerFunc(handler)))
```

## Примеры из проекта

### Main Backend
```go
// main/backend/main.go
import "github.com/zooplatforma/pkg/middleware"

http.Handle("/api/posts", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetPosts)))
```

### Clinic Backend
```go
// clinic/backend/main.go
import "github.com/zooplatforma/pkg/middleware"

http.Handle("/api/patients", middleware.AuthMiddleware(http.HandlerFunc(handlers.GetPatients)))
```

### Admin Backend
```go
// admin/backend/main.go
import "github.com/zooplatforma/pkg/middleware"

http.Handle("/api/users", 
    middleware.AuthMiddleware(
        middleware.RequireSuperAdmin(
            http.HandlerFunc(handlers.GetUsers),
        ),
    ),
)
```

## Тестирование

```bash
# Получить токен
TOKEN=$(curl -X POST http://localhost:7100/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq -r '.data.token')

# Использовать токен
curl http://localhost:8000/api/protected \
  -H "Authorization: Bearer $TOKEN"

# Или через cookie
curl http://localhost:8000/api/protected \
  --cookie "auth_token=$TOKEN"
```
