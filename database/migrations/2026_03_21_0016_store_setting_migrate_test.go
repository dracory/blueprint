package migrations

import "testing"

func TestStoreSettingMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreSettingMigrate{}

	if migration.Signature() != "2026_03_21_0016_store_setting_migrate" {
		t.Errorf("Expected signature '2026_03_21_0016_store_setting_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run setting store MigrateUp to create setting tables" {
		t.Errorf("Expected description 'Run setting store MigrateUp to create setting tables', got '%s'", migration.Description())
	}
}

func TestStoreSettingMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreSettingMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreSettingMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreSettingMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
