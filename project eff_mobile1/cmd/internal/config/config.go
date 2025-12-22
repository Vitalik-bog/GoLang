package config

import (
    "fmt"
    "os"
    "strconv"
    "gopkg.in/yaml.v3"
)

type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
}

type ServerConfig struct {
    Port string `yaml:"port"`
    Host string `yaml:"host"`
    Mode string `yaml:"mode"`
}

type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    DBName   string `yaml:"dbname"`
    SSLMode  string `yaml:"sslmode"`
}

func LoadConfig(configPath string) (*Config, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("error reading config file: %w", err)
    }
    
    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("error parsing config file: %w", err)
    }
    
    // Заполняем из переменных окружения если они есть
    if port := os.Getenv("SERVER_PORT"); port != "" {
        config.Server.Port = port
    }
    
    if host := os.Getenv("SERVER_HOST"); host != "" {
        config.Server.Host = host
    }
    
    if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
        config.Database.Host = dbHost
    }
    
    if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
        if port, err := strconv.Atoi(dbPort); err == nil {
            config.Database.Port = port
        }
    }
    
    if dbUser := os.Getenv("DB_USER"); dbUser != "" {
        config.Database.User = dbUser
    }
    
    if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
        config.Database.Password = dbPass
    }
    
    if dbName := os.Getenv("DB_NAME"); dbName != "" {
        config.Database.DBName = dbName
    }
    
    return &config, nil
}
