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

var _ migrate.MigrationInterface = (*StoreFeedMigrate)(nil)

type StoreFeedMigrate struct {
	app app.AppInterface
}

func (m *StoreFeedMigrate) ID() string {
	return "2026_03_21_0011_store_feed_migrate"
}

func (m *StoreFeedMigrate) Description() string {
	return "Run feed store MigrateUp to create feed tables"
}

func (m *StoreFeedMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetFeedStore()
	if store == nil {
		return errors.New("feed store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreFeedMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.app.GetFeedStore()
	if store == nil {
		return errors.New("feed store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreFeedMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:11:00", "UTC").StdTime()
}
