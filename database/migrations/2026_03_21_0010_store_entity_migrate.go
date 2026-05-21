package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreEntityMigrate)(nil)

type StoreEntityMigrate struct {
	registry RegistryInterface
}

func (m *StoreEntityMigrate) ID() string {
	return "2026_03_21_0010_store_entity_migrate"
}

func (m *StoreEntityMigrate) Description() string {
	return "Run entity store MigrateUp to create entity tables"
}

func (m *StoreEntityMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetEntityStore()
	if store == nil {
		return errors.New("entity store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreEntityMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.registry.GetEntityStore()
	if store == nil {
		return errors.New("entity store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreEntityMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:10:00", "UTC").StdTime()
}
