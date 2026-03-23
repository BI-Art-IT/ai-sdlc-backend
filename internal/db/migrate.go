package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateUp applies all pending up migrations from the given source directory.
// migrationsPath should be a file URL, e.g. "file://migrations".
func MigrateUp(dsn, migrationsPath string) error {
	m, err := migrate.New(migrationsPath, "pgx5://"+dsn)
	if err != nil {
		return fmt.Errorf("migrate: failed to initialise: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate: failed to apply migrations: %w", err)
	}

	return nil
}

// MigrateDown rolls back all applied migrations.
func MigrateDown(dsn, migrationsPath string) error {
	m, err := migrate.New(migrationsPath, "pgx5://"+dsn)
	if err != nil {
		return fmt.Errorf("migrate: failed to initialise: %w", err)
	}
	defer m.Close()

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate: failed to roll back migrations: %w", err)
	}

	return nil
}
