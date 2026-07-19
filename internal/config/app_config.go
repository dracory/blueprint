package config

// appConfig reads application configuration from environment variables.
func appConfig(env *envValidator) appSettings {
	// Application Name
	//
	// This value is the name of your application, used in notifications
	// and other places where the application name is displayed.
	name := env.GetStringOrDefault(KEY_APP_NAME, "Blueprint")

	// Application URL
	//
	// The base URL of your application, used for generating links.
	// Example: https://example.com
	url := env.GetStringOrDefault(KEY_APP_URL, "http://localhost:8080")

	// Application Host
	//
	// The host address the server will listen on.
	// Example: 0.0.0.0 (all interfaces) or 127.0.0.1 (localhost only)
	host := env.GetStringOrError(KEY_APP_HOST, "set the application host address")

	// Application Port
	//
	// The port the server will listen on.
	// Example: 8080
	port := env.GetStringOrError(KEY_APP_PORT, "set the application port")

	// Application Environment
	//
	// Determines the environment your application is running in.
	// Supported values: local, development, staging, testing, production
	appEnv := env.GetStringOrError(KEY_APP_ENVIRONMENT, "set the application environment")

	// Application Debug Mode
	//
	// When enabled, detailed error messages are shown. Disable in production.
	debug := env.GetBool(KEY_APP_DEBUG)

	// CMS MCP API Key
	//
	// API key for the CMS Model Context Protocol integration.
	// Leave empty to disable CMS MCP integration.
	cmsMcpApiKey := env.GetString(KEY_MCP_API_KEY)

	// Maintenance Mode
	//
	// When enabled, the app returns 503 Service Unavailable for all requests.
	// This can also be toggled via CLI: go run ./cmd/server maintenance down/up
	maintenanceEnabled := env.GetBool(KEY_APP_MAINTENANCE_ENABLED)
	maintenanceFilePath := env.GetStringOrDefault(KEY_APP_MAINTENANCE_FILE_PATH, "maintenance_mode_state.json")

	return appSettings{
		name:                name,
		url:                 url,
		host:                host,
		port:                port,
		env:                 appEnv,
		debug:               debug,
		cmsMcpApiKey:        cmsMcpApiKey,
		maintenanceEnabled:  maintenanceEnabled,
		maintenanceFilePath: maintenanceFilePath,
	}
}

type appSettings struct {
	name                string
	url                 string
	host                string
	port                string
	env                 string
	debug               bool
	cmsMcpApiKey        string
	maintenanceEnabled  bool
	maintenanceFilePath string
}
