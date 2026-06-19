package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreCustomMigrate)(nil)

type StoreCustomMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreCustomMigrate) Signature() string {
	return "2026_03_21_0009_store_custom_migrate"
}

func (m *StoreCustomMigrate) Description() string {
	return "Run custom store MigrateUp to create custom tables"
}

func (m *StoreCustomMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetCustomStore()
	if store == nil {
		return errors.New("custom store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreCustomMigrate) Down() error {
	store := m.app.GetCustomStore()
	if store == nil {
		return errors.New("custom store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

