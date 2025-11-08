package config

import (
	"os"
	"project/internal/resources"
	"project/internal/types"
	"strings"

	"github.com/dracory/env"
	"github.com/dracory/envenc"
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

func Load() (types.ConfigInterface, error) {
	env.Load(".env")

	acc := &loadAccumulator{}

	app := loadAppConfig(acc)
	envEnc := loadEnvEncryptionConfig(acc)
	db := loadDatabaseConfig(acc)
	mail := loadMailConfig()
	reg := loadRegistrationConfig()
	stores := loadStoresConfig(acc)
	stripe := loadStripeConfig()
	llms := loadLLMConfig(acc)
	trans := loadTranslationConfig()

	if envEnc.used {
		if err := intializeEnvEncVariables(app.env); err != nil {
			acc.add(err)
		}
	}

	if err := acc.err(); err != nil {
		return nil, err
	}

	cfg := &types.Config{}

	cfg.SetAppName(app.name)
	cfg.SetAppUrl(app.url)
	cfg.SetAppHost(app.host)
	cfg.SetAppPort(app.port)
	cfg.SetAppEnv(app.env)
	cfg.SetAppDebug(app.debug)

	if envEnc.key != "" {
		cfg.SetEnvEncryptionKey(envEnc.key)
	}

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

// initializeEnvEncVariables initializes the envenc variables
// based on the app environment
//
// Business logic:
//   - check if the app environment is testing, skipped as not needed
//   - requires the ENV_ENCRYPTION_KEY env variable
//   - looks for file the file name is .env.<app_environment>.vault
//     both in the local file system and in the resources folder
//   - if none found, it will panic
//   - if it fails for other reasons, it will panic
//
// Parameters:
// - appEnvironment: the app environment
//
// Returns:
// - none
func intializeEnvEncVariables(appEnvironment string) error {
	if appEnvironment == APP_ENVIRONMENT_TESTING {
		return nil
	}

	if strings.TrimSpace(appEnvironment) == "" {
		return MissingEnvError{Key: KEY_APP_ENVIRONMENT, Context: "required to initialize EnvEnc variables"}
	}

	appEnvironment = strings.ToLower(appEnvironment)
	envEncryptionKey := env.GetString(KEY_ENVENC_KEY_PRIVATE)

	if err := ensureRequired(envEncryptionKey, KEY_ENVENC_KEY_PRIVATE, "required to hydrate EnvEnc variables"); err != nil {
		return err
	}

	vaultFilePath := ".env." + appEnvironment + ".vault"

	vaultContent, err := resources.Resource(".env." + appEnvironment + ".vault")

	if err != nil {
		return err
	}

	derivedEnvEncKey, err := deriveEnvEncKey(envEncryptionKey)

	if err != nil {
		return err
	}

	if fileExists(vaultFilePath) {
		err := envenc.HydrateEnvFromFile(vaultFilePath, derivedEnvEncKey)

		if err != nil {
			return err
		}
	}

	if vaultContent != "" {
		err = envenc.HydrateEnvFromString(vaultContent, derivedEnvEncKey)

		if err != nil {
			return err
		}
	}

	return nil
}

// fileExists checks if a file exists at the given path.
func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return !os.IsNotExist(err)
	}
}
