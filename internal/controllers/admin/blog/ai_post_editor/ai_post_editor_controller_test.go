package aiposteditor

import (
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

func TestNewAiPostEditorController(t *testing.T) {
	// Test with nil registry
	controller := NewAiPostEditorController(nil)
	if controller == nil {
		t.Error("NewAiPostEditorController() should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	controller = NewAiPostEditorController(registry)
	if controller == nil {
		t.Error("NewAiPostEditorController() should not return nil")
	}
	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
	}
}

func TestConstants(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"ACTION_REGENERATE_SECTION", ACTION_REGENERATE_SECTION},
		{"ACTION_REGENERATE_IMAGE", ACTION_REGENERATE_IMAGE},
		{"ACTION_CREATE_FINAL_POST", ACTION_CREATE_FINAL_POST},
		{"ACTION_SAVE_DRAFT", ACTION_SAVE_DRAFT},
		{"ACTION_REGENERATE_PARAGRAPH", ACTION_REGENERATE_PARAGRAPH},
		{"ACTION_LOAD_POST", ACTION_LOAD_POST},
		{"ACTION_REGENERATE_SUMMARY", ACTION_REGENERATE_SUMMARY},
		{"ACTION_REGENERATE_METAS", ACTION_REGENERATE_METAS},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("%s should not be empty", tt.name)
			}
		})
	}
}

func TestAiPostEditorController_Handler_DefaultAction(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostEditorController(registry)

	req := httptest.NewRequest("GET", "/admin/blog/ai-post-editor", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)
	if result == "" {
		t.Error("Handler() should return non-empty string")
	}
}

func TestAiPostEditorController_Handler_WithActions(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostEditorController(registry)

	actions := []string{
		ACTION_REGENERATE_SECTION,
		ACTION_REGENERATE_IMAGE,
		ACTION_REGENERATE_PARAGRAPH,
		ACTION_CREATE_FINAL_POST,
		ACTION_SAVE_DRAFT,
		ACTION_LOAD_POST,
		ACTION_REGENERATE_SUMMARY,
		ACTION_REGENERATE_METAS,
	}

	for _, action := range actions {
		t.Run(action, func(t *testing.T) {
			url := "/admin/blog/ai-post-editor?action=" + action
			req := httptest.NewRequest("POST", url, nil)
			w := httptest.NewRecorder()

			result := controller.Handler(w, req)
			_ = result
		})
	}
}

func TestAiPostEditorController_Handler_NilRegistry(t *testing.T) {
	controller := NewAiPostEditorController(nil)

	req := httptest.NewRequest("GET", "/admin/blog/ai-post-editor", nil)
	w := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Logf("Handler() with nil registry panicked as expected: %v", r)
		}
	}()

	controller.Handler(w, req)
}

func TestAiPostEditorController_NilRegistry(t *testing.T) {
	controller := NewAiPostEditorController(nil)
	if controller == nil {
		t.Error("NewAiPostEditorController(nil) should not return nil")
	}
	if controller.registry != nil {
		t.Error("Controller registry should be nil when passed nil")
	}
}
