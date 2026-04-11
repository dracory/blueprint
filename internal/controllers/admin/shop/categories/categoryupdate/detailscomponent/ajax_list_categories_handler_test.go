package detailscomponent

import (
	"testing"

	"project/internal/testutils"
)

// TestHandleAjaxListCategoriesNotNil verifies HandleAjaxListCategories returns a string
func TestHandleAjaxListCategoriesNotNil(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	result := HandleAjaxListCategories(app)

	if result == "" {
		t.Error("HandleAjaxListCategories() should return a non-empty string")
	}
}

// TestHandleAjaxListCategoriesNilStore verifies HandleAjaxListCategories handles nil shop store
func TestHandleAjaxListCategoriesNilStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}

	result := HandleAjaxListCategories(app)

	if result == "" {
		t.Error("HandleAjaxListCategories() should return a non-empty string even with nil shop store")
	}
}
