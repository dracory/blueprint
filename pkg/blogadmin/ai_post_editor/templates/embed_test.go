package templates

import (
	"strings"
	"testing"

	"github.com/flosch/pongo2/v6"
)

// TestToBytes tests the ToBytes function
func TestToBytes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		path        string
		expectError bool
	}{
		{
			name:        "Read app.html template",
			path:        "app.html",
			expectError: false,
		},
		{
			name:        "Read app.js template",
			path:        "app.js",
			expectError: false,
		},
		{
			name:        "Read app.css template",
			path:        "app.css",
			expectError: false,
		},
		{
			name:        "Read non-existent file",
			path:        "nonexistent.txt",
			expectError: true,
		},
		{
			name:        "Read empty path",
			path:        "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := ToBytes(tt.path)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for path %q, got nil", tt.path)
				}
				return
			}
			if err != nil {
				t.Errorf("ToBytes(%q) unexpected error: %v", tt.path, err)
				return
			}
			if len(data) == 0 {
				t.Errorf("ToBytes(%q) returned empty data", tt.path)
			}
		})
	}
}

// TestToString tests the ToString function
func TestToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		expectError    bool
		expectedPrefix string
	}{
		{
			name:           "Read app.html as string",
			path:           "app.html",
			expectError:    false,
			expectedPrefix: "<div id=\"post-editor-app\"",
		},
		{
			name:           "Read app.js as string",
			path:           "app.js",
			expectError:    false,
			expectedPrefix: "",
		},
		{
			name:        "Read non-existent file",
			path:        "nonexistent.txt",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str, err := ToString(tt.path)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for path %q, got nil", tt.path)
				}
				return
			}
			if err != nil {
				t.Errorf("ToString(%q) unexpected error: %v", tt.path, err)
				return
			}
			if tt.expectedPrefix != "" && !strings.HasPrefix(str, tt.expectedPrefix) {
				t.Errorf("ToString(%q) expected prefix %q, got %q", tt.path, tt.expectedPrefix, str[:min(len(str), len(tt.expectedPrefix))])
			}
		})
	}
}

// TestResourceExists tests the ResourceExists function
func TestResourceExists(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "Check existing file app.html",
			path:     "app.html",
			expected: true,
		},
		{
			name:     "Check existing file app.js",
			path:     "app.js",
			expected: true,
		},
		{
			name:     "Check existing file app.css",
			path:     "app.css",
			expected: true,
		},
		{
			name:     "Check non-existent file",
			path:     "nonexistent.txt",
			expected: false,
		},
		{
			name:     "Check empty path",
			path:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists := ResourceExists(tt.path)
			if exists != tt.expected {
				t.Errorf("ResourceExists(%q) = %v, want %v", tt.path, exists, tt.expected)
			}
		})
	}
}

// TestTemplate tests the Template function
func TestTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		path        string
		data        map[string]any
		expectError bool
	}{
		{
			name:        "Template with valid HTML file",
			path:        "app.html",
			data:        map[string]any{},
			expectError: false,
		},
		{
			name:        "Template with valid JS file",
			path:        "app.js",
			data:        map[string]any{"postJSON": "{}", "id": "123"},
			expectError: false,
		},
		{
			name:        "Template with non-existent file",
			path:        "nonexistent.html",
			data:        map[string]any{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Template(tt.path, tt.data)
			if tt.expectError {
				if err == nil {
					t.Errorf("Template(%q) expected error, got nil", tt.path)
				}
				return
			}
			if err != nil {
				t.Errorf("Template(%q) unexpected error: %v", tt.path, err)
				return
			}
			if result == "" && tt.path == "app.html" {
				t.Errorf("Template(%q) returned empty string", tt.path)
			}
		})
	}
}

// TestTemplateCaching tests that templates are properly cached
func TestTemplateCaching(t *testing.T) {
	// Note: Not parallel due to global cache manipulation

	// First call should parse and cache
	result1, err := Template("app.css", map[string]any{})
	if err != nil {
		t.Fatalf("First Template call failed: %v", err)
	}

	// Second call should use cached version
	result2, err := Template("app.css", map[string]any{})
	if err != nil {
		t.Fatalf("Second Template call failed: %v", err)
	}

	// Results should be identical
	if result1 != result2 {
		t.Error("Template caching produced different results")
	}

	// Verify cache has entry
	cacheMutex.RLock()
	_, found := templateCache["app.css"]
	cacheMutex.RUnlock()

	if !found {
		t.Error("Template was not cached")
	}
}

// TestTpl tests the Tpl function (shortcut that ignores errors)
func TestTpl(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		path        string
		data        map[string]any
		expectEmpty bool
	}{
		{
			name:        "Tpl with valid file",
			path:        "app.css",
			data:        map[string]any{},
			expectEmpty: false,
		},
		{
			name:        "Tpl with non-existent file",
			path:        "nonexistent.html",
			data:        map[string]any{},
			expectEmpty: true,
		},
		{
			name:        "Tpl with invalid template syntax",
			path:        "app.js",
			data:        map[string]any{"postJSON": "test"},
			expectEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Tpl(tt.path, tt.data)
			if tt.expectEmpty && result != "" {
				t.Errorf("Tpl(%q) expected empty string on error, got %q", tt.path, result)
			}
		})
	}
}

// TestTemplateConcurrentAccess tests thread safety
func TestTemplateConcurrentAccess(t *testing.T) {
	// Note: Not parallel due to global cache manipulation

	// Clear cache first
	cacheMutex.Lock()
	templateCache = make(map[string]*pongo2.Template)
	cacheMutex.Unlock()

	// Run concurrent template renders
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			_, err := Template("app.css", map[string]any{})
			if err != nil {
				t.Errorf("Concurrent Template call failed: %v", err)
			}
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
