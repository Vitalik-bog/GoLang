package main

import (
	"employee-service/config"
	"employee-service/database"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // ← ЭТА СТРОКА ОБЯЗАТЕЛЬНА!(без нее НЕ  рабоатл постргресс)
)

func main() {
	cfg := config.Load()

	if err := database.InitDB(cfg.DatabaseURL); err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}

	router := mux.NewRouter()

	// TODO: добавить обработчики позже
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("🚀 Сервер работает!"))
	})

	log.Println("🚀 Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
