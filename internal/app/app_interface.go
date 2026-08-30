package app

import (
	"database/sql"
	"log/slog"

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
	"github.com/jellydator/ttlcache/v3"
)

// AppInterface defines accessors for app-scoped runtime services.
// It mirrors the style of ConfigInterface, enabling DI and testability.
//
// This interface is intended to be used at the composition root (startup wiring)
// and for edge integration points. Prefer passing narrower dependency
// interfaces to controllers/tasks.
type AppInterface interface {
	Close() error

	// Logger
	GetLogger() *slog.Logger
	SetLogger(l *slog.Logger)

	// Config
	GetConfig() config.ConfigInterface
	SetConfig(c config.ConfigInterface)

	// Caches
	GetMemoryCache() *ttlcache.Cache[string, any]
	SetMemoryCache(c *ttlcache.Cache[string, any])
	GetFileCache() cachego.Cache
	SetFileCache(c cachego.Cache)

	// DB
	GetDatabase() *sql.DB
	SetDatabase(db *sql.DB)

	GetNeatDatabase() *neatdatabase.Database
	SetNeatDatabase(db *neatdatabase.Database)

	GetDatabaseConnection(name string) *sql.DB

	// ========================================================================
	// == Stores (all specific data stores)
	// ========================================================================

	// Audit store
	GetAuditStore() auditstore.StoreInterface
	SetAuditStore(s auditstore.StoreInterface)
	IsEnabledAuditStore() bool
	IsDisabledAuditStore() bool

	// Blog store
	GetBlogStore() blogstore.StoreInterface
	SetBlogStore(s blogstore.StoreInterface)
	IsEnabledBlogStore() bool
	IsDisabledBlogStore() bool

	// Chat store
	GetChatStore() chatstore.StoreInterface
	SetChatStore(s chatstore.StoreInterface)
	IsEnabledChatStore() bool
	IsDisabledChatStore() bool

	// Blind index store
	GetBlindIndexStoreEmail() blindindexstore.StoreInterface
	SetBlindIndexStoreEmail(s blindindexstore.StoreInterface)
	IsEnabledBlindIndexStoreEmail() bool
	IsDisabledBlindIndexStoreEmail() bool

	// Blind index store
	GetBlindIndexStoreFirstName() blindindexstore.StoreInterface
	SetBlindIndexStoreFirstName(s blindindexstore.StoreInterface)
	IsEnabledBlindIndexStoreFirstName() bool
	IsDisabledBlindIndexStoreFirstName() bool

	// Blind index store
	GetBlindIndexStoreLastName() blindindexstore.StoreInterface
	SetBlindIndexStoreLastName(s blindindexstore.StoreInterface)
	IsEnabledBlindIndexStoreLastName() bool
	IsDisabledBlindIndexStoreLastName() bool

	// Cache store
	GetCacheStore() cachestore.StoreInterface
	SetCacheStore(s cachestore.StoreInterface)
	IsEnabledCacheStore() bool
	IsDisabledCacheStore() bool

	// CMS store
	GetCmsStore() cmsstore.StoreInterface
	SetCmsStore(s cmsstore.StoreInterface)
	IsEnabledCmsStore() bool
	IsDisabledCmsStore() bool

	// Custom store
	GetCustomStore() customstore.StoreInterface
	SetCustomStore(s customstore.StoreInterface)
	IsEnabledCustomStore() bool
	IsDisabledCustomStore() bool

	// Entity store
	GetEntityStore() entitystore.StoreInterface
	SetEntityStore(s entitystore.StoreInterface)
	IsEnabledEntityStore() bool
	IsDisabledEntityStore() bool

	// Feed store
	GetFeedStore() feedstore.StoreInterface
	SetFeedStore(s feedstore.StoreInterface)
	IsEnabledFeedStore() bool
	IsDisabledFeedStore() bool

	// Geo store
	GetGeoStore() geostore.StoreInterface
	SetGeoStore(s geostore.StoreInterface)
	IsEnabledGeoStore() bool
	IsDisabledGeoStore() bool

	// Log store
	GetLogStore() logstore.StoreInterface
	SetLogStore(s logstore.StoreInterface)
	IsEnabledLogStore() bool
	IsDisabledLogStore() bool

	// Meta store
	GetMetaStore() metastore.StoreInterface
	SetMetaStore(s metastore.StoreInterface)
	IsEnabledMetaStore() bool
	IsDisabledMetaStore() bool

	// Session store
	GetSessionStore() sessionstore.StoreInterface
	SetSessionStore(s sessionstore.StoreInterface)
	IsEnabledSessionStore() bool
	IsDisabledSessionStore() bool

	// Setting store
	GetSettingStore() settingstore.StoreInterface
	SetSettingStore(s settingstore.StoreInterface)
	IsEnabledSettingStore() bool
	IsDisabledSettingStore() bool

	// Shop store
	GetShopStore() shopstore.StoreInterface
	SetShopStore(s shopstore.StoreInterface)
	IsEnabledShopStore() bool
	IsDisabledShopStore() bool

	// SQL file storage
	GetSqlFileStorage() filesystem.StorageInterface
	SetSqlFileStorage(s filesystem.StorageInterface)
	IsEnabledSqlFileStorage() bool
	IsDisabledSqlFileStorage() bool

	// Stats store
	GetStatsStore() statsstore.StoreInterface
	SetStatsStore(s statsstore.StoreInterface)
	IsEnabledStatsStore() bool
	IsDisabledStatsStore() bool

	// Subscription store
	GetSubscriptionStore() subscriptionstore.StoreInterface
	SetSubscriptionStore(s subscriptionstore.StoreInterface)
	IsEnabledSubscriptionStore() bool
	IsDisabledSubscriptionStore() bool

	// Task store
	GetTaskStore() taskstore.StoreInterface
	SetTaskStore(s taskstore.StoreInterface)
	IsEnabledTaskStore() bool
	IsDisabledTaskStore() bool

	// User store
	GetUserStore() userstore.StoreInterface
	SetUserStore(s userstore.StoreInterface)
	IsEnabledUserStore() bool
	IsDisabledUserStore() bool

	// Vault store
	GetVaultStore() vaultstore.StoreInterface
	SetVaultStore(s vaultstore.StoreInterface)
	IsEnabledVaultStore() bool
	IsDisabledVaultStore() bool
}
