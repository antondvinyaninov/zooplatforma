# Design Document: Auth & CORS Architecture

## Overview

This document describes the authentication and CORS architecture for the Main Service in ZooPlatforma. The system uses a microservices architecture with an API Gateway as the single entry point, handling authentication and CORS for all services.

## Architecture Diagrams

### Development Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Development Mode                          │
└─────────────────────────────────────────────────────────────┘

Frontend (localhost:3000)
    │
    │ fetch('http://localhost:8000/api/posts', {
    │   headers: { Authorization: 'Bearer TOKEN' },
    │   credentials: 'include'
    │ })
    │
    ↓
Backend (localhost:8000)
    │
    │ DevAuthMiddleware:
    │ 1. Читает токен из Authorization или cookie
    │ 2. Валидирует JWT локально
    │ 3. Добавляет userID в context
    │
    ↓
Handler
    │ userID := r.Context().Value("userID").(int)
    │ Обрабатывает запрос
    │
    ↓
PostgreSQL (88.218.121.213:5432)
```

**Особенности Development:**
- Frontend обращается напрямую к Backend (localhost:8000)
- Backend проверяет JWT самостоятельно через `DevAuthMiddleware`
- CORS настроен в Backend для localhost:3000
- Токен передается через `Authorization: Bearer TOKEN` или cookie
- `NEXT_PUBLIC_API_URL=http://localhost:8000`

### Production Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Production Mode                           │
└─────────────────────────────────────────────────────────────┘

Frontend (Easypanel)
    │
    │ fetch('/api/posts', {  // Относительный путь!
    │   credentials: 'include'
    │ })
    │
    ↓
Gateway (Easypanel)
    │
    │ CORSMiddleware:
    │ 1. Проверяет Origin
    │ 2. Устанавливает Access-Control-* заголовки
    │
    │ AuthMiddleware:
    │ 1. Читает cookie auth_token
    │ 2. Валидирует JWT
    │ 3. Добавляет X-User-ID, X-User-Email, X-User-Role
    │
    ↓
Backend (Easypanel)
    │
    │ AuthMiddleware:
    │ 1. Читает X-User-* заголовки от Gateway
    │ 2. Добавляет userID в context
    │ 3. НЕ проверяет JWT (Gateway уже проверил)
    │
    ↓
Handler
    │ userID := r.Context().Value("userID").(int)
    │ Обрабатывает запрос
    │
    ↓
PostgreSQL (88.218.121.213:5432)
```

**Особенности Production:**
- Frontend использует относительные пути (/api/...)
- Gateway проверяет JWT и добавляет X-User-* заголовки
- Backend читает заголовки (не проверяет JWT)
- CORS управляется только Gateway
- `NEXT_PUBLIC_API_URL=` (пустая строка)

---

## Authentication Flow

### Development Authentication Flow

```
┌──────────────────────────────────────────────────────────────┐
│                    Login Flow (Development)                   │
└──────────────────────────────────────────────────────────────┘

1. User enters credentials
   │
   ↓
2. Frontend: apiClient.post('/api/auth/login', { email, password })
   │ URL: http://localhost:8000/api/auth/login
   │
   ↓
3. Backend: handlers/auth.go
   │ - Проверяет credentials в БД
   │ - Создает JWT токен
   │ - Устанавливает cookie: auth_token
   │ - Возвращает: { success: true, data: { user, token } }
   │
   ↓
4. Frontend: сохраняет токен
   │ localStorage.setItem('auth_token', token)
   │
   ↓
5. Frontend: редирект на главную
   │
   ↓
6. Frontend: загружает данные
   │ apiClient.get('/api/posts')
   │ Headers: { Authorization: 'Bearer TOKEN' }
   │
   ↓
7. Backend: DevAuthMiddleware
   │ - Читает токен из Authorization
   │ - Валидирует JWT с JWT_SECRET
   │ - Извлекает user_id, email, role
   │ - Добавляет в context
   │
   ↓
8. Handler: получает userID из context
   │ userID := r.Context().Value("userID").(int)
   │ Возвращает данные пользователя
```

### Production Authentication Flow

```
┌──────────────────────────────────────────────────────────────┐
│                    Login Flow (Production)                    │
└──────────────────────────────────────────────────────────────┘

1. User enters credentials
   │
   ↓
2. Frontend: apiClient.post('/api/auth/login', { email, password })
   │ URL: /api/auth/login (относительный путь)
   │
   ↓
3. Gateway: /api/auth/login
   │ - Проверяет credentials в БД
   │ - Создает JWT токен
   │ - Устанавливает cookie: auth_token (Domain: .easypanel.host)
   │ - Возвращает: { success: true, data: { user } }
   │
   ↓
4. Frontend: cookie установлена автоматически
   │ Не нужно сохранять в localStorage
   │
   ↓
5. Frontend: редирект на главную
   │
   ↓
6. Frontend: загружает данные
   │ apiClient.get('/api/posts')
   │ URL: /api/posts (относительный путь)
   │ Cookie: auth_token автоматически добавляется браузером
   │
   ↓
7. Gateway: AuthMiddleware
   │ - Читает cookie auth_token
   │ - Валидирует JWT с JWT_SECRET
   │ - Извлекает user_id, email, role
   │ - Добавляет заголовки:
   │   X-User-ID: 123
   │   X-User-Email: user@example.com
   │   X-User-Role: user
   │
   ↓
8. Gateway: проксирует к Backend
   │ GET http://backend:8000/api/posts
   │ Headers: { X-User-ID: 123, X-User-Email: ..., X-User-Role: ... }
   │
   ↓
9. Backend: AuthMiddleware
   │ - Читает X-User-ID из заголовков
   │ - Добавляет в context
   │ - НЕ проверяет JWT (Gateway уже проверил)
   │
   ↓
10. Handler: получает userID из context
    │ userID := r.Context().Value("userID").(int)
    │ Возвращает данные пользователя
```

---

## CORS Architecture

### Why Gateway Manages CORS

**Problem:** Duplicate CORS headers cause browser errors

```
❌ НЕПРАВИЛЬНО:

Frontend → Gateway (устанавливает CORS) → Backend (тоже устанавливает CORS)

Response headers:
Access-Control-Allow-Origin: localhost:3000, localhost:3000
                             ↑ от Gateway    ↑ от Backend

Browser: ❌ CORS error! Duplicate header!
```

**Solution:** Gateway is the single source of CORS headers

```
✅ ПРАВИЛЬНО:

Frontend → Gateway (устанавливает CORS) → Backend (НЕ устанавливает CORS)

Response headers:
Access-Control-Allow-Origin: localhost:3000
                             ↑ только от Gateway

Browser: ✅ OK!
```

### CORS Flow

```
┌──────────────────────────────────────────────────────────────┐
│                    CORS Flow (Production)                     │
└──────────────────────────────────────────────────────────────┘

1. Browser: Preflight request (OPTIONS)
   │ OPTIONS /api/posts
   │ Origin: https://frontend.com
   │
   ↓
2. Gateway: CORSMiddleware
   │ - Проверяет Origin в allowedOrigins
   │ - Если разрешен:
   │   Access-Control-Allow-Origin: https://frontend.com
   │   Access-Control-Allow-Credentials: true
   │   Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
   │   Access-Control-Allow-Headers: Content-Type, Authorization, Cookie
   │   Access-Control-Max-Age: 3600
   │ - Возвращает 200 OK
   │
   ↓
3. Browser: ✅ Preflight OK, делает основной запрос
   │ GET /api/posts
   │ Origin: https://frontend.com
   │ Cookie: auth_token=...
   │
   ↓
4. Gateway: CORSMiddleware
   │ - Устанавливает CORS заголовки
   │
   ↓
5. Gateway: AuthMiddleware
   │ - Проверяет JWT
   │ - Добавляет X-User-* заголовки
   │
   ↓
6. Gateway: ProxyHandler
   │ - Проксирует к Backend
   │ GET http://backend:8000/api/posts
   │ Headers: { X-User-ID: 123, ... }
   │
   ↓
7. Backend: обрабатывает запрос
   │ - НЕ устанавливает CORS заголовки
   │ - Возвращает данные
   │
   ↓
8. Gateway: ProxyHandler
   │ - Копирует заголовки от Backend
   │ - ФИЛЬТРУЕТ Access-Control-* заголовки
   │ - Оставляет только CORS заголовки от Gateway
   │
   ↓
9. Browser: получает ответ
   │ Access-Control-Allow-Origin: https://frontend.com (только один!)
   │ ✅ OK!
```

### CORS Configuration

**Gateway: middleware/cors.go**

```go
func CORSMiddleware(next http.Handler) http.Handler {
    return func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        
        // Список разрешенных origins
        allowedOrigins := map[string]bool{
            "http://localhost:3000": true,  // Development
            "https://my-projects-zooplatforma.crv1ic.easypanel.host": true,  // Production
        }
        
        // Если origin разрешен - устанавливаем CORS заголовки
        if allowedOrigins[origin] {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Credentials", "true")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie, X-User-ID, X-User-Email, X-User-Role")
            w.Header().Set("Access-Control-Max-Age", "3600")
            
            // Preflight request
            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusOK)
                return
            }
        } else if origin != "" && r.Method == "OPTIONS" {
            // Origin не разрешен
            w.WriteHeader(http.StatusForbidden)
            return
        }
        
        next.ServeHTTP(w, r)
    }
}
```

**Gateway: proxy.go**

```go
func ProxyHandler(service *Service) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // ... проксирование запроса ...
        
        // Копируем заголовки от Backend
        for key, values := range resp.Header {
            // ✅ ФИЛЬТРУЕМ CORS заголовки от Backend
            if strings.HasPrefix(key, "Access-Control-") {
                continue  // Пропускаем!
            }
            
            for _, value := range values {
                w.Header().Add(key, value)
            }
        }
        
        // ... копируем body ...
    }
}
```

**Backend: НЕ устанавливает CORS**

```go
// ❌ НЕ ДЕЛАЙ ТАК:
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")  // НЕТ!
        next(w, r)
    }
}

// ✅ ПРАВИЛЬНО:
// Backend вообще не думает о CORS
// Gateway управляет всем
```

---

## API Client Architecture

### How apiClient Works

**File: frontend/lib/api.ts**

```typescript
// Определяем API_URL в зависимости от окружения
const API_URL = process.env.NEXT_PUBLIC_API_URL || 
  (typeof window !== 'undefined' && process.env.NODE_ENV === 'production'
    ? '' // Production: пустая строка = относительные пути
    : 'http://localhost:8000'); // Development: прямой URL

export class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async getHeaders(): Promise<HeadersInit> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    // Добавляем Authorization если есть токен в localStorage
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('auth_token');
      if (token && token !== 'authenticated') {
        headers['Authorization'] = `Bearer ${token}`;
      }
    }

    return headers;
  }

  async get<T>(endpoint: string): Promise<ApiResponse<T>> {
    const url = `${this.baseUrl}${endpoint}`;
    const headers = await this.getHeaders();
    
    const response = await fetch(url, {
      method: 'GET',
      headers,
      credentials: 'include',  // ✅ Важно для cookies!
    });
    
    return this.handleResponse<T>(response);
  }

  async post<T>(endpoint: string, data: any): Promise<ApiResponse<T>> {
    const url = `${this.baseUrl}${endpoint}`;
    const headers = await this.getHeaders();
    
    const response = await fetch(url, {
      method: 'POST',
      headers,
      credentials: 'include',  // ✅ Важно для cookies!
      body: JSON.stringify(data),
    });
    
    return this.handleResponse<T>(response);
  }
  
  // ... put, delete методы ...
}

// Создаем экземпляр
export const apiClient = new ApiClient(API_URL);
```

### Environment Configuration

**Development: frontend/.env.local**

```env
# Backend API URL
NEXT_PUBLIC_API_URL=http://localhost:8000

# Auth Service URL (если используется отдельный Auth Service)
NEXT_PUBLIC_AUTH_URL=http://localhost:7100
```

**Production: frontend/.env.production**

```env
# Backend API URL
# Пустая строка = относительные пути (/api/...)
NEXT_PUBLIC_API_URL=

# Auth Service URL
NEXT_PUBLIC_AUTH_URL=
```

### Why Use apiClient Instead of Direct fetch

**❌ НЕПРАВИЛЬНО: Direct fetch**

```typescript
// Проблема 1: Хардкод URL
const response = await fetch('http://localhost:8000/api/posts');
// ❌ Не работает в production!

// Проблема 2: Забыли credentials
const response = await fetch('/api/posts');
// ❌ Cookie не отправляется!

// Проблема 3: Забыли Authorization
const response = await fetch('/api/posts', {
  credentials: 'include',
});
// ❌ Токен не отправляется!

// Проблема 4: Разная обработка ошибок
const response = await fetch('/api/posts');
if (!response.ok) {
  // Каждый раз пишем свою логику
}
```

**✅ ПРАВИЛЬНО: apiClient**

```typescript
import { apiClient } from '@/lib/api';

// ✅ Автоматически использует правильный URL
// ✅ Автоматически добавляет credentials: 'include'
// ✅ Автоматически добавляет Authorization если есть токен
// ✅ Единообразная обработка ошибок
const response = await apiClient.get<Post[]>('/api/posts');

if (response.success && response.data) {
  console.log(response.data);
} else {
  console.error(response.error);
}
```

---

## Middleware Architecture

### DevAuthMiddleware (Development)

**File: backend/middleware/dev_auth.go**

```go
func DevAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. Проверяем есть ли X-User-ID от Gateway
        if r.Header.Get("X-User-ID") != "" {
            // Работаем через Gateway - используем его заголовки
            AuthMiddleware(next)(w, r)
            return
        }
        
        // 2. Dev режим - проверяем JWT локально
        var tokenString string
        
        // Читаем из Authorization заголовка
        authHeader := r.Header.Get("Authorization")
        if strings.HasPrefix(authHeader, "Bearer ") {
            tokenString = strings.TrimPrefix(authHeader, "Bearer ")
        }
        
        // Или из cookie
        if tokenString == "" {
            cookie, err := r.Cookie("auth_token")
            if err == nil {
                tokenString = cookie.Value
            }
        }
        
        // 3. Валидируем JWT
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(os.Getenv("JWT_SECRET")), nil
        })
        
        if err != nil || !token.Valid {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        
        // 4. Извлекаем claims
        claims := token.Claims.(jwt.MapClaims)
        userID := int(claims["user_id"].(float64))
        userEmail := claims["email"].(string)
        userRole := claims["role"].(string)
        
        // 5. Добавляем в context
        ctx := context.WithValue(r.Context(), "userID", userID)
        ctx = context.WithValue(ctx, "userEmail", userEmail)
        ctx = context.WithValue(ctx, "userRole", userRole)
        
        next(w, r.WithContext(ctx))
    }
}
```

**Когда использовать:**
- ✅ Локальная разработка (без Gateway)
- ✅ Тестирование Backend отдельно от Gateway
- ❌ Production (используй AuthMiddleware)

### AuthMiddleware (Production)

**File: backend/middleware/auth.go**

```go
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // 1. Читаем заголовки от Gateway
        userIDStr := r.Header.Get("X-User-ID")
        userEmail := r.Header.Get("X-User-Email")
        userRole := r.Header.Get("X-User-Role")
        
        // 2. Если заголовков нет - не авторизован
        if userIDStr == "" {
            w.WriteHeader(http.StatusUnauthorized)
            return
        }
        
        // 3. Конвертируем userID
        var userID int
        fmt.Sscanf(userIDStr, "%d", &userID)
        
        // 4. Добавляем в context
        ctx := context.WithValue(r.Context(), "userID", userID)
        ctx = context.WithValue(ctx, "userEmail", userEmail)
        ctx = context.WithValue(ctx, "userRole", userRole)
        
        next(w, r.WithContext(ctx))
    }
}
```

**Когда использовать:**
- ✅ Production (через Gateway)
- ✅ Когда Gateway уже проверил JWT
- ❌ Локальная разработка без Gateway

---

## Security Considerations

### JWT Secret

**КРИТИЧНО:** `JWT_SECRET` должен быть одинаковым везде!

```bash
# Gateway
JWT_SECRET=your-super-secret-key

# Backend (development)
JWT_SECRET=your-super-secret-key  # ✅ Тот же!

# Backend (production)
JWT_SECRET=your-super-secret-key  # ✅ Тот же!
```

**Почему:**
- Gateway создает JWT с этим секретом
- Backend (dev) проверяет JWT с этим секретом
- Если секреты разные → токен невалидный → 401 Unauthorized

### Cookie Security

**Development:**
```go
http.SetCookie(w, &http.Cookie{
    Name:     "auth_token",
    Value:    token,
    Path:     "/",
    Domain:   "localhost",  // ✅ Работает для всех портов localhost
    HttpOnly: true,
    Secure:   false,  // HTTP OK для localhost
    SameSite: http.SameSiteLaxMode,
    MaxAge:   86400 * 7,  // 7 дней
})
```

**Production:**
```go
http.SetCookie(w, &http.Cookie{
    Name:     "auth_token",
    Value:    token,
    Path:     "/",
    Domain:   ".easypanel.host",  // ✅ Работает для всех поддоменов
    HttpOnly: true,
    Secure:   true,  // ✅ Только HTTPS
    SameSite: http.SameSiteLaxMode,
    MaxAge:   86400 * 7,
})
```

### CORS Security

**Никогда не используй `*` в production:**

```go
// ❌ ОПАСНО:
w.Header().Set("Access-Control-Allow-Origin", "*")
w.Header().Set("Access-Control-Allow-Credentials", "true")
// Браузер заблокирует! Нельзя использовать * с credentials

// ✅ ПРАВИЛЬНО:
allowedOrigins := map[string]bool{
    "https://frontend.com": true,
}
if allowedOrigins[origin] {
    w.Header().Set("Access-Control-Allow-Origin", origin)
    w.Header().Set("Access-Control-Allow-Credentials", "true")
}
```

---

## Component Integration

### Frontend Component Example

```typescript
'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { apiClient } from '@/lib/api';

interface Post {
  id: number;
  content: string;
  author_id: number;
}

export default function PostsPage() {
  const router = useRouter();
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadPosts();
  }, []);

  const loadPosts = async () => {
    try {
      // ✅ Используем apiClient
      const response = await apiClient.get<Post[]>('/api/posts');
      
      if (response.success && response.data) {
        setPosts(response.data);
      } else if (response.error === 'Unauthorized') {
        // Не авторизован - редирект на логин
        router.push('/auth');
      }
    } catch (error) {
      console.error('Failed to load posts:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div>Загрузка...</div>;

  return (
    <div>
      {posts.map(post => (
        <div key={post.id}>{post.content}</div>
      ))}
    </div>
  );
}
```

### Backend Handler Example

```go
package handlers

import (
    "encoding/json"
    "log"
    "net/http"
    "database"
)

func HandlePosts(w http.ResponseWriter, r *http.Request) {
    // ✅ Получаем userID из контекста (установлен middleware)
    userID, ok := r.Context().Value("userID").(int)
    if !ok || userID == 0 {
        http.Error(w, `{"success":false,"error":"Unauthorized"}`, http.StatusUnauthorized)
        return
    }
    
    // Можно также получить email и role
    userEmail, _ := r.Context().Value("userEmail").(string)
    userRole, _ := r.Context().Value("userRole").(string)
    
    log.Printf("✅ User: id=%d, email=%s, role=%s", userID, userEmail, userRole)
    
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
        http.Error(w, `{"success":false,"error":"Database error"}`, http.StatusInternalServerError)
        return
    }
    defer rows.Close()
    
    var posts []map[string]interface{}
    // ... scan rows ...
    
    response := map[string]interface{}{
        "success": true,
        "data": posts,
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

---

## Summary

### Key Principles

1. **Gateway is the single entry point** - все запросы идут через Gateway
2. **Gateway manages CORS** - Backend не устанавливает CORS заголовки
3. **Gateway validates JWT** - Backend доверяет X-User-* заголовкам
4. **Use apiClient** - не используй прямой fetch
5. **Relative paths in production** - используй относительные пути
6. **Same JWT_SECRET everywhere** - секрет должен совпадать

### Development vs Production

| Aspect | Development | Production |
|--------|-------------|------------|
| Frontend URL | http://localhost:3000 | https://frontend.com |
| Backend URL | http://localhost:8000 | http://backend:8000 |
| API requests | Direct to Backend | Through Gateway |
| JWT validation | Backend (DevAuthMiddleware) | Gateway |
| CORS | Backend | Gateway |
| Paths | Absolute (http://...) | Relative (/api/...) |
| Cookie Domain | localhost | .easypanel.host |

### Next Steps

1. Review this design document
2. Create implementation tasks
3. Update existing code to follow architecture
4. Add tests for authentication and CORS
5. Document troubleshooting procedures
