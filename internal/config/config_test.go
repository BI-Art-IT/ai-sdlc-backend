package config_test

import (
	"os"
	"testing"

	"github.com/BI-Art-IT/ai-sdlc-backend/internal/config"
)

func TestLoad_RequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error when DATABASE_URL is not set, got nil")
	}
}

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/testdb")
	os.Unsetenv("MIGRATIONS_PATH")
	os.Unsetenv("SERVER_PORT")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.MigrationsPath != "migrations" {
		t.Errorf("expected default MigrationsPath 'migrations', got %q", cfg.MigrationsPath)
	}

	if cfg.ServerPort != "8080" {
		t.Errorf("expected default ServerPort '8080', got %q", cfg.ServerPort)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/mydb")
	t.Setenv("MIGRATIONS_PATH", "/app/migrations")
	t.Setenv("SERVER_PORT", "9090")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cfg.DatabaseURL != "postgres://localhost/mydb" {
		t.Errorf("unexpected DatabaseURL: %q", cfg.DatabaseURL)
	}
	if cfg.MigrationsPath != "/app/migrations" {
		t.Errorf("unexpected MigrationsPath: %q", cfg.MigrationsPath)
	}
	if cfg.ServerPort != "9090" {
		t.Errorf("unexpected ServerPort: %q", cfg.ServerPort)
	}
}
