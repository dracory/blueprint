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

var _ migrate.MigrationInterface = (*StoreLogMigrate)(nil)

type StoreLogMigrate struct {
	app app.AppInterface
}

func (m *StoreLogMigrate) ID() string {
	return "2026_03_21_0013_store_log_migrate"
}

func (m *StoreLogMigrate) Description() string {
	return "Run log store MigrateUp to create log tables"
}

func (m *StoreLogMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetLogStore()
	if store == nil {
		return errors.New("log store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreLogMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.app.GetLogStore()
	if store == nil {
		return errors.New("log store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreLogMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:13:00", "UTC").StdTime()
}
