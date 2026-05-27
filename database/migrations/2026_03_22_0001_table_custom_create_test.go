package migrations

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dracory/migrate"
	_ "modernc.org/sqlite"
)

func TestCustomTableCreate_Up(t *testing.T) {
	t.Skip("Enable if migrations are used")

	// Create in-memory SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	// Create migrator
	migrator, err := migrate.New(db, nil)
	if err != nil {
		t.Fatalf("Failed to create migrator: %v", err)
	}

	// Add migration
	migration := &TableCustomCreate{}
	migrator.AddMigration(migration)

	err = migrator.Up(context.Background())
	if err != nil {
		t.Fatalf("Migration Up failed: %v", err)
	}

	// Verify table exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to check table existence: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected users table to exist, but got count: %d", count)
	}

	// Verify columns exist
	rows, err := db.Query("PRAGMA table_info(users)")
	if err != nil {
		t.Fatalf("Failed to get table info: %v", err)
	}
	defer rows.Close()

	expectedColumns := map[string]bool{
		"id":            false,
		"email":         false,
		"name":          false,
		"password_hash": false,
		"created_at":    false,
		"updated_at":    false,
	}

	for rows.Next() {
		var cid interface{}
		var name, datatype string
		var notnull, pk interface{}
		var dflt_value interface{}

		err = rows.Scan(&cid, &name, &datatype, &notnull, &dflt_value, &pk)
		if err != nil {
			t.Fatalf("Failed to scan column info: %v", err)
		}

		if _, exists := expectedColumns[name]; exists {
			expectedColumns[name] = true
		}
	}

	for column, found := range expectedColumns {
		if !found {
			t.Errorf("Expected column %s not found", column)
		}
	}

	// Verify indexes exist
	indexRows, err := db.Query("SELECT name FROM sqlite_master WHERE type='index' AND name LIKE 'idx_users_%'")
	if err != nil {
		t.Fatalf("Failed to get index info: %v", err)
	}
	defer indexRows.Close()

	expectedIndexes := map[string]bool{
		"idx_users_email":      false,
		"idx_users_created_at": false,
	}

	for indexRows.Next() {
		var indexName string
		err = indexRows.Scan(&indexName)
		if err != nil {
			t.Fatalf("Failed to scan index info: %v", err)
		}

		if _, exists := expectedIndexes[indexName]; exists {
			expectedIndexes[indexName] = true
		}
	}

	for index, found := range expectedIndexes {
		if !found {
			t.Errorf("Expected index %s not found", index)
		}
	}
}

func TestCustomTableCreate_Down(t *testing.T) {
	t.Skip("Enable if migrations are used")

	// Create in-memory SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	// Create migrator
	migrator, err := migrate.New(db, nil)
	if err != nil {
		t.Fatalf("Failed to create migrator: %v", err)
	}

	// Add migration
	migration := &TableCustomCreate{}
	migrator.AddMigration(migration)

	// First run Up to create the table
	err = migrator.Up(context.Background())
	if err != nil {
		t.Fatalf("Migration Up failed: %v", err)
	}

	// Test Down migration
	err = migrator.Down(context.Background())
	if err != nil {
		t.Fatalf("Migration Down failed: %v", err)
	}

	// Verify table doesn't exist
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='users'").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to check table existence: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected users table to be dropped, but got count: %d", count)
	}
}

func TestCustomTableCreate_InterfaceMethods(t *testing.T) {
	t.Skip("Enable if migrations are used")

	migration := &TableCustomCreate{}

	if migration.ID() != "2026_03_22_0001_table_custom_create" {
		t.Errorf("Expected ID '2026_03_22_0001_table_custom_create', got '%s'", migration.ID())
	}

	if migration.Description() != "Example: Create custom table with indexes" {
		t.Errorf("Expected Description 'Example: Create custom table with indexes', got '%s'", migration.Description())
	}

	// Test CreatedAt returns a valid time
	createdAt := migration.CreatedAt()
	if createdAt.IsZero() {
		t.Error("Expected CreatedAt to return a non-zero time")
	}

	// Verify it's the expected date
	expectedYear := 2026
	expectedMonth := time.Month(3)
	expectedDay := 22
	if createdAt.Year() != expectedYear || createdAt.Month() != expectedMonth || createdAt.Day() != expectedDay {
		t.Errorf("Expected CreatedAt to be %d-%02d-%02d, got %s",
			expectedYear, expectedMonth, expectedDay, createdAt.Format("2006-01-02"))
	}
}
