package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreBlindIndexLastNameMigrate)(nil)

type StoreBlindIndexLastNameMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreBlindIndexLastNameMigrate) Signature() string {
	return "2026_03_21_0005_store_blindindex_last_name_migrate"
}

func (m *StoreBlindIndexLastNameMigrate) Description() string {
	return "Run blind index last name store MigrateUp to create blind index last name tables"
}

func (m *StoreBlindIndexLastNameMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetBlindIndexStoreLastName()
	if store == nil {
		return errors.New("blind index last name store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreBlindIndexLastNameMigrate) Down() error {
	store := m.app.GetBlindIndexStoreLastName()
	if store == nil {
		return errors.New("blind index last name store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

