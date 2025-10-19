package config

import (
	"fmt"
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
	googleGeminiApiDefaultModel := env.GetString(KEY_GEMINI_API_DEFAULT_MODEL)

	mailDriver := env.GetString(KEY_MAIL_DRIVER)
	mailFromAddress := env.GetString(KEY_MAIL_FROM_ADDRESS)
	mailFromName := env.GetString(KEY_MAIL_FROM_NAME)
	mailHost := env.GetString(KEY_MAIL_HOST)
	mailPassword := env.GetString(KEY_MAIL_PASSWORD)
	mailPort := env.GetString(KEY_MAIL_PORT)
	mailUsername := env.GetString(KEY_MAIL_USERNAME)

	registrationEnabled := env.GetBool(KEY_AUTH_REGISTRATION_ENABLED)

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
		return nil, fmt.Errorf("%s requires %s to be true", KEY_USER_STORE_USE_VAULT, KEY_VAULT_STORE_USED)
	}

	// mediaBucket := env.Value("MEDIA_BUCKET")
	// mediaDriver := env.Value("MEDIA_DRIVER")
	// mediaEndpoint := env.Value("MEDIA_ENDPOINT")
	// mediaKey := env.Value("MEDIA_KEY")
	// mediaRoot := env.Value("MEDIA_ROOT")
	// mediaSecret := env.Value("MEDIA_SECRET")
	// mediaRegion := env.Value("MEDIA_REGION")
	// mediaUrl := env.Value("MEDIA_URL")

	// Payments
	stripeKeyPrivate := env.GetString(KEY_STRIPE_KEY_PRIVATE)
	stripeKeyPublic := env.GetString(KEY_STRIPE_KEY_PUBLIC)
	stripeUsed := stripeKeyPrivate != "" && stripeKeyPublic != ""

	// LLM: Anthropic
	anthropicApiUsed := env.GetBool(KEY_ANTHROPIC_API_USED)
	anthropicApiKey := env.GetString(KEY_ANTHROPIC_API_KEY)
	anthropicApiDefaultModel := env.GetString(KEY_ANTHROPIC_API_DEFAULT_MODEL)

	// LLM: OpenAI
	openAiApiUsed := env.GetBool(KEY_OPENAI_API_USED)
	openAiApiKey := env.GetString(KEY_OPENAI_API_KEY)
	openAiApiDefaultModel := env.GetString(KEY_OPENAI_API_DEFAULT_MODEL)

	// LLM: OpenRouter
	openRouterApiUsed := env.GetBool(KEY_OPENROUTER_API_USED)
	openRouterApiKey := env.GetString(KEY_OPENROUTER_API_KEY)
	openRouterApiDefaultModel := env.GetString(KEY_OPENROUTER_API_DEFAULT_MODEL)

	// LLM: Vertex AI
	vertexAiApiUsed := env.GetBool(KEY_VERTEX_AI_API_USED)
	vertexAiApiModelID := env.GetString(KEY_VERTEX_AI_API_MODEL_ID)
	vertexAiApiProjectID := env.GetString(KEY_VERTEX_AI_API_PROJECT_ID)
	vertexAiApiRegionID := env.GetString(KEY_VERTEX_AI_API_REGION_ID)
	vertexAiApiDefaultModel := env.GetString(KEY_VERTEX_AI_API_DEFAULT_MODEL)

	// Check required variables

	// Enable if you use CMS template
	if err := requireWhen(cmsStoreUsed, KEY_CMS_STORE_TEMPLATE_ID, "required when `CMS_STORE_USED` is true", cmsStoreTemplateID); err != nil {
		return nil, err
	}

	// LLM: Anthropic
	if err := requireWhen(anthropicApiUsed, KEY_ANTHROPIC_API_KEY, "required when `ANTHROPIC_API_USED` is true", anthropicApiKey); err != nil {
		return nil, err
	}
	if err := requireWhen(anthropicApiUsed, KEY_ANTHROPIC_API_DEFAULT_MODEL, "required when `ANTHROPIC_API_USED` is true", anthropicApiDefaultModel); err != nil {
		return nil, err
	}

	// LLM: Google Gemini
	if err := requireWhen(googleGeminiApiUsed, KEY_GEMINI_API_KEY, "required when `GEMINI_API_USED` is true", googleGeminiApiKey); err != nil {
		return nil, err
	}
	if err := requireWhen(googleGeminiApiUsed, KEY_GEMINI_API_DEFAULT_MODEL, "required when `GEMINI_API_USED` is true", googleGeminiApiDefaultModel); err != nil {
		return nil, err
	}

	// LLM: OpenAI
	if err := requireWhen(openAiApiUsed, KEY_OPENAI_API_KEY, "required when `OPENAI_API_USED` is true", openAiApiKey); err != nil {
		return nil, err
	}
	if err := requireWhen(openAiApiUsed, KEY_OPENAI_API_DEFAULT_MODEL, "required when `OPENAI_API_USED` is true", openAiApiDefaultModel); err != nil {
		return nil, err
	}

	// LLM: OpenRouter
	if err := requireWhen(openRouterApiUsed, KEY_OPENROUTER_API_KEY, "required when `OPENROUTER_API_USED` is true", openRouterApiKey); err != nil {
		return nil, err
	}

	if err := requireWhen(openRouterApiUsed, KEY_OPENROUTER_API_DEFAULT_MODEL, "required when `OPENROUTER_API_USED` is true", openRouterApiDefaultModel); err != nil {
		return nil, err
	}

	// LLM: Vertex AI
	if err := requireWhen(vertexAiApiUsed, KEY_VERTEX_AI_API_MODEL_ID, "required when `VERTEX_AI_API_USED` is true", vertexAiApiModelID); err != nil {
		return nil, err
	}
	if err := requireWhen(vertexAiApiUsed, KEY_VERTEX_AI_API_PROJECT_ID, "required when `VERTEX_AI_API_USED` is true", vertexAiApiProjectID); err != nil {
		return nil, err
	}
	if err := requireWhen(vertexAiApiUsed, KEY_VERTEX_AI_API_REGION_ID, "required when `VERTEX_AI_API_USED` is true", vertexAiApiRegionID); err != nil {
		return nil, err
	}
	if err := requireWhen(vertexAiApiUsed, KEY_VERTEX_AI_API_DEFAULT_MODEL, "required when `VERTEX_AI_API_USED` is true", vertexAiApiDefaultModel); err != nil {
		return nil, err
	}

	// os.Setenv("TZ", "UTC")

	// err = initializeDatabase()

	// if err != nil {
	// 	return nil, err
	// }

	config := types.Config{}

	// Store configurations
	config.SetAuditStoreUsed(auditStoreUsed)
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
	config.SetSubscriptionStoreUsed(subscriptionStoreUsed)
	config.SetTaskStoreUsed(taskStoreUsed)
	config.SetUserStoreUsed(userStoreUsed)
	config.SetUserStoreVaultEnabled(userStoreVaultEnabled)
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
	if openRouterApiDefaultModel != "" {
		config.SetOpenRouterApiDefaultModel(openRouterApiDefaultModel)
	}

	// OpenAI
	config.SetOpenAiApiUsed(openAiApiUsed)
	config.SetOpenAiApiKey(openAiApiKey)
	if openAiApiDefaultModel != "" {
		config.SetOpenAiApiDefaultModel(openAiApiDefaultModel)
	}

	// Anthropic
	config.SetAnthropicApiUsed(anthropicApiUsed)
	config.SetAnthropicApiKey(anthropicApiKey)
	if anthropicApiDefaultModel != "" {
		config.SetAnthropicApiDefaultModel(anthropicApiDefaultModel)
	}

	// Google Gemini
	config.SetGoogleGeminiApiUsed(googleGeminiApiUsed)
	config.SetGoogleGeminiApiKey(googleGeminiApiKey)
	if googleGeminiApiDefaultModel != "" {
		config.SetGoogleGeminiApiDefaultModel(googleGeminiApiDefaultModel)
	}

	// Vertex AI
	config.SetVertexAiApiUsed(vertexAiApiUsed)
	if vertexAiApiDefaultModel != "" {
		config.SetVertexAiApiDefaultModel(vertexAiApiDefaultModel)
	}
	if vertexAiApiProjectID != "" {
		config.SetVertexAiApiProjectID(vertexAiApiProjectID)
	}
	if vertexAiApiRegionID != "" {
		config.SetVertexAiApiRegionID(vertexAiApiRegionID)
	}
	if vertexAiApiModelID != "" {
		config.SetVertexAiApiModelID(vertexAiApiModelID)
	}

	// Stripe
	config.SetStripeUsed(stripeUsed)
	config.SetStripeKeyPrivate(stripeKeyPrivate)
	config.SetStripeKeyPublic(stripeKeyPublic)

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
