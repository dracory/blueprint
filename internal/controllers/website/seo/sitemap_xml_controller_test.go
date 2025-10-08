package seo

import (
	"net/http"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/blogstore"
	"github.com/dracory/test"
)

func TestSitemapXmlController_NoBlogStore(t *testing.T) {
	app := testutils.Setup()
	controller := NewSitemapXmlController(app.GetBlogStore())

	body, response, err := test.CallStringEndpoint(http.MethodGet, controller.Handler, test.NewRequestOptions{})

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	if !strings.Contains(body, "<urlset") {
		t.Fatalf("expected sitemap xml body, got: %s", body)
	}

	if strings.Contains(body, "/blog/post/") {
		t.Fatalf("expected no blog post entries when blog store is nil, got: %s", body)
	}
}

func TestSitemapXmlController_WithBlogStore(t *testing.T) {
	app := testutils.Setup(testutils.WithBlogStore(true))
	store := app.GetBlogStore()
	if store == nil {
		t.Fatal("expected blog store to be initialized")
	}

	post := blogstore.NewPost().
		SetID("post-1").
		SetTitle("first-post").
		SetStatus(blogstore.POST_STATUS_PUBLISHED)

	if err := store.PostCreate(post); err != nil {
		t.Fatalf("failed to create post: %v", err)
	}

	controller := NewSitemapXmlController(store)
	body, response, err := test.CallStringEndpoint(http.MethodGet, controller.Handler, test.NewRequestOptions{})

	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, response.StatusCode)
	}

	expectedLoc := "/blog/post/post-1/first-post"
	if !strings.Contains(body, expectedLoc) {
		t.Fatalf("expected sitemap to contain %s, got: %s", expectedLoc, body)
	}
}
