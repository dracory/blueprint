package cmsblocks

import (
	"project/internal/cmsblocks/blogpost"
	"project/internal/cmsblocks/blogpostlist"
	"project/internal/cmsblocks/search"
	"project/internal/registry"

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
//   - registry: Registry interface for accessing services
//
// Returns:
//   - None
func CmsAddBlockTypes(registry registry.RegistryInterface) {
	if !registry.GetConfig().GetCmsStoreUsed() {
		return
	}

	if registry.GetCmsStore() == nil {
		return
	}

	// Register the Blog Post List block type as a custom block
	if registry.GetBlogStore() != nil {
		blogPostListBlock := blogpostlist.NewBlogPostListBlockType(registry.GetBlogStore())
		cmsstore.RegisterCustomBlockType(blogPostListBlock)

		// Register the Blog Post block type as a custom block
		blogPostBlock := blogpost.NewBlogPostBlockType(registry.GetBlogStore())
		cmsstore.RegisterCustomBlockType(blogPostBlock)

		// Register the Search block type as a custom block
		searchBlock := search.NewSearchBlockType(registry.GetCmsStore(), registry.GetBlogStore())
		cmsstore.RegisterCustomBlockType(searchBlock)
	}
}
