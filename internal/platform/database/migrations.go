package database

import (
	"context"
	"log"
)

// MigrateDatabase creates the necessary database tables if they don't exist
func (d *Database) MigrateDatabase() error {
	// Create users table
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(255) PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		first_name VARCHAR(255) NOT NULL,
		last_name VARCHAR(255) NOT NULL,
		is_active BOOLEAN NOT NULL DEFAULT TRUE,
		created_at TIMESTAMP NOT NULL,
		updated_at TIMESTAMP NOT NULL,
		deleted_at TIMESTAMP NULL
	);
	`

	// Execute the migration queries
	ctx := context.Background()
	_, err := d.ExecuteContext(ctx, usersTable)
	if err != nil {
		log.Printf("Error creating users table: %v", err)
		return err
	}

	return nil
}
