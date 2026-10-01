package layout

import (
	"net/http"
	"strings"
	"testing"

	"project/internal/testutils"

	baselayouts "github.com/dracory/base/layouts"
	"github.com/dracory/hb"
)

func TestNew(t *testing.T) {
	app := testutils.Setup()
	r := &http.Request{}
	opts := baselayouts.Options{
		Title:   "Admin Test",
		Content: hb.Div().Text("Admin content"),
	}

	layout := New(app, r, opts)
	if layout == nil {
		t.Fatal("New() should return non-nil")
	}

	html := layout.ToHTML()
	if html == "" {
		t.Error("New().ToHTML() should return non-empty HTML")
	}
	if !strings.Contains(html, "Admin Test") {
		t.Error("New() should contain the title")
	}
	if !strings.Contains(html, "SectionNavbar") {
		t.Error("New() should render the admin navbar")
	}
}

func TestNewCrud(t *testing.T) {
	app := testutils.Setup()
	r := &http.Request{}

	html := NewCrud(app, r, "Admin CRUD Test", "<p>Admin CRUD content</p>", []string{}, "", []string{}, "")
	if html == "" {
		t.Error("NewCrud() should return non-empty HTML")
	}
	if !strings.Contains(html, "Admin CRUD Test") {
		t.Error("NewCrud() should contain the title")
	}
}

func TestPage(t *testing.T) {
	// Test with no elements
	result := Page()
	if result == nil {
		t.Fatal("Page() with no elements should not return nil")
	}
	html := result.ToHTML()
	if !strings.Contains(html, "container") {
		t.Error("Page() should contain container class")
	}
	if !strings.Contains(html, "py-4") {
		t.Error("Page() should contain py-4 class")
	}

	// Test with single element
	element := hb.Div().Text("Test")
	result = Page(element)
	if result == nil {
		t.Fatal("Page(element) should not return nil")
	}
	html = result.ToHTML()
	if !strings.Contains(html, "Test") {
		t.Error("Page(element) should contain the element text")
	}

	// Test with multiple elements
	element1 := hb.Div().Text("First")
	element2 := hb.Div().Text("Second")
	result = Page(element1, element2)
	if result == nil {
		t.Fatal("Page(elements) should not return nil")
	}
	html = result.ToHTML()
	if !strings.Contains(html, "First") {
		t.Error("Page(elements) should contain the first element")
	}
	if !strings.Contains(html, "Second") {
		t.Error("Page(elements) should contain the second element")
	}

	// Test with nil element
	result = Page(nil)
	if result == nil {
		t.Fatal("Page(nil) should not return nil")
	}
}
