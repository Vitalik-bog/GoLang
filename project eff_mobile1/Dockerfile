FROM golang:1.21-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости для миграций
RUN apk add --no-cache git

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник из builder
COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .
COPY --from=builder /app/internal/database/migrations ./migrations

# Экспонируем порт
EXPOSE 8080

# Запуск приложения
CMD ["./main"]
