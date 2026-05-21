package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreAuditMigrate)(nil)

type StoreAuditMigrate struct {
	registry RegistryInterface
}

func (m *StoreAuditMigrate) ID() string {
	return "2026_03_21_0001_store_audit_migrate"
}

func (m *StoreAuditMigrate) Description() string {
	return "Run audit store AutoMigrate to create audit tables"
}

func (m *StoreAuditMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetAuditStore()
	if store == nil {
		return errors.New("audit store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreAuditMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.registry.GetAuditStore()
	if store == nil {
		return errors.New("audit store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreAuditMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:01:00", "UTC").StdTime()
}
