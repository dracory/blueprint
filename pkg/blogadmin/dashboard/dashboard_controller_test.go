package dashboard

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"

	"project/internal/config"
	"project/internal/registry"
	"project/internal/testutils"

	"github.com/dracory/blogstore"
	"github.com/dracory/test"
)

func TestDashboardController_RequiresAuthentication(t *testing.T) {
	registry := testutils.Setup(
		testutils.WithBlogStore(true),
		testutils.WithCacheStore(true),
		testutils.WithUserStore(true),
	)

	// Test without authentication
	response, responseObj, err := test.CallStringEndpoint(http.MethodGet, NewDashboardController(registry).Handler, test.NewRequestOptions{})
	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if responseObj.StatusCode != http.StatusSeeOther {
		t.Errorf("Should redirect when unauthenticated, got %d", responseObj.StatusCode)
	}
	if !strings.Contains(response, "See Other") {
		t.Error("Should show redirect response")
	}

	// Test with authentication
	user, err := testutils.SeedUser(registry.GetUserStore(), test.USER_01)
	if err != nil {
		t.Errorf("Should create test user: %v", err)
	}

	authResponse, authResponseObj, err := test.CallStringEndpoint(http.MethodGet, NewDashboardController(registry).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})
	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if authResponseObj.StatusCode != http.StatusOK {
		t.Errorf("Should return 200 when authenticated, got %d", authResponseObj.StatusCode)
	}
	if strings.Contains(authResponse, "See Other") {
		t.Error("Should not redirect when authenticated")
	}
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
	if err != nil {
		t.Errorf("Should create test user: %v", err)
	}

	responseHTML, response, err := test.CallStringEndpoint(http.MethodGet, NewDashboardController(registry).Handler, test.NewRequestOptions{
		Context: map[any]any{
			config.AuthenticatedUserContextKey{}: user,
		},
	})

	if err != nil {
		t.Errorf("Handler should not return error: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("Should return 200 status, got %d", response.StatusCode)
	}

	// Verify counts are displayed
	if !strings.Contains(responseHTML, "1") {
		t.Error("Should show post count")
	}
	if !strings.Contains(responseHTML, "Categories") {
		t.Error("Should show Categories tab when taxonomy enabled")
	}
	if !strings.Contains(responseHTML, "Tags") {
		t.Error("Should show Tags tab when taxonomy enabled")
	}
	if !strings.Contains(responseHTML, "Total Posts") {
		t.Error("Should show Total Posts label")
	}
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
	if err != nil {
		t.Errorf("Should create test user: %v", err)
	}

	controller := NewDashboardController(registry)
	req, _ := http.NewRequest(http.MethodGet, "/admin/blog/dashboard", nil)
	req = req.WithContext(context.WithValue(req.Context(), config.AuthenticatedUserContextKey{}, user))

	data, errMsg := controller.prepareData(req)

	if errMsg != "" {
		t.Errorf("Should not return error message, got %s", errMsg)
	}
	if !data.taxonomyEnabled {
		t.Error("Should detect taxonomy as enabled")
	}
	if data.categoryCount != 1 {
		t.Errorf("Should count categories correctly, got %d", data.categoryCount)
	}
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
	if err != nil {
		t.Errorf("Should create test user: %v", err)
	}

	controller := NewDashboardController(registry)
	req, _ := http.NewRequest(http.MethodGet, "/admin/blog/dashboard", nil)
	req = req.WithContext(context.WithValue(req.Context(), config.AuthenticatedUserContextKey{}, user))

	data, errMsg := controller.prepareData(req)

	if errMsg != "" {
		t.Errorf("Should not return error message, got %s", errMsg)
	}
	if data.taxonomyEnabled {
		t.Error("Should detect taxonomy as disabled")
	}
	if !strings.Contains(data.taxonomyErrorMsg, "not available") {
		t.Error("Should set helpful error message")
	}
	if data.categoryCount != 0 {
		t.Errorf("Should have 0 categories when taxonomy disabled, got %d", data.categoryCount)
	}
	if data.tagCount != 0 {
		t.Errorf("Should have 0 tags when taxonomy disabled, got %d", data.tagCount)
	}
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

	if !strings.Contains(html, "Dashboard") {
		t.Error("Should show Dashboard tab")
	}
	if !strings.Contains(html, "Posts") {
		t.Error("Should show Posts tab")
	}
	if !strings.Contains(html, "Categories") {
		t.Error("Should show Categories tab when taxonomy enabled")
	}
	if !strings.Contains(html, "Tags") {
		t.Error("Should show Tags tab when taxonomy enabled")
	}
	if !strings.Contains(html, "5") {
		t.Error("Should show post count badge")
	}
	if !strings.Contains(html, "3") {
		t.Error("Should show category count badge")
	}
	if !strings.Contains(html, "7") {
		t.Error("Should show tag count badge")
	}
}

func TestDashboardController_navTabs_WithTaxonomyDisabled(t *testing.T) {
	data := dashboardControllerData{
		postCount:       5,
		taxonomyEnabled: false,
	}

	controller := &dashboardController{}
	tabs := controller.navTabs(data)
	html := tabs.ToHTML()

	if !strings.Contains(html, "Dashboard") {
		t.Error("Should show Dashboard tab")
	}
	if !strings.Contains(html, "Posts") {
		t.Error("Should show Posts tab")
	}
	if strings.Contains(html, "Categories") {
		t.Error("Should NOT show Categories tab when taxonomy disabled")
	}
	if strings.Contains(html, "Tags") {
		t.Error("Should NOT show Tags tab when taxonomy disabled")
	}
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

	if !strings.Contains(html, "Total Posts") {
		t.Error("Should show Posts card")
	}
	if !strings.Contains(html, "10") {
		t.Error("Should show correct post count")
	}
	if !strings.Contains(html, "Categories") {
		t.Error("Should show Categories card when taxonomy enabled")
	}
	if !strings.Contains(html, "5") {
		t.Error("Should show correct category count")
	}
	if !strings.Contains(html, "Tags") {
		t.Error("Should show Tags card when taxonomy enabled")
	}
	if !strings.Contains(html, "8") {
		t.Error("Should show correct tag count")
	}
}

func TestDashboardController_dashboardCards_WithTaxonomyDisabled(t *testing.T) {
	data := dashboardControllerData{
		postCount:       10,
		taxonomyEnabled: false,
	}

	controller := &dashboardController{}
	cards := controller.dashboardCards(data)
	html := cards.ToHTML()

	if !strings.Contains(html, "Total Posts") {
		t.Error("Should show Posts card")
	}
	if !strings.Contains(html, "10") {
		t.Error("Should show correct post count")
	}
	if strings.Contains(html, "Categories") {
		t.Error("Should NOT show Categories card when taxonomy disabled")
	}
	if strings.Contains(html, "Tags") {
		t.Error("Should NOT show Tags card when taxonomy disabled")
	}
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
