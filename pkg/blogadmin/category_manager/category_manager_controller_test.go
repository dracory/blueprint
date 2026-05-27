package category_manager

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

func TestCategoryManagerController_Functional(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, _ := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	controller := NewCategoryManagerController(registry)

	// Context with auth user
	ctx := context.WithValue(context.Background(), config.AuthenticatedUserContextKey{}, user)

	t.Run("renderPage", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/blog/categories", nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "Category Manager") {
			t.Error("expected Category Manager in response")
		}
	})

	t.Run("handleLoadCategories", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/blog/categories?action=load-categories", nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}
		if !strings.Contains(resp, "categories") {
			t.Error("expected categories in response")
		}
	})

	t.Run("handleCreateCategory", func(t *testing.T) {
		catData := map[string]string{
			"name":        "New Category",
			"slug":        "new-category",
			"description": "Test Description",
		}
		body, _ := json.Marshal(catData)
		req := httptest.NewRequest(http.MethodPost, "/admin/blog/categories?action=create-category", bytes.NewBuffer(body)).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}
		if !strings.Contains(resp, "New Category") {
			t.Error("expected New Category in response")
		}

		// Verify it exists in store
		terms, _ := registry.GetBlogStore().TermList(ctx, blogstore.TermQueryOptions{})
		if len(terms) != 1 {
			t.Errorf("expected 1 term, got %d", len(terms))
		}
		if terms[0].GetName() != "New Category" {
			t.Errorf("expected New Category, got %s", terms[0].GetName())
		}
	})

	t.Run("handleUpdateCategory", func(t *testing.T) {
		// First get the category ID
		terms, _ := registry.GetBlogStore().TermList(ctx, blogstore.TermQueryOptions{})
		categoryID := terms[0].GetID()

		updateData := map[string]string{
			"name": "Updated Category",
		}
		body, _ := json.Marshal(updateData)
		req := httptest.NewRequest(http.MethodPost, "/admin/blog/categories?action=update-category&category_id="+categoryID, bytes.NewBuffer(body)).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}
		if !strings.Contains(resp, "Updated Category") {
			t.Error("expected Updated Category in response")
		}

		// Verify update
		term, _ := registry.GetBlogStore().TermFindByID(ctx, categoryID)
		if term.GetName() != "Updated Category" {
			t.Errorf("expected Updated Category, got %s", term.GetName())
		}
	})

	t.Run("handleReorderCategories", func(t *testing.T) {
		terms, _ := registry.GetBlogStore().TermList(ctx, blogstore.TermQueryOptions{})
		categoryID := terms[0].GetID()

		reorderData := map[string][]string{
			"category_ids": {categoryID},
		}
		body, _ := json.Marshal(reorderData)
		req := httptest.NewRequest(http.MethodPost, "/admin/blog/categories?action=reorder-categories", bytes.NewBuffer(body)).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}
	})

	t.Run("handleDeleteCategory", func(t *testing.T) {
		terms, _ := registry.GetBlogStore().TermList(ctx, blogstore.TermQueryOptions{})
		categoryID := terms[0].GetID()

		deleteData := map[string]string{
			"category_id": categoryID,
		}
		body, _ := json.Marshal(deleteData)
		req := httptest.NewRequest(http.MethodPost, "/admin/blog/categories?action=delete-category", bytes.NewBuffer(body)).WithContext(ctx)
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
