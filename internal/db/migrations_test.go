package db_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const migrationsDir = "../../migrations"

// expectedUpFiles lists all expected migration up files in order.
var expectedUpFiles = []struct {
	filename string
	table    string
	columns  []string
}{
	{
		filename: "000001_create_airports_table.up.sql",
		table:    "airports",
		columns:  []string{"id", "iata_code", "name", "city", "country", "country_code", "timezone", "latitude", "longitude", "active", "created_at", "updated_at"},
	},
	{
		filename: "000002_create_airlines_table.up.sql",
		table:    "airlines",
		columns:  []string{"id", "iata_code", "name", "logo_url", "active", "created_at", "updated_at"},
	},
	{
		filename: "000003_create_fare_classes_table.up.sql",
		table:    "fare_classes",
		columns:  []string{"id", "code", "name", "sort_order"},
	},
	{
		filename: "000004_create_flights_table.up.sql",
		table:    "flights",
		columns:  []string{"id", "flight_number", "airline_id", "origin_airport_id", "destination_airport_id", "departure_time", "arrival_time", "duration_minutes", "stops", "aircraft_type", "total_seat_capacity", "status", "estimated_departure_time", "estimated_arrival_time", "departure_gate", "departure_terminal", "arrival_gate", "arrival_terminal", "created_at", "updated_at"},
	},
	{
		filename: "000005_create_fares_table.up.sql",
		table:    "fares",
		columns:  []string{"id", "flight_id", "fare_class_id", "base_fare", "taxes", "fees", "total_price", "currency", "available_seats", "valid_from", "valid_until", "created_at", "updated_at"},
	},
	{
		filename: "000006_create_search_logs_table.up.sql",
		table:    "search_logs",
		columns:  []string{"id", "user_id_hash", "session_id", "origin_code", "destination_code", "departure_date", "return_date", "trip_type", "adults", "children", "infants", "cabin_class", "flexible_dates", "results_returned", "cache_hit", "response_time_ms", "created_at"},
	},
}

// expectedDownFiles lists all expected migration down files.
var expectedDownFiles = []string{
	"000001_create_airports_table.down.sql",
	"000002_create_airlines_table.down.sql",
	"000003_create_fare_classes_table.down.sql",
	"000004_create_flights_table.down.sql",
	"000005_create_fares_table.down.sql",
	"000006_create_search_logs_table.down.sql",
}

func TestMigrationUpFilesExist(t *testing.T) {
	for _, tc := range expectedUpFiles {
		t.Run(tc.filename, func(t *testing.T) {
			path := filepath.Join(migrationsDir, tc.filename)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("migration file %q does not exist", path)
			}
		})
	}
}

func TestMigrationDownFilesExist(t *testing.T) {
	for _, filename := range expectedDownFiles {
		t.Run(filename, func(t *testing.T) {
			path := filepath.Join(migrationsDir, filename)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("migration file %q does not exist", path)
			}
		})
	}
}

func TestMigrationUpFilesContainCreateTable(t *testing.T) {
	for _, tc := range expectedUpFiles {
		t.Run(tc.filename, func(t *testing.T) {
			path := filepath.Join(migrationsDir, tc.filename)
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read migration file %q: %v", path, err)
			}

			sql := strings.ToUpper(string(content))
			expectedTable := strings.ToUpper(tc.table)

			if !strings.Contains(sql, "CREATE TABLE") {
				t.Errorf("migration %q does not contain CREATE TABLE statement", tc.filename)
			}

			if !strings.Contains(sql, expectedTable) {
				t.Errorf("migration %q does not reference table %q", tc.filename, tc.table)
			}
		})
	}
}

func TestMigrationUpFilesContainExpectedColumns(t *testing.T) {
	for _, tc := range expectedUpFiles {
		t.Run(tc.filename, func(t *testing.T) {
			path := filepath.Join(migrationsDir, tc.filename)
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read migration file %q: %v", path, err)
			}

			sql := string(content)
			for _, col := range tc.columns {
				if !strings.Contains(sql, col) {
					t.Errorf("migration %q is missing column %q", tc.filename, col)
				}
			}
		})
	}
}

func TestMigrationDownFilesContainDropTable(t *testing.T) {
	for _, filename := range expectedDownFiles {
		t.Run(filename, func(t *testing.T) {
			path := filepath.Join(migrationsDir, filename)
			content, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("failed to read migration file %q: %v", path, err)
			}

			sql := strings.ToUpper(string(content))
			if !strings.Contains(sql, "DROP TABLE") {
				t.Errorf("down migration %q does not contain DROP TABLE statement", filename)
			}
		})
	}
}

func TestFareClassesMigrationContainsSeedData(t *testing.T) {
	path := filepath.Join(migrationsDir, "000003_create_fare_classes_table.up.sql")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fare_classes migration: %v", err)
	}

	sql := string(content)
	expectedCodes := []string{"economy", "premium-economy", "business", "first"}
	for _, code := range expectedCodes {
		if !strings.Contains(sql, code) {
			t.Errorf("fare_classes migration is missing seed data for class %q", code)
		}
	}

	if !strings.Contains(strings.ToUpper(sql), "INSERT INTO") {
		t.Errorf("fare_classes migration does not contain INSERT INTO for seed data")
	}
}

func TestFlightsMigrationContainsForeignKeys(t *testing.T) {
	path := filepath.Join(migrationsDir, "000004_create_flights_table.up.sql")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read flights migration: %v", err)
	}

	sql := strings.ToUpper(string(content))
	expectedFKRefs := []string{"AIRLINES", "AIRPORTS"}
	for _, ref := range expectedFKRefs {
		if !strings.Contains(sql, ref) {
			t.Errorf("flights migration is missing FK reference to %q", ref)
		}
	}
	if !strings.Contains(sql, "FOREIGN KEY") {
		t.Errorf("flights migration does not define FOREIGN KEY constraints")
	}
}

func TestFlightsMigrationContainsStatusCheck(t *testing.T) {
	path := filepath.Join(migrationsDir, "000004_create_flights_table.up.sql")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read flights migration: %v", err)
	}

	sql := string(content)
	expectedStatuses := []string{"scheduled", "boarding", "departed", "en_route", "landed", "arrived", "delayed", "cancelled"}
	for _, status := range expectedStatuses {
		if !strings.Contains(sql, status) {
			t.Errorf("flights migration is missing status value %q in CHECK constraint", status)
		}
	}
}

func TestFaresMigrationContainsForeignKeys(t *testing.T) {
	path := filepath.Join(migrationsDir, "000005_create_fares_table.up.sql")
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fares migration: %v", err)
	}

	sql := strings.ToUpper(string(content))
	expectedFKRefs := []string{"FLIGHTS", "FARE_CLASSES"}
	for _, ref := range expectedFKRefs {
		if !strings.Contains(sql, ref) {
			t.Errorf("fares migration is missing FK reference to %q", ref)
		}
	}
	if !strings.Contains(sql, "FOREIGN KEY") {
		t.Errorf("fares migration does not define FOREIGN KEY constraints")
	}
}

func TestMigrationFilesAreNotEmpty(t *testing.T) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("failed to read migrations directory: %v", err)
	}

	if len(entries) == 0 {
		t.Fatal("migrations directory is empty")
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		path := filepath.Join(migrationsDir, entry.Name())
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("failed to stat %q: %v", path, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("migration file %q is empty", entry.Name())
		}
	}
}
