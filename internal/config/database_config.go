package config

import "time"

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

	// Connection Pool - Max Open Connections
	//
	// Maximum number of open connections to the database.
	// SQLite should stay at 1 to avoid concurrent write issues.
	// For postgres/mysql, 25 is a reasonable default for most apps.
	maxOpenConns := env.GetIntOrDefault(KEY_DB_MAX_OPEN_CONNS, 25)
	if driver == driverSQLite {
		maxOpenConns = 1
	}

	// Connection Pool - Max Idle Connections
	//
	// Maximum number of idle connections kept in the pool.
	// Should be less than or equal to MaxOpenConns.
	maxIdleConns := env.GetIntOrDefault(KEY_DB_MAX_IDLE_CONNS, 5)
	if driver == driverSQLite {
		maxIdleConns = 1
	}

	// Connection Pool - Max Connection Lifetime
	//
	// Maximum time a connection may be reused. Connections older than this
	// are closed and replaced. 0 means no limit.
	// Unit: seconds. Default: 300 (5 minutes)
	connMaxLifetime := time.Duration(env.GetIntOrDefault(KEY_DB_CONN_MAX_LIFETIME_SECONDS, 300)) * time.Second
	if driver == driverSQLite {
		connMaxLifetime = 30 * time.Second
	}

	// Connection Pool - Max Connection Idle Time
	//
	// Maximum time a connection may be idle before being closed.
	// 0 means no limit.
	// Unit: seconds. Default: 5
	connMaxIdleTime := time.Duration(env.GetIntOrDefault(KEY_DB_CONN_MAX_IDLE_TIME_SECONDS, 5)) * time.Second

	// Database Charset
	//
	// Character set for the database connection. Only used for MySQL.
	// Example: utf8mb4, utf8
	charset := env.GetStringOrDefault(KEY_DB_CHARSET, "utf8mb4")

	// Database Timezone
	//
	// Timezone for the database connection.
	// Example: UTC, America/New_York, Europe/London
	timezone := env.GetStringOrDefault(KEY_DB_TIMEZONE, "UTC")

	if driver != driverSQLite {
		env.RequireWhen(true, KEY_DB_HOST, "required when `DB_DRIVER` is not sqlite", host)
		env.RequireWhen(true, KEY_DB_PORT, "required when `DB_DRIVER` is not sqlite", port)
		env.RequireWhen(true, KEY_DB_USERNAME, "required when `DB_DRIVER` is not sqlite", user)
		env.RequireWhen(true, KEY_DB_PASSWORD, "required when `DB_DRIVER` is not sqlite", pass)
	}

	return databaseSettings{
		driver:          driver,
		host:            host,
		port:            port,
		name:            name,
		user:            user,
		pass:            pass,
		maxOpenConns:    maxOpenConns,
		maxIdleConns:    maxIdleConns,
		connMaxLifetime: connMaxLifetime,
		connMaxIdleTime: connMaxIdleTime,
		charset:         charset,
		timezone:        timezone,
	}
}

type databaseSettings struct {
	driver          string
	host            string
	port            string
	name            string
	user            string
	pass            string
	maxOpenConns    int
	maxIdleConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
	charset         string
	timezone        string
}
