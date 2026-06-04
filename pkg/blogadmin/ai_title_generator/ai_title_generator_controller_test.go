package aititlegenerator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/config"
	"project/internal/app"
	"project/internal/testutils"

	"github.com/dracory/test"
)

// TestNewAiTitleGeneratorController tests the constructor
func TestNewAiTitleGeneratorController(t *testing.T) {
	t.Parallel()

	// Test with interface - we can only test with nil since AppInterface is not directly instantiable
	var mockRegistry app.AppInterface
	controller := NewAiTitleGeneratorController(mockRegistry)
	if controller == nil {
		t.Fatal("NewAiTitleGeneratorController should not return nil")
	}
	if controller.app != mockRegistry {
		t.Error("Controller should store the app")
	}

	// Note: Nil app acceptance is tested but methods will panic
	// if app is nil when accessing stores. This is acceptable
	// as long as production code never passes nil.
}

// TestAiTitleGeneratorController_Struct tests the controller struct fields
func TestAiTitleGeneratorController_Struct(t *testing.T) {
	t.Parallel()

	controller := &AiTitleGeneratorController{}

	// Test that app field exists and can be set
	var reg app.AppInterface
	controller.app = reg

	if controller.app != reg {
		t.Error("Should be able to set app field")
	}
}

// TestPageData_Struct tests the pageData struct
func TestPageData_Struct(t *testing.T) {
	t.Parallel()

	data := pageData{
		Action:          "test_action",
		HasSystemPrompt: true,
	}

	if data.Action != "test_action" {
		t.Error("Action field should be settable")
	}
	if !data.HasSystemPrompt {
		t.Error("HasSystemPrompt field should be settable")
	}
}

// TestConstants_ACTION_ADD_TITLE tests the ACTION_ADD_TITLE constant
func TestConstants_ACTION_ADD_TITLE(t *testing.T) {
	t.Parallel()

	if ACTION_ADD_TITLE != "add_title" {
		t.Errorf("ACTION_ADD_TITLE = %q, want %q", ACTION_ADD_TITLE, "add_title")
	}
}

// TestConstants_ACTION_GENERATE_TITLES tests the ACTION_GENERATE_TITLES constant
func TestConstants_ACTION_GENERATE_TITLES(t *testing.T) {
	t.Parallel()

	if ACTION_GENERATE_TITLES != "generate_titles" {
		t.Errorf("ACTION_GENERATE_TITLES = %q, want %q", ACTION_GENERATE_TITLES, "generate_titles")
	}
}

// TestConstants_ACTION_APPROVE_TITLE tests the ACTION_APPROVE_TITLE constant
func TestConstants_ACTION_APPROVE_TITLE(t *testing.T) {
	t.Parallel()

	if ACTION_APPROVE_TITLE != "approve_title" {
		t.Errorf("ACTION_APPROVE_TITLE = %q, want %q", ACTION_APPROVE_TITLE, "approve_title")
	}
}

// TestConstants_ACTION_REJECT_TITLE tests the ACTION_REJECT_TITLE constant
func TestConstants_ACTION_REJECT_TITLE(t *testing.T) {
	t.Parallel()

	if ACTION_REJECT_TITLE != "reject_title" {
		t.Errorf("ACTION_REJECT_TITLE = %q, want %q", ACTION_REJECT_TITLE, "reject_title")
	}
}

// TestConstants_ACTION_GENERATE_POST tests the ACTION_GENERATE_POST constant
func TestConstants_ACTION_GENERATE_POST(t *testing.T) {
	t.Parallel()

	if ACTION_GENERATE_POST != "generate_post" {
		t.Errorf("ACTION_GENERATE_POST = %q, want %q", ACTION_GENERATE_POST, "generate_post")
	}
}

// TestConstants_ACTION_DELETE_TITLE tests the ACTION_DELETE_TITLE constant
func TestConstants_ACTION_DELETE_TITLE(t *testing.T) {
	t.Parallel()

	if ACTION_DELETE_TITLE != "delete_title" {
		t.Errorf("ACTION_DELETE_TITLE = %q, want %q", ACTION_DELETE_TITLE, "delete_title")
	}
}

// TestSettingKeyConstant tests the setting key constant
func TestSettingKeyConstant(t *testing.T) {
	t.Parallel()

	if SETTING_KEY_BLOG_TOPIC != "title_generator.blog_topic" {
		t.Errorf("SETTING_KEY_BLOG_TOPIC = %q, want %q", SETTING_KEY_BLOG_TOPIC, "title_generator.blog_topic")
	}
}

// TestAiTitleGeneratorController_MultipleInstances tests creating multiple controllers
func TestAiTitleGeneratorController_MultipleInstances(t *testing.T) {
	t.Parallel()

	// Test with nil registries - each should be independent
	controller1 := NewAiTitleGeneratorController(nil)
	controller2 := NewAiTitleGeneratorController(nil)

	if controller1 == controller2 {
		t.Error("Multiple instances should be independent")
	}

	// Test that each has its own app reference
	var mockRegistry1 app.AppInterface
	var mockRegistry2 app.AppInterface

	controller1.app = mockRegistry1
	controller2.app = mockRegistry2

	if controller1.app != mockRegistry1 {
		t.Error("First controller should have correct app")
	}

	if controller2.app != mockRegistry2 {
		t.Error("Second controller should have correct app")
	}
}

// TestAiTitleGeneratorController_Handler_MethodExists verifies Handler method exists
func TestAiTitleGeneratorController_Handler_MethodExists(t *testing.T) {
	t.Parallel()

	// This test verifies the method signature exists
	// The actual handler requires HTTP request/response which would need integration testing
	controller := NewAiTitleGeneratorController(nil)
	if controller == nil {
		t.Fatal("Controller should not be nil")
	}

	// Method existence is verified by compilation
	// We can't easily test the actual handler without a full HTTP setup
}

// TestAiTitleGeneratorController_RenderPage tests rendering the page
func TestAiTitleGeneratorController_RenderPage(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithCustomStore(true),
	)

	user, _ := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	controller := NewAiTitleGeneratorController(app)

	// Context with auth user
	ctx := context.WithValue(context.Background(), config.AuthenticatedUserContextKey{}, user)

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-title-generator", nil).WithContext(ctx)
	resp := controller.Handler(httptest.NewRecorder(), req)
	if !strings.Contains(resp, "AI Title Generator") {
		t.Error("expected AI Title Generator in response")
	}
}

// TestAiTitleGeneratorController_OnAddTitleModal tests the add title modal
func TestAiTitleGeneratorController_OnAddTitleModal(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithCustomStore(true),
	)

	user, _ := testutils.SeedUser(app.GetUserStore(), test.USER_01)
	controller := NewAiTitleGeneratorController(app)

	// Context with auth user
	ctx := context.WithValue(context.Background(), config.AuthenticatedUserContextKey{}, user)

	req := httptest.NewRequest(http.MethodGet, "/admin/blog/ai-title-generator?action="+ACTION_ADD_TITLE, nil).WithContext(ctx)
	resp := controller.Handler(httptest.NewRecorder(), req)
	if !strings.Contains(resp, "Add Custom Title") {
		t.Error("expected Add Custom Title in response")
	}
}
