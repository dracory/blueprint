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

func setupContentTestAppAndPost(t *testing.T) (types.AppInterface, *blogstore.Post) {
	t.Helper()

	app := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithBlogStore(true),
	)

	post := blogstore.NewPost()
	post.SetTitle("Content Test Post")
	post.SetSummary("Original summary")
	post.SetContent("Original content")
	post.SetStatus(blogstore.POST_STATUS_DRAFT)

	if err := app.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	return app, post
}

func TestPostContentComponent_MountRequiresPostID(t *testing.T) {
	app, _ := setupContentTestAppAndPost(t)

	c := &postContentComponent{App: app}

	err := c.Mount(context.Background(), map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, "Post ID is required", c.FormErrorMessage)
}

func TestPostContentComponent_MountLoadsPostFields(t *testing.T) {
	app, post := setupContentTestAppAndPost(t)

	c := &postContentComponent{App: app}

	err := c.Mount(context.Background(), map[string]string{
		"post_id": post.ID(),
	})

	assert.NoError(t, err)
	assert.Equal(t, post.ID(), c.PostID)
	assert.Equal(t, post.Title(), c.FormTitle)
	assert.Equal(t, post.Summary(), c.FormSummary)
	assert.Equal(t, post.Content(), c.FormContent)
	assert.Equal(t, "", c.FormErrorMessage)
}

func TestPostContentComponent_HandleSave_ValidatesTitleRequired(t *testing.T) {
	app, post := setupContentTestAppAndPost(t)

	c := &postContentComponent{App: app, PostID: post.ID()}

	values := url.Values{
		"post_title":   {""},
		"post_summary": {"New summary"},
		"post_content": {"New content"},
	}

	err := c.Handle(context.Background(), "save", values)

	assert.NoError(t, err)
	assert.Equal(t, "Title is required", c.FormErrorMessage)
	assert.Equal(t, "", c.FormSuccessMessage)
}

func TestPostContentComponent_HandleSave_UpdatesPostAndSetsSuccess(t *testing.T) {
	app, post := setupContentTestAppAndPost(t)

	c := &postContentComponent{App: app, PostID: post.ID()}

	values := url.Values{
		"post_title":   {"Updated Title"},
		"post_summary": {"Updated summary"},
		"post_content": {"Updated content"},
	}

	err := c.Handle(context.Background(), "save", values)

	assert.NoError(t, err)
	assert.Equal(t, "", c.FormErrorMessage)
	assert.Equal(t, "Post saved successfully", c.FormSuccessMessage)

	updated, err := app.GetBlogStore().PostFindByID(context.Background(), post.ID())
	assert.NoError(t, err)
	if assert.NotNil(t, updated) {
		assert.Equal(t, "Updated Title", updated.Title())
		assert.Equal(t, "Updated summary", updated.Summary())
		assert.Equal(t, "Updated content", updated.Content())
	}
}
