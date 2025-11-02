package testutils

import (
	"fmt"
	"log/slog"
	"project/internal/app"
	"project/internal/types"
	"time"

	//smtpmock "github.com/mocktools/go-smtp-mock"

	_ "github.com/glebarez/sqlite"
	"github.com/samber/lo"
)

// setupOptions holds configuration flags for Setup
type setupOptions struct {
	WithAuditStore        bool
	WithBlogStore         bool
	WithCmsStore          bool
	WithCacheStore        bool
	WithGeoStore          bool
	WithLogStore          bool
	WithSettingStore      bool
	WithSessionStore      bool
	WithShopStore         bool
	WithSubscriptionStore bool
	WithTaskStore         bool
	WithUserStore         bool
	WithVaultStore        bool
	WithUserStoreVault    bool
	VaultStoreKey         string
	CmsStoreTemplateID    string

	cfg types.ConfigInterface
}

// SetupOption is a functional option for Setup
type SetupOption func(*setupOptions)

// WithCfg allows providing a custom config for Setup
func WithCfg(cfg types.ConfigInterface) SetupOption {
	return func(opts *setupOptions) {
		opts.cfg = cfg
	}
}

// WithBlogStore enables the blog store during test setup
func WithBlogStore(enable bool) SetupOption {
	return func(opts *setupOptions) {
		opts.WithBlogStore = enable
	}
}

// WithCmsStore enables the CMS store during test setup
// Optional templateID sets the CMS template ID used by the layout renderer
func WithCmsStore(enable bool, templateID ...string) SetupOption {
	return func(opts *setupOptions) {
		opts.WithCmsStore = enable
		if enable {
			opts.CmsStoreTemplateID = lo.FirstOr(templateID, opts.CmsStoreTemplateID)
		}
	}
}

// WithCacheStore enables the cache store during test setup
func WithCacheStore(enable bool) SetupOption {
	return func(opts *setupOptions) {
		opts.WithCacheStore = enable
	}
}

// WithSettingStore enables the setting store during test setup
func WithSettingStore(enable bool) SetupOption {
	return func(opts *setupOptions) {
		opts.WithSettingStore = enable
	}
}

// WithGeoStore enables the geo store during test setup
func WithGeoStore(enable bool) SetupOption {
	return func(opts *setupOptions) {
		opts.WithGeoStore = enable
	}
}

// WithLogStore enables the log store during test setup
func WithLogStore(enable bool) SetupOption {
	return func(opts *setupOptions) {
		opts.WithLogStore = enable
	}
}

// WithSessionStore enables the session store during test setup
func WithSessionStore(enable bool) SetupOption {
	return func(opts *setupOptions) {
		opts.WithSessionStore = enable
	}
}

// WithShopStore enables the shop store during test setup
func WithShopStore(enable bool) SetupOption {
	return func(opts *setupOptions) {
		opts.WithShopStore = enable
	}
}

// WithTaskStore enables the task store during test setup
func WithTaskStore(enable bool) SetupOption {
	return func(opts *setupOptions) {
		opts.WithTaskStore = enable
	}
}

// WithUserStore enables the user store during test setup
// Optional bool flag enables user-store vault integration
func WithUserStore(enable bool, enableVault ...bool) SetupOption {
	return func(opts *setupOptions) {
		opts.WithUserStore = enable
		if len(enableVault) > 0 && enableVault[0] {
			opts.WithUserStoreVault = true
		}
	}
}

// WithVaultStore enables the vault store during test setup
// Optional string parameter sets the vault store key
func WithVaultStore(enable bool, vaultKeys ...string) SetupOption {
	return func(opts *setupOptions) {
		opts.WithVaultStore = enable
		opts.VaultStoreKey = lo.FirstOr(vaultKeys, "test-key")
	}
}

// WithSubscriptionStore enables the subscription store during test setup
func WithSubscriptionStore(enable bool) SetupOption {
	return func(opts *setupOptions) {
		opts.WithSubscriptionStore = enable
	}
}

// WithAuditStore enables the audit store during test setup
func WithAuditStore(enable bool) SetupOption {
	return func(opts *setupOptions) {
		opts.WithAuditStore = enable
	}
}

func DefaultConf() *types.Config {
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
	cfg.SetRegistrationEnabled(true)
	cfg.SetMailFromAddress("test@test.com")
	cfg.SetMailFromName("TestName")

	// All stores are disabled by default in tests to ensure explicit configuration
	// Enable only the stores you need in your test using the appropriate With* methods
	// or by directly setting the flags on the config
	return cfg
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
		opts.cfg = DefaultConf()

		// Only set stores if explicitly set when using default config
		if opts.WithAuditStore {
			opts.cfg.SetAuditStoreUsed(true)
		}
		if opts.WithBlogStore {
			opts.cfg.SetBlogStoreUsed(true)
		}
		if opts.WithCmsStore {
			opts.cfg.SetCmsStoreUsed(true)
			if opts.CmsStoreTemplateID != "" {
				opts.cfg.SetCmsStoreTemplateID(opts.CmsStoreTemplateID)
			}
		}
		if opts.WithCacheStore {
			opts.cfg.SetCacheStoreUsed(true)
		}
		if opts.WithGeoStore {
			opts.cfg.SetGeoStoreUsed(true)
		}
		if opts.WithLogStore {
			opts.cfg.SetLogStoreUsed(true)
		}
		if opts.WithSessionStore {
			opts.cfg.SetSessionStoreUsed(true)
		}
		if opts.WithSettingStore {
			opts.cfg.SetSettingStoreUsed(true)
		}
		if opts.WithShopStore {
			opts.cfg.SetShopStoreUsed(true)
		}
		if opts.WithSubscriptionStore {
			opts.cfg.SetSubscriptionStoreUsed(true)
		}
		if opts.WithTaskStore {
			opts.cfg.SetTaskStoreUsed(true)
		}
		if opts.WithUserStore {
			opts.cfg.SetUserStoreUsed(true)
		}
		if opts.WithUserStoreVault {
			opts.cfg.SetUserStoreVaultEnabled(true)
		}
		if opts.WithVaultStore {
			opts.cfg.SetVaultStoreUsed(true)
		}
		if opts.VaultStoreKey != "" {
			opts.cfg.SetVaultStoreKey(opts.VaultStoreKey)
		}
	}

	// Apply optional toggles to provided configs
	if opts.cfg != nil {
		if opts.WithAuditStore {
			opts.cfg.SetAuditStoreUsed(true)
		}
		if opts.WithBlogStore {
			opts.cfg.SetBlogStoreUsed(true)
		}
		if opts.WithCmsStore {
			opts.cfg.SetCmsStoreUsed(true)
			if opts.CmsStoreTemplateID != "" {
				opts.cfg.SetCmsStoreTemplateID(opts.CmsStoreTemplateID)
			}
		}
		if opts.WithCacheStore {
			opts.cfg.SetCacheStoreUsed(true)
		}
		if opts.WithGeoStore {
			opts.cfg.SetGeoStoreUsed(true)
		}
		if opts.WithSessionStore {
			opts.cfg.SetSessionStoreUsed(true)
		}
		if opts.WithShopStore {
			opts.cfg.SetShopStoreUsed(true)
		}
		if opts.WithSubscriptionStore {
			opts.cfg.SetSubscriptionStoreUsed(true)
		}
		if opts.WithTaskStore {
			opts.cfg.SetTaskStoreUsed(true)
		}
		if opts.WithUserStore {
			opts.cfg.SetUserStoreUsed(true)
		}
		if opts.WithVaultStore {
			opts.cfg.SetVaultStoreUsed(true)
		}
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
