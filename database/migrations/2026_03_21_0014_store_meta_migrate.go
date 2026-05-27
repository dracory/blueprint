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

var _ migrate.MigrationInterface = (*StoreMetaMigrate)(nil)

type StoreMetaMigrate struct {
	registry registry.RegistryInterface
}

func (m *StoreMetaMigrate) ID() string {
	return "2026_03_21_0014_store_meta_migrate"
}

func (m *StoreMetaMigrate) Description() string {
	return "Run meta store MigrateUp to create meta tables"
}

func (m *StoreMetaMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetMetaStore()
	if store == nil {
		return errors.New("meta store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreMetaMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.registry.GetMetaStore()
	if store == nil {
		return errors.New("meta store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreMetaMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:14:00", "UTC").StdTime()
}
