package config

import "github.com/dracory/env"

// ============================================================================
// Interface
// ============================================================================

// AuthConfigInterface defines authentication configuration methods.
type AuthConfigInterface interface {
	SetRegistrationEnabled(bool)
	GetRegistrationEnabled() bool
}

// ============================================================================
// Types
// ============================================================================

// registrationConfig captures authentication registration toggle.
type registrationConfig struct {
	enabled bool // Flag indicating if user registration is allowed
}

// ============================================================================
// Loader
// ============================================================================

// loadRegistrationConfig loads registration configuration from environment variables.
func loadRegistrationConfig() registrationConfig {
	return registrationConfig{
		enabled: env.GetBool(KEY_AUTH_REGISTRATION_ENABLED),
	}
}

// ============================================================================
// Implementation (Getters/Setters)
// ============================================================================

func (c *configImplementation) SetRegistrationEnabled(v bool) {
	c.registrationEnabled = v
}

func (c *configImplementation) GetRegistrationEnabled() bool {
	return c.registrationEnabled
}
