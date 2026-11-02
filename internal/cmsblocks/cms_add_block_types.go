package cmsblocks

import (
	"project/internal/app"
	"project/internal/links"

	"github.com/dracory/cmsblocks/blogpost"
	"github.com/dracory/cmsblocks/blogpostlist"
	"github.com/dracory/cmsblocks/search"
	"github.com/dracory/cmsstore"
)

// CmsAddBlockTypes adds custom block types to the CMS store
//
// Business Logic:
//   - Check if the CMS store is used
//   - Check if the CMS store is nil
//   - Add custom block types to the CMS store
//   - Loaded in the cmd/server/main.go file
//
// The block type implementations live in the external
// github.com/dracory/cmsblocks package; this file only injects the
// project's stores and URL builder.
//
// Parameters:
//   - app: app interface for accessing services
//
// Returns:
//   - None
func CmsAddBlockTypes(app app.AppInterface) {
	if !app.GetConfig().GetCmsStoreUsed() {
		return
	}

	if app.IsDisabledCmsStore() {
		return
	}

	postURL := func(postID, postSlug string) string {
		return links.Website().BlogPost(postID, postSlug)
	}

	// Register the Blog Post List block type as a custom block
	if app.IsEnabledBlogStore() {
		blogPostListBlock := blogpostlist.NewBlogPostListBlockType(app.GetBlogStore())
		blogPostListBlock.SetPostURLFunc(postURL)
		cmsstore.RegisterCustomBlockType(blogPostListBlock)

		// Register the Blog Post block type as a custom block
		blogPostBlock := blogpost.NewBlogPostBlockType(app.GetBlogStore())
		blogPostBlock.SetPostURLFunc(postURL)
		cmsstore.RegisterCustomBlockType(blogPostBlock)

		// Register the Search block type as a custom block
		searchBlock := search.NewSearchBlockType(app.GetCmsStore(), app.GetBlogStore())
		searchBlock.SetPostURLFunc(postURL)
		cmsstore.RegisterCustomBlockType(searchBlock)
	}
}
