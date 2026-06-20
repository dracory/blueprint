package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreVaultMigrate)(nil)

type StoreVaultMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreVaultMigrate) Signature() string {
	return "2026_03_21_0022_store_vault_migrate"
}

func (m *StoreVaultMigrate) Description() string {
	return "Run vault store MigrateUp to create vault tables"
}

func (m *StoreVaultMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetVaultStore()
	if store == nil {
		return errors.New("vault store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreVaultMigrate) Down() error {
	store := m.app.GetVaultStore()
	if store == nil {
		return errors.New("vault store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

