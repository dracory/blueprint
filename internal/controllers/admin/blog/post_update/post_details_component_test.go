package post_update

import (
	"context"
	"net/url"
	"testing"

	"project/internal/registry"
	"project/internal/testutils"

	"github.com/dracory/blogstore"
)

func setupDetailsTestAppAndPost(t *testing.T) (registry.RegistryInterface, blogstore.PostInterface) {
	t.Helper()

	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithBlogStore(true),
	)

	post := blogstore.NewPost()
	post.SetTitle("Details Test Post")
	post.SetStatus(blogstore.POST_STATUS_DRAFT)
	post.SetImageUrl("https://example.com/original.jpg")
	post.SetFeatured("no")
	post.SetPublishedAt("2024-01-02 03:04:05")
	post.SetEditor(blogstore.POST_EDITOR_HTMLAREA)
	post.SetMemo("Initial memo")

	if err := registry.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	return registry, post
}

func TestPostDetailsComponent_MountRequiresPostID(t *testing.T) {
	registry, _ := setupDetailsTestAppAndPost(t)

	c := &postDetailsComponent{registry: registry}

	err := c.Mount(context.Background(), map[string]string{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.FormErrorMessage != "Post ID is required" {
		t.Errorf("expected FormErrorMessage 'Post ID is required', got %s", c.FormErrorMessage)
	}
}

func TestPostDetailsComponent_MountLoadsPostFields(t *testing.T) {
	registry, post := setupDetailsTestAppAndPost(t)

	c := &postDetailsComponent{registry: registry}

	err := c.Mount(context.Background(), map[string]string{
		"post_id": post.GetID(),
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.PostID != post.GetID() {
		t.Errorf("expected PostID %s, got %s", post.GetID(), c.PostID)
	}
	if c.FormStatus != post.GetStatus() {
		t.Errorf("expected FormStatus %s, got %s", post.GetStatus(), c.FormStatus)
	}
	if c.FormImageUrl != post.GetImageUrl() {
		t.Errorf("expected FormImageUrl %s, got %s", post.GetImageUrl(), c.FormImageUrl)
	}
	if c.FormFeatured != post.GetFeatured() {
		t.Errorf("expected FormFeatured %s, got %s", post.GetFeatured(), c.FormFeatured)
	}
	if c.FormEditor != post.GetEditor() {
		t.Errorf("expected FormEditor %s, got %s", post.GetEditor(), c.FormEditor)
	}
	if c.FormMemo != post.GetMemo() {
		t.Errorf("expected FormMemo %s, got %s", post.GetMemo(), c.FormMemo)
	}
	if c.FormErrorMessage != "" {
		t.Errorf("expected empty FormErrorMessage, got %s", c.FormErrorMessage)
	}
}

func TestPostDetailsComponent_HandleSave_ValidatesStatusRequired(t *testing.T) {
	registry, post := setupDetailsTestAppAndPost(t)

	c := &postDetailsComponent{registry: registry, PostID: post.GetID()}

	values := url.Values{
		"post_status":       {""},
		"post_image_url":    {"https://example.com/updated.jpg"},
		"post_featured":     {"yes"},
		"post_published_at": {"2024-02-03 04:05"},
		"post_editor":       {blogstore.POST_EDITOR_MARKDOWN},
		"post_memo":         {"Updated memo"},
	}

	err := c.Handle(context.Background(), "save", values)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.FormErrorMessage != "Status is required" {
		t.Errorf("expected FormErrorMessage 'Status is required', got %s", c.FormErrorMessage)
	}
	if c.FormSuccessMessage != "" {
		t.Errorf("expected empty FormSuccessMessage, got %s", c.FormSuccessMessage)
	}
}

func TestPostDetailsComponent_HandleSave_UpdatesPostAndSetsSuccess(t *testing.T) {
	registry, post := setupDetailsTestAppAndPost(t)

	c := &postDetailsComponent{registry: registry, PostID: post.GetID()}

	values := url.Values{
		"post_status":       {blogstore.POST_STATUS_PUBLISHED},
		"post_image_url":    {"https://example.com/updated.jpg"},
		"post_featured":     {"yes"},
		"post_published_at": {"2024-02-03 04:05"},
		"post_editor":       {blogstore.POST_EDITOR_MARKDOWN},
		"post_memo":         {"Updated memo"},
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

	updated, err := registry.GetBlogStore().PostFindByID(context.Background(), post.GetID())
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if updated == nil {
		t.Fatal("expected updated to not be nil")
	}
	if updated.GetStatus() != blogstore.POST_STATUS_PUBLISHED {
		t.Errorf("expected status %s, got %s", blogstore.POST_STATUS_PUBLISHED, updated.GetStatus())
	}
	if updated.GetImageUrl() != "https://example.com/updated.jpg" {
		t.Errorf("expected image URL 'https://example.com/updated.jpg', got %s", updated.GetImageUrl())
	}
	if updated.GetFeatured() != "yes" {
		t.Errorf("expected featured 'yes', got %s", updated.GetFeatured())
	}
	if updated.GetEditor() != blogstore.POST_EDITOR_MARKDOWN {
		t.Errorf("expected editor %s, got %s", blogstore.POST_EDITOR_MARKDOWN, updated.GetEditor())
	}
	if updated.GetMemo() != "Updated memo" {
		t.Errorf("expected memo 'Updated memo', got %s", updated.GetMemo())
	}
	if updated.GetPublishedAtCarbon() == nil {
		t.Error("expected non-nil PublishedAtCarbon")
	}
}

func TestPostDetailsComponent_HandleRegenerateImage_BlogStoreNotAvailable(t *testing.T) {
	// Ensure we hit the early error branch without requiring AI wiring.
	registry := testutils.Setup(testutils.WithCacheStore(true))

	c := &postDetailsComponent{registry: registry, PostID: "some-id"}

	err := c.Handle(context.Background(), "regenerate_image", nil)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.FormErrorMessage != "Blog store not available" {
		t.Errorf("expected FormErrorMessage 'Blog store not available', got %s", c.FormErrorMessage)
	}
	if c.FormSuccessMessage != "" {
		t.Errorf("expected empty FormSuccessMessage, got %s", c.FormSuccessMessage)
	}
}
