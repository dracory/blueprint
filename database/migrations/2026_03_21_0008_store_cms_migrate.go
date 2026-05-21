package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreCmsMigrate)(nil)

type StoreCmsMigrate struct {
	registry RegistryInterface
}

func (m *StoreCmsMigrate) ID() string {
	return "2026_03_21_0008_store_cms_migrate"
}

func (m *StoreCmsMigrate) Description() string {
	return "Run CMS store AutoMigrate to create CMS tables"
}

func (m *StoreCmsMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetCmsStore()
	if store == nil {
		return errors.New("cms store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreCmsMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.registry.GetCmsStore()
	if store == nil {
		return errors.New("cms store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreCmsMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:08:00", "UTC").StdTime()
}
