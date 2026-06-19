package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreCacheMigrate)(nil)

type StoreCacheMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreCacheMigrate) Signature() string {
	return "2026_03_21_0006_store_cache_migrate"
}

func (m *StoreCacheMigrate) Description() string {
	return "Run cache store MigrateUp to create cache tables"
}

func (m *StoreCacheMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetCacheStore()
	if store == nil {
		return errors.New("cache store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreCacheMigrate) Down() error {
	store := m.app.GetCacheStore()
	if store == nil {
		return errors.New("cache store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

