package config

import (
	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// envEncryptionConfig captures environment encryption settings.
type envEncryptionConfig struct {
	used       bool   // Flag indicating if encryption is enabled
	privateKey string // Private key for decryption
}

// loadEnvEncryptionConfig loads encryption configuration from environment variables.
func loadEnvEncryptionConfig(acc *baseCfg.LoadAccumulator) envEncryptionConfig {
	privateKey := env.GetString(KEY_ENVENC_KEY_PRIVATE)
	used := privateKey != ""

	if used {
		acc.MustWhen(true, KEY_ENVENC_KEY_PRIVATE,
			"required when encryption is enabled", privateKey)
	}

	return envEncryptionConfig{
		used:       used,
		privateKey: privateKey,
	}
}
