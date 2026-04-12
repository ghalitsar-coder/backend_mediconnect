package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	AppEnv      string
	ServerPort  string
	DBURL       string
	RedisURL    string
	RabbitMQURL string
}

// LoadConfig reads configuration from the environment (or .env file as fallback).
func LoadConfig() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	return Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		DBURL:       getEnv("DB_URL", "postgres://mediconnect_user:mediconnect_password@localhost:5432/mediconnect_db?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
