package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	// Подключение к PostgreSQL
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=subscriptions_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(" Database connection failed:", err)
	}
	defer db.Close()

	// Проверка подключения
	if err := db.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL successfully!")

	// Создаем таблицу
	createTable(db)

	// Настройка HTTP сервера
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>Subscription Service</h1><p>Use API endpoints:</p><ul><li>GET /health</li><li>GET/POST /api/v1/subscriptions</li><li>GET /api/v1/subscriptions/total-cost</li></ul>")
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := db.Ping(); err != nil {
			json.NewEncoder(w).Encode(map[string]string{"status": "unhealthy"})
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	http.HandleFunc("/api/v1/subscriptions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch r.Method {
		case "GET":
			// Получение всех подписок
			rows, err := db.Query("SELECT id, service_name, price, user_id, start_date FROM subscriptions ORDER BY id DESC")
			if err != nil {
				http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			var subscriptions []map[string]interface{}
			for rows.Next() {
				var id int
				var serviceName, userID, startDate string
				var price int

				rows.Scan(&id, &serviceName, &price, &userID, &startDate)

				subscriptions = append(subscriptions, map[string]interface{}{
					"id":           id,
					"service_name": serviceName,
					"price":        price,
					"user_id":      userID,
					"start_date":   startDate,
				})
			}

			json.NewEncoder(w).Encode(map[string]interface{}{
				"count":         len(subscriptions),
				"subscriptions": subscriptions,
			})

		case "POST":
			// Создание подписки
			var req map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
				return
			}

			// Извлекаем поля
			serviceName, _ := req["service_name"].(string)
			priceFloat, _ := req["price"].(float64)
			price := int(priceFloat)
			userID, _ := req["user_id"].(string)
			startDate, _ := req["start_date"].(string)
			endDate, _ := req["end_date"].(string)

			// Валидация
			if serviceName == "" || price <= 0 || userID == "" || startDate == "" {
				http.Error(w, `{"error":"Missing required fields"}`, http.StatusBadRequest)
				return
			}

			// Сохранение в БД
			var id int
			err := db.QueryRow(
				"INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id",
				serviceName, price, userID, startDate, endDate,
			).Scan(&id)

			if err != nil {
				http.Error(w, `{"error":"Database error: `+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}

			response := map[string]interface{}{
				"id":           id,
				"service_name": serviceName,
				"price":        price,
				"user_id":      userID,
				"start_date":   startDate,
				"message":      "Subscription created",
			}

			if endDate != "" {
				response["end_date"] = endDate
			}

			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(response)

		default:
			http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/v1/subscriptions/total-cost", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID := r.URL.Query().Get("user_id")
		serviceName := r.URL.Query().Get("service_name")

		query := "SELECT COALESCE(SUM(price), 0) FROM subscriptions WHERE 1=1"
		var args []interface{}
		argNum := 1

		if userID != "" {
			query += fmt.Sprintf(" AND user_id = $%d", argNum)
			args = append(args, userID)
			argNum++
		}

		if serviceName != "" {
			query += fmt.Sprintf(" AND service_name ILIKE $%d", argNum)
			args = append(args, "%"+serviceName+"%")
			argNum++
		}

		var totalCost int
		err := db.QueryRow(query, args...).Scan(&totalCost)
		if err != nil {
			http.Error(w, `{"error":"Database error"}`, http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"total_cost": totalCost,
			"filters": map[string]interface{}{
				"user_id":      userID,
				"service_name": serviceName,
			},
		})
	})

	fmt.Println("🚀 Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
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
		log.Printf("Could not create table: %v", err)
	} else {
		fmt.Println("Table 'subscriptions' created")
	}
}
