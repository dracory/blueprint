package migrations

import (
	"context"
	"errors"

	"project/internal/app"

	"github.com/dracory/neat/database/migrator"
)

var _ migrator.MigrationInterface = (*StoreSubscriptionMigrate)(nil)

type StoreSubscriptionMigrate struct {
	migrator.BaseMigration
	app app.AppInterface
}

func (m *StoreSubscriptionMigrate) Signature() string {
	return "2026_03_21_0019_store_subscription_migrate"
}

func (m *StoreSubscriptionMigrate) Description() string {
	return "Run subscription store AutoMigrate to create subscription tables"
}

func (m *StoreSubscriptionMigrate) Up() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetSubscriptionStore()
	if store == nil {
		return errors.New("subscription store is not initialized")
	}

	return store.MigrateUp(context.Background())
}

func (m *StoreSubscriptionMigrate) Down() error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetSubscriptionStore()
	if store == nil {
		return errors.New("subscription store is not initialized")
	}

	return store.MigrateDown(context.Background())
}
