package config

import (
	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// ============================================================================
// Interface
// ============================================================================

// EncryptionConfigInterface defines encryption configuration methods.
type EncryptionConfigInterface interface {
	SetEnvEncryptionKey(string)
	GetEnvEncryptionKey() string
}

// ============================================================================
// Types
// ============================================================================

// envEncryptionConfig captures environment encryption settings.
type envEncryptionConfig struct {
	used       bool   // Flag indicating if encryption is enabled
	privateKey string // Private key for decryption
}

// ============================================================================
// Loader
// ============================================================================

// loadEnvEncryptionConfig loads encryption configuration from environment variables.
func loadEnvEncryptionConfig(acc *baseCfg.LoadAccumulator) envEncryptionConfig {
	privateKey := env.GetString(KEY_ENVENC_KEY_PRIVATE)
	used := privateKey != ""

	if used {
		acc.MustWhen(true, KEY_ENVENC_KEY_PRIVATE, "required when encryption is enabled", privateKey)
	}

	return envEncryptionConfig{
		used:       used,
		privateKey: privateKey,
	}
}

// ============================================================================
// Implementation (Getters/Setters)
// ============================================================================

func (c *configImplementation) SetEnvEncryptionKey(v string) {
	c.envEncryptionKey = v
}

func (c *configImplementation) GetEnvEncryptionKey() string {
	return c.envEncryptionKey
}
