package testutils

import (
	"fmt"
	"log/slog"
	"project/internal/app"
	"project/internal/types"
	"time"

	//smtpmock "github.com/mocktools/go-smtp-mock"

	_ "modernc.org/sqlite"
)

// setupOptions holds configuration flags for Setup
type setupOptions struct {
	withVault bool
	cfg       types.ConfigInterface
}

// SetupOption is a functional option for Setup
type SetupOption func(*setupOptions)

// WithVault enables in-memory VaultStore during test setup
func WithVault(enable bool) SetupOption {
	return func(opts *setupOptions) {
		opts.withVault = enable
	}
}

// WithCfg allows providing a custom config for Setup
func WithCfg(cfg types.ConfigInterface) SetupOption {
	return func(opts *setupOptions) {
		opts.cfg = cfg
	}
}

// Setup initializes a default in-memory SQLite application for tests,
// unless overridden via options. It returns the initialized application.
func Setup(options ...SetupOption) types.AppInterface {
	// collect options
	opts := &setupOptions{}
	for _, opt := range options {
		opt(opts)
	}

	// If no config provided, create a reasonable default testing config
	if opts.cfg == nil {
		cfg := &types.Config{}
		cfg.SetAppEnv("testing")
		cfg.SetAppDebug(true)
		cfg.SetDatabaseDriver("sqlite")
		cfg.SetDatabaseHost("")
		cfg.SetDatabasePort("")
		// Use a unique in-memory DB per Setup call to avoid cross-test leakage when running package tests
		// Example DSN: file:mp_test_123456789?mode=memory&cache=shared
		uniqueDSN := fmt.Sprintf("file:mp_test_%d?mode=memory&cache=shared", time.Now().UnixNano())
		cfg.SetDatabaseName(uniqueDSN)
		cfg.SetDatabaseUsername("")
		cfg.SetDatabasePassword("")
		opts.cfg = cfg
	}

	// Apply vault toggle BEFORE building the app so New() initializes stores accordingly
	if opts.withVault {
		opts.cfg.SetVaultStoreUsed(true)
	} else {
		opts.cfg.SetVaultStoreUsed(false)
	}

	// Build application using app.New (opens DB and initializes stores)
	application, err := app.New(opts.cfg)
	if err != nil {
		panic("testutils.Setup: failed to build application: " + err.Error())
	}

	if application.GetLogger() == nil {
		application.SetLogger(slog.Default())
	}

	return application
}

// func setupMailServer() {
// 	mailServer := smtpmock.New(smtpmock.ConfigurationAttr{
// 		LogToStdout:       false, // enable if you have errors sending emails
// 		LogServerActivity: true,
// 		PortNumber:        32435,
// 		HostAddress:       "127.0.0.1",
// 	})

// 	if err := mailServer.Start(); err != nil {
// 		fmt.Println(err)
// 	}
// }
