package migrations

import (
	"context"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestStoreBlindIndexFirstNameMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreBlindIndexFirstNameMigrate{}

	if migration.ID() != "2026_03_21_0004_store_blindindex_first_name_migrate" {
		t.Errorf("Expected ID '2026_03_21_0004_store_blindindex_first_name_migrate', got '%s'", migration.ID())
	}

	if migration.Description() != "Run blind index first name store MigrateUp to create blind index first name tables" {
		t.Errorf("Expected description 'Run blind index first name store MigrateUp to create blind index first name tables', got '%s'", migration.Description())
	}

	createdAt := migration.CreatedAt()
	if createdAt.IsZero() {
		t.Error("Expected CreatedAt to return a non-zero time")
	}

	expectedTime := carbon.Parse("2026-03-21 00:04:00", "UTC").StdTime()
	if !createdAt.Equal(expectedTime) {
		t.Errorf("Expected CreatedAt to be %v, got %v", expectedTime, createdAt)
	}
}

func TestStoreBlindIndexFirstNameMigrate_UpWithNilRegistry(t *testing.T) {
	migration := &StoreBlindIndexFirstNameMigrate{}
	err := migration.Up(context.Background(), nil)
	if err == nil {
		t.Error("Expected error when registry is nil")
	}
	if err.Error() != "registry is nil" {
		t.Errorf("Expected error 'registry is nil', got '%s'", err.Error())
	}
}

func TestStoreBlindIndexFirstNameMigrate_DownWithNilRegistry(t *testing.T) {
	migration := &StoreBlindIndexFirstNameMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil registry
		} else {
			t.Error("Expected panic when registry is nil")
		}
	}()
	migration.Down(context.Background(), nil)
}
