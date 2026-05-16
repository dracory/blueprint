package file_manager

import (
	"project/internal/config"
	"project/internal/registry"

	_ "modernc.org/sqlite"
)

func setupTestRegistry() (registry.RegistryInterface, func()) {
	cfg := config.New()
	cfg.SetAppEnv("testing")
	cfg.SetAppDebug(true)
	cfg.SetAppName("Test file manager")
	cfg.SetDatabaseDriver("sqlite")
	cfg.SetDatabaseHost("")
	cfg.SetDatabasePort("")
	cfg.SetDatabaseName("file:test_db?mode=memory&cache=shared")
	cfg.SetDatabaseUsername("")
	cfg.SetDatabasePassword("")
	cfg.SetRegistrationEnabled(true)
	cfg.SetMailFromAddress("test@test.com")
	cfg.SetMailFromName("TestName")
	cfg.SetSqlFileStoreUsed(true)
	cfg.SetMediaRoot("/uploads")

	reg, err := registry.New(cfg)
	if err != nil {
		panic("failed to create registry: " + err.Error())
	}

	cleanup := func() {
		if reg != nil {
			reg.Close()
		}
	}

	return reg, cleanup
}
