package post_update_v1

import (
	"testing"

	"github.com/dracory/blogstore"
	"github.com/stretchr/testify/assert"
)

func TestFormPostUpdate(t *testing.T) {
	// Test cases for different views
	tests := []struct {
		name     string
		view     string
		contains []string
	}{
		{
			name: "DETAILS view",
			view: VIEW_DETAILS,
			contains: []string{
				"name=\"post_status\"",
				"name=\"post_image_url\"",
				"name=\"post_featured\"",
			},
		},
		{
			name: "CONTENT view",
			view: VIEW_CONTENT,
			contains: []string{
				"name=\"post_title\"",
				"name=\"post_summary\"",
				"name=\"post_content\"",
			},
		},
		{
			name: "SEO view",
			view: VIEW_SEO,
			contains: []string{
				"name=\"post_meta_description\"",
				"name=\"post_meta_keywords\"",
				"name=\"post_meta_robots\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test data
			data := postUpdateControllerData{
				view:        tt.view,
				postID:      "test123",
				formTitle:   "Test Title",
				formContent: "Test Content",
				formStatus:  blogstore.POST_STATUS_DRAFT,
				formEditor:  blogstore.POST_EDITOR_TEXTAREA,
				post:        blogstore.NewPost(),
			}

			// Generate the form
			form := formPostUpdate(data)
			html := form.ToHTML()

			// Verify the output contains expected elements
			for _, s := range tt.contains {
				assert.Contains(t, html, s, "Form should contain "+s)
			}

			// Verify the form ID is present
			assert.Contains(t, html, "FormPostUpdate", "Form should have correct ID")
		})
	}
}

func TestFormPostUpdate_ErrorHandling(t *testing.T) {
	// Test with form error message
	t.Run("Form error message", func(t *testing.T) {
		data := postUpdateControllerData{
			view:             VIEW_DETAILS,
			postID:           "test123",
			formErrorMessage: "Test error message",
		}

		form := formPostUpdate(data)
		html := form.ToHTML()

		assert.Contains(t, html, "Swal.fire", "Should show SweetAlert")
		assert.Contains(t, html, "Test error message", "Should show custom error message")
		assert.Contains(t, html, "\"icon\":\"error\"", "Should show error icon")
	})

	// Test with form success message
	t.Run("Form success message", func(t *testing.T) {
		data := postUpdateControllerData{
			view:               VIEW_DETAILS,
			postID:             "test123",
			formSuccessMessage: "Test success message",
		}

		form := formPostUpdate(data)
		html := form.ToHTML()

		assert.Contains(t, html, "Swal.fire", "Should show SweetAlert")
		assert.Contains(t, html, "Test success message", "Should show custom success message")
		assert.Contains(t, html, "\"icon\":\"success\"", "Should show success icon")
	})
}
