package migrations

import "testing"

func TestStoreVaultMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreVaultMigrate{}

	if migration.Signature() != "2026_03_21_0022_store_vault_migrate" {
		t.Errorf("Expected signature '2026_03_21_0022_store_vault_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run vault store MigrateUp to create vault tables" {
		t.Errorf("Expected description 'Run vault store MigrateUp to create vault tables', got '%s'", migration.Description())
	}
}

func TestStoreVaultMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreVaultMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreVaultMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreVaultMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
