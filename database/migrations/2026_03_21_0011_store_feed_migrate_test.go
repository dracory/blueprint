package migrations

import "testing"

func TestStoreFeedMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreFeedMigrate{}

	if migration.Signature() != "2026_03_21_0011_store_feed_migrate" {
		t.Errorf("Expected signature '2026_03_21_0011_store_feed_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run feed store MigrateUp to create feed tables" {
		t.Errorf("Expected description 'Run feed store MigrateUp to create feed tables', got '%s'", migration.Description())
	}
}

func TestStoreFeedMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreFeedMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreFeedMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreFeedMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
