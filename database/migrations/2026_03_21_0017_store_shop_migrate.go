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

var _ migrate.MigrationInterface = (*StoreShopMigrate)(nil)

type StoreShopMigrate struct {
	registry registry.RegistryInterface
}

func (m *StoreShopMigrate) ID() string {
	return "2026_03_21_0017_store_shop_migrate"
}

func (m *StoreShopMigrate) Description() string {
	return "Run shop store MigrateUp to create shop tables"
}

func (m *StoreShopMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetShopStore()
	if store == nil {
		return errors.New("shop store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreShopMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.registry.GetShopStore()
	if store == nil {
		return errors.New("shop store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreShopMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:17:00", "UTC").StdTime()
}
