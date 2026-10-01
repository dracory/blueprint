package layouts

import (
	"net/http"
	"strings"
	"testing"

	baselayouts "github.com/dracory/base/layouts"

	"project/internal/app"
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
	if !strings.Contains(html, "data:image/svg+xml;base64,") {
		t.Error("LogoHTML() should embed the logo as a data URI")
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
	opts := baselayouts.Options{
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

func TestNewBlankLayout(t *testing.T) {
	app := testutils.Setup()
	r := &http.Request{}
	opts := baselayouts.Options{
		AppName:    "TestApp",
		Title:      "Test",
		Content:    hb.Div().Text("Test content"),
		ScriptURLs: []string{"https://example.com/script.js"},
		Scripts:    []string{"alert('test');"},
		StyleURLs:  []string{"https://example.com/style.css"},
		Styles:     []string{"body { color: red; }"},
	}

	layout := NewBlankLayout(app, r, opts)
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
	// Test single baselayouts.Breadcrumb
	bc := baselayouts.Breadcrumb{Name: "Home", URL: "/"}
	if bc.Name != "Home" {
		t.Errorf("Name = %q, want %q", bc.Name, "Home")
	}
	if bc.URL != "/" {
		t.Errorf("URL = %q, want %q", bc.URL, "/")
	}
}

func TestNewPageLayout(t *testing.T) {
	app := testutils.Setup()
	r := &http.Request{}
	opts := baselayouts.Options{
		Title:   "Page Test",
		Content: hb.Div().Text("Page content"),
	}

	layout := NewPageLayout(app, r, opts)
	if layout == nil {
		t.Fatal("NewPageLayout() should return non-nil")
	}

	html := layout.ToHTML()
	if html == "" {
		t.Error("NewPageLayout().ToHTML() should return non-empty HTML")
	}
	if !strings.Contains(html, "Page Test") {
		t.Error("NewPageLayout() should contain the title")
	}
	if !strings.Contains(html, "Page content") {
		t.Error("NewPageLayout() should contain the content")
	}
	if !strings.Contains(html, "SectionNavbar") {
		t.Error("NewPageLayout() should render the default navbar")
	}
}

func TestNewPageLayoutWithConfig(t *testing.T) {
	application := testutils.Setup()
	r := &http.Request{}
	opts := baselayouts.Options{
		Title:   "Custom Test",
		Content: hb.Div().Text("Custom content"),
	}

	layout := NewPageLayoutWithConfig(application, r, opts, PageLayoutConfig{
		IsAdmin: true,
		NavbarFn: func(app app.AppInterface, r *http.Request) hb.TagInterface {
			return hb.Section().ID("CustomNavbar")
		},
	})

	html := layout.ToHTML()
	if html == "" {
		t.Error("NewPageLayoutWithConfig().ToHTML() should return non-empty HTML")
	}
	if !strings.Contains(html, "CustomNavbar") {
		t.Error("NewPageLayoutWithConfig() should render the custom navbar")
	}
	if strings.Contains(html, "<footer") {
		t.Error("NewPageLayoutWithConfig() with IsAdmin should not render the public footer")
	}
}
