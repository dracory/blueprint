package testutils

import (
	"project/internal/config"
	"project/internal/types"
	"testing"

	"github.com/dracory/test"
)

func TestTestKeyIntegration(t *testing.T) {
	// Test that the TestKey function from the blueprint project
	// produces the same result as the TestKey function from the base project
	cfg := &types.Config{}
	cfg.SetDatabaseDriver("sqlite")
	cfg.SetDatabaseHost("")
	cfg.SetDatabasePort("")
	cfg.SetDatabaseName("file::memory:?cache=shared")
	cfg.SetDatabaseUsername("")
	cfg.SetDatabasePassword("")
	blueprintKey := TestKey(cfg)
	baseKey := test.TestKey(cfg.GetDatabaseDriver(), cfg.GetDatabaseHost(), cfg.GetDatabasePort(), cfg.GetDatabaseName(), cfg.GetDatabaseUsername(), cfg.GetDatabasePassword())

	if blueprintKey != baseKey {
		t.Errorf("Blueprint TestKey and base TestKey should produce the same result")
	}
}

func TestTestConfigIntegration(t *testing.T) {
	// Create a test configuration
	testConfig := types.Config{}
	testConfig.SetAppName("Test App")
	testConfig.SetAppUrl("http://localhost:8080")
	testConfig.SetAppEnv(config.APP_ENVIRONMENT_TESTING)
	testConfig.SetDatabaseDriver("sqlite")
	testConfig.SetDatabaseHost("")
	testConfig.SetDatabasePort("")
	testConfig.SetDatabaseName("file::memory:?cache=shared")
	testConfig.SetDatabaseUsername("")
	testConfig.SetDatabasePassword("")
	testConfig.SetAppDebug(true)
	testConfig.SetEnvEncryptionKey("123456")
	testConfig.SetAppHost("localhost")
	testConfig.SetAppPort("8080")
	testConfig.SetMailDriver("smtp")
	testConfig.SetMailHost("127.0.0.1")
	testConfig.SetMailPort(32435)
	testConfig.SetMailUsername("")
	testConfig.SetMailPassword("")
	testConfig.SetMailFromEmail("admintest@test.com")
	testConfig.SetMailFromName("Admin Test")
	testConfig.SetCMSTemplateID("default")
	testConfig.SetVaultKey("abcdefghijklmnopqrstuvwxyz1234567890")
	testConfig.SetOpenAIKey("openai_api_key")
	testConfig.SetStripeKeyPrivate("sk_test_yoursecretkey")
	testConfig.SetStripeKeyPublic("pk_test_yourpublickey")
	testConfig.SetVertexProjectID("vertex_project_id")
	testConfig.SetVertexRegionID("vertex_region_id")
	testConfig.SetVertexModelID("vertex_model_id")

	// os.Setenv("VERTEX_REGION_ID", "TEST_VERTEX_REGION_ID")

	// Verify that the configuration was applied
	if "Test App" != testConfig.GetAppName() {
		t.Errorf("Expected AppName to be %s, got %s", "Test App", testConfig.GetAppName())
	}
}

func TestTestDBIntegration(t *testing.T) {
	// Create a test database
	db, err := test.NewTestDB(nil)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer test.CloseTestDB(db)

	// Create a test table
	err = test.CreateTestTable(db, "test_users", "id INTEGER PRIMARY KEY, name TEXT")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	// Insert test data
	err = test.ExecuteSQLWithArgs(db, "INSERT INTO test_users (name) VALUES (?)", "Test User")
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Query the test data
	var name string
	err = db.QueryRow("SELECT name FROM test_users WHERE name = ?", "Test User").Scan(&name)
	if err != nil {
		t.Fatalf("Failed to query test data: %v", err)
	}

	// Verify the result
	if name != "Test User" {
		t.Errorf("Expected name to be 'Test User', got '%s'", name)
	}
}
