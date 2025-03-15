package config

import (
	"os"

	_ "modernc.org/sqlite"
)

// TestConfig creates and initializes a configuration for testing
func TestConfig() (*Config, error) {
	// Set environment variables for testing
	setTestEnvironmentVariables()

	// Create new config
	config, err := New()
	if err != nil {
		return nil, err
	}

	// Initialize the config
	err = config.Initialize()
	if err != nil {
		return nil, err
	}

	return config, nil
}

// setTestEnvironmentVariables sets environment variables for testing
func setTestEnvironmentVariables() {
	os.Setenv("APP_NAME", "TEST APP NAME")
	os.Setenv("APP_URL", "http://localhost:8080")
	os.Setenv("APP_ENV", "testing")

	os.Setenv("DB_DRIVER", "sqlite")
	os.Setenv("DB_HOST", "")
	os.Setenv("DB_PORT", "")
	os.Setenv("DB_DATABASE", "file::memory:?cache=shared")
	os.Setenv("DB_USERNAME", "")
	os.Setenv("DB_PASSWORD", "")

	// os.Setenv("DEBUG", "yes")

	os.Setenv("ENV_ENCRYPTION_KEY", "123456")

	os.Setenv("SERVER_HOST", "localhost")
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("MAIL_DRIVER", "smtp")
	os.Setenv("MAIL_HOST", "127.0.0.1")
	os.Setenv("MAIL_PORT", "32435")
	os.Setenv("MAIL_USERNAME", "")
	os.Setenv("MAIL_PASSWORD", "")

	os.Setenv("EMAIL_FROM_ADDRESS", "admintest@test.com")
	os.Setenv("EMAIL_FROM_NAME", "Admin Test")

	os.Setenv("CMS_TEMPLATE_ID", "default")

	os.Setenv("VAULT_KEY", "abcdefghijklmnopqrstuvwxyz1234567890")

	os.Setenv("OPENAI_API_KEY", "openai_api_key")

	os.Setenv("STRIPE_KEY_PRIVATE", "sk_test_yoursecretkey")
	os.Setenv("STRIPE_KEY_PUBLIC", "pk_test_yourpublickey")

	os.Setenv("VERTEX_PROJECT_ID", "vertex_project_id")
	os.Setenv("VERTEX_REGION_ID", "vertex_region_id")
	os.Setenv("VERTEX_MODEL_ID", "vertex_model_id")
}
