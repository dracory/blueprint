package testutils

import (
	"project/config"

	"github.com/dracory/base/testutils"
)

// TestKey is a pseudo secret test key used for testing specific unit cases
//
//	where a secret key is required but not available in the testing environment
func TestKey() string {
	// Use the base testutils package's TestKey function
	return testutils.TestKey(config.DbDriver, config.DbHost, config.DbPort, config.DbName, config.DbUser, config.DbPass)
}
