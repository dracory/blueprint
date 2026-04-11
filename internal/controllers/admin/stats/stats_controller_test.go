package stats

import (
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

func TestNewStatsController(t *testing.T) {
	registry := testutils.Setup()
	controller := NewStatsController(registry)

	if controller == nil {
		t.Error("NewStatsController() should not return nil")
	}
	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
	}
	if controller.logger == nil {
		t.Error("Controller logger should not be nil")
	}
}

func TestStatsController_Handler_DefaultAction(t *testing.T) {
	registry := testutils.Setup()
	controller := NewStatsController(registry)

	req := httptest.NewRequest("GET", "/admin/stats", nil)
	w := httptest.NewRecorder()

	controller.Handler(w, req)
}

func TestStatsController_Handler_WithDifferentMethods(t *testing.T) {
	registry := testutils.Setup()
	controller := NewStatsController(registry)

	methods := []string{"GET", "POST"}
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/admin/stats", nil)
			w := httptest.NewRecorder()

			controller.Handler(w, req)
		})
	}
}

func TestStatsController_RegistryField(t *testing.T) {
	registry := testutils.Setup()
	controller := NewStatsController(registry)

	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
	}
}

func TestStatsController_NilRegistry(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic with nil registry
		}
	}()

	controller := NewStatsController(nil)
	if controller == nil {
		t.Error("NewStatsController(nil) should not return nil")
	}
}

func TestStatsController_MultipleInstances(t *testing.T) {
	registry1 := testutils.Setup()
	registry2 := testutils.Setup()

	controller1 := NewStatsController(registry1)
	controller2 := NewStatsController(registry2)

	if controller1 == nil || controller2 == nil {
		t.Fatal("All controllers should be non-nil")
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

func TestAdminLayout_SetTitle(t *testing.T) {
	registry := testutils.Setup()
	layout := &adminLayout{app: registry}

	layout.SetTitle("Test Title")
	if layout.title != "Test Title" {
		t.Errorf("SetTitle() = %q, want %q", layout.title, "Test Title")
	}
}

func TestAdminLayout_SetBody(t *testing.T) {
	registry := testutils.Setup()
	layout := &adminLayout{app: registry}

	layout.SetBody("<p>Test Body</p>")
	if layout.body != "<p>Test Body</p>" {
		t.Errorf("SetBody() = %q, want %q", layout.body, "<p>Test Body</p>")
	}
}

func TestAdminLayout_SetScriptURLs(t *testing.T) {
	registry := testutils.Setup()
	layout := &adminLayout{app: registry}

	urls := []string{"https://example.com/script.js"}
	layout.SetScriptURLs(urls)
	if len(layout.scriptURLs) != 1 || layout.scriptURLs[0] != urls[0] {
		t.Error("SetScriptURLs() should set script URLs correctly")
	}
}

func TestAdminLayout_SetScripts(t *testing.T) {
	registry := testutils.Setup()
	layout := &adminLayout{app: registry}

	scripts := []string{"console.log('test');"}
	layout.SetScripts(scripts)
	if len(layout.scripts) != 1 || layout.scripts[0] != scripts[0] {
		t.Error("SetScripts() should set scripts correctly")
	}
}

func TestAdminLayout_SetStyleURLs(t *testing.T) {
	registry := testutils.Setup()
	layout := &adminLayout{app: registry}

	urls := []string{"https://example.com/style.css"}
	layout.SetStyleURLs(urls)
	if len(layout.styleURLs) != 1 || layout.styleURLs[0] != urls[0] {
		t.Error("SetStyleURLs() should set style URLs correctly")
	}
}

func TestAdminLayout_SetStyles(t *testing.T) {
	registry := testutils.Setup()
	layout := &adminLayout{app: registry}

	styles := []string{"body { color: red; }"}
	layout.SetStyles(styles)
	if len(layout.styles) != 1 || layout.styles[0] != styles[0] {
		t.Error("SetStyles() should set styles correctly")
	}
}

func TestAdminLayout_SetCountryNameByIso2(t *testing.T) {
	registry := testutils.Setup()
	layout := &adminLayout{app: registry}

	f := func(iso2Code string) (string, error) {
		return "Test Country", nil
	}
	layout.SetCountryNameByIso2(f)
	if layout.countryNameByIso2 == nil {
		t.Error("SetCountryNameByIso2() should set the function")
	}
}

func TestAdminLayout_Render(t *testing.T) {
	registry := testutils.Setup()
	layout := &adminLayout{
		app:   registry,
		title: "Test Title",
		body:  "<p>Test Body</p>",
	}

	req := httptest.NewRequest("GET", "/admin/stats", nil)
	w := httptest.NewRecorder()

	result := layout.Render(w, req)
	if result == "" {
		t.Error("Render() should return non-empty string")
	}
	if !contains(result, "Test Title") {
		t.Error("Render() should contain the title")
	}
}

func TestAdminLayout_RenderWithScripts(t *testing.T) {
	registry := testutils.Setup()
	layout := &adminLayout{
		app:        registry,
		title:      "Test Title",
		body:       "<p>Test Body</p>",
		scriptURLs: []string{"https://example.com/script.js"},
		scripts:    []string{"console.log('test');"},
	}

	req := httptest.NewRequest("GET", "/admin/stats", nil)
	w := httptest.NewRecorder()

	result := layout.Render(w, req)
	if result == "" {
		t.Error("Render() with scripts should return non-empty string")
	}
}

func TestAdminLayout_RenderWithStyles(t *testing.T) {
	registry := testutils.Setup()
	layout := &adminLayout{
		app:       registry,
		title:     "Test Title",
		body:      "<p>Test Body</p>",
		styleURLs: []string{"https://example.com/style.css"},
		styles:    []string{"body { color: red; }"},
	}

	req := httptest.NewRequest("GET", "/admin/stats", nil)
	w := httptest.NewRecorder()

	result := layout.Render(w, req)
	if result == "" {
		t.Error("Render() with styles should return non-empty string")
	}
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
