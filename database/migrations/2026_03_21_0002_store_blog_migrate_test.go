package migrations

import "testing"

func TestStoreBlogMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreBlogMigrate{}

	if migration.Signature() != "2026_03_21_0002_store_blog_migrate" {
		t.Errorf("Expected signature '2026_03_21_0002_store_blog_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run blog store MigrateUp to create blog tables" {
		t.Errorf("Expected description 'Run blog store MigrateUp to create blog tables', got '%s'", migration.Description())
	}
}

func TestStoreBlogMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreBlogMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreBlogMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreBlogMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down()
}
