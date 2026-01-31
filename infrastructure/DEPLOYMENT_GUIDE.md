# Руководство по деплою проекта

## 📋 Требования перед деплоем

### Обязательно:
- [ ] Все тесты проходят
- [ ] Нет критических ошибок в логах
- [ ] Auth Service работает
- [ ] PetBase работает
- [ ] Main Backend работает
- [ ] Main Frontend собирается без ошибок
- [ ] Все переменные окружения установлены
- [ ] SSL сертификат готов

### Опционально:
- [ ] Резервная копия БД сделана
- [ ] Мониторинг настроен
- [ ] Логирование настроено
- [ ] CDN настроен для статики

---

## 🚀 Этапы деплоя

### Этап 1: Подготовка сервера

```bash
# 1. Обновить систему
sudo apt update && sudo apt upgrade -y

# 2. Установить Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 3. Установить Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 4. Установить Nginx
sudo apt install nginx -y

# 5. Установить Certbot для SSL
sudo apt install certbot python3-certbot-nginx -y
```

### Этап 2: Клонирование проекта

```bash
# 1. Клонировать репозиторий
git clone https://github.com/yourusername/zooplatforma.git
cd zooplatforma

# 2. Создать .env файл
cp .env.example .env

# 3. Отредактировать .env с production значениями
nano .env
```

### Этап 3: Сборка Docker образов

```bash
# 1. Собрать все образы
docker-compose -f infrastructure/docker-compose.yml build

# 2. Проверить образы
docker images
```

### Этап 4: Запуск контейнеров

```bash
# 1. Запустить все сервисы
docker-compose -f infrastructure/docker-compose.yml up -d

# 2. Проверить статус
docker-compose -f infrastructure/docker-compose.yml ps

# 3. Проверить логи
docker-compose -f infrastructure/docker-compose.yml logs -f
```

### Этап 5: Настройка Nginx

```bash
# 1. Создать конфиг Nginx
sudo nano /etc/nginx/sites-available/zooplatforma

# 2. Включить сайт
sudo ln -s /etc/nginx/sites-available/zooplatforma /etc/nginx/sites-enabled/

# 3. Проверить конфиг
sudo nginx -t

# 4. Перезагрузить Nginx
sudo systemctl reload nginx
```

### Этап 6: SSL сертификат

```bash
# 1. Получить сертификат Let's Encrypt
sudo certbot --nginx -d yourdomain.com -d www.yourdomain.com

# 2. Проверить автоматическое обновление
sudo systemctl enable certbot.timer
sudo systemctl start certbot.timer
```

### Этап 7: Проверка

```bash
# 1. Проверить что все сервисы работают
curl http://localhost:7100/api/health  # Auth Service
curl http://localhost:8100/api/health  # PetBase
curl http://localhost:8000/api/health  # Main Backend

# 2. Проверить фронтенд
curl https://yourdomain.com

# 3. Проверить логи
docker-compose logs -f
```

---

## 🔧 Переменные окружения (.env)

```bash
# Auth Service
AUTH_PORT=7100
AUTH_DB_PATH=/database/auth.db
JWT_SECRET=your-super-secret-key-change-this

# Main Backend
MAIN_BACKEND_PORT=8000
MAIN_DB_PATH=/database/data.db
AUTH_SERVICE_URL=http://auth:7100

# Main Frontend
MAIN_FRONTEND_PORT=3000
NEXT_PUBLIC_API_URL=https://yourdomain.com/api
NEXT_PUBLIC_AUTH_URL=https://yourdomain.com/auth

# PetBase Backend
PETBASE_PORT=8100
PETBASE_DB_PATH=/database/data.db

# Database
DATABASE_PATH=/database/data.db

# General
NODE_ENV=production
LOG_LEVEL=info
```

---

## 📊 Мониторинг и логирование

### Логи контейнеров

```bash
# Все логи
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f auth
docker-compose logs -f main-backend
docker-compose logs -f main-frontend

# Последние 100 строк
docker-compose logs --tail=100
```

### Проверка здоровья сервисов

```bash
# Проверить статус
docker-compose ps

# Проверить ресурсы
docker stats

# Проверить сеть
docker network ls
```

---

## 🔄 Обновление проекта

```bash
# 1. Остановить контейнеры
docker-compose down

# 2. Обновить код
git pull origin main

# 3. Пересобрать образы
docker-compose build

# 4. Запустить снова
docker-compose up -d

# 5. Проверить логи
docker-compose logs -f
```

---

## 🚨 Troubleshooting

### Сервис не запускается

```bash
# 1. Проверить логи
docker-compose logs service-name

# 2. Проверить конфиг
docker-compose config

# 3. Перезапустить сервис
docker-compose restart service-name
```

### Проблемы с БД

```bash
# 1. Проверить том
docker volume ls

# 2. Проверить права доступа
docker exec container-name ls -la /database

# 3. Восстановить из резервной копии
docker cp backup.db container-name:/database/data.db
```

### Проблемы с сетью

```bash
# 1. Проверить сеть
docker network inspect zooplatforma_default

# 2. Проверить DNS
docker exec container-name nslookup auth

# 3. Проверить порты
docker port container-name
```

---

## 📝 Checklist перед production

- [ ] Все сервисы запущены и работают
- [ ] SSL сертификат установлен
- [ ] Nginx настроен и работает
- [ ] Резервная копия БД сделана
- [ ] Логирование настроено
- [ ] Мониторинг настроен
- [ ] Алерты настроены
- [ ] Backup стратегия определена
- [ ] Rollback план готов
- [ ] Документация обновлена

---

## 🔐 Безопасность

### Обязательно:
- [ ] Изменить все default пароли
- [ ] Установить SSL сертификат
- [ ] Настроить firewall
- [ ] Включить логирование
- [ ] Настроить rate limiting
- [ ] Включить CORS только для нужных доменов
- [ ] Скрыть версии сервисов в headers

### Рекомендуется:
- [ ] Настроить WAF (Web Application Firewall)
- [ ] Включить DDoS protection
- [ ] Настроить VPN для админ доступа
- [ ] Включить 2FA для админов
- [ ] Регулярно обновлять зависимости

---

## 📞 Поддержка

При проблемах:
1. Проверить логи: `docker-compose logs -f`
2. Проверить статус: `docker-compose ps`
3. Перезапустить сервис: `docker-compose restart service-name`
4. Проверить документацию: `docs/`
