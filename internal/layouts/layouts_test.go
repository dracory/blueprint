package layouts

import (
	"net/http"
	"strings"
	"testing"

	"project/internal/testutils"

	"github.com/dracory/hb"
)

func TestLogoHTML(t *testing.T) {
	html := LogoHTML()

	if html == "" {
		t.Error("LogoHTML() should return non-empty HTML")
	}
	if !strings.Contains(html, "img") {
		t.Error("LogoHTML() should contain an img tag")
	}
	if !strings.Contains(html, "dracory") {
		t.Error("LogoHTML() should contain the dracory URL")
	}
}

func TestAdminLogoHtml(t *testing.T) {
	html := adminLogoHtml()

	if html == "" {
		t.Error("adminLogoHtml() should return non-empty HTML")
	}
	if !strings.Contains(html, "img") {
		t.Error("adminLogoHtml() should contain an img tag")
	}
}

func TestUserLogoHtml(t *testing.T) {
	html := userLogoHtml()

	if html == "" {
		t.Error("userLogoHtml() should return non-empty HTML")
	}
	if !strings.Contains(html, "img") {
		t.Error("userLogoHtml() should contain an img tag")
	}
}

func TestFaviconURL(t *testing.T) {
	url := FaviconURL()

	if url == "" {
		t.Error("FaviconURL() should return non-empty URL")
	}
	if !strings.HasPrefix(url, "data:image/svg+xml,") {
		t.Error("FaviconURL() should return a data URI")
	}
	if !strings.Contains(url, "svg") {
		t.Error("FaviconURL() should contain SVG content")
	}
}

func TestOptions(t *testing.T) {
	content := hb.Paragraph().Text("Test content")
	opts := Options{
		AppName:         "TestApp",
		WebsiteSection:  "test",
		Title:           "Test Title",
		Content:         content,
		ScriptURLs:      []string{"https://example.com/script.js"},
		Scripts:         []string{"console.log('test');"},
		StyleURLs:       []string{"https://example.com/style.css"},
		Styles:          []string{"body { color: red; }"},
		MetaDescription: "Test description",
		MetaKeywords:    "test, keywords",
		ImageURL:        "https://example.com/image.png",
		CanonicalURL:    "https://example.com/page",
	}

	if opts.AppName != "TestApp" {
		t.Errorf("AppName = %q, want %q", opts.AppName, "TestApp")
	}
	if opts.Title != "Test Title" {
		t.Errorf("Title = %q, want %q", opts.Title, "Test Title")
	}
	if opts.MetaDescription != "Test description" {
		t.Errorf("MetaDescription = %q, want %q", opts.MetaDescription, "Test description")
	}
	if len(opts.ScriptURLs) != 1 {
		t.Errorf("len(ScriptURLs) = %d, want 1", len(opts.ScriptURLs))
	}
	if len(opts.Scripts) != 1 {
		t.Errorf("len(Scripts) = %d, want 1", len(opts.Scripts))
	}
	if len(opts.StyleURLs) != 1 {
		t.Errorf("len(StyleURLs) = %d, want 1", len(opts.StyleURLs))
	}
	if len(opts.Styles) != 1 {
		t.Errorf("len(Styles) = %d, want 1", len(opts.Styles))
	}
}

func TestAdminPage(t *testing.T) {
	// Test with no elements
	result := AdminPage()
	if result == nil {
		t.Fatal("AdminPage() with no elements should not return nil")
	}
	html := result.ToHTML()
	if !strings.Contains(html, "container") {
		t.Error("AdminPage() should contain container class")
	}
	if !strings.Contains(html, "py-4") {
		t.Error("AdminPage() should contain py-4 class")
	}

	// Test with single element
	element := hb.Div().Text("Test")
	result = AdminPage(element)
	if result == nil {
		t.Fatal("AdminPage(element) should not return nil")
	}
	html = result.ToHTML()
	if !strings.Contains(html, "Test") {
		t.Error("AdminPage(element) should contain the element text")
	}

	// Test with multiple elements
	element1 := hb.Div().Text("First")
	element2 := hb.Div().Text("Second")
	result = AdminPage(element1, element2)
	if result == nil {
		t.Fatal("AdminPage(elements) should not return nil")
	}
	html = result.ToHTML()
	if !strings.Contains(html, "First") {
		t.Error("AdminPage(elements) should contain the first element")
	}
	if !strings.Contains(html, "Second") {
		t.Error("AdminPage(elements) should contain the second element")
	}

	// Test with nil element
	result = AdminPage(nil)
	if result == nil {
		t.Fatal("AdminPage(nil) should not return nil")
	}
}

func TestNewBlankLayout(t *testing.T) {
	registry := testutils.Setup()
	r := &http.Request{}
	opts := Options{
		AppName:    "TestApp",
		Title:      "Test",
		Content:    hb.Div().Text("Test content"),
		ScriptURLs: []string{"https://example.com/script.js"},
		Scripts:    []string{"alert('test');"},
		StyleURLs:  []string{"https://example.com/style.css"},
		Styles:     []string{"body { color: red; }"},
	}

	layout := NewBlankLayout(registry, r, opts)
	if layout == nil {
		t.Fatal("NewBlankLayout() should return non-nil")
	}

	// Test ToHTML
	html := layout.ToHTML()
	if html == "" {
		t.Error("ToHTML() should return non-empty HTML")
	}
	if !strings.Contains(html, "Test content") {
		t.Error("ToHTML() should contain the content")
	}
	if !strings.Contains(html, "bootstrap") {
		t.Error("ToHTML() should contain bootstrap CSS")
	}
}

func TestBreadcrumb(t *testing.T) {
	// Test single breadcrumb
	bc := Breadcrumb{Name: "Home", URL: "/"}
	if bc.Name != "Home" {
		t.Errorf("Name = %q, want %q", bc.Name, "Home")
	}
	if bc.URL != "/" {
		t.Errorf("URL = %q, want %q", bc.URL, "/")
	}
}
