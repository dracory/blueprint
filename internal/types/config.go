package types

var _ ConfigInterface = (*Config)(nil)

type Config struct {
	// App configuration
	appName  string
	appType  string
	appEnv   string
	appHost  string
	appPort  string
	appUrl   string
	appDebug bool

	// Email configuration
	emailDriver      string
	emailHost        string
	emailPort        int
	emailUsername    string
	emailPassword    string
	emailFromName    string
	emailFromAddress string

	// Database configuration
	databaseDriver   string
	databaseHost     string
	databasePort     string
	databaseName     string
	databaseUsername string
	databasePassword string
	databaseSSLMode  string

	// LLM configuration
	openRouterApiKey          string
	openRouterApiUsed         bool
	openRouterApiDefaultModel string

	// OpenAI
	openAiApiKey          string
	openAiApiUsed         bool
	openAiApiDefaultModel string

	// Anthropic
	anthropicApiUsed         bool
	anthropicApiKey          string
	anthropicApiDefaultModel string

	// Google Gemini
	googleGeminiApiUsed         bool
	googleGeminiApiKey          string
	googleGeminiApiDefaultModel string

	// Vertex AI
	vertexAiApiUsed         bool
	vertexAiApiDefaultModel string
	vertexAiApiProjectID    string
	vertexAiApiRegionID     string
	vertexAiApiModelID      string

	// Encryption
	envEncryptionKey string

	// Store flags
	auditStoreUsed        bool
	blogStoreUsed      bool
	cacheStoreUsed     bool
	cmsStoreUsed       bool
	cmsStoreTemplateID string
	customStoreUsed    bool
	entityStoreUsed    bool
	feedStoreUsed      bool
	geoStoreUsed       bool
	logStoreUsed       bool
	metaStoreUsed      bool
	sessionStoreUsed   bool
	settingStoreUsed   bool
	shopStoreUsed      bool
	sqlFileStoreUsed   bool
	statsStoreUsed     bool
	subscriptionStoreUsed bool
	taskStoreUsed      bool
	userStoreUsed      bool
	userStoreVaultEnabled bool
	vaultStoreUsed     bool
	vaultStoreKey      string

	// i18n / Translation
	translationLanguageDefault string
	translationLanguageList    map[string]string

	// App-specific settings
	stripeKeyPrivate string
	stripeKeyPublic  string
	stripeUsed       bool

	// Authentication
	registrationEnabled bool

	// Media configuration
	mediaBucket   string
	mediaDriver   string
	mediaKey      string
	mediaEndpoint string
	mediaRegion   string
	mediaRoot     string
	mediaSecret   string
	mediaUrl      string
}

func (c *Config) SetAppName(appName string) {
	c.appName = appName
}

func (c *Config) GetAppName() string {
	return c.appName
}

func (c *Config) SetAppType(appType string) {
	c.appType = appType
}

func (c *Config) GetAppType() string {
	return c.appType
}

func (c *Config) SetAppEnv(appEnv string) {
	c.appEnv = appEnv
}

func (c *Config) GetAppEnv() string {
	return c.appEnv
}

func (c *Config) SetAppHost(appHost string) {
	c.appHost = appHost
}

func (c *Config) GetAppHost() string {
	return c.appHost
}

func (c *Config) SetAppPort(appPort string) {
	c.appPort = appPort
}

func (c *Config) GetAppPort() string {
	return c.appPort
}

func (c *Config) SetAppUrl(appUrl string) {
	c.appUrl = appUrl
}

func (c *Config) GetAppUrl() string {
	return c.appUrl
}

func (c *Config) SetAppDebug(appDebug bool) {
	c.appDebug = appDebug
}

func (c *Config) GetAppDebug() bool {
	return c.appDebug
}

// == Environment Helpers ==
// These methods provide convenient checks for the current app environment.
// They compare the configured environment string to known values.
func (c *Config) IsEnvDevelopment() bool {
	return c.appEnv == "development"
}

func (c *Config) IsEnvLocal() bool {
	return c.appEnv == "local"
}

func (c *Config) IsEnvProduction() bool {
	return c.appEnv == "production"
}

func (c *Config) IsEnvStaging() bool {
	return c.appEnv == "staging"
}

func (c *Config) IsEnvTesting() bool {
	return c.appEnv == "testing"
}

// == Email Getters/Setters ==
func (c *Config) SetMailDriver(v string) {
	c.emailDriver = v
}

func (c *Config) GetMailDriver() string {
	return c.emailDriver
}

func (c *Config) SetMailHost(v string) {
	c.emailHost = v
}

func (c *Config) GetMailHost() string {
	return c.emailHost
}

func (c *Config) SetMailPort(v int) {
	c.emailPort = v
}

func (c *Config) GetMailPort() int {
	return c.emailPort
}

func (c *Config) SetMailUsername(v string) {
	c.emailUsername = v
}

func (c *Config) GetMailUsername() string {
	return c.emailUsername
}

func (c *Config) SetMailPassword(v string) {
	c.emailPassword = v
}

func (c *Config) GetMailPassword() string {
	return c.emailPassword
}

func (c *Config) SetMailFromName(v string) {
	c.emailFromName = v
}

func (c *Config) GetMailFromName() string {
	return c.emailFromName
}

// == Email From Address Getters/Setters ==
func (c *Config) SetMailFromAddress(v string) {
	c.emailFromAddress = v
}

func (c *Config) GetMailFromAddress() string {
	return c.emailFromAddress
}

// == Database Getters/Setters ==
func (c *Config) SetDatabaseDriver(v string) {
	c.databaseDriver = v
}

func (c *Config) GetDatabaseDriver() string {
	return c.databaseDriver
}

func (c *Config) SetDatabaseHost(v string) {
	c.databaseHost = v
}

func (c *Config) GetDatabaseHost() string {
	return c.databaseHost
}

func (c *Config) SetDatabasePort(v string) {
	c.databasePort = v
}

func (c *Config) GetDatabasePort() string {
	return c.databasePort
}

func (c *Config) SetDatabaseName(v string) {
	c.databaseName = v
}

func (c *Config) GetDatabaseName() string {
	return c.databaseName
}

func (c *Config) SetDatabaseUsername(v string) {
	c.databaseUsername = v
}

func (c *Config) GetDatabaseUsername() string {
	return c.databaseUsername
}

func (c *Config) SetDatabasePassword(v string) {
	c.databasePassword = v
}

func (c *Config) GetDatabasePassword() string {
	return c.databasePassword
}

func (c *Config) SetDatabaseSSLMode(v string) {
	c.databaseSSLMode = v
}

func (c *Config) GetDatabaseSSLMode() string {
	return c.databaseSSLMode
}

// == LLM Getters/Setters ==

// LLM: Anthropic
func (c *Config) SetAnthropicApiUsed(v bool) {
	c.anthropicApiUsed = v
}
func (c *Config) GetAnthropicApiUsed() bool {
	return c.anthropicApiUsed
}

func (c *Config) SetAnthropicApiKey(v string) {
	c.anthropicApiKey = v
}
func (c *Config) GetAnthropicApiKey() string {
	return c.anthropicApiKey
}

func (c *Config) SetAnthropicApiDefaultModel(v string) {
	c.anthropicApiDefaultModel = v
}
func (c *Config) GetAnthropicApiDefaultModel() string {
	return c.anthropicApiDefaultModel
}

// LLM: Google Gemini
func (c *Config) SetGoogleGeminiApiUsed(v bool) {
	c.googleGeminiApiUsed = v
}
func (c *Config) GetGoogleGeminiApiUsed() bool {
	return c.googleGeminiApiUsed
}

func (c *Config) SetGoogleGeminiApiKey(v string) {
	c.googleGeminiApiKey = v
}
func (c *Config) GetGoogleGeminiApiKey() string {
	return c.googleGeminiApiKey
}

func (c *Config) SetGoogleGeminiApiDefaultModel(v string) {
	c.googleGeminiApiDefaultModel = v
}
func (c *Config) GetGoogleGeminiApiDefaultModel() string {
	return c.googleGeminiApiDefaultModel
}

// LLM: OpenAI (mapped to existing openAIKey field for key storage)
func (c *Config) SetOpenAiApiUsed(v bool) {
	c.openAiApiUsed = v
}

func (c *Config) GetOpenAiApiUsed() bool {
	return c.openAiApiUsed
}

func (c *Config) SetOpenAiApiKey(v string) {
	c.openAiApiKey = v
}
func (c *Config) GetOpenAiApiKey() string {
	return c.openAiApiKey
}

func (c *Config) SetOpenAiApiDefaultModel(v string) {
	c.openAiApiDefaultModel = v
}
func (c *Config) GetOpenAiApiDefaultModel() string {
	return c.openAiApiDefaultModel
}

func (c *Config) SetOpenRouterApiKey(v string) {
	c.openRouterApiKey = v
}

func (c *Config) GetOpenRouterApiKey() string {
	return c.openRouterApiKey
}

func (c *Config) SetOpenRouterApiUsed(v bool) {
	c.openRouterApiUsed = v
}

func (c *Config) GetOpenRouterApiUsed() bool {
	return c.openRouterApiUsed
}

func (c *Config) SetOpenRouterApiDefaultModel(v string) {
	c.openRouterApiDefaultModel = v
}

func (c *Config) GetOpenRouterApiDefaultModel() string {
	return c.openRouterApiDefaultModel
}

// LLM: Vertex AI
func (c *Config) SetVertexAiApiUsed(v bool) {
	c.vertexAiApiUsed = v
}
func (c *Config) GetVertexAiApiUsed() bool {
	return c.vertexAiApiUsed
}

func (c *Config) SetVertexAiApiDefaultModel(v string) {
	c.vertexAiApiDefaultModel = v
}
func (c *Config) GetVertexAiApiDefaultModel() string {
	return c.vertexAiApiDefaultModel
}

func (c *Config) SetVertexAiApiProjectID(v string) {
	c.vertexAiApiProjectID = v
}
func (c *Config) GetVertexAiApiProjectID() string {
	return c.vertexAiApiProjectID
}

func (c *Config) SetVertexAiApiRegionID(v string) {
	c.vertexAiApiRegionID = v
}
func (c *Config) GetVertexAiApiRegionID() string { return c.vertexAiApiRegionID }

func (c *Config) SetVertexAiApiModelID(v string) { c.vertexAiApiModelID = v }
func (c *Config) GetVertexAiApiModelID() string  { return c.vertexAiApiModelID }

// == Encryption Getters/Setters ==
func (c *Config) SetEnvEncryptionKey(v string) {
	c.envEncryptionKey = v
}

func (c *Config) GetEnvEncryptionKey() string {
	return c.envEncryptionKey
}

// == Cache Store Getters/Setters ==
func (c *Config) SetCacheStoreUsed(v bool) {
	c.cacheStoreUsed = v
}

func (c *Config) GetCacheStoreUsed() bool {
	return c.cacheStoreUsed
}

// == Subscription Store Getters/Setters ==
func (c *Config) SetSubscriptionStoreUsed(v bool) {
	c.subscriptionStoreUsed = v
}

func (c *Config) GetSubscriptionStoreUsed() bool {
	return c.subscriptionStoreUsed
}

// == Blog Store Getters/Setters ==
func (c *Config) SetBlogStoreUsed(v bool) {
	c.blogStoreUsed = v
}

func (c *Config) GetBlogStoreUsed() bool {
	return c.blogStoreUsed
}

// == CMS Store Getters/Setters ==
func (c *Config) SetCmsStoreUsed(v bool) {
	c.cmsStoreUsed = v
}

func (c *Config) GetCmsStoreUsed() bool {
	return c.cmsStoreUsed
}

func (c *Config) SetCmsStoreTemplateID(v string) {
	c.cmsStoreTemplateID = v
}

func (c *Config) GetCmsStoreTemplateID() string {
	return c.cmsStoreTemplateID
}

// == Custom Store Getters/Setters ==
func (c *Config) SetCustomStoreUsed(v bool) {
	c.customStoreUsed = v
}

func (c *Config) GetCustomStoreUsed() bool {
	return c.customStoreUsed
}

// == Entity Store Getters/Setters ==
func (c *Config) SetEntityStoreUsed(v bool) {
	c.entityStoreUsed = v
}

func (c *Config) GetEntityStoreUsed() bool {
	return c.entityStoreUsed
}

// == Feed Store Getters/Setters ==
func (c *Config) SetFeedStoreUsed(v bool) {
	c.feedStoreUsed = v
}

func (c *Config) GetFeedStoreUsed() bool {
	return c.feedStoreUsed
}

// == Geo Store Getters/Setters ==
func (c *Config) SetGeoStoreUsed(v bool) {
	c.geoStoreUsed = v
}

func (c *Config) GetGeoStoreUsed() bool {
	return c.geoStoreUsed
}

// == Log Store Getters/Setters ==
func (c *Config) SetLogStoreUsed(v bool) {
	c.logStoreUsed = v
}

func (c *Config) GetLogStoreUsed() bool {
	return c.logStoreUsed
}

// == Meta Store Getters/Setters ==
func (c *Config) SetMetaStoreUsed(v bool) {
	c.metaStoreUsed = v
}

func (c *Config) GetMetaStoreUsed() bool {
	return c.metaStoreUsed
}

// == Session Store Getters/Setters ==
func (c *Config) SetSessionStoreUsed(v bool) {
	c.sessionStoreUsed = v
}

func (c *Config) GetSessionStoreUsed() bool {
	return c.sessionStoreUsed
}

// == Sql File Store Getters/Setters ==
func (c *Config) SetSqlFileStoreUsed(v bool) {
	c.sqlFileStoreUsed = v
}

func (c *Config) GetSqlFileStoreUsed() bool {
	return c.sqlFileStoreUsed
}

// == Setting Store Getters/Setters ==
func (c *Config) SetSettingStoreUsed(v bool) {
	c.settingStoreUsed = v
}

func (c *Config) GetSettingStoreUsed() bool {
	return c.settingStoreUsed
}

// == Shop Store Getters/Setters ==
func (c *Config) SetShopStoreUsed(v bool) {
	c.shopStoreUsed = v
}

func (c *Config) GetShopStoreUsed() bool {
	return c.shopStoreUsed
}

// == Task Store Getters/Setters ==
func (c *Config) SetTaskStoreUsed(v bool) {
	c.taskStoreUsed = v
}

func (c *Config) GetTaskStoreUsed() bool {
	return c.taskStoreUsed
}

// == User Store Getters/Setters ==
func (c *Config) SetUserStoreUsed(v bool) {
	c.userStoreUsed = v
}

func (c *Config) GetUserStoreUsed() bool {
	return c.userStoreUsed
}

func (c *Config) SetUserStoreVaultEnabled(v bool) {
	c.userStoreVaultEnabled = v
}

func (c *Config) GetUserStoreVaultEnabled() bool {
	return c.userStoreVaultEnabled
}

// == Stats Getters/Setters ==
func (c *Config) SetStatsStoreUsed(v bool) {
	c.statsStoreUsed = v
}

func (c *Config) GetStatsStoreUsed() bool {
	return c.statsStoreUsed
}

// == i18n Getters/Setters ==
func (c *Config) SetTranslationLanguageDefault(v string) {
	c.translationLanguageDefault = v
}

func (c *Config) GetTranslationLanguageDefault() string {
	return c.translationLanguageDefault
}

func (c *Config) SetTranslationLanguageList(v map[string]string) {
	c.translationLanguageList = v
}

func (c *Config) GetTranslationLanguageList() map[string]string {
	return c.translationLanguageList
}

// == App-specific Getters/Setters ==
func (c *Config) SetStripeKeyPrivate(v string) {
	c.stripeKeyPrivate = v
}

func (c *Config) GetStripeKeyPrivate() string {
	return c.stripeKeyPrivate
}

func (c *Config) SetStripeKeyPublic(v string) {
	c.stripeKeyPublic = v
}
func (c *Config) GetStripeKeyPublic() string {
	return c.stripeKeyPublic
}

func (c *Config) SetStripeUsed(v bool) {
	c.stripeUsed = v
}
func (c *Config) GetStripeUsed() bool { return c.stripeUsed }

func (c *Config) SetRegistrationEnabled(v bool) {
	c.registrationEnabled = v
}

func (c *Config) GetRegistrationEnabled() bool {
	return c.registrationEnabled
}

// == Audit Store Getters/Setters ==
func (c *Config) SetAuditStoreUsed(v bool) {
	c.auditStoreUsed = v
}

func (c *Config) GetAuditStoreUsed() bool {
	return c.auditStoreUsed
}

// == Vault Store Getters/Setters ==
func (c *Config) SetVaultStoreUsed(v bool) {
	c.vaultStoreUsed = v
}

func (c *Config) GetVaultStoreUsed() bool {
	return c.vaultStoreUsed
}

func (c *Config) SetVaultStoreKey(v string) {
	c.vaultStoreKey = v
}
func (c *Config) GetVaultStoreKey() string {
	return c.vaultStoreKey
}

// == Media Getters/Setters ==
func (c *Config) SetMediaBucket(v string)   { c.mediaBucket = v }
func (c *Config) GetMediaBucket() string    { return c.mediaBucket }
func (c *Config) SetMediaDriver(v string)   { c.mediaDriver = v }
func (c *Config) GetMediaDriver() string    { return c.mediaDriver }
func (c *Config) SetMediaKey(v string)      { c.mediaKey = v }
func (c *Config) GetMediaKey() string       { return c.mediaKey }
func (c *Config) SetMediaEndpoint(v string) { c.mediaEndpoint = v }
func (c *Config) GetMediaEndpoint() string  { return c.mediaEndpoint }
func (c *Config) SetMediaRegion(v string)   { c.mediaRegion = v }
func (c *Config) GetMediaRegion() string    { return c.mediaRegion }
func (c *Config) SetMediaRoot(v string)     { c.mediaRoot = v }
func (c *Config) GetMediaRoot() string      { return c.mediaRoot }
func (c *Config) SetMediaSecret(v string)   { c.mediaSecret = v }
func (c *Config) GetMediaSecret() string    { return c.mediaSecret }
func (c *Config) SetMediaUrl(v string)      { c.mediaUrl = v }
func (c *Config) GetMediaUrl() string       { return c.mediaUrl }
