package repository

import (
	"database/sql"
	"day19/models"
	"fmt"
	"log"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) (int, error) {
	query := `
        INSERT INTO users (name, email, age)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	var userID int
	err := r.db.QueryRow(query, user.Name, user.Email, user.Age).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("ошибка создания пользователя: %v", err)
	}

	return userID, nil
}

func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	query := `
        SELECT id, name, email, age, created_at
        FROM users
        WHERE id = $1
    `

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь с ID %d не найден", id)
		}
		return nil, fmt.Errorf("ошибка получения пользователя: %v", err)
	}

	return &user, nil
}

func (r *UserRepository) GetAllUsers() ([]*models.User, error) {
	query := `
        SELECT id, name, email, age, created_at
        FROM users
        ORDER BY id
    `

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Age,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %v", err)
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации строк: %v", err)
	}

	return users, nil
}

func (r *UserRepository) UpdateUserAge(id int, newAge int) error {
	query := `UPDATE users SET age = $1 WHERE id = $2`

	result, err := r.db.Exec(query, newAge, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления пользователя: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества обновленных строк: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("пользователь с ID %d не найден", id)
	}

	log.Printf("Возраст пользователя ID %d обновлен. Затронуто строк: %d", id, rowsAffected)
	return nil
}

func (r *UserRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления пользователя: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удаленных строк: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("пользователь с ID %d не найден", id)
	}

	log.Printf("Пользователь ID %d удален. Затронуто строк: %d", id, rowsAffected)
	return nil
}

func (r *UserRepository) InitDatabase() error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(100) NOT NULL,
            email VARCHAR(100) UNIQUE NOT NULL,
            age INTEGER NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы: %v", err)
	}

	log.Println("Таблица users готова к работе")
	return nil
}
