package detailscomponent

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"
)

// TestHandleAjaxSaveDetailsNotNil verifies HandleAjaxSaveDetails returns a string
func TestHandleAjaxSaveDetailsNotNil(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	body := bytes.NewBufferString(`{"title":"Test","description":"Test desc","status":"active","parent_id":""}`)
	req := httptest.NewRequest("POST", "/", body)

	result := HandleAjaxSaveDetails(app, req, "123")

	if result == "" {
		t.Error("HandleAjaxSaveDetails() should return a non-empty string")
	}
}

// TestHandleAjaxSaveDetailsNilStore verifies HandleAjaxSaveDetails handles nil shop store
func TestHandleAjaxSaveDetailsNilStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	body := bytes.NewBufferString(`{"title":"Test","description":"Test desc","status":"active","parent_id":""}`)
	req := httptest.NewRequest("POST", "/", body)

	result := HandleAjaxSaveDetails(app, req, "123")

	if result == "" {
		t.Error("HandleAjaxSaveDetails() should return a non-empty string even with nil shop store")
	}
}

// TestHandleAjaxSaveDetailsInvalidJSON verifies error handling for invalid JSON body
func TestHandleAjaxSaveDetailsInvalidJSON(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	body := bytes.NewBufferString(`invalid json`)
	req := httptest.NewRequest("POST", "/", body)

	result := HandleAjaxSaveDetails(app, req, "123")

	if result == "" {
		t.Error("HandleAjaxSaveDetails() should return a non-empty string for invalid JSON")
	}
	if !strings.Contains(result, "Invalid request body") && !strings.Contains(result, "error") {
		t.Error("HandleAjaxSaveDetails() should return an error response for invalid JSON")
	}
}

// TestHandleAjaxSaveDetailsCategoryNotFound verifies error handling for non-existent category
func TestHandleAjaxSaveDetailsCategoryNotFound(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	body := bytes.NewBufferString(`{"title":"Test","description":"Test desc","status":"active","parent_id":""}`)
	req := httptest.NewRequest("POST", "/", body)

	result := HandleAjaxSaveDetails(app, req, "nonexistent-category-id")

	if result == "" {
		t.Error("HandleAjaxSaveDetails() should return a non-empty string for non-existent category")
	}
	if !strings.Contains(result, "Category not found") && !strings.Contains(result, "not found") {
		t.Error("HandleAjaxSaveDetails() should return 'category not found' error")
	}
}
