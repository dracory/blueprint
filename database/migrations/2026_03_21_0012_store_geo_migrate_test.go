package migrations

import "testing"

func TestStoreGeoMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreGeoMigrate{}

	if migration.Signature() != "2026_03_21_0012_store_geo_migrate" {
		t.Errorf("Expected signature '2026_03_21_0012_store_geo_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run geo store MigrateUp to create geo tables" {
		t.Errorf("Expected description 'Run geo store MigrateUp to create geo tables', got '%s'", migration.Description())
	}
}

func TestStoreGeoMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreGeoMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreGeoMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreGeoMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
