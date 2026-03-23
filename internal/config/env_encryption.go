package config

import (
	"fmt"
	"strings"

	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// envEncryptionConfig captures optional environment encryption key usage state.
// It manages the configuration for environment variable encryption/decryption
// using the envenc package for secure configuration management.
type envEncryptionConfig struct {
	privateKey string // Private key for decrypting encrypted environment variables
	used       bool   // Flag indicating if environment encryption is enabled
}

// loadEnvEncryptionConfig loads environment encryption configuration.
// It validates the private key when encryption is enabled and ensures
// proper security requirements are met before allowing encrypted environment
// variable usage.
//
// Parameters:
//   - acc: LoadAccumulator for collecting validation errors and required field checks
//
// Returns:
//   - envEncryptionConfig: Populated configuration struct with encryption settings
func loadEnvEncryptionConfig(acc *baseCfg.LoadAccumulator) envEncryptionConfig {
	used := env.GetBool(KEY_ENVENC_USED)

	if !used {
		return envEncryptionConfig{privateKey: "", used: used}
	}

	privateKey := strings.TrimSpace(env.GetString(KEY_ENVENC_KEY_PRIVATE))

	if err := baseCfg.EnsureRequired(privateKey, KEY_ENVENC_KEY_PRIVATE, "required when ENVENC_USED is yes"); err != nil {
		acc.Add(err)
		return envEncryptionConfig{used: used}
	}

	if privateKey == "" {
		acc.Add(fmt.Errorf("private key is required when env encryption is enabled"))
		return envEncryptionConfig{used: used, privateKey: privateKey}
	}

	return envEncryptionConfig{privateKey: privateKey, used: used}
}
