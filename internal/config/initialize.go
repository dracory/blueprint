package config

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"project/internal/resources"
	"strings"

	"github.com/faabiosr/cachego/file"

	"github.com/dracory/database"
	"github.com/dracory/env"
	"github.com/gouniverse/logstore"
	"github.com/gouniverse/sb"
	"github.com/gouniverse/utils"
	"github.com/jellydator/ttlcache/v3"
	"github.com/samber/lo"

	// _ "github.com/go-sql-driver/mysql" // Enable MySQL driver if needed
	// _ "github.com/lib/pq" // Enable PostgreSQL driver if needed
	_ "modernc.org/sqlite" // Enable SQLite driver if needed
)

// Initialize initializes the application
//
// Business logic:
//   - initializes the environment variables
//   - initializes the database
//   - migrates the database
//   - initializes the in-memory cache
//   - initializes the logger
//
// Parameters:
// - none
//
// Returns:
// - none
func Initialize() error {
	err := initializeEnvVariables()

	if err != nil {
		return err
	}

	os.Setenv("TZ", "UTC")

	err = initializeDatabase()

	if err != nil {
		return err
	}

	err = migrateDatabase()

	if err != nil {
		return err
	}

	initializeCache()

	Logger = *slog.New(logstore.NewSlogHandler(LogStore))

	return nil
}

// initializeEnvVariables initializes the env variables
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
func initializeEnvVariables() error {
	env.Initialize(".env")

	AppEnvironment = env.Value("APP_ENV")

	// Enable if you use envenc
	// if err := intializeEnvEncVariables(AppEnvironment); err != nil {
	// 	return err
	// }

	AppName = env.Value("APP_NAME")
	AppUrl = env.Value("APP_URL")
	CmsUserTemplateID = env.Value("CMS_TEMPLATE_ID")
	DbDriver = env.Value("DB_DRIVER")
	DbHost = env.Value("DB_HOST")
	DbPort = env.Value("DB_PORT")
	DbName = env.Value("DB_DATABASE")
	DbUser = env.Value("DB_USERNAME")
	DbPass = env.Value("DB_PASSWORD")
	Debug = env.Bool("DEBUG")
	GoogleGeminiApiKey = env.Value("GEMINI_API_KEY")
	MailDriver = env.Value("MAIL_DRIVER")
	MailFromEmailAddress = env.Value("EMAIL_FROM_ADDRESS")
	MailFromName = env.Value("EMAIL_FROM_NAME")
	MailHost = env.Value("MAIL_HOST")
	MailPassword = env.Value("MAIL_PASSWORD")
	MailPort = env.Value("MAIL_PORT")
	MailUsername = env.Value("MAIL_USERNAME")
	MediaBucket = env.Value("MEDIA_BUCKET")
	MediaDriver = env.Value("MEDIA_DRIVER")
	MediaEndpoint = env.Value("MEDIA_ENDPOINT")
	MediaKey = env.Value("MEDIA_KEY")
	MediaRoot = env.Value("MEDIA_ROOT")
	MediaSecret = env.Value("MEDIA_SECRET")
	MediaRegion = env.Value("MEDIA_REGION")
	MediaUrl = env.Value("MEDIA_URL")
	OpenAiApiKey = env.Value("OPENAI_API_KEY")
	StripeKeyPrivate = env.Value("STRIPE_KEY_PRIVATE")
	StripeKeyPublic = env.Value("STRIPE_KEY_PUBLIC")
	VaultKey = env.Value("VAULT_KEY")
	VertexAiModelID = env.Value("VERTEX_MODEL_ID")
	VertexAiProjectID = env.Value("VERTEX_PROJECT_ID")
	VertexAiRegionID = env.Value("VERTEX_REGION_ID")
	WebServerHost = env.Value("SERVER_HOST")
	WebServerPort = env.Value("SERVER_PORT")

	// Check required variables

	if AppEnvironment == "" {
		return errors.New("APP_ENV is required")
	}

	// Enable if you use CMS template
	// if CmsUserTemplateID == "" {
	// 	return errors.New("CMS_TEMPLATE_ID is required")
	// }

	if DbDriver == "" {
		return errors.New("DB_DRIVER is required")
	}

	if DbDriver != "sqlite" && DbHost == "" {
		return errors.New("DB_HOST is required")
	}

	if DbDriver != "sqlite" && DbPort == "" {
		return errors.New("DB_PORT is required")
	}

	if DbName == "" {
		return errors.New("DB_DATABASE is required")
	}

	if DbDriver != "sqlite" && DbUser == "" {
		return errors.New("DB_USERNAME is required")
	}

	if DbDriver != "sqlite" && DbPass == "" {
		return errors.New("DB_PASSWORD is required")
	}

	if GoogleGeminiApiUsed && GoogleGeminiApiKey == "" {
		return errors.New("GEMINI_API_KEY is required")
	}

	if OpenAiApiUsed && OpenAiApiKey == "" {
		return errors.New("OPENAI_API_KEY is required")
	}

	if StripeUsed && StripeKeyPrivate == "" {
		return errors.New("STRIPE_KEY_PRIVATE is required")
	}

	if StripeUsed && StripeKeyPublic == "" {
		return errors.New("STRIPE_KEY_PUBLIC is required")
	}

	if VaultStoreUsed && VaultKey == "" {
		return errors.New("VAULT_KEY is required")
	}

	if VertexAiUsed && VertexAiModelID == "" {
		return errors.New("VERTEX_MODEL_ID is required")
	}
	if VertexAiUsed && VertexAiProjectID == "" {
		return errors.New("VERTEX_PROJECT_ID is required")
	}
	if VertexAiUsed && VertexAiRegionID == "" {
		return errors.New("VERTEX_REGION_ID is required")
	}

	if WebServerHost == "" {
		return errors.New("SERVER_HOST is required")
	}

	if WebServerPort == "" {
		return errors.New("SERVER_PORT is required")
	}

	return nil
}

// initializeEnvEncVariables initializes the envenc variables
// based on the app environment
//
// Business logic:
//   - checkd if the app environment is testing, skipped as not needed
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
		return errors.New("APP_ENV is required")
	}

	appEnvironment = strings.ToLower(appEnvironment)
	envEncryptionKey := env.Value("ENV_ENCRYPTION_KEY")

	if envEncryptionKey == "" {
		return errors.New("ENV_ENCRYPTION_KEY is required")
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

	err = utils.EnvEncInitialize(struct {
		Password      string
		VaultFilePath string
		VaultContent  string
	}{
		Password:      derivedEnvEncKey,
		VaultFilePath: lo.Ternary(vaultContent == "", vaultFilePath, ""),
		VaultContent:  lo.Ternary(vaultContent != "", vaultContent, ""),
	})

	if err != nil {
		return err
	}

	return nil
}

// initializeCache initializes the cache
func initializeCache() {
	CacheMemory = ttlcache.New[string, any]()
	// create a new directory
	_ = os.MkdirAll(".cache", os.ModePerm)
	CacheFile = file.New(".cache")
}

// initializeDatabase initializes the database
//
// Business logic:
//   - opens the database
//   - initializes the required stores
//
// Parameters:
// - none
//
// Returns:
// - error: the error if any
func initializeDatabase() error {
	db, err := database.Open(database.Options().
		SetDatabaseType(DbDriver).
		SetDatabaseHost(DbHost).
		SetDatabasePort(DbPort).
		SetDatabaseName(DbName).
		SetCharset(`utf8mb4`).
		SetTimeZone("UTC").
		SetUserName(DbUser).
		SetPassword(DbPass).
		SetSSLMode("require"))

	if err != nil {
		return err
	}

	if db == nil {
		return errors.New("db is nil")
	}

	// Add connection pool settings
	// db.SetMaxOpenConns(25)                 // Maximum number of open connections
	// db.SetMaxIdleConns(5)                  // Maximum number of idle connections
	// db.SetConnMaxLifetime(5 * time.Minute) // Maximum connection lifetime

	dbInstance := sb.NewDatabase(db, DbDriver)

	if dbInstance == nil {
		return errors.New("dbInstance is nil")
	}

	Database = dbInstance

	for _, init := range databaseInits {
		err = init(db)

		if err != nil {
			return err
		}
	}

	return nil
}

// migrateDatabase migrates the database
//
// Business logic:
//   - migrates the database for each store
//   - a store is only assigned if it is not nil
//
// Parameters:
// - none
//
// Returns:
// - error: the error if any
func migrateDatabase() (err error) {
	for _, migrate := range databaseMigrations {
		err = migrate(context.Background())

		if err != nil {
			return err
		}
	}

	return nil
}
