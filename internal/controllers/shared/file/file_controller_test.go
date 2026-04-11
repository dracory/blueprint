package file

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandlerWithNilStorage verifies that Handler returns the proper error message
// when the controller is constructed without a storage implementation.
func TestHandlerWithNilStorage(t *testing.T) {
	c := &fileController{storage: nil}

	req := httptest.NewRequest(http.MethodGet, "/files/example.txt", nil)
	rec := httptest.NewRecorder()

	res := c.Handler(rec, req)

	if res != "File storage not configured" {
		t.Fatalf("expected 'File storage not configured', got: %q", res)
	}

	// Ensure that nothing was written to the response body in this early-return case
	if rec.Body.Len() != 0 {
		t.Fatalf("expected empty response body when storage not configured, got: %q", rec.Body.String())
	}
}

// TestFindFileNameRoot ensures last segment extraction works from root-level path
func TestFindFileNameRoot(t *testing.T) {
	c := fileController{}
	got := c.findFileName("/hello.txt")
	if got != "hello.txt" {
		t.Fatalf("expected 'hello.txt', got: %q", got)
	}
}

// TestFindFileNameNested ensures last segment extraction works from nested paths
func TestFindFileNameNested(t *testing.T) {
	c := fileController{}
	got := c.findFileName("/a/b/c/image.png")
	if got != "image.png" {
		t.Fatalf("expected 'image.png', got: %q", got)
	}
}

// TestFindFileNameEmpty ensures empty or slash-only path returns empty string
func TestFindFileNameEmpty(t *testing.T) {
	c := fileController{}
	if got := c.findFileName(""); got != "" {
		t.Fatalf("expected '', got: %q", got)
	}
	if got := c.findFileName("/"); got != "" {
		t.Fatalf("expected '', got: %q", got)
	}
}

// TestFindExtensionHappyPath ensures extension is parsed after last dot
func TestFindExtensionHappyPath(t *testing.T) {
	c := fileController{}
	if got := c.findExtension("/dir/file.html"); got != "html" {
		t.Fatalf("expected 'html', got: %q", got)
	}
	if got := c.findExtension("/dir/file.css"); got != "css" {
		t.Fatalf("expected 'css', got: %q", got)
	}
	if got := c.findExtension("/dir/file.jpg"); got != "jpg" {
		t.Fatalf("expected 'jpg', got: %q", got)
	}
}

// TestFindExtensionNoName ensures paths without filename produce empty extension
func TestFindExtensionNoName(t *testing.T) {
	c := fileController{}
	if got := c.findExtension("/"); got != "" {
		t.Fatalf("expected '', got: %q", got)
	}
	if got := c.findExtension(""); got != "" {
		t.Fatalf("expected '', got: %q", got)
	}
}

// TestFindExtensionNoDot ensures filenames without dot produce empty extension
func TestFindExtensionNoDot(t *testing.T) {
	c := fileController{}
	if got := c.findExtension("/dir/readme"); got != "" {
		t.Fatalf("expected '', got: %q", got)
	}
}

// TestFindMIMETypeAllFormats ensures all supported MIME types are correctly detected
func TestFindMIMETypeAllFormats(t *testing.T) {
	c := fileController{}

	tests := []struct {
		ext      string
		expected string
	}{
		{"html", "text/html"},
		{"css", "text/css"},
		{"js", "application/javascript"},
		{"json", "application/json"},
		{"png", "image/png"},
		{"jpg", "image/jpeg"},
		{"jpeg", "image/jpeg"},
		{"gif", "image/gif"},
		{"svg", "image/svg+xml"},
		{"ico", "image/x-icon"},
		{"pdf", "application/pdf"},
		{"zip", "application/zip"},
		{"mp3", "audio/mpeg"},
		{"webm", "video/webm"},
		{"unknown", "application/octet-stream"},
		{"", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			got := c.findMIMEType(tt.ext)
			if got != tt.expected {
				t.Errorf("findMIMEType(%q) = %q, want %q", tt.ext, got, tt.expected)
			}
		})
	}
}

// TestNewFileController verifies the constructor properly initializes the controller
func TestNewFileController(t *testing.T) {
	// Test with nil storage
	c := NewFileController(nil)
	if c == nil {
		t.Fatal("NewFileController(nil) should not return nil")
	}
	if c.storage != nil {
		t.Error("expected storage to be nil")
	}
}

// TestFindFileNameWithSpecialCharacters tests edge cases in filename extraction
func TestFindFileNameWithSpecialCharacters(t *testing.T) {
	c := fileController{}

	tests := []struct {
		path     string
		expected string
	}{
		{"/path/to/file-with-dashes.txt", "file-with-dashes.txt"},
		{"/path/to/file_with_underscores.txt", "file_with_underscores.txt"},
		{"/path/to/file.multiple.dots.txt", "file.multiple.dots.txt"},
		{"filename.txt", "filename.txt"},
		{"/trailing/slash/", "slash"}, // Implementation returns last segment
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := c.findFileName(tt.path)
			if got != tt.expected {
				t.Errorf("findFileName(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}

// TestFindExtensionWithMultipleDots tests extension extraction with multiple dots
// Note: Implementation returns nameParts[1] (first extension after dot)
func TestFindExtensionWithMultipleDots(t *testing.T) {
	c := fileController{}

	tests := []struct {
		path     string
		expected string
	}{
		{"/path/to/file.tar.gz", "tar"},     // Implementation returns first part after dot
		{"/path/to/archive.tar.bz2", "tar"}, // Implementation returns first part after dot
		{"/path/.hiddenfile", "hiddenfile"}, // Hidden file starting with dot
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := c.findExtension(tt.path)
			if got != tt.expected {
				t.Errorf("findExtension(%q) = %q, want %q", tt.path, got, tt.expected)
			}
		})
	}
}
