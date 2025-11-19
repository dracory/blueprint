package post_update

import (
	"context"
	"log"
	"net/url"
	"strings"

	"project/internal/controllers/admin/blog/shared"
	livefluxctl "project/internal/controllers/liveflux"
	"project/internal/types"
	"project/pkg/blogblocks"

	"github.com/dracory/blockeditor"
	"github.com/dracory/blogstore"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/samber/lo"
)

type postContentComponent struct {
	liveflux.Base

	App types.AppInterface

	PostID string
	Post   *blogstore.Post

	FormTitle   string
	FormSummary string
	FormContent string

	Editor string

	FormErrorMessage   string
	FormSuccessMessage string
}

func NewPostContentComponent(app types.AppInterface) liveflux.ComponentInterface {
	inst, err := liveflux.New(&postContentComponent{})
	if err != nil {
		log.Println(err)
		return nil
	}

	if c, ok := inst.(*postContentComponent); ok {
		c.App = app
	}

	return inst
}

func (c *postContentComponent) GetKind() string {
	return "admin_blog_post_content"
}

func (c *postContentComponent) Mount(ctx context.Context, params map[string]string) error {
	if c.App == nil {
		if app, ok := ctx.Value(livefluxctl.AppContextKey).(types.AppInterface); ok {
			c.App = app
		}
	}

	c.PostID = strings.TrimSpace(params["post_id"])
	if c.PostID == "" {
		c.FormErrorMessage = "Post ID is required"
		return nil
	}

	if c.App == nil || c.App.GetBlogStore() == nil {
		c.FormErrorMessage = "Blog store not available"
		return nil
	}

	post, err := c.App.GetBlogStore().PostFindByID(c.PostID)
	if err != nil {
		c.App.GetLogger().Error("Error loading post for content component", "error", err.Error())
		c.FormErrorMessage = "Error loading post"
		return nil
	}

	if post == nil {
		c.FormErrorMessage = "Post not found"
		return nil
	}

	c.Post = post
	c.Editor = post.Editor()
	c.FormTitle = post.Title()
	c.FormSummary = post.Summary()
	c.FormContent = post.Content()

	return nil
}

func (c *postContentComponent) Handle(ctx context.Context, action string, data url.Values) error {
	switch action {
	case "save":
		if data == nil {
			data = url.Values{}
		}

		c.FormTitle = strings.TrimSpace(data.Get("post_title"))
		c.FormSummary = strings.TrimSpace(data.Get("post_summary"))
		c.FormContent = data.Get("post_content")

		if c.FormTitle == "" {
			c.FormErrorMessage = "Title is required"
			c.FormSuccessMessage = ""
			return nil
		}

		if c.App == nil || c.App.GetBlogStore() == nil {
			c.FormErrorMessage = "Blog store not available"
			c.FormSuccessMessage = ""
			return nil
		}

		post, err := c.App.GetBlogStore().PostFindByID(c.PostID)
		if err != nil || post == nil {
			c.FormErrorMessage = "Post not found"
			c.FormSuccessMessage = ""
			return nil
		}

		post.SetTitle(c.FormTitle)
		post.SetSummary(c.FormSummary)
		post.SetContent(c.FormContent)

		if err := c.App.GetBlogStore().PostUpdate(post); err != nil {
			c.App.GetLogger().Error("Error saving post content", "error", err.Error())
			c.FormErrorMessage = "System error. Saving post failed"
			c.FormSuccessMessage = ""
			return nil
		}

		c.FormErrorMessage = ""
		c.FormSuccessMessage = "Post saved successfully"
	}

	return nil
}

func (c *postContentComponent) Render(ctx context.Context) hb.TagInterface {
	frm := form.NewForm(form.FormOptions{
		ID: "FormPostUpdateContent",
	})

	// Determine editor type from loaded post
	editor := lo.IfF(c.Post != nil, func() string { return c.Post.Editor() }).Else("")

	// Base content field
	fieldContent := form.NewField(form.FieldOptions{
		Label:   "Content",
		Name:    "post_content",
		Type:    form.FORM_FIELD_TYPE_TEXTAREA,
		Value:   c.FormContent,
		Help:    "The content of this blog post to display on the post details page.",
		Options: []form.FieldOption{},
	})

	// HTML Area (Trumbowyg configuration)
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
			},
		}
		fieldContent.Type = form.FORM_FIELD_TYPE_HTMLAREA
		fieldContent.Options = htmlAreaFieldOptions
	}

	// BlockEditor
	if editor == blogstore.POST_EDITOR_BLOCKEDITOR {
		value := fieldContent.Value
		if value == "" {
			value = `[]`
		}

		be, err := blockeditor.NewEditor(blockeditor.NewEditorOptions{
			Name:  fieldContent.Name,
			Value: value,
			HandleEndpoint: shared.NewLinks().PostUpdate(map[string]string{
				"post_id": c.PostID,
				"action":  "blockeditor_handle",
			}),
			BlockDefinitions: blogblocks.BlockEditorDefinitions(),
		})
		if err == nil {
			fieldContent.Type = form.FORM_FIELD_TYPE_BLOCKEDITOR
			fieldContent.CustomInput = be
		}
	}

	fields := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Title",
			Name:  "post_title",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: c.FormTitle,
			Help:  "The title of this blog as will be seen everywhere",
		}),
		form.NewField(form.FieldOptions{
			Label: "Summary",
			Name:  "post_summary",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: c.FormSummary,
			Help:  "A short summary of this blog post to display on the post listing page.",
		}),
		fieldContent,
		form.NewField(form.FieldOptions{
			Label:    "Post ID",
			Name:     "post_id",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    c.PostID,
			Readonly: true,
		}),
	}

	// BlockArea script
	if editor == blogstore.POST_EDITOR_BLOCKAREA {
		contentScript := hb.Script(`
setTimeout(() => {
	const textArea = document.querySelector('textarea[name="post_content"]');
	if (!textArea) return;
	const blockArea = new BlockArea(textArea.id);
	blockArea.setParentId('` + c.PostID + `');
	blockArea.registerBlock(BlockAreaHeading);
	blockArea.registerBlock(BlockAreaText);
	blockArea.registerBlock(BlockAreaImage);
	blockArea.registerBlock(BlockAreaCode);
	blockArea.registerBlock(BlockAreaRawHtml);
	blockArea.init();
}, 2000)
		`).
			ToHTML()

		fields = append(fields, form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: contentScript,
		}))
	}

	// Markdown script
	if editor == blogstore.POST_EDITOR_MARKDOWN {
		contentScript := hb.Script(`
setTimeout(() => {
	const textArea = document.querySelector('textarea[name="post_content"]');
	if (!textArea) return;
	textArea.style.height = '300px';
}, 2000)
		`).
			ToHTML()

		fields = append(fields, form.NewField(form.FieldOptions{
			Type:  form.FORM_FIELD_TYPE_RAW,
			Value: contentScript,
		}))
	}

	frm.SetFields(fields)

	if c.FormErrorMessage != "" {
		frm.AddField(form.NewField(form.FieldOptions{
			Type: form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{
				Icon:             "error",
				Text:             c.FormErrorMessage,
				Position:         "top-end",
				Timer:            5000,
				TimerProgressBar: true,
			}).ToHTML(),
		}))
	}

	if c.FormSuccessMessage != "" {
		frm.AddField(form.NewField(form.FieldOptions{
			Type: form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{
				Icon:             "success",
				Text:             c.FormSuccessMessage,
				Position:         "top-end",
				Timer:            3000,
				TimerProgressBar: true,
			}).ToHTML(),
		}))
	}

	saveSpinner := hb.Span().
		Class("post-save-spinner spinner-border spinner-border-sm align-middle ms-2").
		Style("display: none;").
		Attr("role", "status").
		Child(hb.Span().Class("visually-hidden").Text("Loading"))

	saveBtn := hb.Button().
		Type("submit").
		Class("btn btn-primary").
		Attr(liveflux.DataFluxAction, "save").
		Attr(liveflux.DataFluxTargetKind, c.GetKind()).
		Attr(liveflux.DataFluxTargetID, c.GetID()).
		Attr(liveflux.DataFluxIndicator, "this, .post-save-spinner").
		Child(hb.I().Class("bi bi-save me-1")).
		Child(hb.Span().Text("Save")).
		Child(saveSpinner)

	content := hb.Div().
		Child(frm.Build()).
		Child(hb.Div().Class("mt-3 text-end").Child(saveBtn))

	return c.Root(content)
}

func init() {
	if err := liveflux.Register(&postContentComponent{}); err != nil {
		log.Printf("Failed to register postContentComponent: %v", err)
	}
}
