package config

import "os"

type Config struct {
	ServerAddress string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
}

func Load() Config {
	return Config{
		ServerAddress: getEnv("BOOK_SERVICE_ADDR", ":8081"),
		DBHost:        getEnv("DB_HOST", "book-db"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "Password_123"),
		DBName:        getEnv("DB_NAME", "bookdb"),
		JWTSecret:     getEnv("JWT_SECRET", "your_jwt_secret"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
