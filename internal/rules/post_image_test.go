package rules

import (
	"context"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/blogstore"
)

func createPost(t *testing.T, ctx context.Context, store blogstore.StoreInterface, imageURL string) blogstore.PostInterface {
	t.Helper()
	post := blogstore.NewPost().
		SetTitle("Test Post").
		SetStatus(blogstore.POST_STATUS_PUBLISHED).
		SetImageUrl(imageURL)
	if err := store.PostCreate(ctx, post); err != nil {
		t.Fatalf("PostCreate failed: %v", err)
	}
	return post
}

func createMedia(t *testing.T, ctx context.Context, store blogstore.StoreInterface, entityID, url string, sequence int) {
	t.Helper()
	media := blogstore.NewMedia().
		SetEntityID(entityID).
		SetTitle("test.jpg").
		SetURL(url).
		SetType("image/jpeg").
		SetSize("1024").
		SetExtension("jpg").
		SetSequence(sequence)
	if err := store.MediaCreate(ctx, media); err != nil {
		t.Fatalf("MediaCreate failed: %v", err)
	}
}

func TestPostImageURLNilArgs(t *testing.T) {
	if got := PostImageURL(context.Background(), nil, nil); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestPostImageURLFallsBackToImageUrl(t *testing.T) {
	app := testutils.Setup(testutils.WithBlogStore(true))
	defer app.Close()

	ctx := context.Background()
	store := app.GetBlogStore()

	post := createPost(t, ctx, store, "https://example.com/fallback.jpg")

	if got := PostImageURL(ctx, store, post); got != "https://example.com/fallback.jpg" {
		t.Errorf("PostImageURL = %q, want fallback image", got)
	}
}

func TestPostImageURLPrefersFirstMediaBySequence(t *testing.T) {
	app := testutils.Setup(testutils.WithBlogStore(true))
	defer app.Close()

	ctx := context.Background()
	store := app.GetBlogStore()

	post := createPost(t, ctx, store, "https://example.com/fallback.jpg")
	createMedia(t, ctx, store, post.GetID(), "https://example.com/second.jpg", 2)
	createMedia(t, ctx, store, post.GetID(), "https://example.com/first.jpg", 1)

	if got := PostImageURL(ctx, store, post); got != "https://example.com/first.jpg" {
		t.Errorf("PostImageURL = %q, want first media by sequence", got)
	}
}

func TestPostImageURLsEmpty(t *testing.T) {
	app := testutils.Setup(testutils.WithBlogStore(true))
	defer app.Close()

	out := PostImageURLs(context.Background(), app.GetBlogStore(), nil)
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}

func TestPostImageURLsBatch(t *testing.T) {
	app := testutils.Setup(testutils.WithBlogStore(true))
	defer app.Close()

	ctx := context.Background()
	store := app.GetBlogStore()

	postWithMedia := createPost(t, ctx, store, "https://example.com/fallback1.jpg")
	postWithoutMedia := createPost(t, ctx, store, "https://example.com/fallback2.jpg")

	createMedia(t, ctx, store, postWithMedia.GetID(), "https://example.com/media-second.jpg", 2)
	createMedia(t, ctx, store, postWithMedia.GetID(), "https://example.com/media-first.jpg", 1)

	out := PostImageURLs(ctx, store, []blogstore.PostInterface{postWithMedia, postWithoutMedia})

	if got := out[postWithMedia.GetID()]; got != "https://example.com/media-first.jpg" {
		t.Errorf("post with media = %q, want media-first.jpg", got)
	}
	if got := out[postWithoutMedia.GetID()]; got != "https://example.com/fallback2.jpg" {
		t.Errorf("post without media = %q, want fallback2.jpg", got)
	}
}
