package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreBlindIndexEmailMigrate)(nil)

type StoreBlindIndexEmailMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreBlindIndexEmailMigrate) Signature() string {
	return "2026_03_21_0003_store_blindindex_email_migrate"
}

func (m *StoreBlindIndexEmailMigrate) Description() string {
	return "Run blind index email store MigrateUp to create blind index email tables"
}

func (m *StoreBlindIndexEmailMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetBlindIndexStoreEmail()
	if store == nil {
		return errors.New("blind index email store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreBlindIndexEmailMigrate) Down() error {
	store := m.app.GetBlindIndexStoreEmail()
	if store == nil {
		return errors.New("blind index email store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

