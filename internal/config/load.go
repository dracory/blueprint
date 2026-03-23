package config

import (
	"project/internal/resources"

	baseCfg "github.com/dracory/base/config"
	"github.com/dracory/env"
)

// Load loads the configuration
//
// Business logic:
//   - initializes the environment variables from the .env file
//   - initializes envenc variables based on the app environment
//   - checks all the required env variables
//   - panics if any of the required variable is missing
//
// Parameters:
// - none
//
// Returns:
// - none

func Load() (ConfigInterface, error) {
	env.Load(".env")

	acc := &baseCfg.LoadAccumulator{}

	// Load basic config sections first to get app.env and envEnc.privateKey
	app := loadAppConfig(acc)
	envEnc := loadEnvEncryptionConfig(acc)

	// Initialize encrypted environment variables BEFORE other config loaders read them
	if envEnc.used {
		// Use base package config loader with embedded resources support
		if err := baseCfg.InitializeEnvEncVariablesFromResources(app.env, ENVENC_KEY_PUBLIC, envEnc.privateKey, resources.Resource); err != nil {
			acc.Add(err)
		} else {
			envEnc.privateKey = "removed" // reset the private key
		}

		// Reload app config to pick up any encrypted app variables (APP_NAME, APP_URL, etc.)
		app = loadAppConfig(acc)
	}

	// Now load remaining config sections - they will have access to encrypted variables
	db := loadDatabaseConfig(acc)
	mail := loadMailConfig()
	reg := loadRegistrationConfig()
	stores := loadStoresConfig(acc)
	stripe := loadStripeConfig()
	llms := loadLLMConfig(acc)
	trans := loadTranslationConfig()

	if err := acc.Err(); err != nil {
		return nil, err
	}

	cfg := New()

	cfg.SetAppName(app.name)
	cfg.SetAppUrl(app.url)
	cfg.SetAppHost(app.host)
	cfg.SetAppPort(app.port)
	cfg.SetAppEnv(app.env)
	cfg.SetAppDebug(app.debug)

	cfg.SetDatabaseDriver(db.driver)
	cfg.SetDatabaseHost(db.host)
	cfg.SetDatabasePort(db.port)
	cfg.SetDatabaseName(db.name)
	cfg.SetDatabaseUsername(db.username)
	cfg.SetDatabasePassword(db.password)
	cfg.SetDatabaseSSLMode(db.sslMode)

	cfg.SetMailDriver(mail.driver)
	cfg.SetMailFromAddress(mail.fromAddress)
	cfg.SetMailFromName(mail.fromName)
	cfg.SetMailHost(mail.host)
	cfg.SetMailPassword(mail.password)
	cfg.SetMailPort(mail.port)
	cfg.SetMailUsername(mail.username)

	cfg.SetRegistrationEnabled(reg.enabled)

	cfg.SetAuditStoreUsed(stores.auditStoreUsed)
	cfg.SetBlogStoreUsed(stores.blogStoreUsed)
	cfg.SetCacheStoreUsed(stores.cacheStoreUsed)
	cfg.SetCmsStoreUsed(stores.cmsStoreUsed)
	cfg.SetCmsMcpApiKey(app.cmsMcpApiKey)
	cfg.SetCmsStoreTemplateID(stores.cmsStoreTemplateID)
	cfg.SetCustomStoreUsed(stores.customStoreUsed)
	cfg.SetEntityStoreUsed(stores.entityStoreUsed)
	cfg.SetFeedStoreUsed(stores.feedStoreUsed)
	cfg.SetGeoStoreUsed(stores.geoStoreUsed)
	cfg.SetLogStoreUsed(stores.logStoreUsed)
	cfg.SetMetaStoreUsed(stores.metaStoreUsed)
	cfg.SetSessionStoreUsed(stores.sessionStoreUsed)
	cfg.SetSettingStoreUsed(stores.settingStoreUsed)
	cfg.SetShopStoreUsed(stores.shopStoreUsed)
	cfg.SetSqlFileStoreUsed(stores.sqlFileStoreUsed)
	cfg.SetStatsStoreUsed(stores.statsStoreUsed)
	cfg.SetSubscriptionStoreUsed(stores.subscriptionStoreUsed)
	cfg.SetTaskStoreUsed(stores.taskStoreUsed)
	cfg.SetUserStoreUsed(stores.userStoreUsed)
	cfg.SetUserStoreVaultEnabled(stores.userStoreVaultEnabled)
	cfg.SetVaultStoreUsed(stores.vaultStoreUsed)
	cfg.SetVaultStoreKey(stores.vaultStoreKey)

	cfg.SetStripeKeyPrivate(stripe.keyPrivate)
	cfg.SetStripeKeyPublic(stripe.keyPublic)
	cfg.SetStripeUsed(stripe.used)

	cfg.SetAnthropicApiUsed(llms.anthropicUsed)
	cfg.SetAnthropicApiKey(llms.anthropicKey)
	cfg.SetAnthropicApiDefaultModel(llms.anthropicDefaultModel)

	cfg.SetGoogleGeminiApiUsed(llms.googleGeminiUsed)
	cfg.SetGoogleGeminiApiKey(llms.googleGeminiKey)
	cfg.SetGoogleGeminiApiDefaultModel(llms.googleGeminiDefaultModel)

	cfg.SetOpenAiApiUsed(llms.openAiUsed)
	cfg.SetOpenAiApiKey(llms.openAiKey)
	cfg.SetOpenAiApiDefaultModel(llms.openAiDefaultModel)

	cfg.SetOpenRouterApiUsed(llms.openRouterUsed)
	cfg.SetOpenRouterApiKey(llms.openRouterKey)
	cfg.SetOpenRouterApiDefaultModel(llms.openRouterDefaultModel)

	cfg.SetVertexAiApiUsed(llms.vertexAiUsed)
	cfg.SetVertexAiApiModelID(llms.vertexAiModelID)
	cfg.SetVertexAiApiProjectID(llms.vertexAiProjectID)
	cfg.SetVertexAiApiRegionID(llms.vertexAiRegionID)
	cfg.SetVertexAiApiDefaultModel(llms.vertexAiDefaultModel)

	cfg.SetTranslationLanguageDefault(trans.defaultLanguage)
	cfg.SetTranslationLanguageList(trans.languageList)

	return cfg, nil
}
