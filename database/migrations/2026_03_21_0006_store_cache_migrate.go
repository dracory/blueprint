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

var _ migrate.MigrationInterface = (*StoreCacheMigrate)(nil)

type StoreCacheMigrate struct {
	app app.AppInterface
}

func (m *StoreCacheMigrate) ID() string {
	return "2026_03_21_0006_store_cache_migrate"
}

func (m *StoreCacheMigrate) Description() string {
	return "Run cache store MigrateUp to create cache tables"
}

func (m *StoreCacheMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetCacheStore()
	if store == nil {
		return errors.New("cache store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreCacheMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.app.GetCacheStore()
	if store == nil {
		return errors.New("cache store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreCacheMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:06:00", "UTC").StdTime()
}
