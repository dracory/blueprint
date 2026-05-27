package migrations

import (
	"context"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestStoreLogMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreLogMigrate{}

	if migration.ID() != "2026_03_21_0013_store_log_migrate" {
		t.Errorf("Expected ID '2026_03_21_0013_store_log_migrate', got '%s'", migration.ID())
	}

	if migration.Description() != "Run log store MigrateUp to create log tables" {
		t.Errorf("Expected description 'Run log store MigrateUp to create log tables', got '%s'", migration.Description())
	}

	createdAt := migration.CreatedAt()
	if createdAt.IsZero() {
		t.Error("Expected CreatedAt to return a non-zero time")
	}

	expectedTime := carbon.Parse("2026-03-21 00:13:00", "UTC").StdTime()
	if !createdAt.Equal(expectedTime) {
		t.Errorf("Expected CreatedAt to be %v, got %v", expectedTime, createdAt)
	}
}

func TestStoreLogMigrate_UpWithNilRegistry(t *testing.T) {
	migration := &StoreLogMigrate{}
	err := migration.Up(context.Background(), nil)
	if err == nil {
		t.Error("Expected error when registry is nil")
	}
	if err.Error() != "registry is nil" {
		t.Errorf("Expected error 'registry is nil', got '%s'", err.Error())
	}
}

func TestStoreLogMigrate_DownWithNilRegistry(t *testing.T) {
	migration := &StoreLogMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil registry
		} else {
			t.Error("Expected panic when registry is nil")
		}
	}()
	migration.Down(context.Background(), nil)
}
