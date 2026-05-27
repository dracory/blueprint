package post_update

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"project/internal/config"
	"project/internal/testutils"

	"github.com/dracory/blogstore"
	"github.com/dracory/test"
)

func TestPostUpdateController_Functional(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	user, _ := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	controller := NewPostUpdateController(registry)

	// Context with auth user
	ctx := context.WithValue(context.Background(), config.AuthenticatedUserContextKey{}, user)

	// Create test post
	post := blogstore.NewPost()
	post.SetTitle("Test Post")
	post.SetContent("Test Content")
	post.SetStatus(blogstore.POST_STATUS_DRAFT)
	registry.GetBlogStore().PostCreate(ctx, post)
	postID := post.GetID()

	t.Run("renderDetailsView", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/blog/post/update?post_id="+postID+"&view=details", nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "id=\"post-details-app\"") {
			t.Error("expected post-details-app in response")
		}
	})

	t.Run("handleLoadDetails", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/blog/post/update?action=load-details&post_id="+postID, nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}
		if !strings.Contains(resp, "draft") {
			t.Error("expected draft in response")
		}
	})

	t.Run("handleSaveDetails", func(t *testing.T) {
		saveData := map[string]any{
			"post_status": blogstore.POST_STATUS_PUBLISHED,
			"post_editor": "blockeditor",
		}
		body, _ := json.Marshal(saveData)
		req := httptest.NewRequest(http.MethodPost, "/admin/blog/post/update?action=save-details&post_id="+postID, bytes.NewBuffer(body)).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}

		// Verify update
		p, _ := registry.GetBlogStore().PostFindByID(ctx, postID)
		if p.GetStatus() != blogstore.POST_STATUS_PUBLISHED {
			t.Errorf("expected status %s, got %s", blogstore.POST_STATUS_PUBLISHED, p.GetStatus())
		}
	})

	t.Run("handleLoadCategories", func(t *testing.T) {
		// Ensure taxonomy exists
		tax := blogstore.NewTaxonomy()
		tax.SetName("Category")
		tax.SetSlug(blogstore.TAXONOMY_CATEGORY)
		registry.GetBlogStore().TaxonomyCreate(ctx, tax)

		req := httptest.NewRequest(http.MethodGet, "/admin/blog/post/update?action=load-categories&post_id="+postID, nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}
	})

	t.Run("handleAddCategory", func(t *testing.T) {
		tax, _ := registry.GetBlogStore().TaxonomyFindBySlug(ctx, blogstore.TAXONOMY_CATEGORY)
		term := blogstore.NewTerm()
		term.SetName("Cat1")
		term.SetTaxonomyID(tax.GetID())
		registry.GetBlogStore().TermCreate(ctx, term)

		catData := map[string]string{"category_id": term.GetID()}
		body, _ := json.Marshal(catData)
		req := httptest.NewRequest(http.MethodPost, "/admin/blog/post/update?action=add-category&post_id="+postID, bytes.NewBuffer(body)).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}
	})

	t.Run("handleRemoveCategory", func(t *testing.T) {
		terms, _ := registry.GetBlogStore().TermListByPostID(ctx, postID, blogstore.TAXONOMY_CATEGORY)
		catID := terms[0].GetID()

		catData := map[string]string{"category_id": catID}
		body, _ := json.Marshal(catData)
		req := httptest.NewRequest(http.MethodPost, "/admin/blog/post/update?action=remove-category&post_id="+postID, bytes.NewBuffer(body)).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}
	})

	t.Run("handleLoadTags", func(t *testing.T) {
		// Ensure taxonomy exists
		tax := blogstore.NewTaxonomy()
		tax.SetName("Tag")
		tax.SetSlug(blogstore.TAXONOMY_TAG)
		registry.GetBlogStore().TaxonomyCreate(ctx, tax)

		req := httptest.NewRequest(http.MethodGet, "/admin/blog/post/update?action=load-tags&post_id="+postID, nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "success") {
			t.Error("expected success in response")
		}
	})

	t.Run("renderCategoriesView", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/blog/post/update?post_id="+postID+"&view=categories", nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "id=\"post-categories-app\"") {
			t.Error("expected post-categories-app in response")
		}
	})

	t.Run("renderTagsView", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/blog/post/update?post_id="+postID+"&view=tags", nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "id=\"post-tags-app\"") {
			t.Error("expected post-tags-app in response")
		}
	})

	t.Run("renderSEOView", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/admin/blog/post/update?post_id="+postID+"&view=seo", nil).WithContext(ctx)
		resp := controller.Handler(httptest.NewRecorder(), req)
		if !strings.Contains(resp, "id=\"post-seo-app\"") {
			t.Error("expected post-seo-app in response")
		}
	})
}
