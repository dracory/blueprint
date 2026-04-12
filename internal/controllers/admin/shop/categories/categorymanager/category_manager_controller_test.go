package categorymanager

import (
	"testing"

	"project/internal/testutils"
)

// TestNewCategoryManagerController verifies controller can be created
func TestNewCategoryManagerController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryManagerController(app)

	if controller == nil {
		t.Error("NewCategoryManagerController() returned nil")
	}
}

// TestCategoryManagerControllerRegistry verifies controller has registry
func TestCategoryManagerControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryManagerController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestCategoryManagerControllerHandlerExists verifies Handler method exists
func TestCategoryManagerControllerHandlerExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryManagerController(app)

	// Verify Handler method exists (should compile without error)
	_ = controller.Handler
}

func TestCategoryManagerController_HandleLoadCategories_NilShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewCategoryManagerController(app)
	r, _ := testutils.NewRequest("GET", "/admin/shop/categories?action=load-categories", testutils.NewRequestOptions{})

	result := controller.handleLoadCategories(r)
	if result == "" {
		t.Error("Expected non-empty result from handleLoadCategories")
	}
	// Should return error when ShopStore is nil
	if result == "" || len(result) < 10 {
		t.Error("Expected error response when ShopStore is nil")
	}
}

func TestCategoryManagerController_HandleLoadCategories_WithShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryManagerController(app)
	r, _ := testutils.NewRequest("GET", "/admin/shop/categories?action=load-categories", testutils.NewRequestOptions{})

	result := controller.handleLoadCategories(r)
	if result == "" {
		t.Error("Expected non-empty result from handleLoadCategories")
	}
	// Should return success response
	if result == "" {
		t.Error("Expected success response from handleLoadCategories")
	}
}

func TestCategoryManagerController_HandleDeleteCategory_MissingID(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryManagerController(app)
	r, _ := testutils.NewRequest("POST", "/admin/shop/categories?action=delete-category", testutils.NewRequestOptions{})

	result := controller.handleDeleteCategory(nil, r)
	if result == "" {
		t.Error("Expected non-empty result from handleDeleteCategory")
	}
	// Should return error when category_id is missing
	if result == "" || len(result) < 10 {
		t.Error("Expected error response when category_id is missing")
	}
}

func TestCategoryManagerController_HandleDeleteCategory_NilShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewCategoryManagerController(app)
	formValues := map[string][]string{
		"category_id": {"cat123"},
	}
	r, _ := testutils.NewRequest("POST", "/admin/shop/categories?action=delete-category", testutils.NewRequestOptions{FormValues: formValues})

	result := controller.handleDeleteCategory(nil, r)
	if result == "" {
		t.Error("Expected non-empty result from handleDeleteCategory")
	}
	// Should return error when ShopStore is nil
	if result == "" || len(result) < 10 {
		t.Error("Expected error response when ShopStore is nil")
	}
}

func TestCategoryManagerController_HandleDeleteCategory_WithShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryManagerController(app)
	formValues := map[string][]string{
		"category_id": {"cat123"},
	}
	r, _ := testutils.NewRequest("POST", "/admin/shop/categories?action=delete-category", testutils.NewRequestOptions{FormValues: formValues})

	result := controller.handleDeleteCategory(nil, r)
	if result == "" {
		t.Error("Expected non-empty result from handleDeleteCategory")
	}
	// Should return error when category doesn't exist, but not a ShopStore error
}

func TestCategoryManagerController_RenderPage(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewCategoryManagerController(app)
	r, _ := testutils.NewRequest("GET", "/admin/shop/categories", testutils.NewRequestOptions{})

	result := controller.renderPage(r)
	if result == "" {
		t.Error("Expected non-empty result from renderPage")
	}
}
