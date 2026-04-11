package aititlegenerator

import (
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

func TestNewAiTitleGeneratorController(t *testing.T) {
	// Test with nil registry
	controller := NewAiTitleGeneratorController(nil)
	if controller == nil {
		t.Error("NewAiTitleGeneratorController() should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	controller = NewAiTitleGeneratorController(registry)
	if controller == nil {
		t.Error("NewAiTitleGeneratorController() should not return nil")
	}
	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
	}
}

func TestAiTitleGeneratorController_Handler_DefaultAction(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiTitleGeneratorController(registry)

	req := httptest.NewRequest("GET", "/admin/blog/ai-title-generator", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}

func TestAiTitleGeneratorController_Handler_WithActions(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiTitleGeneratorController(registry)

	tests := []struct {
		name   string
		action string
		method string
	}{
		{"add_title_get", ACTION_ADD_TITLE, "GET"},
		{"add_title_post", ACTION_ADD_TITLE, "POST"},
		{"generate_titles", ACTION_GENERATE_TITLES, "POST"},
		{"approve_title", ACTION_APPROVE_TITLE, "POST"},
		{"reject_title", ACTION_REJECT_TITLE, "POST"},
		{"delete_title", ACTION_DELETE_TITLE, "POST"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/admin/blog/ai-title-generator?action=" + tt.action
			req := httptest.NewRequest(tt.method, url, nil)
			w := httptest.NewRecorder()

			controller.Handler(w, req)
		})
	}
}

func TestAiTitleGeneratorController_HandlerMultipleCalls(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiTitleGeneratorController(registry)

	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/admin/blog/ai-title-generator", nil)
		w := httptest.NewRecorder()

		controller.Handler(w, req)
	}
}
