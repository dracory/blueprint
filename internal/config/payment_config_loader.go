package config

import "github.com/dracory/env"

// loadStripeConfig loads Stripe configuration directly into the config.
func loadStripeConfig(cfg ConfigInterface) {
	keyPrivate := env.GetString(KEY_STRIPE_KEY_PRIVATE)
	keyPublic := env.GetString(KEY_STRIPE_KEY_PUBLIC)
	used := keyPrivate != "" && keyPublic != ""

	cfg.SetStripeKeyPrivate(keyPrivate)
	cfg.SetStripeKeyPublic(keyPublic)
	cfg.SetStripeUsed(used)
}
