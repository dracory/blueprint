package layouts

import (
	"strings"
	"testing"
)

func TestBreadcrumbs(t *testing.T) {
	// Test with empty breadcrumbs
	result := Breadcrumbs([]Breadcrumb{})
	if result == nil {
		t.Fatal("Breadcrumbs() should not return nil")
	}
	html := result.ToHTML()
	if !strings.Contains(html, "<nav") {
		t.Error("Expected <nav> element in empty breadcrumbs")
	}
	if !strings.Contains(html, "breadcrumb") {
		t.Error("Expected breadcrumb class in empty breadcrumbs")
	}

	// Test with single breadcrumb
	result = Breadcrumbs([]Breadcrumb{
		{Name: "Home", URL: "/"},
	})
	if result == nil {
		t.Fatal("Breadcrumbs() should not return nil")
	}
	html = result.ToHTML()
	if !strings.Contains(html, "Home") {
		t.Error("Expected breadcrumb name 'Home' in output")
	}
	if !strings.Contains(html, `href="/"`) {
		t.Error("Expected breadcrumb URL in href attribute")
	}
	if !strings.Contains(html, "breadcrumb-item") {
		t.Error("Expected breadcrumb-item class")
	}

	// Test with multiple breadcrumbs
	result = Breadcrumbs([]Breadcrumb{
		{Name: "Home", URL: "/"},
		{Name: "Blog", URL: "/blog"},
		{Name: "Post", URL: "/blog/post"},
	})
	if result == nil {
		t.Fatal("Breadcrumbs() should not return nil")
	}
	html = result.ToHTML()
	if !strings.Contains(html, "Home") {
		t.Error("Expected 'Home' in multiple breadcrumbs")
	}
	if !strings.Contains(html, "Blog") {
		t.Error("Expected 'Blog' in multiple breadcrumbs")
	}
	if !strings.Contains(html, "Post") {
		t.Error("Expected 'Post' in multiple breadcrumbs")
	}
	if !strings.Contains(html, `href="/blog/post"`) {
		t.Error("Expected last breadcrumb URL in href")
	}

	// Test with breadcrumb with icon
	result = Breadcrumbs([]Breadcrumb{
		{Name: "Home", URL: "/", Icon: "bi-house"},
	})
	if result == nil {
		t.Fatal("Breadcrumbs() with icon should not return nil")
	}
	html = result.ToHTML()
	if !strings.Contains(html, "bi-house") {
		t.Error("Expected icon class 'bi-house' in output")
	}
	if !strings.Contains(html, "bi ") {
		t.Error("Expected 'bi' class for icon")
	}

	// Test with breadcrumb without URL
	result = Breadcrumbs([]Breadcrumb{
		{Name: "Current Page", URL: ""},
	})
	if result == nil {
		t.Fatal("Breadcrumbs() without URL should not return nil")
	}
	html = result.ToHTML()
	if !strings.Contains(html, "Current Page") {
		t.Error("Expected breadcrumb name 'Current Page' in output")
	}
}

func TestBreadcrumbStruct(t *testing.T) {
	// Test Breadcrumb struct initialization
	breadcrumb := Breadcrumb{
		Icon: "bi-house",
		Name: "Home",
		URL:  "/",
	}

	if breadcrumb.Name != "Home" {
		t.Errorf("Name = %q, want %q", breadcrumb.Name, "Home")
	}
	if breadcrumb.URL != "/" {
		t.Errorf("URL = %q, want %q", breadcrumb.URL, "/")
	}
	if breadcrumb.Icon != "bi-house" {
		t.Errorf("Icon = %q, want %q", breadcrumb.Icon, "bi-house")
	}

	// Test with empty values
	breadcrumb = Breadcrumb{}
	if breadcrumb.Name != "" {
		t.Errorf("Name = %q, want empty string", breadcrumb.Name)
	}
}

func TestOptionsStruct(t *testing.T) {
	// Test Options struct initialization
	options := Options{
		AppName:         "TestApp",
		WebsiteSection:  "Blog",
		Title:           "Test Title",
		MetaDescription: "Test Description",
		MetaKeywords:    "test, keywords",
		ImageURL:        "https://example.com/image.jpg",
		CanonicalURL:    "https://example.com/canonical",
	}

	if options.AppName != "TestApp" {
		t.Errorf("AppName = %q, want %q", options.AppName, "TestApp")
	}
	if options.WebsiteSection != "Blog" {
		t.Errorf("WebsiteSection = %q, want %q", options.WebsiteSection, "Blog")
	}
	if options.Title != "Test Title" {
		t.Errorf("Title = %q, want %q", options.Title, "Test Title")
	}

	// Test with empty options
	options = Options{}
	if options.AppName != "" {
		t.Errorf("AppName = %q, want empty string", options.AppName)
	}
}
