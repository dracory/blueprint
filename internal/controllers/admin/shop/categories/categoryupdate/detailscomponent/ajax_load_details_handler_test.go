package detailscomponent

import (
	"strings"
	"testing"

	"project/internal/testutils"
)

// TestHandleAjaxLoadDetailsNotNil verifies HandleAjaxLoadDetails returns a string
func TestHandleAjaxLoadDetailsNotNil(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	result := HandleAjaxLoadDetails(app, "123")

	if result == "" {
		t.Error("HandleAjaxLoadDetails() should return a non-empty string")
	}
}

// TestHandleAjaxLoadDetailsNilStore verifies HandleAjaxLoadDetails handles nil shop store
func TestHandleAjaxLoadDetailsNilStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	result := HandleAjaxLoadDetails(app, "123")

	if result == "" {
		t.Error("HandleAjaxLoadDetails() should return a non-empty string even with nil shop store")
	}
}

// TestHandleAjaxLoadDetailsCategoryNotFound verifies HandleAjaxLoadDetails handles non-existent category
func TestHandleAjaxLoadDetailsCategoryNotFound(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	result := HandleAjaxLoadDetails(app, "nonexistent-category-id")

	if result == "" {
		t.Error("HandleAjaxLoadDetails() should return a non-empty string for non-existent category")
	}
	if !strings.Contains(result, "Category not found") && !strings.Contains(result, "not found") {
		t.Error("HandleAjaxLoadDetails() should return 'category not found' error")
	}
}
