package user_update

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
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
	if apiResponse["message"] != "Country code is required" {
		t.Errorf("expected message 'Country code is required', got %v", apiResponse["message"])
	}
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
	timezones, ok := data[FieldTimezones].([]any)
	if !ok {
		t.Fatal("response should contain timezones array")
	}
	if len(timezones) == 0 {
		t.Error("timezones should be loaded from geo store")
	}
}
