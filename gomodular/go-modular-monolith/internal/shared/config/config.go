package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	Database string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	return &Config{
		Port:     getEnv("PORT", "8080"),
		Database: getEnv("DATABASE_URL", "localhost:5432"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}