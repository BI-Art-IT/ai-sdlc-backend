package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Ensure no env vars interfere.
	for _, k := range []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE"} {
		t.Setenv(k, "")
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.DB.Host != "localhost" {
		t.Errorf("DB.Host = %q, want %q", cfg.DB.Host, "localhost")
	}
	if cfg.DB.Port != "5432" {
		t.Errorf("DB.Port = %q, want %q", cfg.DB.Port, "5432")
	}
	if cfg.DB.User != "postgres" {
		t.Errorf("DB.User = %q, want %q", cfg.DB.User, "postgres")
	}
	if cfg.DB.Name != "airlines" {
		t.Errorf("DB.Name = %q, want %q", cfg.DB.Name, "airlines")
	}
	if cfg.DB.SSLMode != "disable" {
		t.Errorf("DB.SSLMode = %q, want %q", cfg.DB.SSLMode, "disable")
	}
}

func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USER", "admin")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "testdb")
	t.Setenv("DB_SSLMODE", "require")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.DB.Host != "db.example.com" {
		t.Errorf("DB.Host = %q, want %q", cfg.DB.Host, "db.example.com")
	}
	if cfg.DB.Port != "5433" {
		t.Errorf("DB.Port = %q, want %q", cfg.DB.Port, "5433")
	}
	if cfg.DB.User != "admin" {
		t.Errorf("DB.User = %q, want %q", cfg.DB.User, "admin")
	}
	if cfg.DB.Password != "secret" {
		t.Errorf("DB.Password = %q, want %q", cfg.DB.Password, "secret")
	}
	if cfg.DB.Name != "testdb" {
		t.Errorf("DB.Name = %q, want %q", cfg.DB.Name, "testdb")
	}
	if cfg.DB.SSLMode != "require" {
		t.Errorf("DB.SSLMode = %q, want %q", cfg.DB.SSLMode, "require")
	}
}

func TestDBConfig_DSN(t *testing.T) {
	db := DBConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "pass",
		Name:     "airlines",
		SSLMode:  "disable",
	}

	want := "host=localhost port=5432 user=postgres password=pass dbname=airlines sslmode=disable"
	got := db.DSN()
	if got != want {
		t.Errorf("DSN() = %q, want %q", got, want)
	}
}

func TestGetEnv_FallbackWhenEmpty(t *testing.T) {
	os.Unsetenv("TEST_KEY_NONEXISTENT")
	got := getEnv("TEST_KEY_NONEXISTENT", "default_val")
	if got != "default_val" {
		t.Errorf("getEnv() = %q, want %q", got, "default_val")
	}
}

func TestGetEnv_ReturnsValue(t *testing.T) {
	t.Setenv("TEST_KEY_EXISTS", "my_value")
	got := getEnv("TEST_KEY_EXISTS", "default_val")
	if got != "my_value" {
		t.Errorf("getEnv() = %q, want %q", got, "my_value")
	}
}
