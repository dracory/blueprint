package aiposteditor

import (
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
