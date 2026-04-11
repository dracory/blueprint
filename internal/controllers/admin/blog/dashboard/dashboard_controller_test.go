package dashboard

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"project/internal/config"
	"project/internal/registry"
	"project/internal/testutils"

	"github.com/dracory/blogstore"
	"github.com/dracory/test"
	"github.com/stretchr/testify/assert"
)

func TestDashboardController_RequiresAuthentication(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Test without authentication
	response, responseObj, err := test.CallStringEndpoint(http.MethodGet, NewDashboardController(registry).Handler, test.NewRequestOptions{})
	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusSeeOther, responseObj.StatusCode, "Should redirect when unauthenticated")
	assert.Contains(t, response, "See Other", "Should show redirect response")

	// Test with authentication
	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	assert.NoError(t, err, "Should create test user")

	authResponse, authResponseObj, err := test.CallStringEndpoint(http.MethodGet, NewDashboardController(registry).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusOK, authResponseObj.StatusCode, "Should return 200 when authenticated")
	assert.NotContains(t, authResponse, "See Other", "Should not redirect when authenticated")
}

func TestDashboardController_WithTaxonomyEnabled(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Create test post
	post := blogstore.NewPost()
	post.SetTitle("Test Post")
	post.SetStatus(blogstore.POST_STATUS_PUBLISHED)
	if err := registry.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		t.Fatalf("failed to create test post: %v", err)
	}

	// Create taxonomy and terms
	categoryTaxonomy := blogstore.NewTaxonomy()
	categoryTaxonomy.SetName("Category")
	categoryTaxonomy.SetSlug(blogstore.TAXONOMY_CATEGORY)
	if err := registry.GetBlogStore().TaxonomyCreate(context.Background(), categoryTaxonomy); err != nil {
		t.Fatalf("failed to create category taxonomy: %v", err)
	}

	tagTaxonomy := blogstore.NewTaxonomy()
	tagTaxonomy.SetName("Tag")
	tagTaxonomy.SetSlug(blogstore.TAXONOMY_TAG)
	if err := registry.GetBlogStore().TaxonomyCreate(context.Background(), tagTaxonomy); err != nil {
		t.Fatalf("failed to create tag taxonomy: %v", err)
	}

	// Create test category
	category := blogstore.NewTerm()
	category.SetName("Test Category")
	category.SetTaxonomyID(categoryTaxonomy.GetID())
	if err := registry.GetBlogStore().TermCreate(context.Background(), category); err != nil {
		t.Fatalf("failed to create test category: %v", err)
	}

	// Create test tag
	tag := blogstore.NewTerm()
	tag.SetName("Test Tag")
	tag.SetTaxonomyID(tagTaxonomy.GetID())
	if err := registry.GetBlogStore().TermCreate(context.Background(), tag); err != nil {
		t.Fatalf("failed to create test tag: %v", err)
	}

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	assert.NoError(t, err, "Should create test user")

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewDashboardController(registry).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	assert.NoError(t, err, "Handler should not return error")
	assert.Equal(t, http.StatusOK, response.StatusCode, "Should return 200 status")

	// Verify counts are displayed
	assert.Contains(t, responseHTML, "1", "Should show post count")
	assert.Contains(t, responseHTML, "Categories", "Should show Categories tab when taxonomy enabled")
	assert.Contains(t, responseHTML, "Tags", "Should show Tags tab when taxonomy enabled")
	assert.Contains(t, responseHTML, "Total Posts", "Should show Total Posts label")
}

func TestDashboardController_prepareData_TaxonomyEnabled(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Create taxonomy and term
	categoryTaxonomy := blogstore.NewTaxonomy()
	categoryTaxonomy.SetName("Category")
	categoryTaxonomy.SetSlug(blogstore.TAXONOMY_CATEGORY)
	if err := registry.GetBlogStore().TaxonomyCreate(context.Background(), categoryTaxonomy); err != nil {
		t.Fatalf("failed to create category taxonomy: %v", err)
	}

	category := blogstore.NewTerm()
	category.SetName("Test Category")
	category.SetTaxonomyID(categoryTaxonomy.GetID())
	if err := registry.GetBlogStore().TermCreate(context.Background(), category); err != nil {
		t.Fatalf("failed to create test category: %v", err)
	}

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	assert.NoError(t, err, "Should create test user")

	controller := NewDashboardController(registry)
	req, _ := http.NewRequest(http.MethodGet, "/admin/blog/dashboard", nil)
	req = req.WithContext(context.WithValue(req.Context(), config.AuthenticatedUserContextKey{}, user))

	data, errMsg := controller.prepareData(req)

	assert.Equal(t, "", errMsg, "Should not return error message")
	assert.True(t, data.taxonomyEnabled, "Should detect taxonomy as enabled")
	assert.Equal(t, int64(1), data.categoryCount, "Should count categories correctly")
}

func TestDashboardController_prepareData_TaxonomyDisabled(t *testing.T) {
	// Create a registry with blog store but taxonomy disabled
	registry := testutils.Setup(
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Manually create blog store without taxonomy enabled
	blogStore, err := createBlogStoreWithoutTaxonomy(registry)
	if err != nil {
		t.Fatalf("failed to create blog store without taxonomy: %v", err)
	}
	registry.SetBlogStore(blogStore)

	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	assert.NoError(t, err, "Should create test user")

	controller := NewDashboardController(registry)
	req, _ := http.NewRequest(http.MethodGet, "/admin/blog/dashboard", nil)
	req = req.WithContext(context.WithValue(req.Context(), config.AuthenticatedUserContextKey{}, user))

	data, errMsg := controller.prepareData(req)

	assert.Equal(t, "", errMsg, "Should not return error message")
	assert.False(t, data.taxonomyEnabled, "Should detect taxonomy as disabled")
	assert.Contains(t, data.taxonomyErrorMsg, "not available", "Should set helpful error message")
	assert.Equal(t, int64(0), data.categoryCount, "Should have 0 categories when taxonomy disabled")
	assert.Equal(t, int64(0), data.tagCount, "Should have 0 tags when taxonomy disabled")
}

func TestDashboardController_navTabs_WithTaxonomyEnabled(t *testing.T) {
	data := dashboardControllerData{
		postCount:       5,
		categoryCount:   3,
		tagCount:        7,
		taxonomyEnabled: true,
	}

	controller := &dashboardController{}
	tabs := controller.navTabs(data)
	html := tabs.ToHTML()

	assert.Contains(t, html, "Dashboard", "Should show Dashboard tab")
	assert.Contains(t, html, "Posts", "Should show Posts tab")
	assert.Contains(t, html, "Categories", "Should show Categories tab when taxonomy enabled")
	assert.Contains(t, html, "Tags", "Should show Tags tab when taxonomy enabled")
	assert.Contains(t, html, "5", "Should show post count badge")
	assert.Contains(t, html, "3", "Should show category count badge")
	assert.Contains(t, html, "7", "Should show tag count badge")
}

func TestDashboardController_navTabs_WithTaxonomyDisabled(t *testing.T) {
	data := dashboardControllerData{
		postCount:       5,
		taxonomyEnabled: false,
	}

	controller := &dashboardController{}
	tabs := controller.navTabs(data)
	html := tabs.ToHTML()

	assert.Contains(t, html, "Dashboard", "Should show Dashboard tab")
	assert.Contains(t, html, "Posts", "Should show Posts tab")
	assert.NotContains(t, html, "Categories", "Should NOT show Categories tab when taxonomy disabled")
	assert.NotContains(t, html, "Tags", "Should NOT show Tags tab when taxonomy disabled")
}

func TestDashboardController_dashboardCards_WithTaxonomyEnabled(t *testing.T) {
	data := dashboardControllerData{
		postCount:       10,
		categoryCount:   5,
		tagCount:        8,
		taxonomyEnabled: true,
	}

	controller := &dashboardController{}
	cards := controller.dashboardCards(data)
	html := cards.ToHTML()

	assert.Contains(t, html, "Total Posts", "Should show Posts card")
	assert.Contains(t, html, "10", "Should show correct post count")
	assert.Contains(t, html, "Categories", "Should show Categories card when taxonomy enabled")
	assert.Contains(t, html, "5", "Should show correct category count")
	assert.Contains(t, html, "Tags", "Should show Tags card when taxonomy enabled")
	assert.Contains(t, html, "8", "Should show correct tag count")
}

func TestDashboardController_dashboardCards_WithTaxonomyDisabled(t *testing.T) {
	data := dashboardControllerData{
		postCount:       10,
		taxonomyEnabled: false,
	}

	controller := &dashboardController{}
	cards := controller.dashboardCards(data)
	html := cards.ToHTML()

	assert.Contains(t, html, "Total Posts", "Should show Posts card")
	assert.Contains(t, html, "10", "Should show correct post count")
	assert.NotContains(t, html, "Categories", "Should NOT show Categories card when taxonomy disabled")
	assert.NotContains(t, html, "Tags", "Should NOT show Tags card when taxonomy disabled")
}

// createBlogStoreWithoutTaxonomy creates a blog store with taxonomy disabled
func createBlogStoreWithoutTaxonomy(r registry.RegistryInterface) (blogstore.StoreInterface, error) {
	if r.GetDatabase() == nil {
		return nil, errors.New("database is not initialized")
	}

	st, err := blogstore.NewStore(blogstore.NewStoreOptions{
		DB:                 r.GetDatabase(),
		PostTableName:      "snv_blogs_post_test",
		TaxonomyEnabled:    false,
		VersioningEnabled:  false,
		AutomigrateEnabled: true,
	})

	if err != nil {
		return nil, err
	}

	return st, nil
}
