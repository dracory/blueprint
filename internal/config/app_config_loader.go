package config

import "strings"

// loadAppConfig loads application configuration directly into the config.
func loadAppConfig(cfg ConfigInterface, v *envValidator) {
	// Application Name
	//
	// This value is the name of your application, used in notifications
	// and other places where the application name is displayed.
	cfg.SetAppName(v.GetString(KEY_APP_NAME))

	// Application URL
	//
	// The base URL of your application, used for generating links.
	// Example: https://example.com
	cfg.SetAppUrl(v.GetString(KEY_APP_URL))

	// Application Host
	//
	// The host address the server will listen on.
	// Example: 0.0.0.0 (all interfaces) or 127.0.0.1 (localhost only)
	cfg.SetAppHost(v.MustString(KEY_APP_HOST, "set the application host address"))

	// Application Port
	//
	// The port the server will listen on.
	// Example: 8080
	cfg.SetAppPort(v.MustString(KEY_APP_PORT, "set the application port"))

	// Application Environment
	//
	// Determines the environment your application is running in.
	// Supported values: local, development, staging, testing, production
	cfg.SetAppEnv(v.MustString(KEY_APP_ENVIRONMENT, "set the application environment"))

	// Application Debug Mode
	//
	// When enabled, detailed error messages are shown. Disable in production.
	cfg.SetAppDebug(v.GetBool(KEY_APP_DEBUG))

	// CMS MCP API Key
	//
	// API key for the CMS Model Context Protocol integration.
	// Leave empty to disable CMS MCP integration.
	cfg.SetCmsMcpApiKey(strings.TrimSpace(v.GetString(KEY_MCP_API_KEY)))
}
