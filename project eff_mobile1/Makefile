.PHONY: build run test migrate-up migrate-down docker-up docker-down clean

# Переменные
APP_NAME=effective_mobile_service
DB_URL=postgres://postgres:postgres@localhost:5432/subscriptions_db?sslmode=disable

# Сборка приложения
build:
	go build -o bin/${APP_NAME} ./cmd/server

# Запуск приложения
run: build
	./bin/${APP_NAME}

# Запуск тестов
test:
	go test ./... -v

# Запуск миграций вверх
migrate-up:
	migrate -path internal/database/migrations -database "${DB_URL}" up

# Запуск миграций вниз
migrate-down:
	migrate -path internal/database/migrations -database "${DB_URL}" down

# Запуск в Docker
docker-up:
	docker-compose up -d --build

# Остановка Docker
docker-down:
	docker-compose down

# Очистка
clean:
	rm -rf bin/
	go clean

# Генерация Swagger документации
swagger:
	swag init -g cmd/server/main.go -o docs

# Запуск с горячей перезагрузкой (установите air: go install github.com/cosmtrek/air@latest)
dev:
	air
