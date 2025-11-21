package post_update

import (
	"log/slog"
	"net/http"

	"project/internal/controllers/admin/blog/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/blogstore"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dracory/req"
)

type postUpdateController struct {
	app types.AppInterface
}

func NewPostUpdateController(app types.AppInterface) *postUpdateController {
	return &postUpdateController{app: app}
}

func (controller *postUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	postID := req.GetStringTrimmed(r, "post_id")
	view := req.GetStringTrimmedOr(r, "view", "content")

	if postID == "" {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, "Post ID is required", links.Admin().Blog(), 10)
	}

	post, err := controller.app.GetBlogStore().PostFindByID(r.Context(), postID)
	if err != nil {
		controller.app.GetLogger().Error(
			"Error. postUpdateController: PostFindByID",
			slog.String("error", err.Error()),
			slog.String("post_id", postID),
		)
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, "Post not found", links.Admin().Blog(), 10)
	}

	if post == nil {
		controller.app.GetLogger().Warn(
			"Warning. postUpdateController: PostFindByID",
			slog.String("error", "Post not found"),
			slog.String("post_id", postID),
		)
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, "Post not found", links.Admin().Blog(), 10)
	}

	pageContent := controller.page(r, post, view)

	return layouts.NewAdminLayout(controller.app, r, layouts.Options{
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

func (controller *postUpdateController) page(r *http.Request, post *blogstore.Post, view string) hb.TagInterface {
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
			URL:  shared.NewLinks().PostUpdate(map[string]string{"post_id": post.ID()}),
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
					"post_id": post.ID(),
					"view":    "details",
				})).
				HTML("Details"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "content", "active").
				Href(shared.NewLinks().PostUpdate(map[string]string{
					"post_id": post.ID(),
					"view":    "content",
				})).
				HTML("Content"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(view == "seo", "active").
				Href(shared.NewLinks().PostUpdate(map[string]string{
					"post_id": post.ID(),
					"view":    "seo",
				})).
				HTML("SEO")))

	postTitle := hb.Heading2().
		Class("mb-3").
		HTML("Post: ").
		HTML(post.Title())

	var body hb.TagInterface

	switch view {
	case "details":
		component := NewPostDetailsComponent(controller.app)
		body = liveflux.Placeholder(component, map[string]string{
			"post_id": post.ID(),
		})
	case "content":
		component := NewPostContentComponent(controller.app)
		body = liveflux.Placeholder(component, map[string]string{
			"post_id": post.ID(),
		})
	case "seo":
		component := NewPostSEOComponent(controller.app)
		body = liveflux.Placeholder(component, map[string]string{
			"post_id": post.ID(),
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
					HTMLIf(view == "seo", "Post SEO").
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
