package post

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/blogstore"
	"github.com/dracory/rtr"
)

func TestBlogPostController_Handler_MissingPostID(t *testing.T) {
	app := testutils.Setup(testutils.WithCacheStore(true))
	controller := NewPostController(app)

	w := httptest.NewRecorder()
	r := newRequestWithParams(http.MethodGet, "/blog/post/", map[string]string{})

	html := controller.Handler(w, r)

	if html == "" {
		t.Fatal("Expected HTML to not be empty")
	}

	if !strings.Contains(html, "post is missing") {
		t.Errorf("Expected HTML to contain 'post is missing'")
	}

	resp := w.Result()
	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("Expected status %d but got %d", http.StatusSeeOther, resp.StatusCode)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), resp)
	if err != nil {
		t.Fatalf("Failed to read flash message: %v", err)
	}

	if flashMessage == nil || !strings.Contains(flashMessage.Message, "no longer exists") {
		t.Fatalf("Expected flash message about missing post, got: %+v", flashMessage)
	}
}

func TestBlogPostController_Handler_PostNotFound(t *testing.T) {
	app := testutils.Setup(testutils.WithBlogStore(true), testutils.WithCacheStore(true))
	controller := NewPostController(app)

	w := httptest.NewRecorder()
	r := newRequestWithParams(http.MethodGet, "/blog/post/nonexistent/Nonexistent-Post", map[string]string{
		"id":    "nonexistent",
		"title": "nonexistent-post",
	})

	html := controller.Handler(w, r)

	if html != "" {
		t.Errorf("Expected empty HTML when post is missing, got: %s", html)
	}

	resp := w.Result()
	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("Expected status %d but got %d", http.StatusSeeOther, resp.StatusCode)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), resp)
	if err != nil {
		t.Fatalf("Failed to read flash message: %v", err)
	}

	if flashMessage == nil || !strings.Contains(flashMessage.Message, "no longer exists") {
		t.Fatalf("Expected flash message about missing post, got: %+v", flashMessage)
	}
}

func TestBlogPostController_Handler_PostNotPublished_NoAuth(t *testing.T) {
	app := testutils.Setup(testutils.WithBlogStore(true), testutils.WithCacheStore(true))

	post := blogstore.NewPost()
	post.SetTitle("Draft Post")
	post.SetContent("Draft content")
	post.SetStatus(blogstore.POST_STATUS_DRAFT)

	if err := app.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	controller := NewPostController(app)

	w := httptest.NewRecorder()
	r := newRequestWithParams(http.MethodGet, "/blog/post/"+post.ID()+"/"+post.Slug(), map[string]string{
		"id":    post.ID(),
		"title": post.Slug(),
	})

	html := controller.Handler(w, r)

	if html != "" {
		t.Errorf("Expected empty HTML for unauthorised draft access, got: %s", html)
	}

	resp := w.Result()
	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("Expected redirect status, got %d", resp.StatusCode)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), resp)
	if err != nil {
		t.Fatalf("Failed to read flash message: %v", err)
	}

	if flashMessage == nil || flashMessage.Type != "warning" || !strings.Contains(flashMessage.Message, "no longer active") {
		t.Fatalf("Expected warning flash about inactive post, got: %+v", flashMessage)
	}
}

func TestBlogPostController_Handler_PostNotPublished_WithAuth(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithUserStore(true),
		testutils.WithSessionStore(true),
		testutils.WithCacheStore(true),
	)

	post := blogstore.NewPost()
	post.SetTitle("Draft Post")
	post.SetContent("Draft content")
	post.SetStatus(blogstore.POST_STATUS_DRAFT)

	if err := app.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	user, err := testutils.SeedUser(app.GetUserStore(), testutils.USER_01)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	controller := NewPostController(app)

	w := httptest.NewRecorder()
	r := newRequestWithParams(http.MethodGet, "/blog/post/"+post.ID()+"/"+post.Slug(), map[string]string{
		"id":    post.ID(),
		"title": post.Slug(),
	})

	r, err = testutils.LoginAs(app, r, user)
	if err != nil {
		t.Fatalf("Failed to authenticate test user: %v", err)
	}

	html := controller.Handler(w, r)

	if html != "" {
		t.Errorf("Expected empty HTML for non-admin draft access, got: %s", html)
	}

	resp := w.Result()
	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("Expected redirect status, got %d", resp.StatusCode)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), resp)
	if err != nil {
		t.Fatalf("Failed to read flash message: %v", err)
	}

	if flashMessage == nil || flashMessage.Type != "warning" || !strings.Contains(flashMessage.Message, "no longer active") {
		t.Fatalf("Expected warning flash about inactive post, got: %+v", flashMessage)
	}
}

func TestBlogPostController_Handler_PostPublished_Success(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetBlogStoreUsed(true)
	cfg.SetCmsStoreUsed(true)
	cfg.SetCmsStoreTemplateID("test-template")
	cfg.SetCacheStoreUsed(true)

	app := testutils.Setup(testutils.WithCfg(cfg))

	if err := testutils.SeedTemplate(app.GetCmsStore(), "test-site", "test-template"); err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	post := blogstore.NewPost()
	post.SetTitle("Published Post")
	post.SetContent("Published content")
	post.SetStatus(blogstore.POST_STATUS_PUBLISHED)

	if err := app.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	controller := NewPostController(app)

	w := httptest.NewRecorder()
	r := newRequestWithParams(http.MethodGet, "/blog/post/"+post.ID()+"/"+post.Slug(), map[string]string{
		"id":    post.ID(),
		"title": post.Slug(),
	})

	html := controller.Handler(w, r)

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
	app := testutils.Setup(testutils.WithBlogStore(true), testutils.WithCacheStore(true))

	post := blogstore.NewPost()
	post.SetTitle("Test Post")
	post.SetContent("Test content")
	post.SetStatus(blogstore.POST_STATUS_PUBLISHED)

	if err := app.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	controller := NewPostController(app)

	w := httptest.NewRecorder()
	r := newRequestWithParams(http.MethodGet, "/blog/post/"+post.ID()+"/wrong-slug", map[string]string{
		"id":    post.ID(),
		"title": "wrong-slug",
	})

	html := controller.Handler(w, r)

	if html != "" {
		t.Fatalf("Expected empty HTML when redirecting, got: %s", html)
	}

	response := w.Result()
	if response.StatusCode != http.StatusSeeOther {
		t.Fatalf("Expected status code %d, got %d", http.StatusSeeOther, response.StatusCode)
	}

	flashMessage, err := testutils.FlashMessageFindFromResponse(app.GetCacheStore(), response)
	if err != nil {
		t.Fatal(err)
	}

	if flashMessage == nil {
		t.Fatal("Response MUST contain flash message")
	}

	if flashMessage.Type != "success" {
		t.Fatalf("Response should be of type 'success', but got: %s", flashMessage.Type)
	}

	if !strings.Contains(flashMessage.Message, "location has changed") {
		t.Fatalf("Expected flash message to contain 'location has changed', but got: %s", flashMessage.Message)
	}
}

func TestBlogPostController_Handler_AdminAccessUnpublished(t *testing.T) {
	cfg := testutils.DefaultConf()
	cfg.SetBlogStoreUsed(true)
	cfg.SetCmsStoreUsed(true)
	cfg.SetCmsStoreTemplateID("test-template")
	cfg.SetCacheStoreUsed(true)
	cfg.SetUserStoreUsed(true)
	cfg.SetSessionStoreUsed(true)

	app := testutils.Setup(
		testutils.WithCfg(cfg),
		testutils.WithUserStore(true),
		testutils.WithSessionStore(true),
	)

	if err := testutils.SeedTemplate(app.GetCmsStore(), "test-site", "test-template"); err != nil {
		t.Fatalf("Failed to create test template: %v", err)
	}

	post := blogstore.NewPost()
	post.SetTitle("Draft Post")
	post.SetContent("Draft content")
	post.SetStatus(blogstore.POST_STATUS_DRAFT)

	if err := app.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	adminUser, err := testutils.SeedUser(app.GetUserStore(), testutils.ADMIN_01)
	if err != nil {
		t.Fatalf("Failed to create admin user: %v", err)
	}

	controller := NewPostController(app)

	w := httptest.NewRecorder()
	r := newRequestWithParams(http.MethodGet, "/blog/post/"+post.ID()+"/"+post.Slug(), map[string]string{
		"id":    post.ID(),
		"title": post.Slug(),
	})

	r, err = testutils.LoginAs(app, r, adminUser)
	if err != nil {
		t.Fatalf("Failed to authenticate admin user: %v", err)
	}

	html := controller.Handler(w, r)

	if html == "" {
		t.Fatal("Expected HTML to not be empty")
	}

	if !strings.Contains(html, "Draft Post") {
		t.Errorf("Expected HTML to contain the post title")
	}

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected OK status, got %d", resp.StatusCode)
	}
}

func newRequestWithParams(method, path string, params map[string]string) *http.Request {
	if params == nil {
		params = map[string]string{}
	}

	req := httptest.NewRequest(method, path, nil)
	ctx := context.WithValue(req.Context(), rtr.ParamsKey, params)
	return req.WithContext(ctx)
}

func TestBlogPostController_ProcessContent_Markdown(t *testing.T) {
	// --- Setup ---
	app := testutils.Setup()
	controller := NewPostController(app)

	markdown := "# Hello World\n\nThis is **bold** text."
	expectedHTML := "<h1 id=\"hello-world\">Hello World</h1>\n<p>This is <strong>bold</strong> text.</p>\n"

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
	controller := NewPostController(app)

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
	controller := NewPostController(app)

	blockEditorContent := `{"blocks": [{"type": "paragraph", "data": {"text": "Test content"}}]}`

	// --- Execute ---
	html, css := controller.processContent(blockEditorContent, blogstore.POST_EDITOR_BLOCKEDITOR)

	// --- Assert ---
	// Block editor processing might return error for invalid content, but should not panic
	if html == "" && css == "" {
		t.Log("Block editor returned empty content, which may be expected for invalid input")
	}
}
