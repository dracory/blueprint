package post_update

import (
	"context"
	"net/url"
	"testing"

	"project/internal/testutils"
	"project/internal/types"

	"github.com/dracory/blogstore"
	"github.com/stretchr/testify/assert"
)

func setupDetailsTestAppAndPost(t *testing.T) (types.AppInterface, *blogstore.Post) {
	t.Helper()

	app := testutils.Setup(
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

	if err := app.GetBlogStore().PostCreate(post); err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	return app, post
}

func TestPostDetailsComponent_MountRequiresPostID(t *testing.T) {
	app, _ := setupDetailsTestAppAndPost(t)

	c := &postDetailsComponent{App: app}

	err := c.Mount(context.Background(), map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, "Post ID is required", c.FormErrorMessage)
}

func TestPostDetailsComponent_MountLoadsPostFields(t *testing.T) {
	app, post := setupDetailsTestAppAndPost(t)

	c := &postDetailsComponent{App: app}

	err := c.Mount(context.Background(), map[string]string{
		"post_id": post.ID(),
	})

	assert.NoError(t, err)
	assert.Equal(t, post.ID(), c.PostID)
	assert.Equal(t, post.Status(), c.FormStatus)
	assert.Equal(t, post.ImageUrl(), c.FormImageUrl)
	assert.Equal(t, post.Featured(), c.FormFeatured)
	assert.Equal(t, post.Editor(), c.FormEditor)
	assert.Equal(t, post.Memo(), c.FormMemo)
	assert.Equal(t, "", c.FormErrorMessage)
}

func TestPostDetailsComponent_HandleSave_ValidatesStatusRequired(t *testing.T) {
	app, post := setupDetailsTestAppAndPost(t)

	c := &postDetailsComponent{App: app, PostID: post.ID()}

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
	app, post := setupDetailsTestAppAndPost(t)

	c := &postDetailsComponent{App: app, PostID: post.ID()}

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

	updated, err := app.GetBlogStore().PostFindByID(post.ID())
	assert.NoError(t, err)
	if assert.NotNil(t, updated) {
		assert.Equal(t, blogstore.POST_STATUS_PUBLISHED, updated.Status())
		assert.Equal(t, "https://example.com/updated.jpg", updated.ImageUrl())
		assert.Equal(t, "yes", updated.Featured())
		assert.Equal(t, blogstore.POST_EDITOR_MARKDOWN, updated.Editor())
		assert.Equal(t, "Updated memo", updated.Memo())
		// PublishedAt is normalized inside the component; just ensure it's non-empty
		assert.NotEmpty(t, updated.PublishedAtCarbon())
	}
}

func TestPostDetailsComponent_HandleRegenerateImage_BlogStoreNotAvailable(t *testing.T) {
	// Ensure we hit the early error branch without requiring AI wiring.
	app := testutils.Setup(testutils.WithCacheStore(true))

	c := &postDetailsComponent{App: app, PostID: "some-id"}

	err := c.Handle(context.Background(), "regenerate_image", nil)

	assert.NoError(t, err)
	assert.Equal(t, "Blog store not available", c.FormErrorMessage)
	assert.Equal(t, "", c.FormSuccessMessage)
}
