package cmsblocks

import (
	"project/internal/cmsblocks/blogpost"
	"project/internal/cmsblocks/blogpostlist"
	"project/internal/cmsblocks/search"
	"project/internal/app"

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
// Parameters:
//   - app: app interface for accessing services
//
// Returns:
//   - None
func CmsAddBlockTypes(app app.AppInterface) {
	if !app.GetConfig().GetCmsStoreUsed() {
		return
	}

	if app.GetCmsStore() == nil {
		return
	}

	// Register the Blog Post List block type as a custom block
	if app.GetBlogStore() != nil {
		blogPostListBlock := blogpostlist.NewBlogPostListBlockType(app.GetBlogStore())
		cmsstore.RegisterCustomBlockType(blogPostListBlock)

		// Register the Blog Post block type as a custom block
		blogPostBlock := blogpost.NewBlogPostBlockType(app.GetBlogStore())
		cmsstore.RegisterCustomBlockType(blogPostBlock)

		// Register the Search block type as a custom block
		searchBlock := search.NewSearchBlockType(app.GetCmsStore(), app.GetBlogStore())
		cmsstore.RegisterCustomBlockType(searchBlock)
	}
}
