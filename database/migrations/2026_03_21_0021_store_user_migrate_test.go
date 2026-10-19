package migrations

import "testing"

func TestStoreUserMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreUserMigrate{}

	if migration.Signature() != "2026_03_21_0021_store_user_migrate" {
		t.Errorf("Expected signature '2026_03_21_0021_store_user_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run user store MigrateUp to create user and role tables and seed default roles" {
		t.Errorf("Expected description 'Run user store MigrateUp to create user and role tables and seed default roles', got '%s'", migration.Description())
	}
}

func TestStoreUserMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreUserMigrate{}
	err := migration.Up()
	if err == nil {
		t.Fatal("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreUserMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreUserMigrate{}
	err := migration.Down()
	if err == nil {
		t.Fatal("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}
