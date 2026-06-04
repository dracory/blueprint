package migrations

import (
	"context"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestStoreEntityMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreEntityMigrate{}

	if migration.ID() != "2026_03_21_0010_store_entity_migrate" {
		t.Errorf("Expected ID '2026_03_21_0010_store_entity_migrate', got '%s'", migration.ID())
	}

	if migration.Description() != "Run entity store MigrateUp to create entity tables" {
		t.Errorf("Expected description 'Run entity store MigrateUp to create entity tables', got '%s'", migration.Description())
	}

	createdAt := migration.CreatedAt()
	if createdAt.IsZero() {
		t.Error("Expected CreatedAt to return a non-zero time")
	}

	expectedTime := carbon.Parse("2026-03-21 00:10:00", "UTC").StdTime()
	if !createdAt.Equal(expectedTime) {
		t.Errorf("Expected CreatedAt to be %v, got %v", expectedTime, createdAt)
	}
}

func TestStoreEntityMigrate_UpWithNilRegistry(t *testing.T) {
	migration := &StoreEntityMigrate{}
	err := migration.Up(context.Background(), nil)
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreEntityMigrate_DownWithNilRegistry(t *testing.T) {
	migration := &StoreEntityMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down(context.Background(), nil)
}
