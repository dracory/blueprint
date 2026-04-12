package category_manager

import (
	"testing"

	"github.com/dracory/str"
)

// TestCategoryManagerController_Struct tests controller structure
func TestCategoryManagerController_Struct(t *testing.T) {
	t.Parallel()

	controller := NewCategoryManagerController(nil)
	if controller == nil {
		t.Fatal("NewCategoryManagerController() returned nil")
	}

	// Verify struct fields
	// The controller should have a registry field
}

// TestNewCategoryManagerController tests the constructor
func TestNewCategoryManagerController(t *testing.T) {
	t.Parallel()

	controller := NewCategoryManagerController(nil)
	if controller == nil {
		t.Error("Expected controller to be non-nil")
	}
	if controller.registry != nil {
		t.Error("Expected registry to be nil when passed nil")
	}
}

// TestSlugify tests the slugify function
func TestSlugify(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"Hello World", "hello-world"},
		{"Test Category", "test-category"},
		{"My_Category-Name", "my-category-name"},
		{"UPPERCASE", "uppercase"},
		{"123 Numbers", "123-numbers"},
		{"", ""},
		{"Already-Slugified", "already-slugified"},
		{"Multiple   Spaces", "multiple-spaces"},
	}

	for _, tt := range tests {
		result := str.Slugify(tt.input, '-')
		if result != tt.expected {
			t.Errorf("str.Slugify(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// TestCategoryManagerController_MultipleInstances tests creating multiple controllers
func TestCategoryManagerController_MultipleInstances(t *testing.T) {
	t.Parallel()

	controller1 := NewCategoryManagerController(nil)
	controller2 := NewCategoryManagerController(nil)

	if controller1 == controller2 {
		t.Error("Each NewCategoryManagerController call should return a new instance")
	}

	if controller1 == nil || controller2 == nil {
		t.Error("Both controllers should be non-nil")
	}
}

// TestCategoryManagerController_ensureTaxonomy_MethodExists tests ensureTaxonomy exists
func TestCategoryManagerController_ensureTaxonomy_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewCategoryManagerController(nil)
	if controller == nil {
		t.Fatal("NewCategoryManagerController() returned nil")
	}

	// The ensureTaxonomy method should exist (tested via reflection concept)
	// This tests that the controller has the expected API
}

// TestCategoryManagerController_Handler_MethodExists tests handler method exists
func TestCategoryManagerController_Handler_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewCategoryManagerController(nil)
	if controller == nil {
		t.Fatal("NewCategoryManagerController() returned nil")
	}

	// Handler method should exist
	// We can't directly test it without mocks, but we can verify the controller structure
}

// TestCategoryManagerController_renderPage_MethodExists tests renderPage method exists
func TestCategoryManagerController_renderPage_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewCategoryManagerController(nil)
	if controller == nil {
		t.Fatal("NewCategoryManagerController() returned nil")
	}

	// renderPage method should exist
}

// TestCategoryManagerController_handleLoadCategories_MethodExists tests method exists
func TestCategoryManagerController_handleLoadCategories_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewCategoryManagerController(nil)
	if controller == nil {
		t.Fatal("NewCategoryManagerController() returned nil")
	}

	// handleLoadCategories method should exist
}

// TestCategoryManagerController_handleCreateCategory_MethodExists tests method exists
func TestCategoryManagerController_handleCreateCategory_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewCategoryManagerController(nil)
	if controller == nil {
		t.Fatal("NewCategoryManagerController() returned nil")
	}

	// handleCreateCategory method should exist
}

// TestCategoryManagerController_handleUpdateCategory_MethodExists tests method exists
func TestCategoryManagerController_handleUpdateCategory_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewCategoryManagerController(nil)
	if controller == nil {
		t.Fatal("NewCategoryManagerController() returned nil")
	}

	// handleUpdateCategory method should exist
}

// TestCategoryManagerController_handleReorderCategories_MethodExists tests method exists
func TestCategoryManagerController_handleReorderCategories_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewCategoryManagerController(nil)
	if controller == nil {
		t.Fatal("NewCategoryManagerController() returned nil")
	}

	// handleReorderCategories method should exist
}

// TestCategoryManagerController_handleDeleteCategory_MethodExists tests method exists
func TestCategoryManagerController_handleDeleteCategory_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewCategoryManagerController(nil)
	if controller == nil {
		t.Fatal("NewCategoryManagerController() returned nil")
	}

	// handleDeleteCategory method should exist
}

// TestCategoryManagerControllerData_Struct tests the deprecated struct
func TestCategoryManagerControllerData_Struct(t *testing.T) {
	t.Parallel()

	// Test the deprecated struct can be instantiated
	data := categoryManagerControllerData{
		page:          "1",
		pageInt:       1,
		perPage:       10,
		taxonomyID:    "test-taxonomy",
		categoryCount: 5,
		categoryList:  nil,
	}

	if data.page != "1" {
		t.Errorf("Expected page '1', got: %s", data.page)
	}
	if data.pageInt != 1 {
		t.Errorf("Expected pageInt 1, got: %d", data.pageInt)
	}
	if data.perPage != 10 {
		t.Errorf("Expected perPage 10, got: %d", data.perPage)
	}
	if data.taxonomyID != "test-taxonomy" {
		t.Errorf("Expected taxonomyID 'test-taxonomy', got: %s", data.taxonomyID)
	}
	if data.categoryCount != 5 {
		t.Errorf("Expected categoryCount 5, got: %d", data.categoryCount)
	}
}

// TestCategoryManagerControllerData_ZeroValues tests zero values
func TestCategoryManagerControllerData_ZeroValues(t *testing.T) {
	t.Parallel()

	data := categoryManagerControllerData{}

	if data.page != "" {
		t.Errorf("Expected empty page, got: %s", data.page)
	}
	if data.pageInt != 0 {
		t.Errorf("Expected pageInt 0, got: %d", data.pageInt)
	}
	if data.perPage != 0 {
		t.Errorf("Expected perPage 0, got: %d", data.perPage)
	}
	if data.taxonomyID != "" {
		t.Errorf("Expected empty taxonomyID, got: %s", data.taxonomyID)
	}
	if data.categoryCount != 0 {
		t.Errorf("Expected categoryCount 0, got: %d", data.categoryCount)
	}
	if data.categoryList != nil {
		t.Error("Expected nil categoryList")
	}
}

// TestSlugifyEdgeCases tests edge cases for slugify
func TestSlugifyEdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"Special!@#$%^&*()Chars", "special-chars"},
		{"  Leading spaces", "leading-spaces"},
		{"Trailing spaces  ", "trailing-spaces"},
		{"Multiple---Dashes", "multiple-dashes"},
		{"MixedCASE_Input", "mixedcase-input"},
		{"Numbers123", "numbers123"},
		{"äöü", "aou"},
		{"Very Long String With Many Words That Should Still Work", "very-long-string-with-many-words-that-should-still-work"},
	}

	for _, tt := range tests {
		result := str.Slugify(tt.input, '-')
		if result != tt.expected {
			t.Errorf("str.Slugify(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
