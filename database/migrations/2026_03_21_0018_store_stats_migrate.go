package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"project/internal/registry"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreStatsMigrate)(nil)

type StoreStatsMigrate struct {
	registry registry.RegistryInterface
}

func (m *StoreStatsMigrate) ID() string {
	return "2026_03_21_0018_store_stats_migrate"
}

func (m *StoreStatsMigrate) Description() string {
	return "Run stats store MigrateUp to create stats tables"
}

func (m *StoreStatsMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetStatsStore()
	if store == nil {
		return errors.New("stats store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreStatsMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.registry.GetStatsStore()
	if store == nil {
		return errors.New("stats store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreStatsMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:18:00", "UTC").StdTime()
}
