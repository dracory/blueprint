package shared

import (
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/blogstore"
)

func SizedThumbnailURL(registry registry.RegistryInterface, post blogstore.Post, width, height, quality string) string {
	postImageURL := post.ImageUrlOrDefault()
	extension := ImageExtension(postImageURL)

	if registry == nil {
		return links.Website().Thumbnail(extension, width, height, quality, postImageURL)
	}
	return links.Website().Thumbnail(extension, width, height, quality, postImageURL, registry.GetFileCache())
}
