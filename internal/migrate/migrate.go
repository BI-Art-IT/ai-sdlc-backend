package migrate

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// newMigrator creates a migrate.Migrate instance for the given source and database.
func newMigrator(databaseURL, migrationsPath string) (*migrate.Migrate, error) {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		databaseURL,
	)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}
	return m, nil
}

// Up applies all pending migrations from the given migrations directory against
// the database at databaseURL.
func Up(databaseURL, migrationsPath string) error {
	m, err := newMigrator(databaseURL, migrationsPath)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}

// Down rolls back the last n migrations. If n is 0 or negative all migrations
// are rolled back (full teardown).
func Down(databaseURL, migrationsPath string, steps int) error {
	m, err := newMigrator(databaseURL, migrationsPath)
	if err != nil {
		return err
	}
	defer m.Close()

	var rollbackErr error
	if steps > 0 {
		rollbackErr = m.Steps(-steps)
	} else {
		rollbackErr = m.Down()
	}

	if rollbackErr != nil && !errors.Is(rollbackErr, migrate.ErrNoChange) {
		return fmt.Errorf("rollback migrations: %w", rollbackErr)
	}

	return nil
}
