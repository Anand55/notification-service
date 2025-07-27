package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL     string
	EmailHost       string
	EmailPort       int
	EmailUsername   string
	EmailPassword   string
	SlackToken      string
	SlackChannel    string
	JWTSecret       string
	Environment     string
}

func Load() *Config {
	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "notifications.db"),
		EmailHost:     getEnv("EMAIL_HOST", "smtp.gmail.com"),
		EmailPort:     getEnvAsInt("EMAIL_PORT", 587),
		EmailUsername: getEnv("EMAIL_USERNAME", ""),
		EmailPassword: getEnv("EMAIL_PASSWORD", ""),
		SlackToken:    getEnv("SLACK_TOKEN", ""),
		SlackChannel:  getEnv("SLACK_CHANNEL", "#general"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		Environment:   getEnv("ENVIRONMENT", "development"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
} 