package migrations

import "testing"

func TestStoreShopMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreShopMigrate{}

	if migration.Signature() != "2026_03_21_0017_store_shop_migrate" {
		t.Errorf("Expected signature '2026_03_21_0017_store_shop_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run shop store MigrateUp to create shop tables" {
		t.Errorf("Expected description 'Run shop store MigrateUp to create shop tables', got '%s'", migration.Description())
	}
}

func TestStoreShopMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreShopMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreShopMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreShopMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
