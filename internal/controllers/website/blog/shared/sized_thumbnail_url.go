package shared

import (
	"context"
	"project/internal/app"
	"project/internal/links"
	"project/internal/rules"

	"github.com/dracory/blogstore"
)

func SizedThumbnailURL(ctx context.Context, app app.AppInterface, post blogstore.PostInterface, width, height, quality string) string {
	var postImageURL string
	if app != nil && app.GetBlogStore() != nil {
		postImageURL = rules.PostImageURL(ctx, app.GetBlogStore(), post)
	} else {
		postImageURL = post.GetImageUrlOrDefault()
	}

	extension := ImageExtension(postImageURL)

	if app == nil {
		return links.Website().Thumbnail(extension, width, height, quality, postImageURL)
	}
	return links.Website().Thumbnail(extension, width, height, quality, postImageURL, app.GetFileCache())
}
