package home

import (
	"context"
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
	err = app.GetBlogStore().PostCreate(context.Background(), post1)
	if err != nil {
		t.Fatalf("Failed to create test post 1: %v", err)
	}

	post2 := blogstore.NewPost()
	post2.SetTitle("Post 2")
	post2.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	err = app.GetBlogStore().PostCreate(context.Background(), post2)
	if err != nil {
		t.Fatalf("Failed to create test post 2: %v", err)
	}

	// Create a draft post that should NOT appear on the page
	draftPost := blogstore.NewPost()
	draftPost.SetTitle("Draft Post")
	draftPost.SetStatus(blogstore.POST_STATUS_DRAFT)
	err = app.GetBlogStore().PostCreate(context.Background(), draftPost)
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
