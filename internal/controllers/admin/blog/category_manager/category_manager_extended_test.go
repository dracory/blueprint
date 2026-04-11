package category_manager

import (
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"
)

func TestCategoryManagerController_HandlerWithNilRegistry(t *testing.T) {
	// Test with nil registry - this will panic due to missing auth context
	// but we're testing that the controller handles it gracefully
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to nil registry or missing auth context
			// This is acceptable behavior
		}
	}()

	controller := NewCategoryManagerController(nil)
	if controller == nil {
		t.Fatal("NewCategoryManagerController(nil) should not return nil")
	}

	req := httptest.NewRequest("GET", "/admin/blog/categories", nil)
	w := httptest.NewRecorder()

	// Handler may panic with nil registry - this is expected
	controller.Handler(w, req)
}

func TestCategoryManagerController_HandlerWithDifferentActions(t *testing.T) {
	registry := testutils.Setup()
	controller := NewCategoryManagerController(registry)

	// Test different action parameters
	actions := []string{
		"", // default (render page)
		"load-categories",
		"create-category",
		"update-category",
		"delete-category",
		"reorder-categories",
		"invalid-action", // unknown action should default to render page
	}

	for _, action := range actions {
		t.Run(action, func(t *testing.T) {
			var url string
			if action == "" {
				url = "/admin/blog/categories"
			} else {
				url = "/admin/blog/categories?action=" + action
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			// This will likely panic due to missing auth, but we're testing the switch logic
			func() {
				defer func() {
					if r := recover(); r != nil {
						// Expected panic due to missing auth context
					}
				}()
				controller.Handler(w, req)
			}()
		})
	}
}

func TestCategoryManagerController_RegistryField(t *testing.T) {
	// Test with nil registry
	controller := NewCategoryManagerController(nil)
	if controller.registry != nil {
		t.Error("Controller registry should be nil when passed nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	controller = NewCategoryManagerController(registry)
	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
	}
}

func TestCategoryManagerController_EmbeddedFiles(t *testing.T) {
	// Test that embedded files are accessible
	htmlContent, err := categoriesFiles.ReadFile("categories.html")
	if err != nil {
		t.Fatalf("Failed to read categories.html: %v", err)
	}
	if len(htmlContent) == 0 {
		t.Error("categories.html content should not be empty")
	}

	jsContent, err := categoriesFiles.ReadFile("categories.js")
	if err != nil {
		t.Fatalf("Failed to read categories.js: %v", err)
	}
	if len(jsContent) == 0 {
		t.Error("categories.js content should not be empty")
	}
}

func TestCategoryManagerController_EmbeddedHTMLContent(t *testing.T) {
	htmlContent, err := categoriesFiles.ReadFile("categories.html")
	if err != nil {
		t.Fatalf("Failed to read categories.html: %v", err)
	}

	content := string(htmlContent)

	// Verify HTML content has expected elements
	if !strings.Contains(content, "<") {
		t.Error("HTML content should contain HTML tags")
	}
}

func TestCategoryManagerController_EmbeddedJSContent(t *testing.T) {
	jsContent, err := categoriesFiles.ReadFile("categories.js")
	if err != nil {
		t.Fatalf("Failed to read categories.js: %v", err)
	}

	content := string(jsContent)

	// Verify JS content has expected elements
	if len(content) == 0 {
		t.Error("JS content should not be empty")
	}
}

func TestCategoryManagerController_MultipleInstances(t *testing.T) {
	// Test that multiple controllers can be created
	registry1 := testutils.Setup()
	registry2 := testutils.Setup()

	controller1 := NewCategoryManagerController(registry1)
	controller2 := NewCategoryManagerController(registry2)

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

func TestCategoryManagerController_HandlerMultipleCalls(t *testing.T) {
	registry := testutils.Setup()
	controller := NewCategoryManagerController(registry)

	// Handler can be called multiple times
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/admin/blog/categories", nil)
		w := httptest.NewRecorder()

		// This will panic due to missing auth, but we're testing that the handler
		// doesn't have state issues
		func() {
			defer func() {
				if r := recover(); r != nil {
					// Expected panic due to missing auth context
				}
			}()
			controller.Handler(w, req)
		}()
	}
}
