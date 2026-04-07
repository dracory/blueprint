package post_create

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"project/internal/config"
	"project/internal/registry"
	"project/internal/testutils"

	"github.com/dracory/blogstore"
	"github.com/dracory/test"
	"github.com/dracory/userstore"
)

func TestPostCreateController_RequiresAuthentication(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	responseHTML, _, err := test.CallStringEndpoint(http.MethodPost, NewPostCreateController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"post_title": {"Test Post"},
		},
	})

	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if !strings.Contains(responseHTML, "You are not logged in") {
		t.Error("Should show login required message")
	}
}

func TestPostCreateController_RequiresPostTitle(t *testing.T) {
	registry, user := setupControllerAppAndUser(t)

	responseHTML, _, err := test.CallStringEndpoint(http.MethodPost, NewPostCreateController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if !strings.Contains(responseHTML, "post title is required") {
		t.Error("Should show title required message")
	}
}

func TestPostCreateController_ShowsFormOnGet(t *testing.T) {
	registry, user := setupControllerAppAndUser(t)

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewPostCreateController(registry).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Should return 200 status, got %d", response.StatusCode)
	}
	if !strings.Contains(responseHTML, "name=\"post_title\"") {
		t.Error("Should show post title input")
	}
}

func TestPostCreateController_CreatesPostSuccessfully(t *testing.T) {
	registry, user := setupControllerAppAndUser(t)
	postTitle := "Test Post Title"

	responseHTML, _, err := test.CallStringEndpoint(http.MethodPost, NewPostCreateController(registry).Handler, test.NewRequestOptions{
		PostValues: url.Values{
			"post_title": {postTitle},
		},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if !strings.Contains(responseHTML, "post created successfully") {
		t.Error("Should show success message")
	}

	// Verify post was created
	posts, err := registry.GetBlogStore().PostList(context.Background(), blogstore.PostQueryOptions{})
	if err != nil {
		t.Errorf("Should list posts without error: %v", err)
	}
	if len(posts) == 0 {
		t.Error("Should have created a post")
	}
	if posts[0].GetTitle() != postTitle {
		t.Errorf("Post title should match, expected %s, got %s", postTitle, posts[0].GetTitle())
	}
}

func setupControllerAppAndUser(t *testing.T) (registry.RegistryInterface, userstore.UserInterface) {
	t.Helper()

	registry := testutils.Setup(
		testutils.WithUserStore(true),
		testutils.WithBlogStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Fatalf("SeedUser returned error: %v", err)
	}

	return registry, user
}
