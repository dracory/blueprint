package user_create

import (
	"strings"
	"testing"

	"project/internal/testutils"
)

func TestNewUserCreateController(t *testing.T) {
	t.Parallel()
	// Test with nil registry
	controller := NewUserCreateController(nil)
	if controller == nil {
		t.Error("NewUserCreateController() should not return nil")
	}

	// Test with valid registry
	registry := testutils.Setup()
	if registry == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = registry.GetDatabase().Close() })
	controller = NewUserCreateController(registry)
	if controller == nil {
		t.Error("NewUserCreateController() should not return nil")
	}
}

func TestUserCreateController_RegistryField(t *testing.T) {
	t.Parallel()
	registry := testutils.Setup()
	if registry == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = registry.GetDatabase().Close() })
	controller := NewUserCreateController(registry)

	if controller.registry != registry {
		t.Error("Controller registry should match the provided registry")
	}
}

func TestUserCreateController_MultipleInstances(t *testing.T) {
	t.Parallel()
	registry1 := testutils.Setup()
	registry2 := testutils.Setup()
	if registry1 == nil || registry2 == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = registry1.GetDatabase().Close() })
	t.Cleanup(func() { _ = registry2.GetDatabase().Close() })

	controller1 := NewUserCreateController(registry1)
	controller2 := NewUserCreateController(registry2)

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

func TestUserCreateController_Handler_Actions(t *testing.T) {
	t.Parallel()
	registry := testutils.Setup()
	if registry == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = registry.GetDatabase().Close() })
	controller := NewUserCreateController(registry)

	// Verify handler exists - methods cannot be nil in Go
	_ = controller.Handler
}

func TestUserCreateController_PrepareDataAndValidate_GET(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserCreateController(app)
	user, _ := testutils.SeedUser(app.GetUserStore(), "test-user")
	r, _ := testutils.NewRequest("GET", "/admin/users/create", testutils.NewRequestOptions{})
	r, _ = testutils.LoginAs(app, r, user)

	data, err := controller.prepareDataAndValidate(r)
	if err != "" {
		t.Errorf("Expected no error for GET request, got: %s", err)
	}
	if data.firstName != "" {
		t.Errorf("Expected empty firstName for GET request, got: %s", data.firstName)
	}
}

func TestUserCreateController_PrepareDataAndValidate_POST_Success(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserCreateController(app)
	user, _ := testutils.SeedUser(app.GetUserStore(), "test-user")
	r, _ := testutils.NewRequest("POST", "/admin/users/create", testutils.NewRequestOptions{
		FormValues: map[string][]string{
			"user_first_name": {"John"},
			"user_last_name":  {"Doe"},
			"user_email":      {"john@example.com"},
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
	if data.firstName != "John" {
		t.Errorf("Expected firstName to be 'John', got: %s", data.firstName)
	}
}

func TestUserCreateController_PrepareDataAndValidate_POST_MissingFirstName(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserCreateController(app)
	user, _ := testutils.SeedUser(app.GetUserStore(), "test-user")
	r, _ := testutils.NewRequest("POST", "/admin/users/create", testutils.NewRequestOptions{
		FormValues: map[string][]string{
			"user_last_name": {"Doe"},
			"user_email":     {"john@example.com"},
		},
	})
	r, _ = testutils.LoginAs(app, r, user)

	_, err := controller.prepareDataAndValidate(r)
	if err == "" {
		t.Error("Expected error when first name is missing")
	}
	if !strings.Contains(err, "first name is required") {
		t.Errorf("Expected error containing 'first name is required', got: %s", err)
	}
}

func TestUserCreateController_PrepareDataAndValidate_POST_MissingLastName(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserCreateController(app)
	user, _ := testutils.SeedUser(app.GetUserStore(), "test-user")
	r, _ := testutils.NewRequest("POST", "/admin/users/create", testutils.NewRequestOptions{
		FormValues: map[string][]string{
			"user_first_name": {"John"},
			"user_email":      {"john@example.com"},
		},
	})
	r, _ = testutils.LoginAs(app, r, user)

	_, err := controller.prepareDataAndValidate(r)
	if err == "" {
		t.Error("Expected error when last name is missing")
	}
	if !strings.Contains(err, "last name is required") {
		t.Errorf("Expected error containing 'last name is required', got: %s", err)
	}
}

func TestUserCreateController_PrepareDataAndValidate_POST_MissingEmail(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserCreateController(app)
	user, _ := testutils.SeedUser(app.GetUserStore(), "test-user")
	r, _ := testutils.NewRequest("POST", "/admin/users/create", testutils.NewRequestOptions{
		FormValues: map[string][]string{
			"user_first_name": {"John"},
			"user_last_name":  {"Doe"},
		},
	})
	r, _ = testutils.LoginAs(app, r, user)

	_, err := controller.prepareDataAndValidate(r)
	if err == "" {
		t.Error("Expected error when email is missing")
	}
	if !strings.Contains(err, "email is required") {
		t.Errorf("Expected error containing 'email is required', got: %s", err)
	}
}

func TestUserCreateController_PrepareDataAndValidate_NilUserStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserCreateController(app)
	r, _ := testutils.NewRequest("GET", "/admin/users/create", testutils.NewRequestOptions{})

	_, err := controller.prepareDataAndValidate(r)
	if err == "" {
		t.Error("Expected error when UserStore is nil")
	}
	if !strings.Contains(err, "User store is not configured") {
		t.Errorf("Expected error containing 'User store is not configured', got: %s", err)
	}
}

func TestUserCreateController_PrepareDataAndValidate_Unauthenticated(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserCreateController(app)
	r, _ := testutils.NewRequest("GET", "/admin/users/create", testutils.NewRequestOptions{})

	_, err := controller.prepareDataAndValidate(r)
	if err == "" {
		t.Error("Expected error when user is not authenticated")
	}
	if !strings.Contains(err, "not logged in") {
		t.Errorf("Expected error containing 'not logged in', got: %s", err)
	}
}

func TestUserCreateController_Modal(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserCreateController(app)

	data := userCreateControllerData{
		firstName: "John",
		lastName:  "Doe",
		email:     "john@example.com",
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

func TestUserCreateController_Modal_EmptyData(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserCreateController(app)

	data := userCreateControllerData{}

	tag := controller.modal(data)
	if tag == nil {
		t.Error("Expected non-nil tag from modal with empty data")
	}
	html := tag.ToHTML()
	if html == "" {
		t.Error("Expected non-empty HTML from modal with empty data")
	}
}
