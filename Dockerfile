# Multi-stage build для zooplatforma проекта
# Этот Dockerfile собирает backend сервисы и Next.js фронтенд

FROM golang:1.25.5-alpine AS go-builder

# Установка зависимостей
RUN apk add --no-cache git make

WORKDIR /app

# Копируем go.mod и go.sum файлы для всех модулей
COPY database/go.mod database/go.sum ./database/
COPY auth/backend/go.mod auth/backend/go.sum ./auth/backend/
COPY main/backend/go.mod main/backend/go.sum ./main/backend/
COPY admin/backend/go.mod admin/backend/go.sum ./admin/backend/
COPY clinic/backend/go.mod clinic/backend/go.sum ./clinic/backend/
COPY owner/backend/go.mod owner/backend/go.sum ./owner/backend/
COPY petbase/backend/go.mod petbase/backend/go.sum ./petbase/backend/
COPY shelter/backend/go.mod shelter/backend/go.sum ./shelter/backend/
COPY volunteer/backend/go.mod volunteer/backend/go.sum ./volunteer/backend/
COPY pkg ./pkg

# Копируем весь проект (нужно для resolve local modules)
COPY . .

# Скачиваем Go зависимости (после копирования всех файлов для resolve local modules)
RUN go mod download -C database && \
    go mod tidy -C database && \
    go mod download -C auth/backend && \
    go mod tidy -C auth/backend && \
    go mod download -C main/backend && \
    go mod tidy -C main/backend && \
    go mod download -C admin/backend && \
    go mod tidy -C admin/backend && \
    go mod download -C clinic/backend && \
    go mod tidy -C clinic/backend && \
    go mod download -C owner/backend && \
    go mod tidy -C owner/backend && \
    go mod download -C petbase/backend && \
    go mod tidy -C petbase/backend && \
    go mod download -C shelter/backend && \
    go mod tidy -C shelter/backend && \
    go mod download -C volunteer/backend && \
    go mod tidy -C volunteer/backend

# Собираем все backend сервисы
RUN cd auth/backend && go build -o /app/bin/auth-backend . && \
    cd /app && \
    cd main/backend && go build -o /app/bin/main-backend . && \
    cd /app && \
    cd admin/backend && go build -o /app/bin/admin-backend . && \
    cd /app && \
    cd clinic/backend && go build -o /app/bin/clinic-backend . && \
    cd /app && \
    cd owner/backend && go build -o /app/bin/owner-backend . && \
    cd /app && \
    cd petbase/backend && go build -o /app/bin/petbase-backend . && \
    cd /app && \
    cd shelter/backend && go build -o /app/bin/shelter-backend . && \
    cd /app && \
    cd volunteer/backend && go build -o /app/bin/volunteer-backend .

# Next.js builder
FROM node:20-alpine AS next-builder

WORKDIR /app

# Копируем весь проект для сборки фронтенда
COPY main/frontend ./main/frontend
COPY shared ./shared

# Устанавливаем зависимости
RUN cd /app/main/frontend && npm install && \
    cd /app/shared && npm install

# Runtime образ
FROM node:20-alpine

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Копируем собранные Go бинарники
COPY --from=go-builder /app/bin/* /app/

# Копируем Next.js
COPY --from=next-builder /app/main/frontend /app/frontend
COPY --from=next-builder /app/main/frontend/node_modules /app/frontend/node_modules
COPY --from=next-builder /app/main/frontend/package.json /app/frontend/package.json

# Копируем миграции БД
COPY database/migrations /app/migrations

# Копируем конфигурационные файлы
COPY infrastructure /app/infrastructure

# Создаем .env файл с переменными окружения
RUN echo "DATABASE_URL=postgres://postgres_zp:7da0905cd3349f58f368@my_projects_bd_zooplatforma:5432/bd_zp?sslmode=disable" > /app/.env && \
    echo "DB_HOST=my_projects_bd_zooplatforma" >> /app/.env && \
    echo "DB_PORT=5432" >> /app/.env && \
    echo "DB_USER=postgres_zp" >> /app/.env && \
    echo "DB_PASSWORD=7da0905cd3349f58f368" >> /app/.env && \
    echo "DB_NAME=bd_zp" >> /app/.env && \
    echo "DB_SSLMODE=disable" >> /app/.env && \
    echo "JWT_SECRET=your-super-secret-key-change-this-in-production" >> /app/.env && \
    echo "ENVIRONMENT=production" >> /app/.env && \
    echo "LOG_LEVEL=info" >> /app/.env

# Создаем скрипт для запуска сервисов
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'SERVICE=${SERVICE:-main}' >> /app/start.sh && \
    echo 'case $SERVICE in' >> /app/start.sh && \
    echo '  auth)' >> /app/start.sh && \
    echo '    exec /app/auth-backend' >> /app/start.sh && \
    echo '    ;;' >> /app/start.sh && \
    echo '  main)' >> /app/start.sh && \
    echo '    /app/main-backend &' >> /app/start.sh && \
    echo '    cd /app/frontend && npm run dev' >> /app/start.sh && \
    echo '    ;;' >> /app/start.sh && \
    echo '  admin)' >> /app/start.sh && \
    echo '    exec /app/admin-backend' >> /app/start.sh && \
    echo '    ;;' >> /app/start.sh && \
    echo '  petbase)' >> /app/start.sh && \
    echo '    exec /app/petbase-backend' >> /app/start.sh && \
    echo '    ;;' >> /app/start.sh && \
    echo '  shelter)' >> /app/start.sh && \
    echo '    exec /app/shelter-backend' >> /app/start.sh && \
    echo '    ;;' >> /app/start.sh && \
    echo '  owner)' >> /app/start.sh && \
    echo '    exec /app/owner-backend' >> /app/start.sh && \
    echo '    ;;' >> /app/start.sh && \
    echo '  volunteer)' >> /app/start.sh && \
    echo '    exec /app/volunteer-backend' >> /app/start.sh && \
    echo '    ;;' >> /app/start.sh && \
    echo '  clinic)' >> /app/start.sh && \
    echo '    exec /app/clinic-backend' >> /app/start.sh && \
    echo '    ;;' >> /app/start.sh && \
    echo '  *)' >> /app/start.sh && \
    echo '    echo "Unknown service: $SERVICE"' >> /app/start.sh && \
    echo '    exit 1' >> /app/start.sh && \
    echo '    ;;' >> /app/start.sh && \
    echo 'esac' >> /app/start.sh && \
    chmod +x /app/start.sh

# Expose все порты
EXPOSE 7100 8000 8100 8200 8400 8500 8600 9000 3000 4000 4100 5100 6100 6200 6300

# Запускаем сервис (по умолчанию main)
CMD ["/app/start.sh"]
