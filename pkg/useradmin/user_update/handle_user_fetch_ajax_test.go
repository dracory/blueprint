package user_update

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"project/internal/app"
	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/dracory/userstore"
)

// TestHandleUserFetchAjaxMethodCheck verifies that handleUserFetchAjax rejects non-POST requests
func TestHandleUserFetchAjaxMethodCheck(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithGeoStore(true),
	)

	controller := NewUserUpdateController(app)
	body, response, err := test.CallStringEndpoint(http.MethodGet, controller.Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionUserFetch},
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var apiResponse map[string]any
	if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if apiResponse["status"] != "error" {
		t.Errorf("expected status error, got %v", apiResponse["status"])
	}
	if apiResponse["message"] != "Method not allowed" {
		t.Errorf("expected message 'Method not allowed', got %v", apiResponse["message"])
	}
}

// TestHandleUserFetchAjaxUserIDValidation verifies that handleUserFetchAjax requires user_id
func TestHandleUserFetchAjaxUserIDValidation(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithGeoStore(true),
	)

	controller := NewUserUpdateController(app)
	body, response, err := test.CallStringEndpoint(http.MethodPost, controller.Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionUserFetch},
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var apiResponse map[string]any
	if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if apiResponse["status"] != "error" {
		t.Errorf("expected status error, got %v", apiResponse["status"])
	}
	if apiResponse["message"] != "User ID is required" {
		t.Errorf("expected message 'User ID is required', got %v", apiResponse["message"])
	}
}

// TestHandleUserFetchAjaxUserLookup verifies that handleUserFetchAjax reports unknown users
func TestHandleUserFetchAjaxUserLookup(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithGeoStore(true),
	)

	controller := NewUserUpdateController(app)
	body, response, err := test.CallStringEndpoint(http.MethodPost, controller.Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action":  {actionUserFetch},
			"user_id": {"non-existent-id"},
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var apiResponse map[string]any
	if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if apiResponse["status"] != "error" {
		t.Errorf("expected status error, got %v", apiResponse["status"])
	}
	if apiResponse["message"] != "User not found" {
		t.Errorf("expected message 'User not found', got %v", apiResponse["message"])
	}
}

// TestHandleUserFetchAjaxGeoStoreNilCheck verifies that handleUserFetchAjax requires GeoStore
func TestHandleUserFetchAjaxGeoStoreNilCheck(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user := seedFetchTestUser(t, app)

	controller := NewUserUpdateController(app)
	body, response, err := test.CallStringEndpoint(http.MethodPost, controller.Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action":  {actionUserFetch},
			"user_id": {user.GetID()},
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var apiResponse map[string]any
	if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if apiResponse["status"] != "error" {
		t.Errorf("expected status error, got %v", apiResponse["status"])
	}
	if apiResponse["message"] != "GeoStore is not configured" {
		t.Errorf("expected message 'GeoStore is not configured', got %v", apiResponse["message"])
	}
}

// TestHandleUserFetchAjaxGeoStoreErrorHandling verifies that handleUserFetchAjax returns user data
// and country/timezone lists when GeoStore is configured
func TestHandleUserFetchAjaxGeoStoreErrorHandling(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithGeoStore(true),
	)

	user := seedFetchTestUser(t, app)

	controller := NewUserUpdateController(app)
	body, response, err := test.CallStringEndpoint(http.MethodPost, controller.Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action":  {actionUserFetch},
			"user_id": {user.GetID()},
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var apiResponse map[string]any
	if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if apiResponse["status"] != "success" {
		t.Errorf("expected status success, got %v", apiResponse["status"])
	}

	data, ok := apiResponse["data"].(map[string]any)
	if !ok {
		t.Fatal("response data should be a map")
	}
	countries, ok := data[FieldCountries].([]any)
	if !ok {
		t.Fatal("response should contain countries array")
	}
	if len(countries) == 0 {
		t.Error("countries should be loaded from geo store")
	}
}

// TestHandleUserFetchAjaxFieldStatusIncludesRole verifies that field_status includes the role field
func TestHandleUserFetchAjaxFieldStatusIncludesRole(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithGeoStore(true),
	)

	user := seedFetchTestUser(t, app)

	controller := NewUserUpdateController(app)
	body, response, err := test.CallStringEndpoint(http.MethodPost, controller.Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action":  {actionUserFetch},
			"user_id": {user.GetID()},
		},
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	var apiResponse map[string]any
	if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if apiResponse["status"] != "success" {
		t.Errorf("expected status success, got %v", apiResponse["status"])
	}

	data, ok := apiResponse["data"].(map[string]any)
	if !ok {
		t.Fatal("response data should be a map")
	}
	fieldStatus, ok := data[FieldStatusField].(map[string]any)
	if !ok {
		t.Fatal("response should contain field_status map")
	}
	if _, ok := fieldStatus["role"]; !ok {
		t.Error("field_status should contain role")
	}
	if fieldStatus["role"] != true {
		t.Errorf("expected role to be true, got %v", fieldStatus["role"])
	}
}

// seedFetchTestUser creates a test user in the in-memory database and returns it
func seedFetchTestUser(t *testing.T, app app.AppInterface) userstore.UserInterface {
	t.Helper()
	user := userstore.NewUser()
	user.SetFirstName("Test")
	user.SetLastName("User")
	user.SetEmail("fetch-test@example.com")
	user.SetStatus(userstore.USER_STATUS_ACTIVE)
	user.SetRole(userstore.USER_ROLE_USER)

	if err := app.GetUserStore().UserCreate(context.Background(), user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}
