package blog

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"project/internal/testutils"
	"strings"
	"testing"

	"github.com/dracory/blogstore"
)

func TestBlogController_Handler_Success(t *testing.T) {
	// --- Setup ---
	cfg := testutils.DefaultConf()
	cfg.SetBlogStoreUsed(true)
	cfg.SetCmsStoreUsed(true)
	cfg.SetCmsStoreTemplateID("test-template")
	app := testutils.Setup(testutils.WithCfg(cfg))
	
	// Create a test template
	err := testutils.SeedTemplate(app.GetCmsStore(), "test-site", "test-template")
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}
	
	controller := NewBlogController(app)

	// Create test posts in the database
	post1 := blogstore.NewPost()
	post1.SetTitle("Post 1")
	post1.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	err = app.GetBlogStore().PostCreate(post1)
	if err != nil {
		t.Fatalf("Failed to create test post 1: %v", err)
	}

	post2 := blogstore.NewPost()
	post2.SetTitle("Post 2")
	post2.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	err = app.GetBlogStore().PostCreate(post2)
	if err != nil {
		t.Fatalf("Failed to create test post 2: %v", err)
	}

	// Create a draft post that should NOT appear on the page
	draftPost := blogstore.NewPost()
	draftPost.SetTitle("Draft Post")
	draftPost.SetStatus(blogstore.POST_STATUS_DRAFT)
	err = app.GetBlogStore().PostCreate(draftPost)
	if err != nil {
		t.Fatalf("Failed to create draft post: %v", err)
	}

	// --- Execute ---
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/blog", nil)
	html := controller.Handler(w, r)

	// --- Assert ---
	if html == "" {
		t.Fatal("Expected HTML to not be empty")
	}
	
	if !strings.Contains(html, "Blog") {
		t.Errorf("Expected HTML to contain 'Blog'")
	}
	if !strings.Contains(html, "Post 1") {
		t.Errorf("Expected HTML to contain 'Post 1'")
	}
	if !strings.Contains(html, "Post 2") {
		t.Errorf("Expected HTML to contain 'Post 2'")
	}
	if strings.Contains(html, "Draft Post") {
		t.Errorf("Expected HTML to NOT contain 'Draft Post'")
	}
	if !strings.Contains(html, "pagination-primary-soft") {
		t.Errorf("Expected HTML to contain 'pagination-primary-soft'")
	}

	// Ensure the flash error redirect was not triggered
	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Result().StatusCode)
	}
}

// fakeBlogStore is a wrapper to simulate database errors.
type fakeBlogStore struct {
	blogstore.StoreInterface // Embed the real blogstore
	postListError            error
	postCountError           error
}

func (f *fakeBlogStore) PostList(options blogstore.PostQueryOptions) ([]blogstore.Post, error) {
	if f.postListError != nil {
		return nil, f.postListError
	}
	return f.StoreInterface.PostList(options)
}

func (f *fakeBlogStore) PostCount(options blogstore.PostQueryOptions) (int64, error) {
	if f.postCountError != nil {
		return 0, f.postCountError
	}
	return f.StoreInterface.PostCount(options)
}

func TestBlogController_Handler_PostListError(t *testing.T) {
	// --- Setup ---
	cfg := testutils.DefaultConf()
	cfg.SetBlogStoreUsed(true)
	cfg.SetCacheStoreUsed(true)
	cfg.SetCmsStoreUsed(true)
	cfg.SetCmsStoreTemplateID("test-template")
	app := testutils.Setup(testutils.WithCfg(cfg))
	
	// Create a test template
	err := testutils.SeedTemplate(app.GetCmsStore(), "test-site", "test-template")
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}
	app.SetBlogStore(&fakeBlogStore{
		StoreInterface: app.GetBlogStore(),
		postListError:  errors.New("simulated database error"),
	})
	controller := NewBlogController(app)

	// --- Execute ---
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/blog", nil)
	html := controller.Handler(w, r)

	// --- Assert ---
	// The handler should return a flash error redirect.
	// The actual HTML will be a redirect link from helpers.ToFlashError.
	if html == "" {
		t.Fatal("Expected HTML to not be empty")
	}
	if !strings.Contains(html, "See Other") {
		t.Errorf("Expected HTML to contain a redirect link 'See Other'")
	}

	// Check that a redirect status code was set on the response writer.
	response := w.Result()
	if response.StatusCode != http.StatusSeeOther {
		t.Errorf("Expected status code %d, got %d", http.StatusSeeOther, response.StatusCode)
	}

	// Check that the flash message was set correctly.
	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)

	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal("Response MUST contain 'flash message'")
	}

	if flashMessage.Type != "error" {
		t.Fatalf("Response be of type 'error', but got: %s %s", flashMessage.Type, flashMessage.Message)
	}

	if !strings.Contains(flashMessage.Message, "error") {
		t.Fatalf("Expected flash message to contain 'error', but got: %s", flashMessage.Message)
	}
}

func TestBlogController_Handler_PostCountError(t *testing.T) {
	// --- Setup ---
	cfg := testutils.DefaultConf()
	cfg.SetBlogStoreUsed(true)
	cfg.SetCacheStoreUsed(true)
	cfg.SetCmsStoreUsed(true)
	cfg.SetCmsStoreTemplateID("test-template")
	app := testutils.Setup(testutils.WithCfg(cfg))
	
	// Create a test template
	err := testutils.SeedTemplate(app.GetCmsStore(), "test-site", "test-template")
	if err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}
	// This fake store will only error on PostCount
	app.SetBlogStore(&fakeBlogStore{
		StoreInterface: app.GetBlogStore(),
		postCountError: errors.New("database count error"),
	})
	controller := NewBlogController(app)

	// --- Execute ---
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/blog", nil)
	html := controller.Handler(w, r)

	// --- Assert ---
	if html == "" {
		t.Fatal("Expected HTML to not be empty")
	}
	if !strings.Contains(html, "See Other") {
		t.Errorf("Expected HTML to contain a redirect link 'See Other'")
	}

	response := w.Result()
	if response.StatusCode != http.StatusSeeOther {
		t.Errorf("Expected status code %d, got %d", http.StatusSeeOther, response.StatusCode)
	}

	// Check that the flash message was set correctly.
	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)

	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal("Response MUST contain 'flash message'")
	}

	if flashMessage.Type != "error" {
		t.Fatalf("Response be of type 'error', but got: %s %s", flashMessage.Type, flashMessage.Message)
	}

	if !strings.Contains(flashMessage.Message, "error") {
		t.Fatalf("Expected flash message to contain 'error', but got: %s", flashMessage.Message)
	}
}

func TestBlogController_PageRendering(t *testing.T) {
	// --- Setup ---
	app := testutils.Setup()
	controller := NewBlogController(app)

	data := blogControllerData{
		postList: []blogstore.Post{
			*blogstore.NewPost().SetID("1").SetTitle("My First Post").SetSummary("A summary.").SetImageUrl("http://example.com/img.png"),
		},
		postCount: 10,
		page:      0,
		perPage:   12,
	}

	// --- Execute ---
	html := controller.page(data)

	// --- Assert ---
	if html == "" {
		t.Fatal("Expected HTML to not be empty")
	}
	
	if !strings.Contains(html, "My First Post") {
		t.Errorf("Expected post title to be in the HTML")
	}
	if !strings.Contains(html, "A summary.") {
		t.Errorf("Expected post summary to be in the HTML")
	}
	if !strings.Contains(html, "/th/png/300x200/80/http/example.com/img.png") {
		t.Errorf("Expected thumbnail URL to be in the HTML")
	}
	if !strings.Contains(html, `pagination-primary-soft`) {
		t.Errorf("Expected pagination to be rendered")
	}
}
