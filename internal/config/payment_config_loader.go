package config

import "github.com/dracory/env"

// stripeConfig captures Stripe payment integration settings.
type stripeConfig struct {
	keyPrivate string // Stripe private API key (secret key)
	keyPublic  string // Stripe public API key (publishable key)
	used       bool   // Flag indicating if Stripe integration is enabled
}

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
