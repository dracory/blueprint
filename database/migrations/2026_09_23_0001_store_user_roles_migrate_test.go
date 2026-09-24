package migrations

import "testing"

func TestStoreUserRolesMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreUserRolesMigrate{}

	if migration.Signature() != "2026_09_23_0001_store_user_roles_migrate" {
		t.Errorf("Expected signature '2026_09_23_0001_store_user_roles_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run user store MigrateUp to create role tables and seed default roles" {
		t.Errorf("Expected description 'Run user store MigrateUp to create role tables and seed default roles', got '%s'", migration.Description())
	}
}

func TestStoreUserRolesMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreUserRolesMigrate{}
	err := migration.Up()
	if err == nil {
		t.Fatal("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreUserRolesMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreUserRolesMigrate{}
	err := migration.Down()
	if err == nil {
		t.Fatal("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}
