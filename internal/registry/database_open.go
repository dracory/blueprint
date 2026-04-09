package registry

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"project/internal/config"

	"github.com/dracory/database"
	// "gorm.io/driver/postgres"
	// "gorm.io/driver/sqlite"
	// "gorm.io/gorm"
)

// databaseOpen opens the database connection using the provided config and returns it.
// The database package now includes SQLite optimizations and connection pool settings automatically.
func databaseOpen(cfg config.ConfigInterface) (*sql.DB, error) {
	if cfg == nil {
		return nil, errors.New("databaseOpen: cfg is nil")
	}

	options := database.Options().
		SetDatabaseType(cfg.GetDatabaseDriver()).
		SetDatabaseHost(cfg.GetDatabaseHost()).
		SetDatabasePort(cfg.GetDatabasePort()).
		SetDatabaseName(cfg.GetDatabaseName()).
		SetCharset(cfg.GetDatabaseCharset()).
		SetTimeZone(cfg.GetDatabaseTimezone()).
		SetUserName(cfg.GetDatabaseUsername()).
		SetPassword(cfg.GetDatabasePassword())

	if v := cfg.GetDatabaseMaxOpenConns(); v > 0 {
		options = options.SetMaxOpenConns(v)
	}
	if v := cfg.GetDatabaseMaxIdleConns(); v > 0 {
		options = options.SetMaxIdleConns(v)
	}
	if v := cfg.GetDatabaseConnMaxLifetimeSeconds(); v > 0 {
		options = options.SetConnMaxLifetime(time.Duration(v) * time.Second)
	}
	if v := cfg.GetDatabaseConnMaxIdleTimeSeconds(); v > 0 {
		options = options.SetConnMaxIdleTime(time.Duration(v) * time.Second)
	}

	// Set SSL mode for non-SQLite databases
	isSQLite := strings.Contains(strings.ToLower(cfg.GetDatabaseDriver()), "sqlite")
	if !isSQLite {
		sslMode := cfg.GetDatabaseSSLMode()
		if sslMode == "" {
			sslMode = "require"
		}
		options = options.SetSSLMode(sslMode)
	}

	return database.Open(options)
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
