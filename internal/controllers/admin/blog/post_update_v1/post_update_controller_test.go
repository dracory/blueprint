package post_update_v1

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

func TestPostUpdateController_ShowsForm(t *testing.T) {
	app, post := setupControllerAppAndPost(t)

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewPostUpdateController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"post_id": {post.ID()},
			"view":    {VIEW_DETAILS},
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusOK, response.StatusCode, "Should return 200 status")
	assert.Contains(t, responseHTML, "FormPostUpdate", "Should show update form")
	assert.Contains(t, responseHTML, post.Title(), "Should show post title")
}

func TestPostUpdateController_HandlesPostRequest(t *testing.T) {
	app, post := setupControllerAppAndPost(t)

	responseHTML, _, err := test.CallStringEndpoint(http.MethodPost, NewPostUpdateController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"post_id":      {post.ID()},
			"view":         {VIEW_DETAILS},
			"post_status":  {blogstore.POST_STATUS_PUBLISHED},
			"post_title":   {post.Title()},
			"post_content": {post.Content()},
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Contains(t, responseHTML, "Post saved successfully", "Should show success message")
}

func TestPostUpdateController_ValidatesRequiredFields(t *testing.T) {
	app, post := setupControllerAppAndPost(t)

	// Test empty title in CONTENT view
	responseHTML, _, err := test.CallStringEndpoint(http.MethodPost, NewPostUpdateController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"post_id":    {post.ID()},
			"view":       {VIEW_CONTENT},
			"post_title": {""}, // Empty title
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Contains(t, responseHTML, "Title is required", "Should show title required error")

	// Test empty status in DETAILS view
	responseHTML, _, err = test.CallStringEndpoint(http.MethodPost, NewPostUpdateController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"post_id":     {post.ID()},
			"view":        {VIEW_DETAILS},
			"post_status": {""}, // Empty status
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Contains(t, responseHTML, "Status is required", "Should show status required error")
}

func setupControllerAppAndPost(t *testing.T) (types.AppInterface, *blogstore.Post) {
	t.Helper()

	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithBlogStore(true),
	)

	post := blogstore.NewPost()
	post.SetTitle("Test Post")
	post.SetContent("Test Content")
	post.SetStatus(blogstore.POST_STATUS_DRAFT)

	err := app.GetBlogStore().PostCreate(context.Background(), post)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	return app, post
}
