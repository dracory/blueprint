package config

import "project/internal/types"

// ConfigInterface is the application configuration contract.
//
// Note: This is currently an alias to the implementation in internal/types.
// It exists to migrate callers to internal/config without a large refactor.
// The long-term goal is to move the implementation out of internal/types.
type ConfigInterface = types.ConfigInterface

// Config is the application configuration implementation.
//
// Note: This is currently an alias to the implementation in internal/types.
type Config = types.Config
