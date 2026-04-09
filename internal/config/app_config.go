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

	return appSettings{
		name:         name,
		url:          url,
		host:         host,
		port:         port,
		env:          appEnv,
		debug:        debug,
		cmsMcpApiKey: cmsMcpApiKey,
	}
}

type appSettings struct {
	name         string
	url          string
	host         string
	port         string
	env          string
	debug        bool
	cmsMcpApiKey string
}
