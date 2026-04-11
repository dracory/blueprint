package post_delete

import (
	"strings"
	"testing"
)

func TestModalPostDelete(t *testing.T) {
	// Create test data
	data := postDeleteControllerData{
		postID: "test123",
	}

	// Generate the modal
	modal := modalPostDelete(data)
	html := modal.ToHTML()

	// Verify basic structure
	if !strings.Contains(html, "ModalPostDelete") {
		t.Error("Should contain modal element")
	}
	if !strings.Contains(html, "Are you sure you want to delete this post?") {
		t.Error("Should show confirmation message")
	}
	if !strings.Contains(html, "This action cannot be undone") {
		t.Error("Should show warning")
	}
	if !strings.Contains(html, data.postID) {
		t.Error("Should contain post ID")
	}

	// Verify buttons
	if !strings.Contains(html, "Delete") {
		t.Error("Should have delete button")
	}
	if !strings.Contains(html, "Close") {
		t.Error("Should have close button")
	}

	// Verify form submission
	if !strings.Contains(html, "hx-post=\"") {
		t.Error("Should have post URL")
	}
	if !strings.Contains(html, "post_id") {
		t.Error("Should have post ID field")
	}
}

