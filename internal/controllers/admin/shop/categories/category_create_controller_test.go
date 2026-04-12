package categories

import (
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

// TestNewCategoryCreateController verifies controller can be created
func TestNewCategoryCreateController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryCreateController(app)

	if controller == nil {
		t.Error("NewCategoryCreateController() returned nil")
	}
}

// TestCategoryCreateControllerRegistry verifies controller has registry
func TestCategoryCreateControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryCreateController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestCategoryCreateControllerHandlerExists verifies Handler method exists
func TestCategoryCreateControllerHandlerExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryCreateController(app)

	// Verify Handler method exists (should compile without error)
	_ = controller.Handler
}

func TestCategoryCreateController_RenderPage(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryCreateController(app)
	r, _ := testutils.NewRequest("GET", "/admin/shop/categories/create", testutils.NewRequestOptions{})

	result := controller.renderPage(r)
	if result == "" {
		t.Error("Expected non-empty result from renderPage")
	}
}

func TestCategoryCreateController_RenderPage_NilShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewCategoryCreateController(app)
	r, _ := testutils.NewRequest("GET", "/admin/shop/categories/create", testutils.NewRequestOptions{})

	result := controller.renderPage(r)
	if result == "" {
		t.Error("Expected non-empty result from renderPage even with nil ShopStore")
	}
}

func TestCategoryCreateController_HandlePost_MissingTitle(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true), testutils.WithCacheStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryCreateController(app)
	formValues := map[string][]string{
		"description": {"Test description"},
		"status":      {"active"},
	}
	r, _ := testutils.NewRequest("POST", "/admin/shop/categories/create", testutils.NewRequestOptions{FormValues: formValues})
	w := httptest.NewRecorder()

	result := controller.handlePost(w, r)
	if result == "" {
		t.Error("Expected non-empty result from handlePost")
	}
	// Should return error when title is missing
}

func TestCategoryCreateController_HandlePost_WithShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true), testutils.WithCacheStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryCreateController(app)
	formValues := map[string][]string{
		"title":       {"Test Category"},
		"description": {"Test description"},
		"status":      {"active"},
	}
	r, _ := testutils.NewRequest("POST", "/admin/shop/categories/create", testutils.NewRequestOptions{FormValues: formValues})
	w := httptest.NewRecorder()

	result := controller.handlePost(w, r)
	if result == "" {
		t.Error("Expected non-empty result from handlePost")
	}
	// Should return success or error depending on whether category creation succeeds
}
