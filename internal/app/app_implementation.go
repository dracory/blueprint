package app

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
	neatdatabase "github.com/dracory/neat/database"
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

// appImplementation is the orchestration facade for starting the app.
// It encapsulates configuration and database (container removed).
type appImplementation struct {
	cfg config.ConfigInterface

	// Database
	neatDB *neatdatabase.Database
	db     *sql.DB

	// Loggers
	databaseLogger *slog.Logger
	consoleLogger  *slog.Logger

	// Caches (instance-scoped)
	memoryCache *ttlcache.Cache[string, any]
	fileCache   cachego.Cache

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

// Ensure appImplementation satisfies the AppInterface contract.
var _ AppInterface = (*appImplementation)(nil)

// New constructs and initializes the app (logger, caches, database).
// It centralizes the boot logic so callers only use this single constructor.
func New(cfg config.ConfigInterface) (AppInterface, error) {
	if cfg == nil {
		return nil, errors.New("cfg is nil")
	}

	// Caches (instance-scoped)
	memoryCache := ttlcache.New[string, any]()
	cacheDir := cacheDirectory()
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return nil, err
	}
	fileCache := file.New(cacheDir)

	consoleLogger := slog.New(tint.NewHandler(os.Stdout, nil))

	// Database open
	neatDB, err := databaseOpen(cfg)
	if err != nil {
		return nil, err
	}

	db, err := neatDB.DB()
	if err != nil {
		return nil, err
	}

	// Build app instance
	app := &appImplementation{cfg: cfg}
	app.SetConsole(consoleLogger)
	app.SetLogger(consoleLogger)
	app.SetMemoryCache(memoryCache)
	app.SetFileCache(fileCache)
	app.SetNeatDatabase(neatDB)
	app.SetDatabase(db)

	// Initialize stores
	if err := app.dataStoresInitialize(); err != nil {
		return nil, err
	}

	if app.GetLogStore() != nil {
		app.SetLogger(slog.New(logstore.NewSlogHandler(app.GetLogStore())))
	}

	// Mirror caches into internal/cache for transitional compatibility
	cache.Memory = memoryCache
	cache.File = fileCache

	return app, nil
}

// Close closes the app and its resources
func (r *appImplementation) Close() error {
	if r == nil {
		return nil
	}

	if r.neatDB == nil {
		return nil
	}

	err := r.neatDB.Close()
	r.neatDB = nil
	r.db = nil
	return err
}

// GetConfig returns the app config
func (r *appImplementation) GetConfig() config.ConfigInterface {
	if r == nil {
		return nil
	}
	return r.cfg
}

func (r *appImplementation) SetConfig(cfg config.ConfigInterface) {
	r.cfg = cfg
}

// GetDatabase returns the app database
func (r *appImplementation) GetDatabase() *sql.DB {
	return r.db
}

// SetDatabase sets the app database
func (r *appImplementation) SetDatabase(db *sql.DB) {
	r.db = db
}

// GetNeatDatabase returns the neat database instance.
func (r *appImplementation) GetNeatDatabase() *neatdatabase.Database {
	return r.neatDB
}

// SetNeatDatabase sets the neat database instance.
func (r *appImplementation) SetNeatDatabase(db *neatdatabase.Database) {
	r.neatDB = db
}

// GetDatabaseConnection returns the underlying *sql.DB for the named connection.
// If the name is empty, it returns the default connection.
func (r *appImplementation) GetDatabaseConnection(name string) *sql.DB {
	if r == nil || r.neatDB == nil {
		return nil
	}
	if name == "" || name == r.cfg.GetDatabaseDefaultConnection() {
		return r.db
	}
	conn, err := r.neatDB.Connection(name)
	if err != nil || conn == nil {
		return nil
	}
	db, err := conn.DB()
	if err != nil {
		return nil
	}
	return db
}

// Logger accessors (delegate to package-level logger singletons)
func (r *appImplementation) GetLogger() *slog.Logger {
	return r.databaseLogger
}

func (r *appImplementation) SetLogger(l *slog.Logger) {
	r.databaseLogger = l
}

func (r *appImplementation) GetConsole() *slog.Logger {
	return r.consoleLogger
}

func (r *appImplementation) SetConsole(l *slog.Logger) {
	r.consoleLogger = l
}

// Cache accessors (instance-scoped)
func (r *appImplementation) GetMemoryCache() *ttlcache.Cache[string, any] {
	if r == nil {
		return nil
	}
	return r.memoryCache
}

func (r *appImplementation) SetMemoryCache(c *ttlcache.Cache[string, any]) {
	if r == nil {
		return
	}
	r.memoryCache = c
}

func (r *appImplementation) GetFileCache() cachego.Cache {
	if r == nil {
		return nil
	}
	return r.fileCache
}

func (r *appImplementation) SetFileCache(c cachego.Cache) {
	if r == nil {
		return
	}
	r.fileCache = c
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

func (r *appImplementation) GetAuditStore() auditstore.StoreInterface {
	return r.auditStore
}
func (r *appImplementation) SetAuditStore(s auditstore.StoreInterface) {
	r.auditStore = s
}

// BlogStore
func (r *appImplementation) GetBlogStore() blogstore.StoreInterface {
	return r.blogStore
}
func (r *appImplementation) SetBlogStore(s blogstore.StoreInterface) {
	r.blogStore = s
}

// ChatStore
func (r *appImplementation) GetChatStore() chatstore.StoreInterface {
	return r.chatStore
}
func (r *appImplementation) SetChatStore(s chatstore.StoreInterface) {
	r.chatStore = s
}

// CacheStore
func (r *appImplementation) GetCacheStore() cachestore.StoreInterface {
	return r.cacheStore
}
func (r *appImplementation) SetCacheStore(s cachestore.StoreInterface) {
	r.cacheStore = s
}

// CmsStore
func (r *appImplementation) GetCmsStore() cmsstore.StoreInterface {
	return r.cmsStore
}
func (r *appImplementation) SetCmsStore(s cmsstore.StoreInterface) {
	r.cmsStore = s
}

// CustomStore
func (r *appImplementation) GetCustomStore() customstore.StoreInterface {
	return r.customStore
}
func (r *appImplementation) SetCustomStore(s customstore.StoreInterface) {
	r.customStore = s
}

// EntityStore
func (r *appImplementation) GetEntityStore() entitystore.StoreInterface {
	return r.entityStore
}
func (r *appImplementation) SetEntityStore(s entitystore.StoreInterface) {
	r.entityStore = s
}

// FeedStore
func (r *appImplementation) GetFeedStore() feedstore.StoreInterface {
	return r.feedStore
}
func (r *appImplementation) SetFeedStore(s feedstore.StoreInterface) {
	r.feedStore = s
}

// GeoStore
func (r *appImplementation) GetGeoStore() geostore.StoreInterface {
	return r.geoStore
}
func (r *appImplementation) SetGeoStore(s geostore.StoreInterface) {
	r.geoStore = s
}

// LogStore
func (r *appImplementation) GetLogStore() logstore.StoreInterface {
	return r.logStore
}
func (r *appImplementation) SetLogStore(s logstore.StoreInterface) {
	r.logStore = s
}

// MetaStore
func (r *appImplementation) GetMetaStore() metastore.StoreInterface {
	return r.metaStore
}

func (r *appImplementation) SetMetaStore(s metastore.StoreInterface) {
	r.metaStore = s
}

// SessionStore

// GetSessionStore returns the session store.
func (r *appImplementation) GetSessionStore() sessionstore.StoreInterface {
	if r == nil {
		return nil
	}
	return r.sessionStore
}
func (r *appImplementation) SetSessionStore(s sessionstore.StoreInterface) {
	r.sessionStore = s
}

// ShopStore
func (r *appImplementation) GetShopStore() shopstore.StoreInterface {
	return r.shopStore
}
func (r *appImplementation) SetShopStore(s shopstore.StoreInterface) {
	r.shopStore = s
}

// SqlFileStorage
func (r *appImplementation) GetSqlFileStorage() filesystem.StorageInterface {
	return r.sqlFileStorage
}
func (r *appImplementation) SetSqlFileStorage(s filesystem.StorageInterface) {
	r.sqlFileStorage = s
}

// StatsStore
func (r *appImplementation) GetStatsStore() statsstore.StoreInterface {
	return r.statsStore
}
func (r *appImplementation) SetStatsStore(s statsstore.StoreInterface) {
	r.statsStore = s
}

// TaskStore
func (r *appImplementation) GetTaskStore() taskstore.StoreInterface {
	return r.taskStore
}
func (r *appImplementation) SetTaskStore(s taskstore.StoreInterface) {
	r.taskStore = s
}

// UserStore
func (r *appImplementation) GetUserStore() userstore.StoreInterface {
	return r.userStore
}
func (r *appImplementation) SetUserStore(s userstore.StoreInterface) {
	r.userStore = s
}

// VaultStore
func (r *appImplementation) GetVaultStore() vaultstore.StoreInterface {
	return r.vaultStore
}
func (r *appImplementation) SetVaultStore(s vaultstore.StoreInterface) {
	r.vaultStore = s
}

// SettingStore
func (r *appImplementation) GetSettingStore() settingstore.StoreInterface {
	return r.settingStore
}
func (r *appImplementation) SetSettingStore(s settingstore.StoreInterface) {
	r.settingStore = s
}

func (r *appImplementation) GetSubscriptionStore() subscriptionstore.StoreInterface {
	return r.subscriptionStore
}
func (r *appImplementation) SetSubscriptionStore(s subscriptionstore.StoreInterface) {
	r.subscriptionStore = s
}

// Blind index stores
func (r *appImplementation) GetBlindIndexStoreEmail() blindindexstore.StoreInterface {
	return r.blindIndexEmail
}
func (r *appImplementation) SetBlindIndexStoreEmail(s blindindexstore.StoreInterface) {
	r.blindIndexEmail = s
}
func (r *appImplementation) GetBlindIndexStoreFirstName() blindindexstore.StoreInterface {
	return r.blindIndexFirstName
}
func (r *appImplementation) SetBlindIndexStoreFirstName(s blindindexstore.StoreInterface) {
	r.blindIndexFirstName = s
}
func (r *appImplementation) GetBlindIndexStoreLastName() blindindexstore.StoreInterface {
	return r.blindIndexLastName
}
func (r *appImplementation) SetBlindIndexStoreLastName(s blindindexstore.StoreInterface) {
	r.blindIndexLastName = s
}
