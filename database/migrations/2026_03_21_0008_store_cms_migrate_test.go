package migrations

import "testing"

func TestStoreCmsMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreCmsMigrate{}

	if migration.Signature() != "2026_03_21_0008_store_cms_migrate" {
		t.Errorf("Expected signature '2026_03_21_0008_store_cms_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run CMS store AutoMigrate to create CMS tables" {
		t.Errorf("Expected description 'Run CMS store AutoMigrate to create CMS tables', got '%s'", migration.Description())
	}
}

func TestStoreCmsMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreCmsMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreCmsMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreCmsMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
