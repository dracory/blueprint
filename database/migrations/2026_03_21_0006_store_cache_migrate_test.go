package migrations

import "testing"

func TestStoreCacheMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreCacheMigrate{}

	if migration.Signature() != "2026_03_21_0006_store_cache_migrate" {
		t.Errorf("Expected signature '2026_03_21_0006_store_cache_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run cache store MigrateUp to create cache tables" {
		t.Errorf("Expected description 'Run cache store MigrateUp to create cache tables', got '%s'", migration.Description())
	}
}

func TestStoreCacheMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreCacheMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreCacheMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreCacheMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
