package app

import (
	"errors"

	"project/internal/config"

	neatdatabase "github.com/dracory/neat/database"
)

// databaseOpen opens the database connection using the provided config and returns the
// neat database instance. The underlying *sql.DB is derived from the neat instance so
// existing stores continue to receive a standard *sql.DB handle.
func databaseOpen(cfg config.ConfigInterface) (*neatdatabase.Database, error) {
	if cfg == nil {
		return nil, errors.New("databaseOpen: cfg is nil")
	}

	neatCfg := config.DatabaseNeatConfig(cfg)
	return neatdatabase.New(neatCfg)
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
