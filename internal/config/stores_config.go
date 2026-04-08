package config

import (
	"fmt"

	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// ============================================================================
// Interfaces
// ============================================================================

// AuditStoreConfigInterface defines audit store configuration methods.
type AuditStoreConfigInterface interface {
	SetAuditStoreUsed(bool)
	GetAuditStoreUsed() bool
}

// BlogStoreConfigInterface defines blog store configuration methods.
type BlogStoreConfigInterface interface {
	SetBlogStoreUsed(bool)
	GetBlogStoreUsed() bool
}

// CacheStoreConfigInterface defines cache store configuration methods.
type CacheStoreConfigInterface interface {
	SetCacheStoreUsed(bool)
	GetCacheStoreUsed() bool
}

// ChatStoreConfigInterface defines chat store configuration methods.
type ChatStoreConfigInterface interface {
	SetChatStoreUsed(bool)
	GetChatStoreUsed() bool
}

// CmsStoreConfigInterface defines CMS store configuration methods.
type CmsStoreConfigInterface interface {
	SetCmsStoreUsed(bool)
	GetCmsStoreUsed() bool

	SetCmsStoreTemplateID(string)
	GetCmsStoreTemplateID() string
}

// CustomStoreConfigInterface defines custom store configuration methods.
type CustomStoreConfigInterface interface {
	SetCustomStoreUsed(bool)
	GetCustomStoreUsed() bool
}

// EntityStoreConfigInterface defines entity store configuration methods.
type EntityStoreConfigInterface interface {
	SetEntityStoreUsed(bool)
	GetEntityStoreUsed() bool
}

// FeedStoreConfigInterface defines feed store configuration methods.
type FeedStoreConfigInterface interface {
	SetFeedStoreUsed(bool)
	GetFeedStoreUsed() bool
}

// GeoStoreConfigInterface defines geo store configuration methods.
type GeoStoreConfigInterface interface {
	SetGeoStoreUsed(bool)
	GetGeoStoreUsed() bool
}

// LogStoreConfigInterface defines log store configuration methods.
type LogStoreConfigInterface interface {
	SetLogStoreUsed(bool)
	GetLogStoreUsed() bool
}

// MetaStoreConfigInterface defines meta store configuration methods.
type MetaStoreConfigInterface interface {
	SetMetaStoreUsed(bool)
	GetMetaStoreUsed() bool
}

// SessionStoreConfigInterface defines session store configuration methods.
type SessionStoreConfigInterface interface {
	SetSessionStoreUsed(bool)
	GetSessionStoreUsed() bool
}

// SettingStoreConfigInterface defines setting store configuration methods.
type SettingStoreConfigInterface interface {
	SetSettingStoreUsed(bool)
	GetSettingStoreUsed() bool
}

// ShopStoreConfigInterface defines shop store configuration methods.
type ShopStoreConfigInterface interface {
	SetShopStoreUsed(bool)
	GetShopStoreUsed() bool
}

// SqlFileStoreConfigInterface defines SQL file store configuration methods.
type SqlFileStoreConfigInterface interface {
	SetSqlFileStoreUsed(bool)
	GetSqlFileStoreUsed() bool
}

// StatsStoreConfigInterface defines stats store configuration methods.
type StatsStoreConfigInterface interface {
	SetStatsStoreUsed(bool)
	GetStatsStoreUsed() bool
}

// SubscriptionStoreConfigInterface defines subscription store configuration methods.
type SubscriptionStoreConfigInterface interface {
	SetSubscriptionStoreUsed(bool)
	GetSubscriptionStoreUsed() bool
}

// TaskStoreConfigInterface defines task store configuration methods.
type TaskStoreConfigInterface interface {
	SetTaskStoreUsed(bool)
	GetTaskStoreUsed() bool
}

// UserStoreConfigInterface defines user store configuration methods.
type UserStoreConfigInterface interface {
	SetUserStoreUsed(bool)
	GetUserStoreUsed() bool

	SetUserStoreVaultEnabled(bool)
	GetUserStoreVaultEnabled() bool
}

// VaultStoreConfigInterface defines vault store configuration methods.
type VaultStoreConfigInterface interface {
	SetVaultStoreUsed(bool)
	GetVaultStoreUsed() bool

	SetVaultStoreKey(string)
	GetVaultStoreKey() string
}

// ============================================================================
// Types
// ============================================================================

// storesConfig captures feature store toggles.
type storesConfig struct {
	auditStoreUsed        bool   // Enable audit logging store
	blogStoreUsed         bool   // Enable blog content store
	cacheStoreUsed        bool   // Enable caching store
	chatStoreUsed         bool   // Enable chat store
	cmsStoreUsed          bool   // Enable CMS content store
	cmsStoreTemplateID    string // Default template ID for CMS store
	customStoreUsed       bool   // Enable custom data store
	entityStoreUsed       bool   // Enable entity management store
	feedStoreUsed         bool   // Enable RSS/Atom feed store
	geoStoreUsed          bool   // Enable geolocation data store
	logStoreUsed          bool   // Enable application logging store
	metaStoreUsed         bool   // Enable metadata store
	sessionStoreUsed      bool   // Enable user session store
	settingStoreUsed      bool   // Enable application settings store
	shopStoreUsed         bool   // Enable e-commerce store
	sqlFileStoreUsed      bool   // Enable SQL file storage store
	statsStoreUsed        bool   // Enable analytics/statistics store
	subscriptionStoreUsed bool   // Enable subscription management store
	taskStoreUsed         bool   // Enable background task store
	userStoreUsed         bool   // Enable user management store
	userStoreVaultEnabled bool   // Enable vault encryption for user store
	vaultStoreUsed        bool   // Enable secure vault store
	vaultStoreKey         string // Encryption key for vault store
}

// ============================================================================
// Loader
// ============================================================================

// loadStoresConfig loads stores configuration from environment variables.
func loadStoresConfig(acc *baseCfg.LoadAccumulator) storesConfig {
	cmsStoreTemplateID := env.GetString(KEY_CMS_STORE_TEMPLATE_ID)
	vaultStoreKey := env.GetString(KEY_VAULT_STORE_KEY)

	if userStoreVaultEnabled && !vaultStoreUsed {
		acc.Add(fmt.Errorf("%v requires %v to be true", userStoreVaultEnabled, vaultStoreUsed))
	}

	acc.MustWhen(cmsStoreUsed, KEY_CMS_STORE_TEMPLATE_ID, "required when `CMS_STORE_USED` is true", cmsStoreTemplateID)

	return storesConfig{
		auditStoreUsed:        auditStoreUsed,
		blogStoreUsed:         blogStoreUsed,
		cacheStoreUsed:        cacheStoreUsed,
		chatStoreUsed:         chatStoreUsed,
		cmsStoreUsed:          cmsStoreUsed,
		cmsStoreTemplateID:    cmsStoreTemplateID,
		customStoreUsed:       customStoreUsed,
		entityStoreUsed:       entityStoreUsed,
		feedStoreUsed:         feedStoreUsed,
		geoStoreUsed:          geoStoreUsed,
		logStoreUsed:          logStoreUsed,
		metaStoreUsed:         metaStoreUsed,
		sessionStoreUsed:      sessionStoreUsed,
		settingStoreUsed:      settingStoreUsed,
		shopStoreUsed:         shopStoreUsed,
		sqlFileStoreUsed:      sqlFileStoreUsed,
		statsStoreUsed:        statsStoreUsed,
		subscriptionStoreUsed: subscriptionStoreUsed,
		taskStoreUsed:         taskStoreUsed,
		userStoreUsed:         userStoreUsed,
		userStoreVaultEnabled: userStoreVaultEnabled,
		vaultStoreUsed:        vaultStoreUsed,
		vaultStoreKey:         vaultStoreKey,
	}
}

// ============================================================================
// Implementation (Getters/Setters)
// ============================================================================

// Audit Store
func (c *configImplementation) SetAuditStoreUsed(v bool) {
	c.auditStoreUsed = v
}

func (c *configImplementation) GetAuditStoreUsed() bool {
	return c.auditStoreUsed
}

// Blog Store
func (c *configImplementation) SetBlogStoreUsed(v bool) {
	c.blogStoreUsed = v
}

func (c *configImplementation) GetBlogStoreUsed() bool {
	return c.blogStoreUsed
}

// Cache Store
func (c *configImplementation) SetCacheStoreUsed(v bool) {
	c.cacheStoreUsed = v
}

func (c *configImplementation) GetCacheStoreUsed() bool {
	return c.cacheStoreUsed
}

// Chat Store
func (c *configImplementation) SetChatStoreUsed(v bool) {
	c.chatStoreUsed = v
}

func (c *configImplementation) GetChatStoreUsed() bool {
	return c.chatStoreUsed
}

// CMS Store
func (c *configImplementation) SetCmsStoreUsed(v bool) {
	c.cmsStoreUsed = v
}

func (c *configImplementation) GetCmsStoreUsed() bool {
	return c.cmsStoreUsed
}

func (c *configImplementation) SetCmsStoreTemplateID(v string) {
	c.cmsStoreTemplateID = v
}

func (c *configImplementation) GetCmsStoreTemplateID() string {
	return c.cmsStoreTemplateID
}

// Custom Store
func (c *configImplementation) SetCustomStoreUsed(v bool) {
	c.customStoreUsed = v
}

func (c *configImplementation) GetCustomStoreUsed() bool {
	return c.customStoreUsed
}

// Entity Store
func (c *configImplementation) SetEntityStoreUsed(v bool) {
	c.entityStoreUsed = v
}

func (c *configImplementation) GetEntityStoreUsed() bool {
	return c.entityStoreUsed
}

// Feed Store
func (c *configImplementation) SetFeedStoreUsed(v bool) {
	c.feedStoreUsed = v
}

func (c *configImplementation) GetFeedStoreUsed() bool {
	return c.feedStoreUsed
}

// Geo Store
func (c *configImplementation) SetGeoStoreUsed(v bool) {
	c.geoStoreUsed = v
}

func (c *configImplementation) GetGeoStoreUsed() bool {
	return c.geoStoreUsed
}

// Log Store
func (c *configImplementation) SetLogStoreUsed(v bool) {
	c.logStoreUsed = v
}

func (c *configImplementation) GetLogStoreUsed() bool {
	return c.logStoreUsed
}

// Meta Store
func (c *configImplementation) SetMetaStoreUsed(v bool) {
	c.metaStoreUsed = v
}

func (c *configImplementation) GetMetaStoreUsed() bool {
	return c.metaStoreUsed
}

// Session Store
func (c *configImplementation) SetSessionStoreUsed(v bool) {
	c.sessionStoreUsed = v
}

func (c *configImplementation) GetSessionStoreUsed() bool {
	return c.sessionStoreUsed
}

// Setting Store
func (c *configImplementation) SetSettingStoreUsed(v bool) {
	c.settingStoreUsed = v
}

func (c *configImplementation) GetSettingStoreUsed() bool {
	return c.settingStoreUsed
}

// Shop Store
func (c *configImplementation) SetShopStoreUsed(v bool) {
	c.shopStoreUsed = v
}

func (c *configImplementation) GetShopStoreUsed() bool {
	return c.shopStoreUsed
}

// SQL File Store
func (c *configImplementation) SetSqlFileStoreUsed(v bool) {
	c.sqlFileStoreUsed = v
}

func (c *configImplementation) GetSqlFileStoreUsed() bool {
	return c.sqlFileStoreUsed
}

// Stats Store
func (c *configImplementation) SetStatsStoreUsed(v bool) {
	c.statsStoreUsed = v
}

func (c *configImplementation) GetStatsStoreUsed() bool {
	return c.statsStoreUsed
}

// Subscription Store
func (c *configImplementation) SetSubscriptionStoreUsed(v bool) {
	c.subscriptionStoreUsed = v
}

func (c *configImplementation) GetSubscriptionStoreUsed() bool {
	return c.subscriptionStoreUsed
}

// Task Store
func (c *configImplementation) SetTaskStoreUsed(v bool) {
	c.taskStoreUsed = v
}

func (c *configImplementation) GetTaskStoreUsed() bool {
	return c.taskStoreUsed
}

// User Store
func (c *configImplementation) SetUserStoreUsed(v bool) {
	c.userStoreUsed = v
}

func (c *configImplementation) GetUserStoreUsed() bool {
	return c.userStoreUsed
}

func (c *configImplementation) SetUserStoreVaultEnabled(v bool) {
	c.userStoreVaultEnabled = v
}

func (c *configImplementation) GetUserStoreVaultEnabled() bool {
	return c.userStoreVaultEnabled
}

// Vault Store
func (c *configImplementation) SetVaultStoreUsed(v bool) {
	c.vaultStoreUsed = v
}

func (c *configImplementation) GetVaultStoreUsed() bool {
	return c.vaultStoreUsed
}

func (c *configImplementation) SetVaultStoreKey(v string) {
	c.vaultStoreKey = v
}

func (c *configImplementation) GetVaultStoreKey() string {
	return c.vaultStoreKey
}
