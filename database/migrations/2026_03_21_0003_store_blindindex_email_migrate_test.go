package migrations

import "testing"

func TestStoreBlindIndexEmailMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreBlindIndexEmailMigrate{}

	if migration.Signature() != "2026_03_21_0003_store_blindindex_email_migrate" {
		t.Errorf("Expected signature '2026_03_21_0003_store_blindindex_email_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run blind index email store MigrateUp to create blind index email tables" {
		t.Errorf("Expected description 'Run blind index email store MigrateUp to create blind index email tables', got '%s'", migration.Description())
	}
}

func TestStoreBlindIndexEmailMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreBlindIndexEmailMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreBlindIndexEmailMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreBlindIndexEmailMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
