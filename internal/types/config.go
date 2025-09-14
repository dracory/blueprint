package types

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
	SetMailFromName(string)
	GetMailFromName() string
	SetMailFromEmail(string)
	GetMailFromEmail() string
}

type llmConfigInterface interface {
	// Anthropic
	SetAnthropicApiUsed(bool)
	GetAnthropicApiUsed() bool
	SetAnthropicApiKey(string)
	GetAnthropicApiKey() string
	SetAnthropicDefaultModel(string)
	GetAnthropicDefaultModel() string

	// Google Gemini
	SetGoogleGeminiApiUsed(bool)
	GetGoogleGeminiApiUsed() bool
	SetGoogleGeminiApiKey(string)
	GetGoogleGeminiApiKey() string
	SetGoogleGeminiDefaultModel(string)
	GetGoogleGeminiDefaultModel() string

	// OpenRouter
	SetOpenRouterApiUsed(bool)
	GetOpenRouterApiUsed() bool
	SetOpenRouterApiKey(string)
	GetOpenRouterApiKey() string
	SetOpenRouterDefaultModel(string)
	GetOpenRouterDefaultModel() string

	// OpenAI
	SetOpenAiApiUsed(bool)
	GetOpenAiApiUsed() bool
	SetOpenAiApiKey(string)
	GetOpenAiApiKey() string
	SetOpenAiDefaultModel(string)
	GetOpenAiDefaultModel() string

	// Vertex AI
	SetVertexAiUsed(bool)
	GetVertexAiUsed() bool
	SetVertexAiDefaultModel(string)
	GetVertexAiDefaultModel() string
	SetVertexAiProjectID(string)
	GetVertexAiProjectID() string
	SetVertexAiRegionID(string)
	GetVertexAiRegionID() string
	SetVertexAiModelID(string)
	GetVertexAiModelID() string
}

type envEncryptionConfigInterface interface {
	SetEnvEncryptionKey(string)
	GetEnvEncryptionKey() string
}

type cmsStoreConfigInterface interface {
	SetCmsStoreUsed(bool)
	GetCmsStoreUsed() bool
	SetCmsStoreTemplateID(string)
	GetCmsStoreTemplateID() string
}

type statsStoreConfigInterface interface {
	SetStatsStoreUsed(bool)
	GetStatsStoreUsed() bool
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

type appSpecificConfigInterface interface {
	// SetCMSTemplateID(string)
	// GetCMSTemplateID() string

	// SetVaultKey(string)
	// GetVaultKey() string

	SetOpenAIKey(string)
	GetOpenAIKey() string

	SetStripeKeyPrivate(string)
	GetStripeKeyPrivate() string

	SetStripeKeyPublic(string)
	GetStripeKeyPublic() string

	SetStripeUsed(bool)
	GetStripeUsed() bool

	SetVertexProjectID(string)
	GetVertexProjectID() string

	SetVertexRegionID(string)
	GetVertexRegionID() string

	SetVertexModelID(string)
	GetVertexModelID() string
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
	appConfigInterface
	emailConfigInterface
	databaseConfigInterface
	llmConfigInterface
	envEncryptionConfigInterface
	cmsStoreConfigInterface
	statsStoreConfigInterface
	mediaConfigInterface
	vaultStoreConfigInterface
	i18nConfigInterface
	appSpecificConfigInterface
}

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
	emailDriver    string
	emailHost      string
	emailPort      int
	emailUsername  string
	emailPassword  string
	emailFromName  string
	emailFromEmail string

	// Database configuration
	databaseDriver   string
	databaseHost     string
	databasePort     string
	databaseName     string
	databaseUsername string
	databasePassword string
	databaseSSLMode  string

	// LLM configuration
	openRouterApiKey       string
	openRouterApiUsed      bool
	openRouterDefaultModel string

	// OpenAI
	openAiApiUsed      bool
	openAiDefaultModel string

	// Anthropic
	anthropicApiUsed      bool
	anthropicApiKey       string
	anthropicDefaultModel string

	// Google Gemini
	googleGeminiApiUsed      bool
	googleGeminiApiKey       string
	googleGeminiDefaultModel string

	// Encryption
	envEncryptionKey string

	// CMS configuration
	cmsStoreUsed       bool
	cmsStoreTemplateID string

	// i18n / Translation
	translationLanguageDefault string
	translationLanguageList    map[string]string

	// App-specific settings
	cmsTemplateID        string
	openAIKey            string
	stripeKeyPrivate     string
	stripeKeyPublic      string
	stripeUsed           bool
	vertexProjectID      string
	vertexRegionID       string
	vertexModelID        string
	vertexAiUsed         bool
	vertexAiDefaultModel string

	// Stats configuration
	statsStoreUsed bool

	// Vault configuration
	vaultStoreUsed bool
	vaultStoreKey  string

	// Media configuration
	mediaBucket   string
	mediaDriver   string
	mediaKey      string
	mediaEndpoint string
	mediaRegion   string
	mediaRoot     string
	mediaSecret   string
	mediaUrl      string

	// Daily analysis configuration
	dailyAnalysisSymbols      []string
	dailyAnalysisTimeUTC      string
	dailyAnalysisCadenceHours int
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

func (c *Config) SetMailFromEmail(v string) {
	c.emailFromEmail = v
}

func (c *Config) GetMailFromEmail() string {
	return c.emailFromEmail
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

func (c *Config) SetOpenRouterDefaultModel(v string) {
	c.openRouterDefaultModel = v
}

func (c *Config) GetOpenRouterDefaultModel() string {
	return c.openRouterDefaultModel
}

// OpenAI (mapped to existing openAIKey field for key storage)
func (c *Config) SetOpenAiApiUsed(v bool) {
	c.openAiApiUsed = v
}

func (c *Config) GetOpenAiApiUsed() bool {
	return c.openAiApiUsed
}

func (c *Config) SetOpenAiApiKey(v string) {
	c.openAIKey = v
}
func (c *Config) GetOpenAiApiKey() string {
	return c.openAIKey
}

func (c *Config) SetOpenAiDefaultModel(v string) {
	c.openAiDefaultModel = v
}
func (c *Config) GetOpenAiDefaultModel() string {
	return c.openAiDefaultModel
}

// Anthropic
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

func (c *Config) SetAnthropicDefaultModel(v string) {
	c.anthropicDefaultModel = v
}
func (c *Config) GetAnthropicDefaultModel() string {
	return c.anthropicDefaultModel
}

// Google Gemini
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

func (c *Config) SetGoogleGeminiDefaultModel(v string) {
	c.googleGeminiDefaultModel = v
}
func (c *Config) GetGoogleGeminiDefaultModel() string {
	return c.googleGeminiDefaultModel
}

// Vertex AI
func (c *Config) SetVertexAiUsed(v bool) {
	c.vertexAiUsed = v
}
func (c *Config) GetVertexAiUsed() bool {
	return c.vertexAiUsed
}

func (c *Config) SetVertexAiDefaultModel(v string) {
	c.vertexAiDefaultModel = v
}
func (c *Config) GetVertexAiDefaultModel() string {
	return c.vertexAiDefaultModel
}

func (c *Config) SetVertexAiProjectID(v string) {
	c.vertexProjectID = v
}
func (c *Config) GetVertexAiProjectID() string {
	return c.vertexProjectID
}

func (c *Config) SetVertexAiRegionID(v string) {
	c.vertexRegionID = v
}
func (c *Config) GetVertexAiRegionID() string { return c.vertexRegionID }

func (c *Config) SetVertexAiModelID(v string) { c.vertexModelID = v }
func (c *Config) GetVertexAiModelID() string  { return c.vertexModelID }

// == Encryption Getters/Setters ==
func (c *Config) SetEnvEncryptionKey(v string) {
	c.envEncryptionKey = v
}

func (c *Config) GetEnvEncryptionKey() string {
	return c.envEncryptionKey
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
func (c *Config) SetOpenAIKey(v string) { c.openAIKey = v }
func (c *Config) GetOpenAIKey() string  { return c.openAIKey }

func (c *Config) SetStripeKeyPrivate(v string) { c.stripeKeyPrivate = v }
func (c *Config) GetStripeKeyPrivate() string  { return c.stripeKeyPrivate }

func (c *Config) SetStripeKeyPublic(v string) { c.stripeKeyPublic = v }
func (c *Config) GetStripeKeyPublic() string  { return c.stripeKeyPublic }

func (c *Config) SetStripeUsed(v bool) { c.stripeUsed = v }
func (c *Config) GetStripeUsed() bool  { return c.stripeUsed }

func (c *Config) SetVertexProjectID(v string) { c.vertexProjectID = v }
func (c *Config) GetVertexProjectID() string  { return c.vertexProjectID }

func (c *Config) SetVertexRegionID(v string) { c.vertexRegionID = v }
func (c *Config) GetVertexRegionID() string  { return c.vertexRegionID }

func (c *Config) SetVertexModelID(v string) { c.vertexModelID = v }
func (c *Config) GetVertexModelID() string  { return c.vertexModelID }

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

// == Daily Analysis Getters/Setters ==
func (c *Config) SetDailyAnalysisSymbols(v []string) { c.dailyAnalysisSymbols = v }
func (c *Config) GetDailyAnalysisSymbols() []string  { return c.dailyAnalysisSymbols }
func (c *Config) SetDailyAnalysisTimeUTC(v string)   { c.dailyAnalysisTimeUTC = v }
func (c *Config) GetDailyAnalysisTimeUTC() string    { return c.dailyAnalysisTimeUTC }
func (c *Config) SetDailyAnalysisCadenceHours(v int) { c.dailyAnalysisCadenceHours = v }
func (c *Config) GetDailyAnalysisCadenceHours() int  { return c.dailyAnalysisCadenceHours }
