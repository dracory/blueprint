package emails

import (
	"strings"
	"testing"

	"project/internal/testutils"
)

func TestCreateEmailTemplate(t *testing.T) {
	// Test with nil app
	result := CreateEmailTemplate(nil, "Test Title", "<p>Test Content</p>")
	if result == "" {
		t.Error("CreateEmailTemplate() with nil app should return non-empty string")
	}
	if !strings.Contains(result, "Test Title") {
		t.Error("CreateEmailTemplate() should contain the title")
	}
	if !strings.Contains(result, "<p>Test Content</p>") {
		t.Error("CreateEmailTemplate() should contain the content")
	}

	// Test with valid app
	app := testutils.Setup()
	result = CreateEmailTemplate(app, "Test Title", "<p>Test Content</p>")
	if result == "" {
		t.Error("CreateEmailTemplate() should return non-empty string")
	}
	if !strings.Contains(result, "Test Title") {
		t.Error("CreateEmailTemplate() should contain the title")
	}
	if !strings.Contains(result, "<p>Test Content</p>") {
		t.Error("CreateEmailTemplate() should contain the content")
	}
	if !strings.Contains(result, "Test app") {
		t.Error("CreateEmailTemplate() should contain the app name from app")
	}

	// Test with empty title
	result = CreateEmailTemplate(app, "", "<p>Test Content</p>")
	if result == "" {
		t.Error("CreateEmailTemplate() with empty title should return non-empty string")
	}
	if !strings.Contains(result, "<p>Test Content</p>") {
		t.Error("CreateEmailTemplate() should contain the content even with empty title")
	}

	// Test with empty content
	result = CreateEmailTemplate(app, "Test Title", "")
	if result == "" {
		t.Error("CreateEmailTemplate() with empty content should return non-empty string")
	}
	if !strings.Contains(result, "Test Title") {
		t.Error("CreateEmailTemplate() should contain the title even with empty content")
	}

	// Test with both empty
	result = CreateEmailTemplate(app, "", "")
	if result == "" {
		t.Error("CreateEmailTemplate() with both empty should return non-empty string")
	}
}
