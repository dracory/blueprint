package admin

import (
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

func TestNewTaskController(t *testing.T) {
	app := testutils.Setup()
	controller := NewTaskController(app)

	if controller == nil {
		t.Error("NewTaskController() should not return nil")
	}
	if controller.app != app {
		t.Error("Controller app should match the provided app")
	}
	if controller.logger == nil {
		t.Error("Controller logger should not be nil")
	}
}

func TestTaskController_Handler_DefaultAction(t *testing.T) {
	app := testutils.Setup()
	controller := NewTaskController(app)

	req := httptest.NewRequest("GET", "/admin/tasks", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)
	if result == "" {
		t.Error("Handler() should return non-empty string")
	}
}

func TestTaskController_Handler_WithTaskStore(t *testing.T) {
	app := testutils.Setup(
		testutils.WithTaskStore(true),
	)
	controller := NewTaskController(app)

	req := httptest.NewRequest("GET", "/admin/tasks", nil)
	w := httptest.NewRecorder()

	result := controller.Handler(w, req)
	if result == "" {
		t.Error("Handler() with task store should return non-empty string")
	}
}

func TestTaskController_RegistryField(t *testing.T) {
	app := testutils.Setup()
	controller := NewTaskController(app)

	if controller.app != app {
		t.Error("Controller app should match the provided app")
	}
}

func TestTaskController_NilRegistry(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			// Expected to panic with nil app
		}
	}()

	controller := NewTaskController(nil)
	if controller == nil {
		t.Error("NewTaskController(nil) should not return nil")
	}
}

func TestTaskController_MultipleInstances(t *testing.T) {
	registry1 := testutils.Setup()
	registry2 := testutils.Setup()

	controller1 := NewTaskController(registry1)
	controller2 := NewTaskController(registry2)

	if controller1 == nil || controller2 == nil {
		t.Fatal("All controllers should be non-nil")
	}

	if controller1 == controller2 {
		t.Error("Controllers should be separate instances")
	}

	if controller1.app != registry1 {
		t.Error("Controller1 should have registry1")
	}

	if controller2.app != registry2 {
		t.Error("Controller2 should have registry2")
	}
}

func TestAdminLayout_SetTitle(t *testing.T) {
	app := testutils.Setup()
	layout := &adminLayout{app: app}

	layout.SetTitle("Test Title")
	if layout.title != "Test Title" {
		t.Errorf("SetTitle() = %q, want %q", layout.title, "Test Title")
	}
}

func TestAdminLayout_SetBody(t *testing.T) {
	app := testutils.Setup()
	layout := &adminLayout{app: app}

	layout.SetBody("<p>Test Body</p>")
	if layout.body != "<p>Test Body</p>" {
		t.Errorf("SetBody() = %q, want %q", layout.body, "<p>Test Body</p>")
	}
}

func TestAdminLayout_SetScriptURLs(t *testing.T) {
	app := testutils.Setup()
	layout := &adminLayout{app: app}

	urls := []string{"https://example.com/script.js"}
	layout.SetScriptURLs(urls)
	if len(layout.scriptURLs) != 1 || layout.scriptURLs[0] != urls[0] {
		t.Error("SetScriptURLs() should set script URLs correctly")
	}
}

func TestAdminLayout_SetScripts(t *testing.T) {
	app := testutils.Setup()
	layout := &adminLayout{app: app}

	scripts := []string{"console.log('test');"}
	layout.SetScripts(scripts)
	if len(layout.scripts) != 1 || layout.scripts[0] != scripts[0] {
		t.Error("SetScripts() should set scripts correctly")
	}
}

func TestAdminLayout_SetStyleURLs(t *testing.T) {
	app := testutils.Setup()
	layout := &adminLayout{app: app}

	urls := []string{"https://example.com/style.css"}
	layout.SetStyleURLs(urls)
	if len(layout.styleURLs) != 1 || layout.styleURLs[0] != urls[0] {
		t.Error("SetStyleURLs() should set style URLs correctly")
	}
}

func TestAdminLayout_SetStyles(t *testing.T) {
	app := testutils.Setup()
	layout := &adminLayout{app: app}

	styles := []string{"body { color: red; }"}
	layout.SetStyles(styles)
	if len(layout.styles) != 1 || layout.styles[0] != styles[0] {
		t.Error("SetStyles() should set styles correctly")
	}
}

func TestAdminLayout_Render(t *testing.T) {
	app := testutils.Setup()
	layout := &adminLayout{
		app:   app,
		title: "Test Title",
		body:  "<p>Test Body</p>",
	}

	req := httptest.NewRequest("GET", "/admin/tasks", nil)
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
	app := testutils.Setup()
	layout := &adminLayout{
		app:        app,
		title:      "Test Title",
		body:       "<p>Test Body</p>",
		scriptURLs: []string{"https://example.com/script.js"},
		scripts:    []string{"console.log('test');"},
	}

	req := httptest.NewRequest("GET", "/admin/tasks", nil)
	w := httptest.NewRecorder()

	result := layout.Render(w, req)
	if result == "" {
		t.Error("Render() with scripts should return non-empty string")
	}
}

func TestAdminLayout_RenderWithStyles(t *testing.T) {
	app := testutils.Setup()
	layout := &adminLayout{
		app:       app,
		title:     "Test Title",
		body:      "<p>Test Body</p>",
		styleURLs: []string{"https://example.com/style.css"},
		styles:    []string{"body { color: red; }"},
	}

	req := httptest.NewRequest("GET", "/admin/tasks", nil)
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
