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

	BrevoApiKey             string
	BrevoSenderEmail        string
	BrevoSenderName         string
	FrontendResetUrl        string
	ResetTokenExpireMinutes int
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

	resetExpire, err := strconv.Atoi(getEnv("RESET_TOKEN_EXPIRE_MINUTES", "30"))
	if err != nil {
		resetExpire = 30
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

		JWTAccessSecret:         getEnv("JWT_ACCESS_SECRET", ""),
		JWTAccessExpireMinutes:  accessExpire,
		JWTRefreshSecret:        getEnv("JWT_REFRESH_SECRET", ""),
		JWTRefreshExpireDays:    refreshExpire,
		BrevoApiKey:             getEnv("BREVO_API_KEY", ""),
		BrevoSenderEmail:        getEnv("BREVO_SENDER_EMAIL", ""),
		BrevoSenderName:         getEnv("BREVO_SENDER_NAME", "Sembako App"),
		FrontendResetUrl:        getEnv("FRONTEND_RESET_URL", "http://localhost:3000/reset-password"),
		ResetTokenExpireMinutes: resetExpire,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
