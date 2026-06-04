package file_manager

import (
	"project/internal/config"
	"project/internal/app"

	_ "modernc.org/sqlite"
)

func setupTestRegistry() (app.AppInterface, func()) {
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

	reg, err := app.New(cfg)
	if err != nil {
		panic("failed to create app: " + err.Error())
	}

	cleanup := func() {
		if reg != nil {
			reg.Close()
		}
	}

	return reg, cleanup
}
