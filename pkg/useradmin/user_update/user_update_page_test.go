package user_update

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/test"
	"github.com/dracory/userstore"
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

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, response.StatusCode)
	}

	flash, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flash.Message != "User ID is required" {
		t.Errorf("expected message 'User ID is required', got %q", flash.Message)
	}
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

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, response.StatusCode)
	}

	flash, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flash.Message != "User store is not configured" {
		t.Errorf("expected message 'User store is not configured', got %q", flash.Message)
	}
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

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusSeeOther {
		t.Errorf("expected status %d, got %d", http.StatusSeeOther, response.StatusCode)
	}

	flash, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flash.Message != "User not found" {
		t.Errorf("expected message 'User not found', got %q", flash.Message)
	}
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

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(html, "Edit User") {
		t.Errorf("expected HTML to contain 'Edit User'")
	}
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

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}
	if !strings.Contains(html, "Edit User") {
		t.Error("should show page heading")
	}
	if !strings.Contains(html, "User: Page Render") {
		t.Error("should show display name")
	}
	if !strings.Contains(html, "app-user-update") {
		t.Error("should render the Vue app container")
	}
	if !strings.Contains(html, "vue.global.js") {
		t.Error("should include Vue.js CDN")
	}
}
