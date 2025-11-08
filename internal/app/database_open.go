package app

import (
	"database/sql"
	"errors"
	"strings"

	"project/internal/types"

	"github.com/dracory/database"
)

// databaseOpen opens the database connection using the provided config and returns it.
func databaseOpen(cfg types.ConfigInterface) (*sql.DB, error) {
	if cfg == nil {
		return nil, errors.New("databaseOpen: cfg is nil")
	}

	isSQLite := strings.Contains(strings.ToLower(cfg.GetDatabaseDriver()), "sqlite")

	options := database.Options().
		SetDatabaseType(cfg.GetDatabaseDriver()).
		SetDatabaseHost(cfg.GetDatabaseHost()).
		SetDatabasePort(cfg.GetDatabasePort()).
		SetDatabaseName(cfg.GetDatabaseName()).
		SetCharset(`utf8mb4`).
		SetTimeZone("UTC").
		SetUserName(cfg.GetDatabaseUsername()).
		SetPassword(cfg.GetDatabasePassword())

	if !isSQLite {
		sslMode := cfg.GetDatabaseSSLMode()
		if sslMode == "" {
			sslMode = "require"
		}
		options = options.SetSSLMode(sslMode)
	}

	db, err := database.Open(options)

	if err != nil {
		return nil, err
	}

	// Add connection pool and driver-specific settings
	// For SQLite, reduce lock contention by enabling WAL and busy timeout,
	// and by constraining pool concurrency.
	if isSQLite {
		// Enable WAL mode for better concurrency; ignore errors if already set.
		_, _ = db.Exec("PRAGMA journal_mode=WAL;")
		// Use NORMAL synchronous for WAL (durable enough, faster writes).
		_, _ = db.Exec("PRAGMA synchronous=NORMAL;")
		// Ensure foreign keys are enforced.
		_, _ = db.Exec("PRAGMA foreign_keys=ON;")
		// Back off up to 5s when the database is busy instead of returning SQLITE_BUSY immediately.
		_, _ = db.Exec("PRAGMA busy_timeout=5000;")

		// Constrain the pool to avoid multiple concurrent writers on SQLite.
		// Increase carefully if needed; 1 is the safest to avoid SQLITE_BUSY.
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
	}

	// Add connection pool settings
	// db.SetMaxOpenConns(25)                 // Maximum number of open connections
	// db.SetMaxIdleConns(5)                  // Maximum number of idle connections
	// db.SetConnMaxLifetime(5 * time.Minute) // Maximum connection lifetime

	return db, nil
}
