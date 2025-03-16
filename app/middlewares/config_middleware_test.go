package middlewares

import (
	"net/http"
	"os"
	"project/config_v2"
	"project/internal/testutils"
	"strings"
	"testing"

	"github.com/dracory/base/test"
)

func TestConfigMiddleware_AttachesConfigToContext(t *testing.T) {
	// Arrange
	testutils.SetupV2SetEnvironmentVariablesOnly()

	// Act
	body, response, err := test.CallMiddleware("GET", NewConfigMiddleware().Handler, func(w http.ResponseWriter, r *http.Request) {
		// Extract the config from the context
		cfg := config_v2.FromContext(r.Context())
		if cfg == nil {
			t.Fatal("Config should not be nil in the context")
		}
		if !cfg.IsEnvTesting() {
			t.Fatal("Environment should be testing")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(cfg.AppEnvironment))
	}, test.NewRequestOptions{})

	// Assert
	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal("Response should not be nil")
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	if body != config_v2.APP_ENVIRONMENT_TESTING {
		t.Fatalf("Expected body %s, got %s", config_v2.APP_ENVIRONMENT_TESTING, body)
	}
}

func TestConfigMiddleware_HandlesErrorGracefully(t *testing.T) {
	// Arrange
	testutils.SetupV2SetEnvironmentVariablesOnly()
	os.Setenv("APP_ENV", "")

	// Act
	body, response, err := test.CallMiddleware("GET", NewConfigMiddleware().Handler, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("Should not be called")
	}, test.NewRequestOptions{})

	// Assert
	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal("Response should not be nil")
	}

	if response.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected status code %d, got %d", http.StatusInternalServerError, response.StatusCode)
	}

	if strings.TrimSpace(body) != ERROR_CONFIG_NOT_FOUND {
		t.Fatalf("Expected body %s, got %s", ERROR_CONFIG_NOT_FOUND, body)
	}
}

func TestConfigMiddleware_PreservesExistingContextValues(t *testing.T) {
	// Set required environment variables for testing
	os.Setenv("APP_ENV", "testing")

	// Arrange
	testutils.Setup()

	// Create a context with a test value
	testKey := struct{}{}
	testValue := "test-value"

	// Act
	var retrievedValue string
	var configFromContext interface{}

	_, response, err := test.CallMiddleware("GET", NewConfigMiddleware().Handler, func(w http.ResponseWriter, r *http.Request) {
		// Extract the test value from the context
		if val, ok := r.Context().Value(testKey).(string); ok {
			retrievedValue = val
		}

		// Extract the config from the context
		configFromContext = config_v2.FromContext(r.Context())

		w.WriteHeader(http.StatusOK)
	}, test.NewRequestOptions{
		Context: map[any]any{
			testKey: testValue,
		},
	})

	// Assert
	if err != nil {
		t.Fatal(err)
	}

	if response == nil {
		t.Fatal("Response should not be nil")
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
	}

	if retrievedValue != testValue {
		t.Fatalf("Expected context value %s, got %s", testValue, retrievedValue)
	}

	if configFromContext == nil {
		t.Fatal("Config should not be nil in the context")
	}
}
