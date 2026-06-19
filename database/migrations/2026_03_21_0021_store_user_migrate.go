package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreUserMigrate)(nil)

type StoreUserMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreUserMigrate) Signature() string {
	return "2026_03_21_0021_store_user_migrate"
}

func (m *StoreUserMigrate) Description() string {
	return "Run user store MigrateUp to create user tables"
}

func (m *StoreUserMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetUserStore()
	if store == nil {
		return errors.New("user store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreUserMigrate) Down() error {
	store := m.app.GetUserStore()
	if store == nil {
		return errors.New("user store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

