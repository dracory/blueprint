package search

import (
	"context"
	"log"

	"github.com/dracory/blogstore"
	"github.com/dracory/cmsstore"
)

// ExampleUsage demonstrates how to use the Search block
func ExampleUsage(cmsStore cmsstore.StoreInterface, blogStore blogstore.StoreInterface) {
	// Create the Search block type
	searchBlock := NewSearchBlockType(cmsStore, blogStore)

	// Register the block type globally as a custom block
	cmsstore.RegisterCustomBlockType(searchBlock)

	log.Println("Search block type registered successfully!")
	log.Println("Now available in CMS admin as 'Search' block type")
}

// ExampleRendering demonstrates how to render the search block
func ExampleRendering(cmsStore cmsstore.StoreInterface, blogStore blogstore.StoreInterface, blockID string) {
	// Get the block
	block, err := cmsStore.BlockFindByID(context.Background(), blockID)
	if err != nil {
		log.Printf("Failed to find block: %v", err)
		return
	}

	// Get the block type
	blockType := cmsstore.GetBlockType(block.Type())
	if blockType == nil {
		log.Printf("Block type not found: %s", block.Type())
		return
	}

	// Render the search block
	html, err := blockType.Render(context.Background(), block)
	if err != nil {
		log.Printf("Failed to render search block: %v", err)
		return
	}

	log.Printf("Rendered search block HTML:\n%s", html)
}

// ExampleConfiguration demonstrates how to configure the search block
func ExampleConfiguration(block cmsstore.BlockInterface) {
	// Set placeholder text for the search input
	block.SetMeta("placeholder", "Search articles, pages...")

	// Set number of results per page
	block.SetMeta("results_per_page", "10")

	// Enable/disable content types in search
	block.SetMeta("show_pages", "true") // Include CMS pages
	block.SetMeta("show_posts", "true") // Include blog posts
}
