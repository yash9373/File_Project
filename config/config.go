package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// AppConfig holds environment config
type AppConfig struct {
	AppEnv              string
	AppPort             string
	JWTSecret           string
	TokenExpiresInHours int
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	SSLMode             string
}

var C AppConfig

// Load loads .env into AppConfig
func Load() {
	_ = godotenv.Load() // load .env if present

	C = AppConfig{
		AppEnv:              getEnv("APP_ENV", "development"),
		AppPort:             getEnv("APP_PORT", "8080"),
		JWTSecret:           getEnv("JWT_SECRET", "change_me"),
		TokenExpiresInHours: getEnvAsInt("TOKEN_EXPIRES_IN_HOURS", 24),
		DBHost:              getEnv("DB_HOST", "localhost"),
		DBPort:              getEnv("DB_PORT", "5432"),
		DBUser:              getEnv("DB_USER", "postgres"),
		DBPassword:          getEnv("DB_PASSWORD", "postgres"),
		DBName:              getEnv("DB_NAME", "fiber_auth"),
		SSLMode:             getEnv("SSL_MODE", "disable"),
	}

	log.Printf("config loaded: env=%s port=%s db=%s@%s:%s/%s", C.AppEnv, C.AppPort, C.DBUser, C.DBHost, C.DBPort, C.DBName)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
