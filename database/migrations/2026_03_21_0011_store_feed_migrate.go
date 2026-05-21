package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreFeedMigrate)(nil)

type StoreFeedMigrate struct {
	registry RegistryInterface
}

func (m *StoreFeedMigrate) ID() string {
	return "2026_03_21_0011_store_feed_migrate"
}

func (m *StoreFeedMigrate) Description() string {
	return "Run feed store MigrateUp to create feed tables"
}

func (m *StoreFeedMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetFeedStore()
	if store == nil {
		return errors.New("feed store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreFeedMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	return nil
}

func (m *StoreFeedMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:11:00", "UTC").StdTime()
}
