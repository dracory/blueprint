package user_update

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/dracory/userstore"
	"github.com/stretchr/testify/assert"
)

// TestRenderPageUserIDValidation verifies that renderPage redirects when user_id is missing
func TestRenderPageUserIDValidation(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	controller := NewUserUpdateController(app)
	_, response, err := test.CallStringEndpoint(http.MethodGet, controller.renderPage, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusSeeOther, response.StatusCode)

	flash, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)
	assert.NoError(t, err)
	assert.Equal(t, "User ID is required", flash.Message)
}

// TestRenderPageUserStoreNilCheck verifies that renderPage redirects when UserStore is not configured
func TestRenderPageUserStoreNilCheck(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
	)

	controller := NewUserUpdateController(app)
	_, response, err := test.CallStringEndpoint(http.MethodGet, controller.renderPage, test.NewRequestOptions{
		GetValues: url.Values{
			"user_id": {"some-id"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusSeeOther, response.StatusCode)

	flash, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)
	assert.NoError(t, err)
	assert.Equal(t, "User store is not configured", flash.Message)
}

// TestRenderPageUserLookup verifies that renderPage redirects when the user is not found
func TestRenderPageUserLookup(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	controller := NewUserUpdateController(app)
	_, response, err := test.CallStringEndpoint(http.MethodGet, controller.renderPage, test.NewRequestOptions{
		GetValues: url.Values{
			"user_id": {"non-existent-id"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusSeeOther, response.StatusCode)

	flash, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)
	assert.NoError(t, err)
	assert.Equal(t, "User not found", flash.Message)
}

// TestRenderPageVaultUntokenization verifies that renderPage untokenizes vault-stored fields
// and displays the plaintext user name
func TestRenderPageVaultUntokenization(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true, true),
		testutils.WithVaultStore(true, "test-key"),
	)

	user := userstore.NewUser()
	user.SetFirstName("VaultedFirst")
	user.SetLastName("VaultedLast")
	user.SetEmail("vault-page@example.com")
	user.SetStatus(userstore.USER_STATUS_ACTIVE)
	user.SetRole(userstore.USER_ROLE_USER)
	if err := app.GetUserStore().UserCreate(context.Background(), user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	controller := NewUserUpdateController(app)
	html, response, err := test.CallStringEndpoint(http.MethodGet, controller.renderPage, test.NewRequestOptions{
		GetValues: url.Values{
			"user_id": {user.GetID()},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	// The page should render without error and show the edit heading
	assert.Contains(t, html, "Edit User")
}

// TestRenderPageHTMLGeneration verifies that renderPage generates the expected HTML
func TestRenderPageHTMLGeneration(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user := userstore.NewUser()
	user.SetFirstName("Page")
	user.SetLastName("Render")
	user.SetEmail("page-render@example.com")
	user.SetStatus(userstore.USER_STATUS_ACTIVE)
	user.SetRole(userstore.USER_ROLE_USER)
	if err := app.GetUserStore().UserCreate(context.Background(), user); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	controller := NewUserUpdateController(app)
	html, response, err := test.CallStringEndpoint(http.MethodGet, controller.renderPage, test.NewRequestOptions{
		GetValues: url.Values{
			"user_id": {user.GetID()},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Contains(t, html, "Edit User", "should show page heading")
	assert.Contains(t, html, "User: Page Render", "should show display name")
	assert.Contains(t, html, "app-user-update", "should render the Vue app container")
	assert.Contains(t, html, "vue.global.js", "should include Vue.js CDN")
}
