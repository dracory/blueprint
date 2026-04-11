package post_manager

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
	"github.com/spf13/cast"
)

func TestManagerController_RequiresAuthentication(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Test without authentication
	response, responseObj, err := test.CallStringEndpoint(http.MethodGet, NewPostManagerController(registry).Handler, test.NewRequestOptions{})
	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if responseObj.StatusCode != http.StatusSeeOther {
		t.Errorf("Should redirect when unauthenticated, expected %d, got %d", http.StatusSeeOther, responseObj.StatusCode)
	}
	if !strings.Contains(response, "See Other") {
		t.Error("Should show redirect response")
	}

	// Test with authentication
	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("Should create test user: %v", err)
	}

	authResponse, authResponseObj, err := test.CallStringEndpoint(http.MethodGet, NewPostManagerController(registry).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if authResponseObj.StatusCode != http.StatusOK {
		t.Errorf("Should return 200 when authenticated, got %d", authResponseObj.StatusCode)
	}
	if strings.Contains(authResponse, "See Other") {
		t.Error("Should not redirect when authenticated")
	}
}

func TestManagerController_ShowsPostList(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Create test posts
	post1 := blogstore.NewPost()
	post1.SetTitle("Test Post 1")
	post1.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	if err := registry.GetBlogStore().PostCreate(context.Background(), post1); err != nil {
		t.Fatalf("failed to create test post1: %v", err)
	}

	post2 := blogstore.NewPost()
	post2.SetTitle("Test Post 2")
	post2.SetStatus(blogstore.POST_STATUS_DRAFT)
	if err := registry.GetBlogStore().PostCreate(context.Background(), post2); err != nil {
		t.Fatalf("failed to create test post2: %v", err)
	}

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("Should create test user: %v", err)
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewPostManagerController(registry).Handler, test.NewRequestOptions{
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
	if !strings.Contains(responseHTML, "Test Post 1") {
		t.Error("Should show first post")
	}
	if !strings.Contains(responseHTML, "Test Post 2") {
		t.Error("Should show second post")
	}
	if !strings.Contains(responseHTML, "New Post") {
		t.Error("Should show create button")
	}
}

func TestManagerController_HandlesFilters(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Create test posts
	post1 := blogstore.NewPost()
	post1.SetTitle("Published Post")
	post1.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	if err := registry.GetBlogStore().PostCreate(context.Background(), post1); err != nil {
		t.Fatalf("failed to create test post1: %v", err)
	}

	post2 := blogstore.NewPost()
	post2.SetTitle("Draft Post")
	post2.SetStatus(blogstore.POST_STATUS_DRAFT)
	if err := registry.GetBlogStore().PostCreate(context.Background(), post2); err != nil {
		t.Fatalf("failed to create test post2: %v", err)
	}

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("Should create test user: %v", err)
	}

	// Test status filter
	responseHTML, _, err := test.CallStringEndpoint(http.MethodGet, NewPostManagerController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"status": {blogstore.POST_STATUS_PUBLISHED},
		},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if !strings.Contains(responseHTML, "Published Post") {
		t.Error("Should show published post")
	}
	if strings.Contains(responseHTML, "Draft Post") {
		t.Error("Should not show draft post with filter")
	}
}

func TestManagerController_HandlesPagination(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Create test posts
	for i := 1; i <= 15; i++ {
		post := blogstore.NewPost()
		post.SetTitle("Post " + cast.ToString(i))
		post.SetStatus(blogstore.POST_STATUS_PUBLISHED)
		if err := registry.GetBlogStore().PostCreate(context.Background(), post); err != nil {
			t.Fatalf("failed to create test post %d: %v", i, err)
		}
	}

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("Should create test user: %v", err)
	}

	// Test pagination by requesting page 2
	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewPostManagerController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"page":     {"1"},
			"per_page": {"10"},
		},
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
	if !strings.Contains(responseHTML, "pagination") {
		t.Error("Should show pagination controls")
	}
}

