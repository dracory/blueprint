package partials

import (
	"testing"

	"project/internal/layouts"
)

func TestPageHeader(t *testing.T) {
	// Test with icon and title only
	result := PageHeader("bi-house", "Test Title")
	if result == nil {
		t.Error("PageHeader() should not return nil")
	}

	// Test with icon, title, and breadcrumbs
	breadcrumbs := []layouts.Breadcrumb{
		{Name: "Home", URL: "/"},
		{Name: "Test", URL: "/test"},
	}
	result = PageHeader("bi-house", "Test Title", breadcrumbs)
	if result == nil {
		t.Error("PageHeader() with breadcrumbs should not return nil")
	}
}
