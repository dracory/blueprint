package migrations

import (
	"database/sql"
	"testing"
	"time"

	"github.com/dracory/migrate"
	_ "modernc.org/sqlite"
)

func TestUsersTableCreate(t *testing.T) {
	t.Skip("Enable if migrations are used")

	// Create in-memory SQLite database
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open in-memory database: %v", err)
	}
	defer db.Close()

	// Create migrator
	migrator := migrate.New(db, nil)

	// Builtin migrations will be added automatically on first Up() call

	// Add migration
	migration := &TableUsersCreate{}
	migrator.AddMigration(migration)

	// Test Up migration
	t.Run("Up", func(t *testing.T) {
		err := migrator.Up()
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
	})

	// Test Down migration
	t.Run("Down", func(t *testing.T) {
		err := migrator.Down()
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
	})

	// Test migration interface methods
	t.Run("Interface Methods", func(t *testing.T) {
		migration := &TableUsersCreate{}

		if migration.ID() != "2026_03_21_table_users_create" {
			t.Errorf("Expected ID '2026_03_21_table_users_create', got '%s'", migration.ID())
		}

		if migration.Description() != "Create users table with email and created_at indexes" {
			t.Errorf("Expected Description 'Create users table with email and created_at indexes', got '%s'", migration.Description())
		}

		// Test CreatedAt returns a valid time
		createdAt := migration.CreatedAt()
		if createdAt.IsZero() {
			t.Error("Expected CreatedAt to return a non-zero time")
		}

		// Verify it's the expected date
		expectedYear := 2026
		expectedMonth := time.Month(3)
		expectedDay := 21
		if createdAt.Year() != expectedYear || createdAt.Month() != expectedMonth || createdAt.Day() != expectedDay {
			t.Errorf("Expected CreatedAt to be %d-%02d-%02d, got %s",
				expectedYear, expectedMonth, expectedDay, createdAt.Format("2006-01-02"))
		}
	})
}
