package migrate_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// migrationsDir returns the absolute path to the migrations directory
// relative to this test file.
func migrationsDir(t *testing.T) string {
	t.Helper()
	// Walk up from this file's location to the project root.
	dir := filepath.Join("..", "..", "migrations")
	abs, err := filepath.Abs(dir)
	if err != nil {
		t.Fatalf("failed to resolve migrations path: %v", err)
	}
	return abs
}

func TestMigrationFiles_ExistAndNonEmpty(t *testing.T) {
	dir := migrationsDir(t)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("cannot read migrations dir %q: %v", dir, err)
	}

	var sqlFiles []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".sql") {
			sqlFiles = append(sqlFiles, e.Name())
		}
	}

	if len(sqlFiles) == 0 {
		t.Fatal("no SQL migration files found")
	}

	for _, name := range sqlFiles {
		path := filepath.Join(dir, name)
		data, err := os.ReadFile(path)
		if err != nil {
			t.Errorf("cannot read %s: %v", name, err)
			continue
		}
		if strings.TrimSpace(string(data)) == "" {
			t.Errorf("migration file is empty: %s", name)
		}
	}
}

func TestMigrationFiles_PairedUpDown(t *testing.T) {
	dir := migrationsDir(t)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("cannot read migrations dir %q: %v", dir, err)
	}

	ups := map[string]bool{}
	downs := map[string]bool{}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		switch {
		case strings.HasSuffix(name, ".up.sql"):
			prefix := strings.TrimSuffix(name, ".up.sql")
			ups[prefix] = true
		case strings.HasSuffix(name, ".down.sql"):
			prefix := strings.TrimSuffix(name, ".down.sql")
			downs[prefix] = true
		}
	}

	for prefix := range ups {
		if !downs[prefix] {
			t.Errorf("migration %q has .up.sql but no .down.sql", prefix)
		}
	}
	for prefix := range downs {
		if !ups[prefix] {
			t.Errorf("migration %q has .down.sql but no .up.sql", prefix)
		}
	}
}

func TestMigrationFiles_NamingConvention(t *testing.T) {
	dir := migrationsDir(t)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("cannot read migrations dir %q: %v", dir, err)
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".sql") {
			continue
		}
		// Expect format: 000NNN_<description>.(up|down).sql
		if !strings.HasSuffix(name, ".up.sql") && !strings.HasSuffix(name, ".down.sql") {
			t.Errorf("migration file %q does not follow *.up.sql or *.down.sql convention", name)
		}
	}
}
