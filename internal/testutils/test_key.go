package testutils

import (
	"project/internal/types"

	"github.com/dracory/test"
)

// TestKey is a pseudo secret test key used for testing specific unit cases
//
//	where a secret key is required but not available in the testing environment
func TestKey(cfg types.ConfigInterface) string {
	// Use the base testutils package's TestKey function
	return test.TestKey(
		cfg.GetDatabaseDriver(),
		cfg.GetDatabaseHost(),
		cfg.GetDatabasePort(),
		cfg.GetDatabaseName(),
		cfg.GetDatabaseUsername(),
		cfg.GetDatabasePassword(),
	)
}
