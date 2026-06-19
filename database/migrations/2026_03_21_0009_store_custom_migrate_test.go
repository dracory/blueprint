package migrations

import "testing"

func TestStoreCustomMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreCustomMigrate{}

	if migration.Signature() != "2026_03_21_0009_store_custom_migrate" {
		t.Errorf("Expected signature '2026_03_21_0009_store_custom_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run custom store MigrateUp to create custom tables" {
		t.Errorf("Expected description 'Run custom store MigrateUp to create custom tables', got '%s'", migration.Description())
	}
}

func TestStoreCustomMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreCustomMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreCustomMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreCustomMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
