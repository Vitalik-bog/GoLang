package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	// Подключение к PostgreSQL
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=subscriptions_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(" Ошибка подключения к БД:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Не удалось подключиться к БД:", err)
	}

	fmt.Println(" Подключение к PostgreSQL успешно!")

	// Создание таблицы
	createTable(db)

	// Настройка Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Русские заголовки
	router.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json; charset=utf-8")
		c.Header("Content-Language", "ru-RU")
		c.Next()
	})

	// Документация на русском
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service":  "Сервис управления подписками",
			"version":  "1.0",
			"language": "русский",
			"автор":    "Виталий",
			"endpoints": gin.H{
				"health":    "GET   /api/v1/health - Проверка работы сервиса",
				"создать":   "POST  /api/v1/subscriptions - Создать новую подписку",
				"список":    "GET   /api/v1/subscriptions - Получить список подписок",
				"стоимость": "GET   /api/v1/subscriptions/total-cost - Общая стоимость",
				"получить":  "GET   /api/v1/subscriptions/:id - Получить подписку по ID",
				"обновить":  "PUT   /api/v1/subscriptions/:id - Обновить подписку",
				"удалить":   "DELETE /api/v1/subscriptions/:id - Удалить подписку",
			},
			"пример_запроса": gin.H{
				"метод": "POST",
				"url":   "/api/v1/subscriptions",
				"тело": gin.H{
					"service_name": "Яндекс Плюс",
					"price":        400,
					"user_id":      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
					"start_date":   "07-2025",
				},
			},
		})
	})

	// API v1 на русском
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", func(c *gin.Context) {
			if err := db.Ping(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"статус":    "неисправен",
					"сообщение": "Ошибка подключения к базе данных",
					"ошибка":    err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"статус":      "исправен",
				"сообщение":   "Сервис работает нормально",
				"сервис":      "Управление подписками",
				"база_данных": "PostgreSQL",
				"время":       time.Now().Format("2006-01-02 15:04:05"),
			})
		})

		// Подписки
		subs := api.Group("/subscriptions")
		{
			// Создание подписки
			subs.POST("", func(c *gin.Context) {
				var request struct {
					ServiceName string `json:"service_name" binding:"required"`
					Price       int    `json:"price" binding:"required,min=1"`
					UserID      string `json:"user_id" binding:"required"`
					StartDate   string `json:"start_date" binding:"required"`
					EndDate     string `json:"end_date,omitempty"`
				}

				if err := c.ShouldBindJSON(&request); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"ошибка":    "Неверный запрос",
						"детали":    err.Error(),
						"сообщение": "Проверьте правильность заполнения полей",
					})
					return
				}

				var id int
				err := db.QueryRow(
					"INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
					request.ServiceName, request.Price, request.UserID, request.StartDate, request.EndDate,
				).Scan(&id)

				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"ошибка": "Ошибка базы данных",
						"детали": err.Error(),
					})
					return
				}

				c.JSON(http.StatusCreated, gin.H{
					"сообщение":     "Подписка успешно создана",
					"id":            id,
					"service_name":  request.ServiceName,
					"price":         request.Price,
					"user_id":       request.UserID,
					"start_date":    request.StartDate,
					"end_date":      request.EndDate,
					"дата_создания": time.Now().Format("2006-01-02 15:04:05"),
				})
			})

			// Список подписок
			subs.GET("", func(c *gin.Context) {
				rows, err := db.Query("SELECT id, service_name, price, user_id, start_date FROM subscriptions ORDER BY id DESC")
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"ошибка": "Ошибка базы данных",
						"детали": err.Error(),
					})
					return
				}
				defer rows.Close()

				var subscriptions []gin.H
				for rows.Next() {
					var id int
					var serviceName, userID, startDate string
					var price int

					rows.Scan(&id, &serviceName, &price, &userID, &startDate)

					subscriptions = append(subscriptions, gin.H{
						"id":           id,
						"service_name": serviceName,
						"price":        price,
						"user_id":      userID,
						"start_date":   startDate,
					})
				}

				c.JSON(http.StatusOK, gin.H{
					"сообщение":  "Список подписок",
					"количество": len(subscriptions),
					"подписки":   subscriptions,
				})
			})

			// Общая стоимость
			subs.GET("/total-cost", func(c *gin.Context) {
				userID := c.Query("user_id")
				serviceName := c.Query("service_name")

				query := "SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE 1=1"
				var args []interface{}

				if userID != "" {
					query += " AND user_id = $1"
					args = append(args, userID)
				}

				if serviceName != "" {
					if len(args) == 0 {
						query += " AND service_name ILIKE $1"
					} else {
						query += " AND service_name ILIKE $2"
					}
					args = append(args, "%"+serviceName+"%")
				}

				var totalCost int
				err := db.QueryRow(query, args...).Scan(&totalCost)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"ошибка": "Ошибка расчета стоимости",
						"детали": err.Error(),
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"сообщение":       "Общая стоимость рассчитана",
					"общая_стоимость": totalCost,
					"валюта":          "RUB",
					"фильтры": gin.H{
						"user_id":      userID,
						"service_name": serviceName,
					},
				})
			})
		}
	}

	// Graceful shutdown
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		fmt.Println("Сервер запущен на http://localhost:8080")
		fmt.Println("Документация: http://localhost:8080/")
		fmt.Println("Health check: http://localhost:8080/api/v1/health")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка завершения: %v", err)
	}

	fmt.Println("Сервер остановлен")
}

func createTable(db *sql.DB) {
	query := `
    CREATE TABLE IF NOT EXISTS subscriptions (
        id SERIAL PRIMARY KEY,
        service_name VARCHAR(100) NOT NULL,
        price INTEGER NOT NULL CHECK (price >= 0),
        user_id VARCHAR(36) NOT NULL,
        start_date VARCHAR(10) NOT NULL,
        end_date VARCHAR(10),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )
    `
	_, err := db.Exec(query)
	if err != nil {
		log.Printf(" Не удалось создать таблицу: %v", err)
	} else {
		fmt.Println("Таблица 'subscriptions' создана")
	}
}
