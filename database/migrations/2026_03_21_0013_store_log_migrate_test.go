package migrations

import "testing"

func TestStoreLogMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreLogMigrate{}

	if migration.Signature() != "2026_03_21_0013_store_log_migrate" {
		t.Errorf("Expected signature '2026_03_21_0013_store_log_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run log store MigrateUp to create log tables" {
		t.Errorf("Expected description 'Run log store MigrateUp to create log tables', got '%s'", migration.Description())
	}
}

func TestStoreLogMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreLogMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreLogMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreLogMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
