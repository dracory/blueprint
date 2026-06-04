package migrations

import (
	"context"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestStoreAuditMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreAuditMigrate{}

	if migration.ID() != "2026_03_21_0001_store_audit_migrate" {
		t.Errorf("Expected ID '2026_03_21_0001_store_audit_migrate', got '%s'", migration.ID())
	}

	if migration.Description() != "Run audit store AutoMigrate to create audit tables" {
		t.Errorf("Expected description 'Run audit store AutoMigrate to create audit tables', got '%s'", migration.Description())
	}

	createdAt := migration.CreatedAt()
	if createdAt.IsZero() {
		t.Error("Expected CreatedAt to return a non-zero time")
	}

	expectedTime := carbon.Parse("2026-03-21 00:01:00", "UTC").StdTime()
	if !createdAt.Equal(expectedTime) {
		t.Errorf("Expected CreatedAt to be %v, got %v", expectedTime, createdAt)
	}
}

func TestStoreAuditMigrate_UpWithNilRegistry(t *testing.T) {
	migration := &StoreAuditMigrate{}
	err := migration.Up(context.Background(), nil)
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreAuditMigrate_DownWithNilRegistry(t *testing.T) {
	migration := &StoreAuditMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down(context.Background(), nil)
}
