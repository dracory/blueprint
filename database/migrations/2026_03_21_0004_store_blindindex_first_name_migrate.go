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

var _ migrate.MigrationInterface = (*StoreBlindIndexFirstNameMigrate)(nil)

type StoreBlindIndexFirstNameMigrate struct {
	registry registry.RegistryInterface
}

func (m *StoreBlindIndexFirstNameMigrate) ID() string {
	return "2026_03_21_0004_store_blindindex_first_name_migrate"
}

func (m *StoreBlindIndexFirstNameMigrate) Description() string {
	return "Run blind index first name store MigrateUp to create blind index first name tables"
}

func (m *StoreBlindIndexFirstNameMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetBlindIndexStoreFirstName()
	if store == nil {
		return errors.New("blind index first name store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreBlindIndexFirstNameMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.registry.GetBlindIndexStoreFirstName()
	if store == nil {
		return errors.New("blind index first name store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreBlindIndexFirstNameMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:04:00", "UTC").StdTime()
}
