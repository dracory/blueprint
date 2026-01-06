package app

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"project/internal/types"

	"github.com/dracory/database"
	// "gorm.io/driver/postgres"
	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"
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
	} else {
		// Provide sensible defaults for production databases; adjust after load testing.
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(5 * time.Minute)
	}

	return db, nil
}

// Enable, if you want to use GORM
//
// func gormOpen(sqlDB *sql.DB, driverName string) (*gorm.DB, error) {
// 	// Open GORM on top of the given *sql.DB
// 	// We, however, need the driver name; default to sqlite if not provided.
// 	driver := strings.ToLower(strings.TrimSpace(driverName))
// 	if driver == "" {
// 		// Try to detect from DB handle
// 		driver = strings.ToLower(strings.TrimSpace(sb.DatabaseDriverName(sqlDB)))
// 		if driver == "" {
// 			// Best-effort default; callers should pass explicit driver when possible.
// 			driver = "sqlite"
// 		}
// 	}

// 	gcfg := &gorm.Config{}
// 	// gcfg.Logger = logger.Default.LogMode(logger.Info) // if we want to see the queries

// 	var (
// 		gdb *gorm.DB
// 		err error
// 	)
// 	switch driver {
// 	case "sqlite":
// 		gdb, err = gorm.Open(sqlite.Dialector{Conn: sqlDB}, gcfg)
// 	// case "mysql":
// 	// 	gdb, err = gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), gcfg)
// 	case "postgres", "postgresql":
// 		gdb, err = gorm.Open(postgres.New(postgres.Config{Conn: sqlDB, DriverName: "postgres"}), gcfg)
// 	default:
// 		return nil, errors.New("agent store: unsupported driver: " + driver)
// 	}
// 	if err != nil {
// 		return nil, err
// 	}

// 	return gdb, nil
// }
