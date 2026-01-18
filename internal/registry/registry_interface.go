package registry

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

// RegistryInterface defines accessors for registry-scoped runtime services.
// It mirrors the style of ConfigInterface, enabling DI and testability.
//
// This interface is intended to be used at the composition root (startup wiring)
// and for edge integration points. Prefer passing narrower dependency
// interfaces to controllers/tasks.
type RegistryInterface interface {
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

	// ========================================================================
	// == Stores (all specific data stores)
	// ========================================================================

    // Audit store
	GetAuditStore() auditstore.StoreInterface
	SetAuditStore(s auditstore.StoreInterface)

	// Blog store
	GetBlogStore() blogstore.StoreInterface
	SetBlogStore(s blogstore.StoreInterface)

	// Chat store
	GetChatStore() chatstore.StoreInterface
	SetChatStore(s chatstore.StoreInterface)

	// Blind index store
	GetBlindIndexStoreEmail() blindindexstore.StoreInterface
	SetBlindIndexStoreEmail(s blindindexstore.StoreInterface)

	// Blind index store
	GetBlindIndexStoreFirstName() blindindexstore.StoreInterface
	SetBlindIndexStoreFirstName(s blindindexstore.StoreInterface)

	// Blind index store
	GetBlindIndexStoreLastName() blindindexstore.StoreInterface
	SetBlindIndexStoreLastName(s blindindexstore.StoreInterface)

	// Cache store
	GetCacheStore() cachestore.StoreInterface
	SetCacheStore(s cachestore.StoreInterface)

	// CMS store
	GetCmsStore() cmsstore.StoreInterface
	SetCmsStore(s cmsstore.StoreInterface)

	// Custom store
	GetCustomStore() customstore.StoreInterface
	SetCustomStore(s customstore.StoreInterface)

	// Entity store
	GetEntityStore() entitystore.StoreInterface
	SetEntityStore(s entitystore.StoreInterface)

	// Feed store
	GetFeedStore() feedstore.StoreInterface
	SetFeedStore(s feedstore.StoreInterface)

	// Geo store
	GetGeoStore() geostore.StoreInterface
	SetGeoStore(s geostore.StoreInterface)

	// Log store
	GetLogStore() logstore.StoreInterface
	SetLogStore(s logstore.StoreInterface)

	// Meta store
	GetMetaStore() metastore.StoreInterface
	SetMetaStore(s metastore.StoreInterface)

	// Session store
	GetSessionStore() sessionstore.StoreInterface
	SetSessionStore(s sessionstore.StoreInterface)

	// Setting store
	GetSettingStore() settingstore.StoreInterface
	SetSettingStore(s settingstore.StoreInterface)

	// Shop store
	GetShopStore() shopstore.StoreInterface
	SetShopStore(s shopstore.StoreInterface)

	// SQL file storage
	GetSqlFileStorage() filesystem.StorageInterface
	SetSqlFileStorage(s filesystem.StorageInterface)

	// Stats store
	GetStatsStore() statsstore.StoreInterface
	SetStatsStore(s statsstore.StoreInterface)

	// Subscription store
	GetSubscriptionStore() subscriptionstore.StoreInterface
	SetSubscriptionStore(s subscriptionstore.StoreInterface)

	// Task store
	GetTaskStore() taskstore.StoreInterface
	SetTaskStore(s taskstore.StoreInterface)

	// User store
	GetUserStore() userstore.StoreInterface
	SetUserStore(s userstore.StoreInterface)

	// Vault store
	GetVaultStore() vaultstore.StoreInterface
	SetVaultStore(s vaultstore.StoreInterface)
}
