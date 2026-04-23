package post_manager

import (
	"embed"
	"encoding/json"
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
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

//go:embed *.html
//go:embed *.js
var postsFiles embed.FS

const (
	actionLoadPosts  = "load-posts"
	actionDeletePost = "delete-post"
	actionCreatePost = "create-post"
)

// == CONTROLLER ==============================================================

type postManagerController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewPostManagerController(registry registry.RegistryInterface) *postManagerController {
	return &postManagerController{registry: registry}
}

func (controller *postManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadPosts:
		return controller.handleLoadPosts(w, r)
	case actionCreatePost:
		return controller.handleCreatePost(w, r)
	case actionDeletePost:
		return controller.handleDeletePost(w, r)
	default:
		return controller.renderPage(w, r)
	}
}

func (controller *postManagerController) renderPage(w http.ResponseWriter, r *http.Request) string {
	authUser := helpers.GetAuthUser(r)
	if authUser == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "You are not logged in. Please login to continue.", links.Admin().Blog(), 10)
	}

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.Admin().Home()},
		{Name: "Blog", URL: links.Admin().Blog()},
		{Name: "Post Manager", URL: ""},
	})

	actionButtons := hb.Div().
		Class("d-flex gap-2 float-end")

	buttonAiHome := hb.Hyperlink().
		Class("btn btn-light text-dark d-inline-flex align-items-center").
		Child(hb.I().Class("bi bi-stars me-2")).
		HTML("AI Tools").
		Href(shared.NewLinks("/admin/blog").AiTools())

	buttonSettings := hb.Hyperlink().
		Class("btn btn-outline-secondary d-inline-flex align-items-center").
		Child(hb.I().Class("bi bi-gear me-2")).
		HTML("Settings").
		Href(shared.NewLinks("/admin/blog").BlogSettings())

	actionButtons = actionButtons.
		Child(buttonAiHome).
		Child(buttonSettings)

	heading := hb.Heading1().HTML("Blog. Post Manager").Child(actionButtons)

	htmlContent, err := postsFiles.ReadFile("posts.html")
	if err != nil {
		slog.Error("Failed to read posts HTML template", "error", err)
		return hb.Div().HTML("Error loading posts component").ToHTML()
	}

	jsContent, err := postsFiles.ReadFile("posts.js")
	if err != nil {
		slog.Error("Failed to read posts JavaScript file", "error", err)
		return hb.Div().HTML("Error loading posts component").ToHTML()
	}

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	initScript := hb.Script(`
		const urlPostsLoad = '` + shared.NewLinks("/admin/blog").PostManager(map[string]string{"action": actionLoadPosts}) + `';
		const urlPostDelete = '` + shared.NewLinks("/admin/blog").PostManager(map[string]string{"action": actionDeletePost}) + `';
		const urlPostCreate = '` + shared.NewLinks("/admin/blog").PostManager(map[string]string{"action": actionCreatePost}) + `';
		const urlAiPostContentUpdate = '` + shared.NewLinks("/admin/blog").AiPostContentUpdate(map[string]string{"post_id": "POST_ID_PLACEHOLDER"}) + `';
		const urlPostUpdate = '` + shared.NewLinks("/admin/blog").PostUpdate(map[string]string{"post_id": "POST_ID_PLACEHOLDER"}) + `';
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
		Title:   "Blog | Post Manager",
		Content: content,
		ScriptURLs: []string{
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *postManagerController) handleLoadPosts(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	blogStore := controller.registry.GetBlogStore()
	if blogStore == nil {
		return api.Error("Blog store not available").ToString()
	}

	// Read parameters from query string for GET requests
	page := cast.ToInt(req.GetStringTrimmedOr(r, "page", "0"))
	perPage := cast.ToInt(req.GetStringTrimmedOr(r, "per_page", "10"))
	sortOrder := req.GetStringTrimmedOr(r, "sort_order", sb.DESC)
	sortBy := req.GetStringTrimmedOr(r, "sort_by", blogstore.COLUMN_CREATED_AT)
	status := req.GetStringTrimmed(r, "status")
	search := req.GetStringTrimmed(r, "search")
	dateFrom := req.GetStringTrimmedOr(r, "date_from", carbon.Now().AddYears(-1).ToDateString())
	dateTo := req.GetStringTrimmedOr(r, "date_to", carbon.Now().ToDateString())

	query := blogstore.PostQueryOptions{
		Search:               search,
		Offset:               page * perPage,
		Limit:                perPage,
		Status:               status,
		CreatedAtGreaterThan: dateFrom + " 00:00:00",
		CreatedAtLessThan:    dateTo + " 23:59:59",
		SortOrder:            sortOrder,
		OrderBy:              sortBy,
	}

	posts, err := blogStore.PostList(ctx, query)
	if err != nil {
		slog.Error("Failed to load posts", "error", err)
		return api.Error("Failed to load posts").ToString()
	}

	postList := []map[string]any{}
	for _, post := range posts {
		postList = append(postList, map[string]any{
			"id":           post.GetID(),
			"title":        post.GetTitle(),
			"status":       post.GetStatus(),
			"featured":     post.GetFeatured(),
			"published_at": post.GetPublishedAt(),
			"created_at":   post.GetCreatedAt(),
			"updated_at":   post.GetUpdatedAt(),
			"slug":         post.GetSlug(),
			"image_url":    post.GetImageUrl(),
		})
	}

	count, err := blogStore.PostCount(ctx, query)
	if err != nil {
		slog.Error("Failed to get posts count", "error", err)
		return api.Error("Failed to get posts count").ToString()
	}

	return api.SuccessWithData("Posts loaded successfully", map[string]any{
		"posts": postList,
		"total": count,
	}).ToString()
}

func (controller *postManagerController) handleDeletePost(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

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

	post, err := blogStore.PostFindByID(ctx, reqData.PostID)
	if err != nil {
		slog.Error("Failed to find post for delete", "error", err)
		return api.Error("Post not found").ToString()
	}

	if err := blogStore.PostDelete(ctx, post); err != nil {
		slog.Error("Failed to delete post", "error", err)
		return api.Error("Failed to delete post").ToString()
	}

	return api.SuccessWithData("Post deleted successfully", map[string]any{}).ToString()
}

func (controller *postManagerController) handleCreatePost(w http.ResponseWriter, r *http.Request) string {
	ctx := r.Context()

	if r.Method != http.MethodPost {
		return api.Error("Method not allowed").ToString()
	}

	var reqData struct {
		Title string `json:"title"`
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

	post := blogstore.NewPost()
	post.SetTitle(reqData.Title)

	if err := blogStore.PostCreate(ctx, post); err != nil {
		slog.Error("Failed to create post", "error", err)
		return api.Error("Failed to create post").ToString()
	}

	return api.SuccessWithData("Post created successfully", map[string]any{
		"id": post.GetID(),
	}).ToString()
}

// Deprecated: Use renderPage with Vue.js instead. This method is kept for backwards compatibility.
func (controller *postManagerController) page(data postManagerControllerData) hb.TagInterface {
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
			URL:  shared.NewLinks("/admin/blog").Home(),
		},
	})

	actionButtons := hb.Div().
		Class("d-flex gap-2 float-end")

	buttonAiHome := hb.Hyperlink().
		Class("btn btn-light text-dark d-inline-flex align-items-center").
		Child(hb.I().Class("bi bi-stars me-2")).
		HTML("AI Tools").
		Href(shared.NewLinks("/admin/blog").AiTools())

	buttonSettings := hb.Hyperlink().
		Class("btn btn-outline-secondary d-inline-flex align-items-center").
		Child(hb.I().Class("bi bi-gear me-2")).
		HTML("Settings").
		Href(shared.NewLinks("/admin/blog").BlogSettings())

	buttonPostNew := hb.Button().
		Class("btn btn-primary d-inline-flex align-items-center").
		Child(hb.I().Class("bi bi-plus-circle me-2")).
		HTML("New Post").
		HxGet(shared.NewLinks("/admin/blog").PostCreate()).
		HxTarget("body").
		HxSwap("beforeend")

	actionButtons = actionButtons.
		Child(buttonAiHome).
		Child(buttonSettings).
		Child(buttonPostNew)

	title := hb.Heading1().
		HTML("Blog. Post Manager").
		Child(actionButtons)

	return layouts.AdminPage(
		title,
		breadcrumbs,
		tablePostList(data),
	)
}

// Deprecated: Use handleLoadPosts API endpoint with Vue.js instead. This method is kept for backwards compatibility.
func (controller *postManagerController) prepareData(r *http.Request) (data postManagerControllerData, errorMessage string) {
	var err error

	authUser := helpers.GetAuthUser(r)

	if authUser == nil {
		return data, "You are not logged in. Please login to continue."
	}

	data.page = req.GetStringTrimmed(r, "page")
	data.pageInt = cast.ToInt(data.page)
	data.perPage = cast.ToInt(req.GetStringTrimmedOr(r, "per_page", "10"))
	data.sortOrder = req.GetStringTrimmedOr(r, "sort_order", sb.DESC)
	data.sortBy = req.GetStringTrimmedOr(r, "by", blogstore.COLUMN_CREATED_AT)
	data.status = req.GetStringTrimmed(r, "status")
	data.search = req.GetStringTrimmed(r, "search")
	data.dateFrom = req.GetStringTrimmedOr(r, "date_from", carbon.Now().AddYears(-1).ToDateString())
	data.dateTo = req.GetStringTrimmedOr(r, "date_to", carbon.Now().ToDateString())
	data.customerID = req.GetStringTrimmed(r, "customer_id")

	query := blogstore.PostQueryOptions{
		Search:               data.search,
		Offset:               data.pageInt * data.perPage,
		Limit:                data.perPage,
		Status:               data.status,
		CreatedAtGreaterThan: data.dateFrom + " 00:00:00",
		CreatedAtLessThan:    data.dateTo + " 23:59:59",
		SortOrder:            data.sortOrder,
		OrderBy:              data.sortBy,
	}

	data.blogList, err = controller.registry.GetBlogStore().
		// EnableDebug(true).
		PostList(r.Context(), query)

	if err != nil {
		controller.registry.GetLogger().Error("At managerController > prepareData", slog.String("error", err.Error()))
		return data, "error retrieving posts"
	}

	// DEBUG: cfmt.Successln("Invoice List: ", blogList)

	data.blogCount, err = controller.registry.GetBlogStore().
		// EnableDebug().
		PostCount(r.Context(), query)

	if err != nil {
		controller.registry.GetLogger().Error("At managerController > prepareData", slog.String("error", err.Error()))
		return data, "Error retrieving posts count"
	}

	return data, ""
}

// Deprecated: Used only by deprecated server-side rendering methods. Use Vue.js API instead.
type postManagerControllerData struct {
	// r            *http.Request
	page       string
	pageInt    int
	perPage    int
	sortOrder  string
	sortBy     string
	status     string
	search     string
	customerID string
	dateFrom   string
	dateTo     string
	blogList   []blogstore.PostInterface
	blogCount  int64
}
