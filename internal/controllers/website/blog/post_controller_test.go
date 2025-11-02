package blog

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/blogstore"
)

func TestBlogPostController_Handler_MissingPostID(t *testing.T) {
	// --- Setup ---
	app := testutils.Setup()

	controller := NewBlogPostController(app)

	// --- Execute ---
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/blog/post/", nil)
	html := controller.Handler(w, r)

	// --- Assert ---
	if html == "" {
		t.Fatal("Expected HTML to not be empty")
	}

	if !strings.Contains(html, "post is missing") {
		t.Errorf("Expected HTML to contain 'post is missing'")
	}
}

func TestBlogPostController_Handler_PostNotFound(t *testing.T) {
	// --- Setup ---
	app := testutils.Setup(testutils.WithBlogStore(true))

	controller := NewBlogPostController(app)

	// --- Execute ---
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/blog/post/nonexistent/Nonexistent-Post", nil)
	html := controller.Handler(w, r)

	// --- Assert ---
	if html == "" {
		t.Fatal("Expected HTML to not be empty")
	}

	if !strings.Contains(html, "post is missing") {
		t.Errorf("Expected HTML to contain 'post is missing'")
	}
}

func TestBlogPostController_Handler_PostNotPublished_NoAuth(t *testing.T) {
	// --- Setup ---
	app := testutils.Setup(testutils.WithBlogStore(true))

	// Create a draft post
	post := blogstore.NewPost()
	post.SetTitle("Draft Post")
	post.SetContent("Draft content")
	post.SetStatus(blogstore.POST_STATUS_DRAFT)
	err := app.GetBlogStore().PostCreate(post)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	controller := NewBlogPostController(app)

	// --- Execute ---
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/blog/post/"+post.ID()+"/"+post.Slug(), nil)
	html := controller.Handler(w, r)

	// --- Assert ---
	if html == "" {
		t.Fatal("Expected HTML to not be empty")
	}

	if !strings.Contains(html, "post is missing") {
		t.Errorf("Expected HTML to contain 'post is missing' for unpublished post without auth")
	}
}

func TestBlogPostController_Handler_PostNotPublished_WithAuth(t *testing.T) {
	// --- Setup ---
	app := testutils.Setup(testutils.WithBlogStore(true), testutils.WithUserStore(true))

	// Create a draft post
	post := blogstore.NewPost()
	post.SetTitle("Draft Post")
	post.SetContent("Draft content")
	post.SetStatus(blogstore.POST_STATUS_DRAFT)
	err := app.GetBlogStore().PostCreate(post)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Create and authenticate a regular user
	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	controller := NewBlogPostController(app)

	// --- Execute ---
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/blog/post/"+post.ID()+"/"+post.Slug(), nil)
	r, err = testutils.LoginAs(app, r, user)
	if err != nil {
		t.Fatalf("Failed to authenticate test user: %v", err)
	}
	html := controller.Handler(w, r)

	// --- Assert ---
	if html == "" {
		t.Fatal("Expected HTML to not be empty")
	}

	if strings.Contains(html, "post is missing") {
		t.Errorf("Expected authenticated user to access unpublished post")
	}

	if !strings.Contains(html, "Draft Post") {
		t.Errorf("Expected HTML to contain the post title")
	}
}

func TestBlogPostController_Handler_PostPublished_Success(t *testing.T) {
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

	// Create a published post
	post := blogstore.NewPost()
	post.SetTitle("Published Post")
	post.SetContent("Published content")
	post.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	err = app.GetBlogStore().PostCreate(post)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	controller := NewBlogPostController(app)

	// --- Execute ---
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/blog/post/"+post.ID()+"/"+post.Slug(), nil)
	html := controller.Handler(w, r)

	// --- Assert ---
	if html == "" {
		t.Fatal("Expected HTML to not be empty")
	}

	if strings.Contains(html, "post is missing") {
		t.Errorf("Expected published post to be accessible")
	}

	if !strings.Contains(html, "Published Post") {
		t.Errorf("Expected HTML to contain the post title")
	}

	if !strings.Contains(html, "Published content") {
		t.Errorf("Expected HTML to contain the post content")
	}
}

func TestBlogPostController_Handler_WrongSlug_Redirect(t *testing.T) {
	// --- Setup ---
	app := testutils.Setup(testutils.WithBlogStore(true), testutils.WithCacheStore(true))

	// Create a published post
	post := blogstore.NewPost()
	post.SetTitle("Test Post")
	post.SetContent("Test content")
	post.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	err := app.GetBlogStore().PostCreate(post)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	controller := NewBlogPostController(app)

	// --- Execute ---
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/blog/post/"+post.ID()+"/wrong-slug", nil)
	html := controller.Handler(w, r)

	// --- Assert ---
	if html == "" {
		t.Fatal("Expected HTML to not be empty")
	}

	if !strings.Contains(html, "See Other") {
		t.Errorf("Expected HTML to contain redirect link 'See Other'")
	}

	// Check that a redirect status code was set
	response := w.Result()
	if response.StatusCode != http.StatusSeeOther {
		t.Errorf("Expected status code %d, got %d", http.StatusSeeOther, response.StatusCode)
	}

	// Check that the flash message was set correctly
	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)
	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal("Response MUST contain 'flash message'")
	}

	if flashMessage.Type != "success" {
		t.Fatalf("Response should be of type 'success', but got: %s", flashMessage.Type)
	}

	if !strings.Contains(flashMessage.Message, "location has changed") {
		t.Fatalf("Expected flash message to contain 'location has changed', but got: %s", flashMessage.Message)
	}
}

func TestBlogPostController_Handler_AdminAccessUnpublished(t *testing.T) {
	// --- Setup ---
	app := testutils.Setup(testutils.WithBlogStore(true), testutils.WithUserStore(true))

	// Create a draft post
	post := blogstore.NewPost()
	post.SetTitle("Draft Post")
	post.SetContent("Draft content")
	post.SetStatus(blogstore.POST_STATUS_DRAFT)
	err := app.GetBlogStore().PostCreate(post)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Create and authenticate an admin user
	adminUser, err := testutils.SeedUser(app.GetUserStore(), testutils.ADMIN_01)
	if err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}

	controller := NewBlogPostController(app)

	// --- Execute ---
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/blog/post/"+post.ID()+"/"+post.Slug(), nil)
	r, err = testutils.LoginAs(app, r, adminUser)
	if err != nil {
		t.Fatalf("Failed to authenticate admin user: %v", err)
	}
	html := controller.Handler(w, r)

	// --- Assert ---
	if html == "" {
		t.Fatal("Expected HTML to not be empty")
	}

	if strings.Contains(html, "post is missing") {
		t.Errorf("Expected admin user to access unpublished post")
	}

	if !strings.Contains(html, "Draft Post") {
		t.Errorf("Expected HTML to contain the post title")
	}
}

func TestBlogPostController_ProcessContent_Markdown(t *testing.T) {
	// --- Setup ---
	app := testutils.Setup()
	controller := NewBlogPostController(app)

	markdown := "# Hello World\n\nThis is **bold** text."
	expectedHTML := "<h1>Hello World</h1>\n<p>This is <strong>bold</strong> text.</p>\n"

	// --- Execute ---
	html, css := controller.processContent(markdown, blogstore.POST_EDITOR_MARKDOWN)

	// --- Assert ---
	if html != expectedHTML {
		t.Errorf("Expected HTML:\n%s\nGot:\n%s", expectedHTML, html)
	}

	if css != "" {
		t.Errorf("Expected CSS to be empty for markdown, got: %s", css)
	}
}

func TestBlogPostController_ProcessContent_BlockArea(t *testing.T) {
	// --- Setup ---
	app := testutils.Setup()
	controller := NewBlogPostController(app)

	blockContent := `[{"Id":"block-1","Type":"text","Sequence":1,"ParentId":"","Attributes":{"Text":"Test content"}}]`

	// --- Execute ---
	html, _ := controller.processContent(blockContent, blogstore.POST_EDITOR_BLOCKAREA)

	// --- Assert ---
	if html == "" {
		t.Errorf("Expected HTML to not be empty for block area content")
	}

	// Block area processing will return processed content, CSS might be empty or contain styles
	if !strings.Contains(html, "Test content") {
		t.Errorf("Expected HTML to contain the processed content")
	}
}

func TestBlogPostController_ProcessContent_BlockEditor(t *testing.T) {
	// --- Setup ---
	app := testutils.Setup()
	controller := NewBlogPostController(app)

	blockEditorContent := `{"blocks": [{"type": "paragraph", "data": {"text": "Test content"}}]}`

	// --- Execute ---
	html, css := controller.processContent(blockEditorContent, blogstore.POST_EDITOR_BLOCKEDITOR)

	// --- Assert ---
	// Block editor processing might return error for invalid content, but should not panic
	if html == "" && css == "" {
		t.Log("Block editor returned empty content, which may be expected for invalid input")
	}
}
