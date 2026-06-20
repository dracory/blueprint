package migrations

import "testing"

func TestStoreSubscriptionMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreSubscriptionMigrate{}

	if migration.Signature() != "2026_03_21_0019_store_subscription_migrate" {
		t.Errorf("Expected signature '2026_03_21_0019_store_subscription_migrate', got '%s'", migration.Signature())
	}

	if migration.Description() != "Run subscription store AutoMigrate to create subscription tables" {
		t.Errorf("Expected description 'Run subscription store AutoMigrate to create subscription tables', got '%s'", migration.Description())
	}
}

func TestStoreSubscriptionMigrate_UpWithNilApp(t *testing.T) {
	migration := &StoreSubscriptionMigrate{}
	err := migration.Up()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreSubscriptionMigrate_DownWithNilApp(t *testing.T) {
	migration := &StoreSubscriptionMigrate{}
	err := migration.Down()
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}
