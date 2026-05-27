package tag_manager

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"project/internal/config"
	"project/internal/testutils"
	"strings"
	"testing"

	"github.com/dracory/blogstore"
	"github.com/dracory/test"
)

func TestTagManagerController_Functional(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, _ := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	controller := NewTagManagerController(registry)

	// Context with auth user
	ctx := context.WithValue(context.Background(), config.AuthenticatedUserContextKey{}, user)

	t.Run("renderPage", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/blog/tags", nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "Tag Manager") {
			t.Error("expected Tag Manager in response")
		}
	})

	t.Run("handleLoadTags", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/blog/tags?action=load-tags", nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}
		if !strings.Contains(resp, "tags") {
			t.Error("expected tags in response")
		}
	})

	t.Run("handleCreateTag", func(t *testing.T) {
		tagData := map[string]string{
			"name": "New Tag",
			"slug": "new-tag",
		}
		body, _ := json.Marshal(tagData)
		req := httptest.NewRequest(http.MethodPost, "/admin/blog/tags?action=create-tag", bytes.NewBuffer(body)).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}
		if !strings.Contains(resp, "New Tag") {
			t.Error("expected New Tag in response")
		}

		// Verify it exists in store
		terms, _ := registry.GetBlogStore().TermList(ctx, blogstore.TermQueryOptions{})
		if len(terms) != 1 {
			t.Errorf("expected 1 term, got %d", len(terms))
		}
		if terms[0].GetName() != "New Tag" {
			t.Errorf("expected New Tag, got %s", terms[0].GetName())
		}
	})

	t.Run("handleLoadTagPosts", func(t *testing.T) {
		terms, _ := registry.GetBlogStore().TermList(ctx, blogstore.TermQueryOptions{})
		tagID := terms[0].GetID()

		req := httptest.NewRequest(http.MethodGet, "/admin/blog/tags?action=load-tag-posts&tag_id="+tagID, nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}
		if !strings.Contains(resp, "posts") {
			t.Error("expected posts in response")
		}
	})

	t.Run("handleUpdateTag", func(t *testing.T) {
		terms, _ := registry.GetBlogStore().TermList(ctx, blogstore.TermQueryOptions{})
		tagID := terms[0].GetID()

		updateData := map[string]string{
			"name": "Updated Tag",
		}
		body, _ := json.Marshal(updateData)
		req := httptest.NewRequest(http.MethodPost, "/admin/blog/tags?action=update-tag&tag_id="+tagID, bytes.NewBuffer(body)).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}

		// Verify update
		term, _ := registry.GetBlogStore().TermFindByID(ctx, tagID)
		if term.GetName() != "Updated Tag" {
			t.Errorf("expected Updated Tag, got %s", term.GetName())
		}
	})

	t.Run("handleDeleteTag", func(t *testing.T) {
		terms, _ := registry.GetBlogStore().TermList(ctx, blogstore.TermQueryOptions{})
		tagID := terms[0].GetID()

		deleteData := map[string]string{
			"tag_id": tagID,
		}
		body, _ := json.Marshal(deleteData)
		req := httptest.NewRequest(http.MethodPost, "/admin/blog/tags?action=delete-tag", bytes.NewBuffer(body)).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}

		// Verify deletion
		termsAfter, _ := registry.GetBlogStore().TermList(ctx, blogstore.TermQueryOptions{})
		if len(termsAfter) != 0 {
			t.Errorf("expected 0 terms after deletion, got %d", len(termsAfter))
		}
	})
}
