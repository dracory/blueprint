package post_update

import (
	"context"
	"embed"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"
	"project/pkg/blogadmin/shared"
	"project/pkg/blogai"

	"github.com/dracory/api"
	"github.com/dracory/blogstore"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/dracory/versionstore"
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
			return controller.handleLoadCategories(w, r)
		case "add-category":
			return controller.handleAddCategory(w, r)
		case "remove-category":
			return controller.handleRemoveCategory(w, r)
		case "load-tags":
			return controller.handleLoadTags(w, r)
		case "add-tag":
			return controller.handleAddTag(w, r)
		case "remove-tag":
			return controller.handleRemoveTag(w, r)
		// Details component actions
		case "load-details":
			return controller.handleLoadDetails(w, r)
		case "save-details":
			return controller.handleSaveDetails(w, r)
		case "regenerate-image":
			return controller.handleRegenerateImage(w, r)
		// Content component actions
		case "load-content":
			return controller.handleLoadContent(w, r)
		case "save-content":
			return controller.handleSaveContent(w, r)
		case "blockeditor-handle":
			return controller.handleBlockEditorHandle(w, r)
		// SEO component actions
		case "load-seo":
			return controller.handleLoadSEO(w, r)
		case "save-seo":
			return controller.handleSaveSEO(w, r)
		// Versioning component actions
		case "load-versions":
			return controller.handleLoadVersions(w, r)
		case "load-version-detail":
			return controller.handleLoadVersionDetail(w, r)
		case "restore-version-attributes":
			return controller.handleRestoreVersionAttributes(w, r)
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
			"https://cdn.jsdelivr.net/npm/summernote@0.8.18/dist/summernote-lite.min.js",   // Summernote
			"https://cdn.jsdelivr.net/npm/easymde/dist/easymde.min.js",                     // EasyMDE
			"https://cdn.jsdelivr.net/npm/codemirror@5.65.5/lib/codemirror.min.js",         // CodeMirror core
			"https://cdn.jsdelivr.net/npm/codemirror@5.65.5/mode/markdown/markdown.min.js", // CodeMirror markdown mode
			"https://cdn.jsdelivr.net/npm/codemirror@5.65.5/mode/xml/xml.min.js",           // CodeMirror XML mode
			cdn.Sweetalert2_10(),
			cdn.JqueryUiJs_1_13_1(), // BlockArea requires jQuery UI
			links.Website().Resource(`/js/blockarea_v0200.js`, map[string]string{}), // BlockArea
			"https://unpkg.com/vue@3/dist/vue.global.js",                            // Vue.js CDN
		},
		Scripts: []string{},
		StyleURLs: []string{
			cdn.JqueryUiCss_1_13_1(), // BlockArea styles
			"https://cdn.jsdelivr.net/npm/summernote@0.8.18/dist/summernote-lite.min.css", // Summernote styles
			"https://cdn.jsdelivr.net/npm/easymde/dist/easymde.min.css",                   // EasyMDE styles
			"https://cdn.jsdelivr.net/npm/codemirror@5.65.5/lib/codemirror.min.css",       // CodeMirror styles
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
			URL:  shared.NewLinks("/admin/blog").PostManager(),
		},
		{
			Name: "Edit Post",
			URL:  shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID()}),
		},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(links.Admin().Blog())

	buttonView := hb.Hyperlink().
		Class("btn btn-info ms-2 float-end").
		Child(hb.I().Class("bi bi-eye").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("View").
		Href("/blog/post/"+post.GetID()+"/"+post.GetSlug()).
		Attr("target", "_blank")

	buttonVersionHistory := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-clock-history").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Version History").
		Attr("data-bs-toggle", "modal").
		Attr("data-bs-target", "#versionHistoryModal")

	heading := hb.Heading1().
		HTML("Edit Post").
		Child(buttonCancel).
		Child(buttonView).
		Child(buttonVersionHistory)

	tabs := bs.NavTabs().
		Class("mb-3").
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "details", "active").
				Href(shared.NewLinks("/admin/blog").PostUpdate(map[string]string{
					"post_id": post.GetID(),
					"view":    "details",
				})).
				HTML("Details"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "content", "active").
				Href(shared.NewLinks("/admin/blog").PostUpdate(map[string]string{
					"post_id": post.GetID(),
					"view":    "content",
				})).
				HTML("Content"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "categories", "active").
				Href(shared.NewLinks("/admin/blog").PostUpdate(map[string]string{
					"post_id": post.GetID(),
					"view":    "categories",
				})).
				HTML("Categories"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "tags", "active").
				Href(shared.NewLinks("/admin/blog").PostUpdate(map[string]string{
					"post_id": post.GetID(),
					"view":    "tags",
				})).
				HTML("Tags"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "seo", "active").
				Href(shared.NewLinks("/admin/blog").PostUpdate(map[string]string{
					"post_id": post.GetID(),
					"view":    "seo",
				})).
				HTML("SEO")))

	postTitle := hb.Heading2().
		Class("mb-3").
		HTML("Post: ").
		HTML(post.GetTitle())

	var body hb.TagInterface

	switch view {
	case "details":
		body = controller.renderDetailsView(r, post)
	case "content":
		body = controller.renderContentView(r, post)
	case "categories":
		body = controller.renderCategoriesView(r, post)
	case "tags":
		body = controller.renderTagsView(r, post)
	case "seo":
		body = controller.renderSEOView(r, post)
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
					Style("margin-bottom:0;display:inline-block;")),
		).
		Child(
			hb.Div().
				Class("card-body").
				Child(body),
		)

	versioningModal := controller.renderVersioningModal(r, post)

	return hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(postTitle).
		Child(tabs).
		Child(card).
		Child(versioningModal)
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
		const urlCategoriesLoad = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "load-categories"}) + `';
		const urlCategoryAdd = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "add-category"}) + `';
		const urlCategoryRemove = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "remove-category"}) + `';
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

func (controller *postUpdateController) handleLoadCategories(w http.ResponseWriter, r *http.Request) string {
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
		const urlTagsLoad = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "load-tags"}) + `';
		const urlTagAdd = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "add-tag"}) + `';
		const urlTagRemove = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "remove-tag"}) + `';
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

func (controller *postUpdateController) renderDetailsView(r *http.Request, post blogstore.PostInterface) hb.TagInterface {
	htmlContent, err := postCategoriesFiles.ReadFile("post_details.html")
	if err != nil {
		slog.Error("Failed to read post details HTML template", "error", err)
		return hb.Div().HTML("Error loading details component")
	}

	jsContent, err := postCategoriesFiles.ReadFile("post_details.js")
	if err != nil {
		slog.Error("Failed to read post details JavaScript file", "error", err)
		return hb.Div().HTML("Error loading details component")
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		const postId = '` + post.GetID() + `';
		const urlDetailsLoad = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "load-details"}) + `';
		const urlDetailsSave = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "save-details"}) + `';
		const urlRegenerateImage = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "regenerate-image"}) + `';
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

func (controller *postUpdateController) renderContentView(r *http.Request, post blogstore.PostInterface) hb.TagInterface {
	htmlContent, err := postCategoriesFiles.ReadFile("post_content.html")
	if err != nil {
		slog.Error("Failed to read post content HTML template", "error", err)
		return hb.Div().HTML("Error loading content component")
	}

	jsContent, err := postCategoriesFiles.ReadFile("post_content.js")
	if err != nil {
		slog.Error("Failed to read post content JavaScript file", "error", err)
		return hb.Div().HTML("Error loading content component")
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		const postId = '` + post.GetID() + `';
		const urlContentLoad = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "load-content"}) + `';
		const urlContentSave = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "save-content"}) + `';
		const urlBlockEditorHandle = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "blockeditor-handle"}) + `';
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

func (controller *postUpdateController) renderSEOView(r *http.Request, post blogstore.PostInterface) hb.TagInterface {
	htmlContent, err := postCategoriesFiles.ReadFile("post_seo.html")
	if err != nil {
		slog.Error("Failed to read post SEO HTML template", "error", err)
		return hb.Div().HTML("Error loading SEO component")
	}

	jsContent, err := postCategoriesFiles.ReadFile("post_seo.js")
	if err != nil {
		slog.Error("Failed to read post SEO JavaScript file", "error", err)
		return hb.Div().HTML("Error loading SEO component")
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		const postId = '` + post.GetID() + `';
		const urlSEOLoad = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "load-seo"}) + `';
		const urlSEOSave = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "save-seo"}) + `';
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

func (controller *postUpdateController) renderVersioningModal(r *http.Request, post blogstore.PostInterface) hb.TagInterface {
	htmlContent, err := postCategoriesFiles.ReadFile("post_versioning.html")
	if err != nil {
		slog.Error("Failed to read post versioning HTML template", "error", err)
		return hb.Div().HTML("Error loading versioning component")
	}

	jsContent, err := postCategoriesFiles.ReadFile("post_versioning.js")
	if err != nil {
		slog.Error("Failed to read post versioning JavaScript file", "error", err)
		return hb.Div().HTML("Error loading versioning component")
	}

	htmlStr := string(htmlContent)
	htmlTemplate := hb.Wrap().HTML(htmlStr)
	componentScript := hb.Script(string(jsContent))

	// Config script for versioning component (uses global object to avoid const redeclaration)
	configScript := hb.Script(`
		window.postVersioningConfig = window.postVersioningConfig || {};
		window.postVersioningConfig.postId = '` + post.GetID() + `';
		window.postVersioningConfig.urlVersionsLoad = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "load-versions"}) + `';
		window.postVersioningConfig.urlVersionDetail = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "load-version-detail"}) + `';
		window.postVersioningConfig.urlVersionRestore = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": post.GetID(), "action": "restore-version-attributes"}) + `';
	`)

	// CSS for v-cloak
	vCloakStyle := hb.Style(`
		[v-cloak] { display: none; }
	`)

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	vueContainer := hb.Div().
		Child(vCloakStyle).
		Child(vueCDN).
		Child(configScript).
		Child(htmlTemplate).
		Child(componentScript)

	return vueContainer
}

func (controller *postUpdateController) handleLoadTags(w http.ResponseWriter, r *http.Request) string {
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

// Details component handlers
func (controller *postUpdateController) handleLoadDetails(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	postID := req.GetStringTrimmed(r, "post_id")
	if postID == "" {
		return api.Error("Post ID is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	post, err := blogStore.PostFindByID(ctx, postID)
	if err != nil || post == nil {
		return api.Error("Post not found").ToString()
	}

	return api.SuccessWithData("Details loaded successfully", map[string]any{
		"status":       post.GetStatus(),
		"image_url":    post.GetImageUrl(),
		"featured":     post.GetFeatured(),
		"published_at": post.GetPublishedAtCarbon().ToDateTimeString(),
		"editor":       post.GetEditor(),
		"memo":         post.GetMemo(),
	}).ToString()
}

func (controller *postUpdateController) handleSaveDetails(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	postID := req.GetStringTrimmed(r, "post_id")

	var reqData struct {
		Status      string `json:"post_status"`
		ImageURL    string `json:"post_image_url"`
		Featured    string `json:"post_featured"`
		PublishedAt string `json:"post_published_at"`
		Editor      string `json:"post_editor"`
		Memo        string `json:"post_memo"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.Status == "" {
		return api.Error("Status is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	post, err := blogStore.PostFindByID(ctx, postID)
	if err != nil || post == nil {
		return api.Error("Post not found").ToString()
	}

	// Normalize published_at - keep existing value if empty
	var publishedAt string
	if strings.TrimSpace(reqData.PublishedAt) == "" {
		// Keep the existing published_at value
		publishedAt = post.GetPublishedAtCarbon().ToDateTimeString()
	} else {
		parsedTime, err := time.Parse("2006-01-02T15:04", reqData.PublishedAt)
		if err != nil {
			return api.Error("Invalid published_at format").ToString()
		}
		publishedAt = parsedTime.Format("2006-01-02 15:04:05")
	}

	post.SetEditor(reqData.Editor)
	post.SetContentType(getContentTypeFromEditor(reqData.Editor))
	post.SetFeatured(reqData.Featured)
	post.SetImageUrl(reqData.ImageURL)
	post.SetMemo(reqData.Memo)
	post.SetPublishedAt(publishedAt)
	post.SetStatus(reqData.Status)

	if err := blogStore.PostUpdate(ctx, post); err != nil {
		controller.registry.GetLogger().Error("Error saving post details", "error", err.Error())
		return api.Error("System error. Saving post failed").ToString()
	}

	if err := createPostVersioning(context.Background(), controller.registry, post); err != nil {
		controller.registry.GetLogger().Error("Error creating post versioning", "error", err.Error())
	}

	return api.Success("Post saved successfully").ToString()
}

func (controller *postUpdateController) handleRegenerateImage(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	postID := req.GetStringTrimmed(r, "post_id")

	var reqData struct {
		PostID string `json:"post_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	post, err := blogStore.PostFindByID(ctx, postID)
	if err != nil || post == nil {
		return api.Error("Post not found").ToString()
	}

	agent := blogai.NewBlogWriterAgent(controller.registry.GetLogger())
	if agent == nil {
		return api.Error("Failed to initialize AI engine").ToString()
	}

	llmEngine, err := shared.LlmEngine(controller.registry)
	if err != nil || llmEngine == nil {
		return api.Error("Failed to initialize AI engine").ToString()
	}

	imageURL, err := agent.GenerateImage(llmEngine, post.GetTitle(), post.GetSummary())
	if err != nil {
		controller.registry.GetLogger().Error("BlogAi.PostUpdateV2.RegenerateImage", "error", err.Error())
		return api.Error("Failed to generate image").ToString()
	}

	post.SetImageUrl(imageURL)
	if err := blogStore.PostUpdate(ctx, post); err != nil {
		controller.registry.GetLogger().Error("BlogAi.PostUpdateV2.RegenerateImage.Save", "error", err.Error())
		return api.Error("Failed to save generated image").ToString()
	}

	return api.SuccessWithData("Image regenerated successfully", map[string]any{
		"image_url": imageURL,
	}).ToString()
}

// Content component handlers
func (controller *postUpdateController) handleLoadContent(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	postID := req.GetStringTrimmed(r, "post_id")
	if postID == "" {
		return api.Error("Post ID is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	post, err := blogStore.PostFindByID(ctx, postID)
	if err != nil || post == nil {
		return api.Error("Post not found").ToString()
	}

	return api.SuccessWithData("Content loaded successfully", map[string]any{
		"title":   post.GetTitle(),
		"summary": post.GetSummary(),
		"content": post.GetContent(),
		"editor":  post.GetEditor(),
	}).ToString()
}

func (controller *postUpdateController) handleSaveContent(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	postID := req.GetStringTrimmed(r, "post_id")

	var reqData struct {
		Title   string `json:"post_title"`
		Summary string `json:"post_summary"`
		Content string `json:"post_content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.Title == "" {
		return api.Error("Title is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	post, err := blogStore.PostFindByID(ctx, postID)
	if err != nil || post == nil {
		return api.Error("Post not found").ToString()
	}

	post.SetTitle(reqData.Title)
	post.SetSummary(reqData.Summary)
	post.SetContent(reqData.Content)

	if err := blogStore.PostUpdate(ctx, post); err != nil {
		controller.registry.GetLogger().Error("Error saving post content", "error", err.Error())
		return api.Error("System error. Saving post failed").ToString()
	}

	if err := createPostVersioning(context.Background(), controller.registry, post); err != nil {
		controller.registry.GetLogger().Error("Error creating post versioning", "error", err.Error())
	}

	return api.Success("Post saved successfully").ToString()
}

func (controller *postUpdateController) handleBlockEditorHandle(w http.ResponseWriter, r *http.Request) string {
	// This is a placeholder for BlockEditor handling
	// The actual implementation would depend on the BlockEditor library
	return api.Error("BlockEditor handle not implemented").ToString()
}

// SEO component handlers
func (controller *postUpdateController) handleLoadSEO(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	postID := req.GetStringTrimmed(r, "post_id")
	if postID == "" {
		return api.Error("Post ID is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	post, err := blogStore.PostFindByID(ctx, postID)
	if err != nil || post == nil {
		return api.Error("Post not found").ToString()
	}

	return api.SuccessWithData("SEO data loaded successfully", map[string]any{
		"slug":             post.GetSlug(),
		"canonical_url":    post.GetCanonicalURL(),
		"meta_description": post.GetMetaDescription(),
		"meta_keywords":    post.GetMetaKeywords(),
		"meta_robots":      post.GetMetaRobots(),
		"old_slugs":        post.GetOldSlugs(),
	}).ToString()
}

func (controller *postUpdateController) handleSaveSEO(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	postID := req.GetStringTrimmed(r, "post_id")

	var reqData struct {
		Slug            string   `json:"post_slug"`
		CanonicalURL    string   `json:"post_canonical_url"`
		MetaDescription string   `json:"post_meta_description"`
		MetaKeywords    string   `json:"post_meta_keywords"`
		MetaRobots      string   `json:"post_meta_robots"`
		OldSlugs        []string `json:"post_old_slugs"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	post, err := blogStore.PostFindByID(ctx, postID)
	if err != nil || post == nil {
		return api.Error("Post not found").ToString()
	}

	post.SetSlug(reqData.Slug)
	post.SetCanonicalURL(reqData.CanonicalURL)
	post.SetMetaDescription(reqData.MetaDescription)
	post.SetMetaKeywords(reqData.MetaKeywords)
	post.SetMetaRobots(reqData.MetaRobots)
	if err := post.SetOldSlugs(reqData.OldSlugs); err != nil {
		controller.registry.GetLogger().Error("Error setting old slugs", "error", err.Error())
		return api.Error("System error. Setting old slugs failed").ToString()
	}

	if err := blogStore.PostUpdate(ctx, post); err != nil {
		controller.registry.GetLogger().Error("Error saving post SEO", "error", err.Error())
		return api.Error("System error. Saving post failed").ToString()
	}

	if err := createPostVersioning(context.Background(), controller.registry, post); err != nil {
		controller.registry.GetLogger().Error("Error creating post versioning", "error", err.Error())
	}

	return api.Success("Post saved successfully").ToString()
}

// Versioning component handlers
func (controller *postUpdateController) handleLoadVersions(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	var reqData struct {
		PostID string `json:"post_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.PostID == "" {
		return api.Error("Post ID is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	if !blogStore.VersioningEnabled() {
		return api.SuccessWithData("Versions loaded successfully", map[string]any{
			"versioning_enabled": false,
			"versions":           []any{},
		}).ToString()
	}

	query := blogstore.NewVersioningQuery()
	query.SetEntityType(blogstore.VERSIONING_TYPE_POST)
	query.SetEntityID(reqData.PostID)
	query.SetOrderBy(versionstore.COLUMN_CREATED_AT)
	query.SetSortOrder(sb.DESC)
	query.SetLimit(50)

	controller.registry.GetLogger().Info("Loading versions for post", "post_id", reqData.PostID)

	versions, err := blogStore.EnableDebug(true).VersioningList(ctx, query)
	if err != nil {
		controller.registry.GetLogger().Error("Failed to load versions", "error", err.Error(), "post_id", reqData.PostID)
		return api.Error("Failed to load versions").ToString()
	}

	controller.registry.GetLogger().Info("Versions loaded from query", "count", len(versions), "post_id", reqData.PostID)

	// Filter versions by entity_id as a safety measure (in case versionstore doesn't filter correctly)
	filteredVersions := []blogstore.VersioningInterface{}
	for _, version := range versions {
		if version.EntityID() == reqData.PostID {
			filteredVersions = append(filteredVersions, version)
		}
	}

	controller.registry.GetLogger().Info("Versions after filtering", "count", len(filteredVersions), "post_id", reqData.PostID)

	// Convert versions to serializable format
	versionList := []map[string]any{}
	for _, version := range filteredVersions {
		versionList = append(versionList, map[string]any{
			"id":         version.ID(),
			"content":    version.Content(),
			"created_at": version.CreatedAt(),
		})
	}

	return api.SuccessWithData("Versions loaded successfully", map[string]any{
		"versioning_enabled": true,
		"versions":           versionList,
	}).ToString()
}

func (controller *postUpdateController) handleLoadVersionDetail(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	var reqData struct {
		VersionID string `json:"version_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.VersionID == "" {
		return api.Error("Version ID is required").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	version, err := blogStore.VersioningFindByID(ctx, reqData.VersionID)
	if err != nil || version == nil {
		return api.Error("Version not found").ToString()
	}

	var versionData map[string]string
	if err := json.Unmarshal([]byte(version.Content()), &versionData); err != nil {
		return api.Error("Invalid version data").ToString()
	}

	// Attribute label mapping for UI display
	attributeLabels := map[string]string{
		"title":            "Title",
		"slug":             "Slug",
		"content":          "Content",
		"summary":          "Summary",
		"status":           "Status",
		"featured":         "Featured",
		"image_url":        "Image URL",
		"memo":             "Memo",
		"canonical_url":    "Canonical URL",
		"meta_description": "Meta Description",
		"meta_keywords":    "Meta Keywords",
		"meta_robots":      "Meta Robots",
		"published_at":     "Published At",
		"author_id":        "Author ID",
		"metas":            "Metadata",
		"id":               "ID",
	}

	attributeList := []map[string]any{}
	for key, value := range versionData {
		label, ok := attributeLabels[key]
		if !ok {
			label = key
		}
		attributeList = append(attributeList, map[string]any{
			"key":   key,
			"label": label,
			"value": value,
		})
	}

	return api.SuccessWithData("Version detail loaded", map[string]any{
		"attributes": attributeList,
		"created_at": version.CreatedAt(),
	}).ToString()
}

func (controller *postUpdateController) handleRestoreVersionAttributes(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	var reqData struct {
		PostID     string   `json:"post_id"`
		VersionID  string   `json:"version_id"`
		Attributes []string `json:"attributes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if reqData.PostID == "" {
		return api.Error("Post ID is required").ToString()
	}

	if reqData.VersionID == "" {
		return api.Error("Version ID is required").ToString()
	}

	if len(reqData.Attributes) == 0 {
		return api.Error("At least one attribute must be selected").ToString()
	}

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	// Get the version to restore from
	version, err := blogStore.VersioningFindByID(ctx, reqData.VersionID)
	if err != nil || version == nil {
		return api.Error("Version not found").ToString()
	}

	// Get the current post
	post, err := blogStore.PostFindByID(ctx, reqData.PostID)
	if err != nil || post == nil {
		return api.Error("Post not found").ToString()
	}

	// Parse version content
	versionData := map[string]string{}
	if err := json.Unmarshal([]byte(version.Content()), &versionData); err != nil {
		return api.Error("Invalid version data").ToString()
	}

	// Apply selected attributes
	for _, attr := range reqData.Attributes {
		if val, ok := versionData[attr]; ok {
			post.Set(attr, val)
		}
	}

	// Update the post
	if err := blogStore.PostUpdate(ctx, post); err != nil {
		controller.registry.GetLogger().Error("Error updating post from version attributes", "error", err.Error())
		return api.Error("Error restoring attributes").ToString()
	}

	// Create a new version for the restoration
	if err := createPostVersioning(context.Background(), controller.registry, post); err != nil {
		controller.registry.GetLogger().Error("Error creating post versioning after restore", "error", err.Error())
	}

	return api.Success("Selected attributes restored successfully").ToString()
}

// Helper function to get content type from editor
func getContentTypeFromEditor(editor string) string {
	switch editor {
	case blogstore.POST_EDITOR_MARKDOWN:
		return blogstore.POST_CONTENT_TYPE_MARKDOWN
	case blogstore.POST_EDITOR_HTMLAREA:
		return blogstore.POST_CONTENT_TYPE_HTML
	case blogstore.POST_EDITOR_TEXTAREA:
		return blogstore.POST_CONTENT_TYPE_PLAIN_TEXT
	case blogstore.POST_EDITOR_BLOCKEDITOR, blogstore.POST_EDITOR_BLOCKAREA:
		return blogstore.POST_CONTENT_TYPE_BLOCKS
	default:
		return blogstore.POST_CONTENT_TYPE_PLAIN_TEXT
	}
}
