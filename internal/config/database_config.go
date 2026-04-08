package config

import (
	"strings"

	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// ============================================================================
// Interface
// ============================================================================

// DatabaseConfigInterface defines database configuration methods.
type DatabaseConfigInterface interface {
	SetDatabaseDriver(string)
	GetDatabaseDriver() string

	SetDatabaseHost(string)
	GetDatabaseHost() string

	SetDatabasePort(string)
	GetDatabasePort() string

	SetDatabaseName(string)
	GetDatabaseName() string

	SetDatabaseUsername(string)
	GetDatabaseUsername() string

	SetDatabasePassword(string)
	GetDatabasePassword() string

	SetDatabaseSSLMode(string)
	GetDatabaseSSLMode() string
}

// ============================================================================
// Types
// ============================================================================

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

// ============================================================================
// Loader
// ============================================================================

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

// ============================================================================
// Implementation (Getters/Setters)
// ============================================================================

func (c *configImplementation) SetDatabaseDriver(v string) {
	c.databaseDriver = v
}

func (c *configImplementation) GetDatabaseDriver() string {
	return c.databaseDriver
}

func (c *configImplementation) SetDatabaseHost(v string) {
	c.databaseHost = v
}

func (c *configImplementation) GetDatabaseHost() string {
	return c.databaseHost
}

func (c *configImplementation) SetDatabasePort(v string) {
	c.databasePort = v
}

func (c *configImplementation) GetDatabasePort() string {
	return c.databasePort
}

func (c *configImplementation) SetDatabaseName(v string) {
	c.databaseName = v
}

func (c *configImplementation) GetDatabaseName() string {
	return c.databaseName
}

func (c *configImplementation) SetDatabaseUsername(v string) {
	c.databaseUsername = v
}

func (c *configImplementation) GetDatabaseUsername() string {
	return c.databaseUsername
}

func (c *configImplementation) SetDatabasePassword(v string) {
	c.databasePassword = v
}

func (c *configImplementation) GetDatabasePassword() string {
	return c.databasePassword
}

func (c *configImplementation) SetDatabaseSSLMode(v string) {
	c.databaseSSLMode = v
}

func (c *configImplementation) GetDatabaseSSLMode() string {
	return c.databaseSSLMode
}
