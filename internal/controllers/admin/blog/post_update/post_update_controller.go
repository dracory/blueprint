package post_update

import (
	"log/slog"
	"net/http"
	"project/internal/controllers/admin/blog/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"
	"project/pkg/blogblocks"
	"strings"

	"github.com/dracory/blockeditor"
	"github.com/dracory/blogstore"
	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

const VIEW_DETAILS = "details"
const VIEW_CONTENT = "content"
const VIEW_SEO = "seo"
const ACTION_BLOCKEDITOR_HANDLE = "blockeditor_handle"

type postUpdateController struct {
	app types.AppInterface
}

func NewPostUpdateController(app types.AppInterface) *postUpdateController {
	return &postUpdateController{app: app}
}

func (controller *postUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareDataAndValidate(r)

	if errorMessage != "" {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, errorMessage, links.Admin().Blog(), 10)
	}

	if data.action == ACTION_BLOCKEDITOR_HANDLE {
		return blockeditor.Handle(w, r, blogblocks.BlockEditorDefinitions())
	}

	if r.Method == http.MethodPost {
		return formPostUpdate(data).ToHTML()
	}

	return layouts.NewAdminLayout(controller.app, r, layouts.Options{
		Title:   "Edit Post | Blog",
		Content: controller.page(data),
		ScriptURLs: []string{
			cdn.Htmx_2_0_0(),
			cdn.Jquery_3_7_1(),
			cdn.TrumbowygJs_2_27_3(),
			cdn.Sweetalert2_10(),
			cdn.JqueryUiJs_1_13_1(), // needed for BlockArea
			links.Website().Resource(`/js/blockarea_v0200.js`, map[string]string{}), // needed for BlockArea
		},
		Scripts: []string{
			controller.script(),
		},
		StyleURLs: []string{
			cdn.TrumbowygCss_2_27_3(),
			cdn.JqueryUiCss_1_13_1(), // needed for BlockArea
		},
	}).ToHTML()
}

func (controller *postUpdateController) script() string {
	js := ``
	return js
}

func (controller *postUpdateController) page(data postUpdateControllerData) hb.TagInterface {
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
			URL:  shared.NewLinks().PostUpdate(map[string]string{"post_id": data.postID}),
		},
	})

	buttonSave := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Save").
		HxInclude("#FormPostUpdate").
		HxPost(shared.NewLinks().PostUpdate(map[string]string{"post_id": data.postID})).
		HxTarget("#FormPostUpdate")

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(links.Admin().Blog())

	heading := hb.Heading1().
		HTML("Edit Post").
		Child(buttonSave).
		Child(buttonCancel)

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header").
				Style(`display:flex;justify-content:space-between;align-items:center;`).
				Child(hb.Heading4().
					HTMLIf(data.view == VIEW_DETAILS, "Post Details").
					HTMLIf(data.view == VIEW_CONTENT, "Post Contents").
					HTMLIf(data.view == VIEW_SEO, "Post SEO").
					Style("margin-bottom:0;display:inline-block;")).
				Child(buttonSave),
		).
		Child(
			hb.Div().
				Class("card-body").
				Child(formPostUpdate(data)))

	tabs := bs.NavTabs().
		Class("mb-3").
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_DETAILS, "active").
				Href(shared.NewLinks().PostUpdate(map[string]string{
					"post_id": data.postID,
					"view":    VIEW_DETAILS,
				})).
				HTML("Details"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_CONTENT, "active").
				Href(shared.NewLinks().PostUpdate(map[string]string{
					"post_id": data.postID,
					"view":    VIEW_CONTENT,
				})).
				HTML("Content"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_SEO, "active").
				Href(shared.NewLinks().PostUpdate(map[string]string{
					"post_id": data.postID,
					"view":    VIEW_SEO,
				})).
				HTML("SEO")))

	postTitle := hb.Heading2().
		Class("mb-3").
		HTML("Post: ").
		HTML(data.post.Title())

	return hb.Div().
		Class("container").
		Child(heading).
		Child(breadcrumbs).
		Child(postTitle).
		Child(tabs).
		Child(card)
}

func (controller *postUpdateController) savePost(r *http.Request, data postUpdateControllerData) (d postUpdateControllerData, errorMessage string) {
	data.formCanonicalURL = req.GetStringTrimmed(r, "post_canonical_url")
	data.formContent = req.GetStringTrimmed(r, "post_content")
	data.formEditor = req.GetStringTrimmed(r, "post_editor")
	data.formFeatured = req.GetStringTrimmed(r, "post_featured")
	data.formImageUrl = req.GetStringTrimmed(r, "post_image_url")
	data.formMemo = req.GetStringTrimmed(r, "post_memo")
	data.formMetaDescription = req.GetStringTrimmed(r, "post_meta_description")
	data.formMetaKeywords = req.GetStringTrimmed(r, "post_meta_keywords")
	data.formMetaRobots = req.GetStringTrimmed(r, "post_meta_robots")
	data.formPublishedAt = req.GetStringTrimmed(r, "post_published_at")
	data.formSummary = req.GetStringTrimmed(r, "post_summary")
	data.formStatus = req.GetStringTrimmed(r, "post_status")
	data.formTitle = req.GetStringTrimmed(r, "post_title")

	if data.view == VIEW_DETAILS {
		if data.formStatus == "" {
			data.formErrorMessage = "Status is required"
			return data, ""
		}
	}

	if data.view == VIEW_CONTENT {
		if data.formTitle == "" {
			data.formErrorMessage = "Title is required"
			return data, ""
		}
	}

	if data.view == VIEW_DETAILS {
		// make sure the date is in the correct format
		data.formPublishedAt = lo.Substring(strings.ReplaceAll(data.formPublishedAt, " ", "T")+":00", 0, 19)
		publishedAt := lo.Ternary(data.formPublishedAt == "", sb.NULL_DATE, carbon.Parse(data.formPublishedAt).ToDateTimeString(carbon.UTC))
		data.post.SetEditor(data.formEditor)
		data.post.SetFeatured(data.formFeatured)
		data.post.SetImageUrl(data.formImageUrl)
		data.post.SetMemo(data.formMemo)
		data.post.SetPublishedAt(publishedAt)
		data.post.SetStatus(data.formStatus)
	}

	if data.view == VIEW_CONTENT {
		data.post.SetContent(data.formContent)
		data.post.SetSummary(data.formSummary)
		data.post.SetTitle(data.formTitle)
	}

	if data.view == VIEW_SEO {
		data.post.SetCanonicalURL(data.formCanonicalURL)
		data.post.SetMetaDescription(data.formMetaDescription)
		data.post.SetMetaKeywords(data.formMetaKeywords)
		data.post.SetMetaRobots(data.formMetaRobots)
	}

	err := controller.app.GetBlogStore().PostUpdate(data.post)

	if err != nil {
		controller.app.GetLogger().Error("At postUpdateController > prepareDataAndValidate", slog.String("error", err.Error()))
		data.formErrorMessage = "System error. Saving post failed"
		return data, ""
	}

	data.formSuccessMessage = "Post saved successfully"

	return data, ""
}

func (controller *postUpdateController) prepareDataAndValidate(r *http.Request) (data postUpdateControllerData, errorMessage string) {
	data.action = req.GetStringTrimmed(r, "action")
	data.postID = req.GetStringTrimmed(r, "post_id")
	data.view = req.GetStringTrimmedOr(r, "view", VIEW_DETAILS)

	if data.view == "" {
		data.view = VIEW_DETAILS
	}

	if data.postID == "" {
		return data, "Post ID is required"
	}

	var err error
	data.post, err = controller.app.GetBlogStore().PostFindByID(data.postID)

	if err != nil {
		controller.app.GetLogger().Error("At postUpdateController > prepareDataAndValidate", slog.String("error", err.Error()))
		return data, "Post not found"
	}

	if data.post == nil {
		return data, "Post not found"
	}

	data.formCanonicalURL = data.post.CanonicalURL()
	data.formContent = data.post.Content()
	data.formEditor = data.post.Editor()
	data.formImageUrl = data.post.ImageUrl()
	data.formFeatured = data.post.Featured()
	data.formMetaDescription = data.post.MetaDescription()
	data.formMetaKeywords = data.post.MetaKeywords()
	data.formMetaRobots = data.post.MetaRobots()
	data.formMemo = data.post.Memo()
	data.formPublishedAt = data.post.PublishedAtCarbon().ToDateTimeString()
	data.formSummary = data.post.Summary()
	data.formStatus = data.post.Status()
	data.formTitle = data.post.Title()

	if r.Method != http.MethodPost {
		return data, ""
	}

	return controller.savePost(r, data)
}

type postUpdateControllerData struct {
	action string
	postID string
	post   *blogstore.Post
	view   string

	formErrorMessage    string
	formSuccessMessage  string
	formCanonicalURL    string
	formContent         string
	formEditor          string
	formFeatured        string
	formImageUrl        string
	formMemo            string
	formMetaDescription string
	formMetaKeywords    string
	formMetaRobots      string
	formPublishedAt     string
	formStatus          string
	formSummary         string
	formTitle           string
}
