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

func setupSEOTestAppAndPost(t *testing.T) (types.AppInterface, *blogstore.Post) {
	t.Helper()

	app := testutils.Setup(
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

	if err := app.GetBlogStore().PostCreate(post); err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	return app, post
}

func TestPostSEOComponent_MountRequiresPostID(t *testing.T) {
	app, _ := setupSEOTestAppAndPost(t)

	c := &postSEOComponent{App: app}

	err := c.Mount(context.Background(), map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, "Post ID is required", c.FormErrorMessage)
}

func TestPostSEOComponent_MountLoadsPostFields(t *testing.T) {
	app, post := setupSEOTestAppAndPost(t)

	c := &postSEOComponent{App: app}

	err := c.Mount(context.Background(), map[string]string{
		"post_id": post.ID(),
	})

	assert.NoError(t, err)
	assert.Equal(t, post.ID(), c.PostID)
	assert.Equal(t, post.CanonicalURL(), c.FormCanonicalURL)
	assert.Equal(t, post.MetaDescription(), c.FormMetaDescription)
	assert.Equal(t, post.MetaKeywords(), c.FormMetaKeywords)
	assert.Equal(t, post.MetaRobots(), c.FormMetaRobots)
	assert.Equal(t, "", c.FormErrorMessage)
}

func TestPostSEOComponent_HandleSave_UpdatesPostAndSetsSuccess(t *testing.T) {
	app, post := setupSEOTestAppAndPost(t)

	c := &postSEOComponent{App: app, PostID: post.ID()}

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
	updated, err := app.GetBlogStore().PostFindByID(post.ID())
	assert.NoError(t, err)
	if assert.NotNil(t, updated) {
		assert.Equal(t, "https://example.com/updated", updated.CanonicalURL())
		assert.Equal(t, "Updated meta", updated.MetaDescription())
		assert.Equal(t, "one,two", updated.MetaKeywords())
		assert.Equal(t, "NOINDEX, FOLLOW", updated.MetaRobots())
	}
}
