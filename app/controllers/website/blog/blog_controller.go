package blog

import (
	"log/slog"
	"net/http"
	"project/config"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"strings"

	"github.com/dracory/base/req"
	"github.com/gouniverse/blogstore"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/hb"
	"github.com/spf13/cast"

	"github.com/samber/lo"
)

type blogController struct{}

type blogControllerData struct {
	postList  []blogstore.Post
	postCount int64
	page      int
	perPage   int
}

func NewBlogController() *blogController {
	return &blogController{}
}

func (controller *blogController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return helpers.ToFlashError(w, r, errorMessage, links.Website().Home(), 10)
	}

	return layouts.NewCmsLayout(layouts.Options{
		Request:        r,
		WebsiteSection: "Blog",
		Title:          "Recent Posts",
		Content:        hb.Wrap().HTML(controller.page(data)),
	}).ToHTML()
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
		postImageURL := post.ImageUrlOrDefault()

		publishedAt := lo.Ternary(post.PublishedAt() == "", "", post.PublishedAtCarbon().Format("d M, Y"))

		postURL := links.NewWebsiteLinks().BlogPost(post.ID(), post.Slug())

		postImage := hb.Image(postImageURL).
			Class("card-img-top rounded-3").
			Style("object-fit: cover;").
			Style("max-height: 180px;").
			Style("aspect-ratio: 9/6;").
			Style("border-radius: 0.5rem").
			Alt("")

		postTitle := hb.Heading5().
			Class("card-title").
			Style("font-size: 16px; color: #224b8e; margin-bottom: 10px; text-align: left; font-weight: 800;").
			Text(post.Title())

		postPublished := hb.Paragraph().
			Style("font-size: 12px;	color: #6c757d;	margin-bottom: 20px; text-align: right;").
			Text(publishedAt)

		// postPublished := hb.Span().
		// 	Class(`small`).
		// 	Style(`font-size:12px;color:#666;display:inline-block;padding-right:10px;padding-top:10px;`).
		// 	HTML(publishedAt)

		// postImage := hb.Div().Class(`overflow-hidden rounded-3`).Children([]hb.TagInterface{
		// 	hb.Image().
		// 		Class(`card-img`).
		// 		Style(`object-fit:cover;max-height:180px;`).
		// 		Src(postImageURL).
		// 		Alt("course image").
		// 		Attr("loading", "lazy"),
		// 	hb.Div().
		// 		Class(`bg-overlay bg-dark opacity-4`),
		// 	// Badge
		// 	// bs.CardImageTop().Class(`d-flex align-items-start`).Children([]hb.TagInterface{
		// 	// 	hb.Div().Class(`badge text-bg-danger`).Style(`position:absolute;top:10px;left:10px;`).HTML("Student life"),
		// 	// }),
		// })

		postSummary := hb.Paragraph().
			Class("card-text").
			Text(post.Summary()).
			Style(`text-align: left;`).
			Style(`font-size: 14px;`).
			Style(`font-weight: 400;`).
			Style(`overflow: hidden;`).
			Style(`text-overflow: ellipsis;`).
			Style(`display: -webkit-box;`).
			Style(`-webkit-line-clamp: 2;`).
			Style(`-webkit-box-orient: vertical;`)

		separator := hb.HR().
			Style(`width: 80%`).
			Style(`margin: 0 auto`).
			Style(`border: 0`).
			Style(`height: 2px`).
			Style(`background-image: linear-gradient(to right, rgba(0, 0, 0, 0), rgba(0, 0, 0, 0.75), rgba(0, 0, 0, 0))`).
			Style(`opacity: 0.25`).
			Style(`margin-bottom: 20px`)

		card := hb.Div().
			Class("card").
			Style("border: none; margin-bottom: 20px;").
			Child(postImage).
			Child(hb.Div().
				Class("card-body").
				Style(`padding: 20px 10px;`).
				Child(postTitle).
				Child(postSummary)).
			Child(hb.Div().
				Class("card-footer").
				Style(`background: none;border: none;padding: 0px;`).
				Child(postPublished).
				Child(separator))

		link := hb.Hyperlink().
			Href(postURL).
			Target("_blank").
			Style("text-decoration: none; color: inherit;").
			Style("display: flex; height: 100%;").
			Child(card)

		return hb.Div().
			Class("col-md-3 col-sm-6 d-flex align-items-stretch").
			Child(link)
	})

	section := hb.Section().
		Style("background:#fff;padding-top:40px; padding-bottom: 40px;").
		Children([]hb.TagInterface{
			bs.Container().Children([]hb.TagInterface{
				bs.Row().Class(`g-4`).Children(columnCards),
				hb.Div().Class(`d-flex justify-content-center mt-5 pagination-primary-soft rounded mb-0`).HTML(pagination),
			}),
		})

	return hb.Wrap().Children([]hb.TagInterface{
		controller.sectionBanner(),
		section,
	}).ToHTML()
}

func (controller blogController) prepareData(r *http.Request) (data blogControllerData, errorMessage string) {
	perPage := 12 // 3 rows x 4 postss
	pageStr := strings.TrimSpace(req.Value(r, "page"))
	page := cast.ToInt(pageStr)

	if page < 0 {
		page = 0
	}

	options := blogstore.PostQueryOptions{
		Status:    blogstore.POST_STATUS_PUBLISHED,
		SortOrder: "DESC",
		OrderBy:   "published_at",
		Offset:    page * perPage,
		Limit:     perPage,
	}

	postList, errList := config.BlogStore.PostList(options)

	if errList != nil {
		config.Logger.Error("Error. At blogController.page", slog.String("error", errList.Error()))
		return data, "Sorry, there was an error loading the posts. Please try again later."
	}

	postCount, errCount := config.BlogStore.PostCount(options)

	if errCount != nil {
		config.Logger.Error("Error. At blogController.page", slog.String("error", errCount.Error()))
		return data, "Sorry, there was an error loading the posts count. Please try again later."
	}

	data.page = page
	data.perPage = perPage
	data.postList = postList
	data.postCount = postCount

	return data, ""
}

func (controller blogController) sectionBanner() *hb.Tag {
	style := hb.Style(`
.fill-success {
	fill: #0cbc87 !important;
}
.fill-orange {
	fill: #fd7e14 !important;
}
.fill-purple {
	fill: #6f42c1 !important;
}
	`)

	nav := hb.Nav().
		Class("breadcrumb mb-0").
		Attr("aria-label", "breadcrumb").
		Child(hb.OL().
			Class("breadcrumb mb-0").
			Child(hb.LI().
				Class("breadcrumb-item").
				Child(hb.A().
					Href(links.NewWebsiteLinks().Home()).
					HTML("Home"),
				)).
			Child(hb.LI().
				Class("breadcrumb-item active").
				Attr("aria-current", "page").
				Child(hb.A().
					Href(links.NewWebsiteLinks().Blog(map[string]string{})).
					HTML("Blog"),
				),
			))

	section := hb.Section().
		Style("background:#1C1626;").
		Style("padding: 30px 0px;").
		Child(hb.Div().
			Class("container").
			Child(hb.Div().
				Class("row").
				Child(hb.Div().
					//HTML(decorationCross).
					Class("col-lg-10 mx-auto text-center").
					Style(`position: relative;`).
					Child(hb.I().Class("bi bi-crosshair").
						Style("color: magenta;").
						Style(`position: absolute; top: 0px; left: 0px;`).
						Style(`font-size: 30px; margin-left: 10px;`)).
					Child(hb.I().Class("bi bi-asterisk").
						Style("color: magenta;").
						Style(`position: absolute; bottom: -10px; left: 100px;`).
						Style(`font-size: 30px; margin-right: 10px; transform: rotate(180deg);`)).
					Child(hb.I().Class("bi bi-star").
						Style("color: magenta;").
						Style(`position: absolute; top: 10px; right: 35px;`).
						Style(`font-size: 30px; margin-left: 10px;`)).
					Child(hb.I().Class("bi bi-star").
						Style("color: magenta;").
						Style(`position: absolute; top: 15px; right: 0px;`).
						Style(`font-size: 30px; margin-left: 10px;`)).
					Child(hb.I().Class("bi bi-star").
						Style("color: magenta;").
						Style(`position: absolute; top: 40px; right: 17px;`).
						Style(`font-size: 30px; margin-left: 10px;`)).
					Child(hb.H1().Style("color:white;").HTML("Blog")).
					Child(hb.Div().
						Class("d-flex justify-content-center position-relative").
						Child(nav),
					),
				),
			),
		)

	return hb.Wrap().
		Child(style).
		Child(section)
}
