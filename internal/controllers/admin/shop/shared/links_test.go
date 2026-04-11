package shared

import (
	"strings"
	"testing"
)

// TestNewLinks verifies Links can be created
func TestNewLinks(t *testing.T) {
	t.Parallel()
	links := NewLinks()

	if links == nil {
		t.Error("NewLinks() returned nil")
	}
}

// TestLinksHome verifies Home method returns correct URL with controller parameter
func TestLinksHome(t *testing.T) {
	t.Parallel()
	links := NewLinks()
	result := links.Home(map[string]string{})

	if result == "" {
		t.Error("Home() returned empty string")
	}
	if !strings.Contains(result, "controller=home") {
		t.Errorf("Home() URL should contain controller=home, got: %s", result)
	}
}

// TestLinksCategories verifies Categories method returns correct URL with controller parameter
func TestLinksCategories(t *testing.T) {
	t.Parallel()
	links := NewLinks()
	result := links.Categories(map[string]string{})

	if result == "" {
		t.Error("Categories() returned empty string")
	}
	if !strings.Contains(result, "controller=categories") {
		t.Errorf("Categories() URL should contain controller=categories, got: %s", result)
	}
}

// TestLinksCategoryCreate verifies CategoryCreate method returns correct URL with controller parameter
func TestLinksCategoryCreate(t *testing.T) {
	t.Parallel()
	links := NewLinks()
	result := links.CategoryCreate(map[string]string{})

	if result == "" {
		t.Error("CategoryCreate() returned empty string")
	}
	if !strings.Contains(result, "controller=category_create") {
		t.Errorf("CategoryCreate() URL should contain controller=category_create, got: %s", result)
	}
}

// TestLinksCategoryUpdate verifies CategoryUpdate method returns correct URL with controller and custom params
func TestLinksCategoryUpdate(t *testing.T) {
	t.Parallel()
	links := NewLinks()
	result := links.CategoryUpdate(map[string]string{"category_id": "123", "action": "edit"})

	if result == "" {
		t.Error("CategoryUpdate() returned empty string")
	}
	if !strings.Contains(result, "controller=category_update") {
		t.Errorf("CategoryUpdate() URL should contain controller=category_update, got: %s", result)
	}
	if !strings.Contains(result, "category_id=123") {
		t.Errorf("CategoryUpdate() URL should contain category_id=123, got: %s", result)
	}
	if !strings.Contains(result, "action=edit") {
		t.Errorf("CategoryUpdate() URL should contain action=edit, got: %s", result)
	}
}

// TestLinksDiscounts verifies Discounts method returns correct URL with controller parameter
func TestLinksDiscounts(t *testing.T) {
	t.Parallel()
	links := NewLinks()
	result := links.Discounts(map[string]string{})

	if result == "" {
		t.Error("Discounts() returned empty string")
	}
	if !strings.Contains(result, "controller=discounts") {
		t.Errorf("Discounts() URL should contain controller=discounts, got: %s", result)
	}
}

// TestLinksOrders verifies Orders method returns correct URL with controller parameter
func TestLinksOrders(t *testing.T) {
	t.Parallel()
	links := NewLinks()
	result := links.Orders(map[string]string{})

	if result == "" {
		t.Error("Orders() returned empty string")
	}
	if !strings.Contains(result, "controller=orders") {
		t.Errorf("Orders() URL should contain controller=orders, got: %s", result)
	}
}

// TestLinksProductCreate verifies ProductCreate method returns correct URL with controller parameter
func TestLinksProductCreate(t *testing.T) {
	t.Parallel()
	links := NewLinks()
	result := links.ProductCreate(map[string]string{})

	if result == "" {
		t.Error("ProductCreate() returned empty string")
	}
	if !strings.Contains(result, "controller=product_create") {
		t.Errorf("ProductCreate() URL should contain controller=product_create, got: %s", result)
	}
}

// TestLinksProductDelete verifies ProductDelete method returns correct URL with controller parameter
func TestLinksProductDelete(t *testing.T) {
	t.Parallel()
	links := NewLinks()
	result := links.ProductDelete(map[string]string{})

	if result == "" {
		t.Error("ProductDelete() returned empty string")
	}
	if !strings.Contains(result, "controller=product_delete") {
		t.Errorf("ProductDelete() URL should contain controller=product_delete, got: %s", result)
	}
}

// TestLinksProducts verifies Products method returns correct URL with controller parameter
func TestLinksProducts(t *testing.T) {
	t.Parallel()
	links := NewLinks()
	result := links.Products(map[string]string{})

	if result == "" {
		t.Error("Products() returned empty string")
	}
	if !strings.Contains(result, "controller=products") {
		t.Errorf("Products() URL should contain controller=products, got: %s", result)
	}
}

// TestLinksProductUpdate verifies ProductUpdate method returns correct URL with controller parameter
func TestLinksProductUpdate(t *testing.T) {
	t.Parallel()
	links := NewLinks()
	result := links.ProductUpdate(map[string]string{})

	if result == "" {
		t.Error("ProductUpdate() returned empty string")
	}
	if !strings.Contains(result, "controller=product_update") {
		t.Errorf("ProductUpdate() URL should contain controller=product_update, got: %s", result)
	}
}
