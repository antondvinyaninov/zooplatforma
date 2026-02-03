# Gateway Dockerfile - минимальный образ для API Gateway
FROM golang:1.25-alpine AS builder

# Установить зависимости для сборки
RUN apk add --no-cache git

WORKDIR /build

# Копировать go.mod и go.sum
COPY gateway/go.mod gateway/go.sum ./

# Скачать зависимости
RUN go mod download

# Копировать исходный код Gateway
COPY gateway/*.go ./

# Собрать Gateway
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway .

# Runtime stage - минимальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Копировать бинарник Gateway
COPY --from=builder /build/gateway .

# Установить timezone
ENV TZ=Europe/Moscow

# Порт Gateway
EXPOSE 80

# Запустить Gateway
CMD ["./gateway"]
