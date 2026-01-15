# Архитектурные правила проекта

## ⚠️ КРИТИЧЕСКОЕ: ЗооБаза (PetBase) - единственный источник данных о питомцах

**Single Source of Truth (SSOT)** - все данные о питомцах хранятся и управляются только в ЗооБазе.

### Правило добавления полей для питомцев

Если нужно добавить новое поле для животных, **ОБЯЗАТЕЛЬНО** следуйте этому порядку:

#### 1. Сначала в ЗооБазе (PetBase)

**Шаг 1: Миграция базы данных**
```sql
-- database/migrations/XXX_add_new_field.sql
ALTER TABLE pets ADD COLUMN new_field TEXT;
```

**Шаг 2: Обновить Go модель**
```go
// petbase/backend/models/pet.go
type Pet struct {
    // ... существующие поля
    NewField string `json:"new_field"`
}
```

**Шаг 3: Обновить API handlers**
```go
// petbase/backend/handlers/pets.go
// Добавить поле в SELECT запросы
// Добавить поле в Scan операции
// Добавить поле в UPDATE запросы
```

**Шаг 4: Обновить TypeScript типы в PetBase frontend**
```typescript
// petbase/frontend/lib/types.ts
interface Pet {
  // ... существующие поля
  new_field?: string;
}
```

#### 2. Только потом в других сервисах

После того как поле добавлено в ЗооБазу и работает через API, можно обновлять другие сервисы:

**Main (localhost:3000)**
```typescript
// main/frontend/lib/api.ts
interface Pet {
  new_field?: string; // Добавить поле
}
```

**Owner (localhost:6100)**
```typescript
// owner/frontend/app/(dashboard)/pets/[id]/edit/page.tsx
// Добавить поле в форму редактирования
```

**Volunteer (localhost:6200)**
```typescript
// volunteer/frontend/...
// Добавить поле где необходимо
```

**Clinic (localhost:6300)**
```typescript
// clinic/frontend/...
// Добавить поле где необходимо
```

### Почему это важно?

1. **Консистентность данных** - одна база данных, одна схема
2. **Упрощение поддержки** - изменения в одном месте
3. **Избежание конфликтов** - нет дублирования логики
4. **Централизованная валидация** - правила в одном месте
5. **Легче отлаживать** - один источник правды

### API endpoints ЗооБазы

Все сервисы работают с питомцами через PetBase API:

- `GET /api/pets` - список питомцев
- `GET /api/pets/:id` - получить данные питомца
- `POST /api/pets` - создать питомца
- `PUT /api/pets/:id` - обновить данные питомца
- `DELETE /api/pets/:id` - удалить питомца

**Порт:** `localhost:8100`

### Пример: Добавление поля "microchip_date"

#### ❌ Неправильно
```typescript
// Добавить поле сразу в Owner кабинете
// owner/frontend/app/(dashboard)/pets/[id]/edit/page.tsx
interface Pet {
  microchip_date?: string; // ❌ Поля нет в БД!
}
```

#### ✅ Правильно

**1. Миграция в ЗооБазе:**
```sql
-- database/migrations/033_add_microchip_date.sql
ALTER TABLE pets ADD COLUMN microchip_date DATE;
```

**2. Модель в ЗооБазе:**
```go
// petbase/backend/models/pet.go
type Pet struct {
    // ...
    MicrochipDate *string `json:"microchip_date"`
}
```

**3. API handler в ЗооБазе:**
```go
// petbase/backend/handlers/pets.go
// Добавить в SELECT и UPDATE
```

**4. Типы в PetBase frontend:**
```typescript
// petbase/frontend/lib/types.ts
interface Pet {
  microchip_date?: string;
}
```

**5. Теперь можно использовать в других сервисах:**
```typescript
// owner/frontend/app/(dashboard)/pets/[id]/edit/page.tsx
interface Pet {
  microchip_date?: string; // ✅ Поле есть в БД и API
}
```

## Другие архитектурные правила

### SSO авторизация

Все сервисы используют SSO через Main API:
- Проверка авторизации: `GET http://localhost:8000/api/auth/me`
- Логин: `POST http://localhost:8000/api/auth/login`
- Логаут: `POST http://localhost:8000/api/auth/logout`

### CORS настройки

Каждый backend должен разрешать запросы от всех frontend сервисов:
```go
allowedOrigins := map[string]bool{
    "http://localhost:3000": true, // Main
    "http://localhost:4000": true, // Admin
    "http://localhost:4100": true, // PetBase
    "http://localhost:5100": true, // Shelter
    "http://localhost:6100": true, // Owner
    "http://localhost:6200": true, // Volunteer
    "http://localhost:6300": true, // Clinic
}
```

### Структура микросервисов

Каждый микросервис имеет:
- `backend/` - Go API сервер
- `frontend/` - Next.js приложение
- `README.md` - документация сервиса

### Порты сервисов

| Сервис | Frontend | Backend |
|--------|----------|---------|
| Main | 3000 | 8000 |
| Admin | 4000 | 9000 |
| PetBase | 4100 | 8100 |
| Shelter | 5100 | 8200 |
| Owner | 6100 | 8400 |
| Volunteer | 6200 | 8500 |
| Clinic | 6300 | 8600 |

---

**Помните:** ЗооБаза - это фундамент для всех данных о питомцах. Всегда начинайте с неё!
