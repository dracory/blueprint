package media

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHandlerWithNilStorage verifies that Handler returns the proper error message
// when the controller is constructed without a storage implementation.
func TestHandlerWithNilStorage(t *testing.T) {
	c := &mediaController{storage: nil}

	req := httptest.NewRequest(http.MethodGet, "/media/example.txt", nil)
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
	c := mediaController{}
	got := c.findFileName("/hello.txt")
	if got != "hello.txt" {
		t.Fatalf("expected 'hello.txt', got: %q", got)
	}
}

// TestFindFileNameNested ensures last segment extraction works from nested paths
func TestFindFileNameNested(t *testing.T) {
	c := mediaController{}
	got := c.findFileName("/a/b/c/image.png")
	if got != "image.png" {
		t.Fatalf("expected 'image.png', got: %q", got)
	}
}

// TestFindFileNameEmpty ensures empty or slash-only path returns empty string
func TestFindFileNameEmpty(t *testing.T) {
	c := mediaController{}
	if got := c.findFileName(""); got != "" {
		t.Fatalf("expected '', got: %q", got)
	}
	if got := c.findFileName("/"); got != "" {
		t.Fatalf("expected '', got: %q", got)
	}
}

// TestFindExtensionHappyPath ensures extension is parsed after last dot
func TestFindExtensionHappyPath(t *testing.T) {
	c := mediaController{}
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
	c := mediaController{}
	if got := c.findExtension("/"); got != "" {
		t.Fatalf("expected '', got: %q", got)
	}
	if got := c.findExtension(""); got != "" {
		t.Fatalf("expected '', got: %q", got)
	}
}

// TestFindExtensionNoDot ensures filenames without dot produce empty extension
func TestFindExtensionNoDot(t *testing.T) {
	c := mediaController{}
	if got := c.findExtension("/dir/readme"); got != "" {
		t.Fatalf("expected '', got: %q", got)
	}
}

// TestFindMIMETypeAllFormats ensures all supported MIME types are correctly detected
func TestFindMIMETypeAllFormats_Html(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("html")
	if got != "text/html" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "html", got, "text/html")
	}
}

func TestFindMIMETypeAllFormats_Css(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("css")
	if got != "text/css" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "css", got, "text/css")
	}
}

func TestFindMIMETypeAllFormats_Js(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("js")
	if got != "application/javascript" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "js", got, "application/javascript")
	}
}

func TestFindMIMETypeAllFormats_Json(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("json")
	if got != "application/json" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "json", got, "application/json")
	}
}

func TestFindMIMETypeAllFormats_Png(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("png")
	if got != "image/png" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "png", got, "image/png")
	}
}

func TestFindMIMETypeAllFormats_Jpg(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("jpg")
	if got != "image/jpeg" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "jpg", got, "image/jpeg")
	}
}

func TestFindMIMETypeAllFormats_Jpeg(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("jpeg")
	if got != "image/jpeg" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "jpeg", got, "image/jpeg")
	}
}

func TestFindMIMETypeAllFormats_Gif(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("gif")
	if got != "image/gif" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "gif", got, "image/gif")
	}
}

func TestFindMIMETypeAllFormats_Svg(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("svg")
	if got != "image/svg+xml" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "svg", got, "image/svg+xml")
	}
}

func TestFindMIMETypeAllFormats_Ico(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("ico")
	if got != "image/x-icon" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "ico", got, "image/x-icon")
	}
}

func TestFindMIMETypeAllFormats_Pdf(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("pdf")
	if got != "application/pdf" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "pdf", got, "application/pdf")
	}
}

func TestFindMIMETypeAllFormats_Zip(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("zip")
	if got != "application/zip" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "zip", got, "application/zip")
	}
}

func TestFindMIMETypeAllFormats_Mp3(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("mp3")
	if got != "audio/mpeg" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "mp3", got, "audio/mpeg")
	}
}

func TestFindMIMETypeAllFormats_Webm(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("webm")
	if got != "video/webm" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "webm", got, "video/webm")
	}
}

func TestFindMIMETypeAllFormats_Unknown(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("unknown")
	if got != "application/octet-stream" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "unknown", got, "application/octet-stream")
	}
}

func TestFindMIMETypeAllFormats_Empty(t *testing.T) {
	c := mediaController{}
	got := c.findMIMEType("")
	if got != "application/octet-stream" {
		t.Errorf("findMIMEType(%q) = %q, want %q", "", got, "application/octet-stream")
	}
}

// TestNewMediaController verifies the constructor properly initializes the controller
func TestNewMediaController(t *testing.T) {
	// Test with nil storage
	c := NewMediaController(nil)
	if c == nil {
		t.Fatal("NewMediaController(nil) should not return nil")
	}
	if c.storage != nil {
		t.Error("expected storage to be nil")
	}
}

// TestFindFileNameWithSpecialCharacters tests edge cases in filename extraction
func TestFindFileNameWithSpecialCharacters_FileWithDashes(t *testing.T) {
	c := mediaController{}
	got := c.findFileName("/path/to/file-with-dashes.txt")
	if got != "file-with-dashes.txt" {
		t.Errorf("findFileName(%q) = %q, want %q", "/path/to/file-with-dashes.txt", got, "file-with-dashes.txt")
	}
}

func TestFindFileNameWithSpecialCharacters_FileWithUnderscores(t *testing.T) {
	c := mediaController{}
	got := c.findFileName("/path/to/file_with_underscores.txt")
	if got != "file_with_underscores.txt" {
		t.Errorf("findFileName(%q) = %q, want %q", "/path/to/file_with_underscores.txt", got, "file_with_underscores.txt")
	}
}

func TestFindFileNameWithSpecialCharacters_FileWithMultipleDots(t *testing.T) {
	c := mediaController{}
	got := c.findFileName("/path/to/file.multiple.dots.txt")
	if got != "file.multiple.dots.txt" {
		t.Errorf("findFileName(%q) = %q, want %q", "/path/to/file.multiple.dots.txt", got, "file.multiple.dots.txt")
	}
}

func TestFindFileNameWithSpecialCharacters_FilenameWithoutPath(t *testing.T) {
	c := mediaController{}
	got := c.findFileName("filename.txt")
	if got != "filename.txt" {
		t.Errorf("findFileName(%q) = %q, want %q", "filename.txt", got, "filename.txt")
	}
}

func TestFindFileNameWithSpecialCharacters_TrailingSlash(t *testing.T) {
	c := mediaController{}
	got := c.findFileName("/trailing/slash/")
	if got != "slash" {
		t.Errorf("findFileName(%q) = %q, want %q", "/trailing/slash/", got, "slash")
	}
}

// TestFindExtensionWithMultipleDots tests extension extraction with multiple dots
// Note: Implementation returns nameParts[1] (first extension after dot)
func TestFindExtensionWithMultipleDots_TarGz(t *testing.T) {
	c := mediaController{}
	got := c.findExtension("/path/to/file.tar.gz")
	if got != "tar" {
		t.Errorf("findExtension(%q) = %q, want %q", "/path/to/file.tar.gz", got, "tar")
	}
}

func TestFindExtensionWithMultipleDots_TarBz2(t *testing.T) {
	c := mediaController{}
	got := c.findExtension("/path/to/archive.tar.bz2")
	if got != "tar" {
		t.Errorf("findExtension(%q) = %q, want %q", "/path/to/archive.tar.bz2", got, "tar")
	}
}

func TestFindExtensionWithMultipleDots_HiddenFile(t *testing.T) {
	c := mediaController{}
	got := c.findExtension("/path/.hiddenfile")
	if got != "hiddenfile" {
		t.Errorf("findExtension(%q) = %q, want %q", "/path/.hiddenfile", got, "hiddenfile")
	}
}
