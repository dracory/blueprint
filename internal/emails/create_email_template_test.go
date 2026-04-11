package emails

import (
	"strings"
	"testing"

	"project/internal/testutils"
)

func TestCreateEmailTemplate(t *testing.T) {
	// Test with nil registry
	result := CreateEmailTemplate(nil, "Test Title", "<p>Test Content</p>")
	if result == "" {
		t.Error("CreateEmailTemplate() with nil registry should return non-empty string")
	}
	if !strings.Contains(result, "Test Title") {
		t.Error("CreateEmailTemplate() should contain the title")
	}
	if !strings.Contains(result, "<p>Test Content</p>") {
		t.Error("CreateEmailTemplate() should contain the content")
	}

	// Test with valid registry
	registry := testutils.Setup()
	result = CreateEmailTemplate(registry, "Test Title", "<p>Test Content</p>")
	if result == "" {
		t.Error("CreateEmailTemplate() should return non-empty string")
	}
	if !strings.Contains(result, "Test Title") {
		t.Error("CreateEmailTemplate() should contain the title")
	}
	if !strings.Contains(result, "<p>Test Content</p>") {
		t.Error("CreateEmailTemplate() should contain the content")
	}
	if !strings.Contains(result, "Test registry") {
		t.Error("CreateEmailTemplate() should contain the app name from registry")
	}

	// Test with empty title
	result = CreateEmailTemplate(registry, "", "<p>Test Content</p>")
	if result == "" {
		t.Error("CreateEmailTemplate() with empty title should return non-empty string")
	}
	if !strings.Contains(result, "<p>Test Content</p>") {
		t.Error("CreateEmailTemplate() should contain the content even with empty title")
	}

	// Test with empty content
	result = CreateEmailTemplate(registry, "Test Title", "")
	if result == "" {
		t.Error("CreateEmailTemplate() with empty content should return non-empty string")
	}
	if !strings.Contains(result, "Test Title") {
		t.Error("CreateEmailTemplate() should contain the title even with empty content")
	}

	// Test with both empty
	result = CreateEmailTemplate(registry, "", "")
	if result == "" {
		t.Error("CreateEmailTemplate() with both empty should return non-empty string")
	}
}
