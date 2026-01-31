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


---

## ⚠️ ТИПИЧНЫЕ ОШИБКИ И ИХ РЕШЕНИЯ

### Ошибка 1: "Выкидывает после обновления страницы"

**Симптомы:**
- Пользователь входит в систему успешно
- После обновления страницы (F5) пользователь не авторизован
- В консоли браузера множество 401 (Unauthorized) ошибок
- В логах Main Backend: `❌ Cookie not found: http: named cookie not present`

**Причина:**
Frontend отправлял запросы авторизации напрямую к Auth Service (7100), а не через Main Backend (8000). Cookie устанавливалась Auth Service для порта 7100, но запросы к Main Backend (8000) не включали эту cookie, так как браузер не передаёт cookies между разными портами.

**Диагностика:**
```bash
# Проверь логи Main Backend
tail -50 /tmp/main-backend.log | grep -E "(Cookie|Token)"

# Если видишь много "Cookie not found" - это проблема
```

**Решение:**

#### Шаг 1: Frontend должен использовать Main Backend для авторизации

**Файл:** `main/frontend/lib/api.ts`

**НЕПРАВИЛЬНО (было):**
```typescript
// API методы для авторизации (используют Auth Service на порту 7100)
export const authApi = {
  register: (name: string, email: string, password: string) =>
    authClient.post<{ user: User }>('/api/auth/register', { name, email, password }),
  
  login: (email: string, password: string) =>
    authClient.post<{ user: User }>('/api/auth/login', { email, password }),
  
  logout: () =>
    authClient.post<{ message: string }>('/api/auth/logout', {}),
  
  me: () =>
    authClient.get<User>('/api/auth/me'),
};
```

**ПРАВИЛЬНО (стало):**
```typescript
// API методы для авторизации (используют Main Backend который проксирует к Auth Service)
export const authApi = {
  register: (name: string, email: string, password: string) =>
    apiClient.post<{ user: User }>('/api/auth/register', { name, email, password }),
  
  login: (email: string, password: string) =>
    apiClient.post<{ user: User }>('/api/auth/login', { email, password }),
  
  logout: () =>
    apiClient.post<{ message: string }>('/api/auth/logout', {}),
  
  me: () =>
    apiClient.get<User>('/api/auth/me'),
};
```

**Изменение:** Заменили `authClient` (7100) на `apiClient` (8000)

#### Шаг 2: Cookie должна быть установлена для всего localhost

**Файл:** `main/backend/handlers/auth.go`

**НЕПРАВИЛЬНО (было):**
```go
http.SetCookie(w, &http.Cookie{
    Name:     "auth_token",
    Value:    authResp.Data.Token,
    Path:     "/",
    Domain:   "", // ❌ Пустой Domain - cookie только для текущего порта
    HttpOnly: true,
    Secure:   false,
    SameSite: http.SameSiteLaxMode,
    MaxAge:   86400 * 7,
})
```

**ПРАВИЛЬНО (стало):**
```go
http.SetCookie(w, &http.Cookie{
    Name:     "auth_token",
    Value:    authResp.Data.Token,
    Path:     "/",
    Domain:   "localhost", // ✅ Cookie работает для всех портов localhost
    HttpOnly: true,
    Secure:   false,
    SameSite: http.SameSiteLaxMode,
    MaxAge:   86400 * 7,
})
```

**Изменение:** Установили `Domain: "localhost"` вместо пустой строки

**Как применить изменение:**
```bash
cd main/backend
sed -i '' 's/Domain:   ""/Domain:   "localhost"/g' handlers/auth.go
```

#### Шаг 3: Проверка исправления

1. Очисти все cookies в браузере:
   - DevTools (F12) → Application → Cookies → Clear all

2. Перезапусти сервисы:
   ```bash
   cd main
   ./run
   ```

3. Войди в систему через http://localhost:3000

4. Проверь cookie в DevTools:
   - Application → Cookies → localhost
   - Должна быть cookie `auth_token` с Domain: `localhost`

5. Обнови страницу (F5)
   - Должен остаться авторизованным ✅

6. Проверь логи:
   ```bash
   tail -20 /tmp/main-backend.log
   ```
   - Должны быть `✅ Cookie found` вместо `❌ Cookie not found`

---

### Ошибка 2: "Middleware не может проверить токен"

**Симптомы:**
- Cookie приходит на Backend
- В логах: `❌ Token parse error` или `❌ Token invalid`
- Все запросы возвращают 401

**Причина:**
JWT токен создан Auth Service с одним секретом, а Main Backend пытается проверить его с другим секретом.

**Решение:**
Убедись что `JWT_SECRET` одинаковый в обоих сервисах:

```bash
# Проверь Auth Service
cat auth/backend/.env | grep JWT_SECRET

# Проверь Main Backend
cat main/backend/.env | grep JWT_SECRET

# Должны быть одинаковые!
```

Если разные - скопируй секрет из Auth Service в Main Backend.

---

### Ошибка 3: "CORS блокирует запросы"

**Симптомы:**
- В консоли браузера: `CORS policy: No 'Access-Control-Allow-Origin' header`
- Запросы не доходят до Backend

**Решение:**
Проверь CORS настройки в `main/backend/main.go`:

```go
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        
        // ✅ Разрешённые origins
        allowedOrigins := map[string]bool{
            "http://localhost:3000": true, // Main frontend
        }
        
        if allowedOrigins[origin] {
            w.Header().Set("Access-Control-Allow-Origin", origin)
        }
        
        // ✅ ВАЖНО: разрешаем credentials (cookies)
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
        
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next(w, r)
    }
}
```

---

## 📋 Checklist: Правильная настройка авторизации

### Frontend (`main/frontend/lib/api.ts`):
- [ ] ✅ `authApi` использует `apiClient` (порт 8000), НЕ `authClient` (порт 7100)
- [ ] ✅ Все fetch запросы имеют `credentials: 'include'`

### Backend (`main/backend/handlers/auth.go`):
- [ ] ✅ Cookie устанавливается с `Domain: "localhost"`
- [ ] ✅ Cookie имеет `HttpOnly: true` и `SameSite: http.SameSiteLaxMode`
- [ ] ✅ Handlers проксируют запросы к Auth Service (7100)

### Backend (`main/backend/main.go`):
- [ ] ✅ CORS настроен с `Access-Control-Allow-Credentials: true`
- [ ] ✅ Роуты `/api/auth/*` используют handlers (НЕ middleware)

### Backend (`main/backend/middleware/auth.go`):
- [ ] ✅ Middleware читает cookie `auth_token`
- [ ] ✅ Middleware парсит JWT с правильным секретом
- [ ] ✅ Middleware добавляет `userID` в контекст

### Environment (`.env`):
- [ ] ✅ `JWT_SECRET` одинаковый в Auth Service и Main Backend

---

## 🔍 Отладка проблем с авторизацией

### 1. Проверь cookie в браузере
```
DevTools (F12) → Application → Cookies → localhost
```
- Должна быть cookie `auth_token`
- Domain должен быть `localhost` (не `localhost:8000`)
- Path должен быть `/`

### 2. Проверь что cookie отправляется
```
DevTools (F12) → Network → Выбери любой запрос к localhost:8000
→ Request Headers → должен быть Cookie: auth_token=...
```

### 3. Проверь логи Backend
```bash
tail -f /tmp/main-backend.log
```
Должны быть:
- `✅ Cookie found: ...` (не `❌ Cookie not found`)
- `✅ Token valid, userID=...` (не `❌ Token parse error`)

### 4. Проверь что JWT секреты совпадают
```bash
diff <(cat auth/backend/.env | grep JWT_SECRET) <(cat main/backend/.env | grep JWT_SECRET)
```
Не должно быть различий!

---

## 💡 Почему это работает

### Правильный flow авторизации:

```
1. Frontend (3000) → Main Backend (8000) POST /api/auth/login
   ↓
2. Main Backend → Auth Service (7100) POST /api/auth/login
   ↓
3. Auth Service проверяет credentials, создаёт JWT токен
   ↓
4. Auth Service → Main Backend (возвращает токен)
   ↓
5. Main Backend устанавливает cookie для Domain: "localhost"
   ↓
6. Frontend получает cookie (автоматически сохраняется браузером)
   ↓
7. Frontend (3000) → Main Backend (8000) GET /api/posts
   ↓ (браузер автоматически добавляет cookie в запрос)
8. Main Backend читает cookie, проверяет JWT, добавляет userID в контекст
   ↓
9. Handler получает userID из контекста и обрабатывает запрос
```

### Почему cookie работает между портами:

- Cookie установлена с `Domain: "localhost"` (без порта)
- Браузер отправляет эту cookie на ВСЕ запросы к `localhost:*`
- Поэтому cookie работает для:
  - `localhost:3000` (Frontend)
  - `localhost:8000` (Main Backend)
  - `localhost:7100` (Auth Service)
  - И любых других портов localhost

### Почему НЕ работало раньше:

- Frontend отправлял запросы напрямую к Auth Service (7100)
- Auth Service устанавливал cookie для своего порта
- Браузер НЕ отправлял эту cookie на запросы к Main Backend (8000)
- Main Backend не видел cookie → 401 Unauthorized

---

## ⚠️ ЗАПОМНИ:

**Всегда используй Main Backend как прокси для авторизации!**

- ✅ Frontend → Main Backend → Auth Service
- ❌ Frontend → Auth Service (напрямую)

**Cookie должна быть установлена с Domain: "localhost"!**

- ✅ `Domain: "localhost"` - работает для всех портов
- ❌ `Domain: ""` - работает только для текущего порта
- ❌ `Domain: "localhost:8000"` - работает только для порта 8000


---

## ⚠️ ТИПИЧНАЯ ОШИБКА 2: "404 Not Found для /api/organizations/X/members"

**Симптомы:**
- В консоли браузера: `GET http://localhost:8000/api/organizations/4/members 404 (Not Found)`
- Пост от организации не показывает кнопки редактирования даже для членов организации
- Frontend пытается загрузить список членов организации через API

**Причина:**
Frontend пытался проверить права на редактирование постов организации, запрашивая список членов через API. Это неправильная архитектура:
- ❌ Проверка прав на Frontend требует загрузки дополнительных данных
- ❌ Нужен запущенный сервис организаций
- ❌ Дополнительные HTTP запросы замедляют работу
- ❌ Логика прав размазана между Frontend и Backend

**Правильная архитектура:**
- ✅ Проверка прав должна быть на Backend
- ✅ Backend возвращает поле `can_edit` для каждого поста
- ✅ Frontend просто использует это поле
- ✅ Проверка членства в организации делается в БД, а не через API

**Решение:**

#### Шаг 1: Добавить поле `CanEdit` в модель Post

**Файл:** `main/backend/models/post.go`

```go
type Post struct {
	ID            int           `json:"id"`
	AuthorID      int           `json:"author_id"`
	AuthorType    string        `json:"author_type"`
	Content       string        `json:"content"`
	// ... другие поля ...
	CanEdit       bool          `json:"can_edit"` // ✅ Может ли текущий пользователь редактировать пост
}
```

#### Шаг 2: Создать функцию проверки прав

**Файл:** `main/backend/handlers/posts.go`

```go
// checkCanEditPost проверяет может ли пользователь редактировать пост
func checkCanEditPost(userID int, post *models.Post) bool {
	if userID == 0 {
		return false
	}

	// Если пост от пользователя - проверяем ID
	if post.AuthorType == "user" && post.AuthorID == userID {
		return true
	}

	// Если пост от организации - проверяем членство с правами
	if post.AuthorType == "organization" {
		var role string
		err := database.DB.QueryRow(`
			SELECT role FROM organization_members 
			WHERE organization_id = ? AND user_id = ?
		`, post.AuthorID, userID).Scan(&role)

		if err == nil && (role == "owner" || role == "admin" || role == "moderator") {
			return true
		}
	}

	return false
}
```

#### Шаг 3: Вызывать проверку во всех функциях загрузки постов

**Файл:** `main/backend/handlers/posts.go`

```go
func getAllPosts(w http.ResponseWriter, r *http.Request) {
	// ... загрузка постов ...
	
	// ✅ Проверяем права на редактирование для каждого поста
	for i := range posts {
		posts[i].CanEdit = checkCanEditPost(userID, &posts[i])
	}

	sendSuccessResponse(w, posts)
}

func getUserPosts(w http.ResponseWriter, r *http.Request, userID int) {
	// ... загрузка постов ...
	
	// ✅ Проверяем права на редактирование для каждого поста
	for i := range posts {
		posts[i].CanEdit = checkCanEditPost(currentUserID, &posts[i])
	}

	sendSuccessResponse(w, posts)
}

func getPetPosts(w http.ResponseWriter, r *http.Request, petID int) {
	// ... загрузка постов ...
	
	// ✅ Проверяем права на редактирование для каждого поста
	for i := range posts {
		posts[i].CanEdit = checkCanEditPost(currentUserID, &posts[i])
	}

	sendSuccessResponse(w, posts)
}

func getOrganizationPosts(w http.ResponseWriter, r *http.Request, orgID int) {
	// ... загрузка постов ...
	
	// ✅ Проверяем права на редактирование для каждого поста
	for i := range posts {
		posts[i].CanEdit = checkCanEditPost(currentUserID, &posts[i])
	}

	sendSuccessResponse(w, posts)
}
```

#### Шаг 4: Обновить Frontend для использования `can_edit`

**Файл:** `main/frontend/app/components/posts/PostCard.tsx`

**НЕПРАВИЛЬНО (было):**
```typescript
interface Post {
  id: number;
  // ... другие поля ...
}

export default function PostCard({ post, onDelete }: PostCardProps) {
  const [canEditPost, setCanEditPost] = useState(false);

  // ❌ Проверка прав на Frontend через API
  useEffect(() => {
    const checkEditPermission = async () => {
      if (!user) {
        setCanEditPost(false);
        return;
      }

      if (post.author_type === 'user' && user.id === post.author_id) {
        setCanEditPost(true);
        return;
      }

      // ❌ Запрос к API организаций
      if (post.author_type === 'organization') {
        const response = await fetch(`http://localhost:8000/api/organizations/${post.author_id}/members`);
        // ... проверка членства ...
      }
    };

    checkEditPermission();
  }, [user, post]);
}
```

**ПРАВИЛЬНО (стало):**
```typescript
interface Post {
  id: number;
  // ... другие поля ...
  can_edit?: boolean; // ✅ Добавлено поле can_edit из Backend
}

export default function PostCard({ post, onDelete }: PostCardProps) {
  // ✅ Используем can_edit из Backend вместо локальной проверки
  const canEditPost = post.can_edit || false;
  
  // Больше не нужен useEffect для проверки прав!
}
```

**Изменения:**
1. Добавили поле `can_edit` в интерфейс `Post`
2. Убрали `useState` для `canEditPost`
3. Убрали `useEffect` с проверкой прав
4. Просто используем `post.can_edit` из Backend

#### Шаг 5: Проверка исправления

1. Перезапусти сервисы:
   ```bash
   cd main
   ./run
   ```

2. Открой консоль браузера (F12) → Network

3. Зайди на главную страницу http://localhost:3000

4. Проверь что НЕТ запросов к `/api/organizations/*/members` ✅

5. Проверь что посты от организаций показывают кнопки редактирования для членов с правами ✅

6. Проверь логи Backend:
   ```bash
   tail -20 /tmp/main-backend.log
   ```
   - Должны быть запросы к `organization_members` в БД, но НЕ через HTTP API

---

### Почему это правильное решение:

**Проблема старого подхода:**
```
Frontend → Main Backend → Organizations Service (8200) → Database
```
- Требует запущенный сервис организаций
- Дополнительный HTTP запрос
- Медленно
- Сложная логика на Frontend

**Новый подход:**
```
Frontend → Main Backend → Database (organization_members)
```
- Не требует сервис организаций
- Один запрос вместо двух
- Быстро
- Простая логика: Backend возвращает `can_edit`, Frontend использует

**Преимущества:**
- ✅ Меньше зависимостей (не нужен сервис организаций)
- ✅ Быстрее (один запрос вместо нескольких)
- ✅ Проще (логика прав в одном месте - на Backend)
- ✅ Безопаснее (проверка прав на Backend, а не на Frontend)
- ✅ Масштабируемее (легко добавить новые правила)

---

## 📋 Checklist: Правильная проверка прав на редактирование

### Backend (`main/backend/models/post.go`):
- [ ] ✅ Добавлено поле `CanEdit bool` в структуру `Post`

### Backend (`main/backend/handlers/posts.go`):
- [ ] ✅ Создана функция `checkCanEditPost(userID int, post *models.Post) bool`
- [ ] ✅ Функция проверяет `author_type == "user" && author_id == userID`
- [ ] ✅ Функция проверяет членство в организации через БД
- [ ] ✅ Функция вызывается в `getAllPosts`
- [ ] ✅ Функция вызывается в `getUserPosts`
- [ ] ✅ Функция вызывается в `getPetPosts`
- [ ] ✅ Функция вызывается в `getOrganizationPosts`

### Frontend (`main/frontend/app/components/posts/PostCard.tsx`):
- [ ] ✅ Добавлено поле `can_edit?: boolean` в интерфейс `Post`
- [ ] ✅ Убран `useState` для `canEditPost`
- [ ] ✅ Убран `useEffect` с проверкой прав
- [ ] ✅ Используется `post.can_edit` из Backend

### Проверка:
- [ ] ✅ Нет запросов к `/api/organizations/*/members` в консоли
- [ ] ✅ Посты от пользователей показывают кнопки редактирования для автора
- [ ] ✅ Посты от организаций показывают кнопки редактирования для членов с правами

---

## ⚠️ ЗАПОМНИ:

**Проверка прав должна быть на Backend, а не на Frontend!**

- ✅ Backend возвращает `can_edit` для каждого поста
- ✅ Frontend просто использует это поле
- ❌ НЕ делай запросы к другим сервисам для проверки прав
- ❌ НЕ проверяй права на Frontend
