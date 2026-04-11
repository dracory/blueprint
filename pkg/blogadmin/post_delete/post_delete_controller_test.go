package post_delete

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"project/internal/config"
	"project/internal/testutils"

	"github.com/dracory/blogstore"
	"github.com/dracory/test"
)

func TestPostDeleteController_RequiresAuthentication(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	response, _, err := test.CallStringEndpoint(http.MethodGet, NewPostDeleteController(registry).Handler, test.NewRequestOptions{})
	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if !strings.Contains(response, "not logged in") {
		t.Error("Should require authentication")
	}
}

func TestPostDeleteController_RequiresPostID(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("Should create test user: %v", err)
	}

	response, _, err := test.CallStringEndpoint(http.MethodGet, NewPostDeleteController(registry).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if !strings.Contains(response, "post id is required") {
		t.Error("Should require post ID")
	}
}

func TestPostDeleteController_HandlesInvalidPostID(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("Should create test user: %v", err)
	}

	response, _, err := test.CallStringEndpoint(http.MethodGet, NewPostDeleteController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"post_id": {"invalid_id"},
		},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if !strings.Contains(response, "Post not found") {
		t.Error("Should handle invalid post ID")
	}
}

func TestPostDeleteController_ShowsDeleteModal(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Create test post
	post := blogstore.NewPost()
	post.SetTitle("Test Post")
	post.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	if err := registry.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("Should create test user: %v", err)
	}

	responseHTML, _, err := test.CallStringEndpoint(http.MethodGet, NewPostDeleteController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"post_id": {post.GetID()},
		},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if !strings.Contains(responseHTML, "ModalPostDelete") {
		t.Error("Should show delete modal")
	}
	if !strings.Contains(responseHTML, post.GetID()) {
		t.Error("Should include post ID in modal")
	}
}

func TestPostDeleteController_DeletesPost(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Create test post
	post := blogstore.NewPost()
	post.SetTitle("Test Post")
	post.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	if err := registry.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("Should create test user: %v", err)
	}

	// Send POST request to delete
	responseHTML, _, err := test.CallStringEndpoint(http.MethodPost, NewPostDeleteController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"post_id": {post.GetID()},
		},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if !strings.Contains(responseHTML, "post deleted successfully") {
		t.Error("Should show success message")
	}

	// Verify post was marked as trash
	deletedPost, err := registry.GetBlogStore().PostFindByID(context.Background(), post.GetID())
	if err != nil {
		t.Errorf("Should not return error when checking post: %v", err)
	}
	if deletedPost.GetStatus() != blogstore.POST_STATUS_TRASH {
		t.Errorf("Post should be marked as trash, got %s", deletedPost.GetStatus())
	}
}

