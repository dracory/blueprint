package migrations

import (
	"context"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestStoreSettingMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreSettingMigrate{}

	if migration.ID() != "2026_03_21_0016_store_setting_migrate" {
		t.Errorf("Expected ID '2026_03_21_0016_store_setting_migrate', got '%s'", migration.ID())
	}

	if migration.Description() != "Run setting store MigrateUp to create setting tables" {
		t.Errorf("Expected description 'Run setting store MigrateUp to create setting tables', got '%s'", migration.Description())
	}

	createdAt := migration.CreatedAt()
	if createdAt.IsZero() {
		t.Error("Expected CreatedAt to return a non-zero time")
	}

	expectedTime := carbon.Parse("2026-03-21 00:16:00", "UTC").StdTime()
	if !createdAt.Equal(expectedTime) {
		t.Errorf("Expected CreatedAt to be %v, got %v", expectedTime, createdAt)
	}
}

func TestStoreSettingMigrate_UpWithNilRegistry(t *testing.T) {
	migration := &StoreSettingMigrate{}
	err := migration.Up(context.Background(), nil)
	if err == nil {
		t.Error("Expected error when registry is nil")
	}
	if err.Error() != "registry is nil" {
		t.Errorf("Expected error 'registry is nil', got '%s'", err.Error())
	}
}

func TestStoreSettingMigrate_DownWithNilRegistry(t *testing.T) {
	migration := &StoreSettingMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil registry
		} else {
			t.Error("Expected panic when registry is nil")
		}
	}()
	migration.Down(context.Background(), nil)
}
