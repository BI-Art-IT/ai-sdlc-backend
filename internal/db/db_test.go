package db_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/BI-Art-IT/ai-sdlc-backend/internal/db"
)

func TestConfigDSN(t *testing.T) {
	cfg := db.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "secret",
		DBName:   "airlines",
		SSLMode:  "disable",
	}

	expected := "host=localhost port=5432 user=postgres password=secret dbname=airlines sslmode=disable"
	got := cfg.DSN()

	if got != expected {
		t.Errorf("DSN() = %q; want %q", got, expected)
	}
}

func TestConfigDSN_AllFields(t *testing.T) {
	tests := []struct {
		name string
		cfg  db.Config
		want string
	}{
		{
			name: "ssl enabled",
			cfg: db.Config{
				Host:     "db.example.com",
				Port:     5433,
				User:     "admin",
				Password: "p@ssw0rd",
				DBName:   "airlines_prod",
				SSLMode:  "require",
			},
			want: "host=db.example.com port=5433 user=admin password=p@ssw0rd dbname=airlines_prod sslmode=require",
		},
		{
			name: "default local config",
			cfg: db.Config{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "",
				DBName:   "airlines",
				SSLMode:  "disable",
			},
			want: "host=localhost port=5432 user=postgres password= dbname=airlines sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.DSN()
			if got != tt.want {
				t.Errorf("DSN() = %q; want %q", got, tt.want)
			}
		})
	}
}

func TestConfigDSN_Format(t *testing.T) {
	cfg := db.Config{
		Host:     "myhost",
		Port:     5432,
		User:     "myuser",
		Password: "mypassword",
		DBName:   "mydb",
		SSLMode:  "disable",
	}

	dsn := cfg.DSN()

	checks := []string{
		fmt.Sprintf("host=%s", cfg.Host),
		fmt.Sprintf("port=%d", cfg.Port),
		fmt.Sprintf("user=%s", cfg.User),
		fmt.Sprintf("password=%s", cfg.Password),
		fmt.Sprintf("dbname=%s", cfg.DBName),
		fmt.Sprintf("sslmode=%s", cfg.SSLMode),
	}

	for _, check := range checks {
		if !strings.Contains(dsn, check) {
			t.Errorf("DSN() %q does not contain %q", dsn, check)
		}
	}
}
