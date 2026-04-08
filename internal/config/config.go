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

// New constructs a new configuration instance.
func New() ConfigInterface {
	return &configImplementation{}
}
