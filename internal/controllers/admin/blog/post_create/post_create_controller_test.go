package post_create

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"project/internal/config"
	"project/internal/testutils"
	"project/internal/types"

	"github.com/dracory/blogstore"
	"github.com/dracory/test"
	"github.com/dracory/userstore"
	"github.com/stretchr/testify/assert"
)

func TestPostCreateController_RequiresAuthentication(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	responseHTML, _, err := test.CallStringEndpoint(http.MethodPost, NewPostCreateController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"post_title": {"Test Post"},
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Contains(t, responseHTML, "You are not logged in", "Should show login required message")
}

func TestPostCreateController_RequiresPostTitle(t *testing.T) {
	app, user := setupControllerAppAndUser(t)

	responseHTML, _, err := test.CallStringEndpoint(http.MethodPost, NewPostCreateController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Contains(t, responseHTML, "post title is required", "Should show title required message")
}

func TestPostCreateController_ShowsFormOnGet(t *testing.T) {
	app, user := setupControllerAppAndUser(t)

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewPostCreateController(app).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusOK, response.StatusCode, "Should return 200 status")
	assert.Contains(t, responseHTML, "name=\"post_title\"", "Should show post title input")
}

func TestPostCreateController_CreatesPostSuccessfully(t *testing.T) {
	app, user := setupControllerAppAndUser(t)
	postTitle := "Test Post Title"

	responseHTML, _, err := test.CallStringEndpoint(http.MethodPost, NewPostCreateController(app).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"post_title": {postTitle},
		},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Contains(t, responseHTML, "post created successfully", "Should show success message")

	// Verify post was created
	posts, err := app.GetBlogStore().PostList(context.Background(), blogstore.PostQueryOptions{})
	assert.NoError(t, err, "Should list posts without error")
	assert.NotEmpty(t, posts, "Should have created a post")
	assert.Equal(t, postTitle, posts[0].Title(), "Post title should match")
}

func setupControllerAppAndUser(t *testing.T) (types.RegistryInterface, userstore.UserInterface) {
	t.Helper()

	app := testutils.Setup(
		testutils.WithUserStore(true),
		testutils.WithBlogStore(true),
	)

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	return app, user
}
