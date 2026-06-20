package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreSessionMigrate)(nil)

type StoreSessionMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreSessionMigrate) Signature() string {
	return "2026_03_21_0015_store_session_migrate"
}

func (m *StoreSessionMigrate) Description() string {
	return "Run session store MigrateUp to create session tables"
}

func (m *StoreSessionMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetSessionStore()
	if store == nil {
		return errors.New("session store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreSessionMigrate) Down() error {
	store := m.app.GetSessionStore()
	if store == nil {
		return errors.New("session store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

