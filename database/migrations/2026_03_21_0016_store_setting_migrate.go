package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreSettingMigrate)(nil)

type StoreSettingMigrate struct {
	registry RegistryInterface
}

func (m *StoreSettingMigrate) ID() string {
	return "2026_03_21_0016_store_setting_migrate"
}

func (m *StoreSettingMigrate) Description() string {
	return "Run setting store MigrateUp to create setting tables"
}

func (m *StoreSettingMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetSettingStore()
	if store == nil {
		return errors.New("setting store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreSettingMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.registry.GetSettingStore()
	if store == nil {
		return errors.New("setting store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreSettingMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:16:00", "UTC").StdTime()
}
