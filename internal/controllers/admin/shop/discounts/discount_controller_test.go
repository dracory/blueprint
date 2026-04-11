package admin

import (
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/shopstore"
)

// TestNewDiscountController verifies controller can be created
func TestNewDiscountController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)

	if controller == nil {
		t.Error("NewDiscountController() returned nil")
	}
}

// TestDiscountControllerRegistry verifies controller has registry
func TestDiscountControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestDiscountControllerAnyIndexExists verifies AnyIndex method exists
func TestDiscountControllerAnyIndexExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)

	// Verify AnyIndex method exists (should compile without error)
	_ = controller.AnyIndex
}

// TestDiscountControllerFuncLayout verifies layout function returns HTML
func TestDiscountControllerFuncLayout(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)

	html := controller.FuncLayout(w, r, "Test Title", "Test Content", []string{}, "", []string{}, "")

	if html == "" {
		t.Error("FuncLayout() returned empty string")
	}

	// Verify HTML contains expected elements
	if !strings.Contains(strings.ToLower(html), "<html") {
		t.Error("FuncLayout() did not return HTML")
	}
}

// TestDiscountControllerFuncRows verifies rows function returns discount data
func TestDiscountControllerFuncRows(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)
	r := httptest.NewRequest("GET", "/", nil)

	// Test with no discounts
	rows, err := controller.FuncRows(r)
	if err != nil {
		t.Errorf("FuncRows() returned error: %v", err)
	}

	if rows == nil {
		t.Error("FuncRows() returned nil")
	}
}

// TestDiscountControllerFuncRowsNoShopStore verifies error when shop store not configured
func TestDiscountControllerFuncRowsNoShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(false))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)
	r := httptest.NewRequest("GET", "/", nil)

	_, err := controller.FuncRows(r)
	if err == nil {
		t.Error("FuncRows() should return error when shop store not configured")
	}

	if err.Error() != "shop store not configured" {
		t.Errorf("FuncRows() returned unexpected error: %v", err)
	}
}

// TestDiscountControllerFuncCreate verifies create function creates a discount
func TestDiscountControllerFuncCreate(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)
	r := httptest.NewRequest("POST", "/", nil)

	data := map[string]string{
		"title": "Test Discount",
	}

	discountID, err := controller.FuncCreate(r, data)
	if err != nil {
		t.Errorf("FuncCreate() returned error: %v", err)
	}

	if discountID == "" {
		t.Error("FuncCreate() returned empty discount ID")
	}
}

// TestDiscountControllerFuncCreateNoShopStore verifies error when shop store not configured
func TestDiscountControllerFuncCreateNoShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(false))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)
	r := httptest.NewRequest("POST", "/", nil)

	data := map[string]string{
		"title": "Test Discount",
	}

	_, err := controller.FuncCreate(r, data)
	if err == nil {
		t.Error("FuncCreate() should return error when shop store not configured")
	}

	if err.Error() != "shop store not configured" {
		t.Errorf("FuncCreate() returned unexpected error: %v", err)
	}
}

// TestDiscountControllerFuncUpdate verifies update function updates a discount
func TestDiscountControllerFuncUpdate(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)

	// First create a discount
	r := httptest.NewRequest("POST", "/", nil)
	createData := map[string]string{
		"title": "Test Discount",
	}
	discountID, err := controller.FuncCreate(r, createData)
	if err != nil {
		t.Fatalf("FuncCreate() failed: %v", err)
	}

	// Update the discount
	r = httptest.NewRequest("PUT", "/", nil)
	updateData := map[string]string{
		"title":       "Updated Discount",
		"status":      shopstore.DISCOUNT_STATUS_DRAFT,
		"type":        shopstore.DISCOUNT_TYPE_PERCENT,
		"code":        "TESTCODE",
		"starts_at":   "2026-01-01 00:00:00",
		"ends_at":     "2026-12-31 23:59:59",
		"amount":      "10.5",
		"description": "Test description",
	}

	err = controller.FuncUpdate(r, discountID, updateData)
	if err != nil {
		t.Errorf("FuncUpdate() returned error: %v", err)
	}
}

// TestDiscountControllerFuncUpdateNoShopStore verifies error when shop store not configured
func TestDiscountControllerFuncUpdateNoShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(false))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)
	r := httptest.NewRequest("PUT", "/", nil)

	data := map[string]string{
		"title": "Test Discount",
	}

	err := controller.FuncUpdate(r, "test-id", data)
	if err == nil {
		t.Error("FuncUpdate() should return error when shop store not configured")
	}

	if err.Error() != "shop store not configured" {
		t.Errorf("FuncUpdate() returned unexpected error: %v", err)
	}
}

// TestDiscountControllerFuncUpdateValidation verifies validation errors
func TestDiscountControllerFuncUpdateValidation(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)

	// First create a discount to update
	r := httptest.NewRequest("POST", "/", nil)
	createData := map[string]string{
		"title": "Test Discount",
	}
	discountID, err := controller.FuncCreate(r, createData)
	if err != nil {
		t.Fatalf("FuncCreate() failed: %v", err)
	}

	tests := []struct {
		name    string
		data    map[string]string
		wantErr string
	}{
		{
			name:    "missing title",
			data:    map[string]string{},
			wantErr: "title is required",
		},
		{
			name: "missing status",
			data: map[string]string{
				"title": "Test",
			},
			wantErr: "status is required",
		},
		{
			name: "missing code",
			data: map[string]string{
				"title":  "Test",
				"status": shopstore.DISCOUNT_STATUS_DRAFT,
			},
			wantErr: "code is required",
		},
		{
			name: "missing type",
			data: map[string]string{
				"title":  "Test",
				"status": shopstore.DISCOUNT_STATUS_DRAFT,
				"code":   "TEST",
			},
			wantErr: "discount type is required",
		},
		{
			name: "missing starts_at",
			data: map[string]string{
				"title":  "Test",
				"status": shopstore.DISCOUNT_STATUS_DRAFT,
				"code":   "TEST",
				"type":   shopstore.DISCOUNT_TYPE_PERCENT,
			},
			wantErr: "starts_at is required",
		},
		{
			name: "missing ends_at",
			data: map[string]string{
				"title":     "Test",
				"status":    shopstore.DISCOUNT_STATUS_DRAFT,
				"code":      "TEST",
				"type":      shopstore.DISCOUNT_TYPE_PERCENT,
				"starts_at": "2026-01-01 00:00:00",
			},
			wantErr: "ends_at is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest("PUT", "/", nil)
			err := controller.FuncUpdate(r, discountID, tt.data)
			if err == nil {
				t.Errorf("FuncUpdate() should return error for %s", tt.name)
			}
			if err.Error() != tt.wantErr {
				t.Errorf("FuncUpdate() returned unexpected error for %s: got %v, want %s", tt.name, err, tt.wantErr)
			}
		})
	}
}

// TestDiscountControllerFuncFetchReadData verifies fetch read data function
func TestDiscountControllerFuncFetchReadData(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)

	// First create a discount
	r := httptest.NewRequest("POST", "/", nil)
	createData := map[string]string{
		"title": "Test Discount",
	}
	discountID, err := controller.FuncCreate(r, createData)
	if err != nil {
		t.Fatalf("FuncCreate() failed: %v", err)
	}

	// Fetch read data
	r = httptest.NewRequest("GET", "/", nil)
	data, err := controller.FuncFetchReadData(r, discountID)
	if err != nil {
		t.Errorf("FuncFetchReadData() returned error: %v", err)
	}

	if data == nil {
		t.Error("FuncFetchReadData() returned nil")
	}
}

// TestDiscountControllerFuncFetchReadDataNoShopStore verifies error when shop store not configured
func TestDiscountControllerFuncFetchReadDataNoShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(false))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)
	r := httptest.NewRequest("GET", "/", nil)

	_, err := controller.FuncFetchReadData(r, "test-id")
	if err == nil {
		t.Error("FuncFetchReadData() should return error when shop store not configured")
	}

	if err.Error() != "shop store not configured" {
		t.Errorf("FuncFetchReadData() returned unexpected error: %v", err)
	}
}

// TestDiscountControllerFuncFetchUpdateData verifies fetch update data function
func TestDiscountControllerFuncFetchUpdateData(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)

	// First create a discount
	r := httptest.NewRequest("POST", "/", nil)
	createData := map[string]string{
		"title": "Test Discount",
	}
	discountID, err := controller.FuncCreate(r, createData)
	if err != nil {
		t.Fatalf("FuncCreate() failed: %v", err)
	}

	// Fetch update data
	r = httptest.NewRequest("GET", "/", nil)
	data, err := controller.FuncFetchUpdateData(r, discountID)
	if err != nil {
		t.Errorf("FuncFetchUpdateData() returned error: %v", err)
	}

	if data == nil {
		t.Error("FuncFetchUpdateData() returned nil")
	}

	if data["title"] != "Test Discount" {
		t.Errorf("FuncFetchUpdateData() returned unexpected title: got %s, want Test Discount", data["title"])
	}
}

// TestDiscountControllerFuncFetchUpdateDataNoShopStore verifies error when shop store not configured
func TestDiscountControllerFuncFetchUpdateDataNoShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(false))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)
	r := httptest.NewRequest("GET", "/", nil)

	_, err := controller.FuncFetchUpdateData(r, "test-id")
	if err == nil {
		t.Error("FuncFetchUpdateData() should return error when shop store not configured")
	}

	if err.Error() != "shop store not configured" {
		t.Errorf("FuncFetchUpdateData() returned unexpected error: %v", err)
	}
}

// TestDiscountControllerFuncTrash verifies trash function deletes a discount
func TestDiscountControllerFuncTrash(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(true))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)

	// First create a discount
	r := httptest.NewRequest("POST", "/", nil)
	createData := map[string]string{
		"title": "Test Discount",
	}
	discountID, err := controller.FuncCreate(r, createData)
	if err != nil {
		t.Fatalf("FuncCreate() failed: %v", err)
	}

	// Trash the discount
	r = httptest.NewRequest("DELETE", "/", nil)
	err = controller.FuncTrash(r, discountID)
	if err != nil {
		t.Errorf("FuncTrash() returned error: %v", err)
	}
}

// TestDiscountControllerFuncTrashNoShopStore verifies error when shop store not configured
func TestDiscountControllerFuncTrashNoShopStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithShopStore(false))
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewDiscountController(app)
	r := httptest.NewRequest("DELETE", "/", nil)

	err := controller.FuncTrash(r, "test-id")
	if err == nil {
		t.Error("FuncTrash() should return error when shop store not configured")
	}

	if err.Error() != "shop store not configured" {
		t.Errorf("FuncTrash() returned unexpected error: %v", err)
	}
}
