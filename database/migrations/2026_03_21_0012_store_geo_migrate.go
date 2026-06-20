package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreGeoMigrate)(nil)

type StoreGeoMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreGeoMigrate) Signature() string {
	return "2026_03_21_0012_store_geo_migrate"
}

func (m *StoreGeoMigrate) Description() string {
	return "Run geo store MigrateUp to create geo tables"
}

func (m *StoreGeoMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetGeoStore()
	if store == nil {
		return errors.New("geo store is not initialized")
	}

	if err := store.MigrateUp(context.Background()); err != nil {
		return err
	}

	// Seed geolocation data (countries, states, timezones)
	return store.Seed(context.Background())
}

func (m *StoreGeoMigrate) Down() error {
	store := m.app.GetGeoStore()
	if store == nil {
		return errors.New("geo store is not initialized")
	}
	return store.MigrateDown(context.Background())
}

