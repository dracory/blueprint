package config

import (
	"strings"

	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// ============================================================================
// Interface
// ============================================================================

// AppConfigInterface defines application-level configuration methods.
type AppConfigInterface interface {
	SetAppName(string)
	GetAppName() string

	SetAppType(string)
	GetAppType() string

	SetAppEnv(string)
	GetAppEnv() string

	SetAppHost(string)
	GetAppHost() string

	SetAppPort(string)
	GetAppPort() string

	SetAppUrl(string)
	GetAppUrl() string

	SetAppDebug(bool)
	GetAppDebug() bool

	// Environment helpers
	IsEnvDevelopment() bool
	IsEnvLocal() bool
	IsEnvProduction() bool
	IsEnvStaging() bool
	IsEnvTesting() bool
}

// ============================================================================
// Types
// ============================================================================

// appConfig captures application-level settings.
type appConfig struct {
	name         string // Application name identifier
	url          string // Base URL for the application
	host         string // Host address for the server
	port         string // Port number for the server
	env          string // Environment (development, staging, production)
	debug        bool   // Debug mode flag
	cmsMcpApiKey string // CMS MCP API key for integration
}

// ============================================================================
// Loader
// ============================================================================

// loadAppConfig loads application configuration from environment variables.
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

// ============================================================================
// Implementation (Getters/Setters)
// ============================================================================

func (c *configImplementation) SetAppName(appName string) {
	c.appName = appName
}

func (c *configImplementation) GetAppName() string {
	return c.appName
}

func (c *configImplementation) SetAppType(appType string) {
	c.appType = appType
}

func (c *configImplementation) GetAppType() string {
	return c.appType
}

func (c *configImplementation) SetAppEnv(appEnv string) {
	c.appEnv = appEnv
}

func (c *configImplementation) GetAppEnv() string {
	return c.appEnv
}

func (c *configImplementation) SetAppHost(appHost string) {
	c.appHost = appHost
}

func (c *configImplementation) GetAppHost() string {
	return c.appHost
}

func (c *configImplementation) SetAppPort(appPort string) {
	c.appPort = appPort
}

func (c *configImplementation) GetAppPort() string {
	return c.appPort
}

func (c *configImplementation) SetAppUrl(appUrl string) {
	c.appUrl = appUrl
}

func (c *configImplementation) GetAppUrl() string {
	return c.appUrl
}

func (c *configImplementation) SetAppDebug(appDebug bool) {
	c.appDebug = appDebug
}

func (c *configImplementation) GetAppDebug() bool {
	return c.appDebug
}

// Environment Helpers
func (c *configImplementation) IsEnvDevelopment() bool {
	return c.appEnv == "development"
}

func (c *configImplementation) IsEnvLocal() bool {
	return c.appEnv == "local"
}

func (c *configImplementation) IsEnvProduction() bool {
	return c.appEnv == "production"
}

func (c *configImplementation) IsEnvStaging() bool {
	return c.appEnv == "staging"
}

func (c *configImplementation) IsEnvTesting() bool {
	return c.appEnv == "testing"
}
