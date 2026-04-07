package post_update

import (
	"context"
	"net/url"
	"testing"

	"project/internal/registry"
	"project/internal/testutils"

	"github.com/dracory/blogstore"
	"github.com/stretchr/testify/assert"
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

	assert.NoError(t, err)
	assert.Equal(t, "Post ID is required", c.FormErrorMessage)
}

func TestPostDetailsComponent_MountLoadsPostFields(t *testing.T) {
	registry, post := setupDetailsTestAppAndPost(t)

	c := &postDetailsComponent{registry: registry}

	err := c.Mount(context.Background(), map[string]string{
		"post_id": post.GetID(),
	})

	assert.NoError(t, err)
	assert.Equal(t, post.GetID(), c.PostID)
	assert.Equal(t, post.GetStatus(), c.FormStatus)
	assert.Equal(t, post.GetImageUrl(), c.FormImageUrl)
	assert.Equal(t, post.GetFeatured(), c.FormFeatured)
	assert.Equal(t, post.GetEditor(), c.FormEditor)
	assert.Equal(t, post.GetMemo(), c.FormMemo)
	assert.Equal(t, "", c.FormErrorMessage)
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

	assert.NoError(t, err)
	assert.Equal(t, "Status is required", c.FormErrorMessage)
	assert.Equal(t, "", c.FormSuccessMessage)
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

	assert.NoError(t, err)
	assert.Equal(t, "", c.FormErrorMessage)
	assert.Equal(t, "Post saved successfully", c.FormSuccessMessage)

	updated, err := registry.GetBlogStore().PostFindByID(context.Background(), post.GetID())
	assert.NoError(t, err)
	if assert.NotNil(t, updated) {
		assert.Equal(t, blogstore.POST_STATUS_PUBLISHED, updated.GetStatus())
		assert.Equal(t, "https://example.com/updated.jpg", updated.GetImageUrl())
		assert.Equal(t, "yes", updated.GetFeatured())
		assert.Equal(t, blogstore.POST_EDITOR_MARKDOWN, updated.GetEditor())
		assert.Equal(t, "Updated memo", updated.GetMemo())
		// PublishedAt is normalized inside the component; just ensure it's non-empty
		assert.NotEmpty(t, updated.GetPublishedAtCarbon())
	}
}

func TestPostDetailsComponent_HandleRegenerateImage_BlogStoreNotAvailable(t *testing.T) {
	// Ensure we hit the early error branch without requiring AI wiring.
	registry := testutils.Setup(testutils.WithCacheStore(true))

	c := &postDetailsComponent{registry: registry, PostID: "some-id"}

	err := c.Handle(context.Background(), "regenerate_image", nil)

	assert.NoError(t, err)
	assert.Equal(t, "Blog store not available", c.FormErrorMessage)
	assert.Equal(t, "", c.FormSuccessMessage)
}
