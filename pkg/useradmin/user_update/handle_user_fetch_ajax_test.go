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
	"github.com/stretchr/testify/assert"
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

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var apiResponse map[string]any
	assert.NoError(t, json.Unmarshal([]byte(body), &apiResponse))
	assert.Equal(t, "error", apiResponse["status"])
	assert.Equal(t, "Method not allowed", apiResponse["message"])
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

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var apiResponse map[string]any
	assert.NoError(t, json.Unmarshal([]byte(body), &apiResponse))
	assert.Equal(t, "error", apiResponse["status"])
	assert.Equal(t, "User ID is required", apiResponse["message"])
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

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var apiResponse map[string]any
	assert.NoError(t, json.Unmarshal([]byte(body), &apiResponse))
	assert.Equal(t, "error", apiResponse["status"])
	assert.Equal(t, "User not found", apiResponse["message"])
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

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var apiResponse map[string]any
	assert.NoError(t, json.Unmarshal([]byte(body), &apiResponse))
	assert.Equal(t, "error", apiResponse["status"])
	assert.Equal(t, "GeoStore is not configured", apiResponse["message"])
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

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var apiResponse map[string]any
	assert.NoError(t, json.Unmarshal([]byte(body), &apiResponse))
	assert.Equal(t, "success", apiResponse["status"])

	data, ok := apiResponse["data"].(map[string]any)
	assert.True(t, ok, "response data should be a map")
	countries, ok := data[FieldCountries].([]any)
	assert.True(t, ok, "response should contain countries array")
	assert.NotEmpty(t, countries, "countries should be loaded from geo store")
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

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var apiResponse map[string]any
	assert.NoError(t, json.Unmarshal([]byte(body), &apiResponse))
	assert.Equal(t, "success", apiResponse["status"])

	data, ok := apiResponse["data"].(map[string]any)
	assert.True(t, ok, "response data should be a map")
	fieldStatus, ok := data[FieldStatusField].(map[string]any)
	assert.True(t, ok, "response should contain field_status map")
	assert.Contains(t, fieldStatus, "role")
	assert.Equal(t, true, fieldStatus["role"])
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
