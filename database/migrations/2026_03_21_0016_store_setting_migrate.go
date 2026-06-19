package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreSettingMigrate)(nil)

type StoreSettingMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreSettingMigrate) Signature() string {
	return "2026_03_21_0016_store_setting_migrate"
}

func (m *StoreSettingMigrate) Description() string {
	return "Run setting store MigrateUp to create setting tables"
}

func (m *StoreSettingMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetSettingStore()
	if store == nil {
		return errors.New("setting store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreSettingMigrate) Down() error {
	store := m.app.GetSettingStore()
	if store == nil {
		return errors.New("setting store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

