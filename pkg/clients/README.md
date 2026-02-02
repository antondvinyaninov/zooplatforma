# API Clients - Клиенты для межсервисного взаимодействия

Этот пакет содержит HTTP клиенты для взаимодействия между микросервисами.

## Структура

```
pkg/clients/
├── base_client.go      # Базовый HTTP клиент
├── auth_client.go      # Клиент для Auth Service
├── petbase_client.go   # Клиент для PetBase Service
├── main_client.go      # Клиент для Main Service
└── README.md           # Этот файл
```

---

## BaseClient

Базовый HTTP клиент с поддержкой:
- GET, POST, PUT, DELETE запросов
- Автоматической сериализации JSON
- Передачи JWT токенов
- Retry механизма
- Timeout

### Использование

```go
import "github.com/yourproject/pkg/clients"

// Создать клиент
client := clients.NewBaseClient("http://localhost:8000", 30*time.Second)

// GET запрос
resp, err := client.Get("/api/users", token)

// POST запрос
data := map[string]interface{}{"name": "John"}
resp, err := client.Post("/api/users", data, token)

// PUT запрос
resp, err := client.Put("/api/users/1", data, token)

// DELETE запрос
resp, err := client.Delete("/api/users/1", token)
```

---

## AuthClient

Клиент для Auth Service (порт 7000).

### Методы

#### VerifyToken
Проверяет JWT токен.

```go
authClient := clients.NewAuthClient("http://localhost:7000")

result, err := authClient.VerifyToken(token)
if err != nil {
    log.Fatal(err)
}

if result.Valid {
    fmt.Printf("User: %s (%s)\n", result.User.Name, result.User.Email)
}
```

#### GetMe
Получает информацию о текущем пользователе.

```go
user, err := authClient.GetMe(token)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("User ID: %d, Email: %s\n", user.ID, user.Email)
```

#### Login
Выполняет вход в систему.

```go
loginResp, err := authClient.Login("user@example.com", "password123")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Token: %s\n", loginResp.Token)
fmt.Printf("User: %s\n", loginResp.User.Name)
```

---

## PetBaseClient

Клиент для PetBase Service (порт 8100).

### Методы

#### GetPet
Получает информацию о питомце.

```go
petbaseClient := clients.NewPetBaseClient("http://localhost:8100")

pet, err := petbaseClient.GetPet(1, token)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Pet: %s (%s)\n", pet.Name, pet.Species)
```

#### GetUserPets
Получает питомцев пользователя.

```go
pets, err := petbaseClient.GetUserPets(userID, token)
if err != nil {
    log.Fatal(err)
}

for _, pet := range pets {
    fmt.Printf("- %s (%s)\n", pet.Name, pet.Breed)
}
```

#### CreatePet
Создает нового питомца.

```go
newPet := &clients.Pet{
    Name:      "Барсик",
    Species:   "cat",
    Breed:     "Британская короткошерстная",
    BirthDate: "2020-05-15",
    Gender:    "male",
    OwnerID:   userID,
}

pet, err := petbaseClient.CreatePet(newPet, token)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Created pet ID: %d\n", pet.ID)
```

#### UpdatePet
Обновляет питомца.

```go
updates := map[string]interface{}{
    "name":  "Барсик Обновленный",
    "color": "Темно-серый",
}

pet, err := petbaseClient.UpdatePet(petID, updates, token)
if err != nil {
    log.Fatal(err)
}
```

#### GetMedicalRecords
Получает медицинские записи питомца.

```go
records, err := petbaseClient.GetMedicalRecords(petID, token)
if err != nil {
    log.Fatal(err)
}

for _, record := range records {
    fmt.Printf("Date: %s, Diagnosis: %s\n", record.Date, record.Diagnosis)
}
```

#### GetVaccinations
Получает вакцинации питомца.

```go
vaccs, err := petbaseClient.GetVaccinations(petID, token)
if err != nil {
    log.Fatal(err)
}

for _, vacc := range vaccs {
    fmt.Printf("Vaccine: %s, Date: %s\n", vacc.VaccineName, vacc.Date)
}
```

---

## MainClient

Клиент для Main Service (порт 8000).

### Методы

#### GetOrganization
Получает информацию об организации.

```go
mainClient := clients.NewMainClient("http://localhost:8000")

org, err := mainClient.GetOrganization(orgID, token)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Organization: %s (%s)\n", org.Name, org.Type)
```

#### GetUserOrganizations
Получает организации пользователя.

```go
orgs, err := mainClient.GetUserOrganizations(userID, token)
if err != nil {
    log.Fatal(err)
}

for _, org := range orgs {
    fmt.Printf("- %s (%s)\n", org.Name, org.Type)
}
```

#### GetOrganizationMembers
Получает участников организации.

```go
members, err := mainClient.GetOrganizationMembers(orgID, token)
if err != nil {
    log.Fatal(err)
}

for _, member := range members {
    fmt.Printf("- %s (%s)\n", member.UserName, member.Role)
}
```

#### CheckMembership
Проверяет является ли пользователь участником организации.

```go
isMember, err := mainClient.CheckMembership(userID, orgID, token)
if err != nil {
    log.Fatal(err)
}

if isMember {
    fmt.Println("User is a member")
} else {
    fmt.Println("User is not a member")
}
```

---

## Пример использования в middleware

### Auth Middleware

```go
package middleware

import (
    "context"
    "net/http"
    "os"
    
    "github.com/yourproject/pkg/clients"
)

type contextKey string

const (
    UserIDKey    contextKey = "userID"
    UserEmailKey contextKey = "userEmail"
    UserRoleKey  contextKey = "userRole"
)

func AuthMiddleware(next http.Handler) http.Handler {
    authClient := clients.NewAuthClient(os.Getenv("AUTH_SERVICE_URL"))
    
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, `{"success":false,"error":"Unauthorized"}`, http.StatusUnauthorized)
            return
        }
        
        // Убрать "Bearer "
        if len(token) > 7 && token[:7] == "Bearer " {
            token = token[7:]
        }
        
        // Проверить токен через Auth Service
        result, err := authClient.VerifyToken(token)
        if err != nil || !result.Valid {
            http.Error(w, `{"success":false,"error":"Invalid token"}`, http.StatusUnauthorized)
            return
        }
        
        // Добавить данные пользователя в контекст
        ctx := context.WithValue(r.Context(), UserIDKey, result.User.ID)
        ctx = context.WithValue(ctx, UserEmailKey, result.User.Email)
        ctx = context.WithValue(ctx, UserRoleKey, result.User.Role)
        
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Tenant Middleware (для Clinic/Shelter)

```go
package middleware

import (
    "context"
    "net/http"
    "os"
    "strconv"
    
    "github.com/yourproject/pkg/clients"
)

const ClinicIDKey contextKey = "clinicID"

func TenantMiddleware(next http.Handler) http.Handler {
    mainClient := clients.NewMainClient(os.Getenv("MAIN_SERVICE_URL"))
    
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Получить user_id из контекста (установлен AuthMiddleware)
        userID, ok := r.Context().Value(UserIDKey).(int)
        if !ok {
            http.Error(w, `{"success":false,"error":"Unauthorized"}`, http.StatusUnauthorized)
            return
        }
        
        // Получить clinic_id из заголовка
        clinicIDStr := r.Header.Get("X-Clinic-ID")
        clinicID, err := strconv.Atoi(clinicIDStr)
        if err != nil {
            http.Error(w, `{"success":false,"error":"Invalid clinic ID"}`, http.StatusBadRequest)
            return
        }
        
        // Получить токен
        token := r.Header.Get("Authorization")
        if len(token) > 7 && token[:7] == "Bearer " {
            token = token[7:]
        }
        
        // Проверить что пользователь - участник клиники
        isMember, err := mainClient.CheckMembership(userID, clinicID, token)
        if err != nil || !isMember {
            http.Error(w, `{"success":false,"error":"Forbidden: not a member"}`, http.StatusForbidden)
            return
        }
        
        // Добавить clinic_id в контекст
        ctx := context.WithValue(r.Context(), ClinicIDKey, clinicID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

## Обработка ошибок

Все клиенты возвращают ошибки в стандартном формате:

```go
pet, err := petbaseClient.GetPet(petID, token)
if err != nil {
    // Логировать ошибку
    log.Printf("❌ Failed to get pet: %v", err)
    
    // Вернуть ошибку пользователю
    http.Error(w, `{"success":false,"error":"Failed to load pet"}`, http.StatusInternalServerError)
    return
}
```

---

## Конфигурация

Все URL сервисов должны быть в .env файле:

```bash
# .env
AUTH_SERVICE_URL=http://localhost:7000
PETBASE_SERVICE_URL=http://localhost:8100
MAIN_SERVICE_URL=http://localhost:8000
```

Использование в коде:

```go
import "os"

authClient := clients.NewAuthClient(os.Getenv("AUTH_SERVICE_URL"))
petbaseClient := clients.NewPetBaseClient(os.Getenv("PETBASE_SERVICE_URL"))
mainClient := clients.NewMainClient(os.Getenv("MAIN_SERVICE_URL"))
```

---

## Тестирование

```go
package main

import (
    "testing"
    
    "github.com/yourproject/pkg/clients"
)

func TestAuthClient(t *testing.T) {
    client := clients.NewAuthClient("http://localhost:7000")
    
    // Тест логина
    resp, err := client.Login("test@example.com", "password123")
    if err != nil {
        t.Fatalf("Login failed: %v", err)
    }
    
    if resp.Token == "" {
        t.Error("Expected token, got empty string")
    }
}
```

---

**Дата обновления:** 29 января 2025
