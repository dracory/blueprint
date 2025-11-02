package post_manager

import (
	"testing"

	"github.com/dracory/blogstore"
	"github.com/stretchr/testify/assert"
)

func TestTablePostList(t *testing.T) {
	// Create test data
	data := postManagerControllerData{
		blogList: []blogstore.Post{
			*createTestPost("post1", "Test Post 1", blogstore.POST_STATUS_PUBLISHED, "yes"),
			*createTestPost("post2", "Test Post 2", blogstore.POST_STATUS_DRAFT, "no"),
		},
		blogCount: 2,
		pageInt:   1,
		perPage:   10,
	}

	// Generate the table
	table := tablePostList(data)
	html := table.ToHTML()

	// Verify basic structure
	assert.Contains(t, html, "<table", "Should contain table element")
	assert.Contains(t, html, "Test Post 1", "Should show first post")
	assert.Contains(t, html, "Test Post 2", "Should show second post")

	// Verify status styling
	assert.Contains(t, html, "text-success", "Published post should use success text styling")
	assert.Contains(t, html, "text-primary", "Draft post should use primary text styling")

	// Verify action buttons
	assert.Contains(t, html, "bi-pencil-square", "Should have edit button")
	assert.Contains(t, html, "bi-trash", "Should have delete button")
}

func TestSortableColumnLabel(t *testing.T) {
	data := postManagerControllerData{
		sortBy:    "title",
		sortOrder: "asc",
	}

	// Test selected column
	label := sortableColumnLabel(data, "Post", "title")
	html := label.ToHTML()
	assert.Contains(t, html, "&#8593;", "Should show up arrow for ascending sort")

	// Test non-selected column
	label = sortableColumnLabel(data, "Status", "status")
	html = label.ToHTML()
	assert.NotContains(t, html, "&#8593;", "Should not show arrow for non-sorted column")
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
	assert.Contains(t, html, "FORM_TRANSACTIONS", "Should contain filter form")
	assert.Contains(t, html, "date_from", "Should have date from input")
	assert.Contains(t, html, "date_to", "Should have date to input")
	assert.Contains(t, html, "Published", "Should show published status selected")
	assert.Contains(t, html, "test", "Should show search term")
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
	assert.Contains(t, html, "pagination", "Should contain pagination")
	assert.Contains(t, html, "page=1", "Should link to previous page")
	assert.Contains(t, html, "page=2", "Should show current page")
}

func createTestPost(id, title, status, featured string) *blogstore.Post {
	post := blogstore.NewPost()
	post.SetID(id)
	post.SetTitle(title)
	post.SetStatus(status)
	post.SetFeatured(featured)
	return post
}
