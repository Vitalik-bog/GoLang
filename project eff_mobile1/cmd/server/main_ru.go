package main

import (
	"context"
	"effective_mobile_service/internal/config"
	"effective_mobile_service/internal/database"
	"effective_mobile_service/internal/router"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("../../config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Подключение к базе данных
	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	fmt.Println("Подключение к PostgreSQL успешно!")

	// Создание таблицы
	createTable(db)

	// Настройка роутера на русском
	router := router.SetupRouterRU(db)

	// Запуск сервера
	server := &http.Server{
		Addr:    cfg.Server.Host + cfg.Server.Port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		fmt.Printf("Сервер запускается на %s (Русская версия)\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Принудительное завершение: %v", err)
	}

	fmt.Println("Сервер остановлен")
}

func createTable(db *database.DB) {
	query := `
    CREATE TABLE IF NOT EXISTS subscriptions (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        service_name VARCHAR(100) NOT NULL,
        price INTEGER NOT NULL CHECK (price >= 0),
        user_id UUID NOT NULL,
        start_date DATE NOT NULL,
        end_date DATE,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
    )
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Printf("Не удалось создать таблицу: %v", err)
	} else {
		fmt.Println("Таблица 'subscriptions' создана")
	}
}
