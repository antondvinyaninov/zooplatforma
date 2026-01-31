---
inclusion: always
---

# Структура проекта - Быстрая справка

## ⚠️ КРИТИЧЕСКОЕ ПРАВИЛО: Всегда проверяй путь перед изменением файла!

**Проект использует микросервисную архитектуру. Каждый сервис имеет свой backend и frontend!**

## Микросервисы и их порты

| Сервис | Backend | Frontend | Назначение |
|--------|---------|----------|------------|
| **Main** | 8000 | 3000 | Главная платформа (лента, профили, питомцы) |
| **Admin** | 9000 | 4000 | Административная панель (модерация) |
| **PetBase** | 8100 | 4100 | ЗооБаза (справочник животных) |
| **Shelter** | 8200 | 5100 | Кабинет приюта |
| **Owner** | 8400 | 6100 | Кабинет владельца питомца |
| **Volunteer** | 8500 | 6200 | Кабинет зоопомощника |
| **Clinic** | 8600 | 6300 | Кабинет ветеринарной клиники |

## Структура папок

```
.
├── main/                    # Главный сервис
│   ├── backend/             # Go API (порт 8000)
│   │   ├── main.go
│   │   ├── handlers/
│   │   ├── middleware/
│   │   └── .air.toml
│   └── frontend/            # Next.js (порт 3000)
│       ├── app/
│       ├── lib/
│       └── package.json
│
├── admin/                   # Админ панель
│   ├── backend/             # Go API (порт 9000)
│   └── frontend/            # Next.js (порт 4000)
│
├── petbase/                 # ЗооБаза
│   ├── backend/             # Go API (порт 8100)
│   └── frontend/            # Next.js (порт 4100)
│
├── shelter/                 # Кабинет приюта
│   ├── backend/             # Go API (порт 8200)
│   └── frontend/            # Next.js (порт 5100)
│
├── owner/                   # Кабинет владельца
│   ├── backend/             # Go API (порт 8400)
│   └── frontend/            # Next.js (порт 6100)
│
├── volunteer/               # Кабинет зоопомощника
│   ├── backend/             # Go API (порт 8500)
│   └── frontend/            # Next.js (порт 6200)
│
├── clinic/                  # Кабинет клиники
│   ├── backend/             # Go API (порт 8600)
│   │   ├── main.go
│   │   ├── handlers/
│   │   │   └── clinic.go    # Handlers для клиники
│   │   ├── middleware/
│   │   │   ├── auth.go      # SSO авторизация
│   │   │   └── tenant.go    # Проверка clinic_id
│   │   ├── models/
│   │   └── .air.toml
│   └── frontend/            # Next.js (порт 6300)
│       ├── app/
│       │   ├── (dashboard)/ # Страницы кабинета
│       │   │   ├── layout.tsx
│       │   │   ├── overview/
│       │   │   ├── patients/
│       │   │   ├── appointments/
│       │   │   └── settings/
│       │   ├── components/
│       │   │   ├── admin/   # AdminLayout, AdminTab
│       │   │   └── MembersManager.tsx
│       │   └── select/      # Выбор клиники
│       └── lib/
│           └── api.ts       # API клиент
│
├── shared/                  # Общие компоненты (@pet/shared)
│   ├── src/
│   │   ├── api/
│   │   ├── components/      # ⚠️ НЕ импортировать! Копировать!
│   │   ├── types/
│   │   └── hooks/
│   └── package.json
│
├── database/                # Общая SQLite база
│   ├── db.go                # Инициализация БД
│   ├── data.db              # Файл базы данных
│   └── migrations/          # SQL миграции
│
├── docs/                    # Документация
│   ├── CHANGELOG.md         # История изменений
│   ├── STRUCTURE.md         # Полная структура
│   ├── ROADMAP.md           # План развития
│   └── *.md                 # Другие документы
│
├── .kiro/                   # Конфигурация Kiro
│   └── steering/            # Правила разработки
│       ├── project-structure.md  # Этот файл
│       ├── changelog-updates.md
│       ├── database-migrations.md
│       └── shared-components.md
│
└── uploads/                 # Загруженные файлы
    └── users/               # Аватарки пользователей
```

## Правила работы с файлами

### ✅ Перед изменением файла ВСЕГДА:

1. **Определи сервис** - Main, Admin, Clinic, Owner и т.д.
2. **Определи слой** - backend или frontend
3. **Проверь полный путь** - например: `clinic/backend/handlers/clinic.go`
4. **Убедись что это правильный файл** - не путай `main/backend/main.go` с `clinic/backend/main.go`!

### ❌ Частые ошибки:

- **Путаница между сервисами:**
  - `main/backend/handlers/auth.go` ≠ `clinic/backend/handlers/auth.go`
  - `main/frontend/app/page.tsx` ≠ `clinic/frontend/app/page.tsx`

- **Путаница между backend и frontend:**
  - `clinic/backend/main.go` - это Go сервер
  - `clinic/frontend/app/page.tsx` - это Next.js страница

- **Неправильные импорты:**
  - ❌ `import { AuthForm } from '@pet/shared'` - стили не применятся (Tailwind v4)
  - ✅ Скопировать компонент в `clinic/frontend/app/components/`

## Когда работаешь с Clinic:

### Backend (Go)
- **Путь:** `clinic/backend/`
- **Порт:** 8600
- **Главный файл:** `clinic/backend/main.go`
- **Handlers:** `clinic/backend/handlers/clinic.go`
- **Middleware:** `clinic/backend/middleware/auth.go`, `tenant.go`
- **Hot reload:** Air автоматически перезагружает при изменениях

### Frontend (Next.js)
- **Путь:** `clinic/frontend/`
- **Порт:** 6300
- **Страницы:** `clinic/frontend/app/(dashboard)/`
- **Компоненты:** `clinic/frontend/app/components/`
- **API клиент:** `clinic/frontend/lib/api.ts`
- **Hot reload:** Next.js автоматически обновляет

### Интеграция с Main
- **Аватарки:** хранятся на Main backend (порт 8000)
- **SSO:** авторизация через Main API `/api/auth/me`
- **База данных:** общая SQLite в `database/data.db`
- **Таблицы:**
  - `organizations` - клиники (type='clinic')
  - `organization_members` - участники клиники
  - `users` - пользователи платформы

## Когда работаешь с Main:

### Backend (Go)
- **Путь:** `main/backend/`
- **Порт:** 8000
- **Главный файл:** `main/backend/main.go`
- **Handlers:** `main/backend/handlers/`

### Frontend (Next.js)
- **Путь:** `main/frontend/`
- **Порт:** 3000
- **Страницы:** `main/frontend/app/(main)/`

## База данных

- **Путь:** `database/data.db`
- **Общая для всех сервисов**
- **Миграции:** `database/migrations/*.sql`
- **Правила:** см. `.kiro/steering/database-migrations.md`

### ⚠️ КРИТИЧЕСКОЕ ПРАВИЛО: PetBase - главная база данных животных

**PetBase (ЗооБаза) - это единственный источник правды о животных!**

#### Архитектура данных о животных:

1. **PetBase Backend (порт 8100)** - главный сервис для работы с данными животных
   - Таблицы: `species` (виды), `breeds` (породы), `pet_cards` (карточки)
   - Все CRUD операции с данными животных
   - API endpoints для других сервисов

2. **Другие сервисы** - ТОЛЬКО читают и используют данные через API
   - Main, Clinic, Shelter, Owner и т.д.
   - НЕ изменяют таблицы животных напрямую
   - Получают данные через `http://localhost:8100/api/...`

#### Правила работы с данными животных:

✅ **ПРАВИЛЬНО:**
```
1. Нужно добавить новое поле для породы?
   → Добавь в PetBase: таблица `breeds`, API endpoint
   → Обнови PetBase Backend handlers
   → Используй новый endpoint в других сервисах

2. Нужно получить список пород в Clinic?
   → Запрос к PetBase API: GET http://localhost:8100/api/breeds

3. Нужно создать карточку животного?
   → Запрос к PetBase API: POST http://localhost:8100/api/cards
```

❌ **НЕПРАВИЛЬНО:**
```
1. Добавлять поля в таблицы животных из Main/Clinic/Shelter
   → Только через PetBase!

2. Изменять данные пород напрямую в БД из других сервисов
   → Только через PetBase API!

3. Дублировать таблицы животных в других сервисах
   → Используй PetBase API!
```

#### Последовательность действий:

**Если нужно добавить новое поле для животных:**

1. **Сначала PetBase:**
   - Добавь поле в таблицу (`database/migrations/`)
   - Обнови модель в `petbase/backend/models/`
   - Обнови handler в `petbase/backend/handlers/`
   - Обнови API endpoint
   - Протестируй через PetBase Frontend (порт 4100)

2. **Потом другие сервисы:**
   - Используй новый endpoint PetBase API
   - Обнови типы в `shared/src/types/`
   - Обнови UI компоненты для отображения нового поля

**Пример: Добавление поля "размер" для породы**

```
Шаг 1: PetBase Backend
- Миграция: ALTER TABLE breeds ADD COLUMN size TEXT
- Модель: type Breed struct { Size string }
- Handler: обновить GetBreeds, CreateBreed
- API: GET/POST /api/breeds теперь возвращает/принимает size

Шаг 2: Clinic Frontend
- Запрос: fetch('http://localhost:8100/api/breeds')
- Отображение: показать поле size в UI
```

#### Интеграция между сервисами:

```
Main Frontend (3000)
    ↓ запрос данных о породах
PetBase API (8100)
    ↓ возвращает данные
Main Backend (8000)
    ↓ передает на фронт
Main Frontend (3000)
    ↓ отображает пользователю
```

**Важно:** PetBase - это микросервис, отвечающий за справочную информацию о животных. Все остальные сервисы - потребители этой информации.

## Shared компоненты

- **Путь:** `shared/src/`
- **Правило:** React компоненты с Tailwind классами КОПИРОВАТЬ, не импортировать!
- **Подробности:** см. `.kiro/steering/shared-components.md`

## Checklist перед изменением файла

- [ ] Определил сервис (Main, Clinic, Admin и т.д.)
- [ ] Определил слой (backend или frontend)
- [ ] Проверил полный путь файла
- [ ] Убедился что это правильный файл
- [ ] Понял назначение файла в контексте сервиса

## Быстрая навигация

**Если нужно изменить:**
- Авторизацию в Clinic → `clinic/backend/middleware/auth.go`
- Список участников клиники → `clinic/backend/handlers/clinic.go` (функция GetMembers)
- Отображение участников → `clinic/frontend/app/components/MembersManager.tsx`
- Страницу настроек клиники → `clinic/frontend/app/(dashboard)/settings/page.tsx`
- Главную авторизацию → `main/backend/handlers/auth.go`
- Профиль пользователя → `main/frontend/app/(main)/[userId]/page.tsx`

---

**Помни:** Каждый микросервис - это отдельное приложение со своим backend и frontend!
