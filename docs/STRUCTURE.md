# Структура проекта

Полное описание архитектуры и организации файлов проекта.

## Общая структура (Микросервисная архитектура)

```
.
├── main/                 # Главный сервис (платформа)
│   ├── backend/          # Go Main API сервер (порт 8000)
│   └── frontend/         # Next.js веб-приложение (порт 3000)
├── admin/                # Административная панель (микросервис)
│   ├── backend/          # Admin API сервер (порт 9000)
│   └── frontend/         # Admin UI (Next.js, порт 4000)
├── petbase/              # ЗооБаза - справочник животных (микросервис)
│   ├── backend/          # PetBase API сервер (порт 8100)
│   └── frontend/         # PetBase UI (Next.js, порт 4100)
├── mobile/               # React Native мобильное приложение (порт 8081)
├── shared/               # Общие компоненты (@pet/shared npm пакет)
│   ├── src/
│   │   ├── api/          # API клиент
│   │   ├── components/   # Переиспользуемые компоненты (AuthForm)
│   │   ├── types/        # TypeScript типы
│   │   └── hooks/        # React хуки
│   └── package.json      # npm пакет
├── database/             # SQLite база данных (общая для всех)
├── docs/                 # Документация проекта
│   ├── README.md         # Обзор документации
│   ├── CHANGELOG.md      # История изменений
│   ├── ROADMAP.md        # План развития
│   ├── STRUCTURE.md      # Этот файл
│   ├── PERFORMANCE.md    # Метрики производительности
│   ├── TEMPLATES.md      # Шаблоны компонентов
│   ├── SSO_ARCHITECTURE.md  # Архитектура SSO
│   ├── SSO_FLOW.md       # Схема SSO
│   └── API_KEYS.md       # API ключи (не в git)
├── tests/                # Тесты и скрипты отладки
├── .kiro/                # Конфигурация Kiro IDE
│   └── steering/         # Правила и гайды разработки
├── run                   # Скрипт запуска всех сервисов
├── .gitignore           # Игнорируемые файлы
└── README.md            # Основная документация
```

---

## Main Service (Главный сервис)

### Backend (Go API)

**Путь:** `main/backend/`

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

**Порт:** 8000

**Endpoints:**
- `GET /` - Welcome message
- `GET /api/health` - Health check
- `POST /api/auth/register` - Регистрация
- `POST /api/auth/login` - Вход
- `POST /api/auth/logout` - Выход
- `GET /api/auth/me` - Текущий пользователь
- `GET /api/auth/verify` - Проверка токена (для SSO)
- `GET /api/users` - Список пользователей (защищено)
- `PUT /api/profile` - Обновление профиля (защищено)
- `GET/POST/PUT/DELETE /api/posts` - CRUD постов (защищено)
- `GET /api/posts/drafts` - Получить черновики (защищено)
- `POST /api/posts/:id/like` - Лайк/анлайк поста (защищено)
- `GET /api/posts/:id/like` - Статус лайка (защищено)
- `GET/POST/DELETE /api/pets` - CRUD питомцев (защищено)
- `GET/POST/PUT/DELETE /api/polls` - CRUD опросов (защищено)
- `POST /api/polls/:id/vote` - Голосование (защищено)
- `DELETE /api/polls/:id/vote` - Отмена голоса (защищено)
- `GET/POST/DELETE /api/comments` - CRUD комментариев (защищено)

**Особенности:**
- CORS настроен для всех клиентов
- JWT авторизация с ролями (user, superadmin)
- Автоматический hot reload через Air
- Интеграция с модулем database
- Логирование событий в system_logs

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
- `users` - Пользователи платформы
  - id, name, email, password, bio, phone, location, avatar, cover_photo, created_at
- `posts` - Публикации пользователей (универсальная система в стиле Threads)
  - id, author_id, author_type, content, attached_pets, attachments, tags, status, scheduled_at, created_at, updated_at, is_deleted
- `pets` - Питомцы пользователей
  - id, user_id, name, species, photo, created_at
- `likes` - Лайки постов
  - id, user_id, post_id, created_at
- `comments` - Комментарии к постам
  - id, post_id, user_id, parent_id, reply_to_user_id, content, created_at, updated_at, is_deleted
- `polls` - Опросы в постах
  - id, post_id, question, multiple_choice, allow_vote_changes, anonymous_voting, expires_at, created_at
- `poll_options` - Варианты ответов в опросах
  - id, poll_id, option_text, display_order
- `poll_votes` - Голоса в опросах
  - id, poll_id, option_id, user_id, created_at
- `admins` - Администраторы
  - id, user_id, role, permissions, created_at, created_by
- `admin_logs` - Логи действий администраторов
  - id, admin_id, action, target_type, target_id, details, ip_address, created_at
- `system_logs` - Системные логи всей платформы
  - id, level, category, action, user_id, target_type, target_id, message, details, ip_address, user_agent, created_at

**Особенности:**
- Отдельный модуль для переиспользования
- Автоматическое создание таблиц при первом запуске
- База создается в `database/data.db`

---

## Admin Panel (Микросервис)

**Путь:** `admin/`

### Admin Backend

**Путь:** `admin/backend/`

**Технологии:**
- Go 1.25.5
- net/http (standard library)
- Air (hot reload)

**Структура:**
```
admin/backend/
├── main.go              # Точка входа, HTTP сервер
├── handlers/            # HTTP handlers
│   ├── auth.go          # Авторизация админов
│   ├── users.go         # Управление пользователями
│   ├── posts.go         # Модерация постов
│   ├── stats.go         # Статистика
│   └── logs.go          # Системные логи
├── middleware/          # Middleware
│   └── admin.go         # Проверка прав суперадмина
├── go.mod               # Go модуль
├── go.sum               # Checksums
├── .air.toml            # Конфигурация Air
├── create-superadmin.sh # Скрипт создания суперадмина
└── tmp/                 # Временные файлы Air
```

**Порт:** 9000

**Endpoints:**
- `GET /api/admin/health` - Health check
- `GET /api/admin/auth/me` - Текущий админ
- `GET /api/admin/users` - Список пользователей
- `PUT /api/admin/users/:id` - Обновление пользователя
- `DELETE /api/admin/users/:id` - Удаление пользователя
- `GET /api/admin/posts` - Список постов
- `PUT /api/admin/posts/:id` - Модерация поста
- `DELETE /api/admin/posts/:id` - Удаление поста
- `GET /api/admin/stats/overview` - Статистика
- `GET /api/admin/logs` - Системные логи

**Особенности:**
- SSO через общий токен `auth_token`
- Проверка ролей из JWT (superadmin)
- Общая база данных с Main Backend
- Hot reload через Air

### Admin Frontend

**Путь:** `admin/frontend/`

**Технологии:**
- Next.js 15.1.3
- React 19.0.0
- TypeScript 5.x
- Tailwind CSS 4.x
- Heroicons 2.x

**Структура:**
```
admin/frontend/
├── app/
│   ├── layout.tsx       # Корневой layout
│   ├── page.tsx         # Дашборд
│   ├── page.module.css  # Стили дашборда
│   ├── globals.css      # Глобальные стили
│   └── auth/            # Страница входа
│       ├── page.tsx
│       └── auth.module.css
├── lib/
│   └── api.ts           # API клиент для Admin Backend
├── package.json
├── tsconfig.json
├── tailwind.config.js
├── next.config.ts
└── .env.local           # NEXT_PUBLIC_API_URL=http://localhost:9000
```

**Порт:** 4000

**Страницы:**
- `/` - Дашборд с вкладками
- `/auth` - Страница входа

**Вкладки дашборда:**
- Пользователи - управление пользователями
- Посты - модерация постов
- Статистика - метрики платформы
- Логирование - системные логи

**Особенности:**
- Дизайн в стиле VK Реклама
- Сворачиваемое боковое меню
- SSO авторизация через Main Backend
- Tailwind CSS + Heroicons
- Hot reload из коробки

---

## PetBase Service (ЗооБаза)

### Backend (Go API)

**Путь:** `petbase/backend/`

**Порт:** 8100

**Endpoints:**
- `GET /api/species` - Список видов животных
- `POST /api/species` - Создать вид
- `GET /api/species/:id` - Получить вид
- `PUT /api/species/:id` - Обновить вид
- `DELETE /api/species/:id` - Удалить вид
- `GET /api/breeds` - Список пород
- `GET /api/breeds/species/:id` - Породы по виду
- `POST /api/breeds` - Создать породу
- `GET /api/breeds/:id` - Получить породу
- `PUT /api/breeds/:id` - Обновить породу
- `DELETE /api/breeds/:id` - Удалить породу
- `GET /api/cards` - Список карточек
- `GET /api/cards/breed/:id` - Карточки по породе
- `POST /api/cards` - Создать карточку
- `GET /api/cards/:id` - Получить карточку
- `PUT /api/cards/:id` - Обновить карточку
- `DELETE /api/cards/:id` - Удалить карточку

**Структура:**
```
petbase/backend/
├── main.go              # Точка входа
├── handlers/
│   ├── species.go       # CRUD для видов
│   ├── breeds.go        # CRUD для пород
│   └── cards.go         # CRUD для карточек
├── models/
│   ├── species.go       # Модель вида
│   ├── breed.go         # Модель породы
│   └── pet_card.go      # Модель карточки
├── seed.sql             # Тестовые данные
├── go.mod
└── go.sum
```

**Таблицы БД:**
- `species` - виды животных (собаки, кошки, птицы и т.д.)
- `breeds` - породы с характеристиками
- `pet_cards` - подробные карточки животных

### Frontend (Next.js Admin UI)

**Путь:** `petbase/frontend/`

**Порт:** 4100

**Страницы:**
- `/` - Дашборд с вкладками
- `/auth` - Страница входа (проверка прав суперадмина)

**Вкладки дашборда:**
- Статистика - общая информация
- Виды - управление видами животных
- Породы - управление породами
- Карточки - управление карточками

**Особенности:**
- Использует AdminLayout из shared (скопирован локально)
- Доступ только для суперадминов
- SSO интеграция через Main + Admin Backend
- Tailwind CSS v4 + Heroicons
- Hot reload из коробки

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
├── app/
│   ├── (main)/          # Основные страницы
│   │   ├── layout.tsx   # Layout с Header и Sidebar
│   │   ├── page.tsx     # Главная страница (лента)
│   │   ├── org/         # Организации
│   │   └── [userId]/    # Профили пользователей
│   ├── components/
│   │   ├── layout/      # Layout компоненты
│   │   │   ├── Header.tsx
│   │   │   ├── Sidebar.tsx
│   │   │   └── RightPanel.tsx
│   │   ├── posts/       # Компоненты постов
│   │   │   ├── CreatePost.tsx    # Создание поста (модальное окно)
│   │   │   ├── PostCard.tsx      # Карточка поста
│   │   │   ├── PostModal.tsx     # Модальное окно поста
│   │   │   ├── PostsFeed.tsx     # Лента постов
│   │   │   └── FeedFilters.tsx   # Фильтры ленты
│   │   ├── polls/       # Компоненты опросов
│   │   │   ├── PollCreator.tsx   # Создание опроса
│   │   │   └── PollDisplay.tsx   # Отображение опроса
│   │   ├── shared/      # Общие компоненты
│   │   │   └── PostComments.tsx  # Комментарии
│   │   └── ui/          # Базовые UI компоненты
│   ├── contexts/        # React контексты
│   │   └── AuthContext.tsx
│   ├── lib/
│   │   └── api.ts       # API клиент
│   ├── layout.tsx       # Корневой layout
│   └── globals.css      # Глобальные стили
├── public/              # Статические файлы
│   └── favicon.svg      # Логотип
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

---

## Tests (Тесты)

**Путь:** `tests/`

**Структура:**
```
tests/
├── performance-test.sh      # Тест производительности
├── compare-with-vk.sh       # Сравнение с VK.com
└── README.md                # Документация по тестам
```

**Типы тестов:**

### 1. Тесты производительности (сейчас)
- `performance-test.sh` - проверка скорости API и страниц
- `compare-with-vk.sh` - сравнение с VK.com

### 2. Unit тесты (планируется для 0.1.0)
- Backend: `go test ./...`
- Frontend: `npm test` (Jest)

### 3. Integration тесты (планируется для 0.2.0)
- API integration tests
- Database tests

### 4. E2E тесты (планируется для 1.0.0)
- Playwright/Cypress
- Критичные пользовательские сценарии

**Запуск:**
```bash
# Тест производительности
./tests/performance-test.sh

# Сравнение с VK.com
./tests/compare-with-vk.sh
```

**Документация:** [tests/README.md](tests/README.md)

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

| Сервис          | Порт | URL                      |
|-----------------|------|--------------------------|
| Main Backend    | 8000 | http://localhost:8000    |
| Admin Backend   | 9000 | http://localhost:9000    |
| PetBase Backend | 8100 | http://localhost:8100    |
| Web Frontend    | 3000 | http://localhost:3000    |
| Admin Frontend  | 4000 | http://localhost:4000    |
| PetBase Frontend| 4100 | http://localhost:4100    |
| Mobile          | 8081 | http://localhost:8081    |

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
