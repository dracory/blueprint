package categoryupdate

import (
	"context"
	"net/http/httptest"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/shopstore"
)

// TestNewCategoryUpdateController verifies controller can be created
func TestNewCategoryUpdateController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryUpdateController(app)

	if controller == nil {
		t.Error("NewCategoryUpdateController() returned nil")
	}
}

// TestCategoryUpdateControllerRegistry verifies controller has registry
func TestCategoryUpdateControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryUpdateController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestCategoryUpdateControllerHandlerExists verifies Handler method exists
func TestCategoryUpdateControllerHandlerExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryUpdateController(app)

	// Verify Handler method exists (should compile without error)
	_ = controller.Handler
}

// TestCategoryUpdateControllerHandlerMissingCategoryID verifies error when category_id is missing
func TestCategoryUpdateControllerHandlerMissingCategoryID(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryUpdateController(app)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	result := controller.Handler(w, r)

	if result == "" {
		t.Error("Handler() should return error message when category_id is missing")
	}
}

// TestCategoryUpdateControllerHandlerInvalidCategoryID verifies error when category_id is invalid
func TestCategoryUpdateControllerHandlerInvalidCategoryID(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryUpdateController(app)
	w := httptest.NewRecorder()

	values := url.Values{}
	values.Add("category_id", "invalid-id")
	r := httptest.NewRequest("GET", "/?"+values.Encode(), nil)

	result := controller.Handler(w, r)

	if result == "" {
		t.Error("Handler() should return error message when category_id is invalid")
	}
}

// TestCategoryUpdateControllerHandlerDefaultAction verifies default action renders page
func TestCategoryUpdateControllerHandlerDefaultAction(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryUpdateController(app)

	// First create a category
	category := shopstore.NewCategory()
	category.SetTitle("Test Category")
	err := app.GetShopStore().CategoryCreate(context.Background(), category)
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	w := httptest.NewRecorder()
	values := url.Values{}
	values.Add("category_id", category.GetID())
	r := httptest.NewRequest("GET", "/?"+values.Encode(), nil)

	result := controller.Handler(w, r)

	if result == "" {
		t.Error("Handler() should return HTML for default action")
	}
}

// TestCategoryUpdateControllerHandlerLoadDetailsAction verifies load-details action
func TestCategoryUpdateControllerHandlerLoadDetailsAction(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryUpdateController(app)

	// First create a category
	category := shopstore.NewCategory()
	category.SetTitle("Test Category")
	err := app.GetShopStore().CategoryCreate(context.Background(), category)
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	w := httptest.NewRecorder()
	values := url.Values{}
	values.Add("category_id", category.GetID())
	values.Add("action", ACTION_LOAD_DETAILS)
	r := httptest.NewRequest("GET", "/?"+values.Encode(), nil)

	result := controller.Handler(w, r)

	// Should return JSON response from detailscomponent
	if result == "" {
		t.Error("Handler() should return response for load-details action")
	}
}

// TestCategoryUpdateControllerHandlerListCategoriesAction verifies list-categories action
func TestCategoryUpdateControllerHandlerListCategoriesAction(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryUpdateController(app)

	w := httptest.NewRecorder()
	values := url.Values{}
	values.Add("category_id", "test-id")
	values.Add("action", ACTION_LIST_CATEGORIES)
	r := httptest.NewRequest("GET", "/?"+values.Encode(), nil)

	result := controller.Handler(w, r)

	// Should return JSON response from detailscomponent
	if result == "" {
		t.Error("Handler() should return response for list-categories action")
	}
}
