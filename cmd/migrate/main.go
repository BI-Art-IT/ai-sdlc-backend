// Command migrate is a standalone CLI for running database migrations.
//
// Usage:
//
//	migrate up            – apply all pending migrations
//	migrate down          – roll back all migrations
//	migrate down <N>      – roll back the last N migrations
//
// Configuration is read from environment variables (or a .env file in the
// working directory):
//
//	DATABASE_URL     – required; PostgreSQL connection URL
//	MIGRATIONS_PATH  – optional; path to the migrations directory (default: "migrations")
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/BI-Art-IT/ai-sdlc-backend/internal/config"
	"github.com/BI-Art-IT/ai-sdlc-backend/internal/migrate"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "up":
		log.Println("applying migrations…")
		if err := migrate.Up(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
			log.Fatalf("migrate up: %v", err)
		}
		log.Println("migrations applied successfully")

	case "down":
		steps := 0
		if len(os.Args) >= 3 {
			steps, err = strconv.Atoi(os.Args[2])
			if err != nil || steps < 1 {
				fmt.Fprintf(os.Stderr, "error: <N> must be a positive integer, got %q\n", os.Args[2])
				os.Exit(1)
			}
		}
		if steps > 0 {
			log.Printf("rolling back last %d migration(s)…", steps)
		} else {
			log.Println("rolling back all migrations…")
		}
		if err := migrate.Down(cfg.DatabaseURL, cfg.MigrationsPath, steps); err != nil {
			log.Fatalf("migrate down: %v", err)
		}
		log.Println("rollback completed successfully")

	default:
		fmt.Fprintf(os.Stderr, "error: unknown command %q\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `Usage:
  migrate up          apply all pending migrations
  migrate down        roll back all migrations
  migrate down <N>    roll back the last N migrations

Environment variables:
  DATABASE_URL      PostgreSQL connection URL (required)
  MIGRATIONS_PATH   path to migration files (default: "migrations")
`)
}
