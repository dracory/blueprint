package post_update

import (
	"context"
	"net/url"
	"testing"

	"project/internal/registry"
	"project/internal/testutils"

	"github.com/dracory/blogstore"
)

func setupContentTestAppAndPost(t *testing.T) (registry.RegistryInterface, blogstore.PostInterface) {
	t.Helper()

	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithBlogStore(true),
	)

	post := blogstore.NewPost()
	post.SetTitle("Content Test Post")
	post.SetSummary("Original summary")
	post.SetContent("Original content")
	post.SetStatus(blogstore.POST_STATUS_DRAFT)

	if err := registry.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	return registry, post
}

func TestPostContentComponent_MountRequiresPostID(t *testing.T) {
	registry, _ := setupContentTestAppAndPost(t)

	c := &postContentComponent{registry: registry}

	err := c.Mount(context.Background(), map[string]string{})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.FormErrorMessage != "Post ID is required" {
		t.Errorf("expected FormErrorMessage 'Post ID is required', got %s", c.FormErrorMessage)
	}
}

func TestPostContentComponent_MountLoadsPostFields(t *testing.T) {
	registry, post := setupContentTestAppAndPost(t)

	c := &postContentComponent{registry: registry}

	err := c.Mount(context.Background(), map[string]string{
		"post_id": post.GetID(),
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.PostID != post.GetID() {
		t.Errorf("expected PostID %s, got %s", post.GetID(), c.PostID)
	}
	if c.FormTitle != post.GetTitle() {
		t.Errorf("expected FormTitle %s, got %s", post.GetTitle(), c.FormTitle)
	}
	if c.FormSummary != post.GetSummary() {
		t.Errorf("expected FormSummary %s, got %s", post.GetSummary(), c.FormSummary)
	}
	if c.FormContent != post.GetContent() {
		t.Errorf("expected FormContent %s, got %s", post.GetContent(), c.FormContent)
	}
	if c.FormErrorMessage != "" {
		t.Errorf("expected empty FormErrorMessage, got %s", c.FormErrorMessage)
	}
}

func TestPostContentComponent_HandleSave_ValidatesTitleRequired(t *testing.T) {
	registry, post := setupContentTestAppAndPost(t)

	c := &postContentComponent{registry: registry, PostID: post.GetID()}

	values := url.Values{
		"post_title":   {""},
		"post_summary": {"New summary"},
		"post_content": {"New content"},
	}

	err := c.Handle(context.Background(), "save", values)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if c.FormErrorMessage != "Title is required" {
		t.Errorf("expected FormErrorMessage 'Title is required', got %s", c.FormErrorMessage)
	}
	if c.FormSuccessMessage != "" {
		t.Errorf("expected empty FormSuccessMessage, got %s", c.FormSuccessMessage)
	}
}

func TestPostContentComponent_HandleSave_UpdatesPostAndSetsSuccess(t *testing.T) {
	registry, post := setupContentTestAppAndPost(t)

	c := &postContentComponent{registry: registry, PostID: post.GetID()}

	values := url.Values{
		"post_title":   {"Updated Title"},
		"post_summary": {"Updated summary"},
		"post_content": {"Updated content"},
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
	if updated.GetTitle() != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %s", updated.GetTitle())
	}
	if updated.GetSummary() != "Updated summary" {
		t.Errorf("expected summary 'Updated summary', got %s", updated.GetSummary())
	}
	if updated.GetContent() != "Updated content" {
		t.Errorf("expected content 'Updated content', got %s", updated.GetContent())
	}
}
