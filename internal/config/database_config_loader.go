package config

import (
	"strings"

	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// databaseConfig captures database connection settings.
type databaseConfig struct {
	driver   string // Database driver type (sqlite, postgres, mysql, etc.)
	host     string // Database server hostname or IP address
	port     string // Database server port number
	name     string // Database name/schema
	username string // Database authentication username
	password string // Database authentication password
	sslMode  string // SSL connection mode (for PostgreSQL)
}

// loadDatabaseConfig loads database configuration from environment variables.
func loadDatabaseConfig(acc *baseCfg.LoadAccumulator) databaseConfig {
	driver := acc.MustString(KEY_DB_DRIVER, "select the database driver (e.g., sqlite, postgres)")
	host := strings.TrimSpace(env.GetString(KEY_DB_HOST))
	port := strings.TrimSpace(env.GetString(KEY_DB_PORT))
	name := acc.MustString(KEY_DB_DATABASE, "set the database name")
	user := strings.TrimSpace(env.GetString(KEY_DB_USERNAME))
	pass := strings.TrimSpace(env.GetString(KEY_DB_PASSWORD))

	if driver != driverSQLite {
		acc.MustWhen(true, KEY_DB_HOST, "required when `DB_DRIVER` is not sqlite", host)
		acc.MustWhen(true, KEY_DB_PORT, "required when `DB_DRIVER` is not sqlite", port)
		acc.MustWhen(true, KEY_DB_USERNAME, "required when `DB_DRIVER` is not sqlite", user)
		acc.MustWhen(true, KEY_DB_PASSWORD, "required when `DB_DRIVER` is not sqlite", pass)
	}

	return databaseConfig{
		driver:   driver,
		host:     host,
		port:     port,
		name:     name,
		username: user,
		password: pass,
		sslMode:  "require",
	}
}
