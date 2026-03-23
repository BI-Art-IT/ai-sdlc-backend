package migrate

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Up applies all pending migrations from the given migrations directory.
func Up(dsn, migrationsPath string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("pgx5://%s", dsn),
	)
	if err != nil {
		return fmt.Errorf("migrate: failed to initialise: %w", err)
	}
	defer m.Close()

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate: up failed: %w", err)
	}
	return nil
}

// Down rolls back all applied migrations.
func Down(dsn, migrationsPath string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("pgx5://%s", dsn),
	)
	if err != nil {
		return fmt.Errorf("migrate: failed to initialise: %w", err)
	}
	defer m.Close()

	if err = m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate: down failed: %w", err)
	}
	return nil
}
