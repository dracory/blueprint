package admin

import (
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/userstore"
)

// TestNewUserManagerController verifies controller can be created
func TestNewUserManagerController(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewUserManagerController(app)

	if controller == nil {
		t.Error("NewUserManagerController() returned nil")
	}
}

// TestUserManagerControllerRegistry verifies controller has registry
func TestUserManagerControllerRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewUserManagerController(app)

	if controller.registry == nil {
		t.Error("Controller registry is nil")
	}
}

// TestUserManagerControllerHandlerExists verifies Handler method exists
func TestUserManagerControllerHandlerExists(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })

	controller := NewUserManagerController(app)

	// Verify Handler method exists (should compile without error)
	_ = controller.Handler
}

func TestUserManagerController_NilRegistry(t *testing.T) {
	t.Parallel()
	controller := NewUserManagerController(nil)
	if controller == nil {
		t.Error("NewUserManagerController(nil) should not return nil")
	}
	if controller.registry != nil {
		t.Error("Controller registry should be nil when passed nil")
	}
	// No cleanup needed as no database was created
}

func TestUserManagerController_MultipleInstances(t *testing.T) {
	t.Parallel()
	app1 := testutils.Setup()
	app2 := testutils.Setup()
	if app1 == nil || app2 == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app1.GetDatabase().Close() })
	t.Cleanup(func() { _ = app2.GetDatabase().Close() })

	controller1 := NewUserManagerController(app1)
	controller2 := NewUserManagerController(app2)

	if controller1 == nil || controller2 == nil {
		t.Fatal("All controllers should be non-nil")
	}

	if controller1 == controller2 {
		t.Error("Controllers should be separate instances")
	}

	if controller1.registry != app1 {
		t.Error("Controller1 should have app1")
	}

	if controller2.registry != app2 {
		t.Error("Controller2 should have app2")
	}
}

func TestUserManagerController_Handler_Actions(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })
	controller := NewUserManagerController(app)

	// Verify handler exists - methods cannot be nil in Go
	// This test ensures the Handler method is callable
	_ = controller.Handler
}

func TestUserManagerController_PrepareData_Basic(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserManagerController(app)
	r, _ := testutils.NewRequest("GET", "/admin/users", testutils.NewRequestOptions{})

	data, err := controller.prepareData(r)
	if err != "" {
		t.Errorf("Expected no error, got: %s", err)
	}
	if data.request == nil {
		t.Error("Expected request to be set in data")
	}
	if data.pageInt != 0 {
		t.Errorf("Expected pageInt to be 0, got %d", data.pageInt)
	}
	if data.perPage != 10 {
		t.Errorf("Expected perPage to be 10, got %d", data.perPage)
	}
}

func TestUserManagerController_PrepareData_WithQueryParams(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserManagerController(app)
	queryParams := url.Values{}
	queryParams.Set("page", "2")
	queryParams.Set("per_page", "25")
	queryParams.Set("sort_order", "asc")
	queryParams.Set("by", "created_at")
	queryParams.Set("status", "active")
	// Don't use email filter as it requires blind index store
	r, _ := testutils.NewRequest("GET", "/admin/users", testutils.NewRequestOptions{QueryParams: queryParams})

	data, err := controller.prepareData(r)
	if err != "" {
		t.Errorf("Expected no error, got: %s", err)
	}
	if data.pageInt != 2 {
		t.Errorf("Expected pageInt to be 2, got %d", data.pageInt)
	}
	if data.perPage != 25 {
		t.Errorf("Expected perPage to be 25, got %d", data.perPage)
	}
	if data.sortOrder != "asc" {
		t.Errorf("Expected sortOrder to be 'asc', got '%s'", data.sortOrder)
	}
	if data.formStatus != "active" {
		t.Errorf("Expected formStatus to be 'active', got '%s'", data.formStatus)
	}
}

func TestUserManagerController_PrepareData_NilRegistry(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserManagerController(app)
	r, _ := testutils.NewRequest("GET", "/admin/users", testutils.NewRequestOptions{})

	data, err := controller.prepareData(r)
	// With nil UserStore, prepareData should return an error
	if err == "" {
		t.Error("Expected error when UserStore is not configured")
	}
	if data.request == nil {
		t.Error("Expected request to be set in data")
	}
}

func TestUserManagerController_OnModalUserFilterShow(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })
	controller := NewUserManagerController(app)

	data := userManagerControllerData{
		formStatus:      "active",
		formFirstName:   "John",
		formLastName:    "Doe",
		formEmail:       "john@example.com",
		formCreatedFrom: "2024-01-01",
		formCreatedTo:   "2024-12-31",
		formUserID:      "user123",
	}

	tag := controller.onModalUserFilterShow(data)
	if tag == nil {
		t.Error("Expected non-nil tag from onModalUserFilterShow")
	}
	html := tag.ToHTML()
	if html == "" {
		t.Error("Expected non-empty HTML from onModalUserFilterShow")
	}
}

func TestUserManagerController_OnModalUserFilterShow_EmptyData(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })
	controller := NewUserManagerController(app)

	data := userManagerControllerData{}

	tag := controller.onModalUserFilterShow(data)
	if tag == nil {
		t.Error("Expected non-nil tag from onModalUserFilterShow with empty data")
	}
	html := tag.ToHTML()
	if html == "" {
		t.Error("Expected non-empty HTML from onModalUserFilterShow with empty data")
	}
}

func TestUserManagerController_Page(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	if app == nil {
		t.Fatal("testutils.Setup() returned nil")
	}
	t.Cleanup(func() { _ = app.GetDatabase().Close() })
	controller := NewUserManagerController(app)

	data := userManagerControllerData{
		userList:  []userstore.UserInterface{},
		userCount: 0,
	}

	tag := controller.page(data)
	if tag == nil {
		t.Error("Expected non-nil tag from page")
	}
	html := tag.ToHTML()
	if html == "" {
		t.Error("Expected non-empty HTML from page")
	}
}

func TestUserManagerController_FetchUserList_NilUserStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup()
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserManagerController(app)
	data := userManagerControllerData{}

	userList, userCount, err := controller.fetchUserList(data)
	if err == nil {
		t.Error("Expected error when UserStore is nil")
	}
	// fetchUserList returns empty list, not nil
	if userList == nil {
		t.Error("Expected empty list when UserStore is nil")
	}
	if userCount != 0 {
		t.Error("Expected userCount to be 0 when UserStore is nil")
	}
}

func TestUserManagerController_FetchUserList_WithUserStore(t *testing.T) {
	t.Parallel()
	app := testutils.Setup(testutils.WithUserStore(true))
	defer func() { _ = app.GetDatabase().Close() }()

	controller := NewUserManagerController(app)
	r, _ := testutils.NewRequest("GET", "/admin/users", testutils.NewRequestOptions{})
	data := userManagerControllerData{
		request:   r,
		pageInt:   0,
		perPage:   10,
		sortOrder: "desc",
		sortBy:    "created_at",
	}

	userList, userCount, err := controller.fetchUserList(data)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if userList == nil {
		t.Error("Expected non-nil userList")
	}
	if userCount < 0 {
		t.Errorf("Expected userCount to be >= 0, got %d", userCount)
	}
}
