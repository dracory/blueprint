package shared

import (
	"project/internal/links"

	"github.com/dracory/blogstore"
)

func SizedThumbnailURL(post blogstore.Post, width, height, quality string) string {
	postImageURL := post.ImageUrlOrDefault()
	extension := ImageExtension(postImageURL)

	return links.Website().Thumbnail(extension, width, height, quality, postImageURL)
}
