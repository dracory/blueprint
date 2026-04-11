package post_create

import (
	"strings"
	"testing"
)

func TestModalPostCreate(t *testing.T) {
	// Test cases
	tests := []struct {
		name     string
		title    string
		contains []string
	}{
		{
			name:  "with empty title",
			title: "",
			contains: []string{
				"<input",
				"name=\"post_title\"",
				"value=\"\"",
				"Create & Edit",
			},
		},
		{
			name:  "with title",
			title: "Test Post",
			contains: []string{
				"value=\"Test Post\"",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate the modal with test data
			data := postCreateControllerData{title: tt.title}
			modal := modalPostCreate(data)
			html := modal.ToHTML()

			// Verify the output contains expected elements
			for _, s := range tt.contains {
				if !strings.Contains(html, s) {
					t.Errorf("HTML output should contain %s", s)
				}
			}

			// Verify the modal ID is present
			if !strings.Contains(html, "ModalPostCreate") {
				t.Error("Modal should have correct ID")
			}

			// Verify the close function script is present
			if !strings.Contains(html, "function closeModal") {
				t.Error("Modal should have close function")
			}
		})
	}
}

