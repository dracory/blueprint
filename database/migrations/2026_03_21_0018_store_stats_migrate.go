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

var _ migrate.MigrationInterface = (*StoreStatsMigrate)(nil)

type StoreStatsMigrate struct {
	app app.AppInterface
}

func (m *StoreStatsMigrate) ID() string {
	return "2026_03_21_0018_store_stats_migrate"
}

func (m *StoreStatsMigrate) Description() string {
	return "Run stats store MigrateUp to create stats tables"
}

func (m *StoreStatsMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetStatsStore()
	if store == nil {
		return errors.New("stats store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreStatsMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.app.GetStatsStore()
	if store == nil {
		return errors.New("stats store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreStatsMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:18:00", "UTC").StdTime()
}
