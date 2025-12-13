package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server    ServerConfig
	WebSocket WebSocketConfig
	App       AppConfig
}

type ServerConfig struct {
	Address string
	Env     string
}

type WebSocketConfig struct {
	ReadBufferSize  int
	WriteBufferSize int
	MaxMessageSize  int64
	PingPeriod      int
}

type AppConfig struct {
	MaxRooms          int
	MaxClientsPerRoom int
	RoomTTL           int
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Address: getEnv("SERVER_ADDRESS", ":8080"),
			Env:     getEnv("ENVIRONMENT", "development"),
		},
		WebSocket: WebSocketConfig{
			ReadBufferSize:  getEnvAsInt("WS_READ_BUFFER_SIZE", 1024),
			WriteBufferSize: getEnvAsInt("WS_WRITE_BUFFER_SIZE", 1024),
			MaxMessageSize:  getEnvAsInt64("WS_MAX_MESSAGE_SIZE", 512),
			PingPeriod:      getEnvAsInt("WS_PING_PERIOD", 60),
		},
		App: AppConfig{
			MaxRooms:          getEnvAsInt("MAX_ROOMS", 100),
			MaxClientsPerRoom: getEnvAsInt("MAX_CLIENTS_PER_ROOM", 50),
			RoomTTL:           getEnvAsInt("ROOM_TTL", 3600),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}
