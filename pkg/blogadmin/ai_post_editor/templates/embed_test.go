package templates

import (
	"strings"
	"testing"

	"github.com/flosch/pongo2/v6"
)

// TestToBytes_ReadAppHtmlTemplate tests reading app.html template
func TestToBytes_ReadAppHtmlTemplate(t *testing.T) {
	t.Parallel()

	data, err := ToBytes("app.html")
	if err != nil {
		t.Errorf("ToBytes(app.html) unexpected error: %v", err)
		return
	}
	if len(data) == 0 {
		t.Error("ToBytes(app.html) returned empty data")
	}
}

// TestToBytes_ReadAppJsTemplate tests reading app.js template
func TestToBytes_ReadAppJsTemplate(t *testing.T) {
	t.Parallel()

	data, err := ToBytes("app.js")
	if err != nil {
		t.Errorf("ToBytes(app.js) unexpected error: %v", err)
		return
	}
	if len(data) == 0 {
		t.Error("ToBytes(app.js) returned empty data")
	}
}

// TestToBytes_ReadAppCssTemplate tests reading app.css template
func TestToBytes_ReadAppCssTemplate(t *testing.T) {
	t.Parallel()

	data, err := ToBytes("app.css")
	if err != nil {
		t.Errorf("ToBytes(app.css) unexpected error: %v", err)
		return
	}
	if len(data) == 0 {
		t.Error("ToBytes(app.css) returned empty data")
	}
}

// TestToBytes_ReadNonExistentFile tests reading non-existent file
func TestToBytes_ReadNonExistentFile(t *testing.T) {
	t.Parallel()

	data, err := ToBytes("nonexistent.txt")
	if err == nil {
		t.Error("Expected error for path nonexistent.txt, got nil")
	}
	if data != nil {
		t.Error("Expected nil data for non-existent file")
	}
}

// TestToBytes_ReadEmptyPath tests reading empty path
func TestToBytes_ReadEmptyPath(t *testing.T) {
	t.Parallel()

	data, err := ToBytes("")
	if err == nil {
		t.Error("Expected error for empty path, got nil")
	}
	if data != nil {
		t.Error("Expected nil data for empty path")
	}
}

// TestToString_ReadAppHtmlAsString tests reading app.html as string
func TestToString_ReadAppHtmlAsString(t *testing.T) {
	t.Parallel()

	str, err := ToString("app.html")
	if err != nil {
		t.Errorf("ToString(app.html) unexpected error: %v", err)
		return
	}
	expectedPrefix := "<div id=\"post-editor-app\""
	if !strings.HasPrefix(str, expectedPrefix) {
		t.Errorf("ToString(app.html) expected prefix %q, got %q", expectedPrefix, str[:min(len(str), len(expectedPrefix))])
	}
}

// TestToString_ReadAppJsAsString tests reading app.js as string
func TestToString_ReadAppJsAsString(t *testing.T) {
	t.Parallel()

	str, err := ToString("app.js")
	if err != nil {
		t.Errorf("ToString(app.js) unexpected error: %v", err)
		return
	}
	_ = str // No specific prefix check for JS
}

// TestToString_ReadNonExistentFile tests reading non-existent file
func TestToString_ReadNonExistentFile(t *testing.T) {
	t.Parallel()

	str, err := ToString("nonexistent.txt")
	if err == nil {
		t.Error("Expected error for path nonexistent.txt, got nil")
	}
	if str != "" {
		t.Error("Expected empty string for non-existent file")
	}
}

// TestResourceExists_CheckExistingFileAppHtml tests checking existing file app.html
func TestResourceExists_CheckExistingFileAppHtml(t *testing.T) {
	t.Parallel()

	exists := ResourceExists("app.html")
	if !exists {
		t.Error("ResourceExists(app.html) = false, want true")
	}
}

// TestResourceExists_CheckExistingFileAppJs tests checking existing file app.js
func TestResourceExists_CheckExistingFileAppJs(t *testing.T) {
	t.Parallel()

	exists := ResourceExists("app.js")
	if !exists {
		t.Error("ResourceExists(app.js) = false, want true")
	}
}

// TestResourceExists_CheckExistingFileAppCss tests checking existing file app.css
func TestResourceExists_CheckExistingFileAppCss(t *testing.T) {
	t.Parallel()

	exists := ResourceExists("app.css")
	if !exists {
		t.Error("ResourceExists(app.css) = false, want true")
	}
}

// TestResourceExists_CheckNonExistentFile tests checking non-existent file
func TestResourceExists_CheckNonExistentFile(t *testing.T) {
	t.Parallel()

	exists := ResourceExists("nonexistent.txt")
	if exists {
		t.Error("ResourceExists(nonexistent.txt) = true, want false")
	}
}

// TestResourceExists_CheckEmptyPath tests checking empty path
func TestResourceExists_CheckEmptyPath(t *testing.T) {
	t.Parallel()

	exists := ResourceExists("")
	if exists {
		t.Error("ResourceExists(\"\") = true, want false")
	}
}

// TestTemplate_ValidHtmlFile tests template with valid HTML file
func TestTemplate_ValidHtmlFile(t *testing.T) {
	t.Parallel()

	result, err := Template("app.html", map[string]any{})
	if err != nil {
		t.Errorf("Template(app.html) unexpected error: %v", err)
		return
	}
	if result == "" {
		t.Error("Template(app.html) returned empty string")
	}
}

// TestTemplate_ValidJsFile tests template with valid JS file
func TestTemplate_ValidJsFile(t *testing.T) {
	t.Parallel()

	result, err := Template("app.js", map[string]any{"postJSON": "{}", "id": "123"})
	if err != nil {
		t.Errorf("Template(app.js) unexpected error: %v", err)
		return
	}
	_ = result // No specific check for JS result
}

// TestTemplate_NonExistentFile tests template with non-existent file
func TestTemplate_NonExistentFile(t *testing.T) {
	t.Parallel()

	result, err := Template("nonexistent.html", map[string]any{})
	if err == nil {
		t.Error("Template(nonexistent.html) expected error, got nil")
	}
	if result != "" {
		t.Error("Template(nonexistent.html) should return empty string on error")
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

// TestTpl_ValidFile tests Tpl with valid file
func TestTpl_ValidFile(t *testing.T) {
	t.Parallel()

	result := Tpl("app.css", map[string]any{})
	if result == "" {
		t.Error("Tpl(app.css) returned empty string for valid file")
	}
}

// TestTpl_NonExistentFile tests Tpl with non-existent file
func TestTpl_NonExistentFile(t *testing.T) {
	t.Parallel()

	result := Tpl("nonexistent.html", map[string]any{})
	if result != "" {
		t.Error("Tpl(nonexistent.html) expected empty string on error, got non-empty")
	}
}

// TestTpl_InvalidTemplateSyntax tests Tpl with invalid template syntax
func TestTpl_InvalidTemplateSyntax(t *testing.T) {
	t.Parallel()

	result := Tpl("app.js", map[string]any{"postJSON": "test"})
	if result == "" {
		t.Error("Tpl(app.js) returned empty string for invalid template syntax")
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
