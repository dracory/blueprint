package testutils

import (
	"project/internal/config"

	"github.com/dracory/base/test"
)

// TestKey is a pseudo secret test key used for testing specific unit cases
//
//	where a secret key is required but not available in the testing environment
func TestKey() string {
	// Use the base testutils package's TestKey function
	return test.TestKey(config.DbDriver, config.DbHost, config.DbPort, config.DbName, config.DbUser, config.DbPass)
}
