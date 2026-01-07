package shared

import (
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/blogstore"
)

func SizedThumbnailURL(app registry.RegistryInterface, post blogstore.Post, width, height, quality string) string {
	postImageURL := post.ImageUrlOrDefault()
	extension := ImageExtension(postImageURL)

	if app == nil {
		return links.Website().Thumbnail(extension, width, height, quality, postImageURL)
	}
	return links.Website().Thumbnail(extension, width, height, quality, postImageURL, app.GetFileCache())
}
