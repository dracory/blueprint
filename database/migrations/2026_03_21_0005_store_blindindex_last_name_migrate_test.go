package migrations

import (
	"context"
	"testing"

	"github.com/dromara/carbon/v2"
)

func TestStoreBlindIndexLastNameMigrate(t *testing.T) {
	migration := &StoreBlindIndexLastNameMigrate{}

	t.Run("Interface Methods", func(t *testing.T) {
		if migration.ID() != "2026_03_21_0005_store_blindindex_last_name_migrate" {
			t.Errorf("Expected ID '2026_03_21_0005_store_blindindex_last_name_migrate', got '%s'", migration.ID())
		}

		if migration.Description() != "Run blind index last name store MigrateUp to create blind index last name tables" {
			t.Errorf("Expected description 'Run blind index last name store MigrateUp to create blind index last name tables', got '%s'", migration.Description())
		}

		createdAt := migration.CreatedAt()
		if createdAt.IsZero() {
			t.Error("Expected CreatedAt to return a non-zero time")
		}

		expectedTime := carbon.Parse("2026-03-21 00:05:00", "UTC").StdTime()
		if !createdAt.Equal(expectedTime) {
			t.Errorf("Expected CreatedAt to be %v, got %v", expectedTime, createdAt)
		}
	})

	t.Run("Up with nil registry", func(t *testing.T) {
		migration := &StoreBlindIndexLastNameMigrate{}
		err := migration.Up(context.Background(), nil)
		if err == nil {
			t.Error("Expected error when registry is nil")
		}
		if err.Error() != "registry is nil" {
			t.Errorf("Expected error 'registry is nil', got '%s'", err.Error())
		}
	})

	t.Run("Down with nil registry", func(t *testing.T) {
		migration := &StoreBlindIndexLastNameMigrate{}
		defer func() {
			if r := recover(); r != nil {
				// Expected panic due to nil registry
			} else {
				t.Error("Expected panic when registry is nil")
			}
		}()
		migration.Down(context.Background(), nil)
	})
}
