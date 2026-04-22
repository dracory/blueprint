package blogpostlist

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dracory/blogstore"
)

func TestRenderBlogPostListHTML_WithSummary(t *testing.T) {
	// Create test posts with and without summaries
	postWithSummary := blogstore.NewPost().
		SetID("1").
		SetTitle("Test Post With Summary").
		SetSummary("This is the post summary")

	postWithoutSummary := blogstore.NewPost().
		SetID("2").
		SetTitle("Test Post Without Summary").
		SetSummary("")

	postList := []blogstore.PostInterface{postWithSummary, postWithoutSummary}

	// Render with showSummary=true
	html, err := renderBlogPostListHTML(postList, 2, 10, 1, false, false, true, false, 2, 150)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Check that summary appears for post with summary
	if !strings.Contains(html, "This is the post summary") {
		t.Errorf("Expected summary text in HTML output\nHTML: %s", html)
	}

	// Check that title appears
	if !strings.Contains(html, "Test Post With Summary") {
		t.Errorf("Expected title in HTML output\nHTML: %s", html)
	}
}

func TestRenderBlogPostListHTML_WithoutSummaryFlag(t *testing.T) {
	postWithSummary := blogstore.NewPost().
		SetID("1").
		SetTitle("Test Post").
		SetSummary("This summary should not appear")

	postList := []blogstore.PostInterface{postWithSummary}

	// Render with showSummary=false
	html, err := renderBlogPostListHTML(postList, 1, 10, 1, false, false, false, false, 1, 150)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Summary should NOT appear when flag is false
	if strings.Contains(html, "This summary should not appear") {
		t.Errorf("Summary should not appear when showSummary=false\nHTML: %s", html)
	}
}

func TestRenderBlogPostListHTML_EmptySummary(t *testing.T) {
	postWithEmptySummary := blogstore.NewPost().
		SetID("1").
		SetTitle("Test Post").
		SetSummary("")

	postList := []blogstore.PostInterface{postWithEmptySummary}

	// Render with showSummary=true but empty summary
	html, err := renderBlogPostListHTML(postList, 1, 10, 1, false, false, true, false, 1, 150)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Empty summary should not create an empty paragraph
	// Just verify it renders without error and has title
	if !strings.Contains(html, "Test Post") {
		t.Errorf("Expected title in HTML output\nHTML: %s", html)
	}
}

func TestRenderBlogPostListHTML_ColumnClasses(t *testing.T) {
	post := blogstore.NewPost().
		SetID("1").
		SetTitle("Test")

	postList := []blogstore.PostInterface{post}

	tests := []struct {
		columns  int
		expected string
	}{
		{1, "col-12"},
		{2, "col-md-6 col-sm-6"},
		{3, "col-md-4 col-sm-6"},
		{4, "col-md-3 col-sm-6"},
		{6, "col-md-2 col-sm-4"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("columns_%d", tt.columns), func(t *testing.T) {
			html, err := renderBlogPostListHTML(postList, 1, 10, 1, false, false, false, false, tt.columns, 150)
			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}

			if !strings.Contains(html, tt.expected) {
				t.Errorf("Expected column class '%s' in HTML\nHTML: %s", tt.expected, html)
			}
		})
	}
}

func TestRenderPagination(t *testing.T) {
	tests := []struct {
		name         string
		totalItems   int
		itemsPerPage int
		currentPage  int
		wantContains []string
	}{
		{
			name:         "single page - no pagination",
			totalItems:   5,
			itemsPerPage: 10,
			currentPage:  1,
			wantContains: []string{},
		},
		{
			name:         "multiple pages - page 1",
			totalItems:   25,
			itemsPerPage: 10,
			currentPage:  1,
			wantContains: []string{"page-item active", "?page=1\">1</a>", "?page=2\">2</a>", "?page=3\">3</a>", "Previous", "Next"},
		},
		{
			name:         "multiple pages - page 2",
			totalItems:   25,
			itemsPerPage: 10,
			currentPage:  2,
			wantContains: []string{"page-item active", "?page=1\">1</a>", "?page=2\">2</a>", "?page=3\">3</a>", "Previous", "Next"},
		},
		{
			name:         "first page - previous disabled",
			totalItems:   25,
			itemsPerPage: 10,
			currentPage:  1,
			wantContains: []string{"page-item disabled", "aria-label=\"Previous\""},
		},
		{
			name:         "last page - next disabled",
			totalItems:   25,
			itemsPerPage: 10,
			currentPage:  3,
			wantContains: []string{"page-item disabled", "aria-label=\"Next\""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderPagination(tt.totalItems, tt.itemsPerPage, tt.currentPage)

			// If single page, should return empty string
			if tt.totalItems <= tt.itemsPerPage {
				if result != "" {
					t.Errorf("Expected empty string for single page, got: %s", result)
				}
				return
			}

			// Check that expected strings are present
			for _, want := range tt.wantContains {
				if !strings.Contains(result, want) {
					t.Errorf("Expected '%s' in pagination output\nGot: %s", want, result)
				}
			}
		})
	}
}
