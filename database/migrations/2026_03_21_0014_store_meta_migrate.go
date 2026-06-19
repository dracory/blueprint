package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreMetaMigrate)(nil)

type StoreMetaMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreMetaMigrate) Signature() string {
	return "2026_03_21_0014_store_meta_migrate"
}

func (m *StoreMetaMigrate) Description() string {
	return "Run meta store MigrateUp to create meta tables"
}

func (m *StoreMetaMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetMetaStore()
	if store == nil {
		return errors.New("meta store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreMetaMigrate) Down() error {
	store := m.app.GetMetaStore()
	if store == nil {
		return errors.New("meta store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

