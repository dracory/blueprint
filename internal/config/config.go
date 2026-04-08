package config

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

	// CMS MCP (small interface, kept here)
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

// Ensure configImplementation satisfies ConfigInterface
var _ ConfigInterface = (*configImplementation)(nil)

// New constructs a new configuration instance.
func New() ConfigInterface {
	return &configImplementation{}
}

// CMS MCP methods (small interface, kept here)
func (c *configImplementation) SetCmsMcpApiKey(v string) {
	c.cmsMcpApiKey = v
}

func (c *configImplementation) GetCmsMcpApiKey() string {
	return c.cmsMcpApiKey
}
