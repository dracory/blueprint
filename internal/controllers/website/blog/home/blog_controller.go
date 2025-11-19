package home

import (
	"log/slog"
	"net/http"
	"project/internal/controllers/website/blog/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/blogstore"
	"github.com/dracory/bs"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/spf13/cast"

	"github.com/samber/lo"
)

type blogController struct {
	app types.AppInterface
}

type blogControllerData struct {
	postList  []blogstore.Post
	postCount int64
	page      int
	perPage   int
}

func NewBlogController(app types.AppInterface) *blogController {
	return &blogController{
		app: app,
	}
}

func (controller *blogController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, errorMessage, links.Website().Home(), 10)
	}

	options := layouts.Options{
		WebsiteSection: "Blog",
		Title:          "Recent Posts",
		Content:        hb.Wrap().HTML(controller.page(data)),
		ScriptURLs: []string{
			"https://cdn.jsdelivr.net/gh/lesichkovm/slazy@latest/dist/slazy.min.js",
		},
	}

	if controller.app.GetConfig().GetCmsStoreUsed() {
		return layouts.NewCmsLayout(
			controller.app,
			r,
			options).ToHTML()
	} else {
		return layouts.NewBlankLayout(
			controller.app,
			r,
			options).ToHTML()
	}
}

func (controller *blogController) page(data blogControllerData) string {
	url := links.Website().Blog(map[string]string{
		"page": "",
	})

	pagination := bs.Pagination(bs.PaginationOptions{
		NumberItems:       cast.ToInt(data.postCount),
		CurrentPageNumber: data.page,
		PagesToShow:       10,
		PerPage:           data.perPage,
		URL:               url,
	})

	columnCards := lo.Map(data.postList, func(post blogstore.Post, index int) hb.TagInterface {
		return cardPost(post)
	})

	section := hb.Section().
		Style("background:#fff;padding-top:40px; padding-bottom: 40px;").
		Child(hb.Div().
			Class(`container`).
			Child(hb.Div().
				Class(`row g-4`).
				Children(columnCards)).
			Child(hb.Div().
				Class(`d-flex justify-content-center mt-5 pagination-primary-soft rounded mb-0`).
				HTML(pagination)))

	return hb.Wrap().Children([]hb.TagInterface{
		shared.SectionBanner(),
		section,
	}).ToHTML()
}

func (controller blogController) prepareData(r *http.Request) (data blogControllerData, errorMessage string) {
	perPage := 12 // 3 rows x 4 postss
	pageStr := req.GetStringTrimmed(r, "page")
	page := cast.ToInt(pageStr)

	if page < 0 {
		page = 0
	}

	blogStore := controller.app.GetBlogStore()

	if blogStore == nil {
		controller.app.GetLogger().Error("Error. At blogController.prepareData", slog.String("error", "blog store is not initialized"))
		return data, "Sorry, the blog is currently unavailable. Please try again later."
	}

	options := blogstore.PostQueryOptions{
		Status:    blogstore.POST_STATUS_PUBLISHED,
		SortOrder: "DESC",
		OrderBy:   "published_at",
		Offset:    page * perPage,
		Limit:     perPage,
	}

	postList, errList := blogStore.PostList(options)

	if errList != nil {
		controller.app.GetLogger().Error("Error. At blogController.page", slog.String("error", errList.Error()))
		return data, "Sorry, there was an error loading the posts. Please try again later."
	}

	postCount, errCount := blogStore.PostCount(options)

	if errCount != nil {
		controller.app.GetLogger().Error("Error. At blogController.page", slog.String("error", errCount.Error()))
		return data, "Sorry, there was an error loading the posts count. Please try again later."
	}

	data.page = page
	data.perPage = perPage
	data.postList = postList
	data.postCount = postCount

	return data, ""
}
