package config

import (
	"context"
	"errors"
)

type configKey struct{}

func GetConfig(ctx context.Context) (*Config, error) {
	config, ok := ctx.Value(configKey{}).(*Config)
	if !ok {
		return nil, errors.New("config not found in context")
	}
	return config, nil
}

func SetConfig(ctx context.Context, config *Config) context.Context {
	return context.WithValue(ctx, configKey{}, config)
}
