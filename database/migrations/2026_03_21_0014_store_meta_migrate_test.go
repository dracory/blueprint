package migrations

import "testing"

func TestStoreMetaMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreMetaMigrate{}

	if migration.Signature() != "2026_03_21_0014_store_meta_migrate" {
		t.Errorf("Expected signature '2026_03_21_0014_store_meta_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run meta store MigrateUp to create meta tables" {
		t.Errorf("Expected description 'Run meta store MigrateUp to create meta tables', got '%s'", migration.Description())
	}
}

func TestStoreMetaMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreMetaMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreMetaMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreMetaMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
