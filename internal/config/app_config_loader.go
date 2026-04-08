package config

import "strings"

// loadAppConfig loads application configuration directly into the config.
func loadAppConfig(cfg ConfigInterface, v *envValidator) {
	cfg.SetAppName(v.GetString(KEY_APP_NAME))
	cfg.SetAppUrl(v.GetString(KEY_APP_URL))
	cfg.SetAppHost(v.MustString(KEY_APP_HOST, "set the application host address"))
	cfg.SetAppPort(v.MustString(KEY_APP_PORT, "set the application port"))
	cfg.SetAppEnv(v.MustString(KEY_APP_ENVIRONMENT, "set the application environment"))
	cfg.SetAppDebug(v.GetBool(KEY_APP_DEBUG))
	cfg.SetCmsMcpApiKey(strings.TrimSpace(v.GetString(KEY_MCP_API_KEY)))
}
