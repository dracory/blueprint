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

var _ migrate.MigrationInterface = (*StoreUserMigrate)(nil)

type StoreUserMigrate struct {
	app app.AppInterface
}

func (m *StoreUserMigrate) ID() string {
	return "2026_03_21_0021_store_user_migrate"
}

func (m *StoreUserMigrate) Description() string {
	return "Run user store MigrateUp to create user tables"
}

func (m *StoreUserMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetUserStore()
	if store == nil {
		return errors.New("user store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreUserMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.app.GetUserStore()
	if store == nil {
		return errors.New("user store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreUserMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:21:00", "UTC").StdTime()
}
