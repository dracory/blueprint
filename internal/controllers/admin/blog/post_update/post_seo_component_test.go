package post_update

import (
	"context"
	"net/url"
	"testing"

	"project/internal/registry"
	"project/internal/testutils"

	"github.com/dracory/blogstore"
)

func setupSEOTestAppAndPost(t *testing.T) (registry.RegistryInterface, blogstore.PostInterface) {
	t.Helper()

	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithBlogStore(true),
	)

	post := blogstore.NewPost()
	post.SetTitle("SEO Test Post")
	post.SetStatus(blogstore.POST_STATUS_DRAFT)
	post.SetCanonicalURL("https://example.com/original")
	post.SetMetaDescription("Original meta description")
	post.SetMetaKeywords("foo,bar")
	post.SetMetaRobots("INDEX, FOLLOW")

	if err := registry.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	return registry, post
}

func TestPostSEOComponent_MountRequiresPostID(t *testing.T) {
	registry, _ := setupSEOTestAppAndPost(t)

	c := &postSEOComponent{registry: registry}

	err := c.Mount(context.Background(), map[string]string{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.FormErrorMessage != "Post ID is required" {
		t.Errorf("expected FormErrorMessage 'Post ID is required', got %s", c.FormErrorMessage)
	}
}

func TestPostSEOComponent_MountLoadsPostFields(t *testing.T) {
	registry, post := setupSEOTestAppAndPost(t)

	c := &postSEOComponent{registry: registry}

	err := c.Mount(context.Background(), map[string]string{
		"post_id": post.GetID(),
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.PostID != post.GetID() {
		t.Errorf("expected PostID %s, got %s", post.GetID(), c.PostID)
	}
	if c.FormCanonicalURL != post.GetCanonicalURL() {
		t.Errorf("expected FormCanonicalURL %s, got %s", post.GetCanonicalURL(), c.FormCanonicalURL)
	}
	if c.FormMetaDescription != post.GetMetaDescription() {
		t.Errorf("expected FormMetaDescription %s, got %s", post.GetMetaDescription(), c.FormMetaDescription)
	}
	if c.FormMetaKeywords != post.GetMetaKeywords() {
		t.Errorf("expected FormMetaKeywords %s, got %s", post.GetMetaKeywords(), c.FormMetaKeywords)
	}
	if c.FormMetaRobots != post.GetMetaRobots() {
		t.Errorf("expected FormMetaRobots %s, got %s", post.GetMetaRobots(), c.FormMetaRobots)
	}
	if c.FormErrorMessage != "" {
		t.Errorf("expected empty FormErrorMessage, got %s", c.FormErrorMessage)
	}
}

func TestPostSEOComponent_HandleSave_UpdatesPostAndSetsSuccess(t *testing.T) {
	registry, post := setupSEOTestAppAndPost(t)

	c := &postSEOComponent{registry: registry, PostID: post.GetID()}

	values := url.Values{
		"post_canonical_url":    {"https://example.com/updated"},
		"post_meta_description": {"Updated meta"},
		"post_meta_keywords":    {"one,two"},
		"post_meta_robots":      {"NOINDEX, FOLLOW"},
	}

	err := c.Handle(context.Background(), "save", values)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.FormErrorMessage != "" {
		t.Errorf("expected empty FormErrorMessage, got %s", c.FormErrorMessage)
	}
	if c.FormSuccessMessage != "Post saved successfully" {
		t.Errorf("expected FormSuccessMessage 'Post saved successfully', got %s", c.FormSuccessMessage)
	}

	// Reload post from store and verify fields were updated
	updated, err := registry.GetBlogStore().PostFindByID(context.Background(), post.GetID())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if updated == nil {
		t.Fatal("expected updated to not be nil")
	}
	if updated.GetCanonicalURL() != "https://example.com/updated" {
		t.Errorf("expected canonical URL 'https://example.com/updated', got %s", updated.GetCanonicalURL())
	}
	if updated.GetMetaDescription() != "Updated meta" {
		t.Errorf("expected meta description 'Updated meta', got %s", updated.GetMetaDescription())
	}
	if updated.GetMetaKeywords() != "one,two" {
		t.Errorf("expected meta keywords 'one,two', got %s", updated.GetMetaKeywords())
	}
	if updated.GetMetaRobots() != "NOINDEX, FOLLOW" {
		t.Errorf("expected meta robots 'NOINDEX, FOLLOW', got %s", updated.GetMetaRobots())
	}
}
