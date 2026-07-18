package blogpostlist

import (
	"context"
	"net/http/httptest"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/cmsstore"
)

// TestBlogPostListBlockType_BasicProperties tests basic properties
func TestBlogPostListBlockType_BasicProperties(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blockType := NewBlogPostListBlockType(app.GetBlogStore())

	if blockType.TypeKey() != "blog_post_list" {
		t.Errorf("Expected type key 'blog_post_list', got '%s'", blockType.TypeKey())
	}

	if blockType.TypeLabel() != "Blog Post List" {
		t.Errorf("Expected type label 'Blog Post List', got '%s'", blockType.TypeLabel())
	}
}

// TestBlogPostListBlockType_GetPreview tests preview
func TestBlogPostListBlockType_GetPreview(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blockType := NewBlogPostListBlockType(app.GetBlogStore())
	block := cmsstore.NewBlock()
	block.SetType("blog_post_list")

	preview := blockType.GetPreview(block)
	if preview != "Blog Post List" {
		t.Errorf("Expected preview 'Blog Post List', got '%s'", preview)
	}
}

// TestBlogPostListBlockType_Validate tests validation
func TestBlogPostListBlockType_Validate(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blockType := NewBlogPostListBlockType(app.GetBlogStore())
	block := cmsstore.NewBlock()
	block.SetType("blog_post_list")

	err := blockType.Validate(block)
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}
}

// TestBlogPostListBlockType_AdminFields tests admin fields
func TestBlogPostListBlockType_AdminFields(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blockType := NewBlogPostListBlockType(app.GetBlogStore())
	block := cmsstore.NewBlock()
	block.SetType("blog_post_list")
	block.SetMeta("posts_per_page", "6")
	block.SetMeta("columns", "3")
	block.SetMeta("show_pagination", "true")
	block.SetMeta("show_images", "true")
	block.SetMeta("show_summary", "true")
	block.SetMeta("show_date", "true")

	req := httptest.NewRequest("GET", "/test", nil)
	fields := blockType.GetAdminFields(block, req)
	if fields == nil {
		t.Error("Expected admin fields, got nil")
		return
	}
}

// TestBlogPostListBlockType_SaveAdminFields tests saving admin fields
func TestBlogPostListBlockType_SaveAdminFields(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blockType := NewBlogPostListBlockType(app.GetBlogStore())
	block := cmsstore.NewBlock()
	block.SetType("blog_post_list")

	req := httptest.NewRequest("POST", "/test", nil)
	req.Form = map[string][]string{
		"posts_per_page":  {"6"},
		"columns":         {"3"},
		"show_pagination": {"true"},
		"show_images":     {"true"},
		"show_summary":    {"true"},
		"show_date":       {"true"},
	}

	err := blockType.SaveAdminFields(req, block)
	if err != nil {
		t.Errorf("Expected no error saving admin fields, got: %v", err)
		return
	}

	if block.Meta("posts_per_page") != "6" {
		t.Errorf("Expected posts_per_page to be '6', got '%s'", block.Meta("posts_per_page"))
	}

	if block.Meta("columns") != "3" {
		t.Errorf("Expected columns to be '3', got '%s'", block.Meta("columns"))
	}
}

// TestExtractTagSlugFromURL tests the tag extraction from URL
func TestExtractTagSlugFromURL(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blockType := NewBlogPostListBlockType(app.GetBlogStore())

	req := httptest.NewRequest("GET", "/tag/journalists-journorequest-pr-media-press/", nil)
	ctx := cmsstore.RequestToContext(context.Background(), req)
	result := blockType.extractTagSlugFromURL(ctx)
	if result != "journalists-journorequest-pr-media-press" {
		t.Errorf("extractTagSlugFromURL() = %v, want journalists-journorequest-pr-media-press", result)
	}
}

// TestExtractTagSlugFromURL_URLWithTagNoTrailingSlash tests URL with tag no trailing slash
func TestExtractTagSlugFromURL_URLWithTagNoTrailingSlash(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blockType := NewBlogPostListBlockType(app.GetBlogStore())

	req := httptest.NewRequest("GET", "/tag/my-tag", nil)
	ctx := cmsstore.RequestToContext(context.Background(), req)
	result := blockType.extractTagSlugFromURL(ctx)
	if result != "my-tag" {
		t.Errorf("extractTagSlugFromURL() = %v, want my-tag", result)
	}
}

// TestExtractTagSlugFromURL_URLWithoutTag tests URL without tag
func TestExtractTagSlugFromURL_URLWithoutTag(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blockType := NewBlogPostListBlockType(app.GetBlogStore())

	req := httptest.NewRequest("GET", "/blog/some-post", nil)
	ctx := cmsstore.RequestToContext(context.Background(), req)
	result := blockType.extractTagSlugFromURL(ctx)
	if result != "" {
		t.Errorf("extractTagSlugFromURL() = %v, want empty string", result)
	}
}

// TestExtractTagSlugFromURL_RootURL tests root URL
func TestExtractTagSlugFromURL_RootURL(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blockType := NewBlogPostListBlockType(app.GetBlogStore())

	req := httptest.NewRequest("GET", "/", nil)
	ctx := cmsstore.RequestToContext(context.Background(), req)
	result := blockType.extractTagSlugFromURL(ctx)
	if result != "" {
		t.Errorf("extractTagSlugFromURL() = %v, want empty string", result)
	}
}

// TestExtractTagSlugFromURL_TagWithAdditionalPath tests tag with additional path
func TestExtractTagSlugFromURL_TagWithAdditionalPath(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blockType := NewBlogPostListBlockType(app.GetBlogStore())

	req := httptest.NewRequest("GET", "/tag/my-tag/extra", nil)
	ctx := cmsstore.RequestToContext(context.Background(), req)
	result := blockType.extractTagSlugFromURL(ctx)
	if result != "my-tag" {
		t.Errorf("extractTagSlugFromURL() = %v, want my-tag", result)
	}
}

// TestExtractTagSlugFromURL_EmptyTag tests empty tag
func TestExtractTagSlugFromURL_EmptyTag(t *testing.T) {
	app := testutils.Setup(
		testutils.WithBlogStore(true),
	)

	blockType := NewBlogPostListBlockType(app.GetBlogStore())

	req := httptest.NewRequest("GET", "/tag/", nil)
	ctx := cmsstore.RequestToContext(context.Background(), req)
	result := blockType.extractTagSlugFromURL(ctx)
	if result != "" {
		t.Errorf("extractTagSlugFromURL() = %v, want empty string", result)
	}
}
