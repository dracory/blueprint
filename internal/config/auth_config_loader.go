package config

import "github.com/dracory/env"

// readRegistrationConfig reads authentication configuration from environment variables.
func readRegistrationConfig(cfg *configImplementation) {
	// User Registration
	//
	// Controls whether new users can register for an account.
	// Set to false to disable public registration (invite-only or closed systems).
	registrationEnabled := env.GetBool(KEY_AUTH_REGISTRATION_ENABLED)

	cfg.setAuthConfig(registrationEnabled)
}
