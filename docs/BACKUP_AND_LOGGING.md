# Резервное копирование и логирование

Документация по системе резервного копирования базы данных и логирования.

---

## 📦 Резервное копирование базы данных

### Автоматическое резервное копирование

**Скрипт:** `scripts/backup-database.sh`

```bash
# Создать резервную копию
./scripts/backup-database.sh
```

**Что делает скрипт:**
1. Создаёт копию файла БД (`data_YYYYMMDD_HHMMSS.db`)
2. Создаёт SQL dump (`backup_YYYYMMDD_HHMMSS.sql`)
3. Проверяет целостность резервной копии
4. Показывает статистику (количество записей)
5. Удаляет старые копии (оставляет последние 10)

**Результат:**
```
database/backups/
├── data_20251228_143022.db      # Копия файла БД
├── backup_20251228_143022.sql   # SQL dump
├── data_20251228_120000.db
└── backup_20251228_120000.sql
```

### Восстановление из резервной копии

**Скрипт:** `scripts/restore-database.sh`

```bash
# Восстановить из резервной копии
./scripts/restore-database.sh
```

**Интерактивный процесс:**
1. Показывает список доступных резервных копий
2. Выбираете нужную копию (по номеру)
3. Подтверждаете восстановление
4. Скрипт создаёт safety backup текущей БД
5. Восстанавливает выбранную копию
6. Проверяет целостность

**Типы резервных копий:**
- `.db` файлы - полная копия файла БД (быстрое восстановление)
- `.sql` файлы - SQL dump (универсальное восстановление)

### Ручное резервное копирование

```bash
# Копия файла БД
cp database/data.db database/backups/manual_backup.db

# SQL dump
sqlite3 database/data.db .dump > database/backups/manual_backup.sql
```

### Восстановление вручную

```bash
# Из копии файла
cp database/backups/data_20251228_143022.db database/data.db

# Из SQL dump
sqlite3 database/data.db < database/backups/backup_20251228_143022.sql
```

### Автоматизация резервного копирования

**Cron (Linux/macOS):**

```bash
# Редактировать crontab
crontab -e

# Добавить строку (каждый день в 3:00 ночи)
0 3 * * * cd /path/to/project && ./scripts/backup-database.sh >> logs/backup.log 2>&1
```

**Systemd Timer (Linux):**

```ini
# /etc/systemd/system/db-backup.timer
[Unit]
Description=Database Backup Timer

[Timer]
OnCalendar=daily
OnCalendar=03:00
Persistent=true

[Install]
WantedBy=timers.target
```

---

## 📝 Система логирования

### Структура логов

```
logs/
├── main/
│   ├── main_2025-12-28.log
│   └── main_2025-12-27.log
├── admin/
│   ├── admin_2025-12-28.log
│   └── admin_2025-12-27.log
├── petid/
│   ├── petid_2025-12-28.log
│   └── petid_2025-12-27.log
└── shelter/
    └── shelter_2025-12-28.log
```

### Использование в Go

#### 1. Инициализация логгера

```go
package main

import (
    "database"
    "log"
)

func main() {
    // Создать логгер для сервиса
    logger, err := database.NewLogger("main")
    if err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }
    defer logger.Close()

    logger.Info("Server starting on port 8000")
}
```

#### 2. Использование глобального логгера

```go
package main

import (
    "database"
    "log"
)

func main() {
    // Инициализировать глобальный логгер
    if err := database.InitGlobalLogger("main"); err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }
    defer database.GlobalLogger.Close()

    // Использовать в любом месте
    database.GlobalLogger.Info("Application started")
}
```

#### 3. Уровни логирования

```go
// INFO - информационные сообщения
logger.Info("User %d logged in", userID)

// WARNING - предупреждения
logger.Warning("High memory usage: %d MB", memUsage)

// ERROR - ошибки
logger.Error("Failed to connect to database: %v", err)

// DEBUG - отладочная информация
logger.Debug("Processing request with params: %+v", params)
```

#### 4. Специализированные методы

```go
// Логирование HTTP запроса
logger.LogRequest("GET", "/api/users", "192.168.1.1", 200, time.Millisecond*150)

// Логирование ошибки с контекстом
logger.LogError("Database query failed", err)

// Логирование SQL запроса
logger.LogDBQuery("SELECT * FROM users WHERE id = ?", time.Millisecond*5)
```

#### 5. Middleware для HTTP

```go
package main

import (
    "database"
    "main/backend/middleware"
    "net/http"
)

func main() {
    logger, _ := database.NewLogger("main")
    defer logger.Close()

    // Применить middleware ко всем routes
    logMiddleware := middleware.LoggingMiddleware(logger)

    http.HandleFunc("/api/users", logMiddleware(handleUsers))
    http.HandleFunc("/api/posts", logMiddleware(handlePosts))

    http.ListenAndServe(":8000", nil)
}
```

### Формат логов

```
[2025-12-28 14:30:22] [main] [INFO] Server starting on port 8000
[2025-12-28 14:30:23] [main] [INFO] HTTP GET /api/users from 192.168.1.1 - Status: 200 - Duration: 150ms
[2025-12-28 14:30:25] [main] [ERROR] Database query failed: connection timeout
[2025-12-28 14:30:26] [main] [DEBUG] DB Query: SELECT * FROM users WHERE id = ? - Duration: 5ms
```

### Просмотр логов

```bash
# Последние 50 строк
tail -n 50 logs/main/main_2025-12-28.log

# Следить за логами в реальном времени
tail -f logs/main/main_2025-12-28.log

# Поиск ошибок
grep ERROR logs/main/main_2025-12-28.log

# Поиск по всем логам
grep -r "connection timeout" logs/

# Статистика по уровням
grep -c INFO logs/main/main_2025-12-28.log
grep -c ERROR logs/main/main_2025-12-28.log
```

### Ротация логов

Логи автоматически разделяются по дням (новый файл каждый день).

**Очистка старых логов (вручную):**

```bash
# Удалить логи старше 30 дней
find logs/ -name "*.log" -mtime +30 -delete

# Архивировать старые логи
find logs/ -name "*.log" -mtime +7 -exec gzip {} \;
```

**Автоматическая очистка (cron):**

```bash
# Каждую неделю удалять логи старше 30 дней
0 0 * * 0 find /path/to/project/logs/ -name "*.log" -mtime +30 -delete
```

---

## 🔧 Интеграция в проект

### Main Backend

```go
// main/backend/main.go
package main

import (
    "database"
    "log"
    "main/backend/middleware"
    "net/http"
)

func main() {
    // Инициализация логгера
    logger, err := database.NewLogger("main")
    if err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }
    defer logger.Close()

    logger.Info("Main backend starting...")

    // Middleware для логирования
    logMiddleware := middleware.LoggingMiddleware(logger)

    // Routes с логированием
    http.HandleFunc("/api/users", logMiddleware(handleUsers))
    http.HandleFunc("/api/posts", logMiddleware(handlePosts))

    logger.Info("Server listening on :8000")
    if err := http.ListenAndServe(":8000", nil); err != nil {
        logger.Error("Server failed: %v", err)
    }
}
```

### Admin Backend

```go
// admin/backend/main.go
package main

import (
    "database"
    "log"
)

func main() {
    logger, err := database.NewLogger("admin")
    if err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }
    defer logger.Close()

    logger.Info("Admin backend starting on :9000")
    // ...
}
```

### PetID Backend

```go
// petid/backend/main.go
package main

import (
    "database"
    "log"
)

func main() {
    logger, err := database.NewLogger("petid")
    if err != nil {
        log.Fatal("Failed to initialize logger:", err)
    }
    defer logger.Close()

    logger.Info("PetID backend starting on :8100")
    // ...
}
```

---

## 📊 Мониторинг

### Анализ логов

```bash
# Количество запросов по endpoints
grep "HTTP" logs/main/main_2025-12-28.log | awk '{print $6}' | sort | uniq -c

# Средняя длительность запросов
grep "Duration:" logs/main/main_2025-12-28.log | awk '{print $NF}' | sed 's/ms//' | awk '{sum+=$1; count++} END {print sum/count "ms"}'

# Топ ошибок
grep ERROR logs/main/main_2025-12-28.log | sort | uniq -c | sort -rn | head -10

# Активность по часам
grep "HTTP" logs/main/main_2025-12-28.log | awk '{print $2}' | cut -d: -f1 | sort | uniq -c
```

---

## ✅ Best Practices

### Резервное копирование:
1. ✅ Делайте бэкап **перед каждой миграцией**
2. ✅ Храните минимум **10 последних копий**
3. ✅ Проверяйте **целостность** резервных копий
4. ✅ Храните копии в **отдельном месте** (не только на сервере)
5. ✅ Тестируйте **процесс восстановления** регулярно

### Логирование:
1. ✅ Используйте правильные **уровни логирования**
2. ✅ Логируйте **все HTTP запросы**
3. ✅ Логируйте **ошибки с контекстом**
4. ✅ Не логируйте **чувствительные данные** (пароли, токены)
5. ✅ Ротируйте логи **автоматически**
6. ✅ Мониторьте **размер логов**

---

**Последнее обновление:** 28 декабря 2025  
**Версия:** 1.0.0
