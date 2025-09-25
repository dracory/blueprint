package app

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"

	"project/internal/cache"
	"project/internal/types"

	"github.com/dracory/blindindexstore"
	"github.com/dracory/blogstore"
	"github.com/dracory/cachestore"
	"github.com/dracory/cmsstore"
	"github.com/dracory/customstore"
	"github.com/dracory/entitystore"
	"github.com/dracory/feedstore"
	"github.com/dracory/geostore"
	"github.com/dracory/logstore"
	"github.com/dracory/metastore"
	"github.com/dracory/sessionstore"
	"github.com/dracory/settingstore"
	"github.com/dracory/shopstore"
	"github.com/dracory/statsstore"
	"github.com/dracory/taskstore"
	"github.com/dracory/tradingstore"
	"github.com/dracory/userstore"
	"github.com/dracory/vaultstore"
	"github.com/faabiosr/cachego"
	"github.com/faabiosr/cachego/file"
	"github.com/gouniverse/filesystem"
	"github.com/jellydator/ttlcache/v3"
	"github.com/lmittmann/tint"
)

// Application is the orchestration facade for starting the app.
// It encapsulates configuration and database (container removed).
type Application struct {
	cfg types.ConfigInterface
	db  *sql.DB

	// Loggers
	databaseLogger *slog.Logger
	consoleLogger  *slog.Logger

	// Database stores
	blogStore           blogstore.StoreInterface
	blindIndexEmail     blindindexstore.StoreInterface
	blindIndexFirstName blindindexstore.StoreInterface
	blindIndexLastName  blindindexstore.StoreInterface
	cacheStore          cachestore.StoreInterface
	cmsStore            cmsstore.StoreInterface
	customStore         customstore.StoreInterface
	entityStore         entitystore.StoreInterface
	feedStore           feedstore.StoreInterface
	geoStore            geostore.StoreInterface
	logStore            logstore.StoreInterface
	metaStore           metastore.StoreInterface
	sessionStore        sessionstore.StoreInterface
	settingStore        settingstore.StoreInterface
	shopStore           shopstore.StoreInterface
	sqlFileStorage      filesystem.StorageInterface
	statsStore          statsstore.StoreInterface
	taskStore           taskstore.StoreInterface
	tradingStore        tradingstore.StoreInterface
	userStore           userstore.StoreInterface
	vaultStore          vaultstore.StoreInterface
}

var _ types.AppInterface = (*Application)(nil)

// New constructs and initializes the Application (logger, caches, database).
// It centralizes the boot logic so callers only use this single constructor.
func New(cfg types.ConfigInterface) (types.AppInterface, error) {
	if cfg == nil {
		return nil, errors.New("cfg is nil")
	}

	// Caches
	if cache.Memory == nil {
		cache.Memory = ttlcache.New[string, any]()
	}
	// Ensure cache directory exists for file cache
	_ = os.MkdirAll(".cache", os.ModePerm)
	if cache.File == nil {
		cache.File = file.New(".cache")
	}

	consoleLogger := slog.New(tint.NewHandler(os.Stdout, nil))

	// Database open
	db, err := databaseOpen(cfg)
	if err != nil {
		return nil, err
	}

	// Build application instance
	application := &Application{cfg: cfg}
	application.SetConsole(consoleLogger)
	application.SetLogger(consoleLogger)
	application.SetMemoryCache(cache.Memory)
	application.SetFileCache(cache.File)
	application.SetDB(db)

	if err := application.dataStoresInitialize(); err != nil {
		return nil, err
	}

	if err := application.dataStoresMigrate(); err != nil {
		return nil, err
	}

	if application.GetLogStore() != nil {
		application.SetLogger(slog.New(logstore.NewSlogHandler(application.GetLogStore())))
	}

	return application, nil
}

// GetConfig returns the application config
func (a *Application) GetConfig() types.ConfigInterface {
	return a.cfg
}
func (a *Application) SetConfig(cfg types.ConfigInterface) {
	a.cfg = cfg
}

// GetDB returns the application database
func (a *Application) GetDB() *sql.DB {
	return a.db
}

// SetDB sets the application database
func (a *Application) SetDB(db *sql.DB) {
	a.db = db
}

// Run remains for future consolidation of boot logic.
func (a *Application) Run() error { return nil }

// Logger accessors (delegate to package-level logger singletons)
func (a *Application) GetLogger() *slog.Logger {
	return a.databaseLogger
}
func (a *Application) SetLogger(l *slog.Logger) {
	a.databaseLogger = l
}

func (a *Application) GetConsole() *slog.Logger {
	return a.consoleLogger
}
func (a *Application) SetConsole(l *slog.Logger) {
	a.consoleLogger = l
}

// Cache accessors (delegate to package-level cache singletons)
func (a *Application) GetMemoryCache() *ttlcache.Cache[string, any] {
	return cache.Memory
}

func (a *Application) SetMemoryCache(c *ttlcache.Cache[string, any]) {
	cache.Memory = c
}

func (a *Application) GetFileCache() cachego.Cache {
	return cache.File
}

func (a *Application) SetFileCache(c cachego.Cache) {
	cache.File = c
}

// ============================================================================
// == Store accessors
// ============================================================================

// LogStore
func (a *Application) GetLogStore() logstore.StoreInterface {
	return a.logStore
}
func (a *Application) SetLogStore(s logstore.StoreInterface) {
	a.logStore = s
}

// BlogStore
func (a *Application) GetBlogStore() blogstore.StoreInterface {
	return a.blogStore
}
func (a *Application) SetBlogStore(s blogstore.StoreInterface) {
	a.blogStore = s
}

// CacheStore
func (a *Application) GetCacheStore() cachestore.StoreInterface {
	return a.cacheStore
}
func (a *Application) SetCacheStore(s cachestore.StoreInterface) {
	a.cacheStore = s
}

// CmsStore
func (a *Application) GetCmsStore() cmsstore.StoreInterface {
	return a.cmsStore
}
func (a *Application) SetCmsStore(s cmsstore.StoreInterface) {
	a.cmsStore = s
}

// CustomStore
func (a *Application) GetCustomStore() customstore.StoreInterface {
	return a.customStore
}
func (a *Application) SetCustomStore(s customstore.StoreInterface) {
	a.customStore = s
}

// EntityStore
func (a *Application) GetEntityStore() entitystore.StoreInterface {
	return a.entityStore
}
func (a *Application) SetEntityStore(s entitystore.StoreInterface) {
	a.entityStore = s
}

// FeedStore
func (a *Application) GetFeedStore() feedstore.StoreInterface {
	return a.feedStore
}
func (a *Application) SetFeedStore(s feedstore.StoreInterface) {
	a.feedStore = s
}

// GeoStore
func (a *Application) GetGeoStore() geostore.StoreInterface {
	return a.geoStore
}
func (a *Application) SetGeoStore(s geostore.StoreInterface) {
	a.geoStore = s
}

// MetaStore
func (a *Application) GetMetaStore() metastore.StoreInterface {
	return a.metaStore
}

func (a *Application) SetMetaStore(s metastore.StoreInterface) {
	a.metaStore = s
}

// SessionStore
func (a *Application) GetSessionStore() sessionstore.StoreInterface {
	return a.sessionStore
}
func (a *Application) SetSessionStore(s sessionstore.StoreInterface) {
	a.sessionStore = s
}

// ShopStore
func (a *Application) GetShopStore() shopstore.StoreInterface {
	return a.shopStore
}
func (a *Application) SetShopStore(s shopstore.StoreInterface) {
	a.shopStore = s
}

// SqlFileStorage
func (a *Application) GetSqlFileStorage() filesystem.StorageInterface {
	return a.sqlFileStorage
}
func (a *Application) SetSqlFileStorage(s filesystem.StorageInterface) {
	a.sqlFileStorage = s
}

// StatsStore
func (a *Application) GetStatsStore() statsstore.StoreInterface {
	return a.statsStore
}
func (a *Application) SetStatsStore(s statsstore.StoreInterface) {
	a.statsStore = s
}

// TaskStore
func (a *Application) GetTaskStore() taskstore.StoreInterface {
	return a.taskStore
}
func (a *Application) SetTaskStore(s taskstore.StoreInterface) {
	a.taskStore = s
}

// TradingStore
func (a *Application) GetTradingStore() tradingstore.StoreInterface {
	return a.tradingStore
}
func (a *Application) SetTradingStore(s tradingstore.StoreInterface) {
	a.tradingStore = s
}

// UserStore
func (a *Application) GetUserStore() userstore.StoreInterface {
	return a.userStore
}
func (a *Application) SetUserStore(s userstore.StoreInterface) {
	a.userStore = s
}

// VaultStore
func (a *Application) GetVaultStore() vaultstore.StoreInterface {
	return a.vaultStore
}
func (a *Application) SetVaultStore(s vaultstore.StoreInterface) {
	a.vaultStore = s
}

// SettingStore
func (a *Application) GetSettingStore() settingstore.StoreInterface {
	return a.settingStore
}
func (a *Application) SetSettingStore(s settingstore.StoreInterface) {
	a.settingStore = s
}

// Blind index stores
func (a *Application) GetBlindIndexStoreEmail() blindindexstore.StoreInterface {
	return a.blindIndexEmail
}
func (a *Application) SetBlindIndexStoreEmail(s blindindexstore.StoreInterface) {
	a.blindIndexEmail = s
}
func (a *Application) GetBlindIndexStoreFirstName() blindindexstore.StoreInterface {
	return a.blindIndexFirstName
}
func (a *Application) SetBlindIndexStoreFirstName(s blindindexstore.StoreInterface) {
	a.blindIndexFirstName = s
}
func (a *Application) GetBlindIndexStoreLastName() blindindexstore.StoreInterface {
	return a.blindIndexLastName
}
func (a *Application) SetBlindIndexStoreLastName(s blindindexstore.StoreInterface) {
	a.blindIndexLastName = s
}
