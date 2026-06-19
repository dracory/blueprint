package migrations

import "testing"

func TestStoreTaskMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreTaskMigrate{}

	if migration.Signature() != "2026_03_21_0020_store_task_migrate" {
		t.Errorf("Expected signature '2026_03_21_0020_store_task_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run task store MigrateUp to create task tables" {
		t.Errorf("Expected description 'Run task store MigrateUp to create task tables', got '%s'", migration.Description())
	}
}

func TestStoreTaskMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreTaskMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreTaskMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreTaskMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
