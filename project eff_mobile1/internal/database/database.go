package database

import (
    "database/sql"
    "fmt"
    "effective_mobile_service/internal/config"
    _ "github.com/lib/pq"
)

type DB struct {
    *sql.DB
}

func Connect(dbConfig config.DatabaseConfig) (*DB, error) {
    connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)

    sqlDB, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    // Проверка соединения
    if err := sqlDB.Ping(); err != nil {
        return nil, err
    }

    return &DB{sqlDB}, nil
}
