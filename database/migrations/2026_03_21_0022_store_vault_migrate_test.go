package migrations

import (
	"context"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestStoreVaultMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreVaultMigrate{}

	if migration.ID() != "2026_03_21_0022_store_vault_migrate" {
		t.Errorf("Expected ID '2026_03_21_0022_store_vault_migrate', got '%s'", migration.ID())
	}

	if migration.Description() != "Run vault store MigrateUp to create vault tables" {
		t.Errorf("Expected description 'Run vault store MigrateUp to create vault tables', got '%s'", migration.Description())
	}

	createdAt := migration.CreatedAt()
	if createdAt.IsZero() {
		t.Error("Expected CreatedAt to return a non-zero time")
	}

	expectedTime := carbon.Parse("2026-03-21 00:22:00", "UTC").StdTime()
	if !createdAt.Equal(expectedTime) {
		t.Errorf("Expected CreatedAt to be %v, got %v", expectedTime, createdAt)
	}
}

func TestStoreVaultMigrate_UpWithNilRegistry(t *testing.T) {
	migration := &StoreVaultMigrate{}
	err := migration.Up(context.Background(), nil)
	if err == nil {
		t.Error("Expected error when registry is nil")
	}
	if err.Error() != "registry is nil" {
		t.Errorf("Expected error 'registry is nil', got '%s'", err.Error())
	}
}

func TestStoreVaultMigrate_DownWithNilRegistry(t *testing.T) {
	migration := &StoreVaultMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil registry
		} else {
			t.Error("Expected panic when registry is nil")
		}
	}()
	migration.Down(context.Background(), nil)
}
