package config

import (
	"project/internal/resources"

	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// envValidator is a local alias for env.Validator for use in config loaders.
type envValidator = env.Validator

func Load() (ConfigInterface, error) {
	env.Load(".env")

	v := &envValidator{}
	cfg := &configImplementation{}

	// Load app config first to get app.env
	readAppConfig(cfg, v)

	// Load encryption config and check if encryption is used
	privateKey := env.GetString(KEY_ENVENC_KEY_PRIVATE)
	encryptionUsed := privateKey != ""

	if encryptionUsed {
		v.RequireWhen(true, KEY_ENVENC_KEY_PRIVATE,
			"required when encryption is enabled", privateKey)
	}

	cfg.SetEnvEncryptionKey(privateKey)

	// Initialize encrypted environment variables BEFORE other config loaders read them
	if encryptionUsed {
		if err := baseCfg.InitializeEnvEncVariablesFromResources(cfg.GetAppEnv(), ENVENC_KEY_PUBLIC, privateKey, resources.Resource); err != nil {
			v.Add(err)
		} else {
			cfg.SetEnvEncryptionKey("removed") // reset the private key
		}

		// Reload app config to pick up any encrypted app variables
		readAppConfig(cfg, v)
	}

	// Now load remaining config sections - they will have access to encrypted variables
	readDatabaseConfig(cfg, v)
	readMailConfig(cfg)
	readRegistrationConfig(cfg)
	readStoresConfig(cfg, v)
	readStripeConfig(cfg)
	readLLMConfig(cfg, v)
	readTranslationConfig(cfg)
	if err := v.Err(); err != nil {
		return nil, err
	}

	return cfg, nil
}
