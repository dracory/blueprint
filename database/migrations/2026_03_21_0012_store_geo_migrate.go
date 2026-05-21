package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreGeoMigrate)(nil)

type StoreGeoMigrate struct {
	registry RegistryInterface
}

func (m *StoreGeoMigrate) ID() string {
	return "2026_03_21_0012_store_geo_migrate"
}

func (m *StoreGeoMigrate) Description() string {
	return "Run geo store MigrateUp to create geo tables"
}

func (m *StoreGeoMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetGeoStore()
	if store == nil {
		return errors.New("geo store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreGeoMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	return nil
}

func (m *StoreGeoMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:12:00", "UTC").StdTime()
}
