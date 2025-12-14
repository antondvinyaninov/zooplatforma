# Infrastructure

Конфигурации для деплоя и инфраструктуры проекта.

## Структура

```
infrastructure/
├── docker/           # Docker конфигурации
│   ├── main/         # Dockerfile для main сервиса
│   ├── admin/        # Dockerfile для admin сервиса
│   └── mobile/       # Dockerfile для mobile (если нужен)
├── kubernetes/       # Kubernetes манифесты
│   ├── main/         # K8s конфиги для main
│   ├── admin/        # K8s конфиги для admin
│   └── shared/       # Общие ресурсы (ConfigMap, Secrets)
└── README.md         # Этот файл
```

## Docker

Для локальной разработки и production деплоя.

### Команды

```bash
# Сборка всех сервисов
docker-compose build

# Запуск всех сервисов
docker-compose up

# Остановка
docker-compose down
```

## Kubernetes

Для production деплоя в кластер.

### Команды

```bash
# Применить все манифесты
kubectl apply -f kubernetes/

# Проверить статус
kubectl get pods
kubectl get services
```

## CI/CD

TODO: Добавить GitHub Actions / GitLab CI конфигурации
