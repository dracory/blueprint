package admin

import (
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/shopstore"
)

// TestNewProductManagerController verifies controller can be created
func TestNewProductManagerController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductManagerController(app)

	if controller == nil {
		t.Error("NewProductManagerController() returned nil")
	}
}

// TestProductManagerControllerRegistry verifies controller has registry
func TestProductManagerControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductManagerController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestProductManagerControllerHandlerExists verifies Handler method exists
func TestProductManagerControllerHandlerExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductManagerController(app)

	// Verify Handler method exists (should compile without error)
	_ = controller.Handler
}

func TestProductManagerController_NilRegistry(t *testing.T) {
	t.Parallel()
	controller := NewProductManagerController(nil)
	if controller == nil {
		t.Error("NewProductManagerController(nil) should not return nil")
	}
	if controller.registry != nil {
		t.Error("Controller registry should be nil when passed nil")
	}
}

func TestProductManagerController_MultipleInstances(t *testing.T) {
	t.Parallel()
	app1 := testutils.Setup(testutils.WithShopStore(true))
	app2 := testutils.Setup(testutils.WithShopStore(true))
	if app1 == nil || app2 == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app1.GetDatabase().Close() })
	t.Cleanup(func() { _ = app2.GetDatabase().Close() })

	controller1 := NewProductManagerController(app1)
	controller2 := NewProductManagerController(app2)

	if controller1 == nil || controller2 == nil {
		t.Fatal("All controllers should be non-nil")
	}

	if controller1 == controller2 {
		t.Error("Controllers should be separate instances")
	}

	if controller1.registry != app1 {
		t.Error("Controller1 should have app1")
	}

	if controller2.registry != app2 {
		t.Error("Controller2 should have app2")
	}
}

func TestProductManagerController_Handler_Actions(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductManagerController(app)

	// Verify handler exists - methods cannot be nil in Go
	// This test ensures the Handler method is callable
	_ = controller.Handler
}

func TestProductManagerController_PrepareData_Basic(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductManagerController(app)
	r, _ := testutils.NewRequest("GET", "/admin/shop/products", testutils.NewRequestOptions{})

	data, err := controller.prepareData(r)
	if err != "" {
		t.Errorf("Expected no error, got: %s", err)
	}
	if data.request == nil {
		t.Error("Expected request to be set in data")
	}
	if data.pageInt != 0 {
		t.Errorf("Expected pageInt to be 0, got %d", data.pageInt)
	}
	if data.perPage != 10 {
		t.Errorf("Expected perPage to be 10, got %d", data.perPage)
	}
}

func TestProductManagerController_PrepareData_WithQueryParams(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductManagerController(app)
	queryParams := url.Values{}
	queryParams.Set("page", "2")
	queryParams.Set("per_page", "25")
	queryParams.Set("sort", "asc")
	queryParams.Set("by", "created_at")
	queryParams.Set("filter_status", "active")
	queryParams.Set("filter_title", "test product")
	queryParams.Set("filter_product_id", "prod123")
	r, _ := testutils.NewRequest("GET", "/admin/shop/products", testutils.NewRequestOptions{QueryParams: queryParams})

	data, err := controller.prepareData(r)
	if err != "" {
		t.Errorf("Expected no error, got: %s", err)
	}
	if data.pageInt != 2 {
		t.Errorf("Expected pageInt to be 2, got %d", data.pageInt)
	}
	if data.perPage != 25 {
		t.Errorf("Expected perPage to be 25, got %d", data.perPage)
	}
	if data.sortOrder != "asc" {
		t.Errorf("Expected sortOrder to be 'asc', got '%s'", data.sortOrder)
	}
	if data.formStatus != "active" {
		t.Errorf("Expected formStatus to be 'active', got '%s'", data.formStatus)
	}
	if data.formTitle != "test product" {
		t.Errorf("Expected formTitle to be 'test product', got '%s'", data.formTitle)
	}
	if data.formProductID != "prod123" {
		t.Errorf("Expected formProductID to be 'prod123', got '%s'", data.formProductID)
	}
}

func TestProductManagerController_PrepareData_NilShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewProductManagerController(app)
	r, _ := testutils.NewRequest("GET", "/admin/shop/products", testutils.NewRequestOptions{})

	data, err := controller.prepareData(r)
	if err == "" {
		t.Error("Expected error when ShopStore is nil")
	}
	if data.request == nil {
		t.Error("Expected request to be set in data")
	}
}

func TestProductManagerController_OnModalProductFilterShow(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductManagerController(app)

	data := productManagerControllerData{
		formStatus:      "active",
		formTitle:       "Test Product",
		formCreatedFrom: "2024-01-01",
		formCreatedTo:   "2024-12-31",
		formProductID:   "prod123",
	}

	tag := controller.onModalProductFilterShow(data)
	if tag == nil {
		t.Error("Expected non-nil tag from onModalProductFilterShow")
	}
	html := tag.ToHTML()
	if html == "" {
		t.Error("Expected non-empty HTML from onModalProductFilterShow")
	}
}

func TestProductManagerController_OnModalProductFilterShow_EmptyData(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductManagerController(app)

	data := productManagerControllerData{}

	tag := controller.onModalProductFilterShow(data)
	if tag == nil {
		t.Error("Expected non-nil tag from onModalProductFilterShow with empty data")
	}
	html := tag.ToHTML()
	if html == "" {
		t.Error("Expected non-empty HTML from onModalProductFilterShow with empty data")
	}
}

func TestProductManagerController_Page(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductManagerController(app)
	r, _ := testutils.NewRequest("GET", "/admin/shop/products", testutils.NewRequestOptions{})

	data := productManagerControllerData{
		request:      r,
		productList:  []shopstore.ProductInterface{},
		productCount: 0,
	}

	tag := controller.page(data)
	if tag == nil {
		t.Error("Expected non-nil tag from page")
	}
	html := tag.ToHTML()
	if html == "" {
		t.Error("Expected non-empty HTML from page")
	}
}

func TestProductManagerController_FetchProductList_NilShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewProductManagerController(app)
	data := productManagerControllerData{}

	productList, productCount, err := controller.fetchProductList(data)
	if err == nil {
		t.Error("Expected error when ShopStore is nil")
	}
	// fetchProductList returns empty list, not nil
	if len(productList) != 0 {
		t.Error("Expected empty list when ShopStore is nil")
	}
	if productCount != 0 {
		t.Error("Expected productCount to be 0 when ShopStore is nil")
	}
}

func TestProductManagerController_FetchProductList_WithShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewProductManagerController(app)
	r, _ := testutils.NewRequest("GET", "/admin/shop/products", testutils.NewRequestOptions{})
	data := productManagerControllerData{
		request:   r,
		pageInt:   0,
		perPage:   10,
		sortOrder: "desc",
		sortBy:    "created_at",
	}

	productList, productCount, err := controller.fetchProductList(data)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if productList == nil {
		t.Error("Expected non-nil productList")
	}
	if productCount < 0 {
		t.Errorf("Expected productCount to be >= 0, got %d", productCount)
	}
}
