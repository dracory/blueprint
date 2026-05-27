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

// TestDiscountControllerFuncUpdateValidation_MissingTitle verifies missing title error
func TestDiscountControllerFuncUpdateValidation_MissingTitle(t *testing.T) {
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

	r = httptest.NewRequest("PUT", "/", nil)
	data := map[string]string{}
	err = controller.FuncUpdate(r, discountID, data)
	if err == nil {
		t.Error("FuncUpdate() should return error for missing title")
	}
	if err.Error() != "title is required" {
		t.Errorf("FuncUpdate() returned unexpected error: got %v, want title is required", err)
	}
}

// TestDiscountControllerFuncUpdateValidation_MissingStatus verifies missing status error
func TestDiscountControllerFuncUpdateValidation_MissingStatus(t *testing.T) {
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

	r = httptest.NewRequest("PUT", "/", nil)
	data := map[string]string{
		"title": "Test",
	}
	err = controller.FuncUpdate(r, discountID, data)
	if err == nil {
		t.Error("FuncUpdate() should return error for missing status")
	}
	if err.Error() != "status is required" {
		t.Errorf("FuncUpdate() returned unexpected error: got %v, want status is required", err)
	}
}

// TestDiscountControllerFuncUpdateValidation_MissingCode verifies missing code error
func TestDiscountControllerFuncUpdateValidation_MissingCode(t *testing.T) {
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

	r = httptest.NewRequest("PUT", "/", nil)
	data := map[string]string{
		"title":  "Test",
		"status": shopstore.DISCOUNT_STATUS_DRAFT,
	}
	err = controller.FuncUpdate(r, discountID, data)
	if err == nil {
		t.Error("FuncUpdate() should return error for missing code")
	}
	if err.Error() != "code is required" {
		t.Errorf("FuncUpdate() returned unexpected error: got %v, want code is required", err)
	}
}

// TestDiscountControllerFuncUpdateValidation_MissingType verifies missing type error
func TestDiscountControllerFuncUpdateValidation_MissingType(t *testing.T) {
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

	r = httptest.NewRequest("PUT", "/", nil)
	data := map[string]string{
		"title":  "Test",
		"status": shopstore.DISCOUNT_STATUS_DRAFT,
		"code":   "TEST",
	}
	err = controller.FuncUpdate(r, discountID, data)
	if err == nil {
		t.Error("FuncUpdate() should return error for missing type")
	}
	if err.Error() != "discount type is required" {
		t.Errorf("FuncUpdate() returned unexpected error: got %v, want discount type is required", err)
	}
}

// TestDiscountControllerFuncUpdateValidation_MissingStartsAt verifies missing starts_at error
func TestDiscountControllerFuncUpdateValidation_MissingStartsAt(t *testing.T) {
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

	r = httptest.NewRequest("PUT", "/", nil)
	data := map[string]string{
		"title":  "Test",
		"status": shopstore.DISCOUNT_STATUS_DRAFT,
		"code":   "TEST",
		"type":   shopstore.DISCOUNT_TYPE_PERCENT,
	}
	err = controller.FuncUpdate(r, discountID, data)
	if err == nil {
		t.Error("FuncUpdate() should return error for missing starts_at")
	}
	if err.Error() != "starts_at is required" {
		t.Errorf("FuncUpdate() returned unexpected error: got %v, want starts_at is required", err)
	}
}

// TestDiscountControllerFuncUpdateValidation_MissingEndsAt verifies missing ends_at error
func TestDiscountControllerFuncUpdateValidation_MissingEndsAt(t *testing.T) {
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

	r = httptest.NewRequest("PUT", "/", nil)
	data := map[string]string{
		"title":     "Test",
		"status":    shopstore.DISCOUNT_STATUS_DRAFT,
		"code":      "TEST",
		"type":      shopstore.DISCOUNT_TYPE_PERCENT,
		"starts_at": "2026-01-01 00:00:00",
	}
	err = controller.FuncUpdate(r, discountID, data)
	if err == nil {
		t.Error("FuncUpdate() should return error for missing ends_at")
	}
	if err.Error() != "ends_at is required" {
		t.Errorf("FuncUpdate() returned unexpected error: got %v, want ends_at is required", err)
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
