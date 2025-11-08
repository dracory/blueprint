package config

import (
	"fmt"
	"strings"

	"github.com/dracory/env"
	"github.com/spf13/cast"
)

// appConfig captures application-level settings.
type appConfig struct {
	name  string
	url   string
	host  string
	port  string
	env   string
	debug bool
}

func loadAppConfig(acc *loadAccumulator) appConfig {
	return appConfig{
		name:  env.GetString(KEY_APP_NAME),
		url:   env.GetString(KEY_APP_URL),
		host:  acc.mustString(KEY_APP_HOST, "set the application host address"),
		port:  acc.mustString(KEY_APP_PORT, "set the application port"),
		env:   acc.mustString(KEY_APP_ENVIRONMENT, "set the application environment"),
		debug: env.GetBool(KEY_APP_DEBUG),
	}
}

// envEncryptionConfig captures optional environment encryption key usage state and derived key.
type envEncryptionConfig struct {
	privateKey string
	derivedKey string
	used       bool
}

func loadEnvEncryptionConfig(acc *loadAccumulator) envEncryptionConfig {
	used := env.GetBool(KEY_ENVENC_USED)

	privateKey := strings.TrimSpace(env.GetString(KEY_ENVENC_KEY_PRIVATE))

	if used {
		if err := ensureRequired(privateKey, KEY_ENVENC_KEY_PRIVATE, "required when ENVENC_USED is yes"); err != nil {
			acc.add(err)
			return envEncryptionConfig{used: used}
		}
	}

	if !used {
		return envEncryptionConfig{privateKey: privateKey, derivedKey: "", used: used}
	}

	derivedKey, err := deriveEnvEncKey(privateKey)
	acc.add(err)
	if err != nil {
		return envEncryptionConfig{used: used}
	}

	return envEncryptionConfig{privateKey: privateKey, derivedKey: derivedKey, used: used}
}

// databaseConfig captures database connection settings.
type databaseConfig struct {
	driver   string
	host     string
	port     string
	name     string
	username string
	password string
	sslMode  string
}

func loadDatabaseConfig(acc *loadAccumulator) databaseConfig {
	driver := acc.mustString(KEY_DB_DRIVER, "select the database driver (e.g., sqlite, postgres)")
	host := strings.TrimSpace(env.GetString(KEY_DB_HOST))
	port := strings.TrimSpace(env.GetString(KEY_DB_PORT))
	name := acc.mustString(KEY_DB_DATABASE, "set the database name")
	user := strings.TrimSpace(env.GetString(KEY_DB_USERNAME))
	pass := strings.TrimSpace(env.GetString(KEY_DB_PASSWORD))

	if driver != driverSQLite {
		acc.mustWhen(true, KEY_DB_HOST, "required when `DB_DRIVER` is not sqlite", host)
		acc.mustWhen(true, KEY_DB_PORT, "required when `DB_DRIVER` is not sqlite", port)
		acc.mustWhen(true, KEY_DB_USERNAME, "required when `DB_DRIVER` is not sqlite", user)
		acc.mustWhen(true, KEY_DB_PASSWORD, "required when `DB_DRIVER` is not sqlite", pass)
	}

	return databaseConfig{
		driver:   driver,
		host:     host,
		port:     port,
		name:     name,
		username: user,
		password: pass,
		sslMode:  "require",
	}
}

// mailConfig captures email delivery settings.
type mailConfig struct {
	driver      string
	fromAddress string
	fromName    string
	host        string
	password    string
	port        int
	username    string
}

func loadMailConfig() mailConfig {
	return mailConfig{
		driver:      env.GetString(KEY_MAIL_DRIVER),
		fromAddress: env.GetString(KEY_MAIL_FROM_ADDRESS),
		fromName:    env.GetString(KEY_MAIL_FROM_NAME),
		host:        env.GetString(KEY_MAIL_HOST),
		password:    env.GetString(KEY_MAIL_PASSWORD),
		port:        cast.ToInt(env.GetString(KEY_MAIL_PORT)),
		username:    env.GetString(KEY_MAIL_USERNAME),
	}
}

// registrationConfig captures authentication registration toggle.
type registrationConfig struct {
	enabled bool
}

func loadRegistrationConfig() registrationConfig {
	return registrationConfig{
		enabled: env.GetBool(KEY_AUTH_REGISTRATION_ENABLED),
	}
}

// storesConfig captures feature store toggles.
type storesConfig struct {
	auditStoreUsed        bool
	blogStoreUsed         bool
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

func loadStoresConfig(acc *loadAccumulator) storesConfig {
	auditStoreUsed := env.GetBool(KEY_AUDIT_STORE_USED)
	blogStoreUsed := env.GetBool(KEY_BLOG_STORE_USED)
	cacheStoreUsed := env.GetBool(KEY_CACHE_STORE_USED)
	cmsStoreUsed := env.GetBool(KEY_CMS_STORE_USED)
	cmsStoreTemplateID := env.GetString(KEY_CMS_STORE_TEMPLATE_ID)
	customStoreUsed := env.GetBool(KEY_CUSTOM_STORE_USED)
	entityStoreUsed := env.GetBool(KEY_ENTITY_STORE_USED)
	feedStoreUsed := env.GetBool(KEY_FEED_STORE_USED)
	geoStoreUsed := env.GetBool(KEY_GEO_STORE_USED)
	logStoreUsed := env.GetBool(KEY_LOG_STORE_USED)
	metaStoreUsed := env.GetBool(KEY_META_STORE_USED)
	sessionStoreUsed := env.GetBool(KEY_SESSION_STORE_USED)
	settingStoreUsed := env.GetBool(KEY_SETTING_STORE_USED)
	shopStoreUsed := env.GetBool(KEY_SHOP_STORE_USED)
	sqlFileStoreUsed := env.GetBool(KEY_SQL_FILE_STORE_USED)
	statsStoreUsed := env.GetBool(KEY_STATS_STORE_USED)
	subscriptionStoreUsed := env.GetBool(KEY_SUBSCRIPTION_STORE_USED)
	taskStoreUsed := env.GetBool(KEY_TASK_STORE_USED)
	userStoreUsed := env.GetBool(KEY_USER_STORE_USED)
	userStoreVaultEnabled := env.GetBool(KEY_USER_STORE_USE_VAULT)
	vaultStoreUsed := env.GetBool(KEY_VAULT_STORE_USED)
	vaultStoreKey := env.GetString(KEY_VAULT_STORE_KEY)

	if userStoreVaultEnabled && !vaultStoreUsed {
		acc.add(fmt.Errorf("%s requires %s to be true", KEY_USER_STORE_USE_VAULT, KEY_VAULT_STORE_USED))
	}

	acc.mustWhen(cmsStoreUsed, KEY_CMS_STORE_TEMPLATE_ID, "required when `CMS_STORE_USED` is true", cmsStoreTemplateID)

	return storesConfig{
		auditStoreUsed:        auditStoreUsed,
		blogStoreUsed:         blogStoreUsed,
		cacheStoreUsed:        cacheStoreUsed,
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

// stripeConfig captures Stripe payment integration settings.
type stripeConfig struct {
	keyPrivate string
	keyPublic  string
	used       bool
}

func loadStripeConfig() stripeConfig {
	keyPrivate := env.GetString(KEY_STRIPE_KEY_PRIVATE)
	keyPublic := env.GetString(KEY_STRIPE_KEY_PUBLIC)
	used := keyPrivate != "" && keyPublic != ""

	return stripeConfig{
		keyPrivate: keyPrivate,
		keyPublic:  keyPublic,
		used:       used,
	}
}

// llmConfig captures LLM provider settings.
type llmConfig struct {
	anthropicUsed         bool
	anthropicKey          string
	anthropicDefaultModel string

	googleGeminiUsed         bool
	googleGeminiKey          string
	googleGeminiDefaultModel string

	openAiUsed         bool
	openAiKey          string
	openAiDefaultModel string

	openRouterUsed         bool
	openRouterKey          string
	openRouterDefaultModel string

	vertexAiUsed         bool
	vertexAiModelID      string
	vertexAiProjectID    string
	vertexAiRegionID     string
	vertexAiDefaultModel string
}

func loadLLMConfig(acc *loadAccumulator) llmConfig {
	anthropicUsed := env.GetBool(KEY_ANTHROPIC_API_USED)
	anthropicKey := env.GetString(KEY_ANTHROPIC_API_KEY)
	anthropicDefaultModel := env.GetString(KEY_ANTHROPIC_API_DEFAULT_MODEL)

	googleGeminiUsed := env.GetBool(KEY_GEMINI_API_USED)
	googleGeminiKey := env.GetString(KEY_GEMINI_API_KEY)
	googleGeminiDefaultModel := env.GetString(KEY_GEMINI_API_DEFAULT_MODEL)

	openAiUsed := env.GetBool(KEY_OPENAI_API_USED)
	openAiKey := env.GetString(KEY_OPENAI_API_KEY)
	openAiDefaultModel := env.GetString(KEY_OPENAI_API_DEFAULT_MODEL)

	openRouterUsed := env.GetBool(KEY_OPENROUTER_API_USED)
	openRouterKey := env.GetString(KEY_OPENROUTER_API_KEY)
	openRouterDefaultModel := env.GetString(KEY_OPENROUTER_API_DEFAULT_MODEL)

	vertexAiUsed := env.GetBool(KEY_VERTEX_AI_API_USED)
	vertexAiModelID := env.GetString(KEY_VERTEX_AI_API_MODEL_ID)
	vertexAiProjectID := env.GetString(KEY_VERTEX_AI_API_PROJECT_ID)
	vertexAiRegionID := env.GetString(KEY_VERTEX_AI_API_REGION_ID)
	vertexAiDefaultModel := env.GetString(KEY_VERTEX_AI_API_DEFAULT_MODEL)

	acc.mustWhen(anthropicUsed, KEY_ANTHROPIC_API_KEY, "required when `ANTHROPIC_API_USED` is true", anthropicKey)
	acc.mustWhen(anthropicUsed, KEY_ANTHROPIC_API_DEFAULT_MODEL, "required when `ANTHROPIC_API_USED` is true", anthropicDefaultModel)

	acc.mustWhen(googleGeminiUsed, KEY_GEMINI_API_KEY, "required when `GEMINI_API_USED` is true", googleGeminiKey)
	acc.mustWhen(googleGeminiUsed, KEY_GEMINI_API_DEFAULT_MODEL, "required when `GEMINI_API_USED` is true", googleGeminiDefaultModel)

	acc.mustWhen(openAiUsed, KEY_OPENAI_API_KEY, "required when `OPENAI_API_USED` is true", openAiKey)
	acc.mustWhen(openAiUsed, KEY_OPENAI_API_DEFAULT_MODEL, "required when `OPENAI_API_USED` is true", openAiDefaultModel)

	acc.mustWhen(openRouterUsed, KEY_OPENROUTER_API_KEY, "required when `OPENROUTER_API_USED` is true", openRouterKey)
	acc.mustWhen(openRouterUsed, KEY_OPENROUTER_API_DEFAULT_MODEL, "required when `OPENROUTER_API_USED` is true", openRouterDefaultModel)

	acc.mustWhen(vertexAiUsed, KEY_VERTEX_AI_API_MODEL_ID, "required when `VERTEX_AI_API_USED` is true", vertexAiModelID)
	acc.mustWhen(vertexAiUsed, KEY_VERTEX_AI_API_PROJECT_ID, "required when `VERTEX_AI_API_USED` is true", vertexAiProjectID)
	acc.mustWhen(vertexAiUsed, KEY_VERTEX_AI_API_REGION_ID, "required when `VERTEX_AI_API_USED` is true", vertexAiRegionID)
	acc.mustWhen(vertexAiUsed, KEY_VERTEX_AI_API_DEFAULT_MODEL, "required when `VERTEX_AI_API_USED` is true", vertexAiDefaultModel)

	return llmConfig{
		anthropicUsed:            anthropicUsed,
		anthropicKey:             anthropicKey,
		anthropicDefaultModel:    anthropicDefaultModel,
		googleGeminiUsed:         googleGeminiUsed,
		googleGeminiKey:          googleGeminiKey,
		googleGeminiDefaultModel: googleGeminiDefaultModel,
		openAiUsed:               openAiUsed,
		openAiKey:                openAiKey,
		openAiDefaultModel:       openAiDefaultModel,
		openRouterUsed:           openRouterUsed,
		openRouterKey:            openRouterKey,
		openRouterDefaultModel:   openRouterDefaultModel,
		vertexAiUsed:             vertexAiUsed,
		vertexAiModelID:          vertexAiModelID,
		vertexAiProjectID:        vertexAiProjectID,
		vertexAiRegionID:         vertexAiRegionID,
		vertexAiDefaultModel:     vertexAiDefaultModel,
	}
}

// translationConfig captures i18n settings.
type translationConfig struct {
	defaultLanguage string
	languageList    map[string]string
}

func loadTranslationConfig() translationConfig {
	defaultLang := env.GetString(KEY_TRANSLATION_LANGUAGE_DEFAULT)
	if defaultLang == "" {
		defaultLang = translationLanguageDefault()
	}

	return translationConfig{
		defaultLanguage: defaultLang,
		languageList:    translationLanguageListDefault(),
	}
}
