package registry

import (
	"database/sql"
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"project/internal/cache"
	"project/internal/config"

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
type registryImplementation struct {
	cfg config.ConfigInterface
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

// Ensure registryImplementation satisfies the RegistryInterface contract.
var _ RegistryInterface = (*registryImplementation)(nil)

// New constructs and initializes the Registry (logger, caches, database).
// It centralizes the boot logic so callers only use this single constructor.
func New(cfg config.ConfigInterface) (RegistryInterface, error) {
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
	registry := &registryImplementation{cfg: cfg}
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
func (r *registryImplementation) GetConfig() config.ConfigInterface {
	if r == nil {
		return nil
	}
	return r.cfg
}
func (r *registryImplementation) SetConfig(cfg config.ConfigInterface) {
	r.cfg = cfg
}

// GetDatabase returns the registry database
func (r *registryImplementation) GetDatabase() *sql.DB {
	return r.db
}

// SetDatabase sets the registry database
func (r *registryImplementation) SetDatabase(db *sql.DB) {
	r.db = db
}

// Logger accessors (delegate to package-level logger singletons)
func (r *registryImplementation) GetLogger() *slog.Logger {
	return r.databaseLogger
}
func (r *registryImplementation) SetLogger(l *slog.Logger) {
	r.databaseLogger = l
}

func (r *registryImplementation) GetConsole() *slog.Logger {
	return r.consoleLogger
}
func (r *registryImplementation) SetConsole(l *slog.Logger) {
	r.consoleLogger = l
}

// Cache accessors (delegate to package-level cache singletons)
func (r *registryImplementation) GetMemoryCache() *ttlcache.Cache[string, any] {
	return cache.Memory
}

func (r *registryImplementation) SetMemoryCache(c *ttlcache.Cache[string, any]) {
	cache.Memory = c
}

func (r *registryImplementation) GetFileCache() cachego.Cache {
	return cache.File
}

func (r *registryImplementation) SetFileCache(c cachego.Cache) {
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

func (r *registryImplementation) GetAuditStore() auditstore.StoreInterface {
	return r.auditStore
}
func (r *registryImplementation) SetAuditStore(s auditstore.StoreInterface) {
	r.auditStore = s
}

// BlogStore
func (r *registryImplementation) GetBlogStore() blogstore.StoreInterface {
	return r.blogStore
}
func (r *registryImplementation) SetBlogStore(s blogstore.StoreInterface) {
	r.blogStore = s
}

// ChatStore
func (r *registryImplementation) GetChatStore() chatstore.StoreInterface {
	return r.chatStore
}
func (r *registryImplementation) SetChatStore(s chatstore.StoreInterface) {
	r.chatStore = s
}

// CacheStore
func (r *registryImplementation) GetCacheStore() cachestore.StoreInterface {
	return r.cacheStore
}
func (r *registryImplementation) SetCacheStore(s cachestore.StoreInterface) {
	r.cacheStore = s
}

// CmsStore
func (r *registryImplementation) GetCmsStore() cmsstore.StoreInterface {
	return r.cmsStore
}
func (r *registryImplementation) SetCmsStore(s cmsstore.StoreInterface) {
	r.cmsStore = s
}

// CustomStore
func (r *registryImplementation) GetCustomStore() customstore.StoreInterface {
	return r.customStore
}
func (r *registryImplementation) SetCustomStore(s customstore.StoreInterface) {
	r.customStore = s
}

// EntityStore
func (r *registryImplementation) GetEntityStore() entitystore.StoreInterface {
	return r.entityStore
}
func (r *registryImplementation) SetEntityStore(s entitystore.StoreInterface) {
	r.entityStore = s
}

// FeedStore
func (r *registryImplementation) GetFeedStore() feedstore.StoreInterface {
	return r.feedStore
}
func (r *registryImplementation) SetFeedStore(s feedstore.StoreInterface) {
	r.feedStore = s
}

// GeoStore
func (r *registryImplementation) GetGeoStore() geostore.StoreInterface {
	return r.geoStore
}
func (r *registryImplementation) SetGeoStore(s geostore.StoreInterface) {
	r.geoStore = s
}

// LogStore
func (r *registryImplementation) GetLogStore() logstore.StoreInterface {
	return r.logStore
}
func (r *registryImplementation) SetLogStore(s logstore.StoreInterface) {
	r.logStore = s
}

// MetaStore
func (r *registryImplementation) GetMetaStore() metastore.StoreInterface {
	return r.metaStore
}

func (r *registryImplementation) SetMetaStore(s metastore.StoreInterface) {
	r.metaStore = s
}

// SessionStore

// GetSessionStore returns the session store.
func (r *registryImplementation) GetSessionStore() sessionstore.StoreInterface {
	if r == nil {
		return nil
	}
	return r.sessionStore
}
func (r *registryImplementation) SetSessionStore(s sessionstore.StoreInterface) {
	r.sessionStore = s
}

// ShopStore
func (r *registryImplementation) GetShopStore() shopstore.StoreInterface {
	return r.shopStore
}
func (r *registryImplementation) SetShopStore(s shopstore.StoreInterface) {
	r.shopStore = s
}

// SqlFileStorage
func (r *registryImplementation) GetSqlFileStorage() filesystem.StorageInterface {
	return r.sqlFileStorage
}
func (r *registryImplementation) SetSqlFileStorage(s filesystem.StorageInterface) {
	r.sqlFileStorage = s
}

// StatsStore
func (r *registryImplementation) GetStatsStore() statsstore.StoreInterface {
	return r.statsStore
}
func (r *registryImplementation) SetStatsStore(s statsstore.StoreInterface) {
	r.statsStore = s
}

// TaskStore
func (r *registryImplementation) GetTaskStore() taskstore.StoreInterface {
	return r.taskStore
}
func (r *registryImplementation) SetTaskStore(s taskstore.StoreInterface) {
	r.taskStore = s
}

// UserStore
func (r *registryImplementation) GetUserStore() userstore.StoreInterface {
	return r.userStore
}
func (r *registryImplementation) SetUserStore(s userstore.StoreInterface) {
	r.userStore = s
}

// VaultStore
func (r *registryImplementation) GetVaultStore() vaultstore.StoreInterface {
	return r.vaultStore
}
func (r *registryImplementation) SetVaultStore(s vaultstore.StoreInterface) {
	r.vaultStore = s
}

// SettingStore
func (r *registryImplementation) GetSettingStore() settingstore.StoreInterface {
	return r.settingStore
}
func (r *registryImplementation) SetSettingStore(s settingstore.StoreInterface) {
	r.settingStore = s
}

func (r *registryImplementation) GetSubscriptionStore() subscriptionstore.StoreInterface {
	return r.subscriptionStore
}
func (r *registryImplementation) SetSubscriptionStore(s subscriptionstore.StoreInterface) {
	r.subscriptionStore = s
}

// Blind index stores
func (r *registryImplementation) GetBlindIndexStoreEmail() blindindexstore.StoreInterface {
	return r.blindIndexEmail
}
func (r *registryImplementation) SetBlindIndexStoreEmail(s blindindexstore.StoreInterface) {
	r.blindIndexEmail = s
}
func (r *registryImplementation) GetBlindIndexStoreFirstName() blindindexstore.StoreInterface {
	return r.blindIndexFirstName
}
func (r *registryImplementation) SetBlindIndexStoreFirstName(s blindindexstore.StoreInterface) {
	r.blindIndexFirstName = s
}
func (r *registryImplementation) GetBlindIndexStoreLastName() blindindexstore.StoreInterface {
	return r.blindIndexLastName
}
func (r *registryImplementation) SetBlindIndexStoreLastName(s blindindexstore.StoreInterface) {
	r.blindIndexLastName = s
}
