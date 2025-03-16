package testutils

import (
	"project/config"
	"testing"

	baseTestUtils "github.com/dracory/base/testutils"
)

func TestTestKeyIntegration(t *testing.T) {
	// Test that the TestKey function from the blueprint project
	// produces the same result as the TestKey function from the base project
	blueprintKey := TestKey()
	baseKey := baseTestUtils.TestKey(config.DbDriver, config.DbHost, config.DbPort, config.DbName, config.DbUser, config.DbPass)
	
	if blueprintKey != baseKey {
		t.Errorf("Blueprint TestKey and base TestKey should produce the same result")
	}
}

func TestTestConfigIntegration(t *testing.T) {
	// Create a test configuration
	testConfig := baseTestUtils.DefaultTestConfig()
	
	// Set up the test environment
	baseTestUtils.SetupTestEnvironment(testConfig)
	
	// Clean up after ourselves
	defer baseTestUtils.CleanupTestEnvironment(testConfig)
	
	// Initialize the application with test configuration
	config.Initialize()
	
	// Verify that the configuration was applied
	if config.AppName != testConfig.AppName {
		t.Errorf("Expected AppName to be %s, got %s", testConfig.AppName, config.AppName)
	}
}

func TestTestDBIntegration(t *testing.T) {
	// Create a test database
	db, err := baseTestUtils.NewTestDB(nil)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer baseTestUtils.CloseTestDB(db)
	
	// Create a test table
	err = baseTestUtils.CreateTestTable(db, "test_users", "id INTEGER PRIMARY KEY, name TEXT")
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}
	
	// Insert test data
	err = baseTestUtils.ExecuteSQLWithArgs(db, "INSERT INTO test_users (name) VALUES (?)", "Test User")
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
