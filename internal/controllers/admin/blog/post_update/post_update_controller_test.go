package post_update

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"project/internal/testutils"
	"project/internal/types"

	"github.com/dracory/blogstore"
	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
)

func TestPostUpdateController_RequiresPostID(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithBlogStore(true),
	)

	_, response, err := test.CallStringEndpoint(http.MethodGet, NewPostUpdateController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusSeeOther, response.StatusCode, "Should redirect with error")

	// Verify flash message was set
	flash, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)
	assert.NoError(t, err, "Should find flash message")
	assert.Equal(t, "Post ID is required", flash.Message, "Should show correct error message")
}

func TestPostUpdateController_InvalidPostID(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithBlogStore(true),
	)

	_, response, err := test.CallStringEndpoint(http.MethodGet, NewPostUpdateController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"post_id": {"invalid_id"},
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusSeeOther, response.StatusCode, "Should redirect with error")

	// Verify flash message was set
	flash, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)
	assert.NoError(t, err, "Should find flash message")
	assert.Equal(t, "Post not found", flash.Message, "Should show correct error message")
}

func TestPostUpdateController_ShowsPage(t *testing.T) {
	app, post := setupControllerAppAndPost(t)

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewPostUpdateController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"post_id": {post.ID()},
			"view":    {"content"},
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusOK, response.StatusCode, "Should return 200 status")
	assert.Contains(t, responseHTML, "Edit Post", "Should show page heading")
	assert.Contains(t, responseHTML, "Post:", "Should show post label")
	assert.Contains(t, responseHTML, post.Title(), "Should show post title")
}

func setupControllerAppAndPost(t *testing.T) (types.RegistryInterface, *blogstore.Post) {
	// Note: we reuse the same pattern as v1 tests but only for GET behavior.
	t.Helper()

	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithBlogStore(true),
	)

	post := blogstore.NewPost()
	post.SetTitle("Test Post")
	post.SetContent("Test Content")
	post.SetStatus(blogstore.POST_STATUS_DRAFT)

	if err := app.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	return app, post
}
