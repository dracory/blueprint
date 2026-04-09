package config

import "strings"

// readAppConfig reads application configuration from environment variables.
func readAppConfig(cfg *configImplementation, v *envValidator) {
	// Application Name
	//
	// This value is the name of your application, used in notifications
	// and other places where the application name is displayed.
	name := v.GetStringOrDefault(KEY_APP_NAME, "Blueprint")

	// Application URL
	//
	// The base URL of your application, used for generating links.
	// Example: https://example.com
	url := v.GetStringOrDefault(KEY_APP_URL, "http://localhost:8080")

	// Application Host
	//
	// The host address the server will listen on.
	// Example: 0.0.0.0 (all interfaces) or 127.0.0.1 (localhost only)
	host := v.GetStringOrError(KEY_APP_HOST, "set the application host address")

	// Application Port
	//
	// The port the server will listen on.
	// Example: 8080
	port := v.GetStringOrError(KEY_APP_PORT, "set the application port")

	// Application Environment
	//
	// Determines the environment your application is running in.
	// Supported values: local, development, staging, testing, production
	appEnv := v.GetStringOrError(KEY_APP_ENVIRONMENT, "set the application environment")

	// Application Debug Mode
	//
	// When enabled, detailed error messages are shown. Disable in production.
	debug := v.GetBool(KEY_APP_DEBUG)

	// CMS MCP API Key
	//
	// API key for the CMS Model Context Protocol integration.
	// Leave empty to disable CMS MCP integration.
	cmsMcpApiKey := strings.TrimSpace(v.GetString(KEY_MCP_API_KEY))

	cfg.setAppConfig(name, url, host, port, appEnv, debug, cmsMcpApiKey)
}
