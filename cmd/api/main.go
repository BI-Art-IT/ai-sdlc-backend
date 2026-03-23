package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BI-Art-IT/ai-sdlc-backend/internal/db"
)

func main() {
	cfg := db.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "airlines"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	ctx := context.Background()

	pool, err := db.Connect(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	migrationsPath := getEnv("MIGRATIONS_PATH", "file://migrations")
	if err := db.MigrateUp(cfg.DSN(), migrationsPath); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	fmt.Println("database migrations applied successfully")
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
