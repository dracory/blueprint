package config

import (
	"os"
	"testing"
)

func TestLoad_Success(t *testing.T) {
	// Setup minimal required env vars
	mustSetenv(t, KEY_APP_HOST, "localhost")
	mustSetenv(t, KEY_APP_PORT, "8080")
	mustSetenv(t, KEY_APP_ENVIRONMENT, "testing")
	mustSetenv(t, KEY_DB_DRIVER, "sqlite")
	mustSetenv(t, KEY_DB_DATABASE, ":memory:")
	defer cleanupEnv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() returned nil config")
	}

	if cfg.GetAppHost() != "localhost" {
		t.Errorf("expected host=localhost, got %s", cfg.GetAppHost())
	}

	if cfg.GetAppPort() != "8080" {
		t.Errorf("expected port=8080, got %s", cfg.GetAppPort())
	}

	if cfg.GetDatabaseDriver() != "sqlite" {
		t.Errorf("expected driver=sqlite, got %s", cfg.GetDatabaseDriver())
	}
}

func TestLoad_MissingRequiredFields(t *testing.T) {
	cleanupEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("Load() should fail with missing required fields")
	}

	verr, ok := err.(validationError)
	if !ok {
		t.Fatalf("expected validationError, got %T", err)
	}

	if len(verr.errs) == 0 {
		t.Error("expected validation errors, got none")
	}
}

func TestLoad_DatabasePostgresRequirements(t *testing.T) {
	mustSetenv(t, KEY_APP_HOST, "localhost")
	mustSetenv(t, KEY_APP_PORT, "8080")
	mustSetenv(t, KEY_APP_ENVIRONMENT, "testing")
	mustSetenv(t, KEY_DB_DRIVER, "postgres")
	mustSetenv(t, KEY_DB_DATABASE, "testdb")
	// Missing host, port, username, password
	defer cleanupEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("Load() should fail when postgres driver missing connection details")
	}

	verr, ok := err.(validationError)
	if !ok {
		t.Fatalf("expected validationError, got %T", err)
	}

	// Should have errors for host, port, username, password
	if len(verr.errs) < 4 {
		t.Errorf("expected at least 4 validation errors for postgres, got %d", len(verr.errs))
	}
}

func TestLoad_EnvEncryptionKeyOptional(t *testing.T) {
	mustSetenv(t, KEY_APP_HOST, "localhost")
	mustSetenv(t, KEY_APP_PORT, "8080")
	mustSetenv(t, KEY_APP_ENVIRONMENT, "testing")
	mustSetenv(t, KEY_DB_DRIVER, "sqlite")
	mustSetenv(t, KEY_DB_DATABASE, ":memory:")
	// No ENVENC_KEY_PRIVATE set
	defer cleanupEnv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() should succeed without env encryption key: %v", err)
	}

	if cfg.GetEnvEncryptionKey() != "" {
		t.Error("expected empty encryption key when not provided")
	}
}

func TestLoad_CMSStoreRequiresTemplateID(t *testing.T) {
	mustSetenv(t, KEY_APP_HOST, "localhost")
	mustSetenv(t, KEY_APP_PORT, "8080")
	mustSetenv(t, KEY_APP_ENVIRONMENT, "testing")
	mustSetenv(t, KEY_DB_DRIVER, "sqlite")
	mustSetenv(t, KEY_DB_DATABASE, ":memory:")
	// Missing CMS_STORE_TEMPLATE_ID
	defer cleanupEnv()

	if !cmsStoreUsed {
		return // CMS store not enabled
	}

	_, err := Load()
	if err == nil {
		t.Fatal("Load() should fail when CMS store enabled without template ID")
	}

	verr, ok := err.(validationError)
	if !ok {
		t.Fatalf("expected validationError, got %T", err)
	}

	found := false
	for _, e := range verr.errs {
		if merr, ok := e.(MissingEnvError); ok && merr.Key == KEY_CMS_STORE_TEMPLATE_ID {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected validation error for missing CMS_STORE_TEMPLATE_ID")
	}
}

func TestLoad_LLMProviderRequirements(t *testing.T) {
	mustSetenv(t, KEY_APP_HOST, "localhost")
	mustSetenv(t, KEY_APP_PORT, "8080")
	mustSetenv(t, KEY_APP_ENVIRONMENT, "testing")
	mustSetenv(t, KEY_DB_DRIVER, "sqlite")
	mustSetenv(t, KEY_DB_DATABASE, ":memory:")
	mustSetenv(t, KEY_OPENAI_API_USED, "true")
	// Missing OPENAI_API_KEY and OPENAI_API_DEFAULT_MODEL
	defer cleanupEnv()

	_, err := Load()
	if err == nil {
		t.Fatal("Load() should fail when OpenAI enabled without credentials")
	}

	verr, ok := err.(validationError)
	if !ok {
		t.Fatalf("expected validationError, got %T", err)
	}

	foundKey := false
	foundModel := false
	for _, e := range verr.errs {
		if merr, ok := e.(MissingEnvError); ok {
			if merr.Key == KEY_OPENAI_API_KEY {
				foundKey = true
			}
			if merr.Key == KEY_OPENAI_API_DEFAULT_MODEL {
				foundModel = true
			}
		}
	}

	if !foundKey {
		t.Error("expected validation error for missing OPENAI_API_KEY")
	}
	if !foundModel {
		t.Error("expected validation error for missing OPENAI_API_DEFAULT_MODEL")
	}
}

func TestLoad_StripeConfiguration(t *testing.T) {
	mustSetenv(t, KEY_APP_HOST, "localhost")
	mustSetenv(t, KEY_APP_PORT, "8080")
	mustSetenv(t, KEY_APP_ENVIRONMENT, "testing")
	mustSetenv(t, KEY_DB_DRIVER, "sqlite")
	mustSetenv(t, KEY_DB_DATABASE, ":memory:")
	mustSetenv(t, KEY_STRIPE_KEY_PRIVATE, "sk_test_123")
	mustSetenv(t, KEY_STRIPE_KEY_PUBLIC, "pk_test_123")
	defer cleanupEnv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if !cfg.GetStripeUsed() {
		t.Error("expected Stripe to be marked as used when keys provided")
	}

	if cfg.GetStripeKeyPrivate() != "sk_test_123" {
		t.Errorf("expected private key=sk_test_123, got %s", cfg.GetStripeKeyPrivate())
	}
}

func TestLoad_MailConfiguration(t *testing.T) {
	mustSetenv(t, KEY_APP_HOST, "localhost")
	mustSetenv(t, KEY_APP_PORT, "8080")
	mustSetenv(t, KEY_APP_ENVIRONMENT, "testing")
	mustSetenv(t, KEY_DB_DRIVER, "sqlite")
	mustSetenv(t, KEY_DB_DATABASE, ":memory:")
	mustSetenv(t, KEY_MAIL_DRIVER, "smtp")
	mustSetenv(t, KEY_MAIL_HOST, "smtp.example.com")
	mustSetenv(t, KEY_MAIL_PORT, "587")
	mustSetenv(t, KEY_MAIL_USERNAME, "user@example.com")
	mustSetenv(t, KEY_MAIL_FROM_ADDRESS, "noreply@example.com")
	defer cleanupEnv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.GetMailDriver() != "smtp" {
		t.Errorf("expected mail driver=smtp, got %s", cfg.GetMailDriver())
	}

	if cfg.GetMailPort() != 587 {
		t.Errorf("expected mail port=587, got %d", cfg.GetMailPort())
	}
}

func TestLoad_TranslationDefaults(t *testing.T) {
	mustSetenv(t, KEY_APP_HOST, "localhost")
	mustSetenv(t, KEY_APP_PORT, "8080")
	mustSetenv(t, KEY_APP_ENVIRONMENT, "testing")
	mustSetenv(t, KEY_DB_DRIVER, "sqlite")
	mustSetenv(t, KEY_DB_DATABASE, ":memory:")
	defer cleanupEnv()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.GetTranslationLanguageDefault() == "" {
		t.Error("expected default translation language to be set")
	}

	if len(cfg.GetTranslationLanguageList()) == 0 {
		t.Error("expected translation language list to be populated")
	}
}

func TestLoad_VaultStoreRequirements(t *testing.T) {
	mustSetenv(t, KEY_APP_HOST, "localhost")
	mustSetenv(t, KEY_APP_PORT, "8080")
	mustSetenv(t, KEY_APP_ENVIRONMENT, "testing")
	mustSetenv(t, KEY_DB_DRIVER, "sqlite")
	mustSetenv(t, KEY_DB_DATABASE, ":memory:")
	defer cleanupEnv()

	if !userStoreVaultEnabled {
		return // User vault not enabled
	}

	_, err := Load()
	if err == nil {
		t.Fatal("Load() should fail when user vault enabled but vault store not used")
	}

	verr, ok := err.(validationError)
	if !ok {
		t.Fatalf("expected validationError, got %T", err)
	}

	found := false
	for _, e := range verr.errs {
		if e.Error() != "" {
			found = true
			break
		}
	}

	if !found {
		t.Error("expected validation error for vault store dependency")
	}
}

// cleanupEnv clears all config-related environment variables
func cleanupEnv() {
	os.Clearenv()
}

// mustSetenv is a test helper that sets an environment variable and fails the test on error.
func mustSetenv(t *testing.T, key, value string) {
	t.Helper()
	if err := os.Setenv(key, value); err != nil {
		t.Fatalf("failed to set env %s: %v", key, err)
	}
}
