package config

import "github.com/dracory/env"

// registrationConfig captures authentication registration toggle.
type registrationConfig struct {
	enabled bool // Flag indicating if user registration is allowed
}

// loadRegistrationConfig loads registration configuration from environment variables.
func loadRegistrationConfig() registrationConfig {
	return registrationConfig{
		enabled: env.GetBool(KEY_AUTH_REGISTRATION_ENABLED),
	}
}
