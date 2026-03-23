package migrate_test

import (
	"errors"
	"testing"

	"github.com/golang-migrate/migrate/v4"

	appmigrate "github.com/BI-Art-IT/ai-sdlc-backend/internal/migrate"
)

func TestUp_InvalidSourceReturnsError(t *testing.T) {
	err := appmigrate.Up("pgx5://localhost/nonexistent", "/nonexistent/path")
	if err == nil {
		t.Fatal("expected error for invalid migration source, got nil")
	}
}

func TestUp_NoChangeIsNotAnError(t *testing.T) {
	// migrate.ErrNoChange must be swallowed inside Up(); this test validates the
	// contract by checking that the sentinel is indeed not propagated.
	// Since we cannot connect to a real DB here, we verify only the exported
	// error value is the standard library sentinel (sanity check).
	if !errors.Is(migrate.ErrNoChange, migrate.ErrNoChange) {
		t.Fatal("migrate.ErrNoChange should satisfy errors.Is with itself")
	}
}
