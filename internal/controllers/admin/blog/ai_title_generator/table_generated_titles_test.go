package aititlegenerator

import (
	"strings"
	"testing"

	"project/pkg/blogai"
)

func TestTableGeneratedTitles_EmptyData(t *testing.T) {
	data := pageData{
		ExistingPostRecords: []blogai.RecordPost{},
	}

	result := tableGeneratedTitles(data)
	html := result.ToHTML()

	if !strings.Contains(html, "No titles generated yet.") {
		t.Error("Should show empty message")
	}
	if !strings.Contains(html, "text-muted") {
		t.Error("Should have muted text class")
	}
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

	if !strings.Contains(html, "Generated Titles") {
		t.Error("Should include generated titles section")
	}
	if !strings.Contains(html, "Published Titles (Reference)") {
		t.Error("Should include published titles section")
	}

	generatedIdx := strings.Index(html, "Generated Titles")
	publishedIdx := strings.Index(html, "Published Titles (Reference)")

	if generatedIdx == -1 {
		t.Error("Generated titles heading should exist")
	}
	if publishedIdx == -1 {
		t.Error("Published titles heading should exist")
	}
	if !(publishedIdx > generatedIdx) {
		t.Error("Published titles section should appear after generated titles")
	}

	publishedSection := html[publishedIdx:]

	if !strings.Contains(html, "Published Title") {
		t.Error("Should render published title content")
	}
	if !strings.Contains(html, "bg-primary") {
		t.Error("Should render published badge with primary color")
	}
	if !strings.Contains(html, "Delete") {
		t.Error("Published titles should still allow delete action")
	}
	if strings.Contains(publishedSection, "action=approve_title") {
		t.Error("Published titles should not include approve button")
	}
	if strings.Contains(publishedSection, "action=reject_title") {
		t.Error("Published titles should not include reject button")
	}
}

func TestTableGeneratedTitles_PublishedOnly(t *testing.T) {
	data := pageData{
		ExistingPostRecords: []blogai.RecordPost{
			createTestRecordPost("record1", "Published Title", blogai.POST_STATUS_PUBLISHED),
		},
	}

	result := tableGeneratedTitles(data)
	html := result.ToHTML()

	if strings.Contains(html, "Generated Titles") {
		t.Error("Should not include generated section when only published titles exist")
	}
	if !strings.Contains(html, "Published Titles (Reference)") {
		t.Error("Should render published section heading")
	}
	if !strings.Contains(html, "Published Title") {
		t.Error("Should render published title content")
	}
	if !strings.Contains(html, "bg-primary") {
		t.Error("Should render published badge with primary color")
	}
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

	if !strings.Contains(html, "Generated Titles") {
		t.Error("Should render generated titles section heading")
	}

	// Verify table structure
	if !strings.Contains(html, "<table") {
		t.Error("Should contain table element")
	}
	if !strings.Contains(html, "table-striped") {
		t.Error("Should have table classes")
	}
	if !strings.Contains(html, "table-hover") {
		t.Error("Should have table classes")
	}

	// Verify headers
	if !strings.Contains(html, "Title") {
		t.Error("Should have Title header")
	}
	if !strings.Contains(html, "Status") {
		t.Error("Should have Status header")
	}
	if !strings.Contains(html, "Actions") {
		t.Error("Should have Actions header")
	}

	// Verify content
	if !strings.Contains(html, "Test Title 1") {
		t.Error("Should show first title")
	}
	if !strings.Contains(html, "Test Title 2") {
		t.Error("Should show second title")
	}

	// Verify status badges
	if !strings.Contains(html, "bg-warning") {
		t.Error("Should have warning badge for pending")
	}
	if !strings.Contains(html, "pending") {
		t.Error("Should show pending status text")
	}

	// Verify action buttons for pending status
	if !strings.Contains(html, "btn-success") {
		t.Error("Should have approve button")
	}
	if !strings.Contains(html, "btn-warning") {
		t.Error("Should have reject button")
	}
	if !strings.Contains(html, "btn-outline-danger") {
		t.Error("Should have delete button")
	}
	if !strings.Contains(html, "Approve") {
		t.Error("Should have approve button text")
	}
	if !strings.Contains(html, "Reject") {
		t.Error("Should have reject button text")
	}
	if !strings.Contains(html, "Delete") {
		t.Error("Should have delete button text")
	}

	// Verify HTMX attributes
	if !strings.Contains(html, "hx-post") {
		t.Error("Should have HTMX post")
	}
	if !strings.Contains(html, "hx-target=\"body\"") {
		t.Error("Should target body")
	}
	if !strings.Contains(html, "hx-swap=\"beforeend\"") {
		t.Error("Should swap beforeend")
	}
	if !strings.Contains(html, "hx-indicator") {
		t.Error("Should have HTMX indicator")
	}
	if !strings.Contains(html, "action=approve_title") {
		t.Error("Should have approve action")
	}
	if !strings.Contains(html, "action=reject_title") {
		t.Error("Should have reject action")
	}
	if !strings.Contains(html, "record_post_id=record1") {
		t.Error("Should have first record ID")
	}
	if !strings.Contains(html, "record_post_id=record2") {
		t.Error("Should have second record ID")
	}

	// Verify spinners
	if !strings.Contains(html, "spinner-border") {
		t.Error("Should have spinner indicators")
	}
}

func TestTableGeneratedTitles_WithApprovedRecords(t *testing.T) {
	data := pageData{
		ExistingPostRecords: []blogai.RecordPost{
			createTestRecordPost("record1", "Approved Title", blogai.POST_STATUS_APPROVED),
		},
	}

	result := tableGeneratedTitles(data)
	html := result.ToHTML()

	if !strings.Contains(html, "Generated Titles") {
		t.Error("Should render generated titles section heading")
	}

	// Verify status badge
	if !strings.Contains(html, "bg-success") {
		t.Error("Should have success badge for approved")
	}
	if !strings.Contains(html, "approved") {
		t.Error("Should show approved status text")
	}

	// Verify generate post button
	if !strings.Contains(html, "btn-primary") {
		t.Error("Should have primary button")
	}
	if !strings.Contains(html, "Generate Post") {
		t.Error("Should have generate post button text")
	}
	if !strings.Contains(html, "action=generate_post") {
		t.Error("Should have generate post action")
	}

	// Verify delete button is present
	if !strings.Contains(html, "btn-outline-danger") {
		t.Error("Should have delete button")
	}
	if !strings.Contains(html, "Delete") {
		t.Error("Should have delete button text")
	}
	if !strings.Contains(html, "action=delete_title") {
		t.Error("Should have delete action")
	}

	// Should not have approve/reject buttons
	if strings.Contains(html, ">Approve<") {
		t.Error("Should not have approve button for approved records")
	}
	if strings.Contains(html, ">Reject<") {
		t.Error("Should not have reject button for approved records")
	}
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

	if !strings.Contains(html, "Generated Titles") {
		t.Error("Should render generated titles section heading")
	}

	// Verify all status badges
	if !strings.Contains(html, "bg-warning") {
		t.Error("Should have warning badge for pending")
	}
	if !strings.Contains(html, "bg-success") {
		t.Error("Should have success badge for approved")
	}
	if !strings.Contains(html, "bg-danger") {
		t.Error("Should have danger badge for rejected")
	}

	// Verify action buttons
	if !strings.Contains(html, "Approve") {
		t.Error("Should have approve button")
	}
	if !strings.Contains(html, "Reject") {
		t.Error("Should have reject button")
	}
	if !strings.Contains(html, "Generate Post") {
		t.Error("Should have generate post button")
	}
	if !strings.Contains(html, "Delete") {
		t.Error("Should have delete button")
	}
}

func TestTableGeneratedTitles_WithEmptyStatus(t *testing.T) {
	data := pageData{
		ExistingPostRecords: []blogai.RecordPost{
			createTestRecordPost("record1", "Title with Empty Status", ""),
		},
	}

	result := tableGeneratedTitles(data)
	html := result.ToHTML()

	if !strings.Contains(html, "Generated Titles") {
		t.Error("Should render generated titles section heading")
	}

	// Should show N/A for empty status
	if !strings.Contains(html, "N/A") {
		t.Error("Should show N/A for empty status")
	}
	if !strings.Contains(html, "bg-secondary") {
		t.Error("Should have secondary badge for default/empty status")
	}
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
			if result != test.expected {
				t.Errorf("Should return correct badge class for status %s: expected %s, got %s", test.status, test.expected, result)
			}
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

	if !strings.Contains(html, "Generated Titles") {
		t.Error("Should render generated titles section heading")
	}
	if strings.Contains(html, "Published Titles (Reference)") {
		t.Error("Should not render published section when there are no published titles")
	}

	// Verify table structure elements
	if !strings.Contains(html, "<thead>") {
		t.Error("Should have table head")
	}
	if !strings.Contains(html, "<tbody>") {
		t.Error("Should have table body")
	}
	if !strings.Contains(html, "<tr>") {
		t.Error("Should have table rows")
	}
	if !strings.Contains(html, "<td>") {
		t.Error("Should have table cells")
	}
	if !strings.Contains(html, "<th") {
		t.Error("Should have table headers")
	}

	// Verify badge structure
	if !strings.Contains(html, "badge rounded-pill") {
		t.Error("Should have badge classes")
	}
	if !strings.Contains(html, "px-3") {
		t.Error("Should have badge padding")
	}

	// Verify button structure
	if !strings.Contains(html, "btn btn-success") {
		t.Error("Should have success button classes")
	}
	if !strings.Contains(html, "btn btn-warning") {
		t.Error("Should have warning button classes")
	}
	if !strings.Contains(html, "btn btn-outline-danger") {
		t.Error("Should have outline danger button classes")
	}
}

// Helper function to create test RecordPost
func createTestRecordPost(id, title, status string) blogai.RecordPost {
	return blogai.RecordPost{
		ID:     id,
		Title:  title,
		Status: status,
	}
}
