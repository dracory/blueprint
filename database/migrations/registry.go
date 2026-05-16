package migrations

import (
	"fmt"
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

// GetAllMigrations returns all available migrations with validation
func GetAllMigrations() ([]migrate.MigrationInterface, error) {
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
