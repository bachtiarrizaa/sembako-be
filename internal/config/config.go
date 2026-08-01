package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string
	AppEnv  string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTAccessSecret        string
	JWTAccessExpireMinutes int
	JWTRefreshSecret       string
	JWTRefreshExpireDays   int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found, reading from system env")
	}

	accessExpire, err := strconv.Atoi(getEnv("JWT_ACCESS_EXPIRE_MINUTES", "15"))
	if err != nil {
		accessExpire = 15
	}

	refreshExpire, err := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRE_DAYS", "7"))
	if err != nil {
		refreshExpire = 7
	}

	return &Config{
		AppPort: getEnv("APP_PORT", "8080"),
		AppEnv:  getEnv("APP_ENV", "development"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "sembako_pos"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		JWTAccessSecret:        getEnv("JWT_ACCESS_SECRET", ""),
		JWTAccessExpireMinutes: accessExpire,
		JWTRefreshSecret:       getEnv("JWT_REFRESH_SECRET", ""),
		JWTRefreshExpireDays:   refreshExpire,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
