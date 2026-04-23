package post_manager

import (
	"context"
	"encoding/json"
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
	// Check that Vue.js app is rendered
	if !strings.Contains(authResponse, "blog-posts-app") {
		t.Error("Should render Vue.js app container")
	}
	if !strings.Contains(authResponse, "vue.global.js") {
		t.Error("Should include Vue.js CDN")
	}
}

func TestManagerController_RendersVueApp(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

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
	// Check for Vue.js app container
	if !strings.Contains(responseHTML, "id=\"blog-posts-app\"") {
		t.Error("Should render Vue.js app container")
	}
	// Check for Vue CDN
	if !strings.Contains(responseHTML, "vue.global.js") {
		t.Error("Should include Vue.js CDN")
	}
	// Check for SweetAlert2 CDN
	if !strings.Contains(responseHTML, "sweetalert2") {
		t.Error("Should include SweetAlert2 CDN")
	}
}

func TestManagerController_HandleLoadPosts(t *testing.T) {
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

	// Test load-posts API endpoint
	response, responseObj, err := test.CallStringEndpoint(http.MethodGet, NewPostManagerController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {"load-posts"},
		},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if responseObj.StatusCode != http.StatusOK {
		t.Errorf("Should return 200 status, got %d", responseObj.StatusCode)
	}

	// Parse JSON response
	var apiResponse map[string]any
	if err := json.Unmarshal([]byte(response), &apiResponse); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if apiResponse["status"] != "success" {
		t.Errorf("API should return success status, got: %s", apiResponse["status"])
	}

	// Check that posts are in the response
	data, ok := apiResponse["data"].(map[string]any)
	if !ok {
		t.Fatal("API response data should be a map")
	}

	posts, ok := data["posts"].([]any)
	if !ok {
		t.Fatal("API response should contain posts array")
	}

	if len(posts) != 2 {
		t.Errorf("Should return 2 posts, got %d", len(posts))
	}
}

func TestManagerController_HandleLoadPostsWithFilters(t *testing.T) {
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

	// Test status filter via API
	response, responseObj, err := test.CallStringEndpoint(http.MethodGet, NewPostManagerController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {"load-posts"},
			"status": {blogstore.POST_STATUS_PUBLISHED},
		},
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if responseObj.StatusCode != http.StatusOK {
		t.Errorf("Should return 200 status, got %d", responseObj.StatusCode)
	}

	// Parse JSON response
	var apiResponse map[string]any
	if err := json.Unmarshal([]byte(response), &apiResponse); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	data, ok := apiResponse["data"].(map[string]any)
	if !ok {
		t.Fatal("API response data should be a map")
	}

	posts, ok := data["posts"].([]any)
	if !ok {
		t.Fatal("API response should contain posts array")
	}

	if len(posts) != 1 {
		t.Errorf("Should return 1 published post, got %d", len(posts))
	}
}

func TestManagerController_HandleLoadPostsWithPagination(t *testing.T) {
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

	// Test pagination via API
	response, responseObj, err := test.CallStringEndpoint(http.MethodGet, NewPostManagerController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action":   {"load-posts"},
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
	if responseObj.StatusCode != http.StatusOK {
		t.Errorf("Should return 200 status, got %d", responseObj.StatusCode)
	}

	// Parse JSON response
	var apiResponse map[string]any
	if err := json.Unmarshal([]byte(response), &apiResponse); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	data, ok := apiResponse["data"].(map[string]any)
	if !ok {
		t.Fatal("API response data should be a map")
	}

	posts, ok := data["posts"].([]any)
	if !ok {
		t.Fatal("API response should contain posts array")
	}

	if len(posts) != 5 {
		t.Errorf("Should return 5 posts on page 2 (15 total, 10 per page), got %d", len(posts))
	}

	total, ok := data["total"]
	if !ok {
		t.Error("API response should contain total count")
	}
	if cast.ToInt(total) != 15 {
		t.Errorf("Total should be 15, got %d", cast.ToInt(total))
	}
}

func TestManagerController_HandleCreatePost(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("Should create test user: %v", err)
	}

	// Test create-post API endpoint
	requestBody := `{"title":"New Test Post"}`
	response, responseObj, err := test.CallStringEndpoint(http.MethodPost, NewPostManagerController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {"create-post"},
		},
		Body: requestBody,
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if responseObj.StatusCode != http.StatusOK {
		t.Errorf("Should return 200 status, got %d", responseObj.StatusCode)
	}

	// Parse JSON response
	var apiResponse map[string]any
	if err := json.Unmarshal([]byte(response), &apiResponse); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if apiResponse["status"] != "success" {
		t.Errorf("API should return success status, got: %s", apiResponse["status"])
	}

	data, ok := apiResponse["data"].(map[string]any)
	if !ok {
		t.Fatal("API response data should be a map")
	}

	postID, ok := data["id"]
	if !ok {
		t.Error("API response should contain post ID")
	}
	if postID == "" {
		t.Error("Post ID should not be empty")
	}
}

func TestManagerController_HandleDeletePost(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Create test post
	post := blogstore.NewPost()
	post.SetTitle("Post to Delete")
	post.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	if err := registry.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("Should create test user: %v", err)
	}

	// Test delete-post API endpoint
	requestBody := `{"post_id":"` + post.GetID() + `"}`
	response, responseObj, err := test.CallStringEndpoint(http.MethodPost, NewPostManagerController(registry).Handler, test.NewRequestOptions{
		GetValues: url.Values{
			"action": {"delete-post"},
		},
		Body: requestBody,
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if responseObj.StatusCode != http.StatusOK {
		t.Errorf("Should return 200 status, got %d", responseObj.StatusCode)
	}

	// Parse JSON response
	var apiResponse map[string]any
	if err := json.Unmarshal([]byte(response), &apiResponse); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if apiResponse["status"] != "success" {
		t.Errorf("API should return success status, got: %s", apiResponse["status"])
	}

	// Verify post was deleted or moved to trash
	deletedPost, err := registry.GetBlogStore().PostFindByID(context.Background(), post.GetID())
	if err == nil && deletedPost != nil {
		// Post still exists, check if it's trashed
		if deletedPost.GetStatus() != blogstore.POST_STATUS_TRASH {
			t.Error("Post should be deleted or moved to trash")
		}
	}
	// If err != nil or deletedPost == nil, the post was successfully deleted (not soft-deleted)
}
