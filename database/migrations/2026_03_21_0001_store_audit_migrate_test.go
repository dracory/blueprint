package migrations

import "testing"

func TestStoreAuditMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreAuditMigrate{}

	if migration.Signature() != "2026_03_21_0001_store_audit_migrate" {
		t.Errorf("Expected signature '2026_03_21_0001_store_audit_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run audit store AutoMigrate to create audit tables" {
		t.Errorf("Expected description 'Run audit store AutoMigrate to create audit tables', got '%s'", migration.Description())
	}
}

func TestStoreAuditMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreAuditMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreAuditMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreAuditMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
