package rules

import (
	"context"
	"sort"

	"github.com/dracory/blogstore"
)

// PostImageURL resolves the image URL for a post.
// It returns the first media URL (sorted by sequence) if media exists,
// otherwise falls back to the post's ImageUrl field.
func PostImageURL(ctx context.Context, store blogstore.StoreInterface, post blogstore.PostInterface) string {
	if store == nil || post == nil {
		return ""
	}

	media, err := store.MediaListByEntityID(ctx, post.GetID())
	if err == nil && len(media) > 0 {
		sort.Slice(media, func(i, j int) bool {
			return media[i].GetSequence() < media[j].GetSequence()
		})
		return media[0].GetURL()
	}

	return post.GetImageUrl()
}
