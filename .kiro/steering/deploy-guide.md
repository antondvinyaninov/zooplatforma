# Деплой на EasyPanel

## 🌐 Адреса и доступы

### Production URL
- **Основной сайт:** https://my-projects-zooplatforma.crv1ic.easypanel.host
- **API Gateway:** https://my-projects-zooplatforma.crv1ic.easypanel.host/api
- **Статические файлы:** https://my-projects-zooplatforma.crv1ic.easypanel.host/uploads

### Внутренние порты (внутри Docker контейнера)
- **nginx:** 80 (API Gateway)
- **Auth Service:** 7100
- **Main Backend:** 8000
- **PetBase Backend:** 8100
- **Main Frontend:** 3000 (Next.js production)
- **PetBase Frontend:** 4100 (Next.js production)

### База данных PostgreSQL
- **Сервис:** zooplatforma-db
- **Хост:** zooplatforma-db (внутри Docker сети)
- **База данных:** zp-db
- **Пользователь:** zp
- **Пароль:** lmLG7k2ed4vas19
- **Порт:** 5432

**Подключение из контейнера:**
```bash
PGPASSWORD=lmLG7k2ed4vas19 psql -h zooplatforma-db -U zp -d zp-db
```

---

## 🏗️ Архитектура деплоя

### Один Docker контейнер со всеми сервисами

```
Docker Container (zooplatforma)
├── nginx (порт 80) - API Gateway
│   ├── / → Main Frontend (3000)
│   ├── /api/* → Main Backend (8000)
│   └── /uploads/* → /app/uploads (статика)
│
├── Auth Service (порт 7100)
│   └── PostgreSQL auth.db → zooplatforma-db
│
├── PetBase Backend (порт 8100)
│   └── Справочник животных (виды, породы, карточки)
│
├── Main Backend (порт 8000)
│   └── Основной API (посты, пользователи, организации)
│
├── Main Frontend (порт 3000)
│   └── Next.js production build
│
└── PetBase Frontend (порт 4100)
    └── Next.js production build (админка справочника)
```

### Внешняя база данных
```
PostgreSQL Container (zooplatforma-db)
└── База данных: zp-db
    ├── Таблицы пользователей
    ├── Таблицы постов
    ├── Таблицы организаций
    ├── Таблицы животных
    └── Таблицы сообщений
```

---

## 📁 Структура файлов в контейнере

```
/app/
├── uploads/                    # Общая папка для всех загруженных файлов
│   ├── users/                  # Аватарки и обложки пользователей
│   │   ├── 1/
│   │   │   ├── avatars/
│   │   │   └── covers/
│   │   └── 2/
│   ├── messages/               # Вложения в сообщениях
│   ├── posts/                  # Медиа в постах
│   ├── pets/                   # Фото питомцев
│   └── temp/                   # Временные файлы
│
├── auth-backend                # Auth Service бинарник
├── main-backend                # Main Backend бинарник
├── petbase-backend             # PetBase Backend бинарник
│
├── frontend/                   # Main Frontend (Next.js production)
│   ├── .next/
│   ├── node_modules/
│   └── package.json
│
├── petbase-frontend/           # PetBase Frontend (Next.js production)
│   ├── .next/
│   ├── node_modules/
│   └── package.json
│
├── migrations/                 # SQL миграции
├── fix_organizations_table.sql
├── fix_posts_table.sql
└── start.sh                    # Скрипт запуска всех сервисов
```

---

## 🔧 Переменные окружения

### В Docker контейнере (Dockerfile)

```bash
# Database
DATABASE_HOST=zooplatforma-db
DATABASE_PORT=5432
DATABASE_USER=zp
DATABASE_PASSWORD=lmLG7k2ed4vas19
DATABASE_NAME=zp-db

# Auth Service
AUTH_SERVICE_URL=http://localhost:7100

# Uploads
UPLOAD_PATH=/app/uploads

# Environment
ENVIRONMENT=production

# DaData API (для автокомплита городов)
NEXT_PUBLIC_DADATA_API_KEY=300ba9e25ef32f0d6ea7c41826b2255b138e19e2
```

### В Next.js Frontend

```bash
# Main Frontend
PORT=3000
NEXT_PUBLIC_API_URL=http://localhost:8000
NEXT_PUBLIC_AUTH_URL=http://localhost:7100

# PetBase Frontend
PORT=4100
```

---

## 🚀 Процесс деплоя

### 1. Локальная подготовка

```bash
# Проверить что всё работает локально
npm run build  # в main/frontend
go build       # в main/backend

# Закоммитить изменения
git add -A
git commit -m "Your changes"
git push origin main
```

### 2. EasyPanel автоматически:

1. **Получает изменения** из GitHub (main branch)
2. **Собирает Docker образ** (~8-10 минут):
   - Go builder stage (компиляция всех backend сервисов)
   - Next.js builder stage (сборка фронтендов)
   - Runtime stage (финальный образ)
3. **Запускает контейнер** с `start.sh`
4. **Применяет SQL fixes** (если ENVIRONMENT=production)
5. **Запускает все сервисы** параллельно

### 3. Проверка деплоя

```bash
# Проверить что сайт доступен
curl https://my-projects-zooplatforma.crv1ic.easypanel.host

# Проверить API
curl https://my-projects-zooplatforma.crv1ic.easypanel.host/api/health

# Проверить Auth Service (через API Gateway)
curl https://my-projects-zooplatforma.crv1ic.easypanel.host/api/auth/health
```

---

## 🐛 Отладка проблем

### Логи контейнера

В EasyPanel:
1. Открыть проект "zooplatforma"
2. Перейти в "Logs"
3. Смотреть вывод всех сервисов

### Подключение к PostgreSQL

Из локальной машины (если открыт порт):
```bash
PGPASSWORD=lmLG7k2ed4vas19 psql -h <easypanel-host> -U zp -d zp-db
```

Из контейнера:
```bash
PGPASSWORD=lmLG7k2ed4vas19 psql -h zooplatforma-db -U zp -d zp-db
```

### Проверка файлов в контейнере

В EasyPanel можно открыть терминал контейнера:
```bash
# Проверить что папка uploads существует
ls -la /app/uploads

# Проверить права доступа
ls -la /app/uploads/users

# Проверить что nginx запущен
ps aux | grep nginx

# Проверить что все сервисы запущены
ps aux | grep backend
```

---

## 📝 Частые проблемы и решения

### 1. 404 на загруженные файлы

**Проблема:** Файлы загружаются, но nginx возвращает 404

**Решение:**
- Проверить что nginx запускается от root: `user root;` в nginx.conf
- Проверить что папка `/app/uploads` создана с правами 777
- Проверить что `UPLOAD_PATH=/app/uploads` установлен для backend сервисов

### 2. PostgreSQL connection refused

**Проблема:** Backend не может подключиться к PostgreSQL

**Решение:**
- Проверить что сервис `zooplatforma-db` запущен
- Проверить переменные окружения: `DATABASE_HOST=zooplatforma-db`
- Проверить что контейнеры в одной Docker сети

### 3. 401 Unauthorized на всех запросах

**Проблема:** Auth Service не работает

**Решение:**
- Проверить что Auth Service запущен: `ps aux | grep auth-backend`
- Проверить что `AUTH_SERVICE_URL=http://localhost:7100` установлен
- Проверить логи Auth Service

### 4. Next.js не запускается

**Проблема:** Frontend возвращает ошибку

**Решение:**
- Проверить что `npm run build` прошёл успешно локально
- Проверить что все зависимости установлены
- Проверить переменные окружения Next.js

### 5. Долгий деплой (>10 минут)

**Причины:**
- Большой размер проекта (6 сервисов)
- Компиляция Go кода
- Сборка Next.js production build
- Установка npm зависимостей

**Нормально:** 8-10 минут для полного деплоя

---

## 🔄 Обновление production

### Обычное обновление (код)

```bash
git add -A
git commit -m "Update: description"
git push origin main
# EasyPanel автоматически задеплоит
```

### Обновление с миграцией БД

1. **Создать SQL файл** с миграцией
2. **Добавить в Dockerfile** (если нужно применить автоматически)
3. **Или применить вручную** через psql:

```bash
# Подключиться к БД
PGPASSWORD=lmLG7k2ed4vas19 psql -h zooplatforma-db -U zp -d zp-db

# Применить миграцию
\i /path/to/migration.sql

# Или скопировать SQL и выполнить
```

### Откат изменений

```bash
# Откатить последний коммит
git revert HEAD
git push origin main

# Или откатить к конкретному коммиту
git reset --hard <commit-hash>
git push origin main --force
```

---

## 📊 Мониторинг

### Проверка здоровья сервисов

```bash
# Main Backend
curl https://my-projects-zooplatforma.crv1ic.easypanel.host/api/health

# Auth Service (через API Gateway)
curl https://my-projects-zooplatforma.crv1ic.easypanel.host/api/auth/health

# PetBase Backend
curl https://my-projects-zooplatforma.crv1ic.easypanel.host/api/species
```

### Метрики

- **Время деплоя:** 8-10 минут
- **Размер Docker образа:** ~1.5GB
- **Количество процессов:** 6 (nginx + 2 backend + 2 frontend + auth)
- **Использование памяти:** ~1GB

---

## 🔐 Безопасность

### Что защищено

- ✅ PostgreSQL доступен только внутри Docker сети
- ✅ Все API запросы идут через nginx (API Gateway)
- ✅ JWT токены для авторизации
- ✅ HTTPS на production (через EasyPanel)

### Что нужно улучшить

- ⚠️ Добавить rate limiting в nginx
- ⚠️ Добавить CORS политики
- ⚠️ Добавить мониторинг и алерты
- ⚠️ Настроить автоматические бэкапы БД

---

## 📚 Полезные команды

### Git

```bash
# Посмотреть последние коммиты
git log --oneline -10

# Посмотреть изменения
git diff

# Откатить файл
git checkout -- <file>
```

### Docker (локально)

```bash
# Собрать образ
docker build -t zooplatforma .

# Запустить контейнер
docker run -p 80:80 zooplatforma

# Посмотреть логи
docker logs <container-id>

# Зайти в контейнер
docker exec -it <container-id> sh
```

### PostgreSQL

```bash
# Список таблиц
\dt

# Описание таблицы
\d table_name

# Выполнить SQL файл
\i /path/to/file.sql

# Выход
\q
```

---

## 📞 Контакты и ссылки

- **GitHub:** https://github.com/antondvinyaninov/zooplatforma
- **EasyPanel:** http://88.218.121.213:3000
- **Production:** https://my-projects-zooplatforma.crv1ic.easypanel.host

---

**Последнее обновление:** 2026-02-03


---

## ⚠️ ТИПИЧНЫЕ ОШИБКИ ПРИ РАЗРАБОТКЕ

### 🔴 ПРОБЛЕМА #1: SQLite vs PostgreSQL

**Описание:**
Локально используется SQLite (`database/data.db`), в production - PostgreSQL. Синтаксис SQL отличается!

**Частые ошибки:**

1. **Плейсхолдеры:**
   ```go
   // ❌ НЕПРАВИЛЬНО (SQLite)
   db.Query("SELECT * FROM users WHERE id = ?", userID)
   
   // ✅ ПРАВИЛЬНО (PostgreSQL)
   db.Query(ConvertPlaceholders("SELECT * FROM users WHERE id = ?"), userID)
   // ConvertPlaceholders заменяет ? на $1, $2, $3...
   ```

2. **LastInsertId() не работает в PostgreSQL:**
   ```go
   // ❌ НЕПРАВИЛЬНО
   result, err := db.Exec("INSERT INTO posts (...) VALUES (...)")
   postID, _ := result.LastInsertId() // Ошибка в PostgreSQL!
   
   // ✅ ПРАВИЛЬНО
   var postID int64
   err := db.QueryRow(ConvertPlaceholders(`
       INSERT INTO posts (...) VALUES (...)
       RETURNING id
   `), ...).Scan(&postID)
   ```

3. **Булевы значения:**
   ```go
   // ❌ НЕПРАВИЛЬНО (SQLite)
   db.Exec("UPDATE users SET is_active = 1 WHERE id = ?", userID)
   
   // ✅ ПРАВИЛЬНО (PostgreSQL)
   db.Exec(ConvertPlaceholders("UPDATE users SET is_active = TRUE WHERE id = ?"), userID)
   ```

4. **Функции даты/времени:**
   ```go
   // ❌ НЕПРАВИЛЬНО (SQLite)
   db.Query("SELECT * FROM posts WHERE created_at > datetime('now', '-7 days')")
   
   // ✅ ПРАВИЛЬНО (PostgreSQL)
   db.Query("SELECT * FROM posts WHERE created_at > NOW() - INTERVAL '7 days'")
   ```

5. **NULL значения:**
   ```go
   // ❌ НЕПРАВИЛЬНО
   var avatar string
   db.QueryRow("SELECT avatar FROM users WHERE id = ?", userID).Scan(&avatar)
   // Ошибка если avatar = NULL!
   
   // ✅ ПРАВИЛЬНО
   var avatar sql.NullString
   db.QueryRow(ConvertPlaceholders("SELECT avatar FROM users WHERE id = ?"), userID).Scan(&avatar)
   if avatar.Valid {
       // Используем avatar.String
   }
   ```

6. **Ошибка: `sql: Scan error on column index X: unsupported Scan, storing driver.Value type <nil> into type *string`**
   
   **Причина:** Попытка записать NULL значение из PostgreSQL в обычную Go переменную.
   
   **Где возникает:**
   - При чтении nullable полей из таблицы `users`: `last_name`, `avatar`, `cover_photo`, `bio`, `location`, `phone`
   - При чтении nullable полей из других таблиц
   
   **Решение:**
   ```go
   // ❌ НЕПРАВИЛЬНО:
   var lastName string
   db.QueryRow("SELECT last_name FROM users WHERE id = ?", userID).Scan(&lastName)
   
   // ✅ ПРАВИЛЬНО:
   var lastName sql.NullString
   db.QueryRow("SELECT last_name FROM users WHERE id = ?", userID).Scan(&lastName)
   if lastName.Valid {
       // Используем lastName.String
   }
   ```
   
   **Файлы где была эта ошибка:**
   - `main/backend/handlers/profile.go` - ИСПРАВЛЕНО ✅
   - `main/backend/handlers/users.go` - ИСПРАВЛЕНО ✅
   - `main/backend/handlers/announcements.go` - ИСПРАВЛЕНО ✅
   - `main/backend/handlers/friends.go` - ИСПРАВЛЕНО ✅
   
   **Как проверить:**
   ```bash
   # Поиск потенциальных проблем
   grep -r "last_name string" main/backend/handlers/
   grep -r "avatar string" main/backend/handlers/
   grep -r "bio string" main/backend/handlers/
   ```

**Решение:**
- ✅ ВСЕГДА оборачивай SQL запросы в `ConvertPlaceholders()`
- ✅ Используй `RETURNING id` вместо `LastInsertId()`
- ✅ Используй `TRUE/FALSE` вместо `1/0`
- ✅ Используй `sql.NullString`, `sql.NullBool` для NULL-able полей
- ✅ Используй PostgreSQL функции для даты/времени

**Файлы где это важно:**
- `main/backend/handlers/*.go` - все handlers
- `auth/backend/handlers.go`
- `petbase/backend/handlers/*.go`

---

### 🔴 ПРОБЛЕМА #2: Локальные адреса вместо относительных путей

**Описание:**
В коде часто прописываются `http://localhost:8000`, `http://localhost:7100` и т.д. В production это не работает!

**Частые ошибки:**

1. **Хардкод localhost в API запросах:**
   ```typescript
   // ❌ НЕПРАВИЛЬНО
   fetch('http://localhost:8000/api/posts')
   
   // ✅ ПРАВИЛЬНО
   fetch('/api/posts') // Относительный путь через nginx
   ```

2. **Хардкод localhost в URL изображений:**
   ```typescript
   // ❌ НЕПРАВИЛЬНО
   setAvatarPreview(`http://localhost:8000${response.data.avatar_url}`)
   
   // ✅ ПРАВИЛЬНО
   setAvatarPreview(response.data.avatar_url) // Уже относительный путь
   ```

3. **Хардкод localhost в SEO meta-тегах:**
   ```typescript
   // ❌ НЕПРАВИЛЬНО
   images: [`http://localhost:8000${image}`]
   
   // ✅ ПРАВИЛЬНО
   images: [image] // Относительный путь
   ```

4. **Хардкод localhost в API клиентах:**
   ```typescript
   // ❌ НЕПРАВИЛЬНО
   const API_URL = 'http://localhost:8000';
   
   // ✅ ПРАВИЛЬНО
   const API_URL = process.env.NEXT_PUBLIC_API_URL || '';
   // В production будет пустая строка = относительные пути
   ```

**Решение:**
- ✅ ВСЕГДА используй относительные пути: `/api/...`, `/uploads/...`
- ✅ В production nginx проксирует всё на правильные сервисы
- ✅ Проверяй перед коммитом: `grep -r "localhost:8000" main/frontend/`
- ✅ Используй переменные окружения для API URL (но в production они пустые)

**Архитектура в production:**
```
User → nginx (80) → /api/* → Main Backend (8000)
                  → /uploads/* → /app/uploads
                  → /* → Main Frontend (3000)
```

**Файлы где это важно:**
- `main/frontend/app/**/*.tsx` - все страницы и компоненты
- `main/frontend/lib/api.ts` - API клиент
- `main/frontend/next.config.ts` - rewrites для API

---

### 🔴 ПРОБЛЕМА #3: Пути к файлам uploads

**Описание:**
Локально: `uploads/` в корне проекта. В Docker: `/app/uploads`. Нужна гибкость!

**Частые ошибки:**

1. **Хардкод абсолютного пути:**
   ```go
   // ❌ НЕПРАВИЛЬНО
   uploadDir := "/app/uploads/users/1/avatars"
   // Не работает локально!
   
   // ✅ ПРАВИЛЬНО
   baseUploadPath := os.Getenv("UPLOAD_PATH")
   if baseUploadPath == "" {
       baseUploadPath = "../../uploads" // Локально
   }
   uploadDir := fmt.Sprintf("%s/users/%d/avatars", baseUploadPath, userID)
   ```

2. **Забыть установить UPLOAD_PATH в Docker:**
   ```bash
   # ❌ НЕПРАВИЛЬНО
   /app/main-backend
   
   # ✅ ПРАВИЛЬНО
   export UPLOAD_PATH=/app/uploads
   /app/main-backend
   ```

**Решение:**
- ✅ Используй переменную окружения `UPLOAD_PATH`
- ✅ В Docker: `UPLOAD_PATH=/app/uploads`
- ✅ Локально: по умолчанию `../../uploads`
- ✅ Все сервисы используют одну общую папку

**Файлы где это важно:**
- `main/backend/handlers/avatar.go`
- `main/backend/handlers/media.go`
- `main/backend/handlers/posts.go`
- `Dockerfile` - установка UPLOAD_PATH

---

### 🔴 ПРОБЛЕМА #4: Авторизация через localhost:8000

**Описание:**
Main Backend (8000) НЕ занимается авторизацией! Только Auth Service (7100)!

**Частые ошибки:**

1. **Запросы авторизации к Main Backend:**
   ```typescript
   // ❌ НЕПРАВИЛЬНО
   fetch('http://localhost:8000/api/auth/login')
   
   // ✅ ПРАВИЛЬНО
   fetch('http://localhost:7100/api/auth/login')
   // Или через authClient который уже настроен на 7100
   ```

2. **Проверка токена через Main Backend:**
   ```typescript
   // ❌ НЕПРАВИЛЬНО
   fetch('http://localhost:8000/api/auth/me')
   
   // ✅ ПРАВИЛЬНО
   fetch('http://localhost:7100/api/auth/me')
   ```

**Решение:**
- ✅ ВСЕ запросы авторизации → Auth Service (7100)
- ✅ Main Backend использует `pkg/middleware.AuthMiddleware`
- ✅ Middleware проверяет токен через Auth Service

**Архитектура SSO:**
```
User → Auth Service (7100) → JWT token
                ↓
        Cookie: auth_token
                ↓
    ┌───────────┴───────────┐
    ↓                       ↓
Main Backend (8000)    PetBase (8100)
    ↓                       ↓
middleware.AuthMiddleware проверяет через Auth Service
```

---

### 🔴 ПРОБЛЕМА #5: Забыть ConvertPlaceholders

**Описание:**
Самая частая ошибка! Забываешь обернуть SQL запрос в `ConvertPlaceholders()`.

**Как найти:**
```bash
# Поиск всех SQL запросов без ConvertPlaceholders
grep -r "db.Query\|db.Exec\|db.QueryRow" main/backend/handlers/ | grep -v "ConvertPlaceholders"
```

**Решение:**
- ✅ Но лучше сразу писать правильно!
- ✅ Используй helper функции из `helpers.go`

---

### 🔴 ПРОБЛЕМА #6: Изменения БД без миграций

**Описание:**
Изменяешь структуру таблиц напрямую в production через psql, но забываешь создать миграцию.

**Решение:**
- ✅ ВСЕГДА создавай миграцию: `database/migrations/XXX_description.sql`
- ✅ Применяй вручную через psql (не автоматически!)
- ✅ Документируй в DEPLOY.md что было изменено

**Процесс:**
1. Создать SQL файл с изменениями
2. Написать инструкции для пользователя
3. Пользователь применяет вручную через psql
4. Проверить что всё работает
5. Закоммитить миграцию

---

## 🛠️ Чеклист перед коммитом

- [ ] Все SQL запросы обернуты в `ConvertPlaceholders()`
- [ ] Нет `LastInsertId()` - используется `RETURNING id`
- [ ] Нет хардкода `localhost:8000` в frontend
- [ ] Нет хардкода `/app/uploads` - используется `UPLOAD_PATH`
- [ ] Авторизация идёт через Auth Service (7100)
- [ ] Используются `sql.NullString` для NULL-able полей
- [ ] Используются `TRUE/FALSE` вместо `1/0`
- [ ] Протестировано локально: `npm run build && go build`
- [ ] Проверены типы TypeScript: `npm run type-check`

---

**Последнее обновление:** 2026-02-03


---

## 🔴 КРИТИЧЕСКАЯ ОШИБКА #7: Проверка ошибки после копирования данных из sql.NullString

**Описание:**
При работе с `sql.NullString` в Go, проверка ошибки `if err != nil` должна идти **ДО** копирования данных из `NullString` в структуру, а не после!

**Неправильно:**
```go
var user User
var lastName sql.NullString

err = db.QueryRow("SELECT name, last_name FROM users WHERE id = ?", userID).Scan(&user.Name, &lastName)

// ❌ НЕПРАВИЛЬНО: копируем данные ДО проверки ошибки
if lastName.Valid {
    user.LastName = lastName.String
}

// Проверка ошибки ПОСЛЕ копирования
if err != nil {
    return err
}
```

**Почему это проблема:**
- Если `Scan()` вернул ошибку, переменные `lastName` и другие `NullString` остаются пустыми
- Мы копируем пустые значения в структуру `user`
- Потом проверяем ошибку и возвращаем её, но данные уже испорчены
- Если ошибка не критичная и мы продолжаем работу, пользователь получит пустые поля

**Правильно:**
```go
var user User
var lastName sql.NullString

err = db.QueryRow("SELECT name, last_name FROM users WHERE id = ?", userID).Scan(&user.Name, &lastName)

// ✅ ПРАВИЛЬНО: сначала проверяем ошибку
if err != nil {
    log.Printf("❌ Failed to get user: %v", err)
    return err
}

// ПОТОМ копируем данные из NullString
if lastName.Valid {
    user.LastName = lastName.String
}
```

**Где была эта ошибка:**
- `auth/backend/handlers.go` - функция `getMeHandler` (ИСПРАВЛЕНО ✅)

**Симптомы:**
- Данные сохраняются в БД правильно
- В логах видно что данные читаются из БД: `🔍 User data from DB: last_name=Двинянинов`
- Но клиент получает пустые поля: `last_name: ""`
- После обновления страницы поля пустые

**Как проверить:**
```bash
# Проверить что данные есть в БД
PGPASSWORD=xxx psql -h host -U user -d db -c "SELECT id, last_name FROM users WHERE id = 1;"

# Если в БД данные есть, но клиент получает пустые - проверь порядок проверки ошибки!
```

**Правило:** ВСЕГДА проверяй `if err != nil` СРАЗУ после `Scan()`, ДО любых операций с данными!

---

**Последнее обновление:** 2026-02-03
