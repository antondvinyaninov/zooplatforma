# Dockerfile for Backend
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/ .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Production stage
FROM alpine:latest

WORKDIR /app

# Copy binary
COPY --from=builder /app/main .

# Create uploads directory
RUN mkdir -p ./uploads

EXPOSE 8000

CMD ["./main"]
