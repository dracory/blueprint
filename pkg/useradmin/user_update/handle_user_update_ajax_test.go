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
	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err)
	return body, response
}

// parseApiResponse decodes a JSON api.Response body into a map
func parseApiResponse(t *testing.T, body string) map[string]any {
	t.Helper()
	var apiResponse map[string]any
	assert.NoError(t, json.Unmarshal([]byte(body), &apiResponse))
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

	assert.Equal(t, http.StatusOK, response.StatusCode)
	apiResponse := parseApiResponse(t, body)
	assert.Equal(t, "error", apiResponse["status"])
	assert.Equal(t, "Invalid request body", apiResponse["message"])
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

	assert.Equal(t, http.StatusOK, response.StatusCode)
	apiResponse := parseApiResponse(t, body)
	assert.Equal(t, "error", apiResponse["status"])
	assert.Equal(t, "User ID is required", apiResponse["message"])
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

	assert.Equal(t, http.StatusOK, response.StatusCode)
	apiResponse := parseApiResponse(t, body)
	assert.Equal(t, "error", apiResponse["status"])
	assert.Equal(t, "User not found", apiResponse["message"])
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

	assert.Equal(t, http.StatusOK, response.StatusCode)
	apiResponse := parseApiResponse(t, body)
	assert.Equal(t, "error", apiResponse["status"])
	assert.Equal(t, "Status is required", apiResponse["message"])
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

	assert.Equal(t, http.StatusOK, response.StatusCode)
	apiResponse := parseApiResponse(t, body)
	assert.Equal(t, "error", apiResponse["status"])
	assert.Equal(t, "First name is required", apiResponse["message"])
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

	assert.Equal(t, http.StatusOK, response.StatusCode)
	apiResponse := parseApiResponse(t, body)
	assert.Equal(t, "error", apiResponse["status"])
	assert.Equal(t, "Invalid role value", apiResponse["message"])
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

	assert.Equal(t, http.StatusOK, response.StatusCode)
	apiResponse := parseApiResponse(t, body)
	assert.Equal(t, "error", apiResponse["status"])
	assert.Equal(t, "Invalid email address", apiResponse["message"])
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

	assert.Equal(t, http.StatusOK, response.StatusCode)
	apiResponse := parseApiResponse(t, body)
	assert.Equal(t, "success", apiResponse["status"])
	assert.Equal(t, "User saved successfully", apiResponse["message"])

	// Verify the user was actually updated in the database
	updated, err := app.GetUserStore().UserFindByID(context.Background(), user.GetID())
	assert.NoError(t, err)
	assert.Equal(t, "Updated", updated.GetFirstName())
	assert.Equal(t, "Name", updated.GetLastName())
	assert.Equal(t, "updated@example.com", updated.GetEmail())
	assert.Equal(t, userstore.USER_ROLE_ADMINISTRATOR, updated.GetRole())
	assert.Equal(t, "updated memo", updated.GetMemo())
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

	assert.Equal(t, http.StatusOK, response.StatusCode)
	apiResponse := parseApiResponse(t, body)
	assert.Equal(t, "success", apiResponse["status"])

	// The stored first name should be a token, not the plaintext value
	updated, err := app.GetUserStore().UserFindByID(context.Background(), user.GetID())
	assert.NoError(t, err)
	assert.NotEqual(t, "Vaulted", updated.GetFirstName(), "first name should be tokenized in storage")
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

	assert.Equal(t, http.StatusOK, response.StatusCode)
	apiResponse := parseApiResponse(t, body)
	assert.Equal(t, "success", apiResponse["status"])

	// Verify a blind index rebuild task was enqueued
	tasks, err := app.GetTaskStore().TaskQueueList(context.Background(), taskstore.TaskQueueQuery())
	assert.NoError(t, err)
	assert.NotEmpty(t, tasks, "a blind index rebuild task should have been enqueued")
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
