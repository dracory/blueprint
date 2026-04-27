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

	assert.NoError(t, err)
	assert.Equal(t, "Post ID is required", c.FormErrorMessage)
}

func TestPostContentComponent_MountLoadsPostFields(t *testing.T) {
	registry, post := setupContentTestAppAndPost(t)

	c := &postContentComponent{registry: registry}

	err := c.Mount(context.Background(), map[string]string{
		"post_id": post.GetID(),
	})

	assert.NoError(t, err)
	assert.Equal(t, post.GetID(), c.PostID)
	assert.Equal(t, post.GetTitle(), c.FormTitle)
	assert.Equal(t, post.GetSummary(), c.FormSummary)
	assert.Equal(t, post.GetContent(), c.FormContent)
	assert.Equal(t, "", c.FormErrorMessage)
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

	assert.NoError(t, err)
	assert.Equal(t, "Title is required", c.FormErrorMessage)
	assert.Equal(t, "", c.FormSuccessMessage)
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

	assert.NoError(t, err)
	assert.Equal(t, "", c.FormErrorMessage)
	assert.Equal(t, "Post saved successfully", c.FormSuccessMessage)

	updated, err := registry.GetBlogStore().PostFindByID(context.Background(), post.GetID())
	assert.NoError(t, err)
	if assert.NotNil(t, updated) {
		assert.Equal(t, "Updated Title", updated.GetTitle())
		assert.Equal(t, "Updated summary", updated.GetSummary())
		assert.Equal(t, "Updated content", updated.GetContent())
	}
}
