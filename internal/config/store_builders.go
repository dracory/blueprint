package config

import (
	"database/sql"

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
)

// ============================================================================
// == START: Store Builders
// ============================================================================
//
// This is where each store is constructed with its table names, options, and
// debug settings. Each builder is a pure function that takes a database
// connection and returns a fully configured store instance.
//
// To customize table names for your schema, edit the values in the
// NewStoreOptions structs below.
//
// ============================================================================

// NewAuditStore creates an audit store with the configured table name.
func NewAuditStore(db *sql.DB) (auditstore.StoreInterface, error) {
	return auditstore.NewStore(auditstore.NewStoreOptions{
		DB:             db,
		AuditTableName: "snv_audit_record",
	})
}

// NewBlogStore creates a blog store with the configured table names.
func NewBlogStore(db *sql.DB, debug bool) (blogstore.StoreInterface, error) {
	st, err := blogstore.NewStore(blogstore.NewStoreOptions{
		DB:                  db,
		PostTableName:       "snv_blogs_post",
		TaxonomyEnabled:     true,
		TaxonomyTableName:   "snv_blogs_taxonomy",
		TermTableName:       "snv_blogs_term",
		VersioningEnabled:   true,
		VersioningTableName: "snv_blogs_version",
		AutomigrateEnabled:  true,
	})
	if err != nil {
		return nil, err
	}
	st.EnableDebug(debug)
	return st, nil
}

// NewBlindIndexEmailStore creates a blind index store for email lookups.
func NewBlindIndexEmailStore(db *sql.DB) (blindindexstore.StoreInterface, error) {
	return blindindexstore.NewStore(blindindexstore.NewStoreOptions{
		DB:          db,
		TableName:   "snv_bindx_email",
		Transformer: &blindindexstore.Sha256Transformer{},
	})
}

// NewBlindIndexFirstNameStore creates a blind index store for first name lookups.
func NewBlindIndexFirstNameStore(db *sql.DB) (blindindexstore.StoreInterface, error) {
	return blindindexstore.NewStore(blindindexstore.NewStoreOptions{
		DB:          db,
		TableName:   "snv_bindx_first_name",
		Transformer: &blindindexstore.Sha256Transformer{},
	})
}

// NewBlindIndexLastNameStore creates a blind index store for last name lookups.
func NewBlindIndexLastNameStore(db *sql.DB) (blindindexstore.StoreInterface, error) {
	return blindindexstore.NewStore(blindindexstore.NewStoreOptions{
		DB:          db,
		TableName:   "snv_bindx_last_name",
		Transformer: &blindindexstore.Sha256Transformer{},
	})
}

// NewCacheStore creates a cache store with the configured table name.
func NewCacheStore(db *sql.DB, debug bool) (cachestore.StoreInterface, error) {
	st, err := cachestore.NewStore(cachestore.NewStoreOptions{
		DB:             db,
		CacheTableName: "snv_caches_cache",
	})
	if err != nil {
		return nil, err
	}
	st.EnableDebug(debug)
	return st, nil
}

// NewChatStore creates a chat store with the configured table names.
func NewChatStore(db *sql.DB) (chatstore.StoreInterface, error) {
	return chatstore.NewStore(chatstore.NewStoreOptions{
		DB:               db,
		TableChatName:    "snv_chat_chats",
		TableMessageName: "snv_chat_messages",
	})
}

// NewCmsStore creates a CMS store with the configured table names.
func NewCmsStore(db *sql.DB, debug bool) (cmsstore.StoreInterface, error) {
	st, err := cmsstore.NewStore(cmsstore.NewStoreOptions{
		DB:                   db,
		BlockTableName:       "snv_cms_block",
		PageTableName:        "snv_cms_page",
		TemplateTableName:    "snv_cms_template",
		SiteTableName:        "snv_cms_site",
		MenusEnabled:         true,
		MenuItemTableName:    "snv_cms_menu_item",
		MenuTableName:        "snv_cms_menu",
		TranslationsEnabled:  true,
		TranslationTableName: "snv_cms_translation",
		VersioningEnabled:    true,
		VersioningTableName:  "snv_cms_version",
	})
	if err != nil {
		return nil, err
	}
	st.EnableDebug(debug)
	return st, nil
}

// NewCustomStore creates a custom store with the configured table name.
func NewCustomStore(db *sql.DB, debug bool) (customstore.StoreInterface, error) {
	st, err := customstore.NewStore(customstore.NewStoreOptions{
		DB:        db,
		TableName: "snv_custom_record",
	})
	if err != nil {
		return nil, err
	}
	st.EnableDebug(debug)
	return st, nil
}

// NewEntityStore creates an entity store with the configured table names.
func NewEntityStore(db *sql.DB) (entitystore.StoreInterface, error) {
	return entitystore.NewStore(entitystore.NewStoreOptions{
		DB:                      db,
		EntityTableName:         "snv_entities_entity",
		EntityTrashTableName:    "snv_entities_entity_trash",
		AttributeTableName:      "snv_entities_attribute",
		AttributeTrashTableName: "snv_entities_attribute_trash",
	})
}

// NewFeedStore creates a feed store with the configured table names.
func NewFeedStore(db *sql.DB) (feedstore.StoreInterface, error) {
	return feedstore.NewStore(feedstore.NewStoreOptions{
		DB:            db,
		FeedTableName: "snv_feeds_feed",
		LinkTableName: "snv_feeds_link",
	})
}

// NewGeoStore creates a geo store with the configured table names.
func NewGeoStore(db *sql.DB) (geostore.StoreInterface, error) {
	return geostore.NewStore(geostore.NewStoreOptions{
		DB:                db,
		CountryTableName:  "snv_geo_country",
		StateTableName:    "snv_geo_state",
		TimezoneTableName: "snv_geo_timezone",
	})
}

// NewLogStore creates a log store with the configured table name.
func NewLogStore(db *sql.DB, debug bool) (logstore.StoreInterface, error) {
	st, err := logstore.NewStore(logstore.NewStoreOptions{
		DB:           db,
		LogTableName: "snv_logs_log",
	})
	if err != nil {
		return nil, err
	}
	st.EnableDebug(debug)
	return st, nil
}

// NewMetaStore creates a meta store with the configured table name.
func NewMetaStore(db *sql.DB, debug bool) (metastore.StoreInterface, error) {
	st, err := metastore.NewStore(metastore.NewStoreOptions{
		DB:            db,
		MetaTableName: "snv_metas_meta",
	})
	if err != nil {
		return nil, err
	}
	st.EnableDebug(debug)
	return st, nil
}

// NewSessionStore creates a session store with the configured table name and timeout.
func NewSessionStore(db *sql.DB, debug bool, isDev bool) (sessionstore.StoreInterface, error) {
	timeoutSeconds := int64(7200) // 2 hours default
	if isDev {
		timeoutSeconds = 14400 // 4 hours in development
	}
	st, err := sessionstore.NewStore(sessionstore.NewStoreOptions{
		DB:               db,
		SessionTableName: "snv_sessions_session",
		TimeoutSeconds:   timeoutSeconds,
	})
	if err != nil {
		return nil, err
	}
	st.EnableDebug(debug)
	return st, nil
}

// NewSettingStore creates a setting store with the configured table name.
func NewSettingStore(db *sql.DB) (settingstore.StoreInterface, error) {
	return settingstore.NewStore(settingstore.NewStoreOptions{
		DB:               db,
		SettingTableName: "snv_settings",
	})
}

// NewShopStore creates a shop store with the configured table names.
func NewShopStore(db *sql.DB, debug bool) (shopstore.StoreInterface, error) {
	st, err := shopstore.NewStore(shopstore.NewStoreOptions{
		DB:                     db,
		CategoryTableName:      "snv_shop_category",
		DiscountTableName:      "snv_shop_discount",
		MediaTableName:         "snv_shop_media",
		OrderTableName:         "snv_shop_order",
		OrderLineItemTableName: "snv_shop_order_line_item",
		ProductTableName:       "snv_shop_product",
	})
	if err != nil {
		return nil, err
	}
	st.EnableDebug(debug)
	return st, nil
}

// NewSqlFileStorage creates a SQL-backed file storage with the configured table name.
func NewSqlFileStorage(db *sql.DB) (filesystem.StorageInterface, error) {
	return filesystem.NewStorage(filesystem.Disk{
		DiskName:  filesystem.DRIVER_SQL,
		Driver:    filesystem.DRIVER_SQL,
		Url:       "/files",
		DB:        db,
		TableName: "snv_files_file",
	})
}

// NewStatsStore creates a stats store with the configured table name.
func NewStatsStore(db *sql.DB, debug bool) (statsstore.StoreInterface, error) {
	st, err := statsstore.NewStore(statsstore.NewStoreOptions{
		DB:               db,
		VisitorTableName: "snv_stats_visitor",
	})
	if err != nil {
		return nil, err
	}
	st.EnableDebug(debug)
	return st, nil
}

// NewSubscriptionStore creates a subscription store with the configured table names.
func NewSubscriptionStore(db *sql.DB) (subscriptionstore.StoreInterface, error) {
	return subscriptionstore.NewStore(subscriptionstore.NewStoreOptions{
		DB:                    db,
		PlanTableName:         "snv_subscriptions_plan",
		SubscriptionTableName: "snv_subscriptions_subscription",
	})
}

// NewTaskStore creates a task store with the configured table names.
func NewTaskStore(db *sql.DB, debug bool) (taskstore.StoreInterface, error) {
	st, err := taskstore.NewStore(taskstore.NewStoreOptions{
		DB:                      db,
		ScheduleTableName:       "snv_tasks_schedule",
		TaskDefinitionTableName: "snv_tasks_task_definition",
		TaskQueueTableName:      "snv_tasks_task_queue",
	})
	if err != nil {
		return nil, err
	}
	st.EnableDebug(debug)
	return st, nil
}

// NewUserStore creates a user store with the configured table name.
func NewUserStore(db *sql.DB) (userstore.StoreInterface, error) {
	return userstore.NewStore(userstore.NewStoreOptions{
		DB:            db,
		UserTableName: "snv_users_user",
	})
}

// NewVaultStore creates a vault store with the configured table names.
func NewVaultStore(db *sql.DB, debug bool) (vaultstore.StoreInterface, error) {
	st, err := vaultstore.NewStore(vaultstore.NewStoreOptions{
		DB:                 db,
		VaultTableName:     "snv_vault_vault",
		VaultMetaTableName: "snv_vault_meta",
		PasswordMinLength:  6,
	})
	if err != nil {
		return nil, err
	}
	st.EnableDebug(debug)
	return st, nil
}

// ============================================================================
// == END: Store Builders
// ============================================================================
