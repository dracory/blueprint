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

var _ migrate.MigrationInterface = (*StoreVaultMigrate)(nil)

type StoreVaultMigrate struct {
	app app.AppInterface
}

func (m *StoreVaultMigrate) ID() string {
	return "2026_03_21_0022_store_vault_migrate"
}

func (m *StoreVaultMigrate) Description() string {
	return "Run vault store MigrateUp to create vault tables"
}

func (m *StoreVaultMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetVaultStore()
	if store == nil {
		return errors.New("vault store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreVaultMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.app.GetVaultStore()
	if store == nil {
		return errors.New("vault store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreVaultMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:22:00", "UTC").StdTime()
}
