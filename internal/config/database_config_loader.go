package config

import "strings"

// loadDatabaseConfig loads database configuration directly into the config.
func loadDatabaseConfig(cfg ConfigInterface, v *envValidator) {
	// Database Driver
	//
	// The database driver to use for the application.
	// Supported values: sqlite, postgres, mysql
	driver := v.GetStringOrError(KEY_DB_DRIVER, "select the database driver (e.g., sqlite, postgres)")

	// Database Host
	//
	// The hostname or IP address of the database server.
	// Not required when using sqlite.
	host := strings.TrimSpace(v.GetString(KEY_DB_HOST))

	// Database Port
	//
	// The port the database server is listening on.
	// Common defaults: postgres=5432, mysql=3306
	// Not required when using sqlite.
	port := strings.TrimSpace(v.GetString(KEY_DB_PORT))

	// Database Name
	//
	// The name of the database to connect to.
	// For sqlite, this is the file path (e.g., ./database.db or :memory:)
	name := v.GetStringOrError(KEY_DB_DATABASE, "set the database name")

	// Database Username
	//
	// The username for authenticating with the database server.
	// Not required when using sqlite.
	user := strings.TrimSpace(v.GetString(KEY_DB_USERNAME))

	// Database Password
	//
	// The password for authenticating with the database server.
	// Not required when using sqlite.
	pass := strings.TrimSpace(v.GetString(KEY_DB_PASSWORD))

	if driver != driverSQLite {
		v.RequireWhen(true, KEY_DB_HOST, "required when `DB_DRIVER` is not sqlite", host)
		v.RequireWhen(true, KEY_DB_PORT, "required when `DB_DRIVER` is not sqlite", port)
		v.RequireWhen(true, KEY_DB_USERNAME, "required when `DB_DRIVER` is not sqlite", user)
		v.RequireWhen(true, KEY_DB_PASSWORD, "required when `DB_DRIVER` is not sqlite", pass)
	}

	cfg.SetDatabaseDriver(driver)
	cfg.SetDatabaseHost(host)
	cfg.SetDatabasePort(port)
	cfg.SetDatabaseName(name)
	cfg.SetDatabaseUsername(user)
	cfg.SetDatabasePassword(pass)
	cfg.SetDatabaseSSLMode("require")
}
