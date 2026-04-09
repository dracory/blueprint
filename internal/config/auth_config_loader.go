package config

import "github.com/dracory/env"

// loadRegistrationConfig loads authentication configuration directly into the config.
func loadRegistrationConfig(cfg ConfigInterface) {
	// User Registration
	//
	// Controls whether new users can register for an account.
	// Set to false to disable public registration (invite-only or closed systems).
	registrationEnabled := env.GetBool(KEY_AUTH_REGISTRATION_ENABLED)

	cfg.SetRegistrationEnabled(registrationEnabled)
}
