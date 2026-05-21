package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreSessionMigrate)(nil)

type StoreSessionMigrate struct {
	registry RegistryInterface
}

func (m *StoreSessionMigrate) ID() string {
	return "2026_03_21_0015_store_session_migrate"
}

func (m *StoreSessionMigrate) Description() string {
	return "Run session store MigrateUp to create session tables"
}

func (m *StoreSessionMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetSessionStore()
	if store == nil {
		return errors.New("session store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreSessionMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	return nil
}

func (m *StoreSessionMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:15:00", "UTC").StdTime()
}
