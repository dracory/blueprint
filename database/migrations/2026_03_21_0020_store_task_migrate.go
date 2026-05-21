package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreTaskMigrate)(nil)

type StoreTaskMigrate struct {
	registry RegistryInterface
}

func (m *StoreTaskMigrate) ID() string {
	return "2026_03_21_0020_store_task_migrate"
}

func (m *StoreTaskMigrate) Description() string {
	return "Run task store MigrateUp to create task tables"
}

func (m *StoreTaskMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetTaskStore()
	if store == nil {
		return errors.New("task store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreTaskMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.registry.GetTaskStore()
	if store == nil {
		return errors.New("task store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreTaskMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:20:00", "UTC").StdTime()
}
