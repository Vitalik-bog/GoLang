package main

import (
	"employee-service/config"
	"employee-service/database"
	"employee-service/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Загрузка конфигурацииыы
	cfg := config.Load()

	// Инициализация базы данных
	if err := database.InitDB(cfg.DatabaseURL); err != nil {
		log.Fatal("❌ Database connection failed:", err)
	}

	// Создание роутера
	router := mux.NewRouter()
	employeeHandler := handlers.NewEmployeeHandler(database.DB)

	// Регистрация маршрутов
	router.HandleFunc("/employees", employeeHandler.AddEmployee).Methods("POST")
	router.HandleFunc("/employees/{id}", employeeHandler.DeleteEmployee).Methods("DELETE")
	router.HandleFunc("/employees/{id}", employeeHandler.UpdateEmployee).Methods("PUT")
	router.HandleFunc("/company/{companyId}/employees", employeeHandler.GetEmployeesByCompany).Methods("GET")
	router.HandleFunc("/company/{companyId}/department/{departmentName}/employees", employeeHandler.GetEmployeesByDepartment).Methods("GET")

	// Health check endpoint, который позволяет проверить что сервер работает и отвечает.
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("🚀 Employee Service API is running!"))
	})

	log.Println("✅ База данных подключена успешно")
	log.Println("🚀 Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
