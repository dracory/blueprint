package config

import (
	"strings"

	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// appConfig captures application-level settings.
// It includes basic application configuration like name, URL, host, port,
// environment, debug mode, and CMS MCP API key for integration.
type appConfig struct {
	name         string // Application name identifier
	url          string // Base URL for the application
	host         string // Host address for the server
	port         string // Port number for the server
	env          string // Environment (development, staging, production)
	debug        bool   // Debug mode flag
	cmsMcpApiKey string // CMS MCP API key for integration
}

// loadAppConfig loads application configuration from environment variables.
// It validates required fields and returns a populated appConfig struct.
// The MCP API key is trimmed of whitespace for security.
//
// Parameters:
//   - acc: LoadAccumulator for collecting validation errors and required field checks
//
// Returns:
//   - appConfig: Populated configuration struct with application settings
func loadAppConfig(acc *baseCfg.LoadAccumulator) appConfig {
	mcpApiKey := strings.TrimSpace(env.GetString(KEY_MCP_API_KEY))

	return appConfig{
		name:         env.GetString(KEY_APP_NAME),
		url:          env.GetString(KEY_APP_URL),
		host:         acc.MustString(KEY_APP_HOST, "set the application host address"),
		port:         acc.MustString(KEY_APP_PORT, "set the application port"),
		env:          acc.MustString(KEY_APP_ENVIRONMENT, "set the application environment"),
		debug:        env.GetBool(KEY_APP_DEBUG),
		cmsMcpApiKey: mcpApiKey,
	}
}
