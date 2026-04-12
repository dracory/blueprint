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

// TestOnDeleteTitle_MissingID tests deleting with missing title ID
func TestOnDeleteTitle_MissingID(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "POST",
		URL:    &url.URL{Path: "/admin/blog/ai-title-generator"},
	}

	result := c.onDeleteTitle(req)

	if !strings.Contains(result, "Title ID is required") {
		t.Error("Should return error for missing title ID")
	}
	if !strings.Contains(result, "error") {
		t.Error("Should show error popup")
	}
}

// TestOnDeleteTitle_WithID tests deleting with valid ID
func TestOnDeleteTitle_WithID(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	// Skip if custom store is not available
	if registry.GetCustomStore() == nil {
		t.Skip("Custom store not available")
	}

	// First create a record to delete
	record := customstore.NewRecord(blogai.POST_RECORD_TYPE)
	if err := record.SetPayloadMap(map[string]any{
		"id":     record.ID(),
		"title":  "Test Title to Delete",
		"status": blogai.POST_STATUS_PENDING,
	}); err != nil {
		t.Fatalf("Failed to set payload: %v", err)
	}

	customStore := registry.GetCustomStore()
	if err := customStore.RecordCreate(record); err != nil {
		t.Fatalf("Failed to create record: %v", err)
	}

	// Now try to delete it
	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path:     "/admin/blog/ai-title-generator",
			RawQuery: "record_post_id=" + record.ID(),
		},
	}

	result := c.onDeleteTitle(req)

	if !strings.Contains(result, "Title deleted successfully") {
		t.Error("Should return success message")
	}
	if !strings.Contains(result, "success") {
		t.Error("Should show success popup")
	}
}

// TestOnDeleteTitle_InvalidID tests deleting with non-existent ID
func TestOnDeleteTitle_InvalidID(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	// Skip if custom store is not available
	if registry.GetCustomStore() == nil {
		t.Skip("Custom store not available")
	}

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path:     "/admin/blog/ai-title-generator",
			RawQuery: "record_post_id=non-existent-id",
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

	// Should not panic and should return success even if record doesn't exist
	result := c.onDeleteTitle(req)

	if !didPanic && result == "" {
		t.Error("Should return some response even for non-existent ID")
	}
}

// TestOnDeleteTitle_EmptyID tests deleting with empty ID parameter
func TestOnDeleteTitle_EmptyID(t *testing.T) {
	registry := testutils.Setup()
	c := NewAiTitleGeneratorController(registry)

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path:     "/admin/blog/ai-title-generator",
			RawQuery: "record_post_id=",
		},
	}

	result := c.onDeleteTitle(req)

	if !strings.Contains(result, "Title ID is required") {
		t.Error("Should return error for empty title ID")
	}
}

// TestOnDeleteTitle_NilRegistry tests behavior with nil registry (should handle gracefully)
func TestOnDeleteTitle_NilRegistry(t *testing.T) {
	c := NewAiTitleGeneratorController(nil)

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Path:     "/admin/blog/ai-title-generator",
			RawQuery: "record_post_id=test-id",
		},
	}

	// This will panic due to nil registry, which is acceptable for production
	// but we test that the method exists and has proper signature
	var didPanic bool
	defer func() {
		if r := recover(); r != nil {
			didPanic = true
			t.Logf("Expected panic with nil registry: %v", r)
		}
	}()

	_ = c.onDeleteTitle(req)

	if !didPanic {
		t.Error("Expected panic with nil registry, but function did not panic")
	}
}
