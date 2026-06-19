package migrations

import "testing"

func TestStoreSessionMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreSessionMigrate{}

	if migration.Signature() != "2026_03_21_0015_store_session_migrate" {
		t.Errorf("Expected signature '2026_03_21_0015_store_session_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run session store MigrateUp to create session tables" {
		t.Errorf("Expected description 'Run session store MigrateUp to create session tables', got '%s'", migration.Description())
	}
}

func TestStoreSessionMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreSessionMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreSessionMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreSessionMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
