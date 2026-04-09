package config

// databaseConfig reads database configuration from environment variables.
func databaseConfig(env *envValidator) databaseSettings {
	// Database Driver
	//
	// The database driver to use for the application.
	// Supported values: sqlite, postgres, mysql
	driver := env.GetStringOrError(KEY_DB_DRIVER, "select the database driver (e.g., sqlite, postgres)")

	// Database Host
	//
	// The hostname or IP address of the database server.
	// Not required when using sqlite.
	host := env.GetString(KEY_DB_HOST)

	// Database Port
	//
	// The port the database server is listening on.
	// Common defaults: postgres=5432, mysql=3306
	// Not required when using sqlite.
	port := env.GetString(KEY_DB_PORT)

	// Database Name
	//
	// The name of the database to connect to.
	// For sqlite, this is the file path (e.g., ./database.db or :memory:)
	name := env.GetStringOrError(KEY_DB_DATABASE, "set the database name")

	// Database Username
	//
	// The username for authenticating with the database server.
	// Not required when using sqlite.
	user := env.GetString(KEY_DB_USERNAME)

	// Database Password
	//
	// The password for authenticating with the database server.
	// Not required when using sqlite.
	pass := env.GetString(KEY_DB_PASSWORD)

	if driver != driverSQLite {
		env.RequireWhen(true, KEY_DB_HOST, "required when `DB_DRIVER` is not sqlite", host)
		env.RequireWhen(true, KEY_DB_PORT, "required when `DB_DRIVER` is not sqlite", port)
		env.RequireWhen(true, KEY_DB_USERNAME, "required when `DB_DRIVER` is not sqlite", user)
		env.RequireWhen(true, KEY_DB_PASSWORD, "required when `DB_DRIVER` is not sqlite", pass)
	}

	return databaseSettings{driver: driver, host: host, port: port, name: name, user: user, pass: pass}
}

type databaseSettings struct {
	driver string
	host   string
	port   string
	name   string
	user   string
	pass   string
}
