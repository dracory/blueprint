package search

import (
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/cmsstore"
)

// TestSearchBlockType_BasicProperties tests basic properties
func TestSearchBlockType_BasicProperties(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCmsStore(true),
		testutils.WithBlogStore(true),
	)

	blockType := NewSearchBlockType(registry.GetCmsStore(), registry.GetBlogStore())

	if blockType.TypeKey() != "search" {
		t.Errorf("Expected type key 'search', got '%s'", blockType.TypeKey())
	}

	if blockType.TypeLabel() != "Search" {
		t.Errorf("Expected type label 'Search', got '%s'", blockType.TypeLabel())
	}
}

// TestSearchBlockType_GetPreview tests preview
func TestSearchBlockType_GetPreview(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCmsStore(true),
		testutils.WithBlogStore(true),
	)

	cmsStore := registry.GetCmsStore()
	blockType := NewSearchBlockType(cmsStore, registry.GetBlogStore())

	// Create a real block for testing
	block := cmsstore.NewBlock()
	block.SetType("search")
	block.SetName("Test Search Block")

	// Test with default settings (both enabled)
	preview := blockType.GetPreview(block)
	if preview != "Search (Pages, Blog Posts)" {
		t.Errorf("Expected preview 'Search (Pages, Blog Posts)', got '%s'", preview)
	}

	// Test with only pages
	block.SetMeta("show_posts", "false")
	preview = blockType.GetPreview(block)
	if preview != "Search (Pages)" {
		t.Errorf("Expected preview 'Search (Pages)', got '%s'", preview)
	}

	// Test with only posts
	block.SetMeta("show_pages", "false")
	block.SetMeta("show_posts", "true")
	preview = blockType.GetPreview(block)
	if preview != "Search (Blog Posts)" {
		t.Errorf("Expected preview 'Search (Blog Posts)', got '%s'", preview)
	}

	// Test with none selected
	block.SetMeta("show_posts", "false")
	preview = blockType.GetPreview(block)
	if preview != "Search (no content types selected)" {
		t.Errorf("Expected preview 'Search (no content types selected)', got '%s'", preview)
	}
}

// TestSearchBlockType_Validate tests validation
func TestSearchBlockType_Validate(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCmsStore(true),
		testutils.WithBlogStore(true),
	)

	cmsStore := registry.GetCmsStore()
	blockType := NewSearchBlockType(cmsStore, registry.GetBlogStore())

	block := cmsstore.NewBlock()
	block.SetType("search")

	err := blockType.Validate(block)
	if err != nil {
		t.Errorf("Expected no validation error, got: %v", err)
	}
}

// TestSearchBlockType_AdminFields tests admin fields
func TestSearchBlockType_AdminFields(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCmsStore(true),
		testutils.WithBlogStore(true),
	)

	cmsStore := registry.GetCmsStore()
	blockType := NewSearchBlockType(cmsStore, registry.GetBlogStore())

	block := cmsstore.NewBlock()
	block.SetType("search")
	block.SetMeta("placeholder", "Find articles...")
	block.SetMeta("results_per_page", "20")
	block.SetMeta("show_pages", "true")
	block.SetMeta("show_posts", "true")

	req := httptest.NewRequest("GET", "/test", nil)
	fields := blockType.GetAdminFields(block, req)
	if fields == nil {
		t.Error("Expected admin fields, got nil")
		return
	}
}

// TestSearchBlockType_SaveAdminFields tests saving admin fields
func TestSearchBlockType_SaveAdminFields(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCmsStore(true),
		testutils.WithBlogStore(true),
	)

	cmsStore := registry.GetCmsStore()
	blockType := NewSearchBlockType(cmsStore, registry.GetBlogStore())

	block := cmsstore.NewBlock()
	block.SetType("search")

	req := httptest.NewRequest("POST", "/test", nil)
	req.Form = map[string][]string{
		"placeholder":      {"Search our site..."},
		"results_per_page": {"15"},
		"show_pages":       {"true"},
		"show_posts":       {"false"},
	}

	err := blockType.SaveAdminFields(req, block)
	if err != nil {
		t.Errorf("Expected no error saving admin fields, got: %v", err)
		return
	}

	if block.Meta("placeholder") != "Search our site..." {
		t.Errorf("Expected placeholder to be 'Search our site...', got '%s'", block.Meta("placeholder"))
	}

	if block.Meta("results_per_page") != "15" {
		t.Errorf("Expected results_per_page to be '15', got '%s'", block.Meta("results_per_page"))
	}

	if block.Meta("show_posts") != "false" {
		t.Errorf("Expected show_posts to be 'false', got '%s'", block.Meta("show_posts"))
	}
}

// TestSearchBlockType_Render_EmptyQuery tests rendering without search query
func TestSearchBlockType_Render_EmptyQuery(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCmsStore(true),
		testutils.WithBlogStore(true),
	)

	cmsStore := registry.GetCmsStore()
	blockType := NewSearchBlockType(cmsStore, registry.GetBlogStore())

	block := cmsstore.NewBlock()
	block.SetType("search")
	block.SetMeta("placeholder", "Search...")
	block.SetMeta("results_per_page", "10")
	block.SetMeta("show_pages", "true")
	block.SetMeta("show_posts", "true")

	// Create request without query
	req, _ := testutils.NewRequest("GET", "/search", testutils.NewRequestOptions{})
	ctx := cmsstore.RequestToContext(req.Context(), req)

	html, err := blockType.Render(ctx, block)
	if err != nil {
		t.Errorf("Expected no error rendering, got: %v", err)
		return
	}

	// Should contain search box
	if !strings.Contains(html, "Search...") {
		t.Error("Expected HTML to contain placeholder text")
	}
}

// TestSearchBlockType_Render_WithQuery tests rendering with search query
func TestSearchBlockType_Render_WithQuery(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithCmsStore(true),
		testutils.WithBlogStore(true),
	)

	cmsStore := registry.GetCmsStore()
	blockType := NewSearchBlockType(cmsStore, registry.GetBlogStore())

	block := cmsstore.NewBlock()
	block.SetType("search")
	block.SetMeta("placeholder", "Search...")
	block.SetMeta("results_per_page", "10")
	block.SetMeta("show_pages", "true")
	block.SetMeta("show_posts", "true")

	// Create request with search query
	req, _ := testutils.NewRequest("GET", "/search", testutils.NewRequestOptions{
		QueryParams: map[string][]string{"q": {"test"}},
	})
	ctx := cmsstore.RequestToContext(req.Context(), req)

	html, err := blockType.Render(ctx, block)
	if err != nil {
		t.Errorf("Expected no error rendering, got: %v", err)
		return
	}

	// Should contain search query or results message
	if !strings.Contains(html, "test") && !strings.Contains(html, "Found") && !strings.Contains(html, "No results") {
		t.Error("Expected HTML to contain search results or no results message")
	}
}

// TestStripHTML tests the HTML stripping function
func TestStripHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<p>Hello World</p>", "Hello World"},
		{"<div><p>Test</p></div>", "Test"},
		{"Plain text", "Plain text"},
		{"<a href='#'>Link</a> text", "Link text"},
		{"", ""},
	}

	for _, test := range tests {
		result := stripHTML(test.input)
		if result != test.expected {
			t.Errorf("stripHTML(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

// TestEscapeHTML tests the HTML escaping function
func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"<script>", "&lt;script&gt;"},
		{"Test & Example", "Test &amp; Example"},
		{"\"quoted\"", "&quot;quoted&quot;"},
		{"Normal text", "Normal text"},
	}

	for _, test := range tests {
		result := escapeHTML(test.input)
		if result != test.expected {
			t.Errorf("escapeHTML(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

// TestBuildPageURL tests the URL building function
func TestBuildPageURL(t *testing.T) {
	tests := []struct {
		query   string
		pageNum int
		want    string
	}{
		{"test", 0, "?q=test"},
		{"test", 1, "?page=2&q=test"},
		{"", 0, ""},
		{"", 1, "?page=2"},
	}

	for _, test := range tests {
		result := buildPageURL(test.query, test.pageNum)
		if result != test.want {
			t.Errorf("buildPageURL(%q, %d) = %q, want %q", test.query, test.pageNum, result, test.want)
		}
	}
}
