package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	// Настройки сервера
	ServerHost string
	ServerPort string

	// Настройки базы данных
	DatabaseURL     string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration

	// Настройки окружения
	Environment string
	Debug       bool

	// Настройки логирования
	LogLevel string

	// Таймауты
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func Load() *Config {
	return &Config{
		ServerHost: getEnv("SERVER_HOST", "localhost"),
		ServerPort: getEnv("SERVER_PORT", "8080"),

		DatabaseURL:     getEnv("DATABASE_URL", ""),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 25),
		ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", time.Hour),

		Environment: getEnv("ENVIRONMENT", "development"),
		Debug:       getEnvAsBool("DEBUG", true),

		LogLevel: getEnv("LOG_LEVEL", "info"),

		ReadTimeout:     getEnvAsDuration("READ_TIMEOUT", 15*time.Second),
		WriteTimeout:    getEnvAsDuration("WRITE_TIMEOUT", 15*time.Second),
		ShutdownTimeout: getEnvAsDuration("SHUTDOWN_TIMEOUT", 10*time.Second),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}
