# Multi-stage build для zooplatforma проекта
# Этот Dockerfile собирает backend сервисы и Next.js фронтенд

FROM golang:1.25.5-alpine AS go-builder

# Установка зависимостей
RUN apk add --no-cache git make build-base

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
# Auth Service требует CGO для sqlite3
RUN cd auth/backend && CGO_ENABLED=1 go build -o /app/bin/auth-backend . && \
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

# Копируем shared package (включая .tgz файл)
COPY shared ./shared

# Копируем main/frontend
COPY main/frontend ./main/frontend

# Устанавливаем зависимости shared (если нужно)
RUN cd /app/shared && npm install || true

# Устанавливаем зависимости main/frontend
RUN cd /app/main/frontend && npm install

# Собираем Next.js (для production)
RUN cd /app/main/frontend && npm run build

# Runtime образ
FROM node:20-alpine

RUN apk add --no-cache ca-certificates nginx postgresql-client

WORKDIR /app

# Копируем собранные Go бинарники
COPY --from=go-builder /app/bin/* /app/

# Копируем Next.js (production build и все необходимые файлы)
COPY --from=next-builder /app/main/frontend/.next /app/frontend/.next
COPY --from=next-builder /app/main/frontend/app /app/frontend/app
COPY --from=next-builder /app/main/frontend/contexts /app/frontend/contexts
COPY --from=next-builder /app/main/frontend/lib /app/frontend/lib
COPY --from=next-builder /app/main/frontend/types /app/frontend/types
COPY --from=next-builder /app/main/frontend/public /app/frontend/public
COPY --from=next-builder /app/main/frontend/node_modules /app/frontend/node_modules
COPY --from=next-builder /app/main/frontend/package.json /app/frontend/package.json
COPY --from=next-builder /app/main/frontend/next.config.ts /app/frontend/next.config.ts
COPY --from=next-builder /app/main/frontend/tsconfig.json /app/frontend/tsconfig.json
COPY --from=next-builder /app/main/frontend/next-env.d.ts /app/frontend/next-env.d.ts
COPY --from=next-builder /app/main/frontend/postcss.config.mjs /app/frontend/postcss.config.mjs
COPY --from=next-builder /app/main/frontend/tailwind.config.ts /app/frontend/tailwind.config.ts

# Копируем миграции БД
COPY database/migrations /app/migrations

# Копируем SQL fix для organizations
COPY fix_organizations_table.sql /app/fix_organizations_table.sql

# Копируем конфигурационные файлы
COPY infrastructure /app/infrastructure

# Копируем nginx конфиг в правильное место
COPY nginx.conf /etc/nginx/nginx.conf

# Создаем скрипт для запуска сервисов
RUN cat > /app/start.sh << 'EOF'
#!/bin/sh
set -e

SERVICE=${SERVICE:-main}

case $SERVICE in
  auth)
    exec /app/auth-backend
    ;;
  main)
    # Применяем SQL fix для organizations (если в production)
    if [ "$ENVIRONMENT" = "production" ]; then
      echo "🔧 Applying organizations table fix..."
      PGPASSWORD=${DATABASE_PASSWORD:-lmLG7k2ed4vas19} psql -h ${DATABASE_HOST:-zooplatforma-db} -U ${DATABASE_USER:-zp} -d ${DATABASE_NAME:-zp-db} -f /app/fix_organizations_table.sql || echo "⚠️ SQL fix failed (maybe already applied)"
    fi
    
    # Запускаем Auth Service (порт 7100)
    echo "🚀 Starting Auth Service..."
    /app/auth-backend &
    AUTH_PID=$!
    
    # Ждем пока Auth Service запустится
    sleep 2
    
    # Запускаем Main Backend (порт 8000)
    echo "🚀 Starting Main Backend..."
    export AUTH_SERVICE_URL=http://localhost:7100
    /app/main-backend &
    BACKEND_PID=$!
    
    # Запускаем frontend (production режим)
    echo "🚀 Starting Main Frontend..."
    cd /app/frontend && NEXT_PUBLIC_API_URL=http://localhost:8000 NEXT_PUBLIC_AUTH_URL=http://localhost:7100 npm start &
    FRONTEND_PID=$!
    
    # Ждем любого процесса
    wait -n
    ;;
  admin)
    exec /app/admin-backend
    ;;
  petbase)
    exec /app/petbase-backend
    ;;
  shelter)
    exec /app/shelter-backend
    ;;
  owner)
    exec /app/owner-backend
    ;;
  volunteer)
    exec /app/volunteer-backend
    ;;
  clinic)
    exec /app/clinic-backend
    ;;
  *)
    echo "Unknown service: $SERVICE"
    exit 1
    ;;
esac
EOF
RUN chmod +x /app/start.sh


# Expose все порты
EXPOSE 7100 8000 8100 8200 8400 8500 8600 9000 3000 4000 4100 5100 6100 6200 6300

# Запускаем сервис (по умолчанию main)
CMD ["/app/start.sh"]
