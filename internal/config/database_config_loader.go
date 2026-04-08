package config

import "strings"

// loadDatabaseConfig loads database configuration directly into the config.
func loadDatabaseConfig(cfg ConfigInterface, v *envValidator) {
	driver := v.MustString(KEY_DB_DRIVER, "select the database driver (e.g., sqlite, postgres)")
	host := strings.TrimSpace(v.GetString(KEY_DB_HOST))
	port := strings.TrimSpace(v.GetString(KEY_DB_PORT))
	name := v.MustString(KEY_DB_DATABASE, "set the database name")
	user := strings.TrimSpace(v.GetString(KEY_DB_USERNAME))
	pass := strings.TrimSpace(v.GetString(KEY_DB_PASSWORD))

	if driver != driverSQLite {
		v.MustWhen(true, KEY_DB_HOST, "required when `DB_DRIVER` is not sqlite", host)
		v.MustWhen(true, KEY_DB_PORT, "required when `DB_DRIVER` is not sqlite", port)
		v.MustWhen(true, KEY_DB_USERNAME, "required when `DB_DRIVER` is not sqlite", user)
		v.MustWhen(true, KEY_DB_PASSWORD, "required when `DB_DRIVER` is not sqlite", pass)
	}

	cfg.SetDatabaseDriver(driver)
	cfg.SetDatabaseHost(host)
	cfg.SetDatabasePort(port)
	cfg.SetDatabaseName(name)
	cfg.SetDatabaseUsername(user)
	cfg.SetDatabasePassword(pass)
	cfg.SetDatabaseSSLMode("require")
}
