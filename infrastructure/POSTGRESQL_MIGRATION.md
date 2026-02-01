# Миграция на PostgreSQL

## 📊 Обзор

Проект был перенесен с SQLite на PostgreSQL для лучшей масштабируемости и производительности в production окружении.

## 🔧 Конфигурация

### Credentials базы данных (EasyPanel)

```
Пользователь: postgres_zp
Пароль: 7da0905cd3349f58f368
База данных: bd_zp
Хост: my_projects_bd_zooplatforma
Порт: 5432
Connection String: postgres://postgres_zp:7da0905cd3349f58f368@my_projects_bd_zooplatforma:5432/bd_zp?sslmode=disable
```

### Переменные окружения

Добавь в `.env` файл каждого backend сервиса:

```bash
# Вариант 1: Использовать полный Connection String
DATABASE_URL=postgres://postgres_zp:7da0905cd3349f58f368@my_projects_bd_zooplatforma:5432/bd_zp?sslmode=disable

# Вариант 2: Использовать отдельные переменные
DB_HOST=my_projects_bd_zooplatforma
DB_PORT=5432
DB_USER=postgres_zp
DB_PASSWORD=7da0905cd3349f58f368
DB_NAME=bd_zp
DB_SSLMODE=disable
```

## 📝 Изменения в коде

### 1. Драйвер БД

**Было (SQLite):**
```go
import _ "github.com/mattn/go-sqlite3"
```

**Стало (PostgreSQL):**
```go
import _ "github.com/lib/pq"
```

### 2. Инициализация БД

**Было (SQLite):**
```go
func InitDB() error {
    return InitDBWithPath("../../database/data.db")
}
```

**Стало (PostgreSQL):**
```go
func InitDB() error {
    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        dbURL = "postgres://postgres_zp:7da0905cd3349f58f368@localhost:5432/bd_zp?sslmode=disable"
    }
    return InitDBWithURL(dbURL)
}
```

### 3. Connection String

**SQLite:**
```
data.db?_busy_timeout=5000&_journal_mode=WAL&cache=shared
```

**PostgreSQL:**
```
postgres://user:password@host:port/database?sslmode=disable
```

### 4. Типы данных

| SQLite | PostgreSQL |
|--------|-----------|
| INTEGER PRIMARY KEY AUTOINCREMENT | SERIAL PRIMARY KEY |
| DATETIME | TIMESTAMP |
| TEXT | TEXT |
| PRAGMA | Не используется |

### 5. Синтаксис SQL

**Было (SQLite):**
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

**Стало (PostgreSQL):**
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🚀 Процесс миграции

### Шаг 1: Обновить зависимости

```bash
# Обновить go.mod во всех backend сервисах
go get github.com/lib/pq@latest

# Скачать зависимости
go mod download
```

### Шаг 2: Применить миграцию

```bash
# Подключиться к PostgreSQL
psql -h my_projects_bd_zooplatforma -U postgres_zp -d bd_zp

# Или использовать Connection String
psql "postgres://postgres_zp:7da0905cd3349f58f368@my_projects_bd_zooplatforma:5432/bd_zp?sslmode=disable"

# Применить миграцию
\i database/migrations/036_migrate_to_postgresql.sql
```

### Шаг 3: Обновить конфигурацию

Добавь `DATABASE_URL` в `.env` файл каждого backend сервиса:

```bash
# main/backend/.env
DATABASE_URL=postgres://postgres_zp:7da0905cd3349f58f368@my_projects_bd_zooplatforma:5432/bd_zp?sslmode=disable

# admin/backend/.env
DATABASE_URL=postgres://postgres_zp:7da0905cd3349f58f368@my_projects_bd_zooplatforma:5432/bd_zp?sslmode=disable

# И так для всех остальных сервисов...
```

### Шаг 4: Перезапустить сервисы

```bash
# Остановить все сервисы
make stop

# Запустить все сервисы
make all
```

## ✅ Проверка

### Проверить подключение

```bash
# Подключиться к БД
psql "postgres://postgres_zp:7da0905cd3349f58f368@my_projects_bd_zooplatforma:5432/bd_zp?sslmode=disable"

# Проверить таблицы
\dt

# Проверить индексы
\di

# Выход
\q
```

### Проверить логи сервисов

```bash
# Main backend
tail -f logs/requests/main.log

# Admin backend
tail -f logs/requests/admin.log

# И так для всех остальных...
```

## 🔄 Откат (если что-то пошло не так)

### Вернуться на SQLite

1. Остановить все сервисы
2. Восстановить старый `database/db.go` из git
3. Восстановить старые `go.mod` файлы
4. Удалить `DATABASE_URL` из `.env` файлов
5. Перезапустить сервисы

```bash
# Откатить изменения
git checkout database/db.go
git checkout */backend/go.mod

# Перезапустить
make stop
make all
```

## 📊 Различия между SQLite и PostgreSQL

### SQLite (локальная разработка)

✅ Простая настройка  
✅ Не требует отдельного сервера  
✅ Хорошо для разработки  
❌ Плохо масштабируется  
❌ Ограниченная конкурентность  
❌ Не подходит для production  

### PostgreSQL (production)

✅ Отличная масштабируемость  
✅ Высокая конкурентность  
✅ Надежность и безопасность  
✅ Подходит для production  
❌ Требует отдельного сервера  
❌ Сложнее настраивать  

## 🆘 Troubleshooting

### Ошибка: "connection refused"

**Причина:** PostgreSQL сервер недоступен

**Решение:**
```bash
# Проверить что PostgreSQL запущен на EasyPanel
# Проверить credentials
# Проверить что хост и порт правильные
```

### Ошибка: "database does not exist"

**Причина:** База данных не создана

**Решение:**
```bash
# Создать базу данных
psql -h my_projects_bd_zooplatforma -U postgres_zp -c "CREATE DATABASE bd_zp;"
```

### Ошибка: "permission denied"

**Причина:** Неправильный пользователь или пароль

**Решение:**
```bash
# Проверить credentials в EasyPanel
# Убедиться что пользователь имеет права на базу данных
```

### Ошибка: "SSL error"

**Причина:** SSL требуется для подключения

**Решение:**
```bash
# Добавить ?sslmode=require в Connection String
# Или использовать ?sslmode=disable если SSL не требуется
```

## 📚 Дополнительные ресурсы

- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Go pq Driver](https://github.com/lib/pq)
- [Database/sql Package](https://golang.org/pkg/database/sql/)
- [EasyPanel PostgreSQL Setup](https://easypanel.io/)

## 🎯 Следующие шаги

1. ✅ Обновить все backend сервисы
2. ✅ Применить миграцию БД
3. ✅ Протестировать все сервисы
4. ✅ Развернуть на production
5. ⏳ Мониторить производительность
6. ⏳ Оптимизировать запросы если нужно
