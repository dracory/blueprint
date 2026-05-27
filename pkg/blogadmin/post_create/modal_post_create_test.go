package post_create

import (
	"strings"
	"testing"
)

func TestModalPostCreate_WithEmptyTitle(t *testing.T) {
	data := postCreateControllerData{title: ""}
	modal := modalPostCreate(data)
	html := modal.ToHTML()

	// Verify the output contains expected elements
	expected := []string{
		"<input",
		"name=\"post_title\"",
		"value=\"\"",
		"Create & Edit",
	}
	for _, s := range expected {
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
}

func TestModalPostCreate_WithTitle(t *testing.T) {
	data := postCreateControllerData{title: "Test Post"}
	modal := modalPostCreate(data)
	html := modal.ToHTML()

	// Verify the output contains expected elements
	if !strings.Contains(html, "value=\"Test Post\"") {
		t.Error("HTML output should contain value=\"Test Post\"")
	}

	// Verify the modal ID is present
	if !strings.Contains(html, "ModalPostCreate") {
		t.Error("Modal should have correct ID")
	}

	// Verify the close function script is present
	if !strings.Contains(html, "function closeModal") {
		t.Error("Modal should have close function")
	}
}
