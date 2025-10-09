package config

import (
	"os"
	"project/internal/resources"
	"project/internal/types"
	"strings"

	"github.com/dracory/env"
	"github.com/dracory/envenc"
	"github.com/samber/lo"
	"github.com/spf13/cast"
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

	appEnvironment, err := requireString(KEY_APP_ENVIRONMENT, "set the application environment")
	if err != nil {
		return nil, err
	}

	// Enable if you use envenc
	// if err := intializeEnvEncVariables(appEnvironment); err != nil {
	// 	return nil, err
	// }

	// == EnvEnc: derive encryption key if private key is provided ==
	envencRealKey := lo.
		IfF(env.GetString(KEY_ENVENC_KEY_PRIVATE) != "", func() string {
			privateKey := strings.TrimSpace(env.GetString(KEY_ENVENC_KEY_PRIVATE))
			realKey, err := deriveEnvEncKey(privateKey)
			if err != nil {
				return ""
			}
			return realKey
		}).
		Else("")

	appName := env.GetString(KEY_APP_NAME)
	appUrl := env.GetString(KEY_APP_URL)
	appHost, err := requireString(KEY_APP_HOST, "set the application host address")
	if err != nil {
		return nil, err
	}
	appPort, err := requireString(KEY_APP_PORT, "set the application port")
	if err != nil {
		return nil, err
	}
	appDebug := env.GetBool(KEY_APP_DEBUG)
	dbDriver, err := requireString(KEY_DB_DRIVER, "select the database driver (e.g., sqlite, postgres)")
	if err != nil {
		return nil, err
	}
	dbHost := strings.TrimSpace(env.GetString(KEY_DB_HOST))
	if err := requireWhen(dbDriver != driverSQLite, KEY_DB_HOST, "required when `DB_DRIVER` is not sqlite", dbHost); err != nil {
		return nil, err
	}
	dbPort := strings.TrimSpace(env.GetString(KEY_DB_PORT))
	if err := requireWhen(dbDriver != driverSQLite, KEY_DB_PORT, "required when `DB_DRIVER` is not sqlite", dbPort); err != nil {
		return nil, err
	}
	dbName, err := requireString(KEY_DB_DATABASE, "set the database name")
	if err != nil {
		return nil, err
	}
	dbUser := strings.TrimSpace(env.GetString(KEY_DB_USERNAME))
	if err := requireWhen(dbDriver != driverSQLite, KEY_DB_USERNAME, "required when `DB_DRIVER` is not sqlite", dbUser); err != nil {
		return nil, err
	}
	dbPass := strings.TrimSpace(env.GetString(KEY_DB_PASSWORD))
	if err := requireWhen(dbDriver != driverSQLite, KEY_DB_PASSWORD, "required when `DB_DRIVER` is not sqlite", dbPass); err != nil {
		return nil, err
	}

	// LLM: Google Gemini
	googleGeminiApiUsed := env.GetBool(KEY_GEMINI_API_USED)
	googleGeminiApiKey := env.GetString(KEY_GEMINI_API_KEY)
	googleGeminiDefaultModel := env.GetString(KEY_GEMINI_DEFAULT_MODEL)

	mailDriver := env.GetString(KEY_MAIL_DRIVER)
	mailFromAddress := env.GetString(KEY_MAIL_FROM_ADDRESS)
	mailFromName := env.GetString(KEY_MAIL_FROM_NAME)
	mailHost := env.GetString(KEY_MAIL_HOST)
	mailPassword := env.GetString(KEY_MAIL_PASSWORD)
	mailPort := env.GetString(KEY_MAIL_PORT)
	mailUsername := env.GetString(KEY_MAIL_USERNAME)

	registrationEnabled := env.GetBool(KEY_AUTH_REGISTRATION_ENABLED)

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
	taskStoreUsed := env.GetBool(KEY_TASK_STORE_USED)
	userStoreUsed := env.GetBool(KEY_USER_STORE_USED)
	vaultStoreUsed := env.GetBool(KEY_VAULT_STORE_USED)
	vaultStoreKey := env.GetString(KEY_VAULT_STORE_KEY)

	// mediaBucket := env.Value("MEDIA_BUCKET")
	// mediaDriver := env.Value("MEDIA_DRIVER")
	// mediaEndpoint := env.Value("MEDIA_ENDPOINT")
	// mediaKey := env.Value("MEDIA_KEY")
	// mediaRoot := env.Value("MEDIA_ROOT")
	// mediaSecret := env.Value("MEDIA_SECRET")
	// mediaRegion := env.Value("MEDIA_REGION")
	// mediaUrl := env.Value("MEDIA_URL")
	// LLM: OpenAI
	openAiApiUsed := env.GetBool(KEY_OPENAI_API_USED)
	openAiApiKey := env.GetString(KEY_OPENAI_API_KEY)
	openAiDefaultModel := env.GetString(KEY_OPENAI_DEFAULT_MODEL)

	// LLM: OpenRouter
	openRouterApiUsed := env.GetBool(KEY_OPENROUTER_API_USED)
	openRouterApiKey := env.GetString(KEY_OPENROUTER_API_KEY)
	openRouterDefaultModel := env.GetString(KEY_OPENROUTER_DEFAULT_MODEL)

	stripeUsed := false
	stripeKeyPrivate := env.GetString(KEY_STRIPE_KEY_PRIVATE)
	stripeKeyPublic := env.GetString(KEY_STRIPE_KEY_PUBLIC)

	// LLM: Vertex AI
	vertexAiUsed := env.GetBool(KEY_VERTEX_AI_USED)
	vertexAiModelID := env.GetString(KEY_VERTEX_MODEL_ID)
	vertexAiProjectID := env.GetString(KEY_VERTEX_PROJECT_ID)
	vertexAiRegionID := env.GetString(KEY_VERTEX_REGION_ID)
	vertexAiDefaultModel := env.GetString(KEY_VERTEX_DEFAULT_MODEL)

	// LLM: Anthropic
	anthropicApiUsed := env.GetBool(KEY_ANTHROPIC_API_USED)
	anthropicApiKey := env.GetString(KEY_ANTHROPIC_API_KEY)
	anthropicDefaultModel := env.GetString(KEY_ANTHROPIC_DEFAULT_MODEL)

	// Check required variables

	// Enable if you use CMS template
	if err := requireWhen(cmsStoreUsed, KEY_CMS_STORE_TEMPLATE_ID, "required when `CMS_STORE_USED` is true", cmsStoreTemplateID); err != nil {
		return nil, err
	}

	if err := requireWhen(googleGeminiApiUsed, KEY_GEMINI_API_KEY, "required when `GEMINI_API_USED` is true", googleGeminiApiKey); err != nil {
		return nil, err
	}

	if err := requireWhen(openAiApiUsed, KEY_OPENAI_API_KEY, "required when `OPENAI_API_USED` is true", openAiApiKey); err != nil {
		return nil, err
	}

	if err := requireWhen(openRouterApiUsed, KEY_OPENROUTER_API_KEY, "required when `OPENROUTER_API_USED` is true", openRouterApiKey); err != nil {
		return nil, err
	}

	if err := requireWhen(anthropicApiUsed, KEY_ANTHROPIC_API_KEY, "required when `ANTHROPIC_API_USED` is true", anthropicApiKey); err != nil {
		return nil, err
	}

	if err := requireWhen(stripeUsed, KEY_STRIPE_KEY_PRIVATE, "required when Stripe integration is enabled", stripeKeyPrivate); err != nil {
		return nil, err
	}

	if err := requireWhen(stripeUsed, KEY_STRIPE_KEY_PUBLIC, "required when Stripe integration is enabled", stripeKeyPublic); err != nil {
		return nil, err
	}

	if err := requireWhen(vaultStoreUsed, KEY_VAULT_STORE_KEY, "required when `VAULT_STORE_USED` is true", vaultStoreKey); err != nil {
		return nil, err
	}

	if err := requireWhen(vertexAiUsed, KEY_VERTEX_MODEL_ID, "required when `VERTEX_AI_USED` is true", vertexAiModelID); err != nil {
		return nil, err
	}
	if err := requireWhen(vertexAiUsed, KEY_VERTEX_PROJECT_ID, "required when `VERTEX_AI_USED` is true", vertexAiProjectID); err != nil {
		return nil, err
	}
	if err := requireWhen(vertexAiUsed, KEY_VERTEX_REGION_ID, "required when `VERTEX_AI_USED` is true", vertexAiRegionID); err != nil {
		return nil, err
	}

	// os.Setenv("TZ", "UTC")

	// err = initializeDatabase()

	// if err != nil {
	// 	return nil, err
	// }

	config := types.Config{}

	// Store configurations
	config.SetBlogStoreUsed(blogStoreUsed)
	config.SetCacheStoreUsed(cacheStoreUsed)
	config.SetCmsStoreUsed(cmsStoreUsed)
	config.SetCmsStoreTemplateID(cmsStoreTemplateID)
	config.SetCustomStoreUsed(customStoreUsed)
	config.SetEntityStoreUsed(entityStoreUsed)
	config.SetFeedStoreUsed(feedStoreUsed)
	config.SetGeoStoreUsed(geoStoreUsed)
	config.SetLogStoreUsed(logStoreUsed)
	config.SetMetaStoreUsed(metaStoreUsed)
	config.SetSessionStoreUsed(sessionStoreUsed)
	config.SetSettingStoreUsed(settingStoreUsed)
	config.SetShopStoreUsed(shopStoreUsed)
	config.SetSqlFileStoreUsed(sqlFileStoreUsed)
	config.SetStatsStoreUsed(statsStoreUsed)
	config.SetTaskStoreUsed(taskStoreUsed)
	config.SetUserStoreUsed(userStoreUsed)
	config.SetVaultStoreUsed(vaultStoreUsed)
	config.SetVaultStoreKey(vaultStoreKey)

	// App configurations
	config.SetAppDebug(appDebug)
	config.SetAppName(appName)
	config.SetAppEnv(appEnvironment)
	config.SetAppHost(appHost)
	config.SetAppPort(appPort)
	config.SetAppUrl(appUrl)

	// Apply EnvEnc key if derived
	if envencRealKey != "" {
		config.SetEnvEncryptionKey(envencRealKey)
	}

	// Mail configurations
	config.SetMailDriver(mailDriver)
	config.SetMailHost(mailHost)
	config.SetMailPort(cast.ToInt(mailPort))
	config.SetMailUsername(mailUsername)
	config.SetMailPassword(mailPassword)
	config.SetMailFromAddress(mailFromAddress)
	config.SetMailFromName(mailFromName)
	config.SetRegistrationEnabled(registrationEnabled)

	// Database configurations
	config.SetDatabaseDriver(dbDriver)
	config.SetDatabaseHost(dbHost)
	config.SetDatabasePort(dbPort)
	config.SetDatabaseName(dbName)
	config.SetDatabaseUsername(dbUser)
	config.SetDatabasePassword(dbPass)
	// config.SetDatabaseCharset(`utf8mb4`)
	// config.SetDatabaseTimeZone("UTC")
	config.SetDatabaseSSLMode("require")

	// i18n defaults (can be overridden via env in future if needed)
	translationDefault := env.GetString(KEY_TRANSLATION_LANGUAGE_DEFAULT)
	if translationDefault == "" {
		translationDefault = translationLanguageDefault()
	}
	translationList := translationLanguageListDefault()
	config.SetTranslationLanguageDefault(translationDefault)
	config.SetTranslationLanguageList(translationList)

	// == LLM configuration ==
	// OpenRouter
	config.SetOpenRouterApiUsed(openRouterApiUsed)
	config.SetOpenRouterApiKey(openRouterApiKey)
	if openRouterDefaultModel != "" {
		config.SetOpenRouterDefaultModel(openRouterDefaultModel)
	}

	// OpenAI
	config.SetOpenAiApiUsed(openAiApiUsed)
	config.SetOpenAiApiKey(openAiApiKey)
	if openAiDefaultModel != "" {
		config.SetOpenAiDefaultModel(openAiDefaultModel)
	}

	// Anthropic
	config.SetAnthropicApiUsed(anthropicApiUsed)
	config.SetAnthropicApiKey(anthropicApiKey)
	if anthropicDefaultModel != "" {
		config.SetAnthropicDefaultModel(anthropicDefaultModel)
	}

	// Google Gemini
	config.SetGoogleGeminiApiUsed(googleGeminiApiUsed)
	config.SetGoogleGeminiApiKey(googleGeminiApiKey)
	if googleGeminiDefaultModel != "" {
		config.SetGoogleGeminiDefaultModel(googleGeminiDefaultModel)
	}

	// Vertex AI
	config.SetVertexAiUsed(vertexAiUsed)
	if vertexAiDefaultModel != "" {
		config.SetVertexAiDefaultModel(vertexAiDefaultModel)
	}
	if vertexAiProjectID != "" {
		config.SetVertexAiProjectID(vertexAiProjectID)
	}
	if vertexAiRegionID != "" {
		config.SetVertexAiRegionID(vertexAiRegionID)
	}
	if vertexAiModelID != "" {
		config.SetVertexAiModelID(vertexAiModelID)
	}

	return &config, nil
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
	envEncryptionKey := env.GetString(KEY_ENV_ENCRYPTION_KEY)

	if err := ensureRequired(envEncryptionKey, KEY_ENV_ENCRYPTION_KEY, "required to hydrate EnvEnc variables"); err != nil {
		return err
	}

	vaultFilePath := ".env." + appEnvironment + ".vault"

	vaultContent, err := resources.Resource(".env." + appEnvironment + ".vault")

	if err != nil {
		panic(err.Error())
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
