# Deployment Guide

Руководство по развёртыванию ЗооПлатформы в production.

## 📋 Предварительные требования

### Сервер

- Ubuntu 22.04 LTS или новее
- Минимум 2 CPU, 4GB RAM
- 50GB свободного места на диске
- Docker 24.0+ и Docker Compose 2.0+

### Домены

Настройте DNS записи для:
- `zooplatform.ru` → Main Frontend
- `api.zooplatform.ru` → Main Backend
- `admin.zooplatform.ru` → Admin Frontend
- `admin-api.zooplatform.ru` → Admin Backend
- `petid.zooplatform.ru` → PetBase Frontend
- `petid-api.zooplatform.ru` → PetBase Backend
- `shelter.zooplatform.ru` → Shelter Frontend
- `shelter-api.zooplatform.ru` → Shelter Backend
- `owner.zooplatform.ru` → Owner Frontend
- `owner-api.zooplatform.ru` → Owner Backend
- `volunteer.zooplatform.ru` → Volunteer Frontend
- `volunteer-api.zooplatform.ru` → Volunteer Backend
- `clinic.zooplatform.ru` → Clinic Frontend
- `clinic-api.zooplatform.ru` → Clinic Backend

## 🚀 Быстрый деплой

### 1. Подготовка сервера

```bash
# Обновление системы
sudo apt update && sudo apt upgrade -y

# Установка Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Установка Docker Compose
sudo apt install docker-compose-plugin

# Добавление пользователя в группу docker
sudo usermod -aG docker $USER
newgrp docker
```

### 2. Клонирование репозитория

```bash
# Клонирование
git clone https://github.com/your-org/zooplatform.git
cd zooplatform

# Переход на production ветку
git checkout production
```

### 3. Настройка переменных окружения

```bash
# Копирование примера
cp infrastructure/.env.example infrastructure/.env

# Генерация JWT Secret
openssl rand -base64 32

# Редактирование .env
nano infrastructure/.env
```

**Обязательно установите:**
- `JWT_SECRET` - сгенерированный ключ
- `ALLOWED_ORIGINS` - ваши production домены
- `NEXT_PUBLIC_*` - production URL'ы

### 4. Настройка SSL (Let's Encrypt)

```bash
# Установка Certbot
sudo apt install certbot python3-certbot-nginx

# Получение сертификатов для всех доменов
sudo certbot certonly --standalone -d zooplatform.ru
sudo certbot certonly --standalone -d api.zooplatform.ru
sudo certbot certonly --standalone -d admin.zooplatform.ru
sudo certbot certonly --standalone -d admin-api.zooplatform.ru
sudo certbot certonly --standalone -d petid.zooplatform.ru
sudo certbot certonly --standalone -d petid-api.zooplatform.ru
sudo certbot certonly --standalone -d shelter.zooplatform.ru
sudo certbot certonly --standalone -d shelter-api.zooplatform.ru
sudo certbot certonly --standalone -d owner.zooplatform.ru
sudo certbot certonly --standalone -d owner-api.zooplatform.ru
sudo certbot certonly --standalone -d volunteer.zooplatform.ru
sudo certbot certonly --standalone -d volunteer-api.zooplatform.ru
sudo certbot certonly --standalone -d clinic.zooplatform.ru
sudo certbot certonly --standalone -d clinic-api.zooplatform.ru

# Автоматическое обновление сертификатов
sudo certbot renew --dry-run
```

### 5. Запуск сервисов

```bash
cd infrastructure

# Сборка образов
docker-compose build

# Запуск в фоновом режиме
docker-compose up -d

# Проверка статуса
docker-compose ps

# Просмотр логов
docker-compose logs -f
```

### 6. Настройка Nginx (Reverse Proxy)

```bash
# Установка Nginx
sudo apt install nginx

# Создание конфигурации
sudo nano /etc/nginx/sites-available/zooplatform
```

**Конфигурация Nginx:**

```nginx
# Main Frontend
server {
    listen 443 ssl http2;
    server_name zooplatform.ru;

    ssl_certificate /etc/letsencrypt/live/zooplatform.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/zooplatform.ru/privkey.pem;

    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}

# Main Backend
server {
    listen 443 ssl http2;
    server_name api.zooplatform.ru;

    ssl_certificate /etc/letsencrypt/live/api.zooplatform.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.zooplatform.ru/privkey.pem;

    location / {
        proxy_pass http://localhost:8000;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Admin Frontend
server {
    listen 443 ssl http2;
    server_name admin.zooplatform.ru;

    ssl_certificate /etc/letsencrypt/live/admin.zooplatform.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/admin.zooplatform.ru/privkey.pem;

    location / {
        proxy_pass http://localhost:4000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}

# Admin Backend
server {
    listen 443 ssl http2;
    server_name admin-api.zooplatform.ru;

    ssl_certificate /etc/letsencrypt/live/admin-api.zooplatform.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/admin-api.zooplatform.ru/privkey.pem;

    location / {
        proxy_pass http://localhost:9000;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# PetBase Frontend
server {
    listen 443 ssl http2;
    server_name petid.zooplatform.ru;

    ssl_certificate /etc/letsencrypt/live/petid.zooplatform.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/petid.zooplatform.ru/privkey.pem;

    location / {
        proxy_pass http://localhost:4100;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}

# PetBase Backend
server {
    listen 443 ssl http2;
    server_name petid-api.zooplatform.ru;

    ssl_certificate /etc/letsencrypt/live/petid-api.zooplatform.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/petid-api.zooplatform.ru/privkey.pem;

    location / {
        proxy_pass http://localhost:8100;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Shelter Frontend
server {
    listen 443 ssl http2;
    server_name shelter.zooplatform.ru;

    ssl_certificate /etc/letsencrypt/live/shelter.zooplatform.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/shelter.zooplatform.ru/privkey.pem;

    location / {
        proxy_pass http://localhost:5100;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}

# Shelter Backend
server {
    listen 443 ssl http2;
    server_name shelter-api.zooplatform.ru;

    ssl_certificate /etc/letsencrypt/live/shelter-api.zooplatform.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/shelter-api.zooplatform.ru/privkey.pem;

    location / {
        proxy_pass http://localhost:8200;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Owner Frontend
server {
    listen 443 ssl http2;
    server_name owner.zooplatform.ru;

    ssl_certificate /etc/letsencrypt/live/owner.zooplatform.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/owner.zooplatform.ru/privkey.pem;

    location / {
        proxy_pass http://localhost:6100;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}

# Owner Backend
server {
    listen 443 ssl http2;
    server_name owner-api.zooplatform.ru;

    ssl_certificate /etc/letsencrypt/live/owner-api.zooplatform.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/owner-api.zooplatform.ru/privkey.pem;

    location / {
        proxy_pass http://localhost:8400;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Volunteer Frontend
server {
    listen 443 ssl http2;
    server_name volunteer.zooplatform.ru;

    ssl_certificate /etc/letsencrypt/live/volunteer.zooplatform.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/volunteer.zooplatform.ru/privkey.pem;

    location / {
        proxy_pass http://localhost:6200;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}

# Volunteer Backend
server {
    listen 443 ssl http2;
    server_name volunteer-api.zooplatform.ru;

    ssl_certificate /etc/letsencrypt/live/volunteer-api.zooplatform.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/volunteer-api.zooplatform.ru/privkey.pem;

    location / {
        proxy_pass http://localhost:8500;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Clinic Frontend
server {
    listen 443 ssl http2;
    server_name clinic.zooplatform.ru;

    ssl_certificate /etc/letsencrypt/live/clinic.zooplatform.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/clinic.zooplatform.ru/privkey.pem;

    location / {
        proxy_pass http://localhost:6300;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }
}

# Clinic Backend
server {
    listen 443 ssl http2;
    server_name clinic-api.zooplatform.ru;

    ssl_certificate /etc/letsencrypt/live/clinic-api.zooplatform.ru/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/clinic-api.zooplatform.ru/privkey.pem;

    location / {
        proxy_pass http://localhost:8600;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# HTTP to HTTPS redirect
server {
    listen 80;
    server_name zooplatform.ru api.zooplatform.ru admin.zooplatform.ru admin-api.zooplatform.ru petid.zooplatform.ru petid-api.zooplatform.ru shelter.zooplatform.ru shelter-api.zooplatform.ru owner.zooplatform.ru owner-api.zooplatform.ru volunteer.zooplatform.ru volunteer-api.zooplatform.ru clinic.zooplatform.ru clinic-api.zooplatform.ru;
    return 301 https://$server_name$request_uri;
}
```

```bash
# Активация конфигурации
sudo ln -s /etc/nginx/sites-available/zooplatform /etc/nginx/sites-enabled/

# Проверка конфигурации
sudo nginx -t

# Перезапуск Nginx
sudo systemctl restart nginx
```

## 🔒 Безопасность

### Firewall

```bash
# Установка UFW
sudo apt install ufw

# Разрешить SSH
sudo ufw allow 22/tcp

# Разрешить HTTP/HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Включить firewall
sudo ufw enable

# Проверить статус
sudo ufw status
```

### Fail2Ban

```bash
# Установка
sudo apt install fail2ban

# Создание конфигурации
sudo cp /etc/fail2ban/jail.conf /etc/fail2ban/jail.local

# Редактирование
sudo nano /etc/fail2ban/jail.local

# Запуск
sudo systemctl enable fail2ban
sudo systemctl start fail2ban
```

## 📊 Мониторинг

### Логи

```bash
# Логи всех сервисов
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f main-backend

# Логи Nginx
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### Метрики

TODO: Добавить Prometheus + Grafana

## 🔄 Обновление

```bash
# Остановка сервисов
docker-compose down

# Обновление кода
git pull origin production

# Пересборка образов
docker-compose build

# Запуск
docker-compose up -d

# Проверка
docker-compose ps
```

## 💾 Бэкапы

### Автоматические бэкапы

```bash
# Создание cron задачи
crontab -e

# Добавить строку (бэкап каждый день в 3:00)
0 3 * * * /path/to/zooplatform/scripts/backup-all.sh
```

### Ручной бэкап

```bash
# База данных
./scripts/backup-database.sh

# Загруженные файлы
./scripts/backup-uploads.sh

# Всё вместе
./scripts/backup-all.sh
```

## 🆘 Troubleshooting

### Сервис не запускается

```bash
# Проверка логов
docker-compose logs main-backend

# Проверка статуса
docker-compose ps

# Перезапуск
docker-compose restart main-backend
```

### База данных недоступна

```bash
# Проверка прав доступа
ls -la database/

# Восстановление из бэкапа
./scripts/restore-database.sh database/backups/latest.db
```

### Проблемы с SSL

```bash
# Проверка сертификатов
sudo certbot certificates

# Обновление сертификатов
sudo certbot renew

# Проверка Nginx
sudo nginx -t
```

---

**Последнее обновление:** 31 декабря 2024 (v0.8.0)  
**Статус:** Готово к использованию