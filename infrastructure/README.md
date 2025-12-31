# Infrastructure

Конфигурации для деплоя и инфраструктуры проекта.

## 📁 Структура

```
infrastructure/
├── docker/                    # Docker конфигурации
│   ├── main/                  # Dockerfiles для main сервиса
│   │   ├── Dockerfile.backend
│   │   └── Dockerfile.frontend
│   ├── admin/                 # Dockerfiles для admin сервиса
│   │   ├── Dockerfile.backend
│   │   └── Dockerfile.frontend
│   ├── petbase/               # Dockerfiles для petbase сервиса
│   │   ├── Dockerfile.backend
│   │   └── Dockerfile.frontend
│   └── shelter/               # Dockerfiles для shelter сервиса
│       ├── Dockerfile.backend
│       └── Dockerfile.frontend
├── kubernetes/                # Kubernetes манифесты (TODO)
├── docker-compose.yml         # Docker Compose для всех сервисов
├── .env.example               # Пример переменных окружения
└── README.md                  # Этот файл
```

## 🐳 Docker

### Быстрый старт

```bash
# 1. Создать .env файл
cp infrastructure/.env.example infrastructure/.env

# 2. Отредактировать .env (установить JWT_SECRET)
nano infrastructure/.env

# 3. Запустить все сервисы
cd infrastructure
docker-compose up -d

# 4. Проверить статус
docker-compose ps

# 5. Просмотр логов
docker-compose logs -f

# 6. Остановка
docker-compose down
```

### Сервисы

| Сервис | Порт | Описание |
|--------|------|----------|
| main-backend | 8000 | Основной backend (Go) |
| main-frontend | 3000 | Основной frontend (Next.js) |
| admin-backend | 9000 | Admin backend (Go) |
| admin-frontend | 4000 | Admin frontend (Next.js) |
| petbase-backend | 8100 | PetBase backend (Go) |
| petbase-frontend | 4100 | PetBase frontend (Next.js) |
| shelter-backend | 8200 | Shelter backend (Go) |
| shelter-frontend | 5100 | Shelter frontend (Next.js) |

### URL сервисов

После запуска `docker-compose up -d` будут доступны:

- **Main (Соцсеть):** http://localhost:3000 (frontend) + :8000 (backend)
- **Admin (Управление):** http://localhost:4000 (frontend) + :9000 (backend)
- **PetBase (Реестр):** http://localhost:4100 (frontend) + :8100 (backend)
- **Shelter (Кабинет приюта):** http://localhost:5100 (frontend) + :8200 (backend)

### Команды

```bash
# Сборка всех сервисов
docker-compose build

# Сборка конкретного сервиса
docker-compose build main-backend

# Запуск всех сервисов
docker-compose up

# Запуск в фоновом режиме
docker-compose up -d

# Запуск конкретного сервиса
docker-compose up main-backend

# Остановка
docker-compose down

# Остановка с удалением volumes
docker-compose down -v

# Просмотр логов
docker-compose logs -f main-backend

# Перезапуск сервиса
docker-compose restart main-backend

# Выполнить команду в контейнере
docker-compose exec main-backend sh
```

### Volumes

- `database/` - общая база данных SQLite для всех сервисов
- `uploads/` - загруженные файлы (аватары, фото, видео)

## ☸️ Kubernetes

TODO: Добавить Kubernetes манифесты для production деплоя

### Планируется

```
kubernetes/
├── main/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── ingress.yaml
├── admin/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── ingress.yaml
├── petbase/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── ingress.yaml
└── shared/
    ├── configmap.yaml
    ├── secrets.yaml
    └── persistent-volume.yaml
```

## 🚀 CI/CD

TODO: Добавить GitHub Actions / GitLab CI конфигурации

### Планируется

- Автоматическая сборка Docker образов
- Автоматический деплой в staging
- Автоматический деплой в production (с подтверждением)
- Автоматические тесты перед деплоем
- Rollback при ошибках

## 🔒 Безопасность

### Production Checklist

- [ ] Изменить JWT_SECRET на сильный ключ (минимум 32 символа)
- [ ] Настроить ALLOWED_ORIGINS только для production доменов
- [ ] Использовать HTTPS для всех сервисов
- [ ] Настроить firewall
- [ ] Регулярные бэкапы базы данных
- [ ] Мониторинг и алерты
- [ ] Rate limiting
- [ ] Логирование всех операций

### Генерация JWT Secret

```bash
# Сгенерировать сильный ключ
openssl rand -base64 32
```

## 📊 Мониторинг

TODO: Добавить Prometheus + Grafana

### Планируется

- Метрики производительности
- Мониторинг ошибок
- Алерты при проблемах
- Дашборды в Grafana

## 🔧 Troubleshooting

### Проблема: Порт уже занят

```bash
# Найти процесс на порту
lsof -i :8000

# Убить процесс
kill -9 <PID>
```

### Проблема: База данных не доступна

```bash
# Проверить volume
docker volume ls

# Проверить права доступа
docker-compose exec database ls -la /app/database
```

### Проблема: Сервис не запускается

```bash
# Просмотр логов
docker-compose logs main-backend

# Проверка статуса
docker-compose ps

# Перезапуск
docker-compose restart main-backend
```

---

**Последнее обновление:** 30 декабря 2024 (v0.7.0)  
**Статус:** В разработке (Docker готов, Kubernetes TODO)
