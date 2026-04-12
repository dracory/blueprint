package user_delete

import (
	"strings"
	"testing"

	"project/internal/testutils"
)

func TestNewUserDeleteController(t *testing.T) {
	t.Parallel()
	// Test with nil registry
	controller := NewUserDeleteController(nil)
	if controller == nil {
		t.Error("NewUserDeleteController() should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	if registry == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = registry.GetDatabase().Close() })
	controller = NewUserDeleteController(registry)
	if controller == nil {
		t.Error("NewUserDeleteController() should not return nil")
	}
}

func TestUserDeleteController_PrepareDataAndValidate_GET(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserDeleteController(app)
	user, _ := testutils.SeedUser(app.GetUserStore(), "test-user")
	r, _ := testutils.NewRequest("GET", "/admin/users/delete", testutils.NewRequestOptions{
		QueryParams: map[string][]string{
			"user_id": {"test-user"},
		},
	})
	r, _ = testutils.LoginAs(app, r, user)

	data, err := controller.prepareDataAndValidate(r)
	if err != "" {
		t.Errorf("Expected no error for GET request, got: %s", err)
	}
	if data.userID != "test-user" {
		t.Errorf("Expected userID to be 'test-user', got: %s", data.userID)
	}
}

func TestUserDeleteController_PrepareDataAndValidate_POST_Success(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserDeleteController(app)
	user, _ := testutils.SeedUser(app.GetUserStore(), "test-user")
	r, _ := testutils.NewRequest("POST", "/admin/users/delete", testutils.NewRequestOptions{
		QueryParams: map[string][]string{
			"user_id": {"test-user"},
		},
	})
	r, _ = testutils.LoginAs(app, r, user)

	data, err := controller.prepareDataAndValidate(r)
	if err != "" {
		t.Errorf("Expected no error for valid POST request, got: %s", err)
	}
	if data.successMessage == "" {
		t.Error("Expected success message for valid POST request")
	}
}

func TestUserDeleteController_PrepareDataAndValidate_MissingUserID(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserDeleteController(app)
	user, _ := testutils.SeedUser(app.GetUserStore(), "test-user")
	r, _ := testutils.NewRequest("GET", "/admin/users/delete", testutils.NewRequestOptions{})
	r, _ = testutils.LoginAs(app, r, user)

	_, err := controller.prepareDataAndValidate(r)
	if err == "" {
		t.Error("Expected error when user_id is missing")
	}
	if !strings.Contains(err, "user id is required") {
		t.Errorf("Expected error containing 'user id is required', got: %s", err)
	}
}

func TestUserDeleteController_PrepareDataAndValidate_NilUserStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserDeleteController(app)
	r, _ := testutils.NewRequest("GET", "/admin/users/delete", testutils.NewRequestOptions{})

	_, err := controller.prepareDataAndValidate(r)
	if err == "" {
		t.Error("Expected error when UserStore is nil")
	}
	if !strings.Contains(err, "User store is not configured") {
		t.Errorf("Expected error containing 'User store is not configured', got: %s", err)
	}
}

func TestUserDeleteController_PrepareDataAndValidate_Unauthenticated(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserDeleteController(app)
	r, _ := testutils.NewRequest("GET", "/admin/users/delete", testutils.NewRequestOptions{})

	_, err := controller.prepareDataAndValidate(r)
	if err == "" {
		t.Error("Expected error when user is not authenticated")
	}
	if !strings.Contains(err, "not logged in") {
		t.Errorf("Expected error containing 'not logged in', got: %s", err)
	}
}

func TestUserDeleteController_PrepareDataAndValidate_UserNotFound(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserDeleteController(app)
	user, _ := testutils.SeedUser(app.GetUserStore(), "test-user")
	r, _ := testutils.NewRequest("GET", "/admin/users/delete", testutils.NewRequestOptions{
		QueryParams: map[string][]string{
			"user_id": {"nonexistent"},
		},
	})
	r, _ = testutils.LoginAs(app, r, user)

	_, err := controller.prepareDataAndValidate(r)
	if err == "" {
		t.Error("Expected error when user is not found")
	}
	if !strings.Contains(err, "User not found") {
		t.Errorf("Expected error containing 'User not found', got: %s", err)
	}
}

func TestUserDeleteController_Modal(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserDeleteController(app)
	user, _ := testutils.SeedUser(app.GetUserStore(), "test-user")

	data := userDeleteControllerData{
		userID: "test-user",
		user:   user,
	}

	tag := controller.modal(data)
	if tag == nil {
		t.Error("Expected non-nil tag from modal")
	}
	html := tag.ToHTML()
	if html == "" {
		t.Error("Expected non-empty HTML from modal")
	}
}

func TestUserDeleteController_Modal_EmptyData(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserDeleteController(app)

	data := userDeleteControllerData{}

	tag := controller.modal(data)
	if tag == nil {
		t.Error("Expected non-nil tag from modal with empty data")
	}
	html := tag.ToHTML()
	if html == "" {
		t.Error("Expected non-empty HTML from modal with empty data")
	}
}
