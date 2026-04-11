package post_manager

import (
	"strings"
	"testing"

	"github.com/dracory/blogstore"
)

func TestTablePostList(t *testing.T) {
	// Create test data
	data := postManagerControllerData{
		blogList: []blogstore.PostInterface{
			createTestPost("post1", "Test Post 1", blogstore.POST_STATUS_PUBLISHED, "yes"),
			createTestPost("post2", "Test Post 2", blogstore.POST_STATUS_DRAFT, "no"),
		},
		blogCount: 2,
		pageInt:   1,
		perPage:   10,
	}

	// Generate the table
	table := tablePostList(data)
	html := table.ToHTML()

	// Verify basic structure
	if !strings.Contains(html, "<table") {
		t.Error("Should contain table element")
	}
	if !strings.Contains(html, "Test Post 1") {
		t.Error("Should show first post")
	}
	if !strings.Contains(html, "Test Post 2") {
		t.Error("Should show second post")
	}

	// Verify status styling
	if !strings.Contains(html, "text-success") {
		t.Error("Published post should use success text styling")
	}
	if !strings.Contains(html, "text-primary") {
		t.Error("Draft post should use primary text styling")
	}

	// Verify action buttons
	if !strings.Contains(html, "bi-pencil-square") {
		t.Error("Should have edit button")
	}
	if !strings.Contains(html, "bi-trash") {
		t.Error("Should have delete button")
	}
}

func TestSortableColumnLabel(t *testing.T) {
	data := postManagerControllerData{
		sortBy:    "title",
		sortOrder: "asc",
	}

	// Test selected column
	label := sortableColumnLabel(data, "Post", "title")
	html := label.ToHTML()
	if !strings.Contains(html, "&#8593;") {
		t.Error("Should show up arrow for ascending sort")
	}

	// Test non-selected column
	label = sortableColumnLabel(data, "Status", "status")
	html = label.ToHTML()
	if strings.Contains(html, "&#8593;") {
		t.Error("Should not show arrow for non-sorted column")
	}
}

func TestTableFilter(t *testing.T) {
	data := postManagerControllerData{
		status:   blogstore.POST_STATUS_PUBLISHED,
		dateFrom: "2023-01-01",
		dateTo:   "2023-12-31",
		search:   "test",
	}

	filter := tableFilter(data)
	html := filter.ToHTML()

	// Verify filter controls
	if !strings.Contains(html, "FORM_TRANSACTIONS") {
		t.Error("Should contain filter form")
	}
	if !strings.Contains(html, "date_from") {
		t.Error("Should have date from input")
	}
	if !strings.Contains(html, "date_to") {
		t.Error("Should have date to input")
	}
	if !strings.Contains(html, "Published") {
		t.Error("Should show published status selected")
	}
	if !strings.Contains(html, "test") {
		t.Error("Should show search term")
	}
}

func TestTablePagination(t *testing.T) {
	data := postManagerControllerData{
		blogCount: 25,
		pageInt:   2,
		perPage:   10,
	}

	pagination := tablePagination(data, int(data.blogCount), data.pageInt, data.perPage)
	html := pagination.ToHTML()

	// Verify pagination controls
	if !strings.Contains(html, "pagination") {
		t.Error("Should contain pagination")
	}
	if !strings.Contains(html, "page=1") {
		t.Error("Should link to previous page")
	}
	if !strings.Contains(html, "page=2") {
		t.Error("Should show current page")
	}
}

func createTestPost(id, title, status, featured string) blogstore.PostInterface {
	post := blogstore.NewPost()
	post.SetID(id)
	post.SetTitle(title)
	post.SetStatus(status)
	post.SetFeatured(featured)
	return post
}

