package config

import "github.com/dracory/env"

// readStripeConfig reads Stripe payment configuration from environment variables.
// Stripe is automatically enabled when both keys are provided.
func paymentConfig(cfg *configImplementation) {
	// Stripe Private Key
	//
	// Your Stripe secret key, used for server-side API calls.
	// Find it at: https://dashboard.stripe.com/apikeys
	// Example: sk_live_... (production) or sk_test_... (testing)
	keyPrivate := env.GetString(KEY_STRIPE_KEY_PRIVATE)

	// Stripe Public Key
	//
	// Your Stripe publishable key, used in client-side code.
	// Find it at: https://dashboard.stripe.com/apikeys
	// Example: pk_live_... (production) or pk_test_... (testing)
	keyPublic := env.GetString(KEY_STRIPE_KEY_PUBLIC)

	// Stripe is automatically enabled when both keys are present.
	used := keyPrivate != "" && keyPublic != ""

	// -------------------------------------------------------------------------
	// Do not edit below this line
	// -------------------------------------------------------------------------
	cfg.setStripeConfig(keyPrivate, keyPublic, used)
}
