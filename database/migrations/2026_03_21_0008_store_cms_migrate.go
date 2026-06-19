package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreCmsMigrate)(nil)

type StoreCmsMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreCmsMigrate) Signature() string {
	return "2026_03_21_0008_store_cms_migrate"
}

func (m *StoreCmsMigrate) Description() string {
	return "Run CMS store AutoMigrate to create CMS tables"
}

func (m *StoreCmsMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetCmsStore()
	if store == nil {
		return errors.New("cms store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreCmsMigrate) Down() error {
	store := m.app.GetCmsStore()
	if store == nil {
		return errors.New("cms store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

