package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	LogLevel         string
	DatabaseURL      string
	RedisURL         string
	KafkaBrokers     string
	JWTSecret        string
	Environment      string
	EncryptionKey    string
	TLSCert          string
	TLSKey           string
	DBMaxConns       int
	DBMinConns       int
	DBMaxLifetime    time.Duration
	DBMaxIdleTime    time.Duration
}

func Load() *Config {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	godotenv.Load(fmt.Sprintf(".env.%s", env))
	godotenv.Load(".env.local")
	godotenv.Load()

	return &Config{
		Port:          getEnv("PORT", "8080"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/paymentdb?sslmode=disable"),
		RedisURL:      getEnv("REDIS_URL", "redis://localhost:6379/0"),
		KafkaBrokers:  getEnv("KAFKA_BROKERS", "localhost:9092"),
		JWTSecret:     getEnv("JWT_SECRET", ""),
		Environment:   getEnv("ENVIRONMENT", "development"),
		EncryptionKey: getEnv("ENCRYPTION_KEY", ""),
		TLSCert:       getEnv("TLS_CERT", ""),
		TLSKey:        getEnv("TLS_KEY", ""),
		DBMaxConns:    getEnvInt("DATABASE_MAX_CONNS", 25),
		DBMinConns:    getEnvInt("DATABASE_MIN_CONNS", 5),
		DBMaxLifetime: getEnvDuration("DATABASE_MAX_LIFETIME", 30*time.Minute),
		DBMaxIdleTime: getEnvDuration("DATABASE_MAX_IDLE_TIME", 5*time.Minute),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return fallback
}
