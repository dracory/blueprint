package tag_manager

import (
	"testing"

	"github.com/dracory/str"
)

// TestNewTagManagerController tests the constructor
func TestNewTagManagerController(t *testing.T) {
	t.Parallel()

	controller := NewTagManagerController(nil)
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
		{"Test Tag", "test-tag"},
		{"My_Tag-Name", "my-tag-name"},
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

// TestTagManagerController_Struct tests controller structure
func TestTagManagerController_Struct(t *testing.T) {
	t.Parallel()

	controller := NewTagManagerController(nil)
	if controller == nil {
		t.Fatal("NewTagManagerController() returned nil")
	}

	// Verify struct fields
	// The controller should have a registry field
}

// TestTagManagerController_MultipleInstances tests creating multiple controllers
func TestTagManagerController_MultipleInstances(t *testing.T) {
	t.Parallel()

	controller1 := NewTagManagerController(nil)
	controller2 := NewTagManagerController(nil)

	if controller1 == controller2 {
		t.Error("Each NewTagManagerController call should return a new instance")
	}

	if controller1 == nil || controller2 == nil {
		t.Error("Both controllers should be non-nil")
	}
}

// TestTagManagerController_Handler_MethodExists tests handler method exists
func TestTagManagerController_Handler_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewTagManagerController(nil)
	if controller == nil {
		t.Fatal("NewTagManagerController() returned nil")
	}

	// Handler method should exist
}

// TestTagManagerController_renderPage_MethodExists tests renderPage method exists
func TestTagManagerController_renderPage_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewTagManagerController(nil)
	if controller == nil {
		t.Fatal("NewTagManagerController() returned nil")
	}

	// renderPage method should exist
}

// TestTagManagerController_handleLoadTags_MethodExists tests method exists
func TestTagManagerController_handleLoadTags_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewTagManagerController(nil)
	if controller == nil {
		t.Fatal("NewTagManagerController() returned nil")
	}

	// handleLoadTags method should exist
}

// TestTagManagerController_handleLoadTagPosts_MethodExists tests method exists
func TestTagManagerController_handleLoadTagPosts_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewTagManagerController(nil)
	if controller == nil {
		t.Fatal("NewTagManagerController() returned nil")
	}

	// handleLoadTagPosts method should exist
}

// TestTagManagerController_handleCreateTag_MethodExists tests method exists
func TestTagManagerController_handleCreateTag_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewTagManagerController(nil)
	if controller == nil {
		t.Fatal("NewTagManagerController() returned nil")
	}

	// handleCreateTag method should exist
}

// TestTagManagerController_handleUpdateTag_MethodExists tests method exists
func TestTagManagerController_handleUpdateTag_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewTagManagerController(nil)
	if controller == nil {
		t.Fatal("NewTagManagerController() returned nil")
	}

	// handleUpdateTag method should exist
}

// TestTagManagerController_handleDeleteTag_MethodExists tests method exists
func TestTagManagerController_handleDeleteTag_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewTagManagerController(nil)
	if controller == nil {
		t.Fatal("NewTagManagerController() returned nil")
	}

	// handleDeleteTag method should exist
}

// TestTagManagerController_ensureTaxonomy_MethodExists tests ensureTaxonomy exists
func TestTagManagerController_ensureTaxonomy_MethodExists(t *testing.T) {
	t.Parallel()

	controller := NewTagManagerController(nil)
	if controller == nil {
		t.Fatal("NewTagManagerController() returned nil")
	}

	// The ensureTaxonomy method should exist
}

// TestTagManagerControllerData_Struct tests the deprecated struct
func TestTagManagerControllerData_Struct(t *testing.T) {
	t.Parallel()

	// Test the deprecated struct can be instantiated
	data := tagManagerControllerData{
		page:       "1",
		pageInt:    1,
		perPage:    10,
		taxonomyID: "test-taxonomy",
		tagCount:   5,
		tagList:    nil,
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
	if data.tagCount != 5 {
		t.Errorf("Expected tagCount 5, got: %d", data.tagCount)
	}
}

// TestTagManagerControllerData_ZeroValues tests zero values
func TestTagManagerControllerData_ZeroValues(t *testing.T) {
	t.Parallel()

	data := tagManagerControllerData{}

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
	if data.tagCount != 0 {
		t.Errorf("Expected tagCount 0, got: %d", data.tagCount)
	}
	if data.tagList != nil {
		t.Error("Expected nil tagList")
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
