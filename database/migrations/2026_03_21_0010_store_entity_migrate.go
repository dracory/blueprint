package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreEntityMigrate)(nil)

type StoreEntityMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreEntityMigrate) Signature() string {
	return "2026_03_21_0010_store_entity_migrate"
}

func (m *StoreEntityMigrate) Description() string {
	return "Run entity store MigrateUp to create entity tables"
}

func (m *StoreEntityMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetEntityStore()
	if store == nil {
		return errors.New("entity store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreEntityMigrate) Down() error {
	store := m.app.GetEntityStore()
	if store == nil {
		return errors.New("entity store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

