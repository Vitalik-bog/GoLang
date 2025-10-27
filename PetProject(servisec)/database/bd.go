package database

import (
	"database/sql"
	"fmt"
)

var DB *sql.DB

func InitDB(connectionString string) error {
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	query := `
    CREATE TABLE IF NOT EXISTS employees (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        surname VARCHAR(100) NOT NULL,
        phone VARCHAR(20),
        company_id INTEGER NOT NULL,
        passport_type VARCHAR(50),
        passport_number VARCHAR(50),
        department_name VARCHAR(100),
        department_phone VARCHAR(20)
    )`

	_, err = DB.Exec(query)
	if err != nil {
		return err
	}
	fmt.Println("✅ База данных подключена успешно")
	return nil
}
