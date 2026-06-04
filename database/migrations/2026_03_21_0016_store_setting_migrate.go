package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"project/internal/app"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreSettingMigrate)(nil)

type StoreSettingMigrate struct {
	app app.AppInterface
}

func (m *StoreSettingMigrate) ID() string {
	return "2026_03_21_0016_store_setting_migrate"
}

func (m *StoreSettingMigrate) Description() string {
	return "Run setting store MigrateUp to create setting tables"
}

func (m *StoreSettingMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetSettingStore()
	if store == nil {
		return errors.New("setting store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreSettingMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.app.GetSettingStore()
	if store == nil {
		return errors.New("setting store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreSettingMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:16:00", "UTC").StdTime()
}
