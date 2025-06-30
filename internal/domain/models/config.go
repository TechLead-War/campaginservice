package models

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	DBHOST    string
	DBPORT    string
	DBUSER    string
	DBPass    string
	DBName    string
	AppPort   string
	RedisAddr string
	RedisPass string
	RedisDB   int
	CacheSize int
	LogLevel  string
}

// LoadConfig loads the application configuration from environment variables or a config file.
func LoadConfig() (*AppConfig, error) {
	// Load .env file if it exists
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: .env file not found, using environment variables")
	}

	cfg := &AppConfig{
		DBHOST:    getEnv("DB_HOST", "localhost"),
		DBPORT:    getEnv("DB_PORT", "5432"),
		DBUSER:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASS", "password"),
		DBName:    getEnv("DB_NAME", "campaign_service"),
		AppPort:   getEnv("APP_PORT", "8080"),
		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPass: getEnv("REDIS_PASS", ""),
		RedisDB:   getEnvAsInt("REDIS_DB", 0),
		CacheSize: getEnvAsInt("CACHE_SIZE", 1000),
		LogLevel:  getEnv("LOG_LEVEL", "info"),
	}

	// Validate configuration
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

func (cfg *AppConfig) validate() error {
	if cfg.DBHOST == "" {
		return fmt.Errorf("DB_HOST cannot be empty")
	}
	if cfg.DBUSER == "" {
		return fmt.Errorf("DB_USER cannot be empty")
	}
	if cfg.DBName == "" {
		return fmt.Errorf("DB_NAME cannot be empty")
	}
	if cfg.AppPort == "" {
		return fmt.Errorf("APP_PORT cannot be empty")
	}

	// Validate port numbers
	if _, err := strconv.Atoi(cfg.DBPORT); err != nil {
		return fmt.Errorf("DB_PORT must be a valid integer: %s", cfg.DBPORT)
	}
	if _, err := strconv.Atoi(cfg.AppPort); err != nil {
		return fmt.Errorf("APP_PORT must be a valid integer: %s", cfg.AppPort)
	}

	// Validate cache size
	if cfg.CacheSize <= 0 {
		return fmt.Errorf("CACHE_SIZE must be greater than 0: %d", cfg.CacheSize)
	}

	// Validate log level
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[strings.ToLower(cfg.LogLevel)] {
		return fmt.Errorf("LOG_LEVEL must be one of: debug, info, warn, error")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.TrimSpace(value)
	}

	log.Printf("Environment variable %s not set, using default: %s", key, defaultValue)
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		log.Printf("Environment variable %s not set, using default: %d", key, defaultValue)
		return defaultValue
	}

	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	log.Printf("Warning: Environment variable %s (value: %s) is not a valid integer, using default: %d", key, valueStr, defaultValue)
	return defaultValue
}
