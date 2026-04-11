package aipostgenerator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"project/internal/testutils"
	"project/pkg/blogai"
)

func TestAiPostGeneratorController_onGeneratePost_MissingRecordID(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	// Test with empty record_post_id - should return error
	req := httptest.NewRequest(http.MethodPost, "/admin/blog/ai-post-generator", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	result := controller.onGeneratePost(req)

	// Should return error popup
	if result == "" {
		t.Error("onGeneratePost with empty record_id should return error HTML")
	}
}

func TestAiPostGeneratorController_onGeneratePost_WithRecordID(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	// Test with record_post_id but invalid store
	formData := url.Values{}
	formData.Set("record_post_id", "test-id-123")

	req := httptest.NewRequest(http.MethodPost, "/admin/blog/ai-post-generator", strings.NewReader(formData.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// This will fail because custom store is not fully configured
	result := controller.onGeneratePost(req)
	_ = result
}

func TestAiPostGeneratorController_stepHandlerGetPostDetails_MissingRecordID(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data := map[string]any{
		// Missing record_post_id
	}

	_, _, err := controller.stepHandlerGetPostDetails(ctx, data)
	if err == nil {
		t.Error("stepHandlerGetPostDetails should error with missing record_post_id")
	}
}

func TestAiPostGeneratorController_stepHandlerGetPostDetails_EmptyRecordID(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data := map[string]any{
		"record_post_id": "",
	}

	_, _, err := controller.stepHandlerGetPostDetails(ctx, data)
	if err == nil {
		t.Error("stepHandlerGetPostDetails should error with empty record_post_id")
	}
}

func TestAiPostGeneratorController_stepHandlerGetPostDetails_NilStore(t *testing.T) {
	// Create controller with registry but no custom store
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data := map[string]any{
		"record_post_id": "test-id",
	}

	_, _, err := controller.stepHandlerGetPostDetails(ctx, data)
	// Will error because custom store is nil in test setup
	// Error is expected, just verify it doesn't panic
	if err == nil {
		t.Error("stepHandlerGetPostDetails should error with nil store")
	}
}

func TestAiPostGeneratorController_stepHandlerGeneratePost_MissingRecordID(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data := map[string]any{
		// Missing record_post_id
	}

	_, _, err := controller.stepHandlerGeneratePost(ctx, data)
	if err == nil {
		t.Error("stepHandlerGeneratePost should error with missing record_post_id")
	}
}

func TestAiPostGeneratorController_stepHandlerGeneratePost_MissingTitle(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data := map[string]any{
		"record_post_id": "test-id",
		// Missing post_title
	}

	_, _, err := controller.stepHandlerGeneratePost(ctx, data)
	if err == nil {
		t.Error("stepHandlerGeneratePost should error with missing post_title")
	}
}

func TestAiPostGeneratorController_stepHandlerGeneratePost_EmptyTitle(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data := map[string]any{
		"record_post_id": "test-id",
		"post_title":     "",
	}

	_, _, err := controller.stepHandlerGeneratePost(ctx, data)
	if err == nil {
		t.Error("stepHandlerGeneratePost should error with empty post_title")
	}
}

func TestAiPostGeneratorController_stepHandlerSavePost_MissingRecordID(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data := map[string]any{
		// Missing record_post_id
	}

	_, _, err := controller.stepHandlerSavePost(ctx, data)
	if err == nil {
		t.Error("stepHandlerSavePost should error with missing record_post_id")
	}
}

func TestAiPostGeneratorController_stepHandlerSavePost_MissingBlogPost(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data := map[string]any{
		"record_post_id": "test-id",
		// Missing blogai_post
	}

	_, _, err := controller.stepHandlerSavePost(ctx, data)
	if err == nil {
		t.Error("stepHandlerSavePost should error with missing blogai_post")
	}
}

func TestAiPostGeneratorController_stepHandlerSavePost_NilStore(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	data := map[string]any{
		"record_post_id": "test-id",
		"blogai_post": blogai.RecordPost{
			ID:    "test-id",
			Title: "Test Title",
		},
	}

	_, _, err := controller.stepHandlerSavePost(ctx, data)
	// Will error because custom store is nil
	// Error is expected, just verify it doesn't panic
	if err == nil {
		t.Error("stepHandlerSavePost should error with nil store")
	}
}

func TestAiPostGeneratorController_tableApprovedTitles_Empty(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	data := pageData{
		ApprovedBlogAiPosts: []blogai.RecordPost{},
	}

	result := controller.tableApprovedTitles(data)
	if result == nil {
		t.Error("tableApprovedTitles with empty data should not return nil")
	}
}

func TestAiPostGeneratorController_tableApprovedTitles_WithPosts(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	data := pageData{
		ApprovedBlogAiPosts: []blogai.RecordPost{
			{
				ID:     "post-1",
				Title:  "Approved Post",
				Status: blogai.POST_STATUS_APPROVED,
			},
			{
				ID:     "post-2",
				Title:  "Draft Post",
				Status: blogai.POST_STATUS_DRAFT,
			},
		},
	}

	result := controller.tableApprovedTitles(data)
	if result == nil {
		t.Error("tableApprovedTitles with posts should not return nil")
	}
}

func TestAiPostGeneratorController_tableApprovedTitles_NilStatus(t *testing.T) {
	registry := testutils.Setup()
	controller := NewAiPostGeneratorController(registry)

	data := pageData{
		ApprovedBlogAiPosts: []blogai.RecordPost{
			{
				ID:     "post-1",
				Title:  "Post with nil status",
				Status: "",
			},
		},
	}

	result := controller.tableApprovedTitles(data)
	if result == nil {
		t.Error("tableApprovedTitles with nil status should not return nil")
	}
}
