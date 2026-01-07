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
//
// This interface is intended to be used at the composition root (startup wiring)
// and for edge integration points. Prefer passing narrower dependency
// interfaces to controllers/tasks.
type RegistryInterface interface {
	Close() error

	GetLogger() *slog.Logger
	SetLogger(l *slog.Logger)

	GetConfig() config.ConfigInterface
	SetConfig(c config.ConfigInterface)

	GetMemoryCache() *ttlcache.Cache[string, any]
	SetMemoryCache(c *ttlcache.Cache[string, any])
	GetFileCache() cachego.Cache
	SetFileCache(c cachego.Cache)

	GetDatabase() *sql.DB
	SetDatabase(db *sql.DB)

	GetAuditStore() auditstore.StoreInterface
	SetAuditStore(s auditstore.StoreInterface)

	GetBlogStore() blogstore.StoreInterface
	SetBlogStore(s blogstore.StoreInterface)

	GetChatStore() chatstore.StoreInterface
	SetChatStore(s chatstore.StoreInterface)

	GetBlindIndexStoreEmail() blindindexstore.StoreInterface
	SetBlindIndexStoreEmail(s blindindexstore.StoreInterface)

	GetBlindIndexStoreFirstName() blindindexstore.StoreInterface
	SetBlindIndexStoreFirstName(s blindindexstore.StoreInterface)

	GetBlindIndexStoreLastName() blindindexstore.StoreInterface
	SetBlindIndexStoreLastName(s blindindexstore.StoreInterface)

	GetCacheStore() cachestore.StoreInterface
	SetCacheStore(s cachestore.StoreInterface)

	GetCmsStore() cmsstore.StoreInterface
	SetCmsStore(s cmsstore.StoreInterface)

	GetCustomStore() customstore.StoreInterface
	SetCustomStore(s customstore.StoreInterface)

	GetEntityStore() entitystore.StoreInterface
	SetEntityStore(s entitystore.StoreInterface)

	GetFeedStore() feedstore.StoreInterface
	SetFeedStore(s feedstore.StoreInterface)

	GetGeoStore() geostore.StoreInterface
	SetGeoStore(s geostore.StoreInterface)

	GetLogStore() logstore.StoreInterface
	SetLogStore(s logstore.StoreInterface)

	GetMetaStore() metastore.StoreInterface
	SetMetaStore(s metastore.StoreInterface)

	GetSessionStore() sessionstore.StoreInterface
	SetSessionStore(s sessionstore.StoreInterface)

	GetSettingStore() settingstore.StoreInterface
	SetSettingStore(s settingstore.StoreInterface)

	GetShopStore() shopstore.StoreInterface
	SetShopStore(s shopstore.StoreInterface)

	GetSqlFileStorage() filesystem.StorageInterface
	SetSqlFileStorage(s filesystem.StorageInterface)

	GetStatsStore() statsstore.StoreInterface
	SetStatsStore(s statsstore.StoreInterface)

	GetSubscriptionStore() subscriptionstore.StoreInterface
	SetSubscriptionStore(s subscriptionstore.StoreInterface)

	GetTaskStore() taskstore.StoreInterface
	SetTaskStore(s taskstore.StoreInterface)

	GetUserStore() userstore.StoreInterface
	SetUserStore(s userstore.StoreInterface)

	GetVaultStore() vaultstore.StoreInterface
	SetVaultStore(s vaultstore.StoreInterface)
}
