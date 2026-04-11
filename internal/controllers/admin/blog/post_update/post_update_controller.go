package post_update

import (
	"embed"
	"encoding/json"
	"log/slog"
	"net/http"

	"project/internal/controllers/admin/blog/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/api"
	"github.com/dracory/blogstore"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dracory/req"
)

//go:embed *.html
//go:embed *.js
var postCategoriesFiles embed.FS

type postUpdateController struct {
	registry registry.RegistryInterface
}

func NewPostUpdateController(registry registry.RegistryInterface) *postUpdateController {
	return &postUpdateController{registry: registry}
}

func (controller *postUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")
	postID := req.GetStringTrimmed(r, "post_id")
	view := req.GetStringTrimmedOr(r, "view", "content")

	// Handle API actions for categories and tags
	if action != "" {
		switch action {
		case "load-categories":
			return controller.handleLoadCategories(r)
		case "add-category":
			return controller.handleAddCategory(w, r)
		case "remove-category":
			return controller.handleRemoveCategory(w, r)
		case "load-tags":
			return controller.handleLoadTags(r)
		case "add-tag":
			return controller.handleAddTag(w, r)
		case "remove-tag":
			return controller.handleRemoveTag(w, r)
		}
	}

	if postID == "" {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "Post ID is required", links.Admin().Blog(), 10)
	}

	post, err := controller.registry.GetBlogStore().PostFindByID(r.Context(), postID)
	if err != nil {
		controller.registry.GetLogger().Error(
			"Error. postUpdateController: PostFindByID",
			slog.String("error", err.Error()),
			slog.String("post_id", postID),
		)
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "Post not found", links.Admin().Blog(), 10)
	}

	if post == nil {
		controller.registry.GetLogger().Warn(
			"Warning. postUpdateController: PostFindByID",
			slog.String("error", "Post not found"),
			slog.String("post_id", postID),
		)
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "Post not found", links.Admin().Blog(), 10)
	}

	pageContent := controller.page(r, post, view)

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Edit Post | Blog",
		Content: pageContent,
		ScriptURLs: []string{
			cdn.Jquery_3_7_1(),
			cdn.TrumbowygJs_2_27_3(),
			cdn.Sweetalert2_10(),
			cdn.JqueryUiJs_1_13_1(), // BlockArea requires jQuery UI
			links.Website().Resource(`/js/blockarea_v0200.js`, map[string]string{}), // BlockArea
		},
		Scripts: []string{
			liveflux.Script().ToHTML(),
		},
		StyleURLs: []string{
			cdn.TrumbowygCss_2_27_3(),
			cdn.JqueryUiCss_1_13_1(), // BlockArea styles
		},
	}).ToHTML()
}

func (controller *postUpdateController) page(r *http.Request, post blogstore.PostInterface, view string) hb.TagInterface {
	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{
			Name: "Home",
			URL:  links.Admin().Home(),
		},
		{
			Name: "Blog",
			URL:  links.Admin().Blog(),
		},
		{
			Name: "Post Manager",
			URL:  shared.NewLinks().PostManager(),
		},
		{
			Name: "Edit Post",
			URL:  shared.NewLinks().PostUpdate(map[string]string{"post_id": post.GetID()}),
		},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(links.Admin().Blog())

	heading := hb.Heading1().
		HTML("Edit Post").
		Child(buttonCancel)

	tabs := bs.NavTabs().
		Class("mb-3").
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "details", "active").
				Href(shared.NewLinks().PostUpdate(map[string]string{
					"post_id": post.GetID(),
					"view":    "details",
				})).
				HTML("Details"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "content", "active").
				Href(shared.NewLinks().PostUpdate(map[string]string{
					"post_id": post.GetID(),
					"view":    "content",
				})).
				HTML("Content"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "categories", "active").
				Href(shared.NewLinks().PostUpdate(map[string]string{
					"post_id": post.GetID(),
					"view":    "categories",
				})).
				HTML("Categories"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "tags", "active").
				Href(shared.NewLinks().PostUpdate(map[string]string{
					"post_id": post.GetID(),
					"view":    "tags",
				})).
				HTML("Tags"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "seo", "active").
				Href(shared.NewLinks().PostUpdate(map[string]string{
					"post_id": post.GetID(),
					"view":    "seo",
				})).
				HTML("SEO"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "versions", "active").
				Href(shared.NewLinks().PostUpdate(map[string]string{
					"post_id": post.GetID(),
					"view":    "versions",
				})).
				HTML("Versions")))

	postTitle := hb.Heading2().
		Class("mb-3").
		HTML("Post: ").
		HTML(post.GetTitle())

	var body hb.TagInterface

	switch view {
	case "details":
		component := NewPostDetailsComponent(controller.registry)
		body = liveflux.Placeholder(component, map[string]string{
			"post_id": post.GetID(),
		})
	case "content":
		component := NewPostContentComponent(controller.registry)
		body = liveflux.Placeholder(component, map[string]string{
			"post_id": post.GetID(),
		})
	case "categories":
		body = controller.renderCategoriesView(r, post)
	case "tags":
		body = controller.renderTagsView(r, post)
	case "seo":
		component := NewPostSEOComponent(controller.registry)
		body = liveflux.Placeholder(component, map[string]string{
			"post_id": post.GetID(),
		})
	case "versions":
		component := NewPostVersioningComponent(controller.registry)
		body = liveflux.Placeholder(component, map[string]string{
			"post_id": post.GetID(),
		})
	default:
		body = hb.Div().Text("Not implemented yet")
	}

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header").
				Child(hb.Heading4().
					HTMLIf(view == "details", "Post Details").
					HTMLIf(view == "content", "Post Contents").
					HTMLIf(view == "categories", "Post Categories").
					HTMLIf(view == "tags", "Post Tags").
					HTMLIf(view == "seo", "Post SEO").
					HTMLIf(view == "versions", "Post Versions").
					Style("margin-bottom:0;display:inline-block;")),
		).
		Child(
			hb.Div().
				Class("card-body").
				Child(body),
		)

	return hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(postTitle).
		Child(tabs).
		Child(card)
}

func (controller *postUpdateController) renderCategoriesView(r *http.Request, post blogstore.PostInterface) hb.TagInterface {
	htmlContent, err := postCategoriesFiles.ReadFile("post_categories.html")
	if err != nil {
		slog.Error("Failed to read post categories HTML template", "error", err)
		return hb.Div().HTML("Error loading categories component")
	}

	jsContent, err := postCategoriesFiles.ReadFile("post_categories.js")
	if err != nil {
		slog.Error("Failed to read post categories JavaScript file", "error", err)
		return hb.Div().HTML("Error loading categories component")
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		const postID = '` + post.GetID() + `';
		const urlCategoriesLoad = '` + shared.NewLinks().PostUpdate(map[string]string{"post_id": post.GetID(), "action": "load-categories"}) + `';
		const urlCategoryAdd = '` + shared.NewLinks().PostUpdate(map[string]string{"post_id": post.GetID(), "action": "add-category"}) + `';
		const urlCategoryRemove = '` + shared.NewLinks().PostUpdate(map[string]string{"post_id": post.GetID(), "action": "remove-category"}) + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	vueContainer := hb.Div().
		Child(vueCDN).
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)

	return vueContainer
}

func (controller *postUpdateController) handleLoadCategories(r *http.Request) string {
	ctx := r.Context()

	postID := req.GetStringTrimmed(r, "post_id")
	if postID == "" {
		return api.Error("Post ID is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	// Get category taxonomy
	categoryTaxonomy, err := blogStore.TaxonomyFindBySlug(ctx, blogstore.TAXONOMY_CATEGORY)
	if err != nil || categoryTaxonomy == nil {
		return api.Error("Category taxonomy not found").ToString()
	}

	// Get all available categories
	allCategories, err := blogStore.TermList(ctx, blogstore.TermQueryOptions{
		TaxonomyID: categoryTaxonomy.GetID(),
		OrderBy:    "sequence",
		SortOrder:  "asc",
	})
	if err != nil {
		slog.Error("Failed to load categories", "error", err)
		return api.Error("Failed to load categories").ToString()
	}

	// Get categories assigned to this post
	assignedCategoryIDs := make(map[string]bool)
	postCategories, err := blogStore.TermListByPostID(ctx, postID, blogstore.TAXONOMY_CATEGORY)
	if err == nil {
		for _, category := range postCategories {
			assignedCategoryIDs[category.GetID()] = true
		}
	}

	categoryList := []map[string]any{}
	for _, category := range allCategories {
		categoryList = append(categoryList, map[string]any{
			"id":          category.GetID(),
			"name":        category.GetName(),
			"slug":        category.GetSlug(),
			"description": category.GetDescription(),
			"assigned":    assignedCategoryIDs[category.GetID()],
		})
	}

	return api.SuccessWithData("Categories loaded successfully", map[string]any{
		"categories": categoryList,
	}).ToString()
}

func (controller *postUpdateController) handleAddCategory(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	postID := req.GetStringTrimmed(r, "post_id")

	var reqData struct {
		CategoryID string `json:"category_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.CategoryID == "" {
		return api.Error("Category ID is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	// Add category to post using PostAddTerm
	if err := blogStore.PostAddTerm(ctx, postID, reqData.CategoryID); err != nil {
		slog.Error("Failed to add category to post", "error", err)
		return api.Error("Failed to add category to post").ToString()
	}

	return api.Success("Category added to post successfully").ToString()
}

func (controller *postUpdateController) handleRemoveCategory(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	postID := req.GetStringTrimmed(r, "post_id")

	var reqData struct {
		CategoryID string `json:"category_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.CategoryID == "" {
		return api.Error("Category ID is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	// Remove category from post using PostRemoveTerm
	if err := blogStore.PostRemoveTerm(ctx, postID, reqData.CategoryID); err != nil {
		slog.Error("Failed to remove category from post", "error", err)
		return api.Error("Failed to remove category from post").ToString()
	}

	return api.Success("Category removed from post successfully").ToString()
}

func (controller *postUpdateController) renderTagsView(r *http.Request, post blogstore.PostInterface) hb.TagInterface {
	htmlContent, err := postCategoriesFiles.ReadFile("post_tags.html")
	if err != nil {
		slog.Error("Failed to read post tags HTML template", "error", err)
		return hb.Div().HTML("Error loading tags component")
	}

	jsContent, err := postCategoriesFiles.ReadFile("post_tags.js")
	if err != nil {
		slog.Error("Failed to read post tags JavaScript file", "error", err)
		return hb.Div().HTML("Error loading tags component")
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		const postID = '` + post.GetID() + `';
		const urlTagsLoad = '` + shared.NewLinks().PostUpdate(map[string]string{"post_id": post.GetID(), "action": "load-tags"}) + `';
		const urlTagAdd = '` + shared.NewLinks().PostUpdate(map[string]string{"post_id": post.GetID(), "action": "add-tag"}) + `';
		const urlTagRemove = '` + shared.NewLinks().PostUpdate(map[string]string{"post_id": post.GetID(), "action": "remove-tag"}) + `';
	`)

	htmlTemplate := hb.Wrap().HTML(string(htmlContent))
	componentScript := hb.Script(string(jsContent))

	vueContainer := hb.Div().
		Child(vueCDN).
		Child(htmlTemplate).
		Child(initScript).
		Child(componentScript)

	return vueContainer
}

func (controller *postUpdateController) handleLoadTags(r *http.Request) string {
	ctx := r.Context()

	postID := req.GetStringTrimmed(r, "post_id")
	if postID == "" {
		return api.Error("Post ID is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	// Get tag taxonomy
	tagTaxonomy, err := blogStore.TaxonomyFindBySlug(ctx, blogstore.TAXONOMY_TAG)
	if err != nil || tagTaxonomy == nil {
		return api.Error("Tag taxonomy not found").ToString()
	}

	// Get all available tags
	allTags, err := blogStore.TermList(ctx, blogstore.TermQueryOptions{
		TaxonomyID: tagTaxonomy.GetID(),
		OrderBy:    "name",
		SortOrder:  "asc",
	})
	if err != nil {
		slog.Error("Failed to load tags", "error", err)
		return api.Error("Failed to load tags").ToString()
	}

	// Get tags assigned to this post
	assignedTagIDs := make(map[string]bool)
	postTags, err := blogStore.TermListByPostID(ctx, postID, blogstore.TAXONOMY_TAG)
	if err == nil {
		for _, tag := range postTags {
			assignedTagIDs[tag.GetID()] = true
		}
	}

	tagList := []map[string]any{}
	for _, tag := range allTags {
		tagList = append(tagList, map[string]any{
			"id":       tag.GetID(),
			"name":     tag.GetName(),
			"slug":     tag.GetSlug(),
			"assigned": assignedTagIDs[tag.GetID()],
		})
	}

	return api.SuccessWithData("Tags loaded successfully", map[string]any{
		"tags": tagList,
	}).ToString()
}

func (controller *postUpdateController) handleAddTag(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	postID := req.GetStringTrimmed(r, "post_id")

	var reqData struct {
		TagID string `json:"tag_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.TagID == "" {
		return api.Error("Tag ID is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	// Add tag to post using PostAddTerm
	if err := blogStore.PostAddTerm(ctx, postID, reqData.TagID); err != nil {
		slog.Error("Failed to add tag to post", "error", err)
		return api.Error("Failed to add tag to post").ToString()
	}

	return api.Success("Tag added to post successfully").ToString()
}

func (controller *postUpdateController) handleRemoveTag(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	postID := req.GetStringTrimmed(r, "post_id")

	var reqData struct {
		TagID string `json:"tag_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.TagID == "" {
		return api.Error("Tag ID is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	// Remove tag from post using PostRemoveTerm
	if err := blogStore.PostRemoveTerm(ctx, postID, reqData.TagID); err != nil {
		slog.Error("Failed to remove tag from post", "error", err)
		return api.Error("Failed to remove tag from post").ToString()
	}

	return api.Success("Tag removed from post successfully").ToString()
}
