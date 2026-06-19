package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreLogMigrate)(nil)

type StoreLogMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreLogMigrate) Signature() string {
	return "2026_03_21_0013_store_log_migrate"
}

func (m *StoreLogMigrate) Description() string {
	return "Run log store MigrateUp to create log tables"
}

func (m *StoreLogMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetLogStore()
	if store == nil {
		return errors.New("log store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreLogMigrate) Down() error {
	store := m.app.GetLogStore()
	if store == nil {
		return errors.New("log store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

