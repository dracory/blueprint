package blogpostlist

import (
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

func TestRenderBlogPostListHTML_ColumnClasses_Column1(t *testing.T) {
	post := blogstore.NewPost().
		SetID("1").
		SetTitle("Test")

	postList := []blogstore.PostInterface{post}

	html, err := renderBlogPostListHTML(postList, 1, 10, 1, false, false, false, false, 1, 150)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(html, "col-12") {
		t.Errorf("Expected column class 'col-12' in HTML\nHTML: %s", html)
	}
}

func TestRenderBlogPostListHTML_ColumnClasses_Column2(t *testing.T) {
	post := blogstore.NewPost().
		SetID("1").
		SetTitle("Test")

	postList := []blogstore.PostInterface{post}

	html, err := renderBlogPostListHTML(postList, 1, 10, 1, false, false, false, false, 2, 150)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(html, "col-md-6 col-sm-6") {
		t.Errorf("Expected column class 'col-md-6 col-sm-6' in HTML\nHTML: %s", html)
	}
}

func TestRenderBlogPostListHTML_ColumnClasses_Column3(t *testing.T) {
	post := blogstore.NewPost().
		SetID("1").
		SetTitle("Test")

	postList := []blogstore.PostInterface{post}

	html, err := renderBlogPostListHTML(postList, 1, 10, 1, false, false, false, false, 3, 150)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(html, "col-md-4 col-sm-6") {
		t.Errorf("Expected column class 'col-md-4 col-sm-6' in HTML\nHTML: %s", html)
	}
}

func TestRenderBlogPostListHTML_ColumnClasses_Column4(t *testing.T) {
	post := blogstore.NewPost().
		SetID("1").
		SetTitle("Test")

	postList := []blogstore.PostInterface{post}

	html, err := renderBlogPostListHTML(postList, 1, 10, 1, false, false, false, false, 4, 150)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(html, "col-md-3 col-sm-6") {
		t.Errorf("Expected column class 'col-md-3 col-sm-6' in HTML\nHTML: %s", html)
	}
}

func TestRenderBlogPostListHTML_ColumnClasses_Column6(t *testing.T) {
	post := blogstore.NewPost().
		SetID("1").
		SetTitle("Test")

	postList := []blogstore.PostInterface{post}

	html, err := renderBlogPostListHTML(postList, 1, 10, 1, false, false, false, false, 6, 150)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(html, "col-md-2 col-sm-4") {
		t.Errorf("Expected column class 'col-md-2 col-sm-4' in HTML\nHTML: %s", html)
	}
}

func TestRenderPagination_SinglePageNoPagination(t *testing.T) {
	result := renderPagination(5, 10, 1)

	// If single page, should return empty string
	if result != "" {
		t.Errorf("Expected empty string for single page, got: %s", result)
	}
}

func TestRenderPagination_MultiplePagesPage1(t *testing.T) {
	result := renderPagination(25, 10, 1)

	wantContains := []string{"page-item active", "?page=1\">1</a>", "?page=2\">2</a>", "?page=3\">3</a>", "Previous", "Next"}

	// Check that expected strings are present
	for _, want := range wantContains {
		if !strings.Contains(result, want) {
			t.Errorf("Expected '%s' in pagination output\nGot: %s", want, result)
		}
	}
}

func TestRenderPagination_MultiplePagesPage2(t *testing.T) {
	result := renderPagination(25, 10, 2)

	wantContains := []string{"page-item active", "?page=1\">1</a>", "?page=2\">2</a>", "?page=3\">3</a>", "Previous", "Next"}

	// Check that expected strings are present
	for _, want := range wantContains {
		if !strings.Contains(result, want) {
			t.Errorf("Expected '%s' in pagination output\nGot: %s", want, result)
		}
	}
}

func TestRenderPagination_FirstPagePreviousDisabled(t *testing.T) {
	result := renderPagination(25, 10, 1)

	wantContains := []string{"page-item disabled", "aria-label=\"Previous\""}

	// Check that expected strings are present
	for _, want := range wantContains {
		if !strings.Contains(result, want) {
			t.Errorf("Expected '%s' in pagination output\nGot: %s", want, result)
		}
	}
}

func TestRenderPagination_LastPageNextDisabled(t *testing.T) {
	result := renderPagination(25, 10, 3)

	wantContains := []string{"page-item disabled", "aria-label=\"Next\""}

	// Check that expected strings are present
	for _, want := range wantContains {
		if !strings.Contains(result, want) {
			t.Errorf("Expected '%s' in pagination output\nGot: %s", want, result)
		}
	}
}
