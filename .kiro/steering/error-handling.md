---
inclusion: always
---

# Error Handling - Правила обработки ошибок

## ⚠️ КРИТИЧЕСКОЕ ПРАВИЛО: Единообразная обработка ошибок

**Все ошибки обрабатываются одинаково во всех микросервисах!**

## Принципы обработки ошибок

1. **Никогда не игнорируй ошибки** - всегда проверяй `err != nil`
2. **Логируй все ошибки** - для отладки и мониторинга
3. **Возвращай понятные сообщения** - пользователь должен понимать что произошло
4. **Не раскрывай внутренние детали** - безопасность превыше всего
5. **Используй правильные HTTP коды** - клиент должен понимать тип ошибки

## Backend (Go) - Обработка ошибок

### Базовый паттерн

```go
func Handler(w http.ResponseWriter, r *http.Request) {
    // 1. Получить данные
    data, err := getData()
    if err != nil {
        // 2. Логировать ошибку (с деталями)
        log.Printf("❌ Handler error: %v", err)
        
        // 3. Вернуть понятное сообщение (без деталей)
        http.Error(w, `{"success":false,"error":"Failed to load data"}`, http.StatusInternalServerError)
        return
    }
    
    // 4. Успешный ответ
    response := map[string]interface{}{
        "success": true,
        "data": data,
    }
    json.NewEncoder(w).Encode(response)
}
```

### Типы ошибок и их обработка

#### 1. Ошибки базы данных

```go
func GetUser(w http.ResponseWriter, r *http.Request) {
    var user User
    err := db.QueryRow("SELECT * FROM users WHERE id = ?", userID).Scan(&user.ID, &user.Name)
    
    if err == sql.ErrNoRows {
        // Ресурс не найден - 404
        log.Printf("⚠️ User not found: id=%d", userID)
        http.Error(w, `{"success":false,"error":"User not found"}`, http.StatusNotFound)
        return
    }
    
    if err != nil {
        // Ошибка БД - 500
        log.Printf("❌ Database error: %v", err)
        http.Error(w, `{"success":false,"error":"Internal server error"}`, http.StatusInternalServerError)
        return
    }
    
    // Успех
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": user,
    })
}
```

#### 2. Ошибки валидации

```go
func CreateUser(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        log.Printf("❌ Invalid JSON: %v", err)
        http.Error(w, `{"success":false,"error":"Invalid JSON format"}`, http.StatusBadRequest)
        return
    }
    
    // Валидация полей
    errors := make(map[string]string)
    
    if user.Email == "" {
        errors["email"] = "Email is required"
    } else if !isValidEmail(user.Email) {
        errors["email"] = "Invalid email format"
    }
    
    if user.Password == "" {
        errors["password"] = "Password is required"
    } else if len(user.Password) < 8 {
        errors["password"] = "Password must be at least 8 characters"
    }
    
    if len(errors) > 0 {
        log.Printf("⚠️ Validation failed: %v", errors)
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "success": false,
            "error": "Validation failed",
            "details": errors,
        })
        return
    }
    
    // Создание пользователя
    // ...
}
```

#### 3. Ошибки авторизации

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        
        if token == "" {
            log.Printf("⚠️ No authorization token")
            http.Error(w, `{"success":false,"error":"Unauthorized"}`, http.StatusUnauthorized)
            return
        }
        
        claims, err := parseJWT(token)
        if err != nil {
            log.Printf("❌ Invalid token: %v", err)
            http.Error(w, `{"success":false,"error":"Invalid token"}`, http.StatusUnauthorized)
            return
        }
        
        // Проверка прав
        if !hasPermission(claims, r.URL.Path) {
            log.Printf("⚠️ Forbidden: user_id=%d, path=%s", claims.UserID, r.URL.Path)
            http.Error(w, `{"success":false,"error":"Forbidden"}`, http.StatusForbidden)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}
```

#### 4. Ошибки конфликтов

```go
func CreateUser(w http.ResponseWriter, r *http.Request) {
    // Проверка существования email
    var exists bool
    err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", user.Email).Scan(&exists)
    
    if err != nil {
        log.Printf("❌ Database error: %v", err)
        http.Error(w, `{"success":false,"error":"Internal server error"}`, http.StatusInternalServerError)
        return
    }
    
    if exists {
        log.Printf("⚠️ Email already exists: %s", user.Email)
        http.Error(w, `{"success":false,"error":"Email already exists"}`, http.StatusConflict)
        return
    }
    
    // Создание пользователя
    // ...
}
```

### Логирование ошибок

**Используй эмодзи для быстрой визуальной идентификации:**

```go
log.Printf("✅ Success: user created, id=%d", userID)
log.Printf("⚠️ Warning: user not found, id=%d", userID)
log.Printf("❌ Error: database connection failed: %v", err)
log.Printf("🔍 Debug: processing request, user_id=%d", userID)
log.Printf("🔐 Auth: user logged in, email=%s", email)
log.Printf("🏥 Tenant: clinic verified, clinic_id=%d", clinicID)
log.Printf("📋 Info: loading data, count=%d", count)
```

**Уровни логирования:**

```go
// DEBUG - детальная информация для отладки
log.Printf("🔍 GetUser: user_id=%d, include_posts=%v", userID, includePosts)

// INFO - обычные операции
log.Printf("📋 User registered: email=%s", email)

// WARNING - потенциальные проблемы
log.Printf("⚠️ User not found: id=%d", userID)

// ERROR - ошибки требующие внимания
log.Printf("❌ Database error: %v", err)

// CRITICAL - критические ошибки
log.Printf("🚨 CRITICAL: Database connection lost")
```

### Helper функции для ошибок

```go
// respondError - универсальная функция для ошибок
func respondError(w http.ResponseWriter, message string, statusCode int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": false,
        "error": message,
    })
}

// respondValidationError - ошибка валидации с деталями
func respondValidationError(w http.ResponseWriter, errors map[string]string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusBadRequest)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": false,
        "error": "Validation failed",
        "details": errors,
    })
}

// respondSuccess - успешный ответ
func respondSuccess(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data": data,
    })
}

// Использование
func GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := db.GetUser(userID)
    
    if err == sql.ErrNoRows {
        respondError(w, "User not found", http.StatusNotFound)
        return
    }
    
    if err != nil {
        log.Printf("❌ Database error: %v", err)
        respondError(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    
    respondSuccess(w, user)
}
```

## Frontend (React/TypeScript) - Обработка ошибок

### Базовый паттерн

```typescript
const [error, setError] = useState<string | null>(null);
const [loading, setLoading] = useState(false);

const fetchData = async () => {
  setLoading(true);
  setError(null);
  
  try {
    const response = await fetch('http://localhost:8000/api/data');
    
    if (!response.ok) {
      // HTTP ошибка
      const errorData = await response.json();
      throw new Error(errorData.error || 'Failed to load data');
    }
    
    const data = await response.json();
    
    if (!data.success) {
      // API вернул success: false
      throw new Error(data.error || 'Unknown error');
    }
    
    // Успех
    setData(data.data);
  } catch (error) {
    // Обработка ошибки
    console.error('Failed to fetch data:', error);
    setError(error instanceof Error ? error.message : 'Unknown error');
  } finally {
    setLoading(false);
  }
};
```

### Отображение ошибок пользователю

#### 1. Inline ошибки (в форме)

```typescript
export default function LoginForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [generalError, setGeneralError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrors({});
    setGeneralError(null);

    try {
      const response = await fetch('http://localhost:8000/api/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password }),
      });

      const data = await response.json();

      if (!response.ok) {
        if (data.details) {
          // Ошибки валидации
          setErrors(data.details);
        } else {
          // Общая ошибка
          setGeneralError(data.error || 'Login failed');
        }
        return;
      }

      // Успех
      router.push('/');
    } catch (error) {
      setGeneralError('Network error. Please try again.');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      {generalError && (
        <div className="p-4 mb-4 bg-red-50 border border-red-200 rounded-lg text-red-800">
          {generalError}
        </div>
      )}

      <div className="mb-4">
        <label>Email</label>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className={errors.email ? 'border-red-500' : ''}
        />
        {errors.email && (
          <p className="text-red-500 text-sm mt-1">{errors.email}</p>
        )}
      </div>

      <div className="mb-4">
        <label>Password</label>
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className={errors.password ? 'border-red-500' : ''}
        />
        {errors.password && (
          <p className="text-red-500 text-sm mt-1">{errors.password}</p>
        )}
      </div>

      <button type="submit">Login</button>
    </form>
  );
}
```

#### 2. Toast уведомления

```typescript
const [message, setMessage] = useState<{
  type: 'success' | 'error';
  text: string;
} | null>(null);

const handleAction = async () => {
  try {
    const response = await fetch('http://localhost:8000/api/action', {
      method: 'POST',
    });

    const data = await response.json();

    if (!response.ok) {
      setMessage({ type: 'error', text: data.error });
      return;
    }

    setMessage({ type: 'success', text: 'Action completed successfully' });
  } catch (error) {
    setMessage({ type: 'error', text: 'Network error' });
  }
};

return (
  <>
    {message && (
      <div
        className={`fixed top-4 right-4 p-4 rounded-lg shadow-lg ${
          message.type === 'success'
            ? 'bg-green-50 text-green-800 border border-green-200'
            : 'bg-red-50 text-red-800 border border-red-200'
        }`}
      >
        {message.text}
      </div>
    )}
    {/* Контент */}
  </>
);
```

#### 3. Error Boundary (для React ошибок)

```typescript
// components/ErrorBoundary.tsx
import React from 'react';

interface Props {
  children: React.ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

export class ErrorBoundary extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('React Error:', error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50">
          <div className="max-w-md p-8 bg-white rounded-lg shadow-lg">
            <h2 className="text-2xl font-bold text-red-600 mb-4">
              Что-то пошло не так
            </h2>
            <p className="text-gray-600 mb-4">
              Произошла ошибка при отображении страницы.
            </p>
            <button
              onClick={() => window.location.reload()}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
            >
              Перезагрузить страницу
            </button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

// Использование в layout.tsx
export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html>
      <body>
        <ErrorBoundary>
          {children}
        </ErrorBoundary>
      </body>
    </html>
  );
}
```

### API клиент с обработкой ошибок

```typescript
// lib/api.ts
class APIError extends Error {
  constructor(
    message: string,
    public statusCode: number,
    public details?: Record<string, string>
  ) {
    super(message);
    this.name = 'APIError';
  }
}

export async function apiRequest<T>(
  url: string,
  options?: RequestInit
): Promise<T> {
  try {
    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    const data = await response.json();

    if (!response.ok) {
      throw new APIError(
        data.error || 'Request failed',
        response.status,
        data.details
      );
    }

    if (!data.success) {
      throw new APIError(data.error || 'Unknown error', response.status);
    }

    return data.data;
  } catch (error) {
    if (error instanceof APIError) {
      throw error;
    }
    
    // Network error
    throw new APIError('Network error. Please check your connection.', 0);
  }
}

// Использование
const fetchUsers = async () => {
  try {
    const users = await apiRequest<User[]>('http://localhost:8000/api/users');
    setUsers(users);
  } catch (error) {
    if (error instanceof APIError) {
      if (error.statusCode === 401) {
        // Не авторизован
        router.push('/login');
      } else if (error.statusCode === 403) {
        // Нет прав
        setError('У вас нет прав для просмотра этой страницы');
      } else {
        setError(error.message);
      }
    } else {
      setError('Неизвестная ошибка');
    }
  }
};
```

## Сообщения об ошибках

### Правила написания сообщений

**✅ Хорошие сообщения:**
- "Email already exists" - понятно что делать
- "Password must be at least 8 characters" - конкретное требование
- "User not found" - понятно что произошло
- "You don't have permission to delete this post" - понятно почему отказано

**❌ Плохие сообщения:**
- "Error" - что за ошибка?
- "Something went wrong" - что именно?
- "sql: no rows in result set" - техническая деталь
- "Invalid input" - что именно невалидно?

### Локализация ошибок

```go
// errors.go
var ErrorMessages = map[string]string{
    "user_not_found": "Пользователь не найден",
    "email_exists": "Email уже используется",
    "invalid_password": "Неверный пароль",
    "unauthorized": "Необходима авторизация",
    "forbidden": "Недостаточно прав",
}

func GetErrorMessage(key string) string {
    if msg, ok := ErrorMessages[key]; ok {
        return msg
    }
    return "Неизвестная ошибка"
}

// Использование
http.Error(w, fmt.Sprintf(`{"success":false,"error":"%s"}`, GetErrorMessage("user_not_found")), http.StatusNotFound)
```

## Мониторинг ошибок

### Логирование в файл

```go
// main.go
func setupLogging() {
    logFile, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatal("Failed to open log file:", err)
    }
    
    // Логировать и в файл, и в консоль
    multiWriter := io.MultiWriter(os.Stdout, logFile)
    log.SetOutput(multiWriter)
    log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
```

### Структурированное логирование

```go
type LogEntry struct {
    Timestamp string `json:"timestamp"`
    Level     string `json:"level"`
    Message   string `json:"message"`
    UserID    int    `json:"user_id,omitempty"`
    Error     string `json:"error,omitempty"`
}

func logError(message string, userID int, err error) {
    entry := LogEntry{
        Timestamp: time.Now().Format(time.RFC3339),
        Level:     "ERROR",
        Message:   message,
        UserID:    userID,
    }
    
    if err != nil {
        entry.Error = err.Error()
    }
    
    jsonData, _ := json.Marshal(entry)
    log.Println(string(jsonData))
}
```

## Checklist обработки ошибок

### Backend:
- [ ] Все ошибки проверяются (`if err != nil`)
- [ ] Ошибки логируются с контекстом
- [ ] Используются правильные HTTP коды
- [ ] Сообщения понятны пользователю
- [ ] Не раскрываются внутренние детали
- [ ] Есть валидация входных данных
- [ ] Есть обработка ошибок БД

### Frontend:
- [ ] Все fetch обернуты в try-catch
- [ ] Ошибки отображаются пользователю
- [ ] Есть loading состояние
- [ ] Есть Error Boundary
- [ ] Сетевые ошибки обрабатываются
- [ ] 401/403 ошибки обрабатываются (редирект)
- [ ] Ошибки валидации показываются у полей

---

**Помни:** Хорошая обработка ошибок - это половина успеха приложения!
