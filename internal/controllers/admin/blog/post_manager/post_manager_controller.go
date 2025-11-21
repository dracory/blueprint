package post_manager

import (
	"log/slog"
	"net/http"

	"project/internal/controllers/admin/blog/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/blogstore"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/spf13/cast"
)

// == CONTROLLER ==============================================================

type postManagerController struct {
	app types.AppInterface
}

// == CONSTRUCTOR =============================================================

func NewPostManagerController(app types.AppInterface) *postManagerController {
	return &postManagerController{app: app}
}

func (controller *postManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, errorMessage, links.Admin().Home(), 10)
	}

	return layouts.NewAdminLayout(controller.app, r, layouts.Options{
		Title:   "Blog | Post Manager",
		Content: controller.page(data),
		ScriptURLs: []string{
			cdn.Htmx_2_0_0(),
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

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
			URL:  shared.NewLinks().Home(),
		},
	})

	actionButtons := hb.Div().
		Class("d-flex gap-2 float-end")

	buttonAiHome := hb.Hyperlink().
		Class("btn btn-light text-dark d-inline-flex align-items-center").
		Child(hb.I().Class("bi bi-stars me-2")).
		HTML("AI Tools").
		Href(shared.NewLinks().AiTools())

	buttonSettings := hb.Hyperlink().
		Class("btn btn-outline-secondary d-inline-flex align-items-center").
		Child(hb.I().Class("bi bi-gear me-2")).
		HTML("Settings").
		Href(shared.NewLinks().BlogSettings())

	buttonPostNew := hb.Button().
		Class("btn btn-primary d-inline-flex align-items-center").
		Child(hb.I().Class("bi bi-plus-circle me-2")).
		HTML("New Post").
		HxGet(shared.NewLinks().PostCreate()).
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

	data.blogList, err = controller.app.GetBlogStore().
		// EnableDebug(true).
		PostList(r.Context(), query)

	if err != nil {
		controller.app.GetLogger().Error("At managerController > prepareData", slog.String("error", err.Error()))
		return data, "error retrieving posts"
	}

	// DEBUG: cfmt.Successln("Invoice List: ", blogList)

	data.blogCount, err = controller.app.GetBlogStore().
		// EnableDebug().
		PostCount(r.Context(), query)

	if err != nil {
		controller.app.GetLogger().Error("At managerController > prepareData", slog.String("error", err.Error()))
		return data, "Error retrieving posts count"
	}

	return data, ""
}

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
	blogList   []blogstore.Post
	blogCount  int64
}
