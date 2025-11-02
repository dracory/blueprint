package post_delete

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Contains(t, html, "ModalPostDelete", "Should contain modal element")
	assert.Contains(t, html, "Are you sure you want to delete this post?", "Should show confirmation message")
	assert.Contains(t, html, "This action cannot be undone", "Should show warning")
	assert.Contains(t, html, data.postID, "Should contain post ID")

	// Verify buttons
	assert.Contains(t, html, "Delete", "Should have delete button")
	assert.Contains(t, html, "Close", "Should have close button")

	// Verify form submission
	assert.Contains(t, html, "hx-post=\"", "Should have post URL")
	assert.Contains(t, html, "post_id", "Should have post ID field")
}
