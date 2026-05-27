package migrations

import (
	"context"
	"errors"
	"fmt"
	"project/internal/registry"

	"github.com/dracory/migrate"
)

// MigrateAll runs all migrations in two phases:
// 1. Store migrations (MigrateUp/AutoMigrate) — run directly, not inside a transaction.
// 2. Custom SQL migrations — run via the migrate framework with transaction support.
func MigrateAll(registry registry.RegistryInterface) error {
	if registry == nil {
		return errors.New("registry is nil")
	}

	// Phase 1: Store-level migrations (run directly outside transactions)
	if err := migrateStores(registry); err != nil {
		return err
	}

	// Phase 2: Custom SQL migrations via the migrate framework
	if err := migrateSQL(registry); err != nil {
		return err
	}

	return nil
}

// migrateStores runs MigrateUp/AutoMigrate for each enabled store.
// These are not wrapped in transactions because the store packages
// manage their own database connections internally.
func migrateStores(registry registry.RegistryInterface) error {
	cfg := registry.GetConfig()

	storeMigrations := getStoreMigrations(cfg, registry)

	ctx := context.Background()
	for _, m := range storeMigrations {
		if err := m.Up(ctx, nil); err != nil {
			return fmt.Errorf("store migration %s failed: %w", m.ID(), err)
		}
	}

	return nil
}

// migrateSQL runs date-prefixed SQL migrations using the migrate framework.
func migrateSQL(registry registry.RegistryInterface) error {
	db := registry.GetDatabase()
	if db == nil {
		return errors.New("database is nil")
	}

	sqlMigrations, err := getSQLMigrations()
	if err != nil {
		return err
	}

	if len(sqlMigrations) == 0 {
		return nil
	}

	migrator, err := migrate.New(db, nil)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}

	for _, m := range sqlMigrations {
		migrator.AddMigration(m)
	}

	return migrator.Up(context.Background())
}
