package config

import "github.com/dracory/env"

// registrationConfig captures authentication registration toggle.
// It controls whether user registration is enabled or disabled
// in the authentication system.
type registrationConfig struct {
	enabled bool // Flag indicating if user registration is allowed
}

// loadRegistrationConfig loads registration configuration from environment variables.
// It reads the registration enabled flag and returns a populated registrationConfig.
//
// Returns:
//   - registrationConfig: Populated configuration struct with registration settings
func loadRegistrationConfig() registrationConfig {
	return registrationConfig{
		enabled: env.GetBool(KEY_AUTH_REGISTRATION_ENABLED),
	}
}
