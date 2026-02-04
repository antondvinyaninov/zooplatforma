# Gateway Deployment Guide

## 🚀 Деплой на EasyPanel

### 1. Создать новое приложение в EasyPanel

1. Открыть EasyPanel: http://88.218.121.213:3000
2. Создать новый проект: **"gateway"**
3. Тип: **GitHub Repository**
4. Repository: `antondvinyaninov/zooplatforma`
5. Branch: **`gateway`** ⚠️ ВАЖНО!
6. Build Path: `.` (корень, т.к. в ветке gateway файлы в корне)
7. Dockerfile Path: `Dockerfile` (в корне ветки gateway)

### 2. Настроить переменные окружения

В EasyPanel добавить следующие переменные:

```bash
# JWT Secret (тот же что и в других сервисах!)
JWT_SECRET=jyjy4VlgOPGIPSG5vJPurXDnd1ZpHj2X2dIBtdWfjJE=

# Gateway Port
GATEWAY_PORT=80

# Environment
ENVIRONMENT=production

# PostgreSQL Database (внутри Docker сети)
DATABASE_HOST=zooplatforma-db
DATABASE_PORT=5432
DATABASE_USER=zp
DATABASE_PASSWORD=lmLG7k2ed4vas19
DATABASE_NAME=zp-db

# Backend Services URLs (внутри Docker сети EasyPanel)
# ВАЖНО: MAIN_SERVICE_URL должен указывать на порт 80 (nginx)
MAIN_SERVICE_URL=http://my-projects-zooplatforma:80
PETBASE_SERVICE_URL=http://petbase-backend:8100
CLINIC_SERVICE_URL=http://clinic-backend:8600
OWNER_SERVICE_URL=http://owner-backend:8400
SHELTER_SERVICE_URL=http://shelter-backend:8200
VOLUNTEER_SERVICE_URL=http://volunteer-backend:8500
ADMIN_SERVICE_URL=http://admin-backend:9000

# Uploads Path
UPLOAD_PATH=/app/uploads
```

### 3. Настроить домен

1. В EasyPanel открыть настройки Gateway
2. Добавить домен: `gateway.your-domain.com`
3. Или использовать EasyPanel домен

### 4. Настроить порты

- **Internal Port:** 80 (Gateway слушает на порту 80)
- **External Port:** 80 или 443 (HTTPS)

### 5. Запустить деплой

1. Нажать **"Deploy"**
2. Дождаться сборки (~2-3 минуты)
3. Проверить логи

### 6. Проверить работу

```bash
# Health check
curl https://gateway.your-domain.com/health

# Регистрация
curl -X POST https://gateway.your-domain.com/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'

# Авторизация
curl https://gateway.your-domain.com/api/auth/me \
  -H "Cookie: auth_token=YOUR_TOKEN"
```

## 🔧 Локальная разработка

### Запуск локально

```bash
cd gateway

# Скопировать .env
cp .env.example .env

# Отредактировать .env для локальной разработки
# DATABASE_HOST=88.218.121.213 (удаленный доступ)
# MAIN_SERVICE_URL=http://localhost:8000
# GATEWAY_PORT=8080

# Запустить
go run .
```

### Тестирование

```bash
# Health check
curl http://localhost:8080/health

# Регистрация
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","name":"Test User"}'
```

## 📋 Архитектура

```
User → Gateway (80) → [JWT Auth] → Backend Services
                                  → Main Backend (8000)
                                  → PetBase (8100)
                                  → Clinic (8600)
                                  → etc.
```

## 🔐 Безопасность

- ✅ JWT токены для авторизации
- ✅ Rate limiting (100 req/s per IP)
- ✅ CORS настроен для всех фронтендов
- ✅ PostgreSQL внутри Docker сети
- ✅ Все backend сервисы недоступны извне (только через Gateway)

## 🐛 Troubleshooting

### Gateway не подключается к PostgreSQL

Проверить что `DATABASE_HOST=zooplatforma-db` (имя сервиса в Docker сети)

### Backend сервисы недоступны

Проверить что URL сервисов правильные (имена контейнеров в EasyPanel)

### 401 Unauthorized на всех запросах

Проверить что `JWT_SECRET` одинаковый во всех сервисах

## 📝 Обновление

```bash
# Локально
git add -A
git commit -m "Update Gateway"
git push origin gateway

# EasyPanel автоматически задеплоит изменения
```

## 🔄 Откат

```bash
# В EasyPanel можно откатиться на предыдущий деплой
# Или откатить коммит в ветке gateway
git revert HEAD
git push origin gateway
```
