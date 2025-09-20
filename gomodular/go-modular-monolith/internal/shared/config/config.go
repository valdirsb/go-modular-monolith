package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	DBHost      string
	DBPort      string
	DBUsername  string
	DBPassword  string
	DBDatabase  string
	JWTSecret   string
	Environment string
}

func LoadConfig() (*Config, error) {
	// Tentar carregar arquivo .env (opcional)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables or defaults")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBUsername:  getEnv("DB_USERNAME", "root"),
		DBPassword:  getEnv("DB_PASSWORD", "123456"),
		DBDatabase:  getEnv("DB_DATABASE", "app_db"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
