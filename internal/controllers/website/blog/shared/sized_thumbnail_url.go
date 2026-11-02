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
	if app != nil && app.IsEnabledBlogStore() {
		postImageURL = rules.PostImageURL(ctx, app.GetBlogStore(), post)
	} else {
		postImageURL = post.GetImageUrlOrDefault()
	}

	return SizedThumbnailURLForImage(ctx, app, postImageURL, width, height, quality)
}

// SizedThumbnailURLForImage is SizedThumbnailURL with a pre-resolved image
// URL — used by list pages that batch-resolve post images in one query
// instead of hitting the media table per card.
func SizedThumbnailURLForImage(ctx context.Context, app app.AppInterface, postImageURL, width, height, quality string) string {
	extension := ImageExtension(postImageURL)

	if app == nil {
		return links.Website().Thumbnail(extension, width, height, quality, postImageURL)
	}
	return links.Website().Thumbnail(extension, width, height, quality, postImageURL, app.GetFileCache())
}
