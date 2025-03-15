package config

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/dracory/base/database"
	"github.com/dracory/base/env"
	"github.com/faabiosr/cachego"
	"github.com/faabiosr/cachego/file"
	"github.com/gouniverse/blindindexstore"
	"github.com/gouniverse/blockeditor"
	"github.com/gouniverse/blogstore"
	"github.com/gouniverse/cachestore"
	"github.com/gouniverse/cms"
	"github.com/gouniverse/cmsstore"
	"github.com/gouniverse/customstore"
	"github.com/gouniverse/entitystore"
	"github.com/gouniverse/geostore"
	"github.com/gouniverse/logstore"
	"github.com/gouniverse/metastore"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/sessionstore"
	"github.com/gouniverse/shopstore"
	"github.com/gouniverse/statsstore"
	"github.com/gouniverse/taskstore"
	"github.com/gouniverse/ui"
	"github.com/gouniverse/userstore"
	"github.com/gouniverse/vaultstore"
	"github.com/jellydator/ttlcache/v3"

	"project/app/config/stores"
)

// Config represents the application configuration
type Config struct {
	// Server configuration
	ServerHost string
	ServerPort string
	AppURL     string
	AppName    string
	AppEnv     string
	AppDebug   bool
	AppSecret  string
	AppVersion string

	// Database configuration
	DatabaseDriver   string
	DatabaseHost     string
	DatabasePort     string
	DatabaseName     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseSSLMode  string

	// Mail configuration
	MailDriver           string
	MailFromEmailAddress string
	MailFromName         string
	MailHost             string
	MailPort             string
	MailPassword         string
	MailUsername         string

	// Media configuration
	MediaBucket   string
	MediaDriver   string
	MediaKey      string
	MediaEndpoint string
	MediaRegion   string
	MediaRoot     string
	MediaSecret   string
	MediaURL      string

	// Auth configuration
	AuthEndpoint string

	// API Keys
	OpenAIAPIKey     string
	StripeKeyPrivate string
	StripeKeyPublic  string
	VaultKey         string

	// Vertex configuration
	VertexModelID   string
	VertexProjectID string
	VertexRegionID  string

	// Translation configuration
	TranslationLanguageDefault string
	TranslationLanguageList    map[string]string

	// Database and server instances
	Database  interface{}
	WebServer interface{}

	// Cache instances
	cacheMemory *ttlcache.Cache[string, any]
	cacheFile   cachego.Cache

	// Store usage flags
	BlindIndexStoreUsed bool
	BlogStoreUsed       bool
	CmsStoreUsed        bool
	CacheStoreUsed      bool
	CustomStoreUsed     bool
	EntityStoreUsed     bool
	GeoStoreUsed        bool
	LogStoreUsed        bool
	MetaStoreUsed       bool
	SessionStoreUsed    bool
	ShopStoreUsed       bool
	SqlFileStoreUsed    bool
	StatsStoreUsed      bool
	TaskStoreUsed       bool
	UserStoreUsed       bool
	VaultStoreUsed      bool
	CmsUsed             bool

	// Store instances - using interface{} to avoid import issues
	BlindIndexStoreEmail     blindindexstore.StoreInterface
	BlindIndexStoreFirstName blindindexstore.StoreInterface
	BlindIndexStoreLastName  blindindexstore.StoreInterface
	BlogStore                blogstore.StoreInterface
	CmsStore                 cmsstore.StoreInterface
	CacheStore               cachestore.StoreInterface
	CustomStore              customstore.StoreInterface
	EntityStore              entitystore.StoreInterface
	GeoStore                 geostore.StoreInterface
	LogStore                 logstore.StoreInterface
	MetaStore                metastore.StoreInterface
	SessionStore             sessionstore.StoreInterface
	ShopStore                shopstore.StoreInterface
	StatsStore               statsstore.StoreInterface
	TaskStore                taskstore.StoreInterface
	UserStore                userstore.StoreInterface
	VaultStore               vaultstore.StoreInterface
	Cms                      *cms.Cms

	// CMS template ID
	CmsUserTemplateID string

	// Logger
	Logger *slog.Logger

	// Database initialization and migration functions
	databaseInits      []func(db *sql.DB) error
	databaseMigrations []func(ctx context.Context) error
}

// AuthenticatedUserContextKey is a context key for the authenticated user.
type AuthenticatedUserContextKey struct{}

// AuthenticatedSessionContextKey is a context key for the authenticated session.
type AuthenticatedSessionContextKey struct{}

// New creates a new Config instance with values loaded from environment variables
func New() (*Config, error) {
	env.Initialize()

	serverHost, err := env.ValueOrError("SERVER_HOST")

	if err != nil {
		return nil, errors.New("SERVER_HOST is required")
	}

	serverPort, err := env.ValueOrError("SERVER_PORT")

	if err != nil {
		return nil, errors.New("SERVER_PORT is required")
	}

	appEnv, err := env.ValueOrError("APP_ENV")

	if err != nil {
		return nil, errors.New("APP_ENV is required")
	}

	databaseDriver, err := env.ValueOrError("DB_DRIVER")

	if err != nil {
		return nil, errors.New("DB_DRIVER is required")
	}

	databaseName, err := env.ValueOrError("DB_DATABASE")

	if err != nil {
		return nil, errors.New("DB_DATABASE is required")
	}

	databaseHost, err := env.ValueOrError("DB_HOST")

	if databaseDriver != "sqlite" && err != nil {
		return nil, errors.New("DB_HOST is required")
	}

	databasePort, err := env.ValueOrError("DB_PORT")

	if databaseDriver != "sqlite" && err != nil {
		return nil, errors.New("DB_PORT is required")
	}

	databaseUser, err := env.ValueOrError("DB_USERNAME")

	if databaseDriver != "sqlite" && err != nil {
		return nil, errors.New("DB_USERNAME is required")
	}

	databasePassword, err := env.ValueOrError("DB_PASSWORD")

	if databaseDriver != "sqlite" && err != nil {
		return nil, errors.New("DB_PASSWORD is required")
	}

	// Optional API keys - uncomment as needed
	// config.OpenAIAPIKey = env.Must("OPENAI_API_KEY")
	// config.StripeKeyPrivate = env.Must("STRIPE_KEY_PRIVATE")
	// config.StripeKeyPublic = env.Must("STRIPE_KEY_PUBLIC")
	// config.VertexModelID = env.Must("VERTEX_MODEL_ID")
	// config.VertexProjectID = env.Must("VERTEX_PROJECT_ID")
	// config.VertexRegionID = env.Must("VERTEX_REGION_ID")

	// CMS template ID - uncomment if needed
	// config.CmsUserTemplateID = env.Must("CMS_TEMPLATE_ID")

	// Create config with default values and environment constants
	config := &Config{
		ServerHost: serverHost,
		ServerPort: serverPort,
		AppURL:     env.ValueOrDefault("APP_URL", "http://localhost:8080"),
		AppName:    env.Value("APP_NAME"),
		AppEnv:     appEnv,
		AppDebug:   env.Value("DEBUG") == "yes",
		AppSecret:  env.ValueOrDefault("APP_SECRET", "change-me-in-production"),
		AppVersion: env.ValueOrDefault("APP_VERSION", "1.0.0"),

		// Database configuration
		DatabaseDriver:   databaseDriver,
		DatabaseName:     databaseName,
		DatabaseHost:     databaseHost,
		DatabasePort:     databasePort,
		DatabaseUser:     databaseUser,
		DatabasePassword: databasePassword,
		DatabaseSSLMode:  env.ValueOrDefault("DB_SSL_MODE", "disable"),

		// Mail configuration
		MailDriver:           env.Value("MAIL_DRIVER"),
		MailFromEmailAddress: env.Value("EMAIL_FROM_ADDRESS"),
		MailFromName:         env.Value("EMAIL_FROM_NAME"),
		MailHost:             env.Value("MAIL_HOST"),
		MailPassword:         env.Value("MAIL_PASSWORD"),
		MailPort:             env.Value("MAIL_PORT"),
		MailUsername:         env.Value("MAIL_USERNAME"),

		// Media configuration
		MediaBucket:   env.Value("MEDIA_BUCKET"),
		MediaDriver:   env.Value("MEDIA_DRIVER"),
		MediaEndpoint: env.Value("MEDIA_ENDPOINT"),
		MediaKey:      env.Value("MEDIA_KEY"),
		MediaRoot:     env.ValueOrDefault("MEDIA_ROOT", "/"),
		MediaSecret:   env.Value("MEDIA_SECRET"),
		MediaRegion:   env.Value("MEDIA_REGION"),
		MediaURL:      env.ValueOrDefault("MEDIA_URL", "/files"),

		AuthEndpoint: "/auth",

		TranslationLanguageDefault: "en",
		TranslationLanguageList:    map[string]string{"en": "English", "bg": "Bulgarian", "de": "German"},

		// Initialize store usage flags with defaults
		BlindIndexStoreUsed: true,
		BlogStoreUsed:       true,
		CmsStoreUsed:        false,
		CacheStoreUsed:      true,
		CustomStoreUsed:     false,
		EntityStoreUsed:     true,
		GeoStoreUsed:        true,
		LogStoreUsed:        true,
		MetaStoreUsed:       true,
		SessionStoreUsed:    true,
		ShopStoreUsed:       true,
		SqlFileStoreUsed:    true,
		StatsStoreUsed:      true,
		TaskStoreUsed:       true,
		UserStoreUsed:       true,
		VaultStoreUsed:      true,
		CmsUsed:             true,

		// Initialize empty slices for database functions
		databaseInits:      []func(db *sql.DB) error{},
		databaseMigrations: []func(ctx context.Context) error{},
	}

	if vaultKey, err := env.ValueOrError("VAULT_KEY"); config.VaultStoreUsed && err != nil {
		return nil, errors.New("VAULT_KEY is required")
	} else {
		config.VaultKey = vaultKey
	}

	return config, nil
}

// NewFromFile creates a new Config instance with values loaded from a specific .env file
// func NewFromFile(envFile string) (*Config, error) {
// 	// Get absolute path to the env file
// 	absPath, err := filepath.Abs(envFile)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Load the specified .env file
// 	err = godotenv.Load(absPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return New()
// }

// Initialize initializes the application configuration
func (c *Config) Initialize() error {
	// Set timezone to UTC
	os.Setenv("TZ", "UTC")

	// Initialize database
	err := c.initializeDatabase()
	if err != nil {
		return err
	}

	// Initialize stores
	db, err := c.GetRawDatabase()
	if err != nil {
		return err
	}

	// Initialize stores based on usage flags
	if c.CacheStoreUsed {
		// Initialize cache store
		cacheStore, err := stores.CacheStoreInitialize(db)
		if err != nil {
			return err
		}
		c.CacheStore = cacheStore
		c.AddDatabaseMigration(func(ctx context.Context) error {
			return stores.CacheStoreAutoMigrate(ctx, cacheStore)
		})
	}

	if c.CmsStoreUsed {
		// Initialize CMS store
		cmsStore, err := stores.CmsStoreInitialize(db)
		if err != nil {
			return err
		}
		c.CmsStore = cmsStore
		c.AddDatabaseMigration(func(ctx context.Context) error {
			return stores.CmsStoreAutoMigrate(ctx, cmsStore)
		})
	}

	if c.EntityStoreUsed {
		// Initialize entity store
		entityStore, err := stores.EntityStoreInitialize(db)
		if err != nil {
			return err
		}
		c.EntityStore = entityStore
		c.AddDatabaseMigration(func(ctx context.Context) error {
			return stores.EntityStoreAutoMigrate(ctx, entityStore)
		})
	}

	if c.LogStoreUsed {
		// Initialize log store
		logStore, err := stores.LogStoreInitialize(db)
		if err != nil {
			return err
		}
		c.LogStore = logStore
		c.AddDatabaseMigration(func(ctx context.Context) error {
			return stores.LogStoreAutoMigrate(ctx, logStore)
		})
	}

	if c.MetaStoreUsed {
		// Initialize meta store
		metaStore, err := stores.MetaStoreInitialize(db)
		if err != nil {
			return err
		}
		c.MetaStore = metaStore
		c.AddDatabaseMigration(func(ctx context.Context) error {
			return stores.MetaStoreAutoMigrate(ctx, metaStore)
		})
	}

	if c.SessionStoreUsed {
		// Initialize session store
		sessionStore, err := stores.SessionStoreInitialize(db)
		if err != nil {
			return err
		}
		c.SessionStore = sessionStore
		c.AddDatabaseMigration(func(ctx context.Context) error {
			return stores.SessionStoreAutoMigrate(ctx, sessionStore)
		})
	}

	if c.ShopStoreUsed {
		// Initialize shop store
		shopStore, err := stores.ShopStoreInitialize(db)
		if err != nil {
			return err
		}
		c.ShopStore = shopStore
		c.AddDatabaseMigration(func(ctx context.Context) error {
			return stores.ShopStoreAutoMigrate(ctx, shopStore)
		})
	}

	if c.StatsStoreUsed {
		// Initialize stats store
		statsStore, err := stores.StatsStoreInitialize(db)
		if err != nil {
			return err
		}
		c.StatsStore = statsStore
		c.AddDatabaseMigration(func(ctx context.Context) error {
			return stores.StatsStoreAutoMigrate(ctx, statsStore)
		})
	}

	if c.TaskStoreUsed {
		// Initialize task store
		taskStore, err := stores.TaskStoreInitialize(db)
		if err != nil {
			return err
		}
		c.TaskStore = taskStore
		c.AddDatabaseMigration(func(ctx context.Context) error {
			return stores.TaskStoreAutoMigrate(ctx, taskStore)
		})
	}

	if c.VaultStoreUsed {
		// Initialize vault store
		vaultStore, err := stores.VaultStoreInitialize(db)
		if err != nil {
			return err
		}
		c.VaultStore = vaultStore
		c.AddDatabaseMigration(func(ctx context.Context) error {
			return stores.VaultStoreAutoMigrate(ctx, vaultStore)
		})
	}

	if c.UserStoreUsed {
		// Initialize user store
		userStore, err := stores.UserStoreInitialize(db)
		if err != nil {
			return err
		}
		c.UserStore = userStore
		c.AddDatabaseMigration(func(ctx context.Context) error {
			return stores.UserStoreAutoMigrate(ctx, userStore)
		})
	}

	if c.CustomStoreUsed {
		// Initialize custom store
		customStore, err := stores.CustomStoreInitialize(db)
		if err != nil {
			return err
		}
		c.CustomStore = customStore
		c.AddDatabaseMigration(func(ctx context.Context) error {
			return stores.CustomStoreAutoMigrate(ctx, customStore)
		})
	}

	if c.GeoStoreUsed {
		// Initialize geo store
		geoStore, err := stores.GeoStoreInitialize(db)
		if err != nil {
			return err
		}
		c.GeoStore = geoStore
		c.AddDatabaseMigration(func(ctx context.Context) error {
			return stores.GeoStoreAutoMigrate(ctx, geoStore)
		})
	}

	if c.CmsUsed {
		// Initialize CMS
		blockEditorDefinitions := []blockeditor.BlockDefinition{}
		blockEditorRenderer := func(blocks []ui.BlockInterface) string {
			return ""
		}

		cms, err := stores.CmsInitialize(db, blockEditorDefinitions, blockEditorRenderer, c.TranslationLanguageDefault, c.TranslationLanguageList)
		if err != nil {
			return err
		}
		c.Cms = cms
	}

	if c.BlogStoreUsed {
		// Initialize blog store
		blogStore, err := stores.BlogStoreInitialize(db)
		if err != nil {
			return err
		}
		c.BlogStore = blogStore
		c.AddDatabaseMigration(func(ctx context.Context) error {
			return stores.BlogStoreAutoMigrate(ctx, blogStore)
		})
	}

	if c.BlindIndexStoreUsed {
		// Initialize email blind index store
		emailStore, err := stores.BlindIndexStoreInitialize(db, "snv_bindx_email")
		if err != nil {
			return err
		}
		c.BlindIndexStoreEmail = emailStore

		// Initialize first name blind index store
		firstNameStore, err := stores.BlindIndexStoreInitialize(db, "snv_bindx_first_name")
		if err != nil {
			return err
		}
		c.BlindIndexStoreFirstName = firstNameStore

		// Initialize last name blind index store
		lastNameStore, err := stores.BlindIndexStoreInitialize(db, "snv_bindx_last_name")
		if err != nil {
			return err
		}
		c.BlindIndexStoreLastName = lastNameStore

		// Add migrations for all blind index stores
		c.AddDatabaseMigration(func(ctx context.Context) error {
			if err := stores.MigrateBlindIndexStore(ctx, c.BlindIndexStoreEmail); err != nil {
				return err
			}
			if err := stores.MigrateBlindIndexStore(ctx, c.BlindIndexStoreFirstName); err != nil {
				return err
			}
			if err := stores.MigrateBlindIndexStore(ctx, c.BlindIndexStoreLastName); err != nil {
				return err
			}
			return nil
		})
	}

	// Migrate database
	err = c.migrateDatabase()
	if err != nil {
		return err
	}

	// Initialize cache
	c.initializeCache()

	// Setup logger
	c.setupLogger()

	return nil
}

// initializeCache initializes the cache
func (c *Config) initializeCache() {
	// Initialize memory cache
	c.setupCache()
}

// initializeDatabase initializes the database
func (c *Config) initializeDatabase() error {
	db, err := database.Open(database.Options().
		SetDatabaseType(c.DatabaseDriver).
		SetDatabaseHost(c.DatabaseHost).
		SetDatabasePort(c.DatabasePort).
		SetDatabaseName(c.DatabaseName).
		SetCharset(`utf8mb4`).
		SetUserName(c.DatabaseUser).
		SetPassword(c.DatabasePassword))

	if err != nil {
		return err
	}

	if db == nil {
		return errors.New("db is nil")
	}

	// Setup database instance
	err = c.setupDatabase(db)
	if err != nil {
		return err
	}

	// Run database initialization functions
	for _, init := range c.databaseInits {
		err = init(db)
		if err != nil {
			return err
		}
	}

	return nil
}

// migrateDatabase migrates the database
func (c *Config) migrateDatabase() error {
	for _, migrate := range c.databaseMigrations {
		err := migrate(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}

// AddDatabaseInit adds a database initialization function
func (c *Config) AddDatabaseInit(init func(db *sql.DB) error) {
	c.databaseInits = append(c.databaseInits, init)
}

// AddDatabaseMigration adds a database migration function
func (c *Config) AddDatabaseMigration(migration func(ctx context.Context) error) {
	c.databaseMigrations = append(c.databaseMigrations, migration)
}

// IsDevelopment returns true if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.AppEnv == APP_ENVIRONMENT_DEVELOPMENT
}

// IsLocal returns true if the application is running in local mode
func (c *Config) IsLocal() bool {
	return c.AppEnv == APP_ENVIRONMENT_LOCAL
}

// IsProduction returns true if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.AppEnv == APP_ENVIRONMENT_PRODUCTION
}

// IsStaging returns true if the application is running in staging mode
func (c *Config) IsStaging() bool {
	return c.AppEnv == APP_ENVIRONMENT_STAGING
}

// IsTesting returns true if the application is running in testing mode
func (c *Config) IsTesting() bool {
	return c.AppEnv == APP_ENVIRONMENT_TESTING
}

// IsDebugEnabled returns true if debug mode is enabled
func (c *Config) IsDebugEnabled() bool {
	return c.AppDebug
}

// GetDatabase returns the database instance as sb.DatabaseInterface
func (c *Config) GetDatabase() (sb.DatabaseInterface, error) {
	if c.Database == nil {
		return nil, errors.New("database is not initialized")
	}

	db, ok := c.Database.(sb.DatabaseInterface)
	if !ok {
		return nil, errors.New("database is not of type sb.DatabaseInterface")
	}

	return db, nil
}

// setupCache sets up the memory cache
func (c *Config) setupCache() {
	// Initialize memory cache
	c.cacheMemory = ttlcache.New[string, any](
		ttlcache.WithTTL[string, any](1 * time.Hour),
	)

	// Initialize file cache
	_ = os.MkdirAll(".cache", os.ModePerm)
	c.cacheFile = file.New(".cache")
}

// setupLogger sets up the logger
func (c *Config) setupLogger() {
	// Setup logger if log store is used
	if c.LogStoreUsed {
		// Create a basic logger that outputs to stdout
		c.Logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}
}

// setupDatabase sets up the database
func (c *Config) setupDatabase(db *sql.DB) error {
	if db == nil {
		return errors.New("db is nil")
	}

	// Create a database instance using sb package
	dbInstance := sb.NewDatabase(db, c.DatabaseDriver)
	if dbInstance == nil {
		return errors.New("failed to create database instance")
	}

	c.Database = dbInstance
	return nil
}

// GetRawDatabase returns the raw sql.DB instance
func (c *Config) GetRawDatabase() (*sql.DB, error) {
	dbInterface, err := c.GetDatabase()
	if err != nil {
		return nil, err
	}

	return dbInterface.DB(), nil
}

// CacheMemory returns the memory cache instance
func (c *Config) GetCacheMemory() *ttlcache.Cache[string, any] {
	return c.cacheMemory
}

// GetCacheFile returns the file cache instance
func (c *Config) GetCacheFile() cachego.Cache {
	return c.cacheFile
}
