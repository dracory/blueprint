package config

import (
	"strings"

	"project/internal/types"

	"github.com/dracory/env"
)

type variableType string

const (
	variableTypeString variableType = "string"
	variableTypeBool   variableType = "bool"
	variableTypeInt    variableType = "int"
	variableTypeFloat  variableType = "float"
)

type setter func(types.ConfigInterface, string) error

type variable struct {
	// key to load
	key string

	// required flag
	required bool

	// required when
	requiredWhen func() bool

	// requiredMessage is shown when the variable must be provided
	requiredMessage string

	// preserveWhitespace controls whether leading/trailing whitespace is retained
	preserveWhitespace bool

	// variable type
	variableType variableType

	// DEPRECATED assign function
	assign setter

	// assign functions
	assignBool   func(cfg types.ConfigInterface, value bool) error
	assignInt    func(cfg types.ConfigInterface, value int) error
	assignFloat  func(cfg types.ConfigInterface, value float64) error
	assignString func(cfg types.ConfigInterface, value string) error
}

func variableDefinitions() []variable {
	return []variable{
		// == App configuration ==
		{
			key: KEY_APP_NAME,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetAppName(value)
				return nil
			},
		},
		{
			key:             KEY_APP_URL,
			requiredMessage: "set the application URL",
			required:        true,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetAppUrl(value)
				return nil
			},
		},
		{
			key:             KEY_APP_ENVIRONMENT,
			requiredMessage: "set the application environment",
			required:        true,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetAppEnv(value)
				return nil
			},
		},
		{
			key:             KEY_APP_HOST,
			requiredMessage: "set the application host address",
			required:        true,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetAppHost(value)
				return nil
			},
		},
		{
			key:             KEY_APP_PORT,
			requiredMessage: "set the application port",
			required:        true,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetAppPort(value)
				return nil
			},
		},
		{
			key: KEY_APP_DEBUG,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetAppDebug(value)
				return nil
			},
		},

		// == Database configuration ==
		{
			key:             KEY_DB_DRIVER,
			requiredMessage: "select the database driver (e.g., sqlite, postgres)",
			required:        true,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetDatabaseDriver(value)
				return nil
			},
		},
		{
			key: KEY_DB_HOST,
			requiredWhen: func() bool {
				return strings.TrimSpace(env.GetString(KEY_DB_DRIVER)) != driverSQLite
			},
			requiredMessage: "required when `DB_DRIVER` is not sqlite",
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetDatabaseHost(value)
				return nil
			},
		},
		{
			key: KEY_DB_PORT,
			requiredWhen: func() bool {
				return strings.TrimSpace(env.GetString(KEY_DB_DRIVER)) != driverSQLite
			},
			requiredMessage: "required when `DB_DRIVER` is not sqlite",
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetDatabasePort(value)
				return nil
			},
		},
		{
			key:             KEY_DB_DATABASE,
			requiredMessage: "set the database name",
			required:        true,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetDatabaseName(value)
				return nil
			},
		},
		{
			key: KEY_DB_USERNAME,
			requiredWhen: func() bool {
				return strings.TrimSpace(env.GetString(KEY_DB_DRIVER)) != driverSQLite
			},
			requiredMessage: "required when `DB_DRIVER` is not sqlite",
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetDatabaseUsername(value)
				return nil
			},
		},
		{
			key: KEY_DB_PASSWORD,
			requiredWhen: func() bool {
				return strings.TrimSpace(env.GetString(KEY_DB_DRIVER)) != driverSQLite
			},
			requiredMessage: "required when `DB_DRIVER` is not sqlite",
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetDatabasePassword(value)
				return nil
			},
		},

		// == Mail configuration ==
		{
			key: KEY_MAIL_DRIVER,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetMailDriver(value)
				return nil
			},
		},
		{
			key: KEY_MAIL_HOST,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetMailHost(value)
				return nil
			},
		},
		{
			key:          KEY_MAIL_PORT,
			variableType: variableTypeInt,
			assignInt: func(cfg types.ConfigInterface, value int) error {
				cfg.SetMailPort(value)
				return nil
			},
		},
		{
			key: KEY_MAIL_USERNAME,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetMailUsername(value)
				return nil
			},
		},
		{
			key: KEY_MAIL_PASSWORD,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetMailPassword(value)
				return nil
			},
		},
		{
			key: KEY_MAIL_FROM_ADDRESS,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetMailFromAddress(value)
				return nil
			},
		},
		{
			key: KEY_MAIL_FROM_NAME,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetMailFromName(value)
				return nil
			},
		},

		// == Store configuration ==
		{
			key:          KEY_AUDIT_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetAuditStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_CACHE_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetCacheStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_CMS_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetCmsStoreUsed(value)
				return nil
			},
		},
		{
			key: KEY_CMS_STORE_TEMPLATE_ID,
			requiredWhen: func() bool {
				return env.GetBool(KEY_CMS_STORE_USED)
			},
			requiredMessage: "required when `CMS_STORE_USED` is true",
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetCmsStoreTemplateID(value)
				return nil
			},
		},
		{
			key:          KEY_CUSTOM_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetCustomStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_ENTITY_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetEntityStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_FEED_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetFeedStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_GEO_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetGeoStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_LOG_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetLogStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_META_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetMetaStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_SESSION_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetSessionStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_SETTING_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetSettingStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_SHOP_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetShopStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_STATS_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetStatsStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_SUBSCRIPTION_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetSubscriptionStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_TASK_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetTaskStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_USER_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetUserStoreUsed(value)
				return nil
			},
		},
		{
			key:          KEY_USER_STORE_USE_VAULT,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetUserStoreVaultEnabled(value)
				return nil
			},
		},
		{
			key:          KEY_VAULT_STORE_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetVaultStoreUsed(value)
				return nil
			},
		},
		{
			key: KEY_VAULT_STORE_KEY,
			requiredWhen: func() bool {
				return env.GetBool(KEY_VAULT_STORE_USED)
			},
			requiredMessage: "required when `VAULT_STORE_USED` is true",
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetVaultStoreKey(value)
				return nil
			},
		},

		// == LLM configuration ==
		{
			key:          KEY_GEMINI_API_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetGoogleGeminiApiUsed(value)
				return nil
			},
		},
		{
			key: KEY_GEMINI_API_KEY,
			requiredWhen: func() bool {
				return env.GetBool(KEY_GEMINI_API_USED)
			},
			requiredMessage: "required when `GEMINI_API_USED` is true",
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetGoogleGeminiApiKey(value)
				return nil
			},
		},
		{
			key: KEY_GEMINI_API_DEFAULT_MODEL,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetGoogleGeminiApiDefaultModel(value)
				return nil
			},
		},
		{
			key:          KEY_OPENAI_API_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetOpenAiApiUsed(value)
				return nil
			},
		},
		{
			key: KEY_OPENAI_API_KEY,
			requiredWhen: func() bool {
				return env.GetBool(KEY_OPENAI_API_USED)
			},
			requiredMessage: "required when `OPENAI_API_USED` is true",
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetOpenAiApiKey(value)
				return nil
			},
		},
		{
			key: KEY_OPENAI_API_DEFAULT_MODEL,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetOpenAiApiDefaultModel(value)
				return nil
			},
		},
		{
			key:          KEY_OPENROUTER_API_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetOpenRouterApiUsed(value)
				return nil
			},
		},
		{
			key: KEY_OPENROUTER_API_KEY,
			requiredWhen: func() bool {
				return env.GetBool(KEY_OPENROUTER_API_USED)
			},
			requiredMessage: "required when `OPENROUTER_API_USED` is true",
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetOpenRouterApiKey(value)
				return nil
			},
		},
		{
			key: KEY_OPENROUTER_API_DEFAULT_MODEL,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetOpenRouterApiDefaultModel(value)
				return nil
			},
		},
		{
			key:          KEY_ANTHROPIC_API_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetAnthropicApiUsed(value)
				return nil
			},
		},
		{
			key: KEY_ANTHROPIC_API_KEY,
			requiredWhen: func() bool {
				return env.GetBool(KEY_ANTHROPIC_API_USED)
			},
			requiredMessage: "required when `ANTHROPIC_API_USED` is true",
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetAnthropicApiKey(value)
				return nil
			},
		},
		{
			key: KEY_ANTHROPIC_API_DEFAULT_MODEL,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetAnthropicApiDefaultModel(value)
				return nil
			},
		},
		{
			key:          KEY_VERTEX_AI_API_USED,
			variableType: variableTypeBool,
			assignBool: func(cfg types.ConfigInterface, value bool) error {
				cfg.SetVertexAiApiUsed(value)
				return nil
			},
		},
		{
			key: KEY_VERTEX_AI_API_MODEL_ID,
			requiredWhen: func() bool {
				return env.GetBool(KEY_VERTEX_AI_API_USED)
			},
			requiredMessage: "required when `VERTEX_AI_API_USED` is true",
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetVertexAiApiModelID(value)
				return nil
			},
		},
		{
			key: KEY_VERTEX_AI_API_PROJECT_ID,
			requiredWhen: func() bool {
				return env.GetBool(KEY_VERTEX_AI_API_USED)
			},
			requiredMessage: "required when `VERTEX_AI_API_USED` is true",
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetVertexAiApiProjectID(value)
				return nil
			},
		},
		{
			key: KEY_VERTEX_AI_API_REGION_ID,
			requiredWhen: func() bool {
				return env.GetBool(KEY_VERTEX_AI_API_USED)
			},
			requiredMessage: "required when `VERTEX_AI_API_USED` is true",
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetVertexAiApiRegionID(value)
				return nil
			},
		},
		{
			key: KEY_VERTEX_AI_API_DEFAULT_MODEL,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetVertexAiApiDefaultModel(value)
				return nil
			},
		},

		// == Payment configuration ==
		{
			key: KEY_STRIPE_KEY_PRIVATE,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetStripeKeyPrivate(value)
				return nil
			},
		},
		{
			key: KEY_STRIPE_KEY_PUBLIC,
			assignString: func(cfg types.ConfigInterface, value string) error {
				cfg.SetStripeKeyPublic(value)
				return nil
			},
		},

		// == Translation configuration ==
		{
			key: KEY_TRANSLATION_LANGUAGE_DEFAULT,
			assignString: func(cfg types.ConfigInterface, value string) error {
				trimmed := strings.TrimSpace(value)
				if trimmed == "" {
					cfg.SetTranslationLanguageDefault(translationLanguageDefault())
					return nil
				}

				cfg.SetTranslationLanguageDefault(trimmed)
				return nil
			},
		},
	}
}
