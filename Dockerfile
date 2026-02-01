# Multi-stage build для zooplatforma проекта
# Этот Dockerfile собирает все backend сервисы

FROM golang:1.25.5-alpine AS builder

# Установка зависимостей
RUN apk add --no-cache git make

WORKDIR /app

# Копируем go.mod и go.sum файлы для всех модулей
COPY database/go.mod database/go.sum ./database/
COPY main/backend/go.mod main/backend/go.sum ./main/backend/
COPY admin/backend/go.mod admin/backend/go.sum ./admin/backend/
COPY clinic/backend/go.mod clinic/backend/go.sum ./clinic/backend/
COPY owner/backend/go.mod owner/backend/go.sum ./owner/backend/
COPY petbase/backend/go.mod petbase/backend/go.sum ./petbase/backend/
COPY shelter/backend/go.mod shelter/backend/go.sum ./shelter/backend/
COPY volunteer/backend/go.mod volunteer/backend/go.sum ./volunteer/backend/
COPY pkg ./pkg

# Скачиваем Go зависимости
RUN go mod download -C database && \
    go mod download -C main/backend && \
    go mod download -C admin/backend && \
    go mod download -C clinic/backend && \
    go mod download -C owner/backend && \
    go mod download -C petbase/backend && \
    go mod download -C shelter/backend && \
    go mod download -C volunteer/backend

# Копируем весь проект
COPY . .

# Собираем все backend сервисы
RUN cd main/backend && go build -o /app/bin/main-backend . && \
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

# Runtime образ
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Копируем собранные бинарники
COPY --from=builder /app/bin/* /app/

# Копируем миграции БД
COPY database/migrations /app/migrations

# Копируем конфигурационные файлы
COPY infrastructure /app/infrastructure

# Копируем .env файл
COPY .env /app/.env

# Expose все порты
EXPOSE 7100 8000 8100 8200 8400 8500 8600 9000 3000 4000 4100 5100 6100 6200 6300

# Запускаем main backend по умолчанию
# EasyPanel может переопределить это через CMD
CMD ["/app/main-backend"]
