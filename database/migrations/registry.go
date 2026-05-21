package migrations

import (
	"fmt"
	"project/internal/config"
	"regexp"
	"strings"

	"github.com/dracory/migrate"
)

// validateMigrationID checks if the migration ID follows the expected format: 2026_03_21_{name}
func validateMigrationID(id string) error {
	// Pattern: YYYY_MM_DD_description
	// Example: 2026_03_21_table_users_create
	pattern := `^\d{4}_\d{2}_\d{2}_[a-zA-Z0-9_]+$`
	matched, err := regexp.MatchString(pattern, id)
	if err != nil {
		return fmt.Errorf("invalid regex pattern for migration ID validation: %w", err)
	}

	if !matched {
		return fmt.Errorf("migration ID '%s' must follow format YYYY_MM_DD_description (e.g., 2026_03_21_table_users_create)", id)
	}

	// Extract date part and validate it's a valid date
	parts := strings.Split(id, "_")
	if len(parts) < 4 {
		return fmt.Errorf("migration ID '%s' must have at least 4 parts separated by underscores", id)
	}

	// Basic date validation - you could use carbon for more sophisticated validation
	if parts[0] < "2020" || parts[0] > "2035" {
		return fmt.Errorf("migration ID '%s' year must be between 2020 and 2035, got %s", id, parts[0])
	}

	if parts[1] < "01" || parts[1] > "12" {
		return fmt.Errorf("migration ID '%s' month must be between 01 and 12, got %s", id, parts[1])
	}

	if parts[2] < "01" || parts[2] > "31" {
		return fmt.Errorf("migration ID '%s' day must be between 01 and 31, got %s", id, parts[2])
	}

	return nil
}

// getStoreMigrations returns store migrations conditionally based on config.
// These are run directly (not inside a transaction) because store packages
// manage their own database connections internally.
func getStoreMigrations(cfg config.ConfigInterface, registry RegistryInterface) []migrate.MigrationInterface {
	migrations := []migrate.MigrationInterface{}

	if cfg.GetAuditStoreUsed() {
		migrations = append(migrations, &StoreAuditMigrate{registry: registry})
	}
	if cfg.GetBlogStoreUsed() {
		migrations = append(migrations, &StoreBlogMigrate{registry: registry})
	}
	if cfg.GetUserStoreUsed() && cfg.GetVaultStoreUsed() {
		migrations = append(migrations, &StoreBlindIndexEmailMigrate{registry: registry})
		migrations = append(migrations, &StoreBlindIndexFirstNameMigrate{registry: registry})
		migrations = append(migrations, &StoreBlindIndexLastNameMigrate{registry: registry})
	}
	if cfg.GetCacheStoreUsed() {
		migrations = append(migrations, &StoreCacheMigrate{registry: registry})
	}
	if cfg.GetChatStoreUsed() {
		migrations = append(migrations, &StoreChatMigrate{registry: registry})
	}
	if cfg.GetCmsStoreUsed() {
		migrations = append(migrations, &StoreCmsMigrate{registry: registry})
	}
	if cfg.GetCustomStoreUsed() {
		migrations = append(migrations, &StoreCustomMigrate{registry: registry})
	}
	if cfg.GetEntityStoreUsed() {
		migrations = append(migrations, &StoreEntityMigrate{registry: registry})
	}
	if cfg.GetFeedStoreUsed() {
		migrations = append(migrations, &StoreFeedMigrate{registry: registry})
	}
	if cfg.GetGeoStoreUsed() {
		migrations = append(migrations, &StoreGeoMigrate{registry: registry})
	}
	if cfg.GetLogStoreUsed() {
		migrations = append(migrations, &StoreLogMigrate{registry: registry})
	}
	if cfg.GetMetaStoreUsed() {
		migrations = append(migrations, &StoreMetaMigrate{registry: registry})
	}
	if cfg.GetSessionStoreUsed() {
		migrations = append(migrations, &StoreSessionMigrate{registry: registry})
	}
	if cfg.GetSettingStoreUsed() {
		migrations = append(migrations, &StoreSettingMigrate{registry: registry})
	}
	if cfg.GetShopStoreUsed() {
		migrations = append(migrations, &StoreShopMigrate{registry: registry})
	}
	if cfg.GetStatsStoreUsed() {
		migrations = append(migrations, &StoreStatsMigrate{registry: registry})
	}
	if cfg.GetSubscriptionStoreUsed() {
		migrations = append(migrations, &StoreSubscriptionMigrate{registry: registry})
	}
	if cfg.GetTaskStoreUsed() {
		migrations = append(migrations, &StoreTaskMigrate{registry: registry})
	}
	if cfg.GetUserStoreUsed() {
		migrations = append(migrations, &StoreUserMigrate{registry: registry})
	}
	if cfg.GetVaultStoreUsed() {
		migrations = append(migrations, &StoreVaultMigrate{registry: registry})
	}

	return migrations
}

// getSQLMigrations returns custom SQL migrations with validation.
// These are run inside transactions via the migrate framework.
func getSQLMigrations() ([]migrate.MigrationInterface, error) {
	migrations := []migrate.MigrationInterface{
		&TableUsersCreate{},
		// &TableTapMessagesCreate{},
		// &TablePointersCreate{},
		// &AddProfileCompletedToUsers{},
		// &TableCacheCreate{},
		// &TableSessionsCreate{},
	}

	// Validate all migration IDs
	for _, migration := range migrations {
		if err := validateMigrationID(migration.ID()); err != nil {
			return nil, fmt.Errorf("invalid migration ID for %s: %w", migration.Description(), err)
		}
	}

	return migrations, nil
}
