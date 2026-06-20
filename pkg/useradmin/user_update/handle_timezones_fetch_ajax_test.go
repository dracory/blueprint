package user_update

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
)

// TestHandleGetTimezonesAjaxMethodCheck verifies that handleTimezonesFetchAjax rejects non-POST requests
func TestHandleGetTimezonesAjaxMethodCheck(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
	)

	controller := NewUserUpdateController(app)
	body, response, err := test.CallStringEndpoint(http.MethodGet, controller.Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionGetTimezones},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var apiResponse map[string]any
	assert.NoError(t, json.Unmarshal([]byte(body), &apiResponse))
	assert.Equal(t, "error", apiResponse["status"])
	assert.Equal(t, "Method not allowed", apiResponse["message"])
}

// TestHandleGetTimezonesAjaxCountryCodeValidation verifies that handleTimezonesFetchAjax requires country_code
func TestHandleGetTimezonesAjaxCountryCodeValidation(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
	)

	controller := NewUserUpdateController(app)
	body, response, err := test.CallStringEndpoint(http.MethodPost, controller.Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {actionGetTimezones},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var apiResponse map[string]any
	assert.NoError(t, json.Unmarshal([]byte(body), &apiResponse))
	assert.Equal(t, "error", apiResponse["status"])
	assert.Equal(t, "Country code is required", apiResponse["message"])
}

// TestHandleGetTimezonesAjaxGeoStoreNilCheck verifies that handleTimezonesFetchAjax requires GeoStore
func TestHandleGetTimezonesAjaxGeoStoreNilCheck(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	controller := NewUserUpdateController(app)
	body, response, err := test.CallStringEndpoint(http.MethodPost, controller.Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action":       {actionGetTimezones},
			"country_code": {"US"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var apiResponse map[string]any
	assert.NoError(t, json.Unmarshal([]byte(body), &apiResponse))
	assert.Equal(t, "error", apiResponse["status"])
	assert.Equal(t, "GeoStore is not configured", apiResponse["message"])
}

// TestHandleGetTimezonesAjaxTimezoneList verifies that handleTimezonesFetchAjax returns timezones
func TestHandleGetTimezonesAjaxTimezoneList(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithGeoStore(true),
	)

	controller := NewUserUpdateController(app)
	body, response, err := test.CallStringEndpoint(http.MethodPost, controller.Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action":       {actionGetTimezones},
			"country_code": {"US"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	var apiResponse map[string]any
	assert.NoError(t, json.Unmarshal([]byte(body), &apiResponse))
	assert.Equal(t, "success", apiResponse["status"])

	data, ok := apiResponse["data"].(map[string]any)
	assert.True(t, ok, "response data should be a map")
	timezones, ok := data[FieldTimezones].([]any)
	assert.True(t, ok, "response should contain timezones array")
	assert.NotEmpty(t, timezones, "timezones should be loaded from geo store")
}
