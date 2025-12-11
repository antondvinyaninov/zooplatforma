# Структура проекта

Полное описание архитектуры и организации файлов проекта.

## Общая структура

```
.
├── backend/              # Go API сервер
├── frontend/             # Next.js веб-приложение
├── mobile/               # React Native мобильное приложение
├── database/             # SQLite база данных (отдельный модуль)
├── shared/               # Общая логика для всех платформ
├── run                   # Скрипт запуска всего проекта
├── .gitignore           # Игнорируемые файлы
├── README.md            # Основная документация
├── CHANGELOG.md         # История изменений
└── STRUCTURE.md         # Этот файл
```

---

## Backend (Go API)

**Путь:** `backend/`

**Технологии:**
- Go 1.25.5
- net/http (standard library)
- Air (hot reload)

**Структура:**
```
backend/
├── main.go              # Точка входа, HTTP сервер
├── go.mod               # Go модуль и зависимости
├── go.sum               # Checksums зависимостей
├── .air.toml            # Конфигурация Air для hot reload
└── tmp/                 # Временные файлы Air (игнорируется)
```

**Порт:** 8080

**Endpoints:**
- `GET /` - Welcome message
- `GET /api/health` - Health check

**Особенности:**
- CORS настроен на `*` (все источники)
- Автоматический hot reload при изменении .go файлов
- Интеграция с модулем database

---

## Database (SQLite)

**Путь:** `database/`

**Технологии:**
- SQLite 3
- go-sqlite3 драйвер

**Структура:**
```
database/
├── db.go                # Инициализация БД, создание таблиц
├── go.mod               # Отдельный Go модуль
├── go.sum               # Checksums зависимостей
└── data.db              # Файл базы данных (игнорируется)
```

**Таблицы:**
- `users` - Пример таблицы пользователей
  - id (INTEGER PRIMARY KEY AUTOINCREMENT)
  - name (TEXT NOT NULL)
  - email (TEXT UNIQUE NOT NULL)
  - created_at (DATETIME DEFAULT CURRENT_TIMESTAMP)

**Особенности:**
- Отдельный модуль для переиспользования
- Автоматическое создание таблиц при первом запуске
- База создается в `database/data.db`

---

## Frontend (Next.js Web)

**Путь:** `frontend/`

**Технологии:**
- Next.js 16.0.8 (App Router)
- React 19.2.1
- TypeScript 5.x
- Tailwind CSS 4.x
- SWC (компилятор)

**Структура:**
```
frontend/
├── src/
│   ├── app/             # Next.js App Router
│   │   ├── layout.tsx   # Корневой layout
│   │   └── page.tsx     # Главная страница
│   ├── components/
│   │   ├── ui/          # Базовые UI компоненты
│   │   │   └── Button.tsx
│   │   ├── shared/      # Переиспользуемые бизнес-компоненты
│   │   └── layout/      # Layout компоненты
│   │       └── Header.tsx
│   └── lib/
│       └── api.ts       # API клиент для связи с backend
├── public/              # Статические файлы
├── .env.local           # Environment variables (игнорируется)
├── package.json         # npm зависимости
├── tsconfig.json        # TypeScript конфигурация
├── tailwind.config.ts   # Tailwind CSS конфигурация
├── next.config.ts       # Next.js конфигурация
└── .next/               # Build output (игнорируется)
```

**Порт:** 3000

**Environment Variables:**
- `NEXT_PUBLIC_API_URL` - URL бекенда (по умолчанию: http://localhost:8080)

**Особенности:**
- App Router (новая архитектура Next.js)
- Turbopack для быстрой сборки
- Tailwind CSS для стилизации
- TypeScript для типизации
- Hot reload из коробки

---

## Mobile (React Native)

**Путь:** `mobile/`

**Технологии:**
- React Native 0.81.5
- Expo 54.0.27
- TypeScript 5.x
- Metro Bundler

**Структура:**
```
mobile/
├── App.tsx              # Главный компонент приложения
├── index.ts             # Точка входа
├── src/
│   └── lib/
│       └── api.ts       # API клиент для связи с backend
├── .env                 # Environment variables (игнорируется)
├── package.json         # npm зависимости
├── tsconfig.json        # TypeScript конфигурация
├── metro.config.js      # Metro bundler конфигурация
└── app.json             # Expo конфигурация
```

**Порт:** 8081 (Metro Bundler)

**Environment Variables:**
- `EXPO_PUBLIC_API_URL` - URL бекенда (по умолчанию: http://192.168.0.12:8080)

**Платформы:**
- Web (браузер)
- iOS (симулятор/устройство через Expo Go)
- Android (эмулятор/устройство через Expo Go)

**Особенности:**
- Expo для упрощенной разработки
- Metro bundler для hot reload
- react-native-web для web поддержки
- QR код для быстрого запуска на устройстве

---

## Shared (Общая логика)

**Путь:** `shared/`

**Технологии:**
- TypeScript 5.x

**Структура:**
```
shared/
├── src/
│   ├── api/
│   │   └── client.ts    # Универсальный API клиент
│   ├── types/
│   │   └── index.ts     # TypeScript типы и интерфейсы
│   ├── hooks/           # Переиспользуемые React хуки
│   ├── utils/           # Утилиты и хелперы
│   └── index.ts         # Экспорты модуля
├── package.json         # Метаданные пакета
└── tsconfig.json        # TypeScript конфигурация
```

**Назначение:**
- API клиент (используется в frontend)
- TypeScript типы (общие для всех платформ)
- React хуки (useAuth, useFetch и т.д.)
- Утилиты (форматирование, валидация)

**Примечание:**
- Frontend импортирует напрямую из `../shared/src`
- Mobile имеет копию API клиента (из-за ограничений Metro)

---

## Скрипты и конфигурация

### run (Скрипт запуска)

**Путь:** `run`

**Что делает:**
1. Очищает порты 3000, 8080, 8081
2. Запускает backend с Air (hot reload)
3. Запускает frontend (Next.js dev server)
4. Запускает mobile (Expo web)

**Использование:**
```bash
./run
```

**Остановка:** `Ctrl+C` (останавливает все три процесса)

### .gitignore

**Игнорируемые файлы:**
- `node_modules/` - npm зависимости
- `.next/` - Next.js build
- `backend/tmp/` - Air временные файлы
- `database/data.db` - SQLite база
- `.env`, `.env.local` - Environment variables
- `.DS_Store` - macOS системные файлы

---

## Порты

| Сервис   | Порт | URL                      |
|----------|------|--------------------------|
| Backend  | 8080 | http://localhost:8080    |
| Frontend | 3000 | http://localhost:3000    |
| Mobile   | 8081 | http://localhost:8081    |

---

## Workflow разработки

### 1. Запуск проекта
```bash
./run
```

### 2. Разработка
- **Backend:** Редактируй `.go` файлы → Air автоматически перезапустит
- **Frontend:** Редактируй `.tsx` файлы → Hot reload автоматически
- **Mobile:** Редактируй `.tsx` файлы → Hot reload автоматически

### 3. База данных
- Создается автоматически при первом запуске backend
- Файл: `database/data.db`
- Добавляй миграции в `database/db.go` → `createTables()`

### 4. Shared логика
- Добавляй типы в `shared/src/types/`
- Добавляй хуки в `shared/src/hooks/`
- Экспортируй в `shared/src/index.ts`
- Импортируй в frontend: `import { ... } from '../../../shared/src'`

### 5. Релиз новой версии
```
Команда: "релиз"
```
Автоматически:
- Обновит CHANGELOG.md
- Создаст git commit и tag
- Инкрементирует версию

---

## Архитектурные решения

### Почему database отдельный модуль?
- Возможность переиспользования в других сервисах
- Четкое разделение ответственности
- Упрощенное тестирование

### Почему shared не используется в mobile?
- Metro bundler имеет ограничения с путями вне проекта
- API клиент продублирован для простоты
- В будущем можно настроить через monorepo (Turborepo/Nx)

### Почему SQLite?
- Не требует отдельного сервера
- Идеален для разработки и MVP
- Легко мигрировать на PostgreSQL/MySQL позже

### Почему три отдельных приложения?
- Web и Mobile имеют разные UI компоненты
- Общая логика вынесена в shared
- Гибкость в развитии каждой платформы

---

## Следующие шаги

### Backend
- [ ] Добавить роутинг (chi/gin)
- [ ] JWT авторизация
- [ ] CRUD endpoints для users
- [ ] Валидация данных
- [ ] Логирование

### Frontend/Mobile
- [ ] Навигация (React Navigation)
- [ ] State management (Zustand)
- [ ] Формы с валидацией
- [ ] UI библиотека компонентов
- [ ] Темная тема

### DevOps
- [ ] Docker контейнеризация
- [ ] CI/CD pipeline
- [ ] Тесты
- [ ] Production конфигурация

### Database
- [ ] Миграции
- [ ] Seeders
- [ ] Модели данных
- [ ] Связи между таблицами
