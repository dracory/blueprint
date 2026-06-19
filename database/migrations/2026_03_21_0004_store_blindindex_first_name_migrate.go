package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreBlindIndexFirstNameMigrate)(nil)

type StoreBlindIndexFirstNameMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreBlindIndexFirstNameMigrate) Signature() string {
	return "2026_03_21_0004_store_blindindex_first_name_migrate"
}

func (m *StoreBlindIndexFirstNameMigrate) Description() string {
	return "Run blind index first name store MigrateUp to create blind index first name tables"
}

func (m *StoreBlindIndexFirstNameMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetBlindIndexStoreFirstName()
	if store == nil {
		return errors.New("blind index first name store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreBlindIndexFirstNameMigrate) Down() error {
	store := m.app.GetBlindIndexStoreFirstName()
	if store == nil {
		return errors.New("blind index first name store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

