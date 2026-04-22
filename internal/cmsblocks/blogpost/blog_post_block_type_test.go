package blogpost

import (
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/blogstore"
	"github.com/dracory/cmsstore"
)

// TestBlogPostBlockType_BasicProperties tests basic properties
func TestBlogPostBlockType_BasicProperties(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blockType := NewBlogPostBlockType(registry.GetBlogStore())

	if blockType.TypeKey() != "blog_post" {
		t.Errorf("Expected type key 'blog_post', got '%s'", blockType.TypeKey())
	}

	if blockType.TypeLabel() != "Blog Post" {
		t.Errorf("Expected type label 'Blog Post', got '%s'", blockType.TypeLabel())
	}
}

// TestBlogPostBlockType_GetPreview tests preview
func TestBlogPostBlockType_GetPreview(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blogStore := registry.GetBlogStore()
	blockType := NewBlogPostBlockType(blogStore)

	// Create a real block for testing
	block := cmsstore.NewBlock()
	block.SetType("blog_post")
	block.SetName("Test Blog Post Block")

	// Test with no post ID
	preview := blockType.GetPreview(block)
	if preview != "Blog Post (auto-detect from URL)" {
		t.Errorf("Expected preview 'Blog Post (auto-detect from URL)', got '%s'", preview)
	}

	// Test with post ID
	block.SetMeta("post_id", "post-123")
	preview = blockType.GetPreview(block)
	if !strings.Contains(preview, "post-123") {
		t.Errorf("Expected preview to contain post ID, got '%s'", preview)
	}
}

// TestBlogPostBlockType_Validate tests validation
func TestBlogPostBlockType_Validate(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blogStore := registry.GetBlogStore()
	blockType := NewBlogPostBlockType(blogStore)

	block := cmsstore.NewBlock()
	block.SetType("blog_post")

	err := blockType.Validate(block)
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}
}

// TestBlogPostBlockType_AdminFields tests admin fields
func TestBlogPostBlockType_AdminFields(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blogStore := registry.GetBlogStore()
	blockType := NewBlogPostBlockType(blogStore)

	block := cmsstore.NewBlock()
	block.SetType("blog_post")
	block.SetMeta("post_id", "post-123")
	block.SetMeta("show_image", "true")
	block.SetMeta("show_title", "true")
	block.SetMeta("show_date", "true")
	block.SetMeta("show_author", "true")
	block.SetMeta("show_summary", "true")

	req := httptest.NewRequest("GET", "/test", nil)
	fields := blockType.GetAdminFields(block, req)
	if fields == nil {
		t.Error("Expected admin fields, got nil")
		return
	}
}

// TestBlogPostBlockType_SaveAdminFields tests saving admin fields
func TestBlogPostBlockType_SaveAdminFields(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blogStore := registry.GetBlogStore()
	blockType := NewBlogPostBlockType(blogStore)

	block := cmsstore.NewBlock()
	block.SetType("blog_post")

	req := httptest.NewRequest("POST", "/test", nil)
	req.Form = map[string][]string{
		"post_id":      {"post-123"},
		"show_image":   {"true"},
		"show_title":   {"true"},
		"show_date":    {"true"},
		"show_author":  {"true"},
		"show_summary": {"true"},
	}

	err := blockType.SaveAdminFields(req, block)
	if err != nil {
		t.Errorf("Expected no error saving admin fields, got: %v", err)
		return
	}

	if block.Meta("post_id") != "post-123" {
		t.Errorf("Expected post_id to be 'post-123', got '%s'", block.Meta("post_id"))
	}

	if block.Meta("show_image") != "true" {
		t.Errorf("Expected show_image to be 'true', got '%s'", block.Meta("show_image"))
	}
}

// TestBlogPostBlockType_Render tests the Render method includes tags and share links
func TestBlogPostBlockType_Render(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blogStore := registry.GetBlogStore()
	blockType := NewBlogPostBlockType(blogStore)

	// Create a published post with tags
	post := blogstore.NewPost()
	post.SetID("test-post-123")
	post.SetTitle("Test Post Title")
	// This is genrated autimatically - post.Set("alias", "test-post-slug")
	post.SetContent("<p>Test post content</p>")
	post.SetSummary("Test summary")
	post.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	post.SetPublishedAt("2024-01-01 12:00:00")
	post.SetImageUrl("https://example.com/test-image.jpg")

	err := blogStore.PostCreate(t.Context(), post)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Create tag taxonomy first
	tagTaxonomy := blogstore.NewTaxonomy()
	tagTaxonomy.SetName("Tag")
	tagTaxonomy.SetSlug(blogstore.TAXONOMY_TAG)
	err = blogStore.TaxonomyCreate(t.Context(), tagTaxonomy)
	if err != nil {
		t.Fatalf("Failed to create tag taxonomy: %v", err)
	}

	// Create tags and associate with post
	tag1 := blogstore.NewTerm()
	tag1.SetID("tag-1")
	tag1.SetName("SEO")
	tag1.SetSlug("seo")
	tag1.SetTaxonomyID(tagTaxonomy.GetID())

	tag2 := blogstore.NewTerm()
	tag2.SetID("tag-2")
	tag2.SetName("Digital Marketing")
	tag2.SetSlug("digital-marketing")
	tag2.SetTaxonomyID(tagTaxonomy.GetID())

	err = blogStore.TermCreate(t.Context(), tag1)
	if err != nil {
		t.Fatalf("Failed to create tag1: %v", err)
	}
	err = blogStore.TermCreate(t.Context(), tag2)
	if err != nil {
		t.Fatalf("Failed to create tag2: %v", err)
	}

	// Associate tags with post
	err = blogStore.PostAddTerm(t.Context(), post.GetID(), tag1.GetID())
	if err != nil {
		t.Fatalf("Failed to add tag1 to post: %v", err)
	}
	err = blogStore.PostAddTerm(t.Context(), post.GetID(), tag2.GetID())
	if err != nil {
		t.Fatalf("Failed to add tag2 to post: %v", err)
	}

	// Create block referencing the post
	block := cmsstore.NewBlock()
	block.SetType("blog_post")
	block.SetMeta("post_id", post.GetID())
	block.SetMeta("show_title", "true")
	block.SetMeta("show_image", "true")
	block.SetMeta("show_date", "false")
	block.SetMeta("show_author", "false")
	block.SetMeta("show_summary", "false")
	block.SetMeta("show_prev_next", "false")

	// Render
	req := httptest.NewRequest("GET", "/blog/test-post-123/test-post-slug", nil)
	ctx := cmsstore.RequestToContext(t.Context(), req)
	ctx = cmsstore.WithVarsContext(ctx)

	html, err := blockType.Render(ctx, block)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Verify post title is in output
	if !strings.Contains(html, "Test Post Title") {
		t.Error("Expected rendered HTML to contain post title")
	}

	// Verify tags are in output
	if !strings.Contains(html, "SEO") {
		t.Error("Expected rendered HTML to contain tag 'SEO'")
	}
	if !strings.Contains(html, "Digital Marketing") {
		t.Error("Expected rendered HTML to contain tag 'Digital Marketing'")
	}

	// Verify tag links use correct slugs
	if !strings.Contains(html, "/blog/tag/seo") {
		t.Error("Expected rendered HTML to contain link to /blog/tag/seo")
	}
	if !strings.Contains(html, "/blog/tag/digital-marketing") {
		t.Error("Expected rendered HTML to contain link to /blog/tag/digital-marketing")
	}

	// Verify share section is in output
	if !strings.Contains(html, "Follow us") {
		t.Error("Expected rendered HTML to contain 'Follow us' text")
	}

	// Verify social share links
	if !strings.Contains(html, "facebook.com/sharer") {
		t.Error("Expected rendered HTML to contain Facebook share link")
	}
	if !strings.Contains(html, "x.com/intent/tweet") {
		t.Error("Expected rendered HTML to contain X/Twitter share link")
	}
	if !strings.Contains(html, "linkedin.com/sharing") {
		t.Error("Expected rendered HTML to contain LinkedIn share link")
	}
	if !strings.Contains(html, "pinterest.com/pin/create") {
		t.Error("Expected rendered HTML to contain Pinterest share link")
	}
}

// TestBlogPostBlockType_RenderWithNoTags tests rendering a post without tags
func TestBlogPostBlockType_RenderWithNoTags(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blogStore := registry.GetBlogStore()
	blockType := NewBlogPostBlockType(blogStore)

	// Create a published post without tags
	post := blogstore.NewPost()
	post.SetID("test-post-no-tags")
	post.SetTitle("Post Without Tags")
	// post.SetSlug("post-without-tags") - generated automatically
	post.SetContent("<p>No tags here</p>")
	post.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	post.SetPublishedAt("2024-01-01 12:00:00")

	err := blogStore.PostCreate(t.Context(), post)
	if err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}

	// Create block
	block := cmsstore.NewBlock()
	block.SetType("blog_post")
	block.SetMeta("post_id", post.GetID())
	block.SetMeta("show_title", "true")
	block.SetMeta("show_image", "false")
	block.SetMeta("show_date", "false")
	block.SetMeta("show_author", "false")
	block.SetMeta("show_summary", "false")
	block.SetMeta("show_prev_next", "false")

	// Render
	req := httptest.NewRequest("GET", "/blog/test-post-no-tags/post-without-tags", nil)
	ctx := cmsstore.RequestToContext(t.Context(), req)
	ctx = cmsstore.WithVarsContext(ctx)

	html, err := blockType.Render(ctx, block)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Verify post renders without tags section
	if !strings.Contains(html, "Post Without Tags") {
		t.Error("Expected rendered HTML to contain post title")
	}

	// When no tags, the tags/share section should not appear
	// (we only render it when len(postTags) > 0)
}
