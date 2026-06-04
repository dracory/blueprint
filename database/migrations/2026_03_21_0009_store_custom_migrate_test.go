package migrations

import (
	"context"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestStoreCustomMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreCustomMigrate{}

	if migration.ID() != "2026_03_21_0009_store_custom_migrate" {
		t.Errorf("Expected ID '2026_03_21_0009_store_custom_migrate', got '%s'", migration.ID())
	}

	if migration.Description() != "Run custom store MigrateUp to create custom tables" {
		t.Errorf("Expected description 'Run custom store MigrateUp to create custom tables', got '%s'", migration.Description())
	}

	createdAt := migration.CreatedAt()
	if createdAt.IsZero() {
		t.Error("Expected CreatedAt to return a non-zero time")
	}

	expectedTime := carbon.Parse("2026-03-21 00:09:00", "UTC").StdTime()
	if !createdAt.Equal(expectedTime) {
		t.Errorf("Expected CreatedAt to be %v, got %v", expectedTime, createdAt)
	}
}

func TestStoreCustomMigrate_UpWithNilRegistry(t *testing.T) {
	migration := &StoreCustomMigrate{}
	err := migration.Up(context.Background(), nil)
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreCustomMigrate_DownWithNilRegistry(t *testing.T) {
	migration := &StoreCustomMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down(context.Background(), nil)
}
