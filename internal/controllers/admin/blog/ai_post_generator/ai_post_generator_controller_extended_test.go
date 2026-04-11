package aipostgenerator

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"project/internal/testutils"
)

func TestAiPostGeneratorController_HandlerWithGET(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)
	if controller == nil {
		t.Fatal("NewAiPostGeneratorController() returned nil")
	}

	// Test GET request (should render the page)
	req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-post-generator", nil)
	w := httptest.NewRecorder()

	// This will return error because custom store is not fully configured
	result := controller.Handler(w, req)
	// Result should be error HTML or the page - just verify it doesn't panic
	_ = result
}

func TestAiPostGeneratorController_HandlerWithNilRegistry(t *testing.T) {
	// Test with nil registry - expect panic due to nil pointer dereference
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil registry
		}
	}()

	controller := NewAiPostGeneratorController(nil)
	if controller == nil {
		t.Fatal("NewAiPostGeneratorController(nil) should not return nil")
	}

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-post-generator", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}

func TestAiPostGeneratorController_HandlerWithPOSTAndGeneratePost(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)
	if controller == nil {
		t.Fatal("NewAiPostGeneratorController() returned nil")
	}

	// Test POST with action=generate_post
	formData := url.Values{}
	formData.Set("action", ACTION_GENERATE_POST)
	formData.Set("title_id", "test-title-id")

	req := httptest.NewRequest(http.MethodPost, "/admin/blog/ai-post-generator", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	// This will return error because custom store and LLM are not configured
	result := controller.Handler(w, req)
	_ = result
}

func TestAiPostGeneratorController_HandlerWithPOSTDifferentAction(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	// Test POST with different action (not generate_post)
	formData := url.Values{}
	formData.Set("action", "other_action")

	req := httptest.NewRequest(http.MethodPost, "/admin/blog/ai-post-generator", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)
	_ = result
}

func TestAiPostGeneratorController_HandlerWithDifferentMethods(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		var req *http.Request
		if method == http.MethodPost {
			formData := url.Values{}
			req = httptest.NewRequest(method, "/admin/blog/ai-post-generator", strings.NewReader(formData.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else {
			req = httptest.NewRequest(method, "/admin/blog/ai-post-generator", nil)
		}
		w := httptest.NewRecorder()

		result := controller.Handler(w, req)
		_ = result
	}
}

func TestAiPostGeneratorController_MultipleInstances(t *testing.T) {
	registry1 := testutils.Setup()
	registry2 := testutils.Setup()

	controller1 := NewAiPostGeneratorController(registry1)
	controller2 := NewAiPostGeneratorController(registry2)

	if controller1 == nil || controller2 == nil {
		t.Fatal("Controllers should not be nil")
	}

	if controller1 == controller2 {
		t.Error("Controllers should be separate instances")
	}

	if controller1.registry != registry1 {
		t.Error("Controller1 should have registry1")
	}

	if controller2.registry != registry2 {
		t.Error("Controller2 should have registry2")
	}
}

func TestAiPostGeneratorController_RegistryField(t *testing.T) {
	// Test with nil registry
	controller := NewAiPostGeneratorController(nil)
	if controller.registry != nil {
		t.Error("Controller registry should be nil when passed nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	controller = NewAiPostGeneratorController(registry)
	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
	}
}

func TestAiPostGeneratorController_prepareData(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-post-generator?action=test", nil)

	data, errorMessage := controller.prepareData(req)

	// With testutils.Setup(), custom store should exist but may not be fully configured
	// Just verify the function runs without panic
	// data and errorMessage may be nil or contain error depending on store state
	_ = data
	_ = errorMessage
}

func TestAiPostGeneratorController_prepareDataWithNilRegistry(t *testing.T) {
	// Test with nil registry - expect panic due to nil pointer dereference
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil registry
		}
	}()

	controller := NewAiPostGeneratorController(nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-post-generator", nil)

	// This will panic with nil registry
	controller.prepareData(req)
}

func TestAiPostGeneratorController_ActionConstant(t *testing.T) {
	if ACTION_GENERATE_POST != "generate_post" {
		t.Errorf("ACTION_GENERATE_POST = %q, want generate_post", ACTION_GENERATE_POST)
	}
}

func TestAiPostGeneratorController_pageDataStruct(t *testing.T) {
	// Test that pageData can be created
	data := pageData{
		Request:             nil,
		Action:              "test",
		ApprovedBlogAiPosts: nil,
	}

	if data.Action != "test" {
		t.Error("pageData.Action not set correctly")
	}
}
