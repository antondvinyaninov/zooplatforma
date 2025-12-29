# Хранение файлов в ЗооПлатформе

**Версия:** 1.0  
**Дата:** 28 декабря 2025  
**Статус:** Планирование Production

---

## 🎯 Принципы хранения файлов

### ⚠️ КРИТИЧЕСКОЕ ПРАВИЛО:

**НИКОГДА не храните файлы в базе данных!**

### Правильная архитектура:

```
База данных (SQLite/PostgreSQL):
├── Пути к файлам (строки)
├── Метаданные (размер, тип, дата)
└── Связи (user_id, post_id)

Файловая система / S3:
├── Сами файлы (изображения, видео)
└── Оптимизированные версии
```

---

## 📁 Текущая структура (Development)

### Папка uploads/

```
uploads/
├── avatars/              # Аватары пользователей
│   ├── user_1.jpg
│   ├── user_2.png
│   └── ...
├── covers/               # Обложки профилей
│   ├── user_1.jpg
│   └── ...
├── posts/                # Медиа в постах
│   ├── photo_123.jpg
│   ├── video_456.mp4
│   ├── video_456_optimized.mp4  # Оптимизированная версия
│   └── ...
└── pets/                 # Фото питомцев
    ├── pet_1.jpg
    └── ...
```

### База данных

```sql
-- Таблица users
CREATE TABLE users (
    id INTEGER PRIMARY KEY,
    name TEXT,
    avatar TEXT,  -- "/uploads/avatars/user_1.jpg"
    cover TEXT    -- "/uploads/covers/user_1.jpg"
);

-- Таблица posts
CREATE TABLE posts (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    content TEXT,
    media_url TEXT,      -- "/uploads/posts/photo_123.jpg"
    media_type TEXT      -- "image" или "video"
);

-- Таблица user_media (галерея)
CREATE TABLE user_media (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    file_path TEXT,      -- "/uploads/posts/photo_123.jpg"
    file_type TEXT,      -- "image" или "video"
    file_size INTEGER,   -- размер в байтах
    created_at DATETIME
);

-- Таблица pets
CREATE TABLE pets (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    name TEXT,
    photo TEXT           -- "/uploads/pets/pet_1.jpg"
);
```

---

## 🚀 Production архитектура

### Вариант 1: Локальное хранилище + CDN (простой)

#### Структура на сервере:

```
/var/www/zooplatform/
├── backend/
├── frontend/
└── uploads/              # Файлы пользователей
    ├── avatars/
    ├── covers/
    ├── posts/
    └── pets/
```

#### Nginx конфигурация:

```nginx
server {
    listen 80;
    server_name zooplatform.ru;

    # Статические файлы (uploads)
    location /uploads/ {
        alias /var/www/zooplatform/uploads/;
        expires 30d;
        add_header Cache-Control "public, immutable";
        
        # Защита от hotlinking
        valid_referers none blocked zooplatform.ru *.zooplatform.ru;
        if ($invalid_referer) {
            return 403;
        }
    }

    # Backend API
    location /api/ {
        proxy_pass http://localhost:8000;
    }

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
    }
}
```

#### CDN (Cloudflare):

```
Пользователь запрашивает:
https://zooplatform.ru/uploads/avatars/user_1.jpg

Cloudflare кеширует и отдаёт с ближайшего сервера:
https://cdn.zooplatform.ru/uploads/avatars/user_1.jpg
```

#### Плюсы:
- ✅ Простая настройка
- ✅ Полный контроль
- ✅ Низкая стоимость (~$0 за хранение)
- ✅ Быстрая отдача через CDN

#### Минусы:
- ❌ Нужно настраивать backup вручную
- ❌ Ограничено дисковым пространством
- ❌ При масштабировании нужно синхронизировать между серверами

#### Стоимость:
- Хранилище: бесплатно (диск сервера)
- CDN: бесплатно (Cloudflare Free Plan)
- **Итого: $0/месяц**

---

### Вариант 2: S3-совместимое хранилище (рекомендуется)

#### Провайдеры:

**Yandex Object Storage (рекомендуется для РФ):**
- Стоимость: ~1.5₽/GB/месяц (~$0.015/GB)
- Встроенный CDN
- Высокая надёжность (99.999%)
- Простая интеграция

**AWS S3:**
- Стоимость: ~$0.023/GB/месяц
- Глобальный CDN (CloudFront)
- Максимальная надёжность

**MinIO (self-hosted):**
- Стоимость: бесплатно (свой сервер)
- S3-совместимый API
- Полный контроль

#### Структура в S3:

```
Bucket: zooplatform-uploads

s3://zooplatform-uploads/
├── avatars/
│   ├── user_1.jpg
│   └── ...
├── covers/
│   └── ...
├── posts/
│   └── ...
└── pets/
    └── ...
```

#### Код интеграции (Go):

```go
// utils/s3.go
package utils

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

var S3Client *s3.S3

func InitS3() {
    sess := session.Must(session.NewSession(&aws.Config{
        Region:      aws.String("ru-central1"),
        Endpoint:    aws.String("https://storage.yandexcloud.net"),
        Credentials: credentials.NewStaticCredentials(
            os.Getenv("S3_ACCESS_KEY"),
            os.Getenv("S3_SECRET_KEY"),
            "",
        ),
    }))
    
    S3Client = s3.New(sess)
}

func UploadToS3(file multipart.File, filename string, folder string) (string, error) {
    key := fmt.Sprintf("%s/%s", folder, filename)
    
    _, err := S3Client.PutObject(&s3.PutObjectInput{
        Bucket: aws.String("zooplatform-uploads"),
        Key:    aws.String(key),
        Body:   file,
        ACL:    aws.String("public-read"),
    })
    
    if err != nil {
        return "", err
    }
    
    // Возвращаем CDN URL
    url := fmt.Sprintf("https://cdn.zooplatform.ru/%s", key)
    return url, nil
}

func DeleteFromS3(fileURL string) error {
    // Извлекаем key из URL
    key := strings.TrimPrefix(fileURL, "https://cdn.zooplatform.ru/")
    
    _, err := S3Client.DeleteObject(&s3.DeleteObjectInput{
        Bucket: aws.String("zooplatform-uploads"),
        Key:    aws.String(key),
    })
    
    return err
}
```

#### Обновление handlers:

```go
// handlers/avatar.go
func UploadAvatar(c *gin.Context) {
    file, header, err := c.Request.FormFile("avatar")
    if err != nil {
        c.JSON(400, gin.H{"error": "No file uploaded"})
        return
    }
    defer file.Close()
    
    userID := c.GetInt("user_id")
    filename := fmt.Sprintf("user_%d_%d.jpg", userID, time.Now().Unix())
    
    // Загружаем в S3 вместо локального диска
    url, err := utils.UploadToS3(file, filename, "avatars")
    if err != nil {
        c.JSON(500, gin.H{"error": "Failed to upload"})
        return
    }
    
    // Сохраняем URL в базу
    db.Exec("UPDATE users SET avatar = ? WHERE id = ?", url, userID)
    
    c.JSON(200, gin.H{"success": true, "url": url})
}
```

#### Плюсы:
- ✅ Автоматический backup
- ✅ Неограниченное пространство
- ✅ Встроенный CDN
- ✅ Высокая надёжность (99.999%)
- ✅ Масштабируемость
- ✅ Не нужно настраивать Nginx для статики

#### Минусы:
- ❌ Стоимость (но небольшая)
- ❌ Зависимость от внешнего сервиса

#### Стоимость (пример для 100GB):
- Хранилище: 100GB × 1.5₽ = 150₽/месяц (~$1.5)
- Трафик: первые 10TB бесплатно
- **Итого: ~150₽/месяц (~$1.5)**

---

### Вариант 3: Гибридный подход

#### Архитектура:

```
Development:
└── Локальная папка /uploads/

Staging:
└── MinIO (self-hosted S3)

Production:
├── Yandex Object Storage (основное хранилище)
├── CDN (быстрая отдача)
└── Локальный кеш (часто запрашиваемые файлы)
```

#### Код с переключением:

```go
// config/storage.go
type StorageProvider interface {
    Upload(file multipart.File, filename string, folder string) (string, error)
    Delete(fileURL string) error
}

var Storage StorageProvider

func InitStorage() {
    env := os.Getenv("ENVIRONMENT") // "development", "staging", "production"
    
    switch env {
    case "production":
        Storage = &S3Storage{Provider: "yandex"}
    case "staging":
        Storage = &S3Storage{Provider: "minio"}
    default:
        Storage = &LocalStorage{BasePath: "./uploads"}
    }
}
```

---

## 📊 Сравнение вариантов

| Критерий | Локальное + CDN | S3 + CDN | Гибридный |
|----------|----------------|----------|-----------|
| **Стоимость** | $0 | ~$1.5/100GB | ~$1.5/100GB |
| **Надёжность** | Средняя | Высокая (99.999%) | Высокая |
| **Масштабируемость** | Низкая | Высокая | Высокая |
| **Простота настройки** | Простая | Средняя | Сложная |
| **Backup** | Вручную | Автоматический | Автоматический |
| **Скорость отдачи** | Быстрая (CDN) | Быстрая (CDN) | Очень быстрая |

---

## 🔄 План миграции на Production

### Этап 1: Подготовка (1 день)

- [ ] Выбрать провайдера (Yandex Object Storage рекомендуется)
- [ ] Создать bucket `zooplatform-uploads`
- [ ] Настроить CDN
- [ ] Получить access keys

### Этап 2: Код (2 дня)

- [ ] Создать абстракцию StorageProvider
- [ ] Реализовать S3Storage
- [ ] Обновить все handlers (avatar, cover, posts, pets)
- [ ] Добавить переменные окружения
- [ ] Тестирование на staging

### Этап 3: Миграция данных (1 день)

```bash
# Скрипт миграции существующих файлов
#!/bin/bash

# Загружаем все файлы из uploads/ в S3
aws s3 sync ./uploads/ s3://zooplatform-uploads/ \
    --endpoint-url https://storage.yandexcloud.net \
    --acl public-read

# Обновляем пути в базе данных
sqlite3 database/data.db <<EOF
UPDATE users SET avatar = REPLACE(avatar, '/uploads/', 'https://cdn.zooplatform.ru/');
UPDATE users SET cover = REPLACE(cover, '/uploads/', 'https://cdn.zooplatform.ru/');
UPDATE posts SET media_url = REPLACE(media_url, '/uploads/', 'https://cdn.zooplatform.ru/');
UPDATE user_media SET file_path = REPLACE(file_path, '/uploads/', 'https://cdn.zooplatform.ru/');
UPDATE pets SET photo = REPLACE(photo, '/uploads/', 'https://cdn.zooplatform.ru/');
EOF
```

### Этап 4: Запуск (1 день)

- [ ] Деплой на production
- [ ] Проверка загрузки новых файлов
- [ ] Проверка отображения старых файлов
- [ ] Мониторинг ошибок

**Итого: 5 дней**

---

## 🔒 Безопасность

### Защита от несанкционированного доступа:

```go
// Проверка прав доступа перед загрузкой
func UploadAvatar(c *gin.Context) {
    userID := c.GetInt("user_id")
    
    // Только владелец может загружать свой аватар
    if userID != targetUserID {
        c.JSON(403, gin.H{"error": "Forbidden"})
        return
    }
    
    // ... загрузка файла
}
```

### Валидация файлов:

```go
func ValidateImage(file multipart.File) error {
    // Проверка типа файла
    buffer := make([]byte, 512)
    file.Read(buffer)
    file.Seek(0, 0)
    
    contentType := http.DetectContentType(buffer)
    if !strings.HasPrefix(contentType, "image/") {
        return errors.New("not an image")
    }
    
    // Проверка размера
    if file.Size() > 10*1024*1024 { // 10MB
        return errors.New("file too large")
    }
    
    return nil
}
```

### Защита от hotlinking (Nginx):

```nginx
location /uploads/ {
    valid_referers none blocked zooplatform.ru *.zooplatform.ru;
    if ($invalid_referer) {
        return 403;
    }
}
```

---

## 📈 Оптимизация

### 1. Сжатие изображений

```go
import "github.com/disintegration/imaging"

func CompressImage(src image.Image) (image.Image, error) {
    // Уменьшаем до максимум 1920x1080
    bounds := src.Bounds()
    if bounds.Dx() > 1920 || bounds.Dy() > 1080 {
        src = imaging.Fit(src, 1920, 1080, imaging.Lanczos)
    }
    
    return src, nil
}
```

### 2. Генерация thumbnails

```go
func GenerateThumbnail(src image.Image) image.Image {
    return imaging.Thumbnail(src, 300, 300, imaging.Lanczos)
}

// Сохраняем две версии:
// - /avatars/user_1.jpg (оригинал)
// - /avatars/user_1_thumb.jpg (thumbnail)
```

### 3. Lazy loading (Frontend)

```tsx
<img 
    src={user.avatar} 
    loading="lazy"
    alt={user.name}
/>
```

### 4. WebP формат

```go
import "github.com/chai2010/webp"

func ConvertToWebP(src image.Image) ([]byte, error) {
    var buf bytes.Buffer
    err := webp.Encode(&buf, src, &webp.Options{Quality: 80})
    return buf.Bytes(), err
}

// Сохраняем обе версии:
// - /avatars/user_1.jpg (для старых браузеров)
// - /avatars/user_1.webp (для современных)
```

---

## 📊 Мониторинг

### Метрики для отслеживания:

- Размер хранилища (GB)
- Количество файлов
- Трафик (GB/месяц)
- Стоимость
- Скорость загрузки
- Ошибки загрузки

### Дашборд в Admin панели:

```
┌─────────────────────────────────────┐
│  Хранилище файлов                   │
├─────────────────────────────────────┤
│  Всего файлов: 12,345               │
│  Размер: 45.6 GB                    │
│  Трафик (месяц): 234 GB             │
│  Стоимость: 150₽/месяц              │
│                                     │
│  По типам:                          │
│  - Аватары: 1,234 (2.3 GB)         │
│  - Обложки: 567 (1.8 GB)           │
│  - Посты: 8,901 (38.5 GB)          │
│  - Питомцы: 1,643 (3.0 GB)         │
└─────────────────────────────────────┘
```

---

## ✅ Checklist для Production

### Перед запуском:
- [ ] Выбран провайдер хранилища
- [ ] Настроен CDN
- [ ] Реализована абстракция StorageProvider
- [ ] Обновлены все handlers
- [ ] Настроена валидация файлов
- [ ] Настроена защита от hotlinking
- [ ] Реализовано сжатие изображений
- [ ] Созданы thumbnails
- [ ] Протестирована загрузка
- [ ] Протестировано удаление
- [ ] Мигрированы существующие файлы
- [ ] Обновлены пути в БД
- [ ] Настроен мониторинг
- [ ] Настроен backup

### После запуска:
- [ ] Проверена загрузка новых файлов
- [ ] Проверено отображение старых файлов
- [ ] Проверена скорость отдачи
- [ ] Мониторинг ошибок
- [ ] Проверка стоимости

---

## 🔮 Будущие улучшения

### Версия 1.1.0+:
- [ ] Автоматическая генерация WebP
- [ ] Адаптивные изображения (srcset)
- [ ] Видео транскодинг в облаке
- [ ] Автоматическое удаление неиспользуемых файлов
- [ ] Дедупликация файлов (одинаковые файлы)
- [ ] Watermark для изображений
- [ ] EXIF данные (геолокация, дата съёмки)

---

**Документ создан:** 28 декабря 2025  
**Автор:** ЗооПлатформа Team  
**Версия:** 1.0
