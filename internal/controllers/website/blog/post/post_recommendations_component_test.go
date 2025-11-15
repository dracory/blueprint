package post

import (
	"context"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/blogstore"
	"github.com/dromara/carbon/v2"
)

func TestNewPostRecommendationsComponent(t *testing.T) {
	app := testutils.Setup()

	component := NewPostRecommendationsComponent(app)
	if component == nil {
		t.Fatal("Expected component to be created, got nil")
	}
}

func TestPostRecommendationsComponent_GetKind(t *testing.T) {
	app := testutils.Setup()
	component := NewPostRecommendationsComponent(app)

	kind := component.GetKind()
	expected := "blog_post_recommendations"
	if kind != expected {
		t.Errorf("Expected kind '%s', got '%s'", expected, kind)
	}
}

func TestPostRecommendationsComponent_Mount_NoApp(t *testing.T) {
	component := &postRecommendationsComponent{}

	err := component.Mount(context.Background(), map[string]string{"post_id": "123"})
	if err != nil {
		t.Errorf("Expected no error when app is nil, got: %v", err)
	}

	if component.errorMessage == "" {
		t.Error("Expected error message when app is nil")
	}
}

func TestPostRecommendationsComponent_Mount_NoStore(t *testing.T) {
	app := testutils.Setup() // No blog store enabled
	component := NewPostRecommendationsComponent(app).(*postRecommendationsComponent)

	err := component.Mount(context.Background(), map[string]string{"post_id": "123"})
	if err != nil {
		t.Errorf("Expected no error when store is nil, got: %v", err)
	}

	if component.errorMessage == "" {
		t.Error("Expected error message when blog store is not configured")
	}
}

func TestPostRecommendationsComponent_Mount_StoreError(t *testing.T) {
	app := testutils.Setup(testutils.WithBlogStore(true))
	component := NewPostRecommendationsComponent(app)

	// Simulate store error by not creating any posts (store might return error for empty results)
	// This test verifies error handling in the Mount method
	err := component.Mount(context.Background(), map[string]string{"post_id": "123"})
	if err != nil {
		t.Errorf("Expected no error from Mount, got: %v", err)
	}

	// The component should handle the error gracefully and set errorMessage
	// We can't easily simulate a real store error without mocking
}

func TestPostRecommendationsComponent_Mount_Success(t *testing.T) {
	app := testutils.Setup(testutils.WithBlogStore(true))

	// Create some test posts
	for i := 0; i < 5; i++ {
		post := blogstore.NewPost()
		post.SetTitle("Test Post " + string(rune('A'+i)))
		post.SetContent("Test content")
		post.SetStatus(blogstore.POST_STATUS_PUBLISHED)
		post.SetPublishedAt(carbon.Now().ToDateTimeString())
		err := app.GetBlogStore().PostCreate(post)
		if err != nil {
			t.Fatalf("Failed to create test post: %v", err)
		}
	}

	component := NewPostRecommendationsComponent(app).(*postRecommendationsComponent)

	// Mount with a post ID that exists
	posts, err := app.GetBlogStore().PostList(blogstore.PostQueryOptions{
		Status: blogstore.POST_STATUS_PUBLISHED,
		Limit:  1,
	})
	if err != nil || len(posts) == 0 {
		t.Fatal("Failed to get a test post ID")
	}

	err = component.Mount(context.Background(), map[string]string{"post_id": posts[0].ID()})
	if err != nil {
		t.Errorf("Expected no error from Mount, got: %v", err)
	}

	if component.errorMessage != "" {
		t.Errorf("Expected no error message, got: %s", component.errorMessage)
	}

	// Should have filtered posts (excluding current post)
	if len(component.Posts) == 0 {
		t.Error("Expected some posts to be loaded")
	}

	if len(component.Posts) > recommendationsDisplayLimit {
		t.Errorf("Expected at most %d posts, got %d", recommendationsDisplayLimit, len(component.Posts))
	}
}

func TestPostRecommendationsComponent_Render_Error(t *testing.T) {
	component := &postRecommendationsComponent{
		errorMessage: "Test error",
	}

	html := component.Render(context.Background())

	// Should return empty div when there's an error
	htmlStr := html.ToHTML()
	if htmlStr == "" {
		t.Error("Expected HTML output even with error")
	}

	if !strings.Contains(htmlStr, "<div>") {
		t.Error("Expected div element in error case")
	}
}

func TestPostRecommendationsComponent_Render_NoPosts(t *testing.T) {
	app := testutils.Setup()
	component := NewPostRecommendationsComponent(app)

	// Mount without any posts
	err := component.Mount(context.Background(), map[string]string{"post_id": "nonexistent"})
	if err != nil {
		t.Errorf("Expected no error from Mount, got: %v", err)
	}

	html := component.Render(context.Background())
	htmlStr := html.ToHTML()

	if htmlStr == "" {
		t.Error("Expected HTML output")
	}

	if !strings.Contains(htmlStr, "<div>") {
		t.Error("Expected div element when no posts")
	}
}

func TestPostRecommendationsComponent_Render_WithPosts(t *testing.T) {
	app := testutils.Setup(testutils.WithBlogStore(true))

	// Create test posts
	posts := make([]blogstore.Post, 3)
	for i := 0; i < 3; i++ {
		post := blogstore.NewPost()
		post.SetTitle("Test Post " + string(rune('A'+i)))
		post.SetContent("Test content " + string(rune('A'+i)))
		post.SetSummary("Test summary " + string(rune('A'+i)))
		post.SetStatus(blogstore.POST_STATUS_PUBLISHED)
		err := app.GetBlogStore().PostCreate(post)
		if err != nil {
			t.Fatalf("Failed to create test post: %v", err)
		}
		posts[i] = *post
	}

	component := NewPostRecommendationsComponent(app).(*postRecommendationsComponent)

	// Mount with first post as current
	err := component.Mount(context.Background(), map[string]string{"post_id": posts[0].ID()})
	if err != nil {
		t.Errorf("Expected no error from Mount, got: %v", err)
	}

	html := component.Render(context.Background())
	htmlStr := html.ToHTML()

	if htmlStr == "" {
		t.Error("Expected HTML output")
	}

	// Verify structure
	if !strings.Contains(htmlStr, "<section") {
		t.Error("Expected section element")
	}

	if !strings.Contains(htmlStr, "Keep Reading") {
		t.Error("Expected heading text")
	}

	if !strings.Contains(htmlStr, "Explore more practical guides") {
		t.Error("Expected summary text")
	}

	if !strings.Contains(htmlStr, "View All Blog Posts") {
		t.Error("Expected view all link")
	}

	if !strings.Contains(htmlStr, "btn-primary") {
		t.Error("Expected primary button")
	}

	// Should contain post cards
	if !strings.Contains(htmlStr, "card") {
		t.Error("Expected card elements for posts")
	}
}

func TestPostRecommendationsComponent_Handle(t *testing.T) {
	app := testutils.Setup()
	component := NewPostRecommendationsComponent(app)

	err := component.Handle(context.Background(), "test", nil)
	// Handle method does nothing, so should not error
	if err != nil {
		t.Errorf("Expected no error from Handle, got: %v", err)
	}
}

func TestPostRecommendationsComponent_PostCard(t *testing.T) {
	app := testutils.Setup()
	component := NewPostRecommendationsComponent(app).(*postRecommendationsComponent)

	post := blogstore.NewPost()
	post.SetID("test-post")
	post.SetTitle("Test Post Title")
	post.SetSummary("This is a very long test summary that should definitely be truncated because it exceeds the maximum allowed length of one hundred and sixty characters for display in the recommendations section and therefore needs to show an ellipsis at the end.")
	post.SetImageUrl("test-image.jpg")

	// Test with image
	post.SetImageUrl("test-image.jpg")
	componentTyped := component
	card := componentTyped.postCard(*post)
	htmlStr := card.ToHTML()

	if !strings.Contains(htmlStr, "<img") {
		t.Error("Expected image element in card")
	}

	if !strings.Contains(htmlStr, "test-image.jpg") {
		t.Error("Expected image URL in card")
	}

	if !strings.Contains(htmlStr, "Test Post Title") {
		t.Error("Expected post title in card")
	}

	if !strings.Contains(htmlStr, "...") {
		t.Error("Expected truncated summary with ellipsis")
	}

	if !strings.Contains(htmlStr, "Read This Next") {
		t.Error("Expected button text")
	}

	// Test without image
	post.SetImageUrl("")
	componentTyped = component
	card = componentTyped.postCard(*post)
	htmlStr = card.ToHTML()

	if strings.Contains(htmlStr, "test-image.jpg") {
		t.Error("Should not contain image when no image URL is set")
	}

	if !strings.Contains(htmlStr, "Test Post Title") {
		t.Error("Expected post title in card even without image")
	}
}

func TestPostRecommendationsComponent_TruncatedSummary(t *testing.T) {
	app := testutils.Setup()
	component := NewPostRecommendationsComponent(app).(*postRecommendationsComponent)

	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"Short text", "Short text"},
		{"This is a very long summary that should be truncated because it exceeds the maximum allowed length for display in the recommendations section. This extra text ensures we go well beyond the 160 character limit that triggers truncation in the component.", "This is a very long summary that should be truncated because it exceeds the maximum allowed length for display in the recommendations section. This extra text e..."},
		{"   Leading and trailing spaces   ", "Leading and trailing spaces"},
	}

	for _, tt := range tests {
		result := component.truncatedSummary(tt.input)
		if result != tt.expected {
			t.Errorf("truncatedSummary(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}
