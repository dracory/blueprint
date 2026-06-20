package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreFeedMigrate)(nil)

type StoreFeedMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreFeedMigrate) Signature() string {
	return "2026_03_21_0011_store_feed_migrate"
}

func (m *StoreFeedMigrate) Description() string {
	return "Run feed store MigrateUp to create feed tables"
}

func (m *StoreFeedMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetFeedStore()
	if store == nil {
		return errors.New("feed store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreFeedMigrate) Down() error {
	store := m.app.GetFeedStore()
	if store == nil {
		return errors.New("feed store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

