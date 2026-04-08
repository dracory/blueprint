package config

import "github.com/dracory/env"

// loadRegistrationConfig loads registration configuration directly into the config.
func loadRegistrationConfig(cfg ConfigInterface) {
	cfg.SetRegistrationEnabled(env.GetBool(KEY_AUTH_REGISTRATION_ENABLED))
}
