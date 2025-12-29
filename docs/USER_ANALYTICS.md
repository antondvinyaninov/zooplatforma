# Аналитика активности пользователей

**Версия:** 1.0  
**Дата:** 28 декабря 2025  
**Статус:** Критично для MVP 1.0.0

---

## 🎯 Зачем нужна аналитика активности?

Аналитика активности пользователей - это **фундамент** для:

1. **NPS опросов** - определение когда показывать опрос
2. **Retention анализа** - кто возвращается, кто уходит
3. **Engagement метрик** - насколько пользователи вовлечены
4. **Персонализации** - что показывать каждому пользователю
5. **Выявления проблем** - где пользователи "застревают"

---

## 📊 Что отслеживаем

### 1. Базовые метрики пользователя

**Хранятся в таблице `users`:**

```sql
ALTER TABLE users ADD COLUMN registered_at DATETIME DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE users ADD COLUMN last_login_at DATETIME;
ALTER TABLE users ADD COLUMN login_count INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN total_sessions INTEGER DEFAULT 0;
```

**Что отслеживаем:**
- Дата регистрации (`registered_at`)
- Последний вход (`last_login_at`)
- Количество входов (`login_count`)
- Общее количество сессий (`total_sessions`)

### 2. Сессии пользователя

**Таблица `user_sessions`:**

```sql
CREATE TABLE user_sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    ended_at DATETIME,
    duration_seconds INTEGER, -- вычисляется при завершении сессии
    pages_viewed INTEGER DEFAULT 0,
    actions_count INTEGER DEFAULT 0,
    device_type TEXT, -- 'desktop', 'mobile', 'tablet'
    browser TEXT,
    ip_address TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_started_at ON user_sessions(started_at);
```

**Что отслеживаем:**
- Начало и конец сессии
- Длительность сессии
- Количество просмотренных страниц
- Количество действий
- Тип устройства и браузер

### 3. События активности

**Таблица `user_activity_log`:**

```sql
CREATE TABLE user_activity_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    session_id INTEGER,
    action_type TEXT NOT NULL, -- 'post_create', 'comment_add', 'like', 'profile_view', etc.
    entity_type TEXT, -- 'post', 'user', 'pet', 'comment'
    entity_id INTEGER,
    metadata TEXT, -- JSON с дополнительными данными
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (session_id) REFERENCES user_sessions(id)
);

CREATE INDEX idx_activity_user_id ON user_activity_log(user_id);
CREATE INDEX idx_activity_action_type ON user_activity_log(action_type);
CREATE INDEX idx_activity_created_at ON user_activity_log(created_at);
```

**Типы действий (action_type):**

**Контент:**
- `post_create` - создание поста
- `post_edit` - редактирование поста
- `post_delete` - удаление поста
- `post_view` - просмотр поста
- `comment_add` - добавление комментария
- `comment_delete` - удаление комментария
- `like_add` - лайк
- `like_remove` - убрать лайк

**Профиль:**
- `profile_view` - просмотр профиля
- `profile_edit` - редактирование профиля
- `avatar_upload` - загрузка аватара
- `cover_upload` - загрузка обложки

**Социальные связи:**
- `friend_request_send` - отправка запроса в друзья
- `friend_request_accept` - принятие запроса
- `friend_request_reject` - отклонение запроса
- `subscription_add` - подписка
- `subscription_remove` - отписка

**Питомцы:**
- `pet_add` - добавление питомца
- `pet_edit` - редактирование питомца
- `pet_delete` - удаление питомца
- `pet_view` - просмотр карточки питомца

**Коммуникации:**
- `message_send` - отправка сообщения (БЕЗ содержимого!)
- `message_read` - прочтение сообщения
- `chat_open` - открытие чата

**⚠️ ВАЖНО:** Содержимое сообщений НЕ логируется в user_activity_log!
Мы отслеживаем ТОЛЬКО факт отправки, но НЕ текст сообщения.

**Поиск:**
- `search_query` - поисковый запрос
- `search_result_click` - клик по результату

**Авторизация:**
- `login` - вход
- `logout` - выход
- `register` - регистрация

### 4. Агрегированная статистика

**Таблица `user_stats` (обновляется ежедневно):**

```sql
CREATE TABLE user_stats (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    date DATE NOT NULL,
    sessions_count INTEGER DEFAULT 0,
    total_time_seconds INTEGER DEFAULT 0,
    posts_created INTEGER DEFAULT 0,
    comments_added INTEGER DEFAULT 0,
    likes_given INTEGER DEFAULT 0,
    profile_views INTEGER DEFAULT 0,
    messages_sent INTEGER DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE(user_id, date)
);

CREATE INDEX idx_user_stats_user_id ON user_stats(user_id);
CREATE INDEX idx_user_stats_date ON user_stats(date);
```

---

## 🔧 Реализация

### Backend (Go)

#### 1. Middleware для отслеживания сессий

```go
// middleware/analytics.go
package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
)

func AnalyticsMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("user_id")
        if !exists {
            c.Next()
            return
        }
        
        // Получаем или создаём сессию
        sessionID := getOrCreateSession(userID.(int), c)
        c.Set("session_id", sessionID)
        
        // Обновляем last_login_at
        updateLastLogin(userID.(int))
        
        c.Next()
        
        // После обработки запроса - логируем активность
        logPageView(userID.(int), sessionID, c.Request.URL.Path)
    }
}

func getOrCreateSession(userID int, c *gin.Context) int {
    // Проверяем, есть ли активная сессия (последние 30 минут)
    var sessionID int
    err := db.QueryRow(`
        SELECT id FROM user_sessions 
        WHERE user_id = ? 
        AND ended_at IS NULL 
        AND started_at > datetime('now', '-30 minutes')
        ORDER BY started_at DESC 
        LIMIT 1
    `, userID).Scan(&sessionID)
    
    if err != nil {
        // Создаём новую сессию
        result, _ := db.Exec(`
            INSERT INTO user_sessions (user_id, device_type, browser, ip_address)
            VALUES (?, ?, ?, ?)
        `, userID, getDeviceType(c), getBrowser(c), c.ClientIP())
        
        id, _ := result.LastInsertId()
        sessionID = int(id)
        
        // Увеличиваем счётчик сессий
        db.Exec("UPDATE users SET total_sessions = total_sessions + 1 WHERE id = ?", userID)
    }
    
    return sessionID
}

func updateLastLogin(userID int) {
    db.Exec(`
        UPDATE users 
        SET last_login_at = CURRENT_TIMESTAMP,
            login_count = login_count + 1
        WHERE id = ?
    `, userID)
}
```

#### 2. Функция логирования активности

```go
// utils/analytics.go
package utils

func LogActivity(userID int, sessionID int, actionType string, entityType string, entityID int, metadata string) {
    db.Exec(`
        INSERT INTO user_activity_log (user_id, session_id, action_type, entity_type, entity_id, metadata)
        VALUES (?, ?, ?, ?, ?, ?)
    `, userID, sessionID, actionType, entityType, entityID, metadata)
    
    // Увеличиваем счётчик действий в сессии
    db.Exec(`
        UPDATE user_sessions 
        SET actions_count = actions_count + 1 
        WHERE id = ?
    `, sessionID)
}
```

#### 3. Использование в handlers

```go
// handlers/posts.go
func CreatePost(c *gin.Context) {
    userID := c.GetInt("user_id")
    sessionID := c.GetInt("session_id")
    
    // ... создание поста ...
    
    // Логируем активность
    utils.LogActivity(userID, sessionID, "post_create", "post", postID, "")
    
    c.JSON(200, gin.H{"success": true, "post_id": postID})
}

func LikePost(c *gin.Context) {
    userID := c.GetInt("user_id")
    sessionID := c.GetInt("session_id")
    postID := c.Param("id")
    
    // ... добавление лайка ...
    
    // Логируем активность
    utils.LogActivity(userID, sessionID, "like_add", "post", postID, "")
    
    c.JSON(200, gin.H{"success": true})
}
```

#### 4. Завершение сессии (cron job)

```go
// jobs/close_sessions.go
func CloseInactiveSessions() {
    // Закрываем сессии, неактивные более 30 минут
    db.Exec(`
        UPDATE user_sessions 
        SET ended_at = datetime(started_at, '+30 minutes'),
            duration_seconds = CAST((julianday(datetime(started_at, '+30 minutes')) - julianday(started_at)) * 86400 AS INTEGER)
        WHERE ended_at IS NULL 
        AND started_at < datetime('now', '-30 minutes')
    `)
}

// Запускать каждые 5 минут
```

#### 5. Агрегация статистики (cron job)

```go
// jobs/aggregate_stats.go
func AggregateUserStats() {
    // Агрегируем статистику за вчерашний день
    db.Exec(`
        INSERT OR REPLACE INTO user_stats (
            user_id, date, sessions_count, total_time_seconds,
            posts_created, comments_added, likes_given, messages_sent
        )
        SELECT 
            user_id,
            DATE(created_at) as date,
            COUNT(DISTINCT session_id) as sessions_count,
            0 as total_time_seconds, -- будет обновлено отдельно
            SUM(CASE WHEN action_type = 'post_create' THEN 1 ELSE 0 END) as posts_created,
            SUM(CASE WHEN action_type = 'comment_add' THEN 1 ELSE 0 END) as comments_added,
            SUM(CASE WHEN action_type = 'like_add' THEN 1 ELSE 0 END) as likes_given,
            SUM(CASE WHEN action_type = 'message_send' THEN 1 ELSE 0 END) as messages_sent
        FROM user_activity_log
        WHERE DATE(created_at) = DATE('now', '-1 day')
        GROUP BY user_id, DATE(created_at)
    `)
    
    // Обновляем total_time_seconds из сессий
    db.Exec(`
        UPDATE user_stats
        SET total_time_seconds = (
            SELECT COALESCE(SUM(duration_seconds), 0)
            FROM user_sessions
            WHERE user_sessions.user_id = user_stats.user_id
            AND DATE(user_sessions.started_at) = user_stats.date
        )
        WHERE date = DATE('now', '-1 day')
    `)
}

// Запускать каждый день в 00:05
```

---

## 📈 Использование для NPS

### Проверка триггеров

```go
// handlers/nps.go

func shouldShowDay7(userID int) bool {
    var registeredAt time.Time
    var activityCount int
    
    // Проверяем дату регистрации и активность
    db.QueryRow(`
        SELECT 
            registered_at,
            (SELECT COUNT(*) FROM user_activity_log WHERE user_id = ? AND created_at > datetime('now', '-7 days')) as activity_count
        FROM users 
        WHERE id = ?
    `, userID, userID).Scan(&registeredAt, &activityCount)
    
    // Зарегистрирован 7 дней назад и был активен минимум 3 раза
    daysSinceRegistration := time.Since(registeredAt).Hours() / 24
    return daysSinceRegistration >= 7 && daysSinceRegistration < 8 && activityCount >= 3
}

func shouldShowDay30(userID int) bool {
    var registeredAt time.Time
    var activityCount int
    var lastNPSDate time.Time
    
    db.QueryRow(`
        SELECT 
            u.registered_at,
            (SELECT COUNT(*) FROM user_activity_log WHERE user_id = ? AND created_at > datetime('now', '-30 days')) as activity_count,
            COALESCE((SELECT MAX(created_at) FROM nps_surveys WHERE user_id = ?), '2000-01-01') as last_nps_date
        FROM users u
        WHERE u.id = ?
    `, userID, userID, userID).Scan(&registeredAt, &activityCount, &lastNPSDate)
    
    daysSinceRegistration := time.Since(registeredAt).Hours() / 24
    daysSinceLastNPS := time.Since(lastNPSDate).Hours() / 24
    
    // Зарегистрирован 30 дней назад, активен минимум 10 раз, прошло >21 день с последнего NPS
    return daysSinceRegistration >= 30 && daysSinceRegistration < 31 && 
           activityCount >= 10 && daysSinceLastNPS > 21
}

func shouldShowQuarterly(userID int) bool {
    var registeredAt time.Time
    var activityCount int
    var lastNPSDate time.Time
    
    db.QueryRow(`
        SELECT 
            u.registered_at,
            (SELECT COUNT(*) FROM user_activity_log WHERE user_id = ? AND created_at > datetime('now', '-90 days')) as activity_count,
            COALESCE((SELECT MAX(created_at) FROM nps_surveys WHERE user_id = ?), '2000-01-01') as last_nps_date
        FROM users u
        WHERE u.id = ?
    `, userID, userID, userID).Scan(&registeredAt, &activityCount, &lastNPSDate)
    
    daysSinceRegistration := time.Since(registeredAt).Hours() / 24
    daysSinceLastNPS := time.Since(lastNPSDate).Hours() / 24
    
    // Зарегистрирован >90 дней, активен минимум 20 раз за 90 дней, прошло >90 дней с последнего NPS
    return daysSinceRegistration > 90 && activityCount >= 20 && daysSinceLastNPS > 90
}
```

---

## 📊 API для аналитики

### GET /api/analytics/user/:id

Получить статистику пользователя.

**Response:**
```json
{
  "user_id": 123,
  "registered_at": "2025-12-01T10:00:00Z",
  "days_since_registration": 27,
  "last_login_at": "2025-12-28T15:30:00Z",
  "total_sessions": 45,
  "total_login_count": 52,
  "last_30_days": {
    "sessions": 15,
    "total_time_seconds": 18000,
    "avg_session_duration": 1200,
    "posts_created": 8,
    "comments_added": 25,
    "likes_given": 120,
    "messages_sent": 45
  },
  "activity_level": "high" // low, medium, high, very_high
}
```

### GET /api/analytics/retention

Retention анализ (только для админов).

**Response:**
```json
{
  "day_1": 0.85,
  "day_7": 0.65,
  "day_30": 0.45,
  "day_90": 0.30
}
```

### GET /api/analytics/engagement

Engagement метрики (только для админов).

**Response:**
```json
{
  "dau": 150,
  "wau": 450,
  "mau": 1200,
  "dau_mau_ratio": 0.125,
  "avg_session_duration": 1350,
  "avg_sessions_per_user": 3.2
}
```

---

## 🎯 Метрики активности

### Уровни активности пользователя

```go
func GetActivityLevel(userID int) string {
    var activityCount int
    db.QueryRow(`
        SELECT COUNT(*) 
        FROM user_activity_log 
        WHERE user_id = ? 
        AND created_at > datetime('now', '-30 days')
    `, userID).Scan(&activityCount)
    
    if activityCount >= 100 {
        return "very_high"
    } else if activityCount >= 50 {
        return "high"
    } else if activityCount >= 20 {
        return "medium"
    } else {
        return "low"
    }
}
```

### Retention Rate

```go
func CalculateRetentionRate(days int) float64 {
    var total, returned int
    
    // Пользователи, зарегистрированные N дней назад
    db.QueryRow(`
        SELECT COUNT(*) 
        FROM users 
        WHERE DATE(registered_at) = DATE('now', '-' || ? || ' days')
    `, days).Scan(&total)
    
    // Из них вернулись хотя бы раз
    db.QueryRow(`
        SELECT COUNT(DISTINCT u.id)
        FROM users u
        JOIN user_activity_log ual ON u.id = ual.user_id
        WHERE DATE(u.registered_at) = DATE('now', '-' || ? || ' days')
        AND DATE(ual.created_at) > DATE(u.registered_at)
    `, days).Scan(&returned)
    
    if total == 0 {
        return 0
    }
    
    return float64(returned) / float64(total)
}
```

### DAU / MAU Ratio

```go
func CalculateDAUMAURatio() float64 {
    var dau, mau int
    
    // Daily Active Users
    db.QueryRow(`
        SELECT COUNT(DISTINCT user_id)
        FROM user_activity_log
        WHERE DATE(created_at) = DATE('now')
    `).Scan(&dau)
    
    // Monthly Active Users
    db.QueryRow(`
        SELECT COUNT(DISTINCT user_id)
        FROM user_activity_log
        WHERE created_at > datetime('now', '-30 days')
    `).Scan(&mau)
    
    if mau == 0 {
        return 0
    }
    
    return float64(dau) / float64(mau)
}
```

---

## 🔒 Приватность и GDPR

### ⚠️ КРИТИЧЕСКОЕ ПРАВИЛО: Приватность сообщений

**Содержимое личных сообщений - это приватная информация!**

#### Что МОЖНО логировать:
- ✅ Факт отправки сообщения (action_type = 'message_send')
- ✅ Кому отправлено (entity_id = to_user_id)
- ✅ Время отправки (created_at)

#### Что НЕЛЬЗЯ логировать:
- ❌ Содержимое сообщения (текст)
- ❌ Анализ текста (sentiment, ключевые слова)
- ❌ Любая информация из содержимого

#### Правильная реализация:

```go
// ✅ ПРАВИЛЬНО
func SendMessage(c *gin.Context) {
    userID := c.GetInt("user_id")
    sessionID := c.GetInt("session_id")
    
    // Сохраняем сообщение в БД
    db.Exec("INSERT INTO messages (from_user_id, to_user_id, content) VALUES (?, ?, ?)", 
        userID, toUserID, content)
    
    // Логируем ТОЛЬКО факт отправки (БЕЗ содержимого)
    utils.LogActivity(userID, sessionID, "message_send", "user", toUserID, "")
    //                                                                        ^^^ пустая metadata
}

// ❌ НЕПРАВИЛЬНО
func SendMessage(c *gin.Context) {
    // НЕ ДЕЛАЙТЕ ТАК!
    utils.LogActivity(userID, sessionID, "message_send", "user", toUserID, content)
    //                                                                       ^^^^^^^ содержимое сообщения!
}
```

### Модерация сообщений

**Когда можно читать сообщения:**

1. **При жалобе пользователя** - модератор видит ТОЛЬКО конкретное сообщение
2. **Автоматическое сканирование** - только для критичных случаев (детская порнография, терроризм)
3. **По запросу правоохранительных органов** - с официальным запросом

**Всё это должно быть указано в Политике конфиденциальности!**

### Анонимизация данных

```sql
-- Удаление персональных данных при запросе пользователя
UPDATE user_sessions 
SET ip_address = NULL, browser = NULL 
WHERE user_id = ?;

-- Или полное удаление аналитики
DELETE FROM user_activity_log WHERE user_id = ?;
DELETE FROM user_sessions WHERE user_id = ?;
DELETE FROM user_stats WHERE user_id = ?;

-- Удаление сообщений
DELETE FROM messages WHERE from_user_id = ? OR to_user_id = ?;
```

### Согласие на сбор аналитики

```sql
ALTER TABLE users ADD COLUMN analytics_consent BOOLEAN DEFAULT TRUE;
ALTER TABLE users ADD COLUMN privacy_policy_accepted BOOLEAN DEFAULT FALSE;
ALTER TABLE users ADD COLUMN privacy_policy_accepted_at DATETIME;
```

В middleware проверять:
```go
if !user.AnalyticsConsent {
    c.Next()
    return // не логируем активность
}
```

### Политика конфиденциальности

**Обязательно создать файл `/privacy` с информацией:**
- Какие данные собираем
- Как используем
- Когда читаем сообщения (модерация)
- Права пользователя (удаление, экспорт)
- Срок хранения данных

---

## ✅ Checklist реализации

### Backend:
- [ ] Добавить поля в таблицу `users` (registered_at, last_login_at, login_count, total_sessions)
- [ ] Создать таблицу `user_sessions`
- [ ] Создать таблицу `user_activity_log`
- [ ] Создать таблицу `user_stats`
- [ ] Реализовать middleware для отслеживания сессий
- [ ] Реализовать функцию `LogActivity()`
- [ ] Добавить логирование во все handlers
- [ ] Реализовать cron job для закрытия сессий
- [ ] Реализовать cron job для агрегации статистики
- [ ] Реализовать API endpoints для аналитики

### NPS интеграция:
- [ ] Обновить функции проверки триггеров (shouldShowDay7, shouldShowDay30, shouldShowQuarterly)
- [ ] Использовать данные активности для определения показа NPS

### Admin панель:
- [ ] Страница аналитики пользователей
- [ ] Графики активности
- [ ] Retention анализ
- [ ] Engagement метрики

### Приватность:
- [ ] Добавить согласие на сбор аналитики
- [ ] Реализовать анонимизацию данных
- [ ] Добавить в политику конфиденциальности

---

## 🚀 План внедрения

### Этап 1: База данных (1 день)
- Создать миграции для всех таблиц
- Обновить существующие данные (registered_at из created_at)

### Этап 2: Backend (2 дня)
- Middleware для сессий
- Функция LogActivity
- Интеграция во все handlers

### Этап 3: Cron jobs (1 день)
- Закрытие сессий
- Агрегация статистики

### Этап 4: API и Admin (1 день)
- API endpoints
- Базовые графики в Admin

**Итого:** 5 дней

---

## 📊 Примеры использования

### Определение "активного" пользователя

```go
func IsActiveUser(userID int) bool {
    var activityCount int
    db.QueryRow(`
        SELECT COUNT(*) 
        FROM user_activity_log 
        WHERE user_id = ? 
        AND created_at > datetime('now', '-7 days')
    `, userID).Scan(&activityCount)
    
    return activityCount >= 5 // минимум 5 действий за неделю
}
```

### Поиск "спящих" пользователей

```sql
-- Пользователи, не заходившие >30 дней
SELECT id, name, email, last_login_at
FROM users
WHERE last_login_at < datetime('now', '-30 days')
ORDER BY last_login_at DESC;
```

### Самые активные пользователи

```sql
-- Топ-10 по активности за месяц
SELECT 
    u.id,
    u.name,
    COUNT(*) as actions_count
FROM user_activity_log ual
JOIN users u ON ual.user_id = u.id
WHERE ual.created_at > datetime('now', '-30 days')
GROUP BY u.id
ORDER BY actions_count DESC
LIMIT 10;
```

---

**Документ создан:** 28 декабря 2025  
**Автор:** ЗооПлатформа Team  
**Версия:** 1.0
