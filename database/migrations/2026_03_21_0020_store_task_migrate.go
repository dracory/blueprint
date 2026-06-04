package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"project/internal/app"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreTaskMigrate)(nil)

type StoreTaskMigrate struct {
	app app.AppInterface
}

func (m *StoreTaskMigrate) ID() string {
	return "2026_03_21_0020_store_task_migrate"
}

func (m *StoreTaskMigrate) Description() string {
	return "Run task store MigrateUp to create task tables"
}

func (m *StoreTaskMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetTaskStore()
	if store == nil {
		return errors.New("task store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreTaskMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.app.GetTaskStore()
	if store == nil {
		return errors.New("task store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreTaskMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:20:00", "UTC").StdTime()
}
