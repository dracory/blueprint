package types

import (
	"database/sql"
	"log/slog"

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
	"github.com/dracory/userstore"
	"github.com/dracory/vaultstore"
	"github.com/faabiosr/cachego"
	"github.com/gouniverse/filesystem"
	"github.com/jellydator/ttlcache/v3"
)

// AppInterface defines accessors for application-scoped runtime services.
// It mirrors the style of ConfigInterface, enabling DI and testability.
//
// Typical implementations will wire these in an Initialize step.
// For now, we also provide a runtime adapter in internal/app that returns
// the current process-level resources.
type AppInterface interface {
	// Logger
	GetLogger() *slog.Logger
	SetLogger(l *slog.Logger)

	// Config
	GetConfig() ConfigInterface
	SetConfig(c ConfigInterface)

	// Caches
	GetMemoryCache() *ttlcache.Cache[string, any]
	SetMemoryCache(c *ttlcache.Cache[string, any])
	GetFileCache() cachego.Cache
	SetFileCache(c cachego.Cache)

	// DB
	GetDB() *sql.DB
	SetDB(db *sql.DB)

	// ========================================================================
	// == Stores (all specific data stores)
	// ========================================================================

	// Blog store
	GetBlogStore() blogstore.StoreInterface
	SetBlogStore(s blogstore.StoreInterface)

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
