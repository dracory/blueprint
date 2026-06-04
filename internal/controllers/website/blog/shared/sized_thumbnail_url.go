package shared

import (
	"project/internal/links"
	"project/internal/app"

	"github.com/dracory/blogstore"
)

func SizedThumbnailURL(app app.AppInterface, post blogstore.PostInterface, width, height, quality string) string {
	postImageURL := post.GetImageUrlOrDefault()
	extension := ImageExtension(postImageURL)

	if app == nil {
		return links.Website().Thumbnail(extension, width, height, quality, postImageURL)
	}
	return links.Website().Thumbnail(extension, width, height, quality, postImageURL, app.GetFileCache())
}
