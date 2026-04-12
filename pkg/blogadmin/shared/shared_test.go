package shared

import (
	"strings"
	"testing"
)

// TestNewLinks tests the Links constructor
func TestNewLinks(t *testing.T) {
	t.Parallel()

	// Test with custom base URL
	links := NewLinks("/custom/base")
	if links.BaseURL != "/custom/base" {
		t.Errorf("Expected BaseURL to be '/custom/base', got: %s", links.BaseURL)
	}

	// Test with empty base URL (should use default)
	links = NewLinks("")
	if links.BaseURL != "/admin/blog" {
		t.Errorf("Expected BaseURL to be '/admin/blog', got: %s", links.BaseURL)
	}

	// Test with trailing slash (should be removed)
	links = NewLinks("/admin/blog/")
	if links.BaseURL != "/admin/blog" {
		t.Errorf("Expected BaseURL to be '/admin/blog', got: %s", links.BaseURL)
	}
}

// TestBuildURL tests the buildURL method
func TestBuildURL(t *testing.T) {
	t.Parallel()

	links := NewLinks("/admin/blog")

	// Test without params
	url := links.buildURL("test-controller", nil)
	if url != "/admin/blog?controller=test-controller" {
		t.Errorf("Expected '/admin/blog?controller=test-controller', got: %s", url)
	}

	// Test with params (check that both controller and param are present, order may vary)
	url = links.buildURL("test-controller", map[string]string{"param1": "value1"})
	if !containsAll(url, []string{"controller=test-controller", "param1=value1"}) {
		t.Errorf("Expected URL to contain controller and param1, got: %s", url)
	}

	// Test with multiple params (check that all params are present, order may vary)
	url = links.buildURL("test-controller", map[string]string{"param1": "value1", "param2": "value2"})
	if !containsAll(url, []string{"controller=test-controller", "param1=value1", "param2=value2"}) {
		t.Errorf("Expected URL to contain all params, got: %s", url)
	}
}

// containsAll checks if a string contains all substrings
func containsAll(s string, substrings []string) bool {
	for _, sub := range substrings {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

// TestMergeParams tests the mergeParams function
func TestMergeParams(t *testing.T) {
	t.Parallel()

	// Test with no params
	result := mergeParams()
	if len(result) != 0 {
		t.Errorf("Expected empty map, got: %v", result)
	}

	// Test with one param
	result = mergeParams(map[string]string{"key1": "value1"})
	if result["key1"] != "value1" {
		t.Errorf("Expected 'value1', got: %s", result["key1"])
	}

	// Test with multiple params (should merge)
	result = mergeParams(
		map[string]string{"key1": "value1"},
		map[string]string{"key2": "value2"},
	)
	if result["key1"] != "value1" || result["key2"] != "value2" {
		t.Errorf("Expected both values, got: %v", result)
	}

	// Test with overlapping keys (later should win)
	result = mergeParams(
		map[string]string{"key1": "value1"},
		map[string]string{"key1": "value2"},
	)
	if result["key1"] != "value2" {
		t.Errorf("Expected 'value2' (later wins), got: %s", result["key1"])
	}
}

// TestLinksMethods tests all the link generation methods
func TestLinksMethods(t *testing.T) {
	t.Parallel()

	links := NewLinks("/admin/blog")

	tests := []struct {
		name     string
		method   func(...map[string]string) string
		expected string
	}{
		{"Home", links.Home, "/admin/blog?controller=post-manager"},
		{"PostCreate", links.PostCreate, "/admin/blog?controller=post-create"},
		{"PostDelete", links.PostDelete, "/admin/blog?controller=post-delete"},
		{"PostManager", links.PostManager, "/admin/blog?controller=post-manager"},
		{"PostUpdate", links.PostUpdate, "/admin/blog?controller=post-update"},
		{"PostUpdateV1", links.PostUpdateV1, "/admin/blog?controller=post-update-v1"},
		{"BlogSettings", links.BlogSettings, "/admin/blog?controller=blog-settings"},
		{"AiTools", links.AiTools, "/admin/blog?controller=ai-tools"},
		{"AiPostContentUpdate", links.AiPostContentUpdate, "/admin/blog?controller=ai-post-content-update"},
		{"AiPostGenerator", links.AiPostGenerator, "/admin/blog?controller=ai-post-generator"},
		{"AiTitleGenerator", links.AiTitleGenerator, "/admin/blog?controller=ai-title-generator"},
		{"AiPostEditor", links.AiPostEditor, "/admin/blog?controller=ai-post-editor"},
		{"AiTest", links.AiTest, "/admin/blog?controller=ai-test"},
		{"Dashboard", links.Dashboard, "/admin/blog?controller=dashboard"},
		{"CategoryManager", links.CategoryManager, "/admin/blog?controller=category-manager"},
		{"TagManager", links.TagManager, "/admin/blog?controller=tag-manager"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.method()
			if result != tt.expected {
				t.Errorf("Expected '%s', got: %s", tt.expected, result)
			}
		})
	}
}

// TestLinksMethodsWithParams tests link generation with parameters
func TestLinksMethodsWithParams(t *testing.T) {
	t.Parallel()

	links := NewLinks("/admin/blog")

	// Test Home with params (order may vary)
	result := links.Home(map[string]string{"page": "2"})
	if !containsAll(result, []string{"controller=post-manager", "page=2"}) {
		t.Errorf("Expected URL to contain controller and page, got: %s", result)
	}

	// Test PostUpdate with params (order may vary)
	result = links.PostUpdate(map[string]string{"post_id": "123"})
	if !containsAll(result, []string{"controller=post-update", "post_id=123"}) {
		t.Errorf("Expected URL to contain controller and post_id, got: %s", result)
	}
}
