package migrations

import (
	"context"
	"errors"
	"fmt"
	"project/internal/app"
	"project/internal/config"

	neatdb "github.com/dracory/neat/database"
	"github.com/dracory/neat/database/migrator"
)

// MigrateAll runs all migrations using the neat migrator package.
// Store migrations and custom SQL migrations are combined and tracked
// in a single migration_tracker table.
func MigrateAll(app app.AppInterface) error {
	if app == nil {
		return errors.New("app is nil")
	}

	db := app.GetDatabase()
	if db == nil {
		return errors.New("database is nil")
	}

	neatDB, err := neatdb.NewFromSQLDB(db)
	if err != nil {
		return fmt.Errorf("failed to create neat database: %w", err)
	}

	m := migrator.NewMigrator(neatDB)
	m.SetTransactionsEnabled(false)

	cfg := app.GetConfig()

	storeMigrations := getStoreMigrations(cfg, app)
	if err := m.AddMigrations(storeMigrations); err != nil {
		return fmt.Errorf("failed to add store migrations: %w", err)
	}

	sqlMigrations, err := getSQLMigrations()
	if err != nil {
		return err
	}
	if err := m.AddMigrations(sqlMigrations); err != nil {
		return fmt.Errorf("failed to add sql migrations: %w", err)
	}

	return m.Up(context.Background())
}

// getStoreMigrations returns store migrations conditionally based on config.
func getStoreMigrations(cfg config.ConfigInterface, reg app.AppInterface) []migrator.MigrationInterface {
	migrations := []migrator.MigrationInterface{}

	if cfg.GetAuditStoreUsed() {
		migrations = append(migrations, &StoreAuditMigrate{app: reg})
	}
	if cfg.GetBlogStoreUsed() {
		migrations = append(migrations, &StoreBlogMigrate{app: reg})
	}
	if cfg.GetUserStoreUsed() && cfg.GetVaultStoreUsed() {
		migrations = append(migrations, &StoreBlindIndexEmailMigrate{app: reg})
		migrations = append(migrations, &StoreBlindIndexFirstNameMigrate{app: reg})
		migrations = append(migrations, &StoreBlindIndexLastNameMigrate{app: reg})
	}
	if cfg.GetCacheStoreUsed() {
		migrations = append(migrations, &StoreCacheMigrate{app: reg})
	}
	if cfg.GetChatStoreUsed() {
		migrations = append(migrations, &StoreChatMigrate{app: reg})
	}
	if cfg.GetCmsStoreUsed() {
		migrations = append(migrations, &StoreCmsMigrate{app: reg})
	}
	if cfg.GetCustomStoreUsed() {
		migrations = append(migrations, &StoreCustomMigrate{app: reg})
	}
	if cfg.GetEntityStoreUsed() {
		migrations = append(migrations, &StoreEntityMigrate{app: reg})
	}
	if cfg.GetFeedStoreUsed() {
		migrations = append(migrations, &StoreFeedMigrate{app: reg})
	}
	if cfg.GetGeoStoreUsed() {
		migrations = append(migrations, &StoreGeoMigrate{app: reg})
	}
	if cfg.GetLogStoreUsed() {
		migrations = append(migrations, &StoreLogMigrate{app: reg})
	}
	if cfg.GetMetaStoreUsed() {
		migrations = append(migrations, &StoreMetaMigrate{app: reg})
	}
	if cfg.GetSessionStoreUsed() {
		migrations = append(migrations, &StoreSessionMigrate{app: reg})
	}
	if cfg.GetSettingStoreUsed() {
		migrations = append(migrations, &StoreSettingMigrate{app: reg})
	}
	if cfg.GetShopStoreUsed() {
		migrations = append(migrations, &StoreShopMigrate{app: reg})
	}
	if cfg.GetStatsStoreUsed() {
		migrations = append(migrations, &StoreStatsMigrate{app: reg})
	}
	if cfg.GetSubscriptionStoreUsed() {
		migrations = append(migrations, &StoreSubscriptionMigrate{app: reg})
	}
	if cfg.GetTaskStoreUsed() {
		migrations = append(migrations, &StoreTaskMigrate{app: reg})
	}
	if cfg.GetUserStoreUsed() {
		migrations = append(migrations, &StoreUserMigrate{app: reg})
	}
	if cfg.GetVaultStoreUsed() {
		migrations = append(migrations, &StoreVaultMigrate{app: reg})
	}

	return migrations
}

// getSQLMigrations returns custom SQL migrations.
func getSQLMigrations() ([]migrator.MigrationInterface, error) {
	migrations := []migrator.MigrationInterface{
		// &TableCustomCreate{},
		// &TableTapMessagesCreate{},
		// &TablePointersCreate{},
		// &AddProfileCompletedToUsers{},
		// &TableCacheCreate{},
		// &TableSessionsCreate{},
	}

	return migrations, nil
}
