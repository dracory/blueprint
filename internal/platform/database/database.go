package database

import (
	"context"
	"database/sql"
	"fmt"

	"project/app/config"

	basedb "github.com/dracory/base/database"
)

// Database represents a database connection wrapper
type Database struct {
	DB     *sql.DB
	Config *config.Config
}

// New creates a new database connection using the base package
func New(cfg *config.Config) (*Database, error) {
	options := basedb.Options()

	// Set database options based on configuration
	options.SetDatabaseType(cfg.DatabaseDriver)
	options.SetDatabaseName(cfg.DatabaseName)

	// Only set these for non-SQLite databases
	if cfg.DatabaseDriver != basedb.DATABASE_TYPE_SQLITE {
		options.SetDatabaseHost(cfg.DatabaseHost)
		options.SetDatabasePort(cfg.DatabasePort)
		options.SetUserName(cfg.DatabaseUser)
		options.SetPassword(cfg.DatabasePassword)
	}

	// Open database connection using the base package
	db, err := basedb.Open(options)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{
		DB:     db,
		Config: cfg,
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.DB.Close()
}

// GetDB returns the underlying sql.DB instance
func (d *Database) GetDB() *sql.DB {
	return d.DB
}

// ExecuteContext executes a query with context without returning any rows
func (d *Database) ExecuteContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.DB.ExecContext(ctx, query, args...)
}
