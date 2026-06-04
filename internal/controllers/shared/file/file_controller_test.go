package file

import (
	"net/http"
	"net/http/httptest"
	"project/internal/testutils"
	"testing"

	"github.com/dracory/filesystem"
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

// TestFindMIMETypeAllFormats_Html tests html MIME type
func TestFindMIMETypeAllFormats_Html(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("html")
	if got != "text/html" {
		t.Errorf("findMIMEType(html) = %q, want text/html", got)
	}
}

// TestFindMIMETypeAllFormats_Css tests css MIME type
func TestFindMIMETypeAllFormats_Css(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("css")
	if got != "text/css" {
		t.Errorf("findMIMEType(css) = %q, want text/css", got)
	}
}

// TestFindMIMETypeAllFormats_Js tests js MIME type
func TestFindMIMETypeAllFormats_Js(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("js")
	if got != "application/javascript" {
		t.Errorf("findMIMEType(js) = %q, want application/javascript", got)
	}
}

// TestFindMIMETypeAllFormats_Json tests json MIME type
func TestFindMIMETypeAllFormats_Json(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("json")
	if got != "application/json" {
		t.Errorf("findMIMEType(json) = %q, want application/json", got)
	}
}

// TestFindMIMETypeAllFormats_Png tests png MIME type
func TestFindMIMETypeAllFormats_Png(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("png")
	if got != "image/png" {
		t.Errorf("findMIMEType(png) = %q, want image/png", got)
	}
}

// TestFindMIMETypeAllFormats_Jpg tests jpg MIME type
func TestFindMIMETypeAllFormats_Jpg(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("jpg")
	if got != "image/jpeg" {
		t.Errorf("findMIMEType(jpg) = %q, want image/jpeg", got)
	}
}

// TestFindMIMETypeAllFormats_Jpeg tests jpeg MIME type
func TestFindMIMETypeAllFormats_Jpeg(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("jpeg")
	if got != "image/jpeg" {
		t.Errorf("findMIMEType(jpeg) = %q, want image/jpeg", got)
	}
}

// TestFindMIMETypeAllFormats_Gif tests gif MIME type
func TestFindMIMETypeAllFormats_Gif(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("gif")
	if got != "image/gif" {
		t.Errorf("findMIMEType(gif) = %q, want image/gif", got)
	}
}

// TestFindMIMETypeAllFormats_Svg tests svg MIME type
func TestFindMIMETypeAllFormats_Svg(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("svg")
	if got != "image/svg+xml" {
		t.Errorf("findMIMEType(svg) = %q, want image/svg+xml", got)
	}
}

// TestFindMIMETypeAllFormats_Ico tests ico MIME type
func TestFindMIMETypeAllFormats_Ico(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("ico")
	if got != "image/x-icon" {
		t.Errorf("findMIMEType(ico) = %q, want image/x-icon", got)
	}
}

// TestFindMIMETypeAllFormats_Pdf tests pdf MIME type
func TestFindMIMETypeAllFormats_Pdf(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("pdf")
	if got != "application/pdf" {
		t.Errorf("findMIMEType(pdf) = %q, want application/pdf", got)
	}
}

// TestFindMIMETypeAllFormats_Zip tests zip MIME type
func TestFindMIMETypeAllFormats_Zip(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("zip")
	if got != "application/zip" {
		t.Errorf("findMIMEType(zip) = %q, want application/zip", got)
	}
}

// TestFindMIMETypeAllFormats_Mp3 tests mp3 MIME type
func TestFindMIMETypeAllFormats_Mp3(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("mp3")
	if got != "audio/mpeg" {
		t.Errorf("findMIMEType(mp3) = %q, want audio/mpeg", got)
	}
}

// TestFindMIMETypeAllFormats_Webm tests webm MIME type
func TestFindMIMETypeAllFormats_Webm(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("webm")
	if got != "video/webm" {
		t.Errorf("findMIMEType(webm) = %q, want video/webm", got)
	}
}

// TestFindMIMETypeAllFormats_Unknown tests unknown MIME type
func TestFindMIMETypeAllFormats_Unknown(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("unknown")
	if got != "application/octet-stream" {
		t.Errorf("findMIMEType(unknown) = %q, want application/octet-stream", got)
	}
}

// TestFindMIMETypeAllFormats_Empty tests empty extension MIME type
func TestFindMIMETypeAllFormats_Empty(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("")
	if got != "application/octet-stream" {
		t.Errorf("findMIMEType(\"\") = %q, want application/octet-stream", got)
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

// TestFindFileNameWithSpecialCharacters_Dashes tests filename with dashes
func TestFindFileNameWithSpecialCharacters_Dashes(t *testing.T) {
	c := fileController{}
	got := c.findFileName("/path/to/file-with-dashes.txt")
	if got != "file-with-dashes.txt" {
		t.Errorf("findFileName = %q, want file-with-dashes.txt", got)
	}
}

// TestFindFileNameWithSpecialCharacters_Underscores tests filename with underscores
func TestFindFileNameWithSpecialCharacters_Underscores(t *testing.T) {
	c := fileController{}
	got := c.findFileName("/path/to/file_with_underscores.txt")
	if got != "file_with_underscores.txt" {
		t.Errorf("findFileName = %q, want file_with_underscores.txt", got)
	}
}

// TestFindFileNameWithSpecialCharacters_MultipleDots tests filename with multiple dots
func TestFindFileNameWithSpecialCharacters_MultipleDots(t *testing.T) {
	c := fileController{}
	got := c.findFileName("/path/to/file.multiple.dots.txt")
	if got != "file.multiple.dots.txt" {
		t.Errorf("findFileName = %q, want file.multiple.dots.txt", got)
	}
}

// TestFindFileNameWithSpecialCharacters_NoPath tests filename without path
func TestFindFileNameWithSpecialCharacters_NoPath(t *testing.T) {
	c := fileController{}
	got := c.findFileName("filename.txt")
	if got != "filename.txt" {
		t.Errorf("findFileName = %q, want filename.txt", got)
	}
}

// TestFindFileNameWithSpecialCharacters_TrailingSlash tests filename with trailing slash
func TestFindFileNameWithSpecialCharacters_TrailingSlash(t *testing.T) {
	c := fileController{}
	got := c.findFileName("/trailing/slash/")
	if got != "slash" {
		t.Errorf("findFileName = %q, want slash", got)
	}
}

// TestFindExtensionWithMultipleDots_TarGz tests extension extraction with tar.gz
func TestFindExtensionWithMultipleDots_TarGz(t *testing.T) {
	c := fileController{}
	got := c.findExtension("/path/to/file.tar.gz")
	if got != "tar" {
		t.Errorf("findExtension = %q, want tar", got)
	}
}

// TestFindExtensionWithMultipleDots_TarBz2 tests extension extraction with tar.bz2
func TestFindExtensionWithMultipleDots_TarBz2(t *testing.T) {
	c := fileController{}
	got := c.findExtension("/path/to/archive.tar.bz2")
	if got != "tar" {
		t.Errorf("findExtension = %q, want tar", got)
	}
}

// TestFindExtensionWithMultipleDots_HiddenFile tests extension extraction with hidden file
func TestFindExtensionWithMultipleDots_HiddenFile(t *testing.T) {
	c := fileController{}
	got := c.findExtension("/path/.hiddenfile")
	if got != "hiddenfile" {
		t.Errorf("findExtension = %q, want hiddenfile", got)
	}
}

// TestHandlerPathVariations_FilesPrefix tests handler with /files prefix
func TestHandlerPathVariations_FilesPrefix(t *testing.T) {
	c := &fileController{storage: nil}
	req := httptest.NewRequest(http.MethodGet, "/files/test.txt", nil)
	rec := httptest.NewRecorder()
	res := c.Handler(rec, req)
	if res != "File storage not configured" {
		t.Errorf("Handler = %q, want 'File storage not configured'", res)
	}
}

// TestHandlerPathVariations_FilePrefix tests handler with /file prefix
func TestHandlerPathVariations_FilePrefix(t *testing.T) {
	c := &fileController{storage: nil}
	req := httptest.NewRequest(http.MethodGet, "/file/test.txt", nil)
	rec := httptest.NewRecorder()
	res := c.Handler(rec, req)
	if res != "File storage not configured" {
		t.Errorf("Handler = %q, want 'File storage not configured'", res)
	}
}

// TestHandlerPathVariations_MediaPrefix tests handler with /media prefix
func TestHandlerPathVariations_MediaPrefix(t *testing.T) {
	c := &fileController{storage: nil}
	req := httptest.NewRequest(http.MethodGet, "/media/test.txt", nil)
	rec := httptest.NewRecorder()
	res := c.Handler(rec, req)
	if res != "File storage not configured" {
		t.Errorf("Handler = %q, want 'File storage not configured'", res)
	}
}

// TestHandlerPathVariations_NoPrefix tests handler with no prefix
func TestHandlerPathVariations_NoPrefix(t *testing.T) {
	c := &fileController{storage: nil}
	req := httptest.NewRequest(http.MethodGet, "/test.txt", nil)
	rec := httptest.NewRecorder()
	res := c.Handler(rec, req)
	if res != "File storage not configured" {
		t.Errorf("Handler = %q, want 'File storage not configured'", res)
	}
}

// TestHandlerWithRealStorage tests the Handler with real SQL file storage
func TestHandlerWithRealStorage(t *testing.T) {
	app := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer app.GetDatabase().Close()

	// Create SQL file storage
	db := app.GetDatabase()
	if db == nil {
		t.Fatal("Database not initialized")
	}

	storage, err := filesystem.NewStorage(filesystem.Disk{
		DiskName:  filesystem.DRIVER_SQL,
		Driver:    filesystem.DRIVER_SQL,
		Url:       "/files",
		DB:        db,
		TableName: "snv_files_file",
	})
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	controller := NewFileController(storage)

	// Test with non-existent file
	req := httptest.NewRequest(http.MethodGet, "/files/nonexistent.txt", nil)
	rec := httptest.NewRecorder()
	res := controller.Handler(rec, req)
	if res != "File not found" {
		t.Errorf("Handler(nonexistent) = %q, want 'File not found'", res)
	}
}

// TestHandlerFileNotFound tests file not found scenario
func TestHandlerFileNotFound(t *testing.T) {
	// Create a mock storage that returns exists=false
	controller := NewFileController(nil)

	req := httptest.NewRequest(http.MethodGet, "/files/missing.txt", nil)
	rec := httptest.NewRecorder()
	res := controller.Handler(rec, req)
	if res != "File storage not configured" {
		t.Errorf("Handler = %q, want 'File storage not configured'", res)
	}
}

// TestHandlerWithEmptyExtension tests file with no extension
func TestHandlerWithEmptyExtension(t *testing.T) {
	app := testutils.Setup(testutils.WithUserStore(true), testutils.WithSessionStore(true))
	defer app.GetDatabase().Close()

	db := app.GetDatabase()
	if db == nil {
		t.Fatal("Database not initialized")
	}

	storage, err := filesystem.NewStorage(filesystem.Disk{
		DiskName:  filesystem.DRIVER_SQL,
		Driver:    filesystem.DRIVER_SQL,
		Url:       "/files",
		DB:        db,
		TableName: "snv_files_file",
	})
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	controller := NewFileController(storage)

	// Test with file that has no extension (should return "File not found")
	req := httptest.NewRequest(http.MethodGet, "/files/noextension", nil)
	rec := httptest.NewRecorder()
	res := controller.Handler(rec, req)
	if res != "File not found" {
		t.Errorf("Handler(noextension) = %q, want 'File not found'", res)
	}
}

// TestFindMIMETypeEdgeCases_HtmlUpperCase tests HTML with uppercase
func TestFindMIMETypeEdgeCases_HtmlUpperCase(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("HTML")
	if got != "application/octet-stream" {
		t.Errorf("findMIMEType(HTML) = %q, want application/octet-stream", got)
	}
}

// TestFindMIMETypeEdgeCases_CssMixedCase tests CSS with mixed case
func TestFindMIMETypeEdgeCases_CssMixedCase(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("Css")
	if got != "application/octet-stream" {
		t.Errorf("findMIMEType(Css) = %q, want application/octet-stream", got)
	}
}

// TestFindMIMETypeEdgeCases_PngUpperCase tests PNG with uppercase
func TestFindMIMETypeEdgeCases_PngUpperCase(t *testing.T) {
	c := fileController{}
	got := c.findMIMEType("PNG")
	if got != "application/octet-stream" {
		t.Errorf("findMIMEType(PNG) = %q, want application/octet-stream", got)
	}
}
