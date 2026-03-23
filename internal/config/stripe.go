package config

import "github.com/dracory/env"

// stripeConfig captures Stripe payment integration settings.
// It manages the configuration for Stripe payment processing,
// including API keys and usage detection.
type stripeConfig struct {
	keyPrivate string // Stripe private API key (secret key)
	keyPublic  string // Stripe public API key (publishable key)
	used       bool   // Flag indicating if Stripe integration is enabled
}

// loadStripeConfig loads Stripe configuration from environment variables.
// It reads Stripe API keys and automatically determines if Stripe is enabled
// based on whether both keys are provided.
//
// Returns:
//   - stripeConfig: Populated configuration struct with Stripe settings
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
