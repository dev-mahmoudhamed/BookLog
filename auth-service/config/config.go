package config

import (
	"fmt"
	"log"
	"os"
)

type Config struct {
	AppPort    string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DBSSLMode  string
	JwtSecret  string
}

func LoadConfig() (*Config, error) {

	cfg := &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "authdb"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		JwtSecret:  getEnv("JWT_SECRET", ""), // no default, must set in env
	}

	if cfg.JwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	log.Println("✅ Configuration loaded successfully")
	return cfg, nil

}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
