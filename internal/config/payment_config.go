package config

import "github.com/dracory/env"

// ============================================================================
// Interface
// ============================================================================

// PaymentConfigInterface defines payment provider configuration methods.
type PaymentConfigInterface interface {
	SetStripeKeyPrivate(string)
	GetStripeKeyPrivate() string

	SetStripeKeyPublic(string)
	GetStripeKeyPublic() string

	SetStripeUsed(bool)
	GetStripeUsed() bool
}

// ============================================================================
// Types
// ============================================================================

// stripeConfig captures Stripe payment integration settings.
type stripeConfig struct {
	keyPrivate string // Stripe private API key (secret key)
	keyPublic  string // Stripe public API key (publishable key)
	used       bool   // Flag indicating if Stripe integration is enabled
}

// ============================================================================
// Loader
// ============================================================================

// loadStripeConfig loads Stripe configuration from environment variables.
func loadStripeConfig() stripeConfig {
	keyPrivate := env.GetString(KEY_STRIPE_KEY_PRIVATE)
	keyPublic := env.GetString(KEY_STRIPE_KEY_PUBLIC)
	used := keyPrivate != "" && keyPublic != ""

	return stripeConfig{
		keyPrivate: keyPrivate,
		keyPublic:  keyPublic,
		used:       used,
	}
}

// ============================================================================
// Implementation (Getters/Setters)
// ============================================================================

func (c *configImplementation) SetStripeKeyPrivate(v string) {
	c.stripeKeyPrivate = v
}

func (c *configImplementation) GetStripeKeyPrivate() string {
	return c.stripeKeyPrivate
}

func (c *configImplementation) SetStripeKeyPublic(v string) {
	c.stripeKeyPublic = v
}

func (c *configImplementation) GetStripeKeyPublic() string {
	return c.stripeKeyPublic
}

func (c *configImplementation) SetStripeUsed(v bool) {
	c.stripeUsed = v
}

func (c *configImplementation) GetStripeUsed() bool {
	return c.stripeUsed
}
