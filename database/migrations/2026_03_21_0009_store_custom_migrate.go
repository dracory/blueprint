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

var _ migrate.MigrationInterface = (*StoreCustomMigrate)(nil)

type StoreCustomMigrate struct {
	app app.AppInterface
}

func (m *StoreCustomMigrate) ID() string {
	return "2026_03_21_0009_store_custom_migrate"
}

func (m *StoreCustomMigrate) Description() string {
	return "Run custom store MigrateUp to create custom tables"
}

func (m *StoreCustomMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetCustomStore()
	if store == nil {
		return errors.New("custom store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreCustomMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.app.GetCustomStore()
	if store == nil {
		return errors.New("custom store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreCustomMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:09:00", "UTC").StdTime()
}
