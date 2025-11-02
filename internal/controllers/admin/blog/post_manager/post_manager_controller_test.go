package post_manager

import (
	"net/http"
	"net/url"
	"testing"

	"project/internal/config"
	"project/internal/testutils"

	"github.com/dracory/blogstore"
	"github.com/dracory/test"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestManagerController_RequiresAuthentication(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Test without authentication
	response, responseObj, err := test.CallStringEndpoint(http.MethodGet, NewPostManagerController(app).Handler, test.NewRequestOptions{})
	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusSeeOther, responseObj.StatusCode, "Should redirect when unauthenticated")
	assert.Contains(t, response, "See Other", "Should show redirect response")

	// Test with authentication
	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	assert.NoError(t, err, "Should create test user")

	authResponse, authResponseObj, err := test.CallStringEndpoint(http.MethodGet, NewPostManagerController(app).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusOK, authResponseObj.StatusCode, "Should return 200 when authenticated")
	assert.NotContains(t, authResponse, "See Other", "Should not redirect when authenticated")
}

func TestManagerController_ShowsPostList(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Create test posts
	post1 := blogstore.NewPost()
	post1.SetTitle("Test Post 1")
	post1.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	app.GetBlogStore().PostCreate(post1)

	post2 := blogstore.NewPost()
	post2.SetTitle("Test Post 2")
	post2.SetStatus(blogstore.POST_STATUS_DRAFT)
	app.GetBlogStore().PostCreate(post2)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	assert.NoError(t, err, "Should create test user")

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewPostManagerController(app).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusOK, response.StatusCode, "Should return 200 status")
	assert.Contains(t, responseHTML, "Test Post 1", "Should show first post")
	assert.Contains(t, responseHTML, "Test Post 2", "Should show second post")
	assert.Contains(t, responseHTML, "New Post", "Should show create button")
}

func TestManagerController_HandlesFilters(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Create test posts
	post1 := blogstore.NewPost()
	post1.SetTitle("Published Post")
	post1.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	app.GetBlogStore().PostCreate(post1)

	post2 := blogstore.NewPost()
	post2.SetTitle("Draft Post")
	post2.SetStatus(blogstore.POST_STATUS_DRAFT)
	app.GetBlogStore().PostCreate(post2)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	assert.NoError(t, err, "Should create test user")

	// Test status filter
	responseHTML, _, err := test.CallStringEndpoint(http.MethodGet, NewPostManagerController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"status": {blogstore.POST_STATUS_PUBLISHED},
		},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Contains(t, responseHTML, "Published Post", "Should show published post")
	assert.NotContains(t, responseHTML, "Draft Post", "Should not show draft post with filter")
}

func TestManagerController_HandlesPagination(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Create test posts
	for i := 1; i <= 15; i++ {
		post := blogstore.NewPost()
		post.SetTitle("Post " + cast.ToString(i))
		post.SetStatus(blogstore.POST_STATUS_PUBLISHED)
		app.GetBlogStore().PostCreate(post)
	}

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	assert.NoError(t, err, "Should create test user")

	// Test pagination by requesting page 2
	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewPostManagerController(app).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page":     {"1"},
			"per_page": {"10"},
		},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusOK, response.StatusCode, "Should return 200 status")
	assert.Contains(t, responseHTML, "pagination", "Should show pagination controls")
}
