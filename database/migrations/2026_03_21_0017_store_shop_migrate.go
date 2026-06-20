package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreShopMigrate)(nil)

type StoreShopMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreShopMigrate) Signature() string {
	return "2026_03_21_0017_store_shop_migrate"
}

func (m *StoreShopMigrate) Description() string {
	return "Run shop store MigrateUp to create shop tables"
}

func (m *StoreShopMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetShopStore()
	if store == nil {
		return errors.New("shop store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreShopMigrate) Down() error {
	store := m.app.GetShopStore()
	if store == nil {
		return errors.New("shop store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

