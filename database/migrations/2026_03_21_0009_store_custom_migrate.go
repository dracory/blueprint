package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreCustomMigrate)(nil)

type StoreCustomMigrate struct {
	registry RegistryInterface
}

func (m *StoreCustomMigrate) ID() string {
	return "2026_03_21_0009_store_custom_migrate"
}

func (m *StoreCustomMigrate) Description() string {
	return "Run custom store MigrateUp to create custom tables"
}

func (m *StoreCustomMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetCustomStore()
	if store == nil {
		return errors.New("custom store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreCustomMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	return nil
}

func (m *StoreCustomMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:09:00", "UTC").StdTime()
}
