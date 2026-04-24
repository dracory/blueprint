// The z_ prefix keeps this file sorted after user-configurable files in the directory listing.

package config

// ============================================================================
// Main Config Interface
// ============================================================================

// ConfigInterface defines the contract for application configuration.
// It composes all domain-specific configuration interfaces.
type ConfigInterface interface {
	AppConfigInterface
	AuthConfigInterface
	DatabaseConfigInterface
	EmailConfigInterface
	EncryptionConfigInterface
	I18nConfigInterface
	LLMConfigInterface
	MediaConfigInterface
	PaymentConfigInterface
	SEOConfigInterface

	// CMS MCP
	SetCmsMcpApiKey(string)
	GetCmsMcpApiKey() string

	// Datastores
	AuditStoreConfigInterface
	BlogStoreConfigInterface
	CacheStoreConfigInterface
	ChatStoreConfigInterface
	CmsStoreConfigInterface
	CustomStoreConfigInterface
	EntityStoreConfigInterface
	FeedStoreConfigInterface
	GeoStoreConfigInterface
	LogStoreConfigInterface
	MetaStoreConfigInterface
	SessionStoreConfigInterface
	SettingStoreConfigInterface
	ShopStoreConfigInterface
	SqlFileStoreConfigInterface
	StatsStoreConfigInterface
	SubscriptionStoreConfigInterface
	TaskStoreConfigInterface
	UserStoreConfigInterface
	VaultStoreConfigInterface
}

// ============================================================================
// App Config Interface
// ============================================================================

// AppConfigInterface defines application-level configuration methods.
type AppConfigInterface interface {
	SetAppName(string)
	GetAppName() string

	SetAppType(string)
	GetAppType() string

	SetAppEnv(string)
	GetAppEnv() string

	SetAppHost(string)
	GetAppHost() string

	SetAppPort(string)
	GetAppPort() string

	SetAppUrl(string)
	GetAppUrl() string

	SetAppDebug(bool)
	GetAppDebug() bool

	// Environment helpers
	IsEnvDevelopment() bool
	IsEnvLocal() bool
	IsEnvProduction() bool
	IsEnvStaging() bool
	IsEnvTesting() bool
}

// ============================================================================
// Auth Config Interface
// ============================================================================

// AuthConfigInterface defines authentication configuration methods.
type AuthConfigInterface interface {
	SetRegistrationEnabled(bool)
	GetRegistrationEnabled() bool

	SetEmailsAllowedAccess([]string)
	GetEmailsAllowedAccess() []string
}

// ============================================================================
// Database Config Interface
// ============================================================================

// DatabaseConfigInterface defines database configuration methods.
type DatabaseConfigInterface interface {
	SetDatabaseDriver(string)
	GetDatabaseDriver() string

	SetDatabaseHost(string)
	GetDatabaseHost() string

	SetDatabasePort(string)
	GetDatabasePort() string

	SetDatabaseName(string)
	GetDatabaseName() string

	SetDatabaseUsername(string)
	GetDatabaseUsername() string

	SetDatabasePassword(string)
	GetDatabasePassword() string

	SetDatabaseSSLMode(string)
	GetDatabaseSSLMode() string

	SetDatabaseMaxOpenConns(int)
	GetDatabaseMaxOpenConns() int

	SetDatabaseMaxIdleConns(int)
	GetDatabaseMaxIdleConns() int

	SetDatabaseConnMaxLifetimeSeconds(int)
	GetDatabaseConnMaxLifetimeSeconds() int

	SetDatabaseConnMaxIdleTimeSeconds(int)
	GetDatabaseConnMaxIdleTimeSeconds() int

	SetDatabaseCharset(string)
	GetDatabaseCharset() string

	SetDatabaseTimezone(string)
	GetDatabaseTimezone() string
}

// ============================================================================
// Email Config Interface
// ============================================================================

// EmailConfigInterface defines email/mail configuration methods.
type EmailConfigInterface interface {
	SetMailDriver(string)
	GetMailDriver() string

	SetMailHost(string)
	GetMailHost() string

	SetMailPort(int)
	GetMailPort() int

	SetMailUsername(string)
	GetMailUsername() string

	SetMailPassword(string)
	GetMailPassword() string

	SetMailFromAddress(string)
	GetMailFromAddress() string

	SetMailFromName(string)
	GetMailFromName() string
}

// ============================================================================
// Encryption Config Interface
// ============================================================================

// EncryptionConfigInterface defines encryption configuration methods.
type EncryptionConfigInterface interface {
	SetEnvEncryptionKey(string)
	GetEnvEncryptionKey() string
	SetEnvEncUsed(bool)
	GetEnvEncUsed() bool
	SetEnvEncPublicKey(string)
	GetEnvEncPublicKey() string
}

// ============================================================================
// i18n Config Interface
// ============================================================================

// I18nConfigInterface defines internationalization configuration methods.
type I18nConfigInterface interface {
	SetTranslationLanguageDefault(string)
	GetTranslationLanguageDefault() string

	SetTranslationLanguageList(map[string]string)
	GetTranslationLanguageList() map[string]string
}

// ============================================================================
// LLM Config Interface
// ============================================================================

// LLMConfigInterface defines LLM provider configuration methods.
type LLMConfigInterface interface {
	// Anthropic
	SetAnthropicApiUsed(bool)
	GetAnthropicApiUsed() bool
	SetAnthropicApiKey(string)
	GetAnthropicApiKey() string
	SetAnthropicApiDefaultModel(string)
	GetAnthropicApiDefaultModel() string

	// Google Gemini
	SetGoogleGeminiApiUsed(bool)
	GetGoogleGeminiApiUsed() bool
	SetGoogleGeminiApiKey(string)
	GetGoogleGeminiApiKey() string
	SetGoogleGeminiApiDefaultModel(string)
	GetGoogleGeminiApiDefaultModel() string

	// OpenRouter
	SetOpenRouterApiUsed(bool)
	GetOpenRouterApiUsed() bool
	SetOpenRouterApiKey(string)
	GetOpenRouterApiKey() string
	SetOpenRouterApiDefaultModel(string)
	GetOpenRouterApiDefaultModel() string

	// OpenAI
	SetOpenAiApiUsed(bool)
	GetOpenAiApiUsed() bool
	SetOpenAiApiKey(string)
	GetOpenAiApiKey() string
	SetOpenAiApiDefaultModel(string)
	GetOpenAiApiDefaultModel() string

	// Vertex AI
	SetVertexAiApiUsed(bool)
	GetVertexAiApiUsed() bool
	SetVertexAiApiDefaultModel(string)
	GetVertexAiApiDefaultModel() string
	SetVertexAiApiProjectID(string)
	GetVertexAiApiProjectID() string
	SetVertexAiApiRegionID(string)
	GetVertexAiApiRegionID() string
	SetVertexAiApiModelID(string)
	GetVertexAiApiModelID() string
}

// ============================================================================
// Media Config Interface
// ============================================================================

// MediaConfigInterface defines media storage configuration methods.
type MediaConfigInterface interface {
	SetMediaBucket(string)
	GetMediaBucket() string

	SetMediaDriver(string)
	GetMediaDriver() string

	SetMediaKey(string)
	GetMediaKey() string

	SetMediaEndpoint(string)
	GetMediaEndpoint() string

	SetMediaRegion(string)
	GetMediaRegion() string

	SetMediaRoot(string)
	GetMediaRoot() string

	SetMediaSecret(string)
	GetMediaSecret() string

	SetMediaUrl(string)
	GetMediaUrl() string
}

// ============================================================================
// Payment Config Interface
// ============================================================================

// PaymentConfigInterface defines payment provider configuration methods.
type PaymentConfigInterface interface {
	SetStripeKeyPrivate(string)
	GetStripeKeyPrivate() string

	SetStripeKeyPublic(string)
	GetStripeKeyPublic() string

	SetStripeUsed(bool)
	GetStripeUsed() bool
}

// ============================================================================
// SEO Config Interface
// ============================================================================

// SEOConfigInterface defines SEO configuration methods.
type SEOConfigInterface interface {
	SetIndexNowKey(string)
	GetIndexNowKey() string
}

// ============================================================================
// Store Config Interfaces
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
