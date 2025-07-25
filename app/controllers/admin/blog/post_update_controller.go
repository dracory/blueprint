package admin

import (
	"log/slog"
	"net/http"
	"project/app/layouts"
	"project/app/links"
	"project/config"
	"project/internal/helpers"
	"project/pkg/blogblocks"
	"strings"

	"github.com/dracory/base/req"
	"github.com/dromara/carbon/v2"
	"github.com/gouniverse/blockeditor"
	"github.com/gouniverse/blogstore"
	"github.com/gouniverse/bs"
	"github.com/gouniverse/cdn"
	"github.com/gouniverse/form"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/sb"
	"github.com/samber/lo"
)

const VIEW_DETAILS = "details"
const VIEW_CONTENT = "content"
const VIEW_SEO = "seo"
const ACTION_BLOCKEDITOR_HANDLE = "blockeditor_handle"

type postUpdateController struct{}

func NewPostUpdateController() *postUpdateController {
	return &postUpdateController{}
}

func (controller postUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareDataAndValidate(r)

	if errorMessage != "" {
		return helpers.ToFlashError(w, r, errorMessage, links.NewAdminLinks().BlogPostManager(map[string]string{}), 10)
	}

	if data.action == ACTION_BLOCKEDITOR_HANDLE {
		return blockeditor.Handle(w, r, blogblocks.BlockEditorDefinitions())
	}

	if r.Method == http.MethodPost {
		return controller.form(data).ToHTML()
	}

	return layouts.NewAdminLayout(r, layouts.Options{
		Title:   "Edit Post | Blog",
		Content: controller.page(data),
		ScriptURLs: []string{
			cdn.Jquery_3_7_1(),
			cdn.TrumbowygJs_2_27_3(),
			cdn.Sweetalert2_10(),
			cdn.JqueryUiJs_1_13_1(), // needed for BlockArea
			links.NewWebsiteLinks().Resource(`/js/blockarea_v0200.js`, map[string]string{}), // needed for BlockArea
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

func (controller postUpdateController) script() string {
	js := ``
	return js
}

func (controller postUpdateController) page(data postUpdateControllerData) hb.TagInterface {
	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{
			Name: "Home",
			URL:  links.NewAdminLinks().Home(map[string]string{}),
		},
		{
			Name: "Post Manager",
			URL:  links.NewAdminLinks().BlogPostManager(map[string]string{}),
		},
		{
			Name: "Edit Post",
			URL:  links.NewAdminLinks().BlogPostUpdate(map[string]string{"post_id": data.postID}),
		},
	})

	buttonSave := hb.Button().
		Class("btn btn-primary ms-2 float-end").
		Child(hb.I().Class("bi bi-save").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Save").
		HxInclude("#FormPostUpdate").
		HxPost(links.NewAdminLinks().BlogPostUpdate(map[string]string{"postID": data.postID})).
		HxTarget("#FormPostUpdate")

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(links.NewAdminLinks().BlogPostManager(map[string]string{}))

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
				Child(controller.form(data)))

	tabs := bs.NavTabs().
		Class("mb-3").
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_DETAILS, "active").
				Href(links.NewAdminLinks().BlogPostUpdate(map[string]string{
					"post_id": data.postID,
					"view":    VIEW_DETAILS,
				})).
				HTML("Details"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_CONTENT, "active").
				Href(links.NewAdminLinks().BlogPostUpdate(map[string]string{
					"post_id": data.postID,
					"view":    VIEW_CONTENT,
				})).
				HTML("Content"))).
		Child(bs.NavItem().
			Child(bs.NavLink().
				ClassIf(data.view == VIEW_SEO, "active").
				Href(links.NewAdminLinks().BlogPostUpdate(map[string]string{
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

func (controller postUpdateController) form(data postUpdateControllerData) hb.TagInterface {
	fieldsDetails := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Status",
			Name:  "post_status",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formStatus,
			Help:  "The status of this blog post. Published posts will be displayed on the website.",
			Options: []form.FieldOption{
				{
					Value: "- not selected -",
					Key:   "",
				},
				{
					Value: "Draft",
					Key:   blogstore.POST_STATUS_DRAFT,
				},
				{
					Value: "Published",
					Key:   blogstore.POST_STATUS_PUBLISHED,
				},
				{
					Value: "Unpublished",
					Key:   blogstore.POST_STATUS_UNPUBLISHED,
				},
				{
					Value: "In Trash Bin",
					Key:   blogstore.POST_STATUS_TRASH,
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Image URL",
			Name:  "post_image_url",
			Type:  form.FORM_FIELD_TYPE_IMAGE,
			Value: data.formImageUrl,
			Help:  "The image that will be displayed on the blog post. If left empty, the default image will be used.",
		}),
		form.NewField(form.FieldOptions{
			Label: "Featured",
			Name:  "post_featured",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formFeatured,
			Help:  "Is this blog post featured? Featured posts will be displayed on the home page.",
			Options: []form.FieldOption{
				{
					Value: "- not selected -",
					Key:   "",
				},
				{
					Value: "No",
					Key:   "no",
				},
				{
					Value: "Yes",
					Key:   "yes",
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Published At",
			Name:  "post_published_at",
			Type:  form.FORM_FIELD_TYPE_DATETIME,
			Value: data.formPublishedAt,
			Help:  "The date this blog post was published.",
		}),
		form.NewField(form.FieldOptions{
			Label: "Editor",
			Name:  "post_editor",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formEditor,
			Help:  "The editor that will be used while editing this blogpost. Once set, this should not be changed, or the content may be lost.",
			Options: []form.FieldOption{
				{
					Value: "- not selected -",
					Key:   "",
				},
				{
					Value: "BlockArea",
					Key:   blogstore.POST_EDITOR_BLOCKAREA,
				},
				{
					Value: "BlockEditor",
					Key:   blogstore.POST_EDITOR_BLOCKEDITOR,
				},
				{
					Value: "Markdown",
					Key:   blogstore.POST_EDITOR_MARKDOWN,
				},
				{
					Value: "HTML Area (WYSIWYG)",
					Key:   blogstore.POST_EDITOR_HTMLAREA,
				},
				{
					Value: "Text Area",
					Key:   blogstore.POST_EDITOR_TEXTAREA,
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Admin Notes",
			Name:  "post_memo",
			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
			Value: data.formMemo,
			Help:  "Admin notes for this blogpost. These notes will not be visible to the public.",
		}),
		form.NewField(form.FieldOptions{
			Label:    "Post ID",
			Name:     "post_id",
			Type:     form.FORM_FIELD_TYPE_STRING,
			Value:    data.postID,
			Readonly: true,
			Help:     "The reference number (ID) of the post.",
		}),
		form.NewField(form.FieldOptions{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    data.view,
			Readonly: true,
		}),
	}

	editor := lo.IfF(data.post != nil, func() string { return data.post.Editor() }).Else("")

	fieldContent := form.NewField(form.FieldOptions{
		Label:   "Content",
		Name:    "post_content",
		Type:    form.FORM_FIELD_TYPE_TEXTAREA,
		Value:   data.formContent,
		Help:    "The content of this blog post to display on the post details page.",
		Options: []form.FieldOption{},
	})

	// For HTML Area editor, configure the Trumbowyg editor
	if editor == blogstore.POST_EDITOR_HTMLAREA {
		htmlAreaFieldOptions := []form.FieldOption{
			{
				Key: "config",
				Value: `{
	btns: [
		['viewHTML'],
		['undo', 'redo'],
		['formatting'],
		['strong', 'em', 'del'],
		['superscript', 'subscript'],
		['link','justifyLeft','justifyRight','justifyCenter','justifyFull'],
		['unorderedList', 'orderedList'],
		['insertImage'],
		['removeformat'],
		['horizontalRule'],
		['fullscreen'],
	],
	autogrow: true,
	removeformatPasted: true,
	tagsToRemove: ['script', 'link', 'embed', 'iframe', 'input'],
	tagsToKeep: ['hr', 'img', 'i'],
	autogrowOnEnter: true,
	linkTargets: ['_blank'],
	}`,
			}}
		fieldContent.Type = form.FORM_FIELD_TYPE_HTMLAREA
		fieldContent.Options = htmlAreaFieldOptions
	}

	if editor == blogstore.POST_EDITOR_BLOCKEDITOR {
		value := fieldContent.Value

		if value == "" {
			value = `[]`
		}

		editor, err := blockeditor.NewEditor(blockeditor.NewEditorOptions{
			// ID:    "blockeditor" + uid.HumanUid(),
			Name:  fieldContent.Name,
			Value: value,
			HandleEndpoint: links.NewAdminLinks().BlogPostUpdate(map[string]string{
				"post_id": data.postID,
				"action":  ACTION_BLOCKEDITOR_HANDLE,
			}),
			BlockDefinitions: blogblocks.BlockEditorDefinitions(),
		})

		if err != nil {
			return hb.Div().Class("alert alert-danger").Text("Error creating blockeditor: ").Text(err.Error())
		}

		fieldContent.Type = form.FORM_FIELD_TYPE_BLOCKEDITOR
		fieldContent.CustomInput = editor
	}

	fieldsContent := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Title",
			Name:  "post_title",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formTitle,
			Help:  "The title of this blog as will be seen everywhere",
		}),
		form.NewField(form.FieldOptions{
			Label: "Summary",
			Name:  "post_summary",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formSummary,
			Help:  "A short summary of this blog post to display on the post listing page.",
		}),
		fieldContent,
		form.NewField(form.FieldOptions{
			Label:    "Post ID",
			Name:     "post_id",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    data.postID,
			Readonly: true,
		}),
		form.NewField(form.FieldOptions{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    VIEW_CONTENT,
			Readonly: true,
		}),
	}

	if editor == blogstore.POST_EDITOR_BLOCKAREA {
		contentScript := hb.Script(`
setTimeout(() => {
	const textArea = document.querySelector('textarea[name="post_content"]');
	// get ID
	console.log(textArea.id);
	const blockArea = new BlockArea(textArea.id);
	blockArea.setParentId('` + data.postID + `');
	blockArea.registerBlock(BlockAreaHeading);
	blockArea.registerBlock(BlockAreaText);
	blockArea.registerBlock(BlockAreaImage);
	blockArea.registerBlock(BlockAreaCode);
	blockArea.registerBlock(BlockAreaRawHtml);
	blockArea.init();
}, 2000)
			`).
			ToHTML()

		fieldsContent = append(fieldsContent, form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: contentScript,
		}))
	}

	if editor == blogstore.POST_EDITOR_MARKDOWN {
		contentScript := hb.Script(`
setTimeout(() => {
	const textArea = document.querySelector('textarea[name="post_content"]');
	textArea.style.height = '300px';
}, 2000)
			`).
			ToHTML()

		fieldsContent = append(fieldsContent, form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: contentScript,
		}))
	}

	fieldsSEO := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Meta Description",
			Name:  "post_meta_description",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formMetaDescription,
			Help:  "The description of this blog as will be seen in search engines.",
		}),
		form.NewField(form.FieldOptions{
			Label: "Meta Keywords",
			Name:  "post_meta_keywords",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formMetaKeywords,
			Help:  "Specifies the keywords that will be used by the search engines to find this blog. Separate keywords with commas.",
		}),
		form.NewField(form.FieldOptions{
			Label: "Meta Robots",
			Name:  "post_meta_robots",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: data.formMetaRobots,
			Help:  "Specifies if this blog should be indexed by the search engines. Index, Follow, means all. NoIndex, NoFollow means none.",
			Options: []form.FieldOption{
				{
					Value: "- not selected -",
					Key:   "",
				},
				{
					Value: "INDEX, FOLLOW",
					Key:   "INDEX, FOLLOW",
				},
				{
					Value: "NOINDEX, FOLLOW",
					Key:   "NOINDEX, FOLLOW",
				},
				{
					Value: "INDEX, NOFOLLOW",
					Key:   "INDEX, NOFOLLOW",
				},
				{
					Value: "NOINDEX, NOFOLLOW",
					Key:   "NOINDEX, NOFOLLOW",
				},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Canonical URL",
			Name:  "post_canonical_url",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: data.formCanonicalURL,
			Help:  "The canonical URL for this blog. This is used by the search engines to display the preferred version of the blog post in search results.",
		}),
		form.NewField(form.FieldOptions{
			Label:    "Post ID",
			Name:     "post_id",
			Type:     form.FORM_FIELD_TYPE_STRING,
			Value:    data.postID,
			Readonly: true,
		}),
		form.NewField(form.FieldOptions{
			Label:    "View",
			Name:     "view",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    VIEW_SEO,
			Readonly: true,
		}),
	}

	formPostUpdate := form.NewForm(form.FormOptions{
		ID: "FormPostUpdate",
	})

	if data.view == VIEW_DETAILS {
		formPostUpdate.SetFields(fieldsDetails)
	}

	if data.view == VIEW_CONTENT {
		formPostUpdate.SetFields(fieldsContent)
	}

	if data.view == VIEW_SEO {
		formPostUpdate.SetFields(fieldsSEO)
	}

	if data.formErrorMessage != "" {
		formPostUpdate.AddField(form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{Icon: "error", Text: data.formErrorMessage}).ToHTML(),
		}))
	}

	if data.formSuccessMessage != "" {
		formPostUpdate.AddField(form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{Icon: "success", Text: data.formSuccessMessage}).ToHTML(),
		}))
	}

	return formPostUpdate.Build()

	// required := hb.Sup().HTML("required").Style("color:red;margin-left:10px;")

	// // Status
	// fomrGroupStatus := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Status").Child(required)).
	// 	Child(bs.FormSelect().
	// 		Name("post_status").
	// 		Child(bs.FormSelectOption("", "").
	// 			AttrIf(data.formStatus == "", "selected", "")).
	// 		Child(bs.FormSelectOption(blogstore.POST_STATUS_DRAFT, "Draft").
	// 			AttrIf(data.formStatus == blogstore.POST_STATUS_DRAFT, "selected", "selected")).
	// 		Child(bs.FormSelectOption(blogstore.POST_STATUS_PUBLISHED, "Published").
	// 			AttrIf(data.formStatus == blogstore.POST_STATUS_PUBLISHED, "selected", "selected")).
	// 		Child(bs.FormSelectOption(blogstore.POST_STATUS_UNPUBLISHED, "Unpublished").
	// 			AttrIf(data.formStatus == blogstore.POST_STATUS_UNPUBLISHED, "selected", "selected")).
	// 		Child(bs.FormSelectOption(blogstore.POST_STATUS_TRASH, "Trashed").
	// 			AttrIf(data.formStatus == blogstore.POST_STATUS_TRASH, "selected", "selected")),
	// 	)

	// // Admin Notes
	// formGroupMemo := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Admin Notes")).
	// 	Child(bs.FormTextArea().
	// 		Name("post_memo").
	// 		HTML(data.formMemo).
	// 		Style("height:100px;"),
	// 	)

	// // Post ID
	// formGroupPostId := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Post ID")).
	// 	Child(bs.FormInput().
	// 		Type(hb.TYPE_TEXT).
	// 		Name("post_id").
	// 		Value(data.postID).
	// 		Attr("readonly", ""),
	// 	)

	// // Title
	// formGroupTitle := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Title").Child(required)).
	// 	Child(bs.FormInput().
	// 		Type("text").
	// 		Name("post_title").
	// 		Value(data.formTitle).
	// 		Style("width:100%;"),
	// 	)

	// // Summary
	// formGroupSummary := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Summary")).
	// 	Child(bs.FormTextArea().
	// 		Type("text").
	// 		Name("post_summary").
	// 		HTML(data.formSummary).
	// 		Style("width:100%;"),
	// 	)

	// // Published Date
	// formGroupPublishedAt := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Published Date")).
	// 	Child(bs.FormInput().
	// 		Type(hb.TYPE_TEXT).
	// 		Name("post_published_at").
	// 		Value(data.formPublishedAt).
	// 		Style("width:100%;"),
	// 	)

	// // Featured
	// formGroupFeatured := bs.FormGroup().
	// 	Class("mb-3").
	// 	Child(bs.FormLabel("Featured")).
	// 	Child(bs.FormSelect().
	// 		Name("post_featured").
	// 		Child(bs.FormSelectOption("", "").
	// 			AttrIf(data.formFeatured == "", "selected", "")).
	// 		Child(bs.FormSelectOption("yes", "Yes").
	// 			AttrIf(data.formFeatured == "yes", "selected", "selected")).
	// 		Child(bs.FormSelectOption("no", "No").
	// 			AttrIf(data.formFeatured == "no", "selected", "selected")),
	// 	)

	// form := hb.Form().
	// 	ID("FormPostUpdate").
	// 	Child(formGroupTitle).
	// 	Child(fomrGroupStatus).
	// 	Child(formGroupSummary).
	// 	Child(formGroupPublishedAt).
	// 	Child(formGroupFeatured).
	// 	Child(formGroupMemo).
	// 	Child(formGroupPostId)

	// if data.formErrorMessage != "" {
	// 	form.Child(hb.Swal(hb.SwalOptions{Icon: "error", Text: data.formErrorMessage}))
	// }

	// if data.formSuccessMessage != "" {
	// 	form.Child(hb.Swal(hb.SwalOptions{Icon: "success", Text: data.formSuccessMessage}))
	// }

	// return form
}

func (controller postUpdateController) savePost(r *http.Request, data postUpdateControllerData) (d postUpdateControllerData, errorMessage string) {
	data.formCanonicalURL = req.Value(r, "post_canonical_url")
	data.formContent = req.Value(r, "post_content")
	data.formEditor = req.Value(r, "post_editor")
	data.formFeatured = req.Value(r, "post_featured")
	data.formImageUrl = req.Value(r, "post_image_url")
	data.formMemo = req.Value(r, "post_memo")
	data.formMetaDescription = req.Value(r, "post_meta_description")
	data.formMetaKeywords = req.Value(r, "post_meta_keywords")
	data.formMetaRobots = req.Value(r, "post_meta_robots")
	data.formPublishedAt = req.Value(r, "post_published_at")
	data.formSummary = req.Value(r, "post_summary")
	data.formStatus = req.Value(r, "post_status")
	data.formTitle = req.Value(r, "post_title")

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

	err := config.BlogStore.PostUpdate(data.post)

	if err != nil {
		config.Logger.Error("At postUpdateController > prepareDataAndValidate", slog.String("error", err.Error()))
		data.formErrorMessage = "System error. Saving post failed"
		return data, ""
	}

	data.formSuccessMessage = "Post saved successfully"

	return data, ""
}

func (controller postUpdateController) prepareDataAndValidate(r *http.Request) (data postUpdateControllerData, errorMessage string) {
	data.action = req.Value(r, "action")
	data.postID = req.Value(r, "post_id")
	data.view = req.ValueOr(r, "view", VIEW_DETAILS)

	if data.view == "" {
		data.view = VIEW_DETAILS
	}

	if data.postID == "" {
		return data, "Post ID is required"
	}

	var err error
	data.post, err = config.BlogStore.PostFindByID(data.postID)

	if err != nil {
		config.Logger.Error("At postUpdateController > prepareDataAndValidate", slog.String("error", err.Error()))
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
