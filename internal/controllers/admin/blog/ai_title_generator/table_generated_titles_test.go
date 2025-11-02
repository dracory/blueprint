package aititlegenerator

import (
	"strings"
	"testing"

	"project/pkg/blogai"

	"github.com/stretchr/testify/assert"
)

func TestTableGeneratedTitles_EmptyData(t *testing.T) {
	data := pageData{
		ExistingPostRecords: []blogai.RecordPost{},
	}

	result := tableGeneratedTitles(data)
	html := result.ToHTML()

	assert.Contains(t, html, "No titles generated yet.", "Should show empty message")
	assert.Contains(t, html, "text-muted", "Should have muted text class")
}

func TestTableGeneratedTitles_WithPublishedRecords(t *testing.T) {
	data := pageData{
		ExistingPostRecords: []blogai.RecordPost{
			createTestRecordPost("record1", "Pending Title", blogai.POST_STATUS_PENDING),
			createTestRecordPost("record2", "Published Title", blogai.POST_STATUS_PUBLISHED),
		},
	}

	result := tableGeneratedTitles(data)
	html := result.ToHTML()

	assert.Contains(t, html, "Generated Titles", "Should include generated titles section")
	assert.Contains(t, html, "Published Titles (Reference)", "Should include published titles section")

	generatedIdx := strings.Index(html, "Generated Titles")
	publishedIdx := strings.Index(html, "Published Titles (Reference)")

	assert.NotEqual(t, -1, generatedIdx, "Generated titles heading should exist")
	assert.NotEqual(t, -1, publishedIdx, "Published titles heading should exist")
	assert.Greater(t, publishedIdx, generatedIdx, "Published titles section should appear after generated titles")

	publishedSection := html[publishedIdx:]

	assert.Contains(t, html, "Published Title", "Should render published title content")
	assert.Contains(t, html, "bg-primary", "Should render published badge with primary color")
	assert.Contains(t, html, "Delete", "Published titles should still allow delete action")
	assert.NotContains(t, publishedSection, "action=approve_title", "Published titles should not include approve button")
	assert.NotContains(t, publishedSection, "action=reject_title", "Published titles should not include reject button")
}

func TestTableGeneratedTitles_PublishedOnly(t *testing.T) {
	data := pageData{
		ExistingPostRecords: []blogai.RecordPost{
			createTestRecordPost("record1", "Published Title", blogai.POST_STATUS_PUBLISHED),
		},
	}

	result := tableGeneratedTitles(data)
	html := result.ToHTML()

	assert.NotContains(t, html, "Generated Titles", "Should not include generated section when only published titles exist")
	assert.Contains(t, html, "Published Titles (Reference)", "Should render published section heading")
	assert.Contains(t, html, "Published Title", "Should render published title content")
	assert.Contains(t, html, "bg-primary", "Should render published badge with primary color")
}

func TestTableGeneratedTitles_WithPendingRecords(t *testing.T) {
	data := pageData{
		ExistingPostRecords: []blogai.RecordPost{
			createTestRecordPost("record1", "Test Title 1", blogai.POST_STATUS_PENDING),
			createTestRecordPost("record2", "Test Title 2", blogai.POST_STATUS_PENDING),
		},
	}

	result := tableGeneratedTitles(data)
	html := result.ToHTML()

	assert.Contains(t, html, "Generated Titles", "Should render generated titles section heading")

	// Verify table structure
	assert.Contains(t, html, "<table", "Should contain table element")
	assert.Contains(t, html, "table-striped", "Should have table classes")
	assert.Contains(t, html, "table-hover", "Should have table classes")

	// Verify headers
	assert.Contains(t, html, "Title", "Should have Title header")
	assert.Contains(t, html, "Status", "Should have Status header")
	assert.Contains(t, html, "Actions", "Should have Actions header")

	// Verify content
	assert.Contains(t, html, "Test Title 1", "Should show first title")
	assert.Contains(t, html, "Test Title 2", "Should show second title")

	// Verify status badges
	assert.Contains(t, html, "bg-warning", "Should have warning badge for pending")
	assert.Contains(t, html, "pending", "Should show pending status text")

	// Verify action buttons for pending status
	assert.Contains(t, html, "btn-success", "Should have approve button")
	assert.Contains(t, html, "btn-warning", "Should have reject button")
	assert.Contains(t, html, "btn-outline-danger", "Should have delete button")
	assert.Contains(t, html, "Approve", "Should have approve button text")
	assert.Contains(t, html, "Reject", "Should have reject button text")
	assert.Contains(t, html, "Delete", "Should have delete button text")

	// Verify HTMX attributes
	assert.Contains(t, html, "hx-post", "Should have HTMX post")
	assert.Contains(t, html, "hx-target=\"body\"", "Should target body")
	assert.Contains(t, html, "hx-swap=\"beforeend\"", "Should swap beforeend")
	assert.Contains(t, html, "hx-indicator", "Should have HTMX indicator")
	assert.Contains(t, html, "action=approve_title", "Should have approve action")
	assert.Contains(t, html, "action=reject_title", "Should have reject action")
	assert.Contains(t, html, "record_post_id=record1", "Should have first record ID")
	assert.Contains(t, html, "record_post_id=record2", "Should have second record ID")

	// Verify spinners
	assert.Contains(t, html, "spinner-border", "Should have spinner indicators")
}

func TestTableGeneratedTitles_WithApprovedRecords(t *testing.T) {
	data := pageData{
		ExistingPostRecords: []blogai.RecordPost{
			createTestRecordPost("record1", "Approved Title", blogai.POST_STATUS_APPROVED),
		},
	}

	result := tableGeneratedTitles(data)
	html := result.ToHTML()

	assert.Contains(t, html, "Generated Titles", "Should render generated titles section heading")

	// Verify status badge
	assert.Contains(t, html, "bg-success", "Should have success badge for approved")
	assert.Contains(t, html, "approved", "Should show approved status text")

	// Verify generate post button
	assert.Contains(t, html, "btn-primary", "Should have primary button")
	assert.Contains(t, html, "Generate Post", "Should have generate post button text")
	assert.Contains(t, html, "action=generate_post", "Should have generate post action")

	// Verify delete button is present
	assert.Contains(t, html, "btn-outline-danger", "Should have delete button")
	assert.Contains(t, html, "Delete", "Should have delete button text")
	assert.Contains(t, html, "action=delete_title", "Should have delete action")

	// Should not have approve/reject buttons
	assert.NotContains(t, html, ">Approve<", "Should not have approve button for approved records")
	assert.NotContains(t, html, ">Reject<", "Should not have reject button for approved records")
}

func TestTableGeneratedTitles_WithMixedStatuses(t *testing.T) {
	data := pageData{
		ExistingPostRecords: []blogai.RecordPost{
			createTestRecordPost("record1", "Pending Title", blogai.POST_STATUS_PENDING),
			createTestRecordPost("record2", "Approved Title", blogai.POST_STATUS_APPROVED),
			createTestRecordPost("record3", "Rejected Title", blogai.POST_STATUS_REJECTED),
		},
	}

	result := tableGeneratedTitles(data)
	html := result.ToHTML()

	assert.Contains(t, html, "Generated Titles", "Should render generated titles section heading")

	// Verify all status badges
	assert.Contains(t, html, "bg-warning", "Should have warning badge for pending")
	assert.Contains(t, html, "bg-success", "Should have success badge for approved")
	assert.Contains(t, html, "bg-danger", "Should have danger badge for rejected")

	// Verify action buttons
	assert.Contains(t, html, "Approve", "Should have approve button")
	assert.Contains(t, html, "Reject", "Should have reject button")
	assert.Contains(t, html, "Generate Post", "Should have generate post button")
	assert.Contains(t, html, "Delete", "Should have delete button")
}

func TestTableGeneratedTitles_WithEmptyStatus(t *testing.T) {
	data := pageData{
		ExistingPostRecords: []blogai.RecordPost{
			createTestRecordPost("record1", "Title with Empty Status", ""),
		},
	}

	result := tableGeneratedTitles(data)
	html := result.ToHTML()

	assert.Contains(t, html, "Generated Titles", "Should render generated titles section heading")

	// Should show N/A for empty status
	assert.Contains(t, html, "N/A", "Should show N/A for empty status")
	assert.Contains(t, html, "bg-secondary", "Should have secondary badge for default/empty status")
}

func TestGetStatusBadgeClass(t *testing.T) {
	tests := []struct {
		status   string
		expected string
	}{
		{blogai.POST_STATUS_PENDING, "bg-warning"},
		{blogai.POST_STATUS_APPROVED, "bg-success"},
		{blogai.POST_STATUS_REJECTED, "bg-danger"},
		{blogai.POST_STATUS_DRAFT, "bg-info"},
		{blogai.POST_STATUS_PUBLISHED, "bg-primary"},
		{"unknown_status", "bg-secondary"},
		{"", "bg-secondary"},
	}

	for _, test := range tests {
		t.Run("Status_"+test.status, func(t *testing.T) {
			result := getStatusBadgeClass(test.status)
			assert.Equal(t, test.expected, result, "Should return correct badge class for status: %s", test.status)
		})
	}
}

func TestTableGeneratedTitles_Structure(t *testing.T) {
	data := pageData{
		ExistingPostRecords: []blogai.RecordPost{
			createTestRecordPost("record1", "Test Title", blogai.POST_STATUS_PENDING),
		},
	}

	result := tableGeneratedTitles(data)
	html := result.ToHTML()

	assert.Contains(t, html, "Generated Titles", "Should render generated titles section heading")
	assert.NotContains(t, html, "Published Titles (Reference)", "Should not render published section when there are no published titles")

	// Verify table structure elements
	assert.Contains(t, html, "<thead>", "Should have table head")
	assert.Contains(t, html, "<tbody>", "Should have table body")
	assert.Contains(t, html, "<tr>", "Should have table rows")
	assert.Contains(t, html, "<td>", "Should have table cells")
	assert.Contains(t, html, "<th", "Should have table headers")

	// Verify badge structure
	assert.Contains(t, html, "badge rounded-pill", "Should have badge classes")
	assert.Contains(t, html, "px-3", "Should have badge padding")

	// Verify button structure
	assert.Contains(t, html, "btn btn-success", "Should have success button classes")
	assert.Contains(t, html, "btn btn-warning", "Should have warning button classes")
	assert.Contains(t, html, "btn btn-outline-danger", "Should have outline danger button classes")
}

// Helper function to create test RecordPost
func createTestRecordPost(id, title, status string) blogai.RecordPost {
	return blogai.RecordPost{
		ID:     id,
		Title:  title,
		Status: status,
	}
}
