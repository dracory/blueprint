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

var _ migrate.MigrationInterface = (*StoreBlindIndexEmailMigrate)(nil)

type StoreBlindIndexEmailMigrate struct {
	registry registry.RegistryInterface
}

func (m *StoreBlindIndexEmailMigrate) ID() string {
	return "2026_03_21_0003_store_blindindex_email_migrate"
}

func (m *StoreBlindIndexEmailMigrate) Description() string {
	return "Run blind index email store MigrateUp to create blind index email tables"
}

func (m *StoreBlindIndexEmailMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetBlindIndexStoreEmail()
	if store == nil {
		return errors.New("blind index email store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreBlindIndexEmailMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.registry.GetBlindIndexStoreEmail()
	if store == nil {
		return errors.New("blind index email store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreBlindIndexEmailMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:03:00", "UTC").StdTime()
}
