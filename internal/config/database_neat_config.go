package config

import (
	"strconv"
	"strings"

	"github.com/dracory/neat/database/db"
)

// DatabaseNeatConfig maps the blueprint configuration to a neat DBConfig.
// It builds the default connection from the existing single-database settings
// and applies the connection pool configuration.
func DatabaseNeatConfig(cfg ConfigInterface) db.DBConfig {
	if cfg == nil {
		return db.DBConfig{}
	}

	defaultConnection := cfg.GetDatabaseDefaultConnection()
	if defaultConnection == "" {
		defaultConnection = "default"
	}

	connections := make(map[string]db.ConnectionConfig)
	for _, conn := range cfg.GetDatabaseConnections() {
		if conn == nil {
			continue
		}
		connections[conn.GetName()] = connectionNeatConfig(conn)
	}

	// Ensure the default connection is always present.
	if _, ok := connections[defaultConnection]; !ok {
		connections[defaultConnection] = connectionNeatConfig(&databaseConnectionSettings{
			name:     defaultConnection,
			driver:   cfg.GetDatabaseDriver(),
			host:     cfg.GetDatabaseHost(),
			port:     cfg.GetDatabasePort(),
			database: cfg.GetDatabaseName(),
			username: cfg.GetDatabaseUsername(),
			password: cfg.GetDatabasePassword(),
			sslMode:  cfg.GetDatabaseSSLMode(),
			charset:  cfg.GetDatabaseCharset(),
			timezone: cfg.GetDatabaseTimezone(),
			dsn:      cfg.GetDatabaseDSN(),
			prefix:   cfg.GetDatabasePrefix(),
		})
	}

	pool := db.PoolConfig{
		MaxOpenConns:    cfg.GetDatabaseMaxOpenConns(),
		MaxIdleConns:    cfg.GetDatabaseMaxIdleConns(),
		ConnMaxLifetime: cfg.GetDatabaseConnMaxLifetimeSeconds(),
		ConnMaxIdleTime: cfg.GetDatabaseConnMaxIdleTimeSeconds(),
		QueryTimeout:    30,
	}

	return db.DBConfig{
		Default:     defaultConnection,
		Connections: connections,
		Pool:        pool,
	}
}

// connectionNeatConfig maps a DatabaseConnectionConfigInterface to a neat
// ConnectionConfig. It converts string port values to integers and applies
// driver-specific defaults.
func connectionNeatConfig(conn DatabaseConnectionConfigInterface) db.ConnectionConfig {
	if conn == nil {
		return db.ConnectionConfig{}
	}

	driver := strings.ToLower(strings.TrimSpace(conn.GetDriver()))

	nc := db.ConnectionConfig{
		Driver:   driver,
		Dsn:      conn.GetDSN(),
		Host:     conn.GetHost(),
		Database: conn.GetDatabase(),
		Username: conn.GetUsername(),
		Password: conn.GetPassword(),
		Charset:  conn.GetCharset(),
		SSLMode:  conn.GetSSLMode(),
		Timezone: conn.GetTimezone(),
		Prefix:   conn.GetPrefix(),
		Port:     portToInt(conn.GetPort(), driver),
	}

	return nc
}

// portToInt converts a string port to an integer with driver-specific defaults.
func portToInt(port, driver string) int {
	port = strings.TrimSpace(port)
	if port == "" {
		switch driver {
		case "mysql":
			return 3306
		case "postgres":
			return 5432
		case "sqlserver":
			return 1433
		case "oracle":
			return 1521
		default:
			return 0
		}
	}

	v, err := strconv.Atoi(port)
	if err != nil {
		return 0
	}
	return v
}
