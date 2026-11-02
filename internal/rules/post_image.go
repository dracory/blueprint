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

// PostImageURLs resolves image URLs for a batch of posts with a single
// MediaList query — calling PostImageURL per post is an N+1 (one media
// query per card) that added a database round-trip for every card on the
// blog index. Returns postID -> image URL, same precedence as
// PostImageURL: first media by sequence, else the post's ImageUrl.
func PostImageURLs(ctx context.Context, store blogstore.StoreInterface, posts []blogstore.PostInterface) map[string]string {
	out := map[string]string{}
	if store == nil || len(posts) == 0 {
		return out
	}

	postIDs := make([]string, 0, len(posts))
	for _, post := range posts {
		if post != nil {
			postIDs = append(postIDs, post.GetID())
		}
	}

	media, err := store.MediaList(ctx, blogstore.MediaQueryOptions{
		EntityIDIn: postIDs,
		OrderBy:    "sequence",
		SortOrder:  "ASC",
	})

	// First media per entity wins — ascending sequence order.
	mediaByEntity := map[string]blogstore.MediaInterface{}
	if err == nil {
		for _, m := range media {
			if _, exists := mediaByEntity[m.GetEntityID()]; !exists {
				mediaByEntity[m.GetEntityID()] = m
			}
		}
	}

	for _, post := range posts {
		if post == nil {
			continue
		}
		if m := mediaByEntity[post.GetID()]; m != nil {
			out[post.GetID()] = m.GetURL()
			continue
		}
		out[post.GetID()] = post.GetImageUrl()
	}
	return out
}
