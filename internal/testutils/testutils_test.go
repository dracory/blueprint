package testutils

import (
	"project/internal/config"
	"testing"

	"github.com/dracory/base/test"
)

func TestTestKeyIntegration(t *testing.T) {
	// Test that the TestKey function from the blueprint project
	// produces the same result as the TestKey function from the base project
	blueprintKey := TestKey()
	baseKey := test.TestKey(config.DbDriver, config.DbHost, config.DbPort, config.DbName, config.DbUser, config.DbPass)

	if blueprintKey != baseKey {
		t.Errorf("Blueprint TestKey and base TestKey should produce the same result")
	}
}

func TestTestConfigIntegration(t *testing.T) {
	// Create a test configuration
	testConfig := test.DefaultTestConfig()

	// Set up the test environment
	test.SetupTestEnvironment(testConfig)

	// Clean up after ourselves
	defer test.CleanupTestEnvironment(testConfig)

	// os.Setenv("CMS_TEMPLATE_ID", "TEST_TEMPLATE_ID")
	// os.Setenv("STRIPE_KEY_PRIVATE", "TEST_STRIPE_KEY_PRIVATE")
	// os.Setenv("STRIPE_KEY_PUBLIC", "TEST_STRIPE_KEY_PUBLIC")
	// os.Setenv("VERTEX_MODEL_ID", "TEST_VERTEX_MODEL_ID")
	// os.Setenv("VERTEX_PROJECT_ID", "TEST_VERTEX_PROJECT_ID")
	// os.Setenv("VERTEX_REGION_ID", "TEST_VERTEX_REGION_ID")

	// Initialize the application with test configuration
	config.Initialize()

	// Verify that the configuration was applied
	if config.AppName != testConfig.AppName {
		t.Errorf("Expected AppName to be %s, got %s", testConfig.AppName, config.AppName)
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
