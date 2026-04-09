// The z_ prefix keeps this file sorted after user-configurable files in the directory listing.

package config

import (
	"project/internal/resources"

	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// envValidator is a local alias for env.Validator for use in config loaders.
type envValidator = env.Validator

// configImplementation holds all configuration values.
type configImplementation struct {
	// App configuration
	appName  string
	appType  string
	appEnv   string
	appHost  string
	appPort  string
	appUrl   string
	appDebug bool

	// Database configuration
	databaseDriver   string
	databaseHost     string
	databasePort     string
	databaseName     string
	databaseUsername string
	databasePassword string
	databaseSSLMode  string

	// Email configuration
	emailDriver      string
	emailHost        string
	emailPort        int
	emailUsername    string
	emailPassword    string
	emailFromName    string
	emailFromAddress string

	// LLM configuration
	openRouterApiKey            string
	openRouterApiUsed           bool
	openRouterApiDefaultModel   string
	openAiApiKey                string
	openAiApiUsed               bool
	openAiApiDefaultModel       string
	anthropicApiUsed            bool
	anthropicApiKey             string
	anthropicApiDefaultModel    string
	googleGeminiApiUsed         bool
	googleGeminiApiKey          string
	googleGeminiApiDefaultModel string
	vertexAiApiUsed             bool
	vertexAiApiDefaultModel     string
	vertexAiApiProjectID        string
	vertexAiApiRegionID         string
	vertexAiApiModelID          string

	// Media configuration
	mediaBucket   string
	mediaDriver   string
	mediaKey      string
	mediaEndpoint string
	mediaRegion   string
	mediaRoot     string
	mediaSecret   string
	mediaUrl      string

	// Payment configuration
	stripeKeyPrivate string
	stripeKeyPublic  string
	stripeUsed       bool

	// Authentication
	registrationEnabled bool

	// i18n / Translation
	translationLanguageDefault string
	translationLanguageList    map[string]string

	// SEO configuration
	indexNowKey string

	// Encryption
	envEncryptionKey string

	// CMS MCP
	cmsMcpApiKey string

	// Store flags
	auditStoreUsed        bool
	blogStoreUsed         bool
	chatStoreUsed         bool
	cacheStoreUsed        bool
	cmsStoreUsed          bool
	cmsStoreTemplateID    string
	customStoreUsed       bool
	entityStoreUsed       bool
	feedStoreUsed         bool
	geoStoreUsed          bool
	logStoreUsed          bool
	metaStoreUsed         bool
	sessionStoreUsed      bool
	settingStoreUsed      bool
	shopStoreUsed         bool
	sqlFileStoreUsed      bool
	statsStoreUsed        bool
	subscriptionStoreUsed bool
	taskStoreUsed         bool
	userStoreUsed         bool
	userStoreVaultEnabled bool
	vaultStoreUsed        bool
	vaultStoreKey         string
}

// New constructs a new configuration instance.
func New() ConfigInterface {
	return &configImplementation{}
}

// NewFromEnv constructs a configuration instance populated from environment variables.
func NewFromEnv() (ConfigInterface, error) {
	env.Load(".env")

	v := &envValidator{}
	cfg := &configImplementation{}

	// Load app config first to get app.env
	cfg.setAppConfig(appConfig(v))

	// Load encryption config and check if encryption is used
	privateKey := env.GetString(KEY_ENVENC_KEY_PRIVATE)
	encryptionUsed := privateKey != ""

	if encryptionUsed {
		v.RequireWhen(true, KEY_ENVENC_KEY_PRIVATE,
			"required when encryption is enabled", privateKey)
	}

	cfg.SetEnvEncryptionKey(privateKey)

	// Initialize encrypted environment variables BEFORE other config loaders read them
	if encryptionUsed {
		if err := baseCfg.InitializeEnvEncVariablesFromResources(cfg.GetAppEnv(), ENVENC_KEY_PUBLIC, privateKey, resources.Resource); err != nil {
			v.Add(err)
		} else {
			cfg.SetEnvEncryptionKey("removed") // reset the private key
		}

		// Reload app config to pick up any encrypted app variables
		cfg.setAppConfig(appConfig(v))
	}

	// Now load remaining config sections - they will have access to encrypted variables
	cfg.setDatabaseConfig(databaseConfig(v))
	cfg.setMailConfig(emailConfig())
	cfg.setAuthConfig(authConfig())
	cfg.setStoresConfig(storesConfig(v))
	cfg.setStripeConfig(paymentConfig())
	cfg.setLLMConfig(llmConfig(v))
	cfg.setTranslationConfig(i18nConfig())

	if err := v.Err(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Ensure configImplementation satisfies ConfigInterface
var _ ConfigInterface = (*configImplementation)(nil)

// ============================================================================
// CMS MCP Config Implementation
// ============================================================================

func (c *configImplementation) SetCmsMcpApiKey(v string) {
	c.cmsMcpApiKey = v
}

func (c *configImplementation) GetCmsMcpApiKey() string {
	return c.cmsMcpApiKey
}

// ============================================================================
// App Config Implementation
// ============================================================================

func (c *configImplementation) setAppConfig(s appSettings) {
	c.appName = s.name
	c.appUrl = s.url
	c.appHost = s.host
	c.appPort = s.port
	c.appEnv = s.env
	c.appDebug = s.debug
	c.cmsMcpApiKey = s.cmsMcpApiKey
}

func (c *configImplementation) SetAppName(appName string) {
	c.appName = appName
}

func (c *configImplementation) GetAppName() string {
	return c.appName
}

func (c *configImplementation) SetAppType(appType string) {
	c.appType = appType
}

func (c *configImplementation) GetAppType() string {
	return c.appType
}

func (c *configImplementation) SetAppEnv(appEnv string) {
	c.appEnv = appEnv
}

func (c *configImplementation) GetAppEnv() string {
	return c.appEnv
}

func (c *configImplementation) SetAppHost(appHost string) {
	c.appHost = appHost
}

func (c *configImplementation) GetAppHost() string {
	return c.appHost
}

func (c *configImplementation) SetAppPort(appPort string) {
	c.appPort = appPort
}

func (c *configImplementation) GetAppPort() string {
	return c.appPort
}

func (c *configImplementation) SetAppUrl(appUrl string) {
	c.appUrl = appUrl
}

func (c *configImplementation) GetAppUrl() string {
	return c.appUrl
}

func (c *configImplementation) SetAppDebug(appDebug bool) {
	c.appDebug = appDebug
}

func (c *configImplementation) GetAppDebug() bool {
	return c.appDebug
}

func (c *configImplementation) IsEnvDevelopment() bool {
	return c.appEnv == "development"
}

func (c *configImplementation) IsEnvLocal() bool {
	return c.appEnv == "local"
}

func (c *configImplementation) IsEnvProduction() bool {
	return c.appEnv == "production"
}

func (c *configImplementation) IsEnvStaging() bool {
	return c.appEnv == "staging"
}

func (c *configImplementation) IsEnvTesting() bool {
	return c.appEnv == "testing"
}

// ============================================================================
// Auth Config Implementation
// ============================================================================

func (c *configImplementation) setAuthConfig(s authSettings) {
	c.registrationEnabled = s.registrationEnabled
}

func (c *configImplementation) SetRegistrationEnabled(v bool) {
	c.registrationEnabled = v
}

func (c *configImplementation) GetRegistrationEnabled() bool {
	return c.registrationEnabled
}

// ============================================================================
// Database Config Implementation
// ============================================================================

func (c *configImplementation) setDatabaseConfig(s databaseSettings) {
	c.databaseDriver = s.driver
	c.databaseHost = s.host
	c.databasePort = s.port
	c.databaseName = s.name
	c.databaseUsername = s.user
	c.databasePassword = s.pass
	c.databaseSSLMode = "require"
}

func (c *configImplementation) SetDatabaseDriver(v string) {
	c.databaseDriver = v
}

func (c *configImplementation) GetDatabaseDriver() string {
	return c.databaseDriver
}

func (c *configImplementation) SetDatabaseHost(v string) {
	c.databaseHost = v
}

func (c *configImplementation) GetDatabaseHost() string {
	return c.databaseHost
}

func (c *configImplementation) SetDatabasePort(v string) {
	c.databasePort = v
}

func (c *configImplementation) GetDatabasePort() string {
	return c.databasePort
}

func (c *configImplementation) SetDatabaseName(v string) {
	c.databaseName = v
}

func (c *configImplementation) GetDatabaseName() string {
	return c.databaseName
}

func (c *configImplementation) SetDatabaseUsername(v string) {
	c.databaseUsername = v
}

func (c *configImplementation) GetDatabaseUsername() string {
	return c.databaseUsername
}

func (c *configImplementation) SetDatabasePassword(v string) {
	c.databasePassword = v
}

func (c *configImplementation) GetDatabasePassword() string {
	return c.databasePassword
}

func (c *configImplementation) SetDatabaseSSLMode(v string) {
	c.databaseSSLMode = v
}

func (c *configImplementation) GetDatabaseSSLMode() string {
	return c.databaseSSLMode
}

// ============================================================================
// Email Config Implementation
// ============================================================================

func (c *configImplementation) setMailConfig(s emailSettings) {
	c.emailDriver = s.driver
	c.emailFromAddress = s.fromAddress
	c.emailFromName = s.fromName
	c.emailHost = s.host
	c.emailPassword = s.password
	c.emailPort = s.port
	c.emailUsername = s.username
}

func (c *configImplementation) SetMailDriver(v string) {
	c.emailDriver = v
}

func (c *configImplementation) GetMailDriver() string {
	return c.emailDriver
}

func (c *configImplementation) SetMailHost(v string) {
	c.emailHost = v
}

func (c *configImplementation) GetMailHost() string {
	return c.emailHost
}

func (c *configImplementation) SetMailPort(v int) {
	c.emailPort = v
}

func (c *configImplementation) GetMailPort() int {
	return c.emailPort
}

func (c *configImplementation) SetMailUsername(v string) {
	c.emailUsername = v
}

func (c *configImplementation) GetMailUsername() string {
	return c.emailUsername
}

func (c *configImplementation) SetMailPassword(v string) {
	c.emailPassword = v
}

func (c *configImplementation) GetMailPassword() string {
	return c.emailPassword
}

func (c *configImplementation) SetMailFromName(v string) {
	c.emailFromName = v
}

func (c *configImplementation) GetMailFromName() string {
	return c.emailFromName
}

func (c *configImplementation) SetMailFromAddress(v string) {
	c.emailFromAddress = v
}

func (c *configImplementation) GetMailFromAddress() string {
	return c.emailFromAddress
}

// ============================================================================
// Encryption Config Implementation
// ============================================================================

func (c *configImplementation) SetEnvEncryptionKey(v string) {
	c.envEncryptionKey = v
}

func (c *configImplementation) GetEnvEncryptionKey() string {
	return c.envEncryptionKey
}

// ============================================================================
// i18n Config Implementation
// ============================================================================

func (c *configImplementation) setTranslationConfig(s i18nSettings) {
	c.translationLanguageDefault = s.defaultLanguage
	c.translationLanguageList = s.languageList
}

func (c *configImplementation) SetTranslationLanguageDefault(v string) {
	c.translationLanguageDefault = v
}

func (c *configImplementation) GetTranslationLanguageDefault() string {
	return c.translationLanguageDefault
}

func (c *configImplementation) SetTranslationLanguageList(v map[string]string) {
	c.translationLanguageList = v
}

func (c *configImplementation) GetTranslationLanguageList() map[string]string {
	return c.translationLanguageList
}

// ============================================================================
// LLM Config Implementation
// ============================================================================

// ============================================================================
// LLM Config Implementation
// ============================================================================

func (c *configImplementation) setLLMConfig(s llmSettings) {
	c.anthropicApiUsed = s.anthropicUsed
	c.anthropicApiKey = s.anthropicKey
	c.anthropicApiDefaultModel = s.anthropicDefaultModel
	c.googleGeminiApiUsed = s.googleGeminiUsed
	c.googleGeminiApiKey = s.googleGeminiKey
	c.googleGeminiApiDefaultModel = s.googleGeminiDefaultModel
	c.openAiApiUsed = s.openAiUsed
	c.openAiApiKey = s.openAiKey
	c.openAiApiDefaultModel = s.openAiDefaultModel
	c.openRouterApiUsed = s.openRouterUsed
	c.openRouterApiKey = s.openRouterKey
	c.openRouterApiDefaultModel = s.openRouterDefaultModel
	c.vertexAiApiUsed = s.vertexAiUsed
	c.vertexAiApiModelID = s.vertexAiModelID
	c.vertexAiApiProjectID = s.vertexAiProjectID
	c.vertexAiApiRegionID = s.vertexAiRegionID
	c.vertexAiApiDefaultModel = s.vertexAiDefaultModel
}

// Anthropic
func (c *configImplementation) SetAnthropicApiUsed(v bool) {
	c.anthropicApiUsed = v
}

func (c *configImplementation) GetAnthropicApiUsed() bool {
	return c.anthropicApiUsed
}

func (c *configImplementation) SetAnthropicApiKey(v string) {
	c.anthropicApiKey = v
}

func (c *configImplementation) GetAnthropicApiKey() string {
	return c.anthropicApiKey
}

func (c *configImplementation) SetAnthropicApiDefaultModel(v string) {
	c.anthropicApiDefaultModel = v
}

func (c *configImplementation) GetAnthropicApiDefaultModel() string {
	return c.anthropicApiDefaultModel
}

// Google Gemini
func (c *configImplementation) SetGoogleGeminiApiUsed(v bool) {
	c.googleGeminiApiUsed = v
}

func (c *configImplementation) GetGoogleGeminiApiUsed() bool {
	return c.googleGeminiApiUsed
}

func (c *configImplementation) SetGoogleGeminiApiKey(v string) {
	c.googleGeminiApiKey = v
}

func (c *configImplementation) GetGoogleGeminiApiKey() string {
	return c.googleGeminiApiKey
}

func (c *configImplementation) SetGoogleGeminiApiDefaultModel(v string) {
	c.googleGeminiApiDefaultModel = v
}

func (c *configImplementation) GetGoogleGeminiApiDefaultModel() string {
	return c.googleGeminiApiDefaultModel
}

// OpenAI
func (c *configImplementation) SetOpenAiApiUsed(v bool) {
	c.openAiApiUsed = v
}

func (c *configImplementation) GetOpenAiApiUsed() bool {
	return c.openAiApiUsed
}

func (c *configImplementation) SetOpenAiApiKey(v string) {
	c.openAiApiKey = v
}

func (c *configImplementation) GetOpenAiApiKey() string {
	return c.openAiApiKey
}

func (c *configImplementation) SetOpenAiApiDefaultModel(v string) {
	c.openAiApiDefaultModel = v
}

func (c *configImplementation) GetOpenAiApiDefaultModel() string {
	return c.openAiApiDefaultModel
}

// OpenRouter
func (c *configImplementation) SetOpenRouterApiKey(v string) {
	c.openRouterApiKey = v
}

func (c *configImplementation) GetOpenRouterApiKey() string {
	return c.openRouterApiKey
}

func (c *configImplementation) SetOpenRouterApiUsed(v bool) {
	c.openRouterApiUsed = v
}

func (c *configImplementation) GetOpenRouterApiUsed() bool {
	return c.openRouterApiUsed
}

func (c *configImplementation) SetOpenRouterApiDefaultModel(v string) {
	c.openRouterApiDefaultModel = v
}

func (c *configImplementation) GetOpenRouterApiDefaultModel() string {
	return c.openRouterApiDefaultModel
}

// Vertex AI
func (c *configImplementation) SetVertexAiApiUsed(v bool) {
	c.vertexAiApiUsed = v
}

func (c *configImplementation) GetVertexAiApiUsed() bool {
	return c.vertexAiApiUsed
}

func (c *configImplementation) SetVertexAiApiDefaultModel(v string) {
	c.vertexAiApiDefaultModel = v
}

func (c *configImplementation) GetVertexAiApiDefaultModel() string {
	return c.vertexAiApiDefaultModel
}

func (c *configImplementation) SetVertexAiApiProjectID(v string) {
	c.vertexAiApiProjectID = v
}

func (c *configImplementation) GetVertexAiApiProjectID() string {
	return c.vertexAiApiProjectID
}

func (c *configImplementation) SetVertexAiApiRegionID(v string) {
	c.vertexAiApiRegionID = v
}

func (c *configImplementation) GetVertexAiApiRegionID() string {
	return c.vertexAiApiRegionID
}

func (c *configImplementation) SetVertexAiApiModelID(v string) {
	c.vertexAiApiModelID = v
}

func (c *configImplementation) GetVertexAiApiModelID() string {
	return c.vertexAiApiModelID
}

// ============================================================================
// Media Config Implementation
// ============================================================================

func (c *configImplementation) SetMediaBucket(v string) {
	c.mediaBucket = v
}

func (c *configImplementation) GetMediaBucket() string {
	return c.mediaBucket
}

func (c *configImplementation) SetMediaDriver(v string) {
	c.mediaDriver = v
}

func (c *configImplementation) GetMediaDriver() string {
	return c.mediaDriver
}

func (c *configImplementation) SetMediaKey(v string) {
	c.mediaKey = v
}

func (c *configImplementation) GetMediaKey() string {
	return c.mediaKey
}

func (c *configImplementation) SetMediaEndpoint(v string) {
	c.mediaEndpoint = v
}

func (c *configImplementation) GetMediaEndpoint() string {
	return c.mediaEndpoint
}

func (c *configImplementation) SetMediaRegion(v string) {
	c.mediaRegion = v
}

func (c *configImplementation) GetMediaRegion() string {
	return c.mediaRegion
}

func (c *configImplementation) SetMediaRoot(v string) {
	c.mediaRoot = v
}

func (c *configImplementation) GetMediaRoot() string {
	return c.mediaRoot
}

func (c *configImplementation) SetMediaSecret(v string) {
	c.mediaSecret = v
}

func (c *configImplementation) GetMediaSecret() string {
	return c.mediaSecret
}

func (c *configImplementation) SetMediaUrl(v string) {
	c.mediaUrl = v
}

func (c *configImplementation) GetMediaUrl() string {
	return c.mediaUrl
}

// ============================================================================
// Payment Config Implementation
// ============================================================================

func (c *configImplementation) setStripeConfig(s paymentSettings) {
	c.stripeKeyPrivate = s.keyPrivate
	c.stripeKeyPublic = s.keyPublic
	c.stripeUsed = s.used
}

func (c *configImplementation) SetStripeKeyPrivate(v string) {
	c.stripeKeyPrivate = v
}

func (c *configImplementation) GetStripeKeyPrivate() string {
	return c.stripeKeyPrivate
}

func (c *configImplementation) SetStripeKeyPublic(v string) {
	c.stripeKeyPublic = v
}

func (c *configImplementation) GetStripeKeyPublic() string {
	return c.stripeKeyPublic
}

func (c *configImplementation) SetStripeUsed(v bool) {
	c.stripeUsed = v
}

func (c *configImplementation) GetStripeUsed() bool {
	return c.stripeUsed
}

// ============================================================================
// SEO Config Implementation
// ============================================================================

func (c *configImplementation) SetIndexNowKey(v string) {
	c.indexNowKey = v
}

func (c *configImplementation) GetIndexNowKey() string {
	return c.indexNowKey
}

// ============================================================================
// Stores Config Implementation
// ============================================================================

// ============================================================================
// Stores Config Implementation
// ============================================================================

func (c *configImplementation) setStoresConfig(s storesSettings) {
	c.auditStoreUsed = auditStoreUsed
	c.blogStoreUsed = blogStoreUsed
	c.cacheStoreUsed = cacheStoreUsed
	c.chatStoreUsed = chatStoreUsed
	c.cmsStoreUsed = cmsStoreUsed
	c.cmsStoreTemplateID = s.cmsStoreTemplateID
	c.customStoreUsed = customStoreUsed
	c.entityStoreUsed = entityStoreUsed
	c.feedStoreUsed = feedStoreUsed
	c.geoStoreUsed = geoStoreUsed
	c.logStoreUsed = logStoreUsed
	c.metaStoreUsed = metaStoreUsed
	c.sessionStoreUsed = sessionStoreUsed
	c.settingStoreUsed = settingStoreUsed
	c.shopStoreUsed = shopStoreUsed
	c.sqlFileStoreUsed = sqlFileStoreUsed
	c.statsStoreUsed = statsStoreUsed
	c.subscriptionStoreUsed = subscriptionStoreUsed
	c.taskStoreUsed = taskStoreUsed
	c.userStoreUsed = userStoreUsed
	c.userStoreVaultEnabled = userStoreVaultEnabled
	c.vaultStoreUsed = vaultStoreUsed
	c.vaultStoreKey = s.vaultStoreKey
}

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
