package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreStatsMigrate)(nil)

type StoreStatsMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreStatsMigrate) Signature() string {
	return "2026_03_21_0018_store_stats_migrate"
}

func (m *StoreStatsMigrate) Description() string {
	return "Run stats store MigrateUp to create stats tables"
}

func (m *StoreStatsMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetStatsStore()
	if store == nil {
		return errors.New("stats store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreStatsMigrate) Down() error {
	store := m.app.GetStatsStore()
	if store == nil {
		return errors.New("stats store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

