package migrations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"project/internal/config"

	"github.com/dracory/auditstore"
	"github.com/dracory/blindindexstore"
	"github.com/dracory/blogstore"
	"github.com/dracory/cachestore"
	"github.com/dracory/chatstore"
	"github.com/dracory/cmsstore"
	"github.com/dracory/customstore"
	"github.com/dracory/entitystore"
	"github.com/dracory/feedstore"
	"github.com/dracory/filesystem"
	"github.com/dracory/geostore"
	"github.com/dracory/logstore"
	"github.com/dracory/metastore"
	"github.com/dracory/migrate"
	"github.com/dracory/sessionstore"
	"github.com/dracory/settingstore"
	"github.com/dracory/shopstore"
	"github.com/dracory/statsstore"
	"github.com/dracory/subscriptionstore"
	"github.com/dracory/taskstore"
	"github.com/dracory/userstore"
	"github.com/dracory/vaultstore"
)

// RegistryInterface is the narrow subset of the application registry needed
// to run migrations. The application's full registry satisfies this interface.
type RegistryInterface interface {
	GetConfig() config.ConfigInterface
	GetDatabase() *sql.DB

	// Stores
	GetAuditStore() auditstore.StoreInterface
	GetBlogStore() blogstore.StoreInterface
	GetBlindIndexStoreEmail() blindindexstore.StoreInterface
	GetBlindIndexStoreFirstName() blindindexstore.StoreInterface
	GetBlindIndexStoreLastName() blindindexstore.StoreInterface
	GetCacheStore() cachestore.StoreInterface
	GetChatStore() chatstore.StoreInterface
	GetCmsStore() cmsstore.StoreInterface
	GetCustomStore() customstore.StoreInterface
	GetEntityStore() entitystore.StoreInterface
	GetFeedStore() feedstore.StoreInterface
	GetGeoStore() geostore.StoreInterface
	GetLogStore() logstore.StoreInterface
	GetMetaStore() metastore.StoreInterface
	GetSessionStore() sessionstore.StoreInterface
	GetSettingStore() settingstore.StoreInterface
	GetShopStore() shopstore.StoreInterface
	GetSqlFileStorage() filesystem.StorageInterface
	GetStatsStore() statsstore.StoreInterface
	GetSubscriptionStore() subscriptionstore.StoreInterface
	GetTaskStore() taskstore.StoreInterface
	GetUserStore() userstore.StoreInterface
	GetVaultStore() vaultstore.StoreInterface
}

// MigrateAll runs all migrations in two phases:
// 1. Store migrations (MigrateUp/AutoMigrate) — run directly, not inside a transaction.
// 2. Custom SQL migrations — run via the migrate framework with transaction support.
func MigrateAll(registry RegistryInterface) error {
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
func migrateStores(registry RegistryInterface) error {
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
func migrateSQL(registry RegistryInterface) error {
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
