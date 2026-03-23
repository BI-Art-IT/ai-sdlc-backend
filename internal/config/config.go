package config

import (
	"fmt"
	"os"
)

// Config holds all application configuration values loaded from environment variables.
type Config struct {
	DatabaseURL    string
	MigrationsPath string
	ServerPort     string
}

// Load reads configuration from environment variables and returns a Config.
// It returns an error if a required variable is missing.
func Load() (*Config, error) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		migrationsPath = "migrations"
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	return &Config{
		DatabaseURL:    dbURL,
		MigrationsPath: migrationsPath,
		ServerPort:     serverPort,
	}, nil
}
