package tag_manager

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"
	"project/pkg/blogadmin/shared"

	"github.com/dracory/api"
	"github.com/dracory/blogstore"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/str"
	"github.com/dracory/uid"
)

//go:embed *.html
//go:embed *.js
var tagsFiles embed.FS

type tagManagerController struct {
	registry registry.RegistryInterface
}

func NewTagManagerController(registry registry.RegistryInterface) *tagManagerController {
	return &tagManagerController{registry: registry}
}

func (controller *tagManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := r.URL.Query().Get("action")

	switch action {
	case "load-tags":
		return controller.handleLoadTags(r)
	case "load-tag-posts":
		return controller.handleLoadTagPosts(r)
	case "create-tag":
		return controller.handleCreateTag(w, r)
	case "update-tag":
		return controller.handleUpdateTag(w, r)
	case "delete-tag":
		return controller.handleDeleteTag(w, r)
	default:
		return controller.renderPage(r)
	}
}

func (controller *tagManagerController) renderPage(r *http.Request) string {
	authUser := helpers.GetAuthUser(r)
	if authUser == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), nil, r, "You are not logged in. Please login to continue.", links.Admin().Blog(), 10)
	}

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.Admin().Home()},
		{Name: "Blog", URL: links.Admin().Blog()},
		{Name: "Tags", URL: ""},
	})

	heading := hb.Heading1().HTML("Blog. Tag Manager")

	htmlContent, err := tagsFiles.ReadFile("tags.html")
	if err != nil {
		slog.Error("Failed to read tags HTML template", "error", err)
		return hb.Div().HTML("Error loading tags component").ToHTML()
	}

	jsContent, err := tagsFiles.ReadFile("tags.js")
	if err != nil {
		slog.Error("Failed to read tags JavaScript file", "error", err)
		return hb.Div().HTML("Error loading tags component").ToHTML()
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		const urlTagsLoad = '` + shared.NewLinks("/admin/blog").TagManager(map[string]string{"action": "load-tags"}) + `';
		const urlTagPostsLoad = '` + shared.NewLinks("/admin/blog").TagManager(map[string]string{"action": "load-tag-posts", "tag_id": "TAG_ID_PLACEHOLDER"}) + `';
		const urlTagCreate = '` + shared.NewLinks("/admin/blog").TagManager(map[string]string{"action": "create-tag"}) + `';
		const urlTagUpdate = '` + shared.NewLinks("/admin/blog").TagManager(map[string]string{"action": "update-tag", "tag_id": "TAG_ID_PLACEHOLDER"}) + `';
		const urlTagDelete = '` + shared.NewLinks("/admin/blog").TagManager(map[string]string{"action": "delete-tag"}) + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	vueContainer := hb.Div().
		Child(vueCDN).
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)

	content := hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(vueContainer)

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Blog | Tag Manager",
		Content: content,
		ScriptURLs: []string{
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *tagManagerController) handleLoadTags(r *http.Request) string {
	ctx := r.Context()

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	tagTaxonomy, err := controller.ensureTaxonomy(ctx, blogStore)
	if err != nil {
		return api.Error("Failed to ensure taxonomy: " + err.Error()).ToString()
	}

	terms, err := blogStore.TermList(ctx, blogstore.TermQueryOptions{
		TaxonomyID: tagTaxonomy.GetID(),
		OrderBy:    "name",
		SortOrder:  "asc",
	})
	if err != nil {
		slog.Error("Failed to load tags", "error", err)
		return api.Error("Failed to load tags").ToString()
	}

	tagList := []map[string]any{}
	for _, term := range terms {
		tagList = append(tagList, map[string]any{
			"id":    term.GetID(),
			"name":  term.GetName(),
			"slug":  term.GetSlug(),
			"count": term.GetCount(),
		})
	}

	return api.SuccessWithData("Tags loaded successfully", map[string]any{
		"tags": tagList,
	}).ToString()
}

func (controller *tagManagerController) handleLoadTagPosts(r *http.Request) string {
	ctx := r.Context()

	tagID := r.URL.Query().Get("tag_id")
	if tagID == "" {
		return api.Error("Tag ID is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	// Get the tag to verify it exists
	tag, err := blogStore.TermFindByID(ctx, tagID)
	if err != nil {
		return api.Error("Tag not found").ToString()
	}

	// Get posts associated with this tag (all statuses)
	posts, err := blogStore.PostListByTermID(ctx, tagID, blogstore.PostQueryOptions{
		OrderBy:   "published_at",
		SortOrder: "desc",
		Limit:     100,
	})
	if err != nil {
		slog.Error("Failed to load posts for tag", "error", err, "tag_id", tagID)
		return api.Error("Failed to load posts for tag").ToString()
	}

	postList := []map[string]any{}
	for _, post := range posts {
		postList = append(postList, map[string]any{
			"id":           post.GetID(),
			"title":        post.GetTitle(),
			"slug":         post.GetSlug(),
			"status":       post.GetStatus(),
			"published_at": post.GetPublishedAt(),
		})
	}

	// Use actual posts count (not stored count which may be from import)
	//actualCount := len(postList)

	return api.SuccessWithData("Tag information loaded", map[string]any{
		"tag": map[string]any{
			"id":    tag.GetID(),
			"name":  tag.GetName(),
			"slug":  tag.GetSlug(),
			"count": tag.GetCount(),
		},
		"posts": postList,
	}).ToString()
}

func (controller *tagManagerController) handleCreateTag(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	var reqData struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.Name == "" {
		return api.Error("Tag name is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	tagTaxonomy, err := controller.ensureTaxonomy(ctx, blogStore)
	if err != nil {
		return api.Error("Failed to ensure taxonomy: " + err.Error()).ToString()
	}

	slug := reqData.Slug
	if slug == "" {
		slug = str.Slugify(reqData.Name, '-')
	}

	term := blogstore.NewTerm()
	term.SetID(uid.HumanUid()[:8])
	term.SetName(reqData.Name)
	term.SetSlug(slug)
	term.SetTaxonomyID(tagTaxonomy.GetID())

	if err := blogStore.TermCreate(ctx, term); err != nil {
		slog.Error("Failed to create tag", "error", err)
		return api.Error("Failed to create tag").ToString()
	}

	return api.SuccessWithData("Tag created successfully", map[string]any{
		"id":   term.GetID(),
		"name": term.GetName(),
		"slug": term.GetSlug(),
	}).ToString()
}

func (controller *tagManagerController) handleUpdateTag(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	tagID := r.URL.Query().Get("tag_id")
	if tagID == "" {
		return api.Error("Tag ID is required").ToString()
	}

	var reqData struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.Name == "" {
		return api.Error("Tag name is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	term, err := blogStore.TermFindByID(ctx, tagID)
	if err != nil {
		return api.Error("Tag not found").ToString()
	}

	slug := reqData.Slug
	if slug == "" {
		slug = str.Slugify(reqData.Name, '-')
	}

	term.SetName(reqData.Name)
	term.SetSlug(slug)

	if err := blogStore.TermUpdate(ctx, term); err != nil {
		slog.Error("Failed to update tag", "error", err)
		return api.Error("Failed to update tag").ToString()
	}

	return api.SuccessWithData("Tag updated successfully", map[string]any{
		"id":   term.GetID(),
		"name": term.GetName(),
		"slug": term.GetSlug(),
	}).ToString()
}

func (controller *tagManagerController) handleDeleteTag(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	var reqData struct {
		TagID string `json:"tag_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		reqData.TagID = r.FormValue("tag_id")
	}

	if reqData.TagID == "" {
		return api.Error("Tag ID is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	term, err := blogStore.TermFindByID(ctx, reqData.TagID)
	if err != nil {
		slog.Error("Failed to find tag for delete", "error", err)
		return api.Error("Tag not found").ToString()
	}

	if err := blogStore.TermDelete(ctx, term); err != nil {
		slog.Error("Failed to delete tag", "error", err)
		return api.Error("Failed to delete tag").ToString()
	}

	return api.SuccessWithData("Tag deleted successfully", map[string]any{}).ToString()
}

func (controller *tagManagerController) ensureTaxonomy(ctx context.Context, store blogstore.StoreInterface) (blogstore.TaxonomyInterface, error) {
	tagTaxonomy, err := store.TaxonomyFindBySlug(ctx, blogstore.TAXONOMY_TAG)
	if err != nil || tagTaxonomy == nil {
		controller.registry.GetLogger().Info("Creating tag taxonomy")
		tagTaxonomy = blogstore.NewTaxonomy()
		tagTaxonomy.SetName("Tag")
		tagTaxonomy.SetSlug(blogstore.TAXONOMY_TAG)
		tagTaxonomy.SetDescription("Blog post tags")
		if err := store.TaxonomyCreate(ctx, tagTaxonomy); err != nil {
			return nil, err
		}
	}

	if tagTaxonomy == nil {
		return nil, errors.New("tag taxonomy is nil after ensure")
	}

	return tagTaxonomy, nil
}

// Deprecated: kept for backwards compatibility
type tagManagerControllerData struct {
	page       string
	pageInt    int
	perPage    int
	taxonomyID string
	tagCount   int64
	tagList    []blogstore.TermInterface
}
