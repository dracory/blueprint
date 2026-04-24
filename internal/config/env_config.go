package config

import "github.com/dracory/env"

// Package env provides option for encrypted environment variables using EnvEnc
//
// EnvEnc is a tool for storing environment variables in encrypted envenc vault
// files, which can be committed to version control and decrypted at runtime.
//
// Important Security Information:
// - A public key, combined with a corresponding private key, is used
//   as input to a secure one-way hashing function to derive the final
//   encryption key.
// - Both the private and public keys must be at least 32-character strings,
//   composed of randomly generated characters, numbers, and symbols.
// - DO NOT store the actual final key anywhere. It should be generated dynamically when needed.
// - DO NOT directly commit the actual PRIVATE key to version control. Use environment variables or secure configuration management.
// - Replace "YOUR_PUBLIC_KEY" with your actual 32-character public key.
// - The associated private key must be kept extremely secure.
// - Ensure that the random number generator used to create the keys is cryptographically secure (CSPRNG).
// - Ideally, the public key should be obfuscated. See envenc for more details.
//
// Example:
// keyPublic := "aBcD123$456!eFgH789%iJkL0mNoPqRsTuVwXyZ" // Replace with your actual key

// envConfig reads EnvEnc encryption configuration from environment variables.
// EnvEnc is automatically enabled when the required keys are provided.
func envConfig() envSettings {
	// EnvEnc Used
	//
	// Determines whether EnvEnc encryption is enabled.
	// Set to "true" or "1" to enable encryption.
	used := env.GetBool(KEY_ENVENC_USED)

	// EnvEnc Private Key
	//
	// Your private key for EnvEnc vault encryption.
	//
	// See package-level documentation for more information.
	//
	// Security Notes:
	// - Should be loaded from an environment variable or secure configuration.
	// - Should be kept extremely secure and NOT committed to version control.
	// - Should be composed of randomly generated characters, numbers, and symbols.
	// - Should be at least 32 characters long.
	keyPrivate := env.GetString(KEY_ENVENC_KEY_PRIVATE)

	// EnvEnc Public Key
	//
	// Your public key for EnvEnc vault encryption.
	//
	// See package-level documentation for more information.
	//
	// Security Notes:
	// - Should be at least 32 characters long.
	// - Should be composed of randomly generated characters, numbers, and symbols.
	// - Can be committed to version control.
	// - Can be obfuscated for additional security (see envenc for more details).
	keyPublic := "YOUR_PUBLIC_KEY"

	return envSettings{
		used:       used,
		keyPrivate: keyPrivate,
		keyPublic:  keyPublic,
	}
}

type envSettings struct {
	used       bool
	keyPrivate string
	keyPublic  string
}
