package cmsblocks

import (
	"testing"

	"project/internal/testutils"

	"github.com/dracory/cmsstore"
)

// The block type registry is global, so tests asserting absence must run
// before the test that registers the types.

func TestCmsAddBlockTypesCmsStoreDisabled(t *testing.T) {
	app := testutils.Setup(testutils.WithBlogStore(true))
	defer app.Close()

	CmsAddBlockTypes(app)

	for _, key := range []string{"blog_post", "blog_post_list", "search"} {
		if cmsstore.GetBlockType(key) != nil {
			t.Errorf("expected block type %q to not be registered", key)
		}
	}
}

func TestCmsAddBlockTypesBlogStoreDisabled(t *testing.T) {
	app := testutils.Setup(testutils.WithCmsStore(true))
	defer app.Close()

	CmsAddBlockTypes(app)

	for _, key := range []string{"blog_post", "blog_post_list", "search"} {
		if cmsstore.GetBlockType(key) != nil {
			t.Errorf("expected block type %q to not be registered", key)
		}
	}
}

func TestCmsAddBlockTypesRegistersBlockTypes(t *testing.T) {
	app := testutils.Setup(
		testutils.WithCmsStore(true),
		testutils.WithBlogStore(true),
	)
	defer app.Close()

	CmsAddBlockTypes(app)

	for _, key := range []string{"blog_post", "blog_post_list", "search"} {
		if cmsstore.GetBlockType(key) == nil {
			t.Errorf("expected block type %q to be registered", key)
		}
	}
}
