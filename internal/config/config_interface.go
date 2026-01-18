package config

type appConfigInterface interface {
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

type databaseConfigInterface interface {
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
}

type emailConfigInterface interface {
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

type authConfigInterface interface {
	SetRegistrationEnabled(bool)
	GetRegistrationEnabled() bool
}

type llmConfigInterface interface {
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

type envEncryptionConfigInterface interface {
	SetEnvEncryptionKey(string)
	GetEnvEncryptionKey() string
}

type blogStoreConfigInterface interface {
	SetBlogStoreUsed(bool)
	GetBlogStoreUsed() bool
}

type chatStoreConfigInterface interface {
	SetChatStoreUsed(bool)
	GetChatStoreUsed() bool
}

type cacheStoreConfigInterface interface {
	SetCacheStoreUsed(bool)
	GetCacheStoreUsed() bool
}

type cmsStoreConfigInterface interface {
	SetCmsStoreUsed(bool)
	GetCmsStoreUsed() bool
	SetCmsStoreTemplateID(string)
	GetCmsStoreTemplateID() string
}

type customStoreConfigInterface interface {
	SetCustomStoreUsed(bool)
	GetCustomStoreUsed() bool
}

type entityStoreConfigInterface interface {
	SetEntityStoreUsed(bool)
	GetEntityStoreUsed() bool
}

type feedStoreConfigInterface interface {
	SetFeedStoreUsed(bool)
	GetFeedStoreUsed() bool
}

type geoStoreConfigInterface interface {
	SetGeoStoreUsed(bool)
	GetGeoStoreUsed() bool
}

type logStoreConfigInterface interface {
	SetLogStoreUsed(bool)
	GetLogStoreUsed() bool
}

type metaStoreConfigInterface interface {
	SetMetaStoreUsed(bool)
	GetMetaStoreUsed() bool
}

type subscriptionStoreConfigInterface interface {
	SetSubscriptionStoreUsed(bool)
	GetSubscriptionStoreUsed() bool
}

type sessionStoreConfigInterface interface {
	SetSessionStoreUsed(bool)
	GetSessionStoreUsed() bool
}

type settingStoreConfigInterface interface {
	SetSettingStoreUsed(bool)
	GetSettingStoreUsed() bool
}

type shopStoreConfigInterface interface {
	SetShopStoreUsed(bool)
	GetShopStoreUsed() bool
}

type sqlFileStoreConfigInterface interface {
	SetSqlFileStoreUsed(bool)
	GetSqlFileStoreUsed() bool
}

type statsStoreConfigInterface interface {
	SetStatsStoreUsed(bool)
	GetStatsStoreUsed() bool
}

type taskStoreConfigInterface interface {
	SetTaskStoreUsed(bool)
	GetTaskStoreUsed() bool
}

type userStoreConfigInterface interface {
	SetUserStoreUsed(bool)
	GetUserStoreUsed() bool
	SetUserStoreVaultEnabled(bool)
	GetUserStoreVaultEnabled() bool
}

type auditStoreConfigInterface interface {
	SetAuditStoreUsed(bool)
	GetAuditStoreUsed() bool
}

type vaultStoreConfigInterface interface {
	SetVaultStoreUsed(bool)
	GetVaultStoreUsed() bool
	SetVaultStoreKey(string)
	GetVaultStoreKey() string
}

type i18nConfigInterface interface {
	SetTranslationLanguageDefault(string)
	GetTranslationLanguageDefault() string
	SetTranslationLanguageList(map[string]string)
	GetTranslationLanguageList() map[string]string
}

type paymentConfigInterface interface {
	SetStripeKeyPrivate(string)
	GetStripeKeyPrivate() string
	SetStripeKeyPublic(string)
	GetStripeKeyPublic() string
	SetStripeUsed(bool)
	GetStripeUsed() bool
}

type mediaConfigInterface interface {
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

type ConfigInterface interface {
	// App-specific settings
	appConfigInterface
	authConfigInterface
	databaseConfigInterface
	emailConfigInterface
	envEncryptionConfigInterface
	i18nConfigInterface
	llmConfigInterface
	mediaConfigInterface
	paymentConfigInterface

	// Datastores
	auditStoreConfigInterface
	blogStoreConfigInterface
	chatStoreConfigInterface
	cacheStoreConfigInterface
	cmsStoreConfigInterface
	customStoreConfigInterface
	entityStoreConfigInterface
	feedStoreConfigInterface
	geoStoreConfigInterface
	logStoreConfigInterface
	metaStoreConfigInterface
	subscriptionStoreConfigInterface
	sessionStoreConfigInterface
	settingStoreConfigInterface
	shopStoreConfigInterface
	sqlFileStoreConfigInterface
	statsStoreConfigInterface
	taskStoreConfigInterface
	userStoreConfigInterface
	vaultStoreConfigInterface
}

