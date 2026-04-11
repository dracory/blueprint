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

	assert.NoError(t, err)
	assert.Equal(t, "Post ID is required", c.FormErrorMessage)
}

func TestPostSEOComponent_MountLoadsPostFields(t *testing.T) {
	registry, post := setupSEOTestAppAndPost(t)

	c := &postSEOComponent{registry: registry}

	err := c.Mount(context.Background(), map[string]string{
		"post_id": post.GetID(),
	})

	assert.NoError(t, err)
	assert.Equal(t, post.GetID(), c.PostID)
	assert.Equal(t, post.GetCanonicalURL(), c.FormCanonicalURL)
	assert.Equal(t, post.GetMetaDescription(), c.FormMetaDescription)
	assert.Equal(t, post.GetMetaKeywords(), c.FormMetaKeywords)
	assert.Equal(t, post.GetMetaRobots(), c.FormMetaRobots)
	assert.Equal(t, "", c.FormErrorMessage)
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

	assert.NoError(t, err)
	assert.Equal(t, "", c.FormErrorMessage)
	assert.Equal(t, "Post saved successfully", c.FormSuccessMessage)

	// Reload post from store and verify fields were updated
	updated, err := registry.GetBlogStore().PostFindByID(context.Background(), post.GetID())
	assert.NoError(t, err)
	if assert.NotNil(t, updated) {
		assert.Equal(t, "https://example.com/updated", updated.GetCanonicalURL())
		assert.Equal(t, "Updated meta", updated.GetMetaDescription())
		assert.Equal(t, "one,two", updated.GetMetaKeywords())
		assert.Equal(t, "NOINDEX, FOLLOW", updated.GetMetaRobots())
	}
}
