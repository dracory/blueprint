package category_manager

import (
	"net/http/httptest"
	"testing"

	"project/internal/testutils"
)

func TestNewCategoryManagerController(t *testing.T) {
	// Test with nil registry
	controller := NewCategoryManagerController(nil)
	if controller == nil {
		t.Error("NewCategoryManagerController() should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	controller = NewCategoryManagerController(registry)
	if controller == nil {
		t.Error("NewCategoryManagerController() should not return nil")
	}
	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
	}
}

func TestCategoryManagerControllerHandler(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithUserStore(true),
		testutils.WithSessionStore(true),
		testutils.WithBlogStore(true),
	)
	controller := NewCategoryManagerController(registry)
	if controller == nil {
		t.Fatal("NewCategoryManagerController() returned nil")
	}

	user, _, err := testutils.SeedUserAndSession(
		registry.GetUserStore(),
		registry.GetSessionStore(),
		"test-admin",
		httptest.NewRequest("GET", "/", nil),
		3600,
	)
	if err != nil {
		t.Fatalf("Failed to seed user and session: %v", err)
	}

	req := httptest.NewRequest("GET", "/admin/blog/categories", nil)
	w := httptest.NewRecorder()

	req, err = testutils.LoginAs(registry, req, user)
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}

	controller.Handler(w, req)
}

func TestCategoryManagerController_HandlerWithActions(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithUserStore(true),
		testutils.WithSessionStore(true),
		testutils.WithBlogStore(true),
	)
	controller := NewCategoryManagerController(registry)

	user, _, err := testutils.SeedUserAndSession(
		registry.GetUserStore(),
		registry.GetSessionStore(),
		"test-admin",
		httptest.NewRequest("GET", "/", nil),
		3600,
	)
	if err != nil {
		t.Fatalf("Failed to seed user and session: %v", err)
	}

	tests := []struct {
		name   string
		action string
		method string
	}{
		{"load-categories", "load-categories", "GET"},
		{"create-category", "create-category", "POST"},
		{"update-category", "update-category", "POST"},
		{"delete-category", "delete-category", "POST"},
		{"reorder-categories", "reorder-categories", "POST"},
		{"default", "", "GET"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url := "/admin/blog/categories"
			if tt.action != "" {
				url += "?action=" + tt.action
			}

			req := httptest.NewRequest(tt.method, url, nil)
			w := httptest.NewRecorder()

			req, err = testutils.LoginAs(registry, req, user)
			if err != nil {
				t.Fatalf("Failed to login: %v", err)
			}

			controller.Handler(w, req)
		})
	}
}

func TestCategoryManagerController_NilRegistry(t *testing.T) {
	controller := NewCategoryManagerController(nil)
	if controller == nil {
		t.Error("NewCategoryManagerController(nil) should not return nil")
	}
	if controller.registry != nil {
		t.Error("Controller registry should be nil when passed nil")
	}
}
