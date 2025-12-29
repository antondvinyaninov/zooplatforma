# Безопасность ЗооБазы

## ✅ ТЕКУЩЕЕ СОСТОЯНИЕ

### Реализованная защита

1. **✅ JWT аутентификация** - реализована с поддержкой Bearer токенов
2. **✅ Проверка владельца** - пользователь может редактировать только своих питомцев
3. **✅ User ID из токена** - используется user_id из JWT, а не из запроса
4. **✅ CORS с проверкой origins** - разрешены только указанные домены
5. **✅ Middleware для защищённых routes** - применён к /api/pets endpoints

### Временные методы для разработки

- **X-User-ID заголовок** - для тестирования без полной системы аутентификации
- **Cookie user_id** - альтернативный метод для разработки

### Что ещё нужно сделать

- [ ] Rate limiting
- [ ] Валидация входных данных
- [ ] Логирование операций
- [ ] HTTPS в продакшене
- [ ] Refresh tokens

## 🛡️ АРХИТЕКТУРА БЕЗОПАСНОСТИ

### Этап 1: Аутентификация через JWT (КРИТИЧНО)

#### Архитектура

```
Основной сайт                          ЗооБаза
    │                                      │
    │  1. Пользователь логинится          │
    │  ← JWT токен                         │
    │                                      │
    │  2. Запрос с токеном                 │
    │  → Authorization: Bearer <token>  →  │
    │                                      │
    │                                   3. Проверка токена
    │                                   4. Извлечение user_id
    │                                   5. Проверка прав
    │                                      │
    │  ← Данные (только свои)           ←  │
```

#### Реализация

**1. Создать middleware для проверки токена**

```go
// petbase/backend/middleware/auth.go
package middleware

import (
    "context"
    "net/http"
    "strings"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("YOUR_SECRET_KEY") // TODO: из .env

type Claims struct {
    UserID int `json:"user_id"`
    jwt.RegisteredClaims
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Получаем токен из заголовка
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Unauthorized: no token", http.StatusUnauthorized)
            return
        }

        // Проверяем формат "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Unauthorized: invalid token format", http.StatusUnauthorized)
            return
        }

        tokenString := parts[1]

        // Парсим и проверяем токен
        token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
            return
        }

        // Извлекаем claims
        claims, ok := token.Claims.(*Claims)
        if !ok {
            http.Error(w, "Unauthorized: invalid claims", http.StatusUnauthorized)
            return
        }

        // Добавляем user_id в контекст
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}

// Опциональная аутентификация (для публичных endpoints)
func OptionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            // Нет токена - продолжаем без user_id
            next.ServeHTTP(w, r)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            next.ServeHTTP(w, r)
            return
        }

        token, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err == nil && token.Valid {
            if claims, ok := token.Claims.(*Claims); ok {
                ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
                r = r.WithContext(ctx)
            }
        }

        next.ServeHTTP(w, r)
    }
}

// Получить user_id из контекста
func GetUserID(r *http.Request) (int, bool) {
    userID, ok := r.Context().Value("user_id").(int)
    return userID, ok
}
```

**2. Обновить handlers для проверки владельца**

```go
// petbase/backend/handlers/pets.go

func createPet(w http.ResponseWriter, r *http.Request) {
    // Получаем user_id из токена (не из запроса!)
    userID, ok := middleware.GetUserID(r)
    if !ok {
        sendError(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var req CreatePetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        sendError(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // ВАЖНО: Используем user_id из токена, а не из запроса!
    req.UserID = userID

    // ... остальной код ...
}

func updatePet(w http.ResponseWriter, r *http.Request, id int) {
    // Получаем user_id из токена
    userID, ok := middleware.GetUserID(r)
    if !ok {
        sendError(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Проверяем, что питомец принадлежит пользователю
    var ownerID int
    err := database.DB.QueryRow("SELECT user_id FROM pets WHERE id = ?", id).Scan(&ownerID)
    if err != nil {
        sendError(w, "Pet not found", http.StatusNotFound)
        return
    }

    if ownerID != userID {
        sendError(w, "Forbidden: not your pet", http.StatusForbidden)
        return
    }

    // ... остальной код ...
}

func deletePet(w http.ResponseWriter, _ *http.Request, id int) {
    // Получаем user_id из токена
    userID, ok := middleware.GetUserID(r)
    if !ok {
        sendError(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Проверяем владельца
    var ownerID int
    err := database.DB.QueryRow("SELECT user_id FROM pets WHERE id = ?", id).Scan(&ownerID)
    if err != nil {
        sendError(w, "Pet not found", http.StatusNotFound)
        return
    }

    if ownerID != userID {
        sendError(w, "Forbidden: not your pet", http.StatusForbidden)
        return
    }

    // ... остальной код ...
}
```

**3. Применить middleware к routes**

```go
// petbase/backend/main.go

func main() {
    // ... существующий код ...

    // Публичные endpoints (без аутентификации)
    http.HandleFunc("/", enableCORS(handleRoot))
    http.HandleFunc("/api/health", enableCORS(handleHealth))
    http.HandleFunc("/api/species", enableCORS(handlers.SpeciesHandler))
    http.HandleFunc("/api/breeds", enableCORS(handlers.BreedsHandler))

    // Защищённые endpoints (требуют аутентификации)
    http.HandleFunc("/api/pets", enableCORS(middleware.AuthMiddleware(handlers.PetsHandler)))
    http.HandleFunc("/api/pets/", enableCORS(middleware.AuthMiddleware(handlers.PetDetailHandler)))

    // ... остальной код ...
}
```

**4. Обновить API клиент на основном сайте**

```typescript
// main/frontend/lib/petbase-api.ts

class PetBaseAPI {
  private baseURL: string;
  private getAuthToken: () => string | null;

  constructor(
    baseURL: string = PETBASE_API_URL,
    getAuthToken: () => string | null = () => localStorage.getItem('token')
  ) {
    this.baseURL = baseURL;
    this.getAuthToken = getAuthToken;
  }

  private getHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    const token = this.getAuthToken();
    if (token) {
      headers['Authorization'] = `Bearer ${token}`;
    }

    return headers;
  }

  async getPets(): Promise<Pet[]> {
    try {
      const response = await fetch(`${this.baseURL}/pets`, {
        headers: this.getHeaders(),
        credentials: 'include',
      });
      
      if (response.status === 401) {
        throw new Error('Unauthorized');
      }
      
      const data = await response.json();
      return data.success ? data.data : [];
    } catch (error) {
      console.error('Error fetching pets:', error);
      return [];
    }
  }

  // ... остальные методы с headers ...
}
```

### Этап 2: CORS - строгая проверка origins

```go
// petbase/backend/main.go

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        
        // Разрешённые origins (из .env в продакшене)
        allowedOrigins := map[string]bool{
            "http://localhost:3000": true,  // Development
            "http://localhost:4100": true,  // Development
            "https://yoursite.com": true,   // Production
            "https://petbase.yoursite.com": true, // Production
        }
        
        // Проверяем origin
        if !allowedOrigins[origin] {
            http.Error(w, "Forbidden origin", http.StatusForbidden)
            return
        }

        w.Header().Set("Access-Control-Allow-Origin", origin)
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        w.Header().Set("Access-Control-Allow-Credentials", "true")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next(w, r)
    }
}
```

### Этап 3: Rate Limiting

```go
// petbase/backend/middleware/ratelimit.go
package middleware

import (
    "net/http"
    "sync"
    "time"
)

type rateLimiter struct {
    requests map[string][]time.Time
    mu       sync.Mutex
    limit    int
    window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *rateLimiter {
    return &rateLimiter{
        requests: make(map[string][]time.Time),
        limit:    limit,
        window:   window,
    }
}

func (rl *rateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        
        rl.mu.Lock()
        defer rl.mu.Unlock()

        now := time.Now()
        windowStart := now.Add(-rl.window)

        // Очищаем старые запросы
        requests := rl.requests[ip]
        validRequests := []time.Time{}
        for _, t := range requests {
            if t.After(windowStart) {
                validRequests = append(validRequests, t)
            }
        }

        // Проверяем лимит
        if len(validRequests) >= rl.limit {
            http.Error(w, "Too many requests", http.StatusTooManyRequests)
            return
        }

        // Добавляем текущий запрос
        validRequests = append(validRequests, now)
        rl.requests[ip] = validRequests

        next(w, r)
    }
}

// Использование:
// limiter := middleware.NewRateLimiter(100, time.Minute) // 100 запросов в минуту
// http.HandleFunc("/api/pets", enableCORS(limiter.Middleware(middleware.AuthMiddleware(handlers.PetsHandler))))
```

### Этап 4: Валидация данных

```go
// petbase/backend/handlers/validation.go
package handlers

import (
    "errors"
    "regexp"
)

func validatePet(pet *CreatePetRequest) error {
    // Проверка имени
    if pet.Name == "" {
        return errors.New("name is required")
    }
    if len(pet.Name) > 100 {
        return errors.New("name is too long")
    }

    // Проверка вида
    if pet.Species == "" {
        return errors.New("species is required")
    }

    // Проверка веса
    if pet.Weight < 0 || pet.Weight > 1000 {
        return errors.New("invalid weight")
    }

    // Проверка email
    if pet.OwnerEmail != "" {
        emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
        if !emailRegex.MatchString(pet.OwnerEmail) {
            return errors.New("invalid email")
        }
    }

    // Проверка телефона
    if pet.OwnerPhone != "" {
        phoneRegex := regexp.MustCompile(`^\+?[0-9\s\-\(\)]{7,20}$`)
        if !phoneRegex.MatchString(pet.OwnerPhone) {
            return errors.New("invalid phone")
        }
    }

    return nil
}
```

### Этап 5: Логирование и мониторинг

```go
// petbase/backend/middleware/logging.go
package middleware

import (
    "log"
    "net/http"
    "time"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Логируем запрос
        log.Printf("[%s] %s %s from %s", 
            r.Method, 
            r.URL.Path, 
            r.RemoteAddr,
            r.Header.Get("User-Agent"),
        )

        // Выполняем запрос
        next(w, r)

        // Логируем время выполнения
        duration := time.Since(start)
        log.Printf("[%s] %s completed in %v", r.Method, r.URL.Path, duration)
    }
}
```

## 📋 CHECKLIST БЕЗОПАСНОСТИ

### ✅ Реализовано
- [x] Добавить JWT аутентификацию
- [x] Проверять владельца при изменении/удалении
- [x] Использовать user_id из токена, а не из запроса
- [x] Строгая проверка CORS origins
- [x] Применить auth middleware к защищённым routes
- [x] Обновить API клиент с поддержкой токенов

### Критично (сделать в ближайшее время)
- [ ] Валидация всех входных данных
- [ ] Rate limiting
- [ ] Логирование всех операций
- [ ] HTTPS в продакшене
- [ ] Хранить JWT secret в .env (пример создан)
- [ ] Добавить refresh tokens

### Желательно
- [ ] Мониторинг подозрительной активности
- [ ] Двухфакторная аутентификация
- [ ] Аудит логи
- [ ] Автоматическое тестирование безопасности

## 🔐 ПЕРЕМЕННЫЕ ОКРУЖЕНИЯ

```env
# petbase/backend/.env
JWT_SECRET=your-super-secret-key-change-in-production
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:4100
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60
```

## 🚀 РАЗВЁРТЫВАНИЕ В ПРОДАКШЕНЕ

### Обязательно
1. Использовать HTTPS (Let's Encrypt)
2. Сгенерировать сильный JWT secret (256 бит)
3. Настроить firewall
4. Регулярные бэкапы базы данных
5. Мониторинг и алерты

### Рекомендуется
1. Использовать reverse proxy (nginx)
2. Настроить fail2ban
3. Регулярные обновления зависимостей
4. Penetration testing

---

**Последнее обновление:** 29 декабря 2024

**Статус:** ✅ БАЗОВАЯ ЗАЩИТА РЕАЛИЗОВАНА

**Следующие шаги:**
1. Валидация входных данных
2. Rate limiting
3. Логирование операций
4. Тестирование безопасности
