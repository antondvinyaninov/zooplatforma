# Правильная логика авторизации через Auth Service (7100)

## ⚠️ КРИТИЧЕСКОЕ ПРАВИЛО: Единая точка авторизации

**ВСЕ микросервисы используют Auth Service (порт 7100) для авторизации!**

**НИКОГДА не используй локальную авторизацию в микросервисах!**

---

## Архитектура авторизации

```
User → Frontend → Auth Service (7100) → JWT Token
                        ↓
                  Cookie: auth_token
                        ↓
        ┌───────────────┴───────────────┐
        ↓                               ↓
   Main Backend (8000)          Admin Backend (9000)
   использует pkg/middleware    использует pkg/middleware
        ↓                               ↓
   Проверяет токен через        Проверяет токен через
   Auth Service (7100)          Auth Service (7100)
```

---

## 1. Frontend: Авторизация пользователя

### Страница логина (Main Frontend)

**Файл:** `main/frontend/app/(auth)/login/page.tsx`

```typescript
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      // ✅ ОБЯЗАТЕЛЬНО: используем Auth Service (7100)
      const response = await fetch('http://localhost:7100/api/auth/login', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include', // ВАЖНО: отправляет и получает cookies
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();

      if (!response.ok) {
        setError(data.error || 'Ошибка входа');
        setLoading(false);
        return;
      }

      if (data.success) {
        // Cookie auth_token установлен автоматически
        // Перенаправляем на главную
        router.push('/');
      } else {
        setError(data.error || 'Ошибка входа');
      }
    } catch (error) {
      console.error('Login error:', error);
      setError('Ошибка сети. Проверьте подключение.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full p-8 bg-white rounded-lg shadow-lg">
        <h1 className="text-2xl font-bold mb-6">Вход в систему</h1>
        
        {error && (
          <div className="mb-4 p-3 bg-red-50 border border-red-200 rounded text-red-700">
            {error}
          </div>
        )}

        <form onSubmit={handleLogin}>
          <div className="mb-4">
            <label className="block text-sm font-medium mb-2">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-3 py-2 border rounded-lg"
              required
            />
          </div>

          <div className="mb-6">
            <label className="block text-sm font-medium mb-2">Пароль</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-3 py-2 border rounded-lg"
              required
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
          >
            {loading ? 'Вход...' : 'Войти'}
          </button>
        </form>
      </div>
    </div>
  );
}
```

---

## 2. Frontend: Проверка авторизации на защищенных страницах

### Любая защищенная страница (например, главная)

**Файл:** `main/frontend/app/(main)/page.tsx`

```typescript
'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';

export default function HomePage() {
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
        credentials: 'include', // ВАЖНО: отправляет cookie
      });

      if (!response.ok) {
        // Не авторизован - редирект на логин
        router.push('/login');
        return;
      }

      const data = await response.json();
      
      if (data.success && data.data) {
        setUser(data.data.user);
        setLoading(false);
      } else {
        router.push('/login');
      }
    } catch (error) {
      console.error('Auth check failed:', error);
      router.push('/login');
    }
  };

  if (loading) {
    return <div>Загрузка...</div>;
  }

  return (
    <div>
      <h1>Добро пожаловать, {user?.name}!</h1>
      {/* Контент главной страницы */}
    </div>
  );
}
```

---

## 3. Backend: Использование pkg/middleware

### Main Backend - правильная настройка

**Файл:** `main/backend/main.go`

```go
package main

import (
    "database"
    "log"
    "net/http"
    
    // ✅ ОБЯЗАТЕЛЬНО: импортируй middleware из pkg
    "github.com/zooplatforma/pkg/middleware"
    
    "main/backend/handlers"
)

func main() {
    // Инициализация БД
    if err := database.InitDB(); err != nil {
        log.Fatal("Failed to initialize database:", err)
    }
    defer database.CloseDB()
    
    // ✅ ОБЯЗАТЕЛЬНО: Инициализируй AuthMiddleware с URL Auth Service
    middleware.InitAuthMiddleware("http://localhost:7100")
    
    // CORS middleware
    http.HandleFunc("/", enableCORS(handleRoot))
    
    // Публичные роуты (без авторизации)
    http.HandleFunc("/api/health", enableCORS(handleHealth))
    
    // ✅ Защищенные роуты - используй middleware.AuthMiddleware
    http.HandleFunc("/api/posts", enableCORS(middleware.AuthMiddleware(handlers.HandlePosts)))
    http.HandleFunc("/api/posts/", enableCORS(middleware.AuthMiddleware(handlers.HandlePostDetail)))
    http.HandleFunc("/api/profile", enableCORS(middleware.AuthMiddleware(handlers.HandleProfile)))
    http.HandleFunc("/api/friends", enableCORS(middleware.AuthMiddleware(handlers.HandleFriends)))
    
    port := ":8000"
    log.Printf("🚀 Main Backend started on port %s", port)
    log.Fatal(http.ListenAndServe(port, nil))
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next(w, r)
    }
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Main Backend API"))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"success":true,"status":"healthy"}`))
}
```

---

## 4. Backend: Использование user_id в handlers

### Получение user_id из контекста

**Файл:** `main/backend/handlers/posts.go`

```go
package handlers

import (
    "database"
    "encoding/json"
    "log"
    "net/http"
)

func HandlePosts(w http.ResponseWriter, r *http.Request) {
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
    
    switch r.Method {
    case "GET":
        getPosts(w, r, userID)
    case "POST":
        createPost(w, r, userID)
    default:
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
    }
}

func getPosts(w http.ResponseWriter, r *http.Request, userID int) {
    db := database.GetDB()
    
    rows, err := db.Query(`
        SELECT id, user_id, content, created_at 
        FROM posts 
        ORDER BY created_at DESC 
        LIMIT 20
    `)
    if err != nil {
        log.Printf("❌ Database error: %v", err)
        http.Error(w, `{"success":false,"error":"Database error"}`, http.StatusInternalServerError)
        return
    }
    defer rows.Close()
    
    var posts []map[string]interface{}
    for rows.Next() {
        var post map[string]interface{}
        // ... scan rows
        posts = append(posts, post)
    }
    
    response := map[string]interface{}{
        "success": true,
        "data": posts,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func createPost(w http.ResponseWriter, r *http.Request, userID int) {
    var req struct {
        Content string `json:"content"`
    }
    
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"success":false,"error":"Invalid JSON"}`, http.StatusBadRequest)
        return
    }
    
    db := database.GetDB()
    
    result, err := db.Exec(`
        INSERT INTO posts (user_id, content, created_at) 
        VALUES (?, ?, datetime('now'))
    `, userID, req.Content)
    
    if err != nil {
        log.Printf("❌ Database error: %v", err)
        http.Error(w, `{"success":false,"error":"Failed to create post"}`, http.StatusInternalServerError)
        return
    }
    
    postID, _ := result.LastInsertId()
    
    response := map[string]interface{}{
        "success": true,
        "data": map[string]interface{}{
            "id": postID,
            "user_id": userID,
            "content": req.Content,
        },
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

---

## 5. pkg/middleware - как это работает

### Структура pkg/middleware

**Файл:** `pkg/middleware/auth.go`

```go
package middleware

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
)

var authServiceURL string

// InitAuthMiddleware инициализирует middleware с URL Auth Service
func InitAuthMiddleware(url string) {
    authServiceURL = url
    log.Printf("🔐 Auth middleware initialized with URL: %s", url)
}

// AuthMiddleware проверяет JWT токен через Auth Service
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Получаем cookie auth_token
        cookie, err := r.Cookie("auth_token")
        if err != nil {
            log.Printf("⚠️ No auth_token cookie")
            http.Error(w, `{"success":false,"error":"Unauthorized"}`, http.StatusUnauthorized)
            return
        }
        
        // Проверяем токен через Auth Service
        client := &http.Client{}
        req, _ := http.NewRequest("GET", authServiceURL+"/api/auth/me", nil)
        req.AddCookie(cookie)
        
        resp, err := client.Do(req)
        if err != nil {
            log.Printf("❌ Auth Service error: %v", err)
            http.Error(w, `{"success":false,"error":"Auth service unavailable"}`, http.StatusServiceUnavailable)
            return
        }
        defer resp.Body.Close()
        
        if resp.StatusCode != http.StatusOK {
            log.Printf("⚠️ Auth failed: status %d", resp.StatusCode)
            http.Error(w, `{"success":false,"error":"Unauthorized"}`, http.StatusUnauthorized)
            return
        }
        
        // Парсим ответ от Auth Service
        var authResp struct {
            Success bool `json:"success"`
            Data    struct {
                User struct {
                    ID    int    `json:"id"`
                    Email string `json:"email"`
                    Role  string `json:"role"`
                } `json:"user"`
            } `json:"data"`
        }
        
        if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
            log.Printf("❌ Failed to parse auth response: %v", err)
            http.Error(w, `{"success":false,"error":"Internal error"}`, http.StatusInternalServerError)
            return
        }
        
        if !authResp.Success {
            http.Error(w, `{"success":false,"error":"Unauthorized"}`, http.StatusUnauthorized)
            return
        }
        
        // ✅ Добавляем данные пользователя в контекст
        ctx := context.WithValue(r.Context(), "user_id", authResp.Data.User.ID)
        ctx = context.WithValue(ctx, "user_email", authResp.Data.User.Email)
        ctx = context.WithValue(ctx, "user_role", authResp.Data.User.Role)
        
        log.Printf("✅ User authenticated: id=%d, email=%s", authResp.Data.User.ID, authResp.Data.User.Email)
        
        // Передаем управление следующему handler
        next(w, r.WithContext(ctx))
    }
}
```

---

## 6. Auth Service - что он делает

### Auth Service endpoints

**Файл:** `auth/backend/main.go`

```go
// POST /api/auth/login - Вход в систему
// - Проверяет email/password
// - Создает JWT токен
// - Устанавливает cookie auth_token
// - Возвращает данные пользователя

// GET /api/auth/me - Проверка токена
// - Читает cookie auth_token
// - Проверяет JWT токен
// - Возвращает данные пользователя
// - Используется всеми микросервисами через pkg/middleware

// POST /api/auth/logout - Выход
// - Удаляет cookie auth_token
```

---

## Checklist: Правильная авторизация

### Frontend:
- [ ] ✅ Логин через `http://localhost:7100/api/auth/login`
- [ ] ✅ Проверка через `http://localhost:7100/api/auth/me`
- [ ] ✅ Используется `credentials: 'include'` в fetch
- [ ] ✅ Редирект на `/login` если не авторизован

### Backend:
- [ ] ✅ Импортирован `github.com/zooplatforma/pkg/middleware`
- [ ] ✅ Вызван `middleware.InitAuthMiddleware("http://localhost:7100")`
- [ ] ✅ Middleware применен к защищенным роутам
- [ ] ✅ В handlers используется `r.Context().Value("user_id")`
- [ ] ✅ CORS настроен с `Access-Control-Allow-Credentials: true`

### Auth Service:
- [ ] ✅ Запущен на порту 7100
- [ ] ✅ Endpoint `/api/auth/me` работает
- [ ] ✅ Возвращает токен в ответе

---

## ⚠️ ЗАПОМНИ:

**Auth Service (7100) - единая точка авторизации!**

- ✅ Frontend → Auth Service (7100) → Cookie
- ✅ Backend → pkg/middleware → Auth Service (7100) → Context
- ❌ НИКОГДА не делай локальную авторизацию в микросервисах!
