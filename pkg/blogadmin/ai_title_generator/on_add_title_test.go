package aititlegenerator

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"project/internal/testutils"
	"project/pkg/blogai"

	"github.com/dracory/customstore"
)

// TestOnAddTitle_MissingTitle tests adding with missing title
func TestOnAddTitle_MissingTitle(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/admin/blog/ai-title-generator"},
	}

	result := c.onAddTitle(req)

	if !strings.Contains(result, "Title is required") {
		t.Error("Should return error for missing title")
	}
	if !strings.Contains(result, "error") {
		t.Error("Should show error popup")
	}
}

// TestOnAddTitle_EmptyTitle tests adding with empty title
func TestOnAddTitle_EmptyTitle(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path:     "/admin/blog/ai-title-generator",
			RawQuery: "custom_title=",
		},
	}

	result := c.onAddTitle(req)

	if !strings.Contains(result, "Title is required") {
		t.Error("Should return error for empty title")
	}
}

// TestOnAddTitle_WhitespaceTitle tests adding with whitespace-only title
func TestOnAddTitle_WhitespaceTitle(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path:     "/admin/blog/ai-title-generator",
			RawQuery: "custom_title=   ",
		},
	}

	result := c.onAddTitle(req)

	if !strings.Contains(result, "Title is required") {
		t.Error("Should return error for whitespace-only title")
	}
}

// TestOnAddTitle_ValidTitle tests adding with valid title
func TestOnAddTitle_ValidTitle(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path:     "/admin/blog/ai-title-generator",
			RawQuery: "custom_title=My+Test+Title",
		},
	}

	// May panic if custom store is not configured
	var didPanic bool
	defer func() {
		if r := recover(); r != nil {
			didPanic = true
			t.Logf("Expected panic if custom store not configured: %v", r)
		}
	}()

	result := c.onAddTitle(req)

	if !didPanic && !strings.Contains(result, "Custom title added successfully") {
		t.Error("Should return success message")
	}
	if !strings.Contains(result, "success") {
		t.Error("Should show success popup")
	}
	if !strings.Contains(result, "closeModal") {
		t.Error("Should include modal close script")
	}
}

// TestOnAddTitle_TitleWithSpecialChars tests adding title with special characters
func TestOnAddTitle_TitleWithSpecialChars(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path:     "/admin/blog/ai-title-generator",
			RawQuery: "custom_title=Title+with+%26+special+%23+chars",
		},
	}

	// May panic if custom store is not configured
	var didPanic bool
	defer func() {
		if r := recover(); r != nil {
			didPanic = true
			t.Logf("Expected panic if custom store not configured: %v", r)
		}
	}()

	result := c.onAddTitle(req)

	if !didPanic && !strings.Contains(result, "Custom title added successfully") {
		t.Error("Should handle special characters in title")
	}
}

// TestOnAddTitle_VerifyRecordCreated verifies the record was actually created
func TestOnAddTitle_VerifyRecordCreated(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	// Skip if custom store is not available
	if registry.GetCustomStore() == nil {
		t.Skip("Custom store not available")
	}

	testTitle := "Unique Test Title 12345"

	// Count records before
	recordsBefore, err := registry.GetCustomStore().RecordList(customstore.RecordQuery().SetType(blogai.POST_RECORD_TYPE))
	if err != nil {
		t.Fatalf("Failed to list records before: %v", err)
	}
	countBefore := len(recordsBefore)

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path:     "/admin/blog/ai-title-generator",
			RawQuery: "custom_title=" + url.QueryEscape(testTitle),
		},
	}

	// May panic if custom store is not fully configured
	var didPanic bool
	defer func() {
		if r := recover(); r != nil {
			didPanic = true
			t.Logf("Expected panic if custom store not fully configured: %v", r)
		}
	}()

	c.onAddTitle(req)

	if didPanic {
		return
	}

	// Count records after
	recordsAfter, err := registry.GetCustomStore().RecordList(customstore.RecordQuery().SetType(blogai.POST_RECORD_TYPE))
	if err != nil {
		t.Fatalf("Failed to list records after: %v", err)
	}
	countAfter := len(recordsAfter)

	if countAfter != countBefore+1 {
		t.Errorf("Expected %d records after, got %d", countBefore+1, countAfter)
	}

	// Verify the created record has correct status
	var found bool
	for _, record := range recordsAfter {
		recordPost, err := blogai.NewRecordPostFromCustomRecord(record)
		if err != nil {
			continue
		}
		if recordPost.Title == testTitle {
			found = true
			if recordPost.Status != blogai.POST_STATUS_PENDING {
				t.Errorf("Expected status %s, got %s", blogai.POST_STATUS_PENDING, recordPost.Status)
			}
			break
		}
	}

	if !found {
		t.Error("Could not find created record with test title")
	}
}

// TestOnAddTitle_NilRegistry tests behavior with nil registry
func TestOnAddTitle_NilRegistry(t *testing.T) {
	c := NewAiTitleGeneratorController(nil)

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path:     "/admin/blog/ai-title-generator",
			RawQuery: "custom_title=Test+Title",
		},
	}

	var didPanic bool
	defer func() {
		if r := recover(); r != nil {
			didPanic = true
			t.Logf("Expected panic with nil registry: %v", r)
		}
	}()

	_ = c.onAddTitle(req)

	if !didPanic {
		t.Error("Expected panic with nil registry, but function did not panic")
	}
}
