package app

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"project/internal/cache"
	"project/internal/types"

	"github.com/dracory/auditstore"
	"github.com/dracory/blindindexstore"
	"github.com/dracory/blogstore"
	"github.com/dracory/cachestore"
	"github.com/dracory/chatstore"
	"github.com/dracory/cmsstore"
	"github.com/dracory/customstore"
	"github.com/dracory/entitystore"
	"github.com/dracory/feedstore"
	"github.com/dracory/filesystem"
	"github.com/dracory/geostore"
	"github.com/dracory/logstore"
	"github.com/dracory/metastore"
	"github.com/dracory/sessionstore"
	"github.com/dracory/settingstore"
	"github.com/dracory/shopstore"
	"github.com/dracory/statsstore"
	"github.com/dracory/subscriptionstore"
	"github.com/dracory/taskstore"
	"github.com/dracory/userstore"
	"github.com/dracory/vaultstore"
	"github.com/faabiosr/cachego"
	"github.com/faabiosr/cachego/file"
	"github.com/jellydator/ttlcache/v3"
	"github.com/lmittmann/tint"
	// "gorm.io/gorm"
)

// Registry is the orchestration facade for starting the app.
// It encapsulates configuration and database (container removed).
type Registry struct {
	cfg types.ConfigInterface
	db  *sql.DB

	// Loggers
	databaseLogger *slog.Logger
	consoleLogger  *slog.Logger

	// Database stores
	auditStore          auditstore.StoreInterface
	blogStore           blogstore.StoreInterface
	blindIndexEmail     blindindexstore.StoreInterface
	blindIndexFirstName blindindexstore.StoreInterface
	blindIndexLastName  blindindexstore.StoreInterface
	cacheStore          cachestore.StoreInterface
	chatStore           chatstore.StoreInterface
	cmsStore            cmsstore.StoreInterface
	customStore         customstore.StoreInterface
	entityStore         entitystore.StoreInterface
	feedStore           feedstore.StoreInterface
	geoStore            geostore.StoreInterface
	logStore            logstore.StoreInterface
	metaStore           metastore.StoreInterface
	subscriptionStore   subscriptionstore.StoreInterface
	sessionStore        sessionstore.StoreInterface
	settingStore        settingstore.StoreInterface
	shopStore           shopstore.StoreInterface
	sqlFileStorage      filesystem.StorageInterface
	statsStore          statsstore.StoreInterface
	taskStore           taskstore.StoreInterface
	userStore           userstore.StoreInterface
	vaultStore          vaultstore.StoreInterface
}

// Ensure Registry satisfies the RegistryInterface contract.
var _ types.RegistryInterface = (*Registry)(nil)

// New constructs and initializes the Registry (logger, caches, database).
// It centralizes the boot logic so callers only use this single constructor.
func New(cfg types.ConfigInterface) (types.RegistryInterface, error) {
	if cfg == nil {
		return nil, errors.New("cfg is nil")
	}

	// Caches
	if cache.Memory == nil {
		cache.Memory = ttlcache.New[string, any]()
	}
	// Ensure cache directory exists for file cache
	cacheDir := cacheDirectory()
	_ = os.MkdirAll(cacheDir, os.ModePerm)
	if cache.File == nil {
		cache.File = file.New(cacheDir)
	}

	consoleLogger := slog.New(tint.NewHandler(os.Stdout, nil))

	// Database open
	db, err := databaseOpen(cfg)
	if err != nil {
		return nil, err
	}

	// Build registry instance
	registry := &Registry{cfg: cfg}
	registry.SetConsole(consoleLogger)
	registry.SetLogger(consoleLogger)
	registry.SetMemoryCache(cache.Memory)
	registry.SetFileCache(cache.File)
	registry.SetDatabase(db)

	if err := registry.dataStoresInitialize(); err != nil {
		return nil, err
	}

	if err := registry.dataStoresMigrate(); err != nil {
		return nil, err
	}

	if registry.GetLogStore() != nil {
		registry.SetLogger(slog.New(logstore.NewSlogHandler(registry.GetLogStore())))
	}

	return registry, nil
}

// GetConfig returns the registry config
func (r *Registry) GetConfig() types.ConfigInterface {
	if r == nil {
		return nil
	}
	return r.cfg
}
func (r *Registry) SetConfig(cfg types.ConfigInterface) {
	r.cfg = cfg
}

// GetDatabase returns the registry database
func (r *Registry) GetDatabase() *sql.DB {
	return r.db
}

// SetDatabase sets the registry database
func (r *Registry) SetDatabase(db *sql.DB) {
	r.db = db
}

// Enable if you want to use GORM
// // GetDatabaseGorm returns the GORM database handle.
// func (r *Registry) GetDatabaseGorm() *gorm.DB {
// 	if r == nil {
// 		return nil
// 	}
// 	return r.gormDb
// }

// // SetDatabaseGorm sets the GORM database handle for the application.
// func (r *Registry) SetDatabaseGorm(gormDb *gorm.DB) {
// 	r.gormDb = gormDb
// }

// Logger accessors (delegate to package-level logger singletons)
func (r *Registry) GetLogger() *slog.Logger {
	return r.databaseLogger
}
func (r *Registry) SetLogger(l *slog.Logger) {
	r.databaseLogger = l
}

func (r *Registry) GetConsole() *slog.Logger {
	return r.consoleLogger
}
func (r *Registry) SetConsole(l *slog.Logger) {
	r.consoleLogger = l
}

// Cache accessors (delegate to package-level cache singletons)
func (r *Registry) GetMemoryCache() *ttlcache.Cache[string, any] {
	return cache.Memory
}

func (r *Registry) SetMemoryCache(c *ttlcache.Cache[string, any]) {
	cache.Memory = c
}

func (r *Registry) GetFileCache() cachego.Cache {
	return cache.File
}

func (r *Registry) SetFileCache(c cachego.Cache) {
	cache.File = c
}

// cacheDirectory returns the path to the shared filesystem cache directory.
// It walks up from the current working directory until it finds go.mod and
// then returns the .cache directory at that project root. If no go.mod is
// found, it falls back to a relative .cache in the current working directory.
func cacheDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		return ".cache"
	}

	dir := wd
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return filepath.Join(dir, ".cache")
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ".cache"
}

// ============================================================================
// == Store accessors
// ============================================================================

func (r *Registry) GetAuditStore() auditstore.StoreInterface {
	return r.auditStore
}
func (r *Registry) SetAuditStore(s auditstore.StoreInterface) {
	r.auditStore = s
}

// BlogStore
func (r *Registry) GetBlogStore() blogstore.StoreInterface {
	return r.blogStore
}
func (r *Registry) SetBlogStore(s blogstore.StoreInterface) {
	r.blogStore = s
}

// ChatStore
func (r *Registry) GetChatStore() chatstore.StoreInterface {
	return r.chatStore
}
func (r *Registry) SetChatStore(s chatstore.StoreInterface) {
	r.chatStore = s
}

// CacheStore
func (r *Registry) GetCacheStore() cachestore.StoreInterface {
	return r.cacheStore
}
func (r *Registry) SetCacheStore(s cachestore.StoreInterface) {
	r.cacheStore = s
}

// CmsStore
func (r *Registry) GetCmsStore() cmsstore.StoreInterface {
	return r.cmsStore
}
func (r *Registry) SetCmsStore(s cmsstore.StoreInterface) {
	r.cmsStore = s
}

// CustomStore
func (r *Registry) GetCustomStore() customstore.StoreInterface {
	return r.customStore
}
func (r *Registry) SetCustomStore(s customstore.StoreInterface) {
	r.customStore = s
}

// EntityStore
func (r *Registry) GetEntityStore() entitystore.StoreInterface {
	return r.entityStore
}
func (r *Registry) SetEntityStore(s entitystore.StoreInterface) {
	r.entityStore = s
}

// FeedStore
func (r *Registry) GetFeedStore() feedstore.StoreInterface {
	return r.feedStore
}
func (r *Registry) SetFeedStore(s feedstore.StoreInterface) {
	r.feedStore = s
}

// GeoStore
func (r *Registry) GetGeoStore() geostore.StoreInterface {
	return r.geoStore
}
func (r *Registry) SetGeoStore(s geostore.StoreInterface) {
	r.geoStore = s
}

// LogStore
func (r *Registry) GetLogStore() logstore.StoreInterface {
	return r.logStore
}
func (r *Registry) SetLogStore(s logstore.StoreInterface) {
	r.logStore = s
}

// MetaStore
func (r *Registry) GetMetaStore() metastore.StoreInterface {
	return r.metaStore
}

func (r *Registry) SetMetaStore(s metastore.StoreInterface) {
	r.metaStore = s
}

// SessionStore

// GetSessionStore returns the session store.
func (r *Registry) GetSessionStore() sessionstore.StoreInterface {
	if r == nil {
		return nil
	}
	return r.sessionStore
}
func (r *Registry) SetSessionStore(s sessionstore.StoreInterface) {
	r.sessionStore = s
}

// ShopStore
func (r *Registry) GetShopStore() shopstore.StoreInterface {
	return r.shopStore
}
func (r *Registry) SetShopStore(s shopstore.StoreInterface) {
	r.shopStore = s
}

// SqlFileStorage
func (r *Registry) GetSqlFileStorage() filesystem.StorageInterface {
	return r.sqlFileStorage
}
func (r *Registry) SetSqlFileStorage(s filesystem.StorageInterface) {
	r.sqlFileStorage = s
}

// StatsStore
func (r *Registry) GetStatsStore() statsstore.StoreInterface {
	return r.statsStore
}
func (r *Registry) SetStatsStore(s statsstore.StoreInterface) {
	r.statsStore = s
}

// TaskStore
func (r *Registry) GetTaskStore() taskstore.StoreInterface {
	return r.taskStore
}
func (r *Registry) SetTaskStore(s taskstore.StoreInterface) {
	r.taskStore = s
}

// UserStore
func (r *Registry) GetUserStore() userstore.StoreInterface {
	return r.userStore
}
func (r *Registry) SetUserStore(s userstore.StoreInterface) {
	r.userStore = s
}

// VaultStore
func (r *Registry) GetVaultStore() vaultstore.StoreInterface {
	return r.vaultStore
}
func (r *Registry) SetVaultStore(s vaultstore.StoreInterface) {
	r.vaultStore = s
}

// SettingStore
func (r *Registry) GetSettingStore() settingstore.StoreInterface {
	return r.settingStore
}
func (r *Registry) SetSettingStore(s settingstore.StoreInterface) {
	r.settingStore = s
}

func (r *Registry) GetSubscriptionStore() subscriptionstore.StoreInterface {
	return r.subscriptionStore
}
func (r *Registry) SetSubscriptionStore(s subscriptionstore.StoreInterface) {
	r.subscriptionStore = s
}

// Blind index stores
func (r *Registry) GetBlindIndexStoreEmail() blindindexstore.StoreInterface {
	return r.blindIndexEmail
}
func (r *Registry) SetBlindIndexStoreEmail(s blindindexstore.StoreInterface) {
	r.blindIndexEmail = s
}
func (r *Registry) GetBlindIndexStoreFirstName() blindindexstore.StoreInterface {
	return r.blindIndexFirstName
}
func (r *Registry) SetBlindIndexStoreFirstName(s blindindexstore.StoreInterface) {
	r.blindIndexFirstName = s
}
func (r *Registry) GetBlindIndexStoreLastName() blindindexstore.StoreInterface {
	return r.blindIndexLastName
}
func (r *Registry) SetBlindIndexStoreLastName(s blindindexstore.StoreInterface) {
	r.blindIndexLastName = s
}
