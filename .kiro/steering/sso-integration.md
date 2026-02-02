---
inclusion: always
---

# SSO Integration - Правила интеграции авторизации

## 🚨 КРИТИЧЕСКОЕ ПРАВИЛО #1: НИКОГДА НЕ ИСПОЛЬЗУЙ ПОРТ 8000 ДЛЯ АВТОРИЗАЦИИ!

**⛔️ ЗАПРЕЩЕНО:**
- ❌ `http://localhost:8000/api/auth/login`
- ❌ `http://localhost:8000/api/auth/me`
- ❌ `apiClient.post('/api/auth/login', ...)`
- ❌ Main Backend (8000) для любых операций авторизации

**✅ ПРАВИЛЬНО:**
- ✅ `http://localhost:7100/api/auth/login`
- ✅ `http://localhost:7100/api/auth/me`
- ✅ `authClient.post('/api/auth/login', ...)` (использует порт 7100)
- ✅ Auth Service (7100) для ВСЕХ операций авторизации

## ⚠️ КРИТИЧЕСКОЕ ПРАВИЛО #2: Единая система авторизации

**Все микросервисы используют SSO (Single Sign-On) через Auth Service (порт 7100)!**

**НИКОГДА НЕ ИСПОЛЬЗУЙ Main Backend (8000) для авторизации!**

## Архитектура SSO

```
User → Main Frontend (3000) → Auth Service (7100) → JWT Token
                                      ↓
                              Cookie: auth_token
                                      ↓
                    ┌─────────────────┴─────────────────┐
                    ↓                 ↓                 ↓
            Admin (9000)      Clinic (8600)      Owner (8400)
            Shelter (8200)    Volunteer (8500)   PetBase (8100)
            Main (8000)
```

**ВАЖНО:** Main Backend (8000) НЕ участвует в процессе авторизации!

## 📝 Примеры: ПРАВИЛЬНО vs НЕПРАВИЛЬНО

### Frontend: API клиент (main/frontend/lib/api.ts)

**❌ НЕПРАВИЛЬНО:**
```typescript
// НЕ ДЕЛАЙ ТАК!
export const authApi = {
  login: (email: string, password: string) =>
    apiClient.post('/api/auth/login', { email, password }), // ❌ apiClient = порт 8000
  
  me: () =>
    apiClient.get('/api/auth/me'), // ❌ apiClient = порт 8000
};
```

**✅ ПРАВИЛЬНО:**
```typescript
// ДЕЛАЙ ТАК!
export const authApi = {
  login: (email: string, password: string) =>
    authClient.post('/api/auth/login', { email, password }), // ✅ authClient = порт 7100
  
  me: () =>
    authClient.get('/api/auth/me'), // ✅ authClient = порт 7100
};
```

### Frontend: Прямые fetch запросы

**❌ НЕПРАВИЛЬНО:**
```typescript
// НЕ ДЕЛАЙ ТАК!
const response = await fetch('http://localhost:8000/api/auth/login', { // ❌ порт 8000
  method: 'POST',
  body: JSON.stringify({ email, password }),
});
```

**✅ ПРАВИЛЬНО:**
```typescript
// ДЕЛАЙ ТАК!
const response = await fetch('http://localhost:7100/api/auth/login', { // ✅ порт 7100
  method: 'POST',
  credentials: 'include', // ✅ ВАЖНО для cookies
  body: JSON.stringify({ email, password }),
});
```

### Backend: Проверка авторизации

**❌ НЕПРАВИЛЬНО:**
```go
// НЕ ДЕЛАЙ ТАК!
func handlePets(w http.ResponseWriter, r *http.Request) {
    // ❌ Локальная проверка JWT
    tokenString := r.Header.Get("Authorization")
    token, err := jwt.Parse(tokenString, ...)
    // ...
}
```

**✅ ПРАВИЛЬНО:**
```go
// ДЕЛАЙ ТАК!
import "github.com/zooplatforma/pkg/middleware"

func main() {
    // ✅ Инициализация middleware с Auth Service
    middleware.InitAuthMiddleware("http://localhost:7100")
    
    // ✅ Применение middleware к защищенным роутам
    http.HandleFunc("/api/pets", enableCORS(middleware.AuthMiddleware(handlePets)))
}

func handlePets(w http.ResponseWriter, r *http.Request) {
    // ✅ Получение user_id из контекста (установлен middleware)
    userID, ok := r.Context().Value("user_id").(int)
    if !ok {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    // ...
}
```

## Как работает SSO

### 1. Авторизация (Auth Service - 7100)

**Endpoint:** `POST http://localhost:7100/api/auth/login`

**⚠️ НЕ `http://localhost:8000/api/auth/login`!**

**Процесс:**
1. Пользователь вводит email/password на Main Frontend
2. **Frontend отправляет запрос НАПРЯМУЮ к Auth Service (7100)**
3. Auth Service проверяет credentials в `auth/backend/auth.db`
4. Создается JWT token с данными пользователя и ролями из `user_roles`
5. Token сохраняется в cookie `auth_token` для домена `localhost`
6. Frontend получает данные пользователя

**JWT Payload:**
```go
{
  "user_id": 1,
  "email": "user@example.com",
  "role": "superadmin", // главная роль из user_roles
  "all_roles": ["superadmin", "user"], // все роли
  "exp": 1234567890
}
```

### 2. Проверка токена (Auth Service - 7100)

**Endpoint:** `GET http://localhost:7100/api/auth/me`

**Назначение:** Все микросервисы проверяют токен через этот endpoint

**Ответ:**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": 1,
      "email": "user@example.com",
      "name": "Антон",
      "role": "superadmin",
      "all_roles": ["superadmin", "user"]
    },
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

## Правила для Backend микросервисов

### ⚠️ ОБЯЗАТЕЛЬНО: Используй pkg/middleware

**ВСЕ backend микросервисы ДОЛЖНЫ использовать `pkg/middleware.AuthMiddleware`!**

**НЕ создавай свой middleware! Используй готовый из `pkg/`!**

### Подключение pkg/middleware

**Файл:** `<service>/backend/main.go`

```go
package main

import (
    "database"
    "log"
    "net/http"
    
    // ✅ ОБЯЗАТЕЛЬНО: импортируй middleware из pkg
    "github.com/zooplatforma/pkg/middleware"
)

func main() {
    // Инициализация БД
    if err := database.InitDB(); err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer database.CloseDB()
    
    // ✅ ОБЯЗАТЕЛЬНО: Инициализируй AuthMiddleware с URL Auth Service
    middleware.InitAuthMiddleware("http://localhost:7100")
    
    // Публичные роуты (без авторизации)
    http.HandleFunc("/api/health", enableCORS(handleHealth))
    
    // ✅ Защищенные роуты - используй middleware.AuthMiddleware
    http.HandleFunc("/api/pets", enableCORS(middleware.AuthMiddleware(handlePets)))
    http.HandleFunc("/api/pets/", enableCORS(middleware.AuthMiddleware(handlePetDetail)))
    
    port := ":8100"
    log.Printf("🚀 Service started on port %s", port)
    log.Fatal(http.ListenAndServe(port, nil))
}
```

### Использование в handlers

**Файл:** `<service>/backend/handlers/pets.go`

```go
package handlers

import (
    "net/http"
    "log"
    
    "github.com/zooplatforma/pkg/middleware"
)

func CreatePet(w http.ResponseWriter, r *http.Request) {
    // ✅ Получаем user_id из контекста (установлен middleware)
    userID, ok := r.Context().Value("user_id").(int)
    if !ok || userID == 0 {
        http.Error(w, `{"success":false,"error":"Unauthorized"}`, http.StatusUnauthorized)
        return
    }
    
    // ✅ Можно также получить email и role
    email, _ := r.Context().Value("user_email").(string)
    role, _ := r.Context().Value("user_role").(string)
    
    log.Printf("✅ User authenticated: id=%d, email=%s, role=%s", userID, email, role)
    
    // Твоя логика...
}
```

### Context Keys (из pkg/middleware)

```go
// Эти ключи устанавливаются автоматически pkg/middleware.AuthMiddleware
"user_id"    // int - ID пользователя
"user_email" // string - Email пользователя  
"user_role"  // string - Главная роль пользователя
```

## Правила для Frontend микросервисов

### ⚠️ ВСЕГДА используй Auth Service (7100)

**Паттерн для всех страниц:**

```typescript
'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export default function ProtectedPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState(null);

  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      // ✅ ОБЯЗАТЕЛЬНО: проверяем через Auth Service (7100)
      const response = await fetch('http://localhost:7100/api/auth/me', {
        credentials: 'include', // Отправляет cookie
      });

      if (!response.ok) {
        // Не авторизован - редирект на главный сайт
        window.location.href = 'http://localhost:3000';
        return;
      }

      const data = await response.json();
      if (data.success && data.data) {
        setUser(data.data.user);
        setLoading(false);
      } else {
        window.location.href = 'http://localhost:3000';
      }
    } catch (error) {
      console.error('Auth check failed:', error);
      window.location.href = 'http://localhost:3000';
    }
  };

  if (loading) {
    return <div>Загрузка...</div>;
  }

  return (
    <div>
      {/* Защищенный контент */}
    </div>
  );
}
```

### API запросы с авторизацией

**Проблема:** Cookie `auth_token` создан для `localhost:3000`, но не доступен для `localhost:6300`

**Решение:** Получить токен через Auth Service и передавать в заголовке

```typescript
// ✅ Helper функция для получения токена
const getAuthToken = async () => {
  try {
    // ОБЯЗАТЕЛЬНО: используй Auth Service (7100)
    const response = await fetch('http://localhost:7100/api/auth/me', {
      credentials: 'include',
    });
    
    if (response.ok) {
      const data = await response.json();
      return data.data?.token || null;
    }
  } catch (error) {
    console.error('Failed to get token:', error);
  }
  return null;
};

// Использование в API запросах
const fetchData = async () => {
  const token = await getAuthToken();
  
  if (!token) {
    console.error('No auth token');
    return;
  }
  
  const response = await fetch('http://localhost:8600/api/data', {
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
  });
  
  // ...
};
```

## Проверка ролей

### Суперадмин (superadmin)

**Доступ к:** Admin, PetBase

**Проверка в frontend:**
```typescript
useEffect(() => {
  checkAuth();
}, []);

const checkAuth = async () => {
  // ✅ ОБЯЗАТЕЛЬНО: Auth Service (7100)
  const response = await fetch('http://localhost:7100/api/auth/me', {
    credentials: 'include',
  });
  
  const data = await response.json();
  
  if (data.data?.user?.role !== 'superadmin') {
    alert('Доступ запрещен. Требуются права суперадмина.');
    window.location.href = 'http://localhost:3000';
    return;
  }
  
  setUser(data.data.user);
};
```

## CORS настройки

### Auth Service (7100)

**Файл:** `auth/backend/main.go`

```go
func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        allowedOrigins := []string{
            "http://localhost:3000",  // Main Frontend
            "http://localhost:4000",  // Admin Frontend
            "http://localhost:4100",  // PetBase Frontend
            "http://localhost:5100",  // Shelter Frontend
            "http://localhost:6100",  // Owner Frontend
            "http://localhost:6200",  // Volunteer Frontend
            "http://localhost:6300",  // Clinic Frontend
        }
        
        origin := r.Header.Get("Origin")
        for _, allowed := range allowedOrigins {
            if origin == allowed {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                break
            }
        }
        
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

## Checklist для нового микросервиса

### Backend:

- [ ] ✅ Используется `pkg/middleware.AuthMiddleware` (НЕ свой!)
- [ ] ✅ Вызван `middleware.InitAuthMiddleware("http://localhost:7100")`
- [ ] ✅ Middleware применен к защищенным роутам
- [ ] ✅ В handlers используется `r.Context().Value("user_id")`
- [ ] ✅ Логирование авторизации для отладки

### Frontend:

- [ ] ✅ Проверка авторизации через `http://localhost:7100/api/auth/me`
- [ ] ✅ Редирект на `http://localhost:3000` если не авторизован
- [ ] ✅ Получение токена через Auth Service (7100)
- [ ] ✅ Передача токена в заголовке `Authorization: Bearer {token}`

### Auth Service:

- [ ] ✅ Добавлен origin нового микросервиса в CORS
- [ ] ✅ Endpoint `/api/auth/me` возвращает токен в ответе

## Troubleshooting

### Проблема: 401 Unauthorized

**Причины:**
1. Не используется Auth Service (7100)
2. Не вызван `middleware.InitAuthMiddleware`
3. Токен истек
4. Middleware не применен к роуту

**Решение:**
- ✅ Проверь что используется `http://localhost:7100`
- ✅ Проверь что вызван `middleware.InitAuthMiddleware("http://localhost:7100")`
- ✅ Добавь логирование в handler

### Проблема: CORS error

**Причины:**
1. Origin не добавлен в Auth Service CORS
2. Credentials не включены

**Решение:**
- ✅ Добавь origin в `auth/backend/main.go`
- ✅ Используй `credentials: 'include'` в fetch

## Примеры из проекта

### PetBase Backend (ПРАВИЛЬНЫЙ ПРИМЕР)
- Файл: `petbase/backend/main.go`
- Использует: `pkg/middleware.AuthMiddleware`
- Auth Service: `http://localhost:7100`

### Volunteer Frontend (ПРАВИЛЬНЫЙ ПРИМЕР)
- Файл: `volunteer/frontend/app/dashboard/page.tsx`
- Использует: `http://localhost:7100/api/auth/me`

---

## ⚠️ ЗАПОМНИ:

1. **ВСЕГДА** используй Auth Service (7100) для авторизации
2. **НИКОГДА** не используй Main Backend (8000) для авторизации
3. **ВСЕГДА** используй `pkg/middleware.AuthMiddleware` в backend
4. **НИКОГДА** не создавай свой middleware для авторизации

**Auth Service (7100) - это единая точка авторизации для ВСЕХ микросервисов!**
