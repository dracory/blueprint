package category_manager

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"project/internal/config"
	"project/internal/testutils"
	"testing"

	"github.com/dracory/blogstore"
	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
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
		assert.Contains(t, resp, "Category Manager")
	})

	t.Run("handleLoadCategories", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/blog/categories?action=load-categories", nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		assert.Contains(t, resp, "success")
		assert.Contains(t, resp, "categories")
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
		assert.Contains(t, resp, "success")
		assert.Contains(t, resp, "New Category")

		// Verify it exists in store
		terms, _ := registry.GetBlogStore().TermList(ctx, blogstore.TermQueryOptions{})
		assert.Len(t, terms, 1)
		assert.Equal(t, "New Category", terms[0].GetName())
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
		assert.Contains(t, resp, "success")
		assert.Contains(t, resp, "Updated Category")

		// Verify update
		term, _ := registry.GetBlogStore().TermFindByID(ctx, categoryID)
		assert.Equal(t, "Updated Category", term.GetName())
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
		assert.Contains(t, resp, "success")
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
		assert.Contains(t, resp, "success")

		// Verify deletion
		termsAfter, _ := registry.GetBlogStore().TermList(ctx, blogstore.TermQueryOptions{})
		assert.Len(t, termsAfter, 0)
	})
}
