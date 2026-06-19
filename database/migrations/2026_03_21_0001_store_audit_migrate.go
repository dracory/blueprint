package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreAuditMigrate)(nil)

type StoreAuditMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreAuditMigrate) Signature() string {
	return "2026_03_21_0001_store_audit_migrate"
}

func (m *StoreAuditMigrate) Description() string {
	return "Run audit store AutoMigrate to create audit tables"
}

func (m *StoreAuditMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetAuditStore()
	if store == nil {
		return errors.New("audit store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreAuditMigrate) Down() error {
	store := m.app.GetAuditStore()
	if store == nil {
		return errors.New("audit store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

