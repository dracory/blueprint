package migrations

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"project/internal/registry"

	"github.com/dracory/migrate"
	"github.com/dromara/carbon/v2"
)

var _ migrate.MigrationInterface = (*StoreSubscriptionMigrate)(nil)

type StoreSubscriptionMigrate struct {
	registry registry.RegistryInterface
}

func (m *StoreSubscriptionMigrate) ID() string {
	return "2026_03_21_0019_store_subscription_migrate"
}

func (m *StoreSubscriptionMigrate) Description() string {
	return "Run subscription store AutoMigrate to create subscription tables"
}

func (m *StoreSubscriptionMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.registry == nil {
		return errors.New("registry is nil")
	}

	store := m.registry.GetSubscriptionStore()
	if store == nil {
		return errors.New("subscription store is not initialized")
	}

	return store.AutoMigrate(ctx)
}

func (m *StoreSubscriptionMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	// Subscription store uses AutoMigrate which doesn't support rollback
	return nil
}

func (m *StoreSubscriptionMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:19:00", "UTC").StdTime()
}
