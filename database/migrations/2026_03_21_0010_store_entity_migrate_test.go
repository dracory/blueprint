package migrations

import "testing"

func TestStoreEntityMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreEntityMigrate{}

	if migration.Signature() != "2026_03_21_0010_store_entity_migrate" {
		t.Errorf("Expected signature '2026_03_21_0010_store_entity_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run entity store MigrateUp to create entity tables" {
		t.Errorf("Expected description 'Run entity store MigrateUp to create entity tables', got '%s'", migration.Description())
	}
}

func TestStoreEntityMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreEntityMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreEntityMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreEntityMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
