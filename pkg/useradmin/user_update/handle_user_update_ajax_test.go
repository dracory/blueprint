package user_update

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"project/internal/app"
	"project/internal/tasks/blind_index_rebuild"
	"project/internal/testutils"

	"github.com/dracory/taskstore"
	"github.com/dracory/test"
	"github.com/dracory/userstore"
)

// updateTestPayload mirrors the JSON body expected by handleUserUpdateAjax
type updateTestPayload struct {
	UserID       string `json:"user_id"`
	Status       string `json:"status"`
	Role         string `json:"role"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	BusinessName string `json:"business_name"`
	Phone        string `json:"phone"`
	Country      string `json:"country"`
	Timezone     string `json:"timezone"`
	Memo         string `json:"memo"`
}

// callUpdateAjax invokes the controller's Handler with the user-update-ajax action and a JSON body
func callUpdateAjax(t *testing.T, controller *userUpdateController, method string, payload any, query url.Values) (string, *http.Response) {
	t.Helper()
	if query == nil {
		query = url.Values{}
	}
	query.Set("action", actionUserUpdate)

	opts := test.NewRequestOptions{
		GetValues: query,
	}

	if rawBody, ok := payload.(string); ok {
		opts.Body = rawBody
		opts.ContentType = "application/json"
	} else if payload != nil {
		opts.JSONData = payload
	}

	body, response, err := test.CallStringEndpoint(method, controller.Handler, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return body, response
}

// parseApiResponse decodes a JSON api.Response body into a map
func parseApiResponse(t *testing.T, body string) map[string]any {
	t.Helper()
	var apiResponse map[string]any
	if err := json.Unmarshal([]byte(body), &apiResponse); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	return apiResponse
}

// TestHandleUserUpdateAjaxMethodCheck verifies that handleUserUpdateAjax rejects invalid JSON bodies
func TestHandleUserUpdateAjaxMethodCheck(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithGeoStore(true),
	)

	controller := NewUserUpdateController(app)
	body, response := callUpdateAjax(t, controller, http.MethodPost, "not-valid-json", nil)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	apiResponse := parseApiResponse(t, body)
	if apiResponse["status"] != "error" {
		t.Errorf("expected status error, got %v", apiResponse["status"])
	}
	if apiResponse["message"] != "Invalid request body" {
		t.Errorf("expected message 'Invalid request body', got %v", apiResponse["message"])
	}
}

// TestHandleUserUpdateAjaxPayloadValidation verifies that handleUserUpdateAjax requires user_id
func TestHandleUserUpdateAjaxPayloadValidation(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithGeoStore(true),
	)

	controller := NewUserUpdateController(app)
	body, response := callUpdateAjax(t, controller, http.MethodPost, updateTestPayload{}, nil)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	apiResponse := parseApiResponse(t, body)
	if apiResponse["status"] != "error" {
		t.Errorf("expected status error, got %v", apiResponse["status"])
	}
	if apiResponse["message"] != "User ID is required" {
		t.Errorf("expected message 'User ID is required', got %v", apiResponse["message"])
	}
}

// TestHandleUserUpdateAjaxUserIDValidation verifies that handleUserUpdateAjax reports unknown users
func TestHandleUserUpdateAjaxUserIDValidation(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithGeoStore(true),
	)

	controller := NewUserUpdateController(app)
	body, response := callUpdateAjax(t, controller, http.MethodPost, updateTestPayload{UserID: "non-existent-id"}, nil)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	apiResponse := parseApiResponse(t, body)
	if apiResponse["status"] != "error" {
		t.Errorf("expected status error, got %v", apiResponse["status"])
	}
	if apiResponse["message"] != "User not found" {
		t.Errorf("expected message 'User not found', got %v", apiResponse["message"])
	}
}

// TestHandleUserUpdateAjaxUserLookup verifies that handleUserUpdateAjax validates required fields
// after a successful user lookup
func TestHandleUserUpdateAjaxUserLookup(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithGeoStore(true),
	)

	user := seedUpdateTestUser(t, app)

	controller := NewUserUpdateController(app)
	body, response := callUpdateAjax(t, controller, http.MethodPost, updateTestPayload{UserID: user.GetID()}, nil)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	apiResponse := parseApiResponse(t, body)
	if apiResponse["status"] != "error" {
		t.Errorf("expected status error, got %v", apiResponse["status"])
	}
	if apiResponse["message"] != "Status is required" {
		t.Errorf("expected message 'Status is required', got %v", apiResponse["message"])
	}
}

// TestHandleUserUpdateAjaxFieldValidation verifies that handleUserUpdateAjax requires first/last name and email
func TestHandleUserUpdateAjaxFieldValidation(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithGeoStore(true),
	)

	user := seedUpdateTestUser(t, app)

	controller := NewUserUpdateController(app)
	body, response := callUpdateAjax(t, controller, http.MethodPost, updateTestPayload{
		UserID: user.GetID(),
		Status: userstore.USER_STATUS_ACTIVE,
	}, nil)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	apiResponse := parseApiResponse(t, body)
	if apiResponse["status"] != "error" {
		t.Errorf("expected status error, got %v", apiResponse["status"])
	}
	if apiResponse["message"] != "First name is required" {
		t.Errorf("expected message 'First name is required', got %v", apiResponse["message"])
	}
}

// TestHandleUserUpdateAjaxRoleValidation verifies that handleUserUpdateAjax rejects invalid role values
func TestHandleUserUpdateAjaxRoleValidation(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithGeoStore(true),
	)

	user := seedUpdateTestUser(t, app)

	controller := NewUserUpdateController(app)
	body, response := callUpdateAjax(t, controller, http.MethodPost, updateTestPayload{
		UserID:    user.GetID(),
		Status:    userstore.USER_STATUS_ACTIVE,
		FirstName: "Updated",
		LastName:  "User",
		Email:     "updated@example.com",
		Country:   "US",
		Timezone:  "America/New_York",
		Role:      "super-admin",
	}, nil)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	apiResponse := parseApiResponse(t, body)
	if apiResponse["status"] != "error" {
		t.Errorf("expected status error, got %v", apiResponse["status"])
	}
	if apiResponse["message"] != "Invalid role value" {
		t.Errorf("expected message 'Invalid role value', got %v", apiResponse["message"])
	}
}

// TestHandleUserUpdateAjaxEmailValidation verifies that handleUserUpdateAjax rejects invalid emails
func TestHandleUserUpdateAjaxEmailValidation(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithGeoStore(true),
	)

	user := seedUpdateTestUser(t, app)

	controller := NewUserUpdateController(app)
	body, response := callUpdateAjax(t, controller, http.MethodPost, updateTestPayload{
		UserID:    user.GetID(),
		Status:    userstore.USER_STATUS_ACTIVE,
		FirstName: "Updated",
		LastName:  "User",
		Email:     "not-an-email",
		Country:   "US",
		Timezone:  "America/New_York",
	}, nil)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	apiResponse := parseApiResponse(t, body)
	if apiResponse["status"] != "error" {
		t.Errorf("expected status error, got %v", apiResponse["status"])
	}
	if apiResponse["message"] != "Invalid email address" {
		t.Errorf("expected message 'Invalid email address', got %v", apiResponse["message"])
	}
}

// TestHandleUserUpdateAjaxUserUpdate verifies that handleUserUpdateAjax updates a user successfully
func TestHandleUserUpdateAjaxUserUpdate(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
		testutils.WithGeoStore(true),
	)

	user := seedUpdateTestUser(t, app)

	controller := NewUserUpdateController(app)
	body, response := callUpdateAjax(t, controller, http.MethodPost, updateTestPayload{
		UserID:    user.GetID(),
		Status:    userstore.USER_STATUS_ACTIVE,
		Role:      userstore.USER_ROLE_ADMINISTRATOR,
		FirstName: "Updated",
		LastName:  "Name",
		Email:     "updated@example.com",
		Country:   "US",
		Timezone:  "America/New_York",
		Phone:     "+1234567890",
		Memo:      "updated memo",
	}, nil)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	apiResponse := parseApiResponse(t, body)
	if apiResponse["status"] != "success" {
		t.Errorf("expected status success, got %v", apiResponse["status"])
	}
	if apiResponse["message"] != "User saved successfully" {
		t.Errorf("expected message 'User saved successfully', got %v", apiResponse["message"])
	}

	// Verify the user was actually updated in the database
	updated, err := app.GetUserStore().UserFindByID(context.Background(), user.GetID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.GetFirstName() != "Updated" {
		t.Errorf("expected first name 'Updated', got %q", updated.GetFirstName())
	}
	if updated.GetLastName() != "Name" {
		t.Errorf("expected last name 'Name', got %q", updated.GetLastName())
	}
	if updated.GetEmail() != "updated@example.com" {
		t.Errorf("expected email 'updated@example.com', got %q", updated.GetEmail())
	}
	if updated.GetRole() != userstore.USER_ROLE_ADMINISTRATOR {
		t.Errorf("expected role %q, got %q", userstore.USER_ROLE_ADMINISTRATOR, updated.GetRole())
	}
	if updated.GetMemo() != "updated memo" {
		t.Errorf("expected memo 'updated memo', got %q", updated.GetMemo())
	}
}

// TestHandleUserUpdateAjaxVaultTokenization verifies that handleUserUpdateAjax tokenizes fields
// when the vault store is enabled
func TestHandleUserUpdateAjaxVaultTokenization(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
		testutils.WithGeoStore(true),
	)

	user := seedUpdateTestUser(t, app)
	tokenizeUpdateUser(t, app, user, "Original", "User", "original@example.com")

	controller := NewUserUpdateController(app)
	body, response := callUpdateAjax(t, controller, http.MethodPost, updateTestPayload{
		UserID:    user.GetID(),
		Status:    userstore.USER_STATUS_ACTIVE,
		FirstName: "Vaulted",
		LastName:  "User",
		Email:     "vaulted@example.com",
		Country:   "US",
		Timezone:  "America/New_York",
	}, nil)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	apiResponse := parseApiResponse(t, body)
	if apiResponse["status"] != "success" {
		t.Errorf("expected status success, got %v", apiResponse["status"])
	}

	// The stored first name should be a token, not the plaintext value
	updated, err := app.GetUserStore().UserFindByID(context.Background(), user.GetID())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.GetFirstName() == "Vaulted" {
		t.Error("first name should be tokenized in storage")
	}
}

// TestHandleUserUpdateAjaxBlindIndexUpdate verifies that handleUserUpdateAjax enqueues a blind index
// rebuild task when the email changes and vault + task stores are configured
func TestHandleUserUpdateAjaxBlindIndexUpdate(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
		testutils.WithTaskStore(true),
		testutils.WithGeoStore(true),
	)

	user := seedUpdateTestUser(t, app)
	tokenizeUpdateUser(t, app, user, "Original", "User", "original@example.com")

	// Register the BlindIndexUpdate task definition so the controller can enqueue it
	if err := app.GetTaskStore().TaskHandlerAdd(context.Background(), blind_index_rebuild.NewBlindIndexRebuildTask(app), true); err != nil {
		t.Fatalf("TaskHandlerAdd returned error: %v", err)
	}

	controller := NewUserUpdateController(app)
	body, response := callUpdateAjax(t, controller, http.MethodPost, updateTestPayload{
		UserID:    user.GetID(),
		Status:    userstore.USER_STATUS_ACTIVE,
		FirstName: "Changed",
		LastName:  "Email",
		Email:     "changed-email@example.com",
		Country:   "US",
		Timezone:  "America/New_York",
	}, nil)

	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	apiResponse := parseApiResponse(t, body)
	if apiResponse["status"] != "success" {
		t.Errorf("expected status success, got %v", apiResponse["status"])
	}

	// Verify a blind index rebuild task was enqueued
	tasks, err := app.GetTaskStore().TaskQueueList(context.Background(), taskstore.TaskQueueQuery())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) == 0 {
		t.Error("a blind index rebuild task should have been enqueued")
	}
}

// seedUpdateTestUser creates a test user in the in-memory database and returns it
func seedUpdateTestUser(t *testing.T, app interface {
	GetUserStore() userstore.StoreInterface
}) userstore.UserInterface {
	t.Helper()
	user := userstore.NewUser()
	user.SetFirstName("Original")
	user.SetLastName("User")
	user.SetEmail("original@example.com")
	user.SetStatus(userstore.USER_STATUS_ACTIVE)
	user.SetRole(userstore.USER_ROLE_USER)

	if err := app.GetUserStore().UserCreate(context.Background(), user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}
	return user
}

// tokenizeUpdateUser stores vault tokens for the user's fields so the update handler can
// untokenize the existing values before applying changes. It creates fresh tokens for each
// field (the user starts with no tokens) and persists the tokenized user to the store.
func tokenizeUpdateUser(t *testing.T, app app.AppInterface, user userstore.UserInterface, first, last, email string) {
	t.Helper()
	ctx := context.Background()
	vaultStore := app.GetVaultStore()
	vaultKey := app.GetConfig().GetVaultStoreKey()

	ensureToken := func(value, field string) string {
		token, err := vaultStore.TokenCreate(ctx, value, vaultKey, 20)
		if err != nil {
			t.Fatalf("TokenCreate (%s) returned error: %v", field, err)
		}
		return token
	}

	firstToken := ensureToken(first, "first name")
	lastToken := ensureToken(last, "last name")
	emailToken := ensureToken(email, "email")

	user.SetFirstName(firstToken)
	user.SetLastName(lastToken)
	user.SetEmail(emailToken)
	if err := app.GetUserStore().UserUpdate(ctx, user); err != nil {
		t.Fatalf("UserUpdate returned error: %v", err)
	}
}
