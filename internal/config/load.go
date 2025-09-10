package config

import (
	"errors"
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

	appEnvironment, err := env.GetStringOrError(KEY_APP_ENVIRONMENT)
	if err != nil {
		return nil, errors.New(KEY_APP_ENVIRONMENT + " is required")
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
	appHost, err := env.GetStringOrError(KEY_APP_HOST)
	if err != nil {
		return nil, errors.New(KEY_APP_HOST + " is required")
	}
	appPort, err := env.GetStringOrError(KEY_APP_PORT)
	if err != nil {
		return nil, errors.New(KEY_APP_PORT + " is required")
	}
	appDebug := env.GetBool(KEY_APP_DEBUG)

	// cmsUserTemplateID := env.Value("CMS_TEMPLATE_ID")

	dbDriver, err := env.GetStringOrError(KEY_DB_DRIVER)
	if err != nil {
		return nil, errors.New(KEY_DB_DRIVER + " is required")
	}
	dbHost, err := env.GetStringOrError(KEY_DB_HOST)
	if dbDriver != "sqlite" && err != nil {
		return nil, errors.New(KEY_DB_HOST + " is required")
	}
	dbPort, err := env.GetStringOrError(KEY_DB_PORT)
	if dbDriver != "sqlite" && err != nil {
		return nil, errors.New(KEY_DB_PORT + " is required")
	}
	dbName, err := env.GetStringOrError(KEY_DB_DATABASE)
	if err != nil {
		return nil, errors.New(KEY_DB_DATABASE + " is required")
	}
	dbUser, err := env.GetStringOrError(KEY_DB_USERNAME)
	if dbDriver != "sqlite" && err != nil {
		return nil, errors.New(KEY_DB_USERNAME + " is required")
	}
	dbPass, err := env.GetStringOrError(KEY_DB_PASSWORD)
	if dbDriver != "sqlite" && err != nil {
		return nil, errors.New(KEY_DB_PASSWORD + " is required")
	}

	// LLM: Google Gemini
	googleGeminiApiUsed := env.GetBool(KEY_GEMINI_API_USED)
	googleGeminiApiKey := env.GetString(KEY_GEMINI_API_KEY)
	googleGeminiDefaultModel := env.GetString(KEY_GEMINI_DEFAULT_MODEL)

	mailDriver := env.GetString(KEY_MAIL_DRIVER)
	mailFromEmailAddress := env.GetString(KEY_EMAIL_FROM_ADDRESS)
	mailFromName := env.GetString(KEY_EMAIL_FROM_NAME)
	mailHost := env.GetString(KEY_MAIL_HOST)
	mailPassword := env.GetString(KEY_MAIL_PASSWORD)
	mailPort := env.GetString(KEY_MAIL_PORT)
	mailUsername := env.GetString(KEY_MAIL_USERNAME)

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

	vaultKey := env.GetString(KEY_VAULT_KEY)
	vaultStoreUsed := env.GetBool(KEY_VAULT_STORE_USED)

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

	// Daily Analysis
	dailySymbolsCSV := env.GetString(KEY_DAILY_ANALYSIS_SYMBOLS)
	dailyTimeUTC := env.GetString(KEY_DAILY_ANALYSIS_TIME_UTC)
	dailyCadenceHours := env.GetString(KEY_DAILY_ANALYSIS_CADENCE_HOURS)

	// Check required variables

	// Enable if you use CMS template
	// if cmsUserTemplateID == "" {
	// 	return nil, errors.New("CMS_TEMPLATE_ID is required")
	// }

	if googleGeminiApiUsed && googleGeminiApiKey == "" {
		return nil, errors.New(KEY_GEMINI_API_KEY + " is required")
	}

	if openAiApiUsed && openAiApiKey == "" {
		return nil, errors.New(KEY_OPENAI_API_KEY + " is required")
	}

	if openRouterApiUsed && openRouterApiKey == "" {
		return nil, errors.New(KEY_OPENROUTER_API_KEY + " is required")
	}

	if anthropicApiUsed && anthropicApiKey == "" {
		return nil, errors.New(KEY_ANTHROPIC_API_KEY + " is required")
	}

	if stripeUsed && stripeKeyPrivate == "" {
		return nil, errors.New(KEY_STRIPE_KEY_PRIVATE + " is required")
	}

	if stripeUsed && stripeKeyPublic == "" {
		return nil, errors.New(KEY_STRIPE_KEY_PUBLIC + " is required")
	}

	if vaultStoreUsed && vaultKey == "" {
		return nil, errors.New(KEY_VAULT_KEY + " is required")
	}

	if vertexAiUsed && vertexAiModelID == "" {
		return nil, errors.New(KEY_VERTEX_MODEL_ID + " is required")
	}
	if vertexAiUsed && vertexAiProjectID == "" {
		return nil, errors.New(KEY_VERTEX_PROJECT_ID + " is required")
	}
	if vertexAiUsed && vertexAiRegionID == "" {
		return nil, errors.New(KEY_VERTEX_REGION_ID + " is required")
	}

	// os.Setenv("TZ", "UTC")

	// err = initializeDatabase()

	// if err != nil {
	// 	return nil, err
	// }

	config := types.Config{}

	config.SetAppDebug(appDebug)
	config.SetAppName(appName)
	// config.SetAppType(AppType)
	config.SetAppEnv(appEnvironment)
	config.SetAppHost(appHost)
	config.SetAppPort(appPort)
	config.SetAppUrl(appUrl)

	// Apply EnvEnc key if derived
	if envencRealKey != "" {
		config.SetEnvEncryptionKey(envencRealKey)
	}

	config.SetMailDriver(mailDriver)
	config.SetMailHost(mailHost)
	config.SetMailPort(cast.ToInt(mailPort))
	config.SetMailUsername(mailUsername)
	config.SetMailPassword(mailPassword)
	config.SetMailFromEmail(mailFromEmailAddress)
	config.SetMailFromName(mailFromName)

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
		translationDefault = "en"
	}
	translationList := map[string]string{
		"en": "English",
		"bg": "Bulgarian",
		"de": "German",
	}
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

	// == Daily Analysis defaults and config ==
	if strings.TrimSpace(dailySymbolsCSV) == "" {
		dailySymbolsCSV = "US30,US2000,US500,SP500,NASDAQ,NL25,UK100,DAX,FR40,SMI,NIKKEI225,XAGUSD,XAUUSD,XPTUSD,XPDUSD"
	}
	if strings.TrimSpace(dailyTimeUTC) == "" {
		dailyTimeUTC = "06:30"
	}
	cadence := cast.ToInt(dailyCadenceHours)
	if cadence == 0 {
		cadence = 24
	}
	// parse symbols
	var symList []string
	for _, s := range strings.Split(dailySymbolsCSV, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			symList = append(symList, s)
		}
	}
	config.SetDailyAnalysisSymbols(symList)
	config.SetDailyAnalysisTimeUTC(dailyTimeUTC)
	config.SetDailyAnalysisCadenceHours(cadence)

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

	if appEnvironment == "" {
		return errors.New(KEY_APP_ENVIRONMENT + " is required")
	}

	appEnvironment = strings.ToLower(appEnvironment)
	envEncryptionKey := env.GetString(KEY_ENV_ENCRYPTION_KEY)

	if envEncryptionKey == "" {
		return errors.New(KEY_ENV_ENCRYPTION_KEY + " is required")
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
