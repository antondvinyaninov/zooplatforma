# API Gateway

Централизованный API Gateway для всех микросервисов проекта ZooPlatforma.

## 🎯 Назначение

Gateway выполняет следующие функции:
- **Роутинг** - направляет запросы к нужным backend сервисам
- **Авторизация** - проверяет JWT токены централизованно
- **Rate Limiting** - защита от DDoS атак
- **Логирование** - все запросы логируются в одном месте
- **CORS** - управление cross-origin запросами
- **Health Checks** - мониторинг здоровья всех сервисов

## 🏗️ Архитектура

```
Frontend → Gateway (80) 
                        → Main Backend (8000)
                        → PetBase Backend (8100)
                        → Clinic Backend (8600)
                        → Owner Backend (8400)
                        → Shelter Backend (8200)
                        → Volunteer Backend (8500)
                        → Admin Backend (9000)
```

## 📋 Роутинг

| Path | Service | Auth Required | Notes |
|------|---------|---------------|-------|
| `/api/auth/*` | Auth Service | ❌ | Регистрация, логин |
| `/api/posts` (GET) | Main Backend | ❌ | Публичный просмотр |
| `/api/posts` (POST/PUT/DELETE) | Main Backend | ✅ | Требует авторизации |
| `/api/users/{id}` (GET) | Main Backend | ❌ | Публичный профиль |
| `/api/profile` | Main Backend | ✅ | Редактирование профиля |
| `/api/petbase/*` | PetBase Backend | ❌ | Справочник животных |
| `/api/clinic/*` | Clinic Backend | ✅ | Кабинет клиники |
| `/api/owner/*` | Owner Backend | ✅ | Кабинет владельца |
| `/api/shelter/*` | Shelter Backend | ✅ | Кабинет приюта |
| `/api/volunteer/*` | Volunteer Backend | ✅ | Кабинет волонтера |
| `/api/admin/*` | Admin Backend | ✅ + Admin | Админ панель |
| `/uploads/*` | Static Files | ❌ | Загруженные файлы |

## 🔐 Авторизация

Gateway проверяет JWT токен и добавляет заголовки для backend сервисов:

```
X-User-ID: 123
X-User-Email: user@example.com
X-User-Role: user
```

Backend сервисы **НЕ проверяют** JWT - они доверяют Gateway и используют заголовки.

## 🚀 Запуск

### Development

```bash
cd gateway
go run .
```

### Production

```bash
cd gateway
go build -o gateway .
./gateway
```

## ⚙️ Переменные окружения

```bash
# JWT Secret (обязательно!)
JWT_SECRET=your-secret-key

# Gateway
GATEWAY_PORT=80

# Backend Services
AUTH_SERVICE_URL=http://localhost:7100
MAIN_SERVICE_URL=http://localhost:8000
PETBASE_SERVICE_URL=http://localhost:8100
CLINIC_SERVICE_URL=http://localhost:8600
OWNER_SERVICE_URL=http://localhost:8400
SHELTER_SERVICE_URL=http://localhost:8200
VOLUNTEER_SERVICE_URL=http://localhost:8500
ADMIN_SERVICE_URL=http://localhost:9000

# Uploads
UPLOAD_PATH=/app/uploads
```

## 📊 Мониторинг

### Health Check

```bash
curl http://localhost/health
```

Ответ:
```json
{
  "success": true,
  "status": "healthy",
  "gateway": "API Gateway",
  "version": "1.0.0",
  "services": {
    "auth_service": {
      "url": "http://localhost:7100",
      "healthy": true
    },
    "main_backend": {
      "url": "http://localhost:8000",
      "healthy": true
    }
    // ... другие сервисы
  }
}
```

## 🔧 Настройка Rate Limiting

По умолчанию:
- **100 запросов/секунду** с одного IP
- **Burst до 200 запросов**

Изменить в `middleware.go`:
```go
rate:  rate.Limit(100), // запросов в секунду
burst: 200,             // максимальный burst
```

## 📝 Логирование

Все запросы логируются в формате:
```
📋 GET /api/posts 200 15ms 127.0.0.1
✅ Proxied to Main Backend: GET /api/posts → 200 (took 15ms)
✅ Authenticated: user_id=1, email=user@example.com, role=user
```

## 🐛 Отладка

### Проверить что Gateway запущен
```bash
curl http://localhost/health
```

### Проверить авторизацию
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost/api/profile
```

### Проверить проксирование
```bash
curl http://localhost/api/posts
```

## 🔄 Миграция с nginx

### Что меняется:

**Было (nginx):**
```
Frontend → nginx (80) → Main Backend (8000) → Auth Service (7100)
```

**Стало (Gateway):**
```
Frontend → Gateway (80) → Auth Service (7100)
                        → Main Backend (8000)
```

### Что удалить:

1. ❌ `nginx.conf` - больше не нужен
2. ❌ `pkg/middleware/auth.go` - логика в Gateway
3. ❌ Импорты middleware из backend сервисов

### Что обновить:

1. ✅ Backend сервисы - убрать проверку JWT, использовать заголовки `X-User-ID`
2. ✅ Dockerfile - собирать Gateway вместо nginx
3. ✅ docker-compose.yml - запускать Gateway

## 📚 Структура файлов

```
gateway/
├── main.go          # Точка входа, роутинг
├── auth.go          # JWT авторизация
├── middleware.go    # CORS, rate limiting, логирование
├── proxy.go         # Проксирование запросов
├── services.go      # Конфигурация сервисов
├── go.mod           # Зависимости
├── .env.example     # Пример переменных окружения
└── README.md        # Эта документация
```

## 🎯 Преимущества Gateway

1. ✅ **Централизованная авторизация** - JWT проверяется в одном месте
2. ✅ **Меньше кода** - backend сервисы не проверяют JWT
3. ✅ **Rate limiting** - защита от DDoS
4. ✅ **Единая точка входа** - проще мониторить
5. ✅ **Логирование** - все запросы в одном месте
6. ✅ **CORS** - управление в одном месте
7. ✅ **Health checks** - мониторинг всех сервисов

## 🚨 Важно

- Gateway **НЕ заменяет** Auth Service - это разные вещи!
- Auth Service - бизнес-логика авторизации (регистрация, логин, JWT)
- Gateway - роутер и проверка JWT

## 📞 Поддержка

При проблемах проверьте:
1. Все backend сервисы запущены
2. JWT_SECRET установлен
3. Порты не заняты
4. Логи Gateway: `docker logs gateway`
