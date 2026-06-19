package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreTaskMigrate)(nil)

type StoreTaskMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreTaskMigrate) Signature() string {
	return "2026_03_21_0020_store_task_migrate"
}

func (m *StoreTaskMigrate) Description() string {
	return "Run task store MigrateUp to create task tables"
}

func (m *StoreTaskMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetTaskStore()
	if store == nil {
		return errors.New("task store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreTaskMigrate) Down() error {
	store := m.app.GetTaskStore()
	if store == nil {
		return errors.New("task store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

