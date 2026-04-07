package post_update

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"project/internal/registry"
	"project/internal/testutils"

	"github.com/dracory/blogstore"
	"github.com/dracory/test"
)

func TestPostUpdateController_RequiresPostID(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithBlogStore(true),
	)

	_, response, err := test.CallStringEndpoint(http.MethodGet, NewPostUpdateController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if response.StatusCode != http.StatusSeeOther {
		t.Errorf("Should redirect with error, expected %d, got %d", http.StatusSeeOther, response.StatusCode)
	}

	// Verify flash message was set
	flash, err := testutils.FlashMessageFindFromResponse(registry.GetCacheStore(), response)
	if err != nil {
		t.Errorf("Should find flash message: %v", err)
	}
	if flash.Message != "Post ID is required" {
		t.Errorf("Should show correct error message, got %s", flash.Message)
	}
}

func TestPostUpdateController_InvalidPostID(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithBlogStore(true),
	)

	_, response, err := test.CallStringEndpoint(http.MethodGet, NewPostUpdateController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"post_id": {"invalid_id"},
		},
	})

	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if response.StatusCode != http.StatusSeeOther {
		t.Errorf("Should redirect with error, expected %d, got %d", http.StatusSeeOther, response.StatusCode)
	}

	// Verify flash message was set
	flash, err := testutils.FlashMessageFindFromResponse(registry.GetCacheStore(), response)
	if err != nil {
		t.Errorf("Should find flash message: %v", err)
	}
	if flash.Message != "Post not found" {
		t.Errorf("Should show correct error message, got %s", flash.Message)
	}
}

func TestPostUpdateController_ShowsPage(t *testing.T) {
	registry, post := setupControllerAppAndPost(t)

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewPostUpdateController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"post_id": {post.GetID()},
			"view":    {"content"},
		},
	})

	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Should return 200 status, got %d", response.StatusCode)
	}
	if !strings.Contains(responseHTML, "Edit Post") {
		t.Error("Should show page heading")
	}
	if !strings.Contains(responseHTML, "Post:") {
		t.Error("Should show post label")
	}
	if !strings.Contains(responseHTML, post.GetTitle()) {
		t.Error("Should show post title")
	}
}

func setupControllerAppAndPost(t *testing.T) (registry.RegistryInterface, blogstore.PostInterface) {
	// Note: we reuse the same pattern as v1 tests but only for GET behavior.
	t.Helper()

	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithBlogStore(true),
	)

	post := blogstore.NewPost()
	post.SetTitle("Test Post")
	post.SetContent("Test Content")
	post.SetStatus(blogstore.POST_STATUS_DRAFT)

	if err := registry.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	return registry, post
}
