package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	LogLevel    string
	DatabaseURL string
	RedisURL    string
	KafkaBrokers string
	JWTSecret   string
	Environment string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		Port:        getEnv("PORT", "8080"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/paymentdb?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379/0"),
		KafkaBrokers: getEnv("KAFKA_BROKERS", "localhost:9092"),
		JWTSecret:   getEnv("JWT_SECRET", "dev-secret-change-in-production"),
		Environment: getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
