package migrations

import (
	"context"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestStoreSessionMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreSessionMigrate{}

	if migration.ID() != "2026_03_21_0015_store_session_migrate" {
		t.Errorf("Expected ID '2026_03_21_0015_store_session_migrate', got '%s'", migration.ID())
	}

	if migration.Description() != "Run session store MigrateUp to create session tables" {
		t.Errorf("Expected description 'Run session store MigrateUp to create session tables', got '%s'", migration.Description())
	}

	createdAt := migration.CreatedAt()
	if createdAt.IsZero() {
		t.Error("Expected CreatedAt to return a non-zero time")
	}

	expectedTime := carbon.Parse("2026-03-21 00:15:00", "UTC").StdTime()
	if !createdAt.Equal(expectedTime) {
		t.Errorf("Expected CreatedAt to be %v, got %v", expectedTime, createdAt)
	}
}

func TestStoreSessionMigrate_UpWithNilRegistry(t *testing.T) {
	migration := &StoreSessionMigrate{}
	err := migration.Up(context.Background(), nil)
	if err == nil {
		t.Error("Expected error when app is nil")
	}
	if err.Error() != "app is nil" {
		t.Errorf("Expected error 'app is nil', got '%s'", err.Error())
	}
}

func TestStoreSessionMigrate_DownWithNilRegistry(t *testing.T) {
	migration := &StoreSessionMigrate{}
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil app
		} else {
			t.Error("Expected panic when app is nil")
		}
	}()
	migration.Down(context.Background(), nil)
}
