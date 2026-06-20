package migrations

import "testing"

func TestStoreBlindIndexLastNameMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreBlindIndexLastNameMigrate{}

	if migration.Signature() != "2026_03_21_0005_store_blindindex_last_name_migrate" {
		t.Errorf("Expected signature '2026_03_21_0005_store_blindindex_last_name_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run blind index last name store MigrateUp to create blind index last name tables" {
		t.Errorf("Expected description 'Run blind index last name store MigrateUp to create blind index last name tables', got '%s'", migration.Description())
	}
}

func TestStoreBlindIndexLastNameMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreBlindIndexLastNameMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreBlindIndexLastNameMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreBlindIndexLastNameMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
