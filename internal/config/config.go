package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the service.
type Config struct {
	HTTPPort   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

// Load reads configuration from environment variables with sensible defaults
// for local development.
func Load() *Config {
	return &Config{
		HTTPPort:   envOrDefault("HTTP_PORT", "8080"),
		DBHost:     envOrDefault("DB_HOST", "localhost"),
		DBPort:     envOrDefault("DB_PORT", "5432"),
		DBUser:     envOrDefault("DB_USER", "postgres"),
		DBPassword: envOrDefault("DB_PASSWORD", "postgres"),
		DBName:     envOrDefault("DB_NAME", "familydb"),
		DBSSLMode:  envOrDefault("DB_SSLMODE", "disable"),
	}
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
