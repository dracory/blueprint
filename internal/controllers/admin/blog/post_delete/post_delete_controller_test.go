package post_delete

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"project/internal/config"
	"project/internal/testutils"

	"github.com/dracory/blogstore"
	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
)

func TestPostDeleteController_RequiresAuthentication(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	response, _, err := test.CallStringEndpoint(http.MethodGet, NewPostDeleteController(registry).Handler, test.NewRequestOptions{})
	assert.NoError(t, err, "Handler should not return error")
	assert.Contains(t, response, "not logged in", "Should require authentication")
}

func TestPostDeleteController_RequiresPostID(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	assert.NoError(t, err, "Should create test user")

	response, _, err := test.CallStringEndpoint(http.MethodGet, NewPostDeleteController(registry).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	assert.NoError(t, err, "Handler should not return error")
	assert.Contains(t, response, "post id is required", "Should require post ID")
}

func TestPostDeleteController_HandlesInvalidPostID(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	assert.NoError(t, err, "Should create test user")

	response, _, err := test.CallStringEndpoint(http.MethodGet, NewPostDeleteController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"post_id": {"invalid_id"},
		},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	assert.NoError(t, err, "Handler should not return error")
	assert.Contains(t, response, "Post not found", "Should handle invalid post ID")
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
	assert.NoError(t, err, "Should create test user")

	responseHTML, _, err := test.CallStringEndpoint(http.MethodGet, NewPostDeleteController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"post_id": {post.GetID()},
		},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	assert.NoError(t, err, "Handler should not return error")
	assert.Contains(t, responseHTML, "ModalPostDelete", "Should show delete modal")
	assert.Contains(t, responseHTML, post.GetID(), "Should include post ID in modal")
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
	assert.NoError(t, err, "Should create test user")

	// Send POST request to delete
	responseHTML, _, err := test.CallStringEndpoint(http.MethodPost, NewPostDeleteController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"post_id": {post.GetID()},
		},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	assert.NoError(t, err, "Handler should not return error")
	assert.Contains(t, responseHTML, "post deleted successfully", "Should show success message")

	// Verify post was marked as trash
	deletedPost, err := registry.GetBlogStore().PostFindByID(context.Background(), post.GetID())
	assert.NoError(t, err, "Should not return error when checking post")
	assert.Equal(t, blogstore.POST_STATUS_TRASH, deletedPost.GetStatus(), "Post should be marked as trash")
}
