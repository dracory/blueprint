package migrations

import "testing"

func TestStoreChatMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreChatMigrate{}

	if migration.Signature() != "2026_03_21_0007_store_chat_migrate" {
		t.Errorf("Expected signature '2026_03_21_0007_store_chat_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run chat store MigrateUp to create chat tables" {
		t.Errorf("Expected description 'Run chat store MigrateUp to create chat tables', got '%s'", migration.Description())
	}
}

func TestStoreChatMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreChatMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreChatMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreChatMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
