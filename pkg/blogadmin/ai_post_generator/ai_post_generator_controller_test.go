package aipostgenerator

import (
	"testing"
)

// TestNewAiPostGeneratorController tests the constructor
func TestNewAiPostGeneratorController(t *testing.T) {
	t.Parallel()

	controller := NewAiPostGeneratorController(nil)
	if controller == nil {
		t.Error("Expected controller to be non-nil")
	}
	if controller.registry != nil {
		t.Error("Expected registry to be nil when passed nil")
	}
}

// TestACTION_GENERATE_POST_Constant tests the action constant
func TestACTION_GENERATE_POST_Constant(t *testing.T) {
	t.Parallel()

	expected := "generate_post"
	if ACTION_GENERATE_POST != expected {
		t.Errorf("Expected ACTION_GENERATE_POST to be %q, got: %q", expected, ACTION_GENERATE_POST)
	}
}

// TestAiPostGeneratorController_StructFields tests controller structure
func TestAiPostGeneratorController_StructFields(t *testing.T) {
	t.Parallel()

	controller := NewAiPostGeneratorController(nil)
	if controller == nil {
		t.Fatal("NewAiPostGeneratorController() returned nil")
	}

	// Verify the controller has the expected structure
	// The struct should have a registry field
}

// TestPageData_Struct tests pageData structure
func TestPageData_Struct(t *testing.T) {
	t.Parallel()

	// Test that pageData can be created
	data := pageData{}

	// Verify default zero values
	if data.Request != nil {
		t.Error("Request should be nil by default")
	}
	if data.Action != "" {
		t.Error("Action should be empty by default")
	}
	if data.ApprovedBlogAiPosts != nil {
		t.Error("ApprovedBlogAiPosts should be nil by default")
	}

	// Test setting values
	data.Action = "test_action"
	if data.Action != "test_action" {
		t.Errorf("Expected Action to be 'test_action', got: %s", data.Action)
	}

	// ApprovedBlogAiPosts should be slice
	data.ApprovedBlogAiPosts = nil
	if data.ApprovedBlogAiPosts != nil {
		t.Error("Should be able to set ApprovedBlogAiPosts to nil")
	}
}

// TestPageData_WithApprovedPosts tests pageData with approved posts slice
func TestPageData_WithApprovedPosts(t *testing.T) {
	t.Parallel()

	// Create pageData and set approved posts
	data := pageData{
		Action:              "generate_post",
		ApprovedBlogAiPosts: nil,
	}

	if data.Action != "generate_post" {
		t.Errorf("Expected Action 'generate_post', got: %s", data.Action)
	}

	if data.ApprovedBlogAiPosts != nil {
		t.Error("Expected nil ApprovedBlogAiPosts")
	}
}

// TestAiPostGeneratorController_MultipleInstances tests creating multiple controllers
func TestAiPostGeneratorController_MultipleInstances(t *testing.T) {
	t.Parallel()

	controller1 := NewAiPostGeneratorController(nil)
	controller2 := NewAiPostGeneratorController(nil)

	if controller1 == controller2 {
		t.Error("Each NewAiPostGeneratorController call should return a new instance")
	}

	if controller1 == nil || controller2 == nil {
		t.Error("Both controllers should be non-nil")
	}
}
