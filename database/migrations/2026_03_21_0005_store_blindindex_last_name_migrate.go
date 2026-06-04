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

var _ migrate.MigrationInterface = (*StoreBlindIndexLastNameMigrate)(nil)

type StoreBlindIndexLastNameMigrate struct {
	app app.AppInterface
}

func (m *StoreBlindIndexLastNameMigrate) ID() string {
	return "2026_03_21_0005_store_blindindex_last_name_migrate"
}

func (m *StoreBlindIndexLastNameMigrate) Description() string {
	return "Run blind index last name store MigrateUp to create blind index last name tables"
}

func (m *StoreBlindIndexLastNameMigrate) Up(ctx context.Context, tx *sql.Tx) error {
	if m.app == nil {
		return errors.New("app is nil")
	}

	store := m.app.GetBlindIndexStoreLastName()
	if store == nil {
		return errors.New("blind index last name store is not initialized")
	}

	return store.MigrateUp(ctx)
}

func (m *StoreBlindIndexLastNameMigrate) Down(ctx context.Context, tx *sql.Tx) error {
	store := m.app.GetBlindIndexStoreLastName()
	if store == nil {
		return errors.New("blind index last name store is not initialized")
	}
	return store.MigrateDown(ctx, tx)
}

func (m *StoreBlindIndexLastNameMigrate) CreatedAt() time.Time {
	return carbon.Parse("2026-03-21 00:05:00", "UTC").StdTime()
}
