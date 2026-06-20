package migrations

import "testing"

func TestStoreUserMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreUserMigrate{}

	if migration.Signature() != "2026_03_21_0021_store_user_migrate" {
		t.Errorf("Expected signature '2026_03_21_0021_store_user_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run user store MigrateUp to create user tables" {
		t.Errorf("Expected description 'Run user store MigrateUp to create user tables', got '%s'", migration.Description())
	}
}

func TestStoreUserMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreUserMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreUserMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreUserMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
