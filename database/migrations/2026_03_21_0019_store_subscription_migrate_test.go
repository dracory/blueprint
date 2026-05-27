package migrations

import (
	"context"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestStoreSubscriptionMigrate_InterfaceMethods(t *testing.T) {
	migration := &StoreSubscriptionMigrate{}

	if migration.ID() != "2026_03_21_0019_store_subscription_migrate" {
		t.Errorf("Expected ID '2026_03_21_0019_store_subscription_migrate', got '%s'", migration.ID())
	}

	if migration.Description() != "Run subscription store AutoMigrate to create subscription tables" {
		t.Errorf("Expected description 'Run subscription store AutoMigrate to create subscription tables', got '%s'", migration.Description())
	}

	createdAt := migration.CreatedAt()
	if createdAt.IsZero() {
		t.Error("Expected CreatedAt to return a non-zero time")
	}

	expectedTime := carbon.Parse("2026-03-21 00:19:00", "UTC").StdTime()
	if !createdAt.Equal(expectedTime) {
		t.Errorf("Expected CreatedAt to be %v, got %v", expectedTime, createdAt)
	}
}

func TestStoreSubscriptionMigrate_UpWithNilRegistry(t *testing.T) {
	migration := &StoreSubscriptionMigrate{}
	err := migration.Up(context.Background(), nil)
	if err == nil {
		t.Error("Expected error when registry is nil")
	}
	if err.Error() != "registry is nil" {
		t.Errorf("Expected error 'registry is nil', got '%s'", err.Error())
	}
}

func TestStoreSubscriptionMigrate_DownWithNilRegistry(t *testing.T) {
	migration := &StoreSubscriptionMigrate{}
	err := migration.Down(context.Background(), nil)
	if err != nil {
		t.Errorf("Expected nil error (no rollback), got '%s'", err.Error())
	}
}
