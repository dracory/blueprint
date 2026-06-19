package migrations

import "testing"

func TestStoreBlindIndexFirstNameMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreBlindIndexFirstNameMigrate{}

	if migration.Signature() != "2026_03_21_0004_store_blindindex_first_name_migrate" {
		t.Errorf("Expected signature '2026_03_21_0004_store_blindindex_first_name_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run blind index first name store MigrateUp to create blind index first name tables" {
		t.Errorf("Expected description 'Run blind index first name store MigrateUp to create blind index first name tables', got '%s'", migration.Description())
	}
}

func TestStoreBlindIndexFirstNameMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreBlindIndexFirstNameMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreBlindIndexFirstNameMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreBlindIndexFirstNameMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
