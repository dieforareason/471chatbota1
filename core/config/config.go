// Package config provides configuration management for the LLM bot
package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration settings for the application
type Config struct {
	// LLM Configuration
	GroqAPIKey     string
	ModelName      string
	SystemPrompt   string
	RequestTimeout time.Duration

	// Database Configuration
	DBPath          string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// LoadConfig loads configuration from environment variables with sensible defaults
func LoadConfig() *Config {
	return &Config{
		// LLM Config
		GroqAPIKey:     os.Getenv("GROQ_API_KEY"),
		ModelName:      getEnvOrDefault("MODEL_NAME", "llama3-8b-8192"),
		SystemPrompt:   getEnvOrDefault("SYSTEM_PROMPT", "Kamu adalah seorang personal assistent yang baik, manja dan supportive."),
		RequestTimeout: getDurationOrDefault("REQUEST_TIMEOUT", 30*time.Second),

		// Database Config
		DBPath:          getEnvOrDefault("DB_PATH", "D:/db/test.db"),
		MaxOpenConns:    getIntOrDefault("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getIntOrDefault("DB_MAX_IDLE_CONNS", 25),
		ConnMaxLifetime: getDurationOrDefault("DB_CONN_MAX_LIFETIME", 5*time.Minute),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
