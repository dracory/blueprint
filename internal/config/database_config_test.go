package config_test

import (
	"testing"

	"project/internal/config"

	"github.com/dracory/neat/database/db"
)

func TestDatabaseNeatConfig_MapsDefaultConnection(t *testing.T) {
	cfg := config.New()
	cfg.SetDatabaseDefaultConnection("default")
	cfg.SetDatabaseDriver("sqlite")
	cfg.SetDatabaseHost("")
	cfg.SetDatabasePort("")
	cfg.SetDatabaseName("database.db")
	cfg.SetDatabaseUsername("")
	cfg.SetDatabasePassword("")
	cfg.SetDatabaseCharset("utf8mb4")
	cfg.SetDatabaseTimezone("UTC")
	cfg.SetDatabaseMaxOpenConns(1)
	cfg.SetDatabaseMaxIdleConns(1)
	cfg.SetDatabaseConnMaxLifetimeSeconds(30)
	cfg.SetDatabaseConnMaxIdleTimeSeconds(5)

	neatCfg := config.DatabaseNeatConfig(cfg)

	if neatCfg.Default != "default" {
		t.Errorf("expected default connection 'default', got %q", neatCfg.Default)
	}

	conn, ok := neatCfg.Connections["default"]
	if !ok {
		t.Fatal("expected 'default' connection to exist")
	}

	if conn.Driver != "sqlite" {
		t.Errorf("expected driver 'sqlite', got %q", conn.Driver)
	}
	if conn.Database != "database.db" {
		t.Errorf("expected database 'database.db', got %q", conn.Database)
	}
	if conn.Host != "" {
		t.Errorf("expected empty host, got %q", conn.Host)
	}
	if conn.Port != 0 {
		t.Errorf("expected port 0 for sqlite, got %d", conn.Port)
	}
	if conn.Charset != "utf8mb4" {
		t.Errorf("expected charset 'utf8mb4', got %q", conn.Charset)
	}
	if conn.Timezone != "UTC" {
		t.Errorf("expected timezone 'UTC', got %q", conn.Timezone)
	}
	if conn.SSLMode != "" {
		t.Errorf("expected empty sslmode for sqlite, got %q", conn.SSLMode)
	}

	if neatCfg.Pool.MaxOpenConns != 1 {
		t.Errorf("expected max open conns 1, got %d", neatCfg.Pool.MaxOpenConns)
	}
	if neatCfg.Pool.MaxIdleConns != 1 {
		t.Errorf("expected max idle conns 1, got %d", neatCfg.Pool.MaxIdleConns)
	}
	if neatCfg.Pool.ConnMaxLifetime != 30 {
		t.Errorf("expected conn max lifetime 30, got %d", neatCfg.Pool.ConnMaxLifetime)
	}
	if neatCfg.Pool.ConnMaxIdleTime != 5 {
		t.Errorf("expected conn max idle time 5, got %d", neatCfg.Pool.ConnMaxIdleTime)
	}
	if neatCfg.Pool.QueryTimeout != 30 {
		t.Errorf("expected query timeout 30, got %d", neatCfg.Pool.QueryTimeout)
	}
}

func TestDatabaseNeatConfig_MapsPostgresConnection(t *testing.T) {
	cfg := config.New()
	cfg.SetDatabaseDefaultConnection("default")
	cfg.SetDatabaseDriver("postgres")
	cfg.SetDatabaseHost("localhost")
	cfg.SetDatabasePort("5432")
	cfg.SetDatabaseName("blueprint")
	cfg.SetDatabaseUsername("user")
	cfg.SetDatabasePassword("pass")
	cfg.SetDatabaseSSLMode("require")
	cfg.SetDatabaseCharset("utf8mb4")
	cfg.SetDatabaseTimezone("UTC")

	neatCfg := config.DatabaseNeatConfig(cfg)

	conn := neatCfg.Connections["default"]
	if conn.Driver != "postgres" {
		t.Errorf("expected driver 'postgres', got %q", conn.Driver)
	}
	if conn.Host != "localhost" {
		t.Errorf("expected host 'localhost', got %q", conn.Host)
	}
	if conn.Port != 5432 {
		t.Errorf("expected port 5432, got %d", conn.Port)
	}
	if conn.Database != "blueprint" {
		t.Errorf("expected database 'blueprint', got %q", conn.Database)
	}
	if conn.Username != "user" {
		t.Errorf("expected username 'user', got %q", conn.Username)
	}
	if conn.Password != "pass" {
		t.Errorf("expected password 'pass', got %q", conn.Password)
	}
	if conn.SSLMode != "require" {
		t.Errorf("expected sslmode 'require', got %q", conn.SSLMode)
	}
}

func TestDatabaseNeatConfig_AppliesPortDefaults(t *testing.T) {
	cfg := config.New()
	cfg.SetDatabaseDefaultConnection("default")
	cfg.SetDatabaseDriver("mysql")
	cfg.SetDatabaseHost("localhost")
	cfg.SetDatabasePort("")
	cfg.SetDatabaseName("blueprint")
	cfg.SetDatabaseUsername("user")
	cfg.SetDatabasePassword("pass")

	neatCfg := config.DatabaseNeatConfig(cfg)
	conn := neatCfg.Connections["default"]
	if conn.Port != 3306 {
		t.Errorf("expected default mysql port 3306, got %d", conn.Port)
	}
}

func TestDatabaseNeatConfig_KeepsDSNAndPrefix(t *testing.T) {
	cfg := config.New()
	cfg.SetDatabaseDefaultConnection("default")
	cfg.SetDatabaseDriver("postgres")
	cfg.SetDatabaseName("blueprint")
	cfg.SetDatabaseDSN("postgres://user:pass@localhost:5432/blueprint?sslmode=require")
	cfg.SetDatabasePrefix("bp_")

	neatCfg := config.DatabaseNeatConfig(cfg)
	conn := neatCfg.Connections["default"]
	if conn.Dsn != "postgres://user:pass@localhost:5432/blueprint?sslmode=require" {
		t.Errorf("expected dsn to be preserved, got %q", conn.Dsn)
	}
	if conn.Prefix != "bp_" {
		t.Errorf("expected prefix 'bp_', got %q", conn.Prefix)
	}
}

func TestDatabaseNeatConfig_NilConfig(t *testing.T) {
	var cfg config.ConfigInterface
	neatCfg := config.DatabaseNeatConfig(cfg)
	if neatCfg.Default != "" {
		t.Errorf("expected empty default for nil config, got %q", neatCfg.Default)
	}
	if len(neatCfg.Connections) != 0 {
		t.Errorf("expected no connections for nil config, got %d", len(neatCfg.Connections))
	}
}

func TestDatabaseNeatConfig_DefaultConnectionFallback(t *testing.T) {
	cfg := config.New()
	cfg.SetDatabaseDriver("sqlite")
	cfg.SetDatabaseName(":memory:")

	neatCfg := config.DatabaseNeatConfig(cfg)
	if neatCfg.Default != "default" {
		t.Errorf("expected default connection name 'default', got %q", neatCfg.Default)
	}
	if _, ok := neatCfg.Connections["default"]; !ok {
		t.Fatal("expected fallback default connection to be created")
	}
}

func TestDatabaseNeatConfig_Validate(t *testing.T) {
	cfg := config.New()
	cfg.SetDatabaseDefaultConnection("default")
	cfg.SetDatabaseDriver("sqlite")
	cfg.SetDatabaseName(":memory:")

	neatCfg := config.DatabaseNeatConfig(cfg)
	if err := neatCfg.Validate(); err != nil {
		t.Errorf("expected sqlite config to validate, got: %v", err)
	}
}

func TestDatabaseNeatConfig_ConnectionValidate(t *testing.T) {
	conn := db.ConnectionConfig{
		Driver:   "sqlite",
		Database: ":memory:",
	}
	if err := conn.Validate(); err != nil {
		t.Errorf("expected sqlite connection to validate, got: %v", err)
	}
}
