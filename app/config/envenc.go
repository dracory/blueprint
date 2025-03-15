package config

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strings"

	"github.com/dracory/base/env"
)

// EnvEncryptionKey returns the encryption key for environment variables
func EnvEncryptionKey() string {
	return env.ValueOrDefault("ENV_ENCRYPTION_KEY", "")
}

// EnvEncryptionKeyHashed returns the hashed encryption key for environment variables
func EnvEncryptionKeyHashed() string {
	key := EnvEncryptionKey()
	if key == "" {
		return ""
	}

	// Create a new SHA-256 hash
	hasher := sha256.New()
	hasher.Write([]byte(key))

	// Get the hashed value as a hex string
	return hex.EncodeToString(hasher.Sum(nil))
}

// EnvEncryptionKeyDerived returns a derived encryption key based on the environment
// This is useful for having different encryption keys for different environments
func EnvEncryptionKeyDerived() string {
	key := EnvEncryptionKey()
	if key == "" {
		return ""
	}

	// Get the environment
	appEnv := strings.ToLower(os.Getenv("APP_ENV"))
	if appEnv == "" {
		appEnv = "development"
	}

	// Create a new SHA-256 hash with the environment as a salt
	hasher := sha256.New()
	hasher.Write([]byte(key + ":" + appEnv))

	// Get the hashed value as a hex string
	return hex.EncodeToString(hasher.Sum(nil))
}
