package config

import "github.com/dracory/env"

// authConfig reads authentication configuration from environment variables.
func authConfig() authSettings {
	// User Registration
	//
	// Controls whether new users can register for an account.
	// Set to false to disable public registration (invite-only or closed systems).
	registrationEnabled := env.GetBool(KEY_AUTH_REGISTRATION_ENABLED)

	return authSettings{registrationEnabled: registrationEnabled}
}

type authSettings struct {
	registrationEnabled bool
}
