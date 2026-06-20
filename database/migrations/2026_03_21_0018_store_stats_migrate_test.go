package migrations

import "testing"

func TestStoreStatsMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreStatsMigrate{}

	if migration.Signature() != "2026_03_21_0018_store_stats_migrate" {
		t.Errorf("Expected signature '2026_03_21_0018_store_stats_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run stats store MigrateUp to create stats tables" {
		t.Errorf("Expected description 'Run stats store MigrateUp to create stats tables', got '%s'", migration.Description())
	}
}

func TestStoreStatsMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreStatsMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreStatsMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreStatsMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
