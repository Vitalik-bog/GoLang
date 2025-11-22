package main

import (
	"database/sql"
	"day19/models"
	"day19/repository"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Подключение к базе данных
	connStr := "user=postgres password=your_password dbname=testdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	// Проверка подключения
	err = db.Ping()
	if err != nil {
		log.Fatal("Ошибка ping БД:", err)
	}
	fmt.Println("Успешное подключение к базе данных!")

	// Инициализация репозитория
	userRepo := repository.NewUserRepository(db)

	// Создание таблицы
	err = userRepo.InitDatabase()
	if err != nil {
		log.Fatal("Ошибка инициализации БД:", err)
	}

	// Демонстрация работы всех методов
	demonstrateExec(userRepo)
	demonstrateQueryRow(userRepo)
	demonstrateQuery(userRepo)
}

func demonstrateExec(userRepo *repository.UserRepository) {
	fmt.Println("\n=== Демонстрация Exec (CREATE, UPDATE, DELETE) ===")

	// Создание пользователей
	users := []*models.User{
		{Name: "Алексей Петров", Email: "alex@example.com", Age: 28},
		{Name: "Мария Сидорова", Email: "maria@example.com", Age: 25},
		{Name: "Иван Козлов", Email: "ivan@example.com", Age: 32},
	}

	for _, user := range users {
		userID, err := userRepo.CreateUser(user)
		if err != nil {
			log.Printf("Ошибка создания пользователя %s: %v", user.Name, err)
			continue
		}
		fmt.Printf("Создан пользователь: %s (ID: %d)\n", user.Name, userID)
	}
}

func demonstrateQueryRow(userRepo *repository.UserRepository) {
	fmt.Println("\n=== Демонстрация QueryRow (получение одной строки) ===")

	// Получение пользователя по ID
	user, err := userRepo.GetUserByID(1)
	if err != nil {
		log.Printf("Ошибка получения пользователя: %v", err)
	} else {
		fmt.Printf("Найден пользователь: ID=%d, Name=%s, Email=%s, Age=%d\n",
			user.ID, user.Name, user.Email, user.Age)
	}

	// Попытка получить несуществующего пользователя
	_, err = userRepo.GetUserByID(999)
	if err != nil {
		fmt.Printf("Ожидаемая ошибка: %v\n", err)
	}
}

func demonstrateQuery(userRepo *repository.UserRepository) {
	fmt.Println("\n=== Демонстрация Query (получение множества строк) ===")

	// Получение всех пользователей
	users, err := userRepo.GetAllUsers()
	if err != nil {
		log.Printf("Ошибка получения пользователей: %v", err)
		return
	}

	fmt.Printf("Найдено пользователей: %d\n", len(users))
	for _, user := range users {
		fmt.Printf("  - ID: %d, Имя: %s, Email: %s, Возраст: %d, Создан: %s\n",
			user.ID, user.Name, user.Email, user.Age,
			user.CreatedAt.Format("2006-01-02 15:04"))
	}

	// Демонстрация UPDATE и DELETE
	fmt.Println("\n=== Демонстрация UPDATE и DELETE ===")

	// Обновление возраста
	err = userRepo.UpdateUserAge(1, 29)
	if err != nil {
		log.Printf("Ошибка обновления: %v", err)
	}

	// Удаление пользователя
	err = userRepo.DeleteUser(3)
	if err != nil {
		log.Printf("Ошибка удаления: %v", err)
	}

	// Покажем обновленный список
	fmt.Println("\n=== Обновленный список пользователей ===")
	users, err = userRepo.GetAllUsers()
	if err != nil {
		log.Printf("Ошибка получения пользователей: %v", err)
		return
	}

	for _, user := range users {
		fmt.Printf("  - ID: %d, Имя: %s, Возраст: %d\n",
			user.ID, user.Name, user.Age)
	}
}
