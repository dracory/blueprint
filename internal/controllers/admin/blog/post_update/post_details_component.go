package post_update

import (
	"context"
	"log"
	"net/url"
	"strings"

	"project/internal/controllers/admin/blog/shared"
	livefluxctl "project/internal/controllers/liveflux"
	"project/internal/registry"
	"project/pkg/blogai"

	"github.com/dracory/blogstore"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dromara/carbon/v2"
	"github.com/samber/lo"
)

type postDetailsComponent struct {
	liveflux.Base

	registry registry.RegistryInterface

	PostID string
	Post   *blogstore.Post

	FormStatus      string
	FormImageUrl    string
	FormFeatured    string
	FormPublishedAt string
	FormEditor      string
	FormMemo        string

	FormErrorMessage   string
	FormSuccessMessage string
}

func NewPostDetailsComponent(registry registry.RegistryInterface) liveflux.ComponentInterface {
	inst, err := liveflux.New(&postDetailsComponent{})
	if err != nil {
		log.Println(err)
		return nil
	}

	if c, ok := inst.(*postDetailsComponent); ok {
		c.registry = registry
	}

	return inst
}

func (c *postDetailsComponent) GetKind() string {
	return "admin_blog_post_details"
}

func (c *postDetailsComponent) Mount(ctx context.Context, params map[string]string) error {
	if c.registry == nil {
		if app, ok := ctx.Value(livefluxctl.AppContextKey).(registry.RegistryInterface); ok {
			c.registry = app
		}
	}

	c.PostID = strings.TrimSpace(params["post_id"])
	if c.PostID == "" {
		c.FormErrorMessage = "Post ID is required"
		return nil
	}

	if c.registry == nil || c.registry.GetBlogStore() == nil {
		c.FormErrorMessage = "Blog store not available"
		return nil
	}

	post, err := c.registry.GetBlogStore().PostFindByID(ctx, c.PostID)
	if err != nil {
		c.registry.GetLogger().Error("Error loading post for details component", "error", err.Error())
		c.FormErrorMessage = "Error loading post"
		return nil
	}

	if post == nil {
		c.FormErrorMessage = "Post not found"
		return nil
	}

	c.Post = post
	c.FormStatus = post.Status()
	c.FormImageUrl = post.ImageUrl()
	c.FormFeatured = post.Featured()
	c.FormPublishedAt = lo.Ternary(post.PublishedAtCarbon() != nil, post.PublishedAtCarbon().ToDateTimeString(), "")
	c.FormEditor = post.Editor()
	c.FormMemo = post.Memo()

	return nil
}

func (c *postDetailsComponent) Handle(ctx context.Context, action string, data url.Values) error {
	switch action {
	case "save":
		if data == nil {
			data = url.Values{}
		}

		c.FormStatus = strings.TrimSpace(data.Get("post_status"))
		c.FormImageUrl = strings.TrimSpace(data.Get("post_image_url"))
		c.FormFeatured = strings.TrimSpace(data.Get("post_featured"))
		c.FormPublishedAt = strings.TrimSpace(data.Get("post_published_at"))
		c.FormEditor = strings.TrimSpace(data.Get("post_editor"))
		c.FormMemo = strings.TrimSpace(data.Get("post_memo"))

		if c.FormStatus == "" {
			c.FormErrorMessage = "Status is required"
			c.FormSuccessMessage = ""
			return nil
		}

		if c.registry == nil || c.registry.GetBlogStore() == nil {
			c.FormErrorMessage = "Blog store not available"
			c.FormSuccessMessage = ""
			return nil
		}

		post, err := c.registry.GetBlogStore().PostFindByID(ctx, c.PostID)
		if err != nil || post == nil {
			c.FormErrorMessage = "Post not found"
			c.FormSuccessMessage = ""
			return nil
		}

		// normalize published_at similar to v1
		publishedAtStr := strings.ReplaceAll(c.FormPublishedAt, " ", "T") + ":00"
		publishedAtStr = lo.Substring(publishedAtStr, 0, 19)
		var publishedAt string
		if strings.TrimSpace(c.FormPublishedAt) == "" {
			publishedAt = "0000-00-00 00:00:00"
		} else {
			publishedAt = carbon.Parse(publishedAtStr).ToDateTimeString(carbon.UTC)
		}

		post.SetEditor(c.FormEditor)
		post.SetFeatured(c.FormFeatured)
		post.SetImageUrl(c.FormImageUrl)
		post.SetMemo(c.FormMemo)
		post.SetPublishedAt(publishedAt)
		post.SetStatus(c.FormStatus)

		if err := c.registry.GetBlogStore().PostUpdate(ctx, post); err != nil {
			c.registry.GetLogger().Error("Error saving post details", "error", err.Error())
			c.FormErrorMessage = "System error. Saving post failed"
			c.FormSuccessMessage = ""
			return nil
		}

		c.FormErrorMessage = ""
		c.FormSuccessMessage = "Post saved successfully"
	case "regenerate_image":
		if c.registry == nil || c.registry.GetBlogStore() == nil {
			c.FormErrorMessage = "Blog store not available"
			c.FormSuccessMessage = ""
			return nil
		}

		post, err := c.registry.GetBlogStore().PostFindByID(ctx, c.PostID)
		if err != nil || post == nil {
			c.FormErrorMessage = "Post not found"
			c.FormSuccessMessage = ""
			return nil
		}

		agent := blogai.NewBlogWriterAgent(c.registry.GetLogger())
		if agent == nil {
			c.FormErrorMessage = "Failed to initialize AI engine"
			c.FormSuccessMessage = ""
			return nil
		}

		llmEngine, err := shared.LlmEngine(c.registry)
		if err != nil || llmEngine == nil {
			c.FormErrorMessage = "Failed to initialize AI engine"
			c.FormSuccessMessage = ""
			return nil
		}

		imageURL, err := agent.GenerateImage(llmEngine, post.Title(), post.Summary())
		if err != nil {
			c.registry.GetLogger().Error("BlogAi.PostUpdateV2.RegenerateImage", "error", err.Error())
			c.FormErrorMessage = "Failed to generate image"
			c.FormSuccessMessage = ""
			return nil
		}

		post.SetImageUrl(imageURL)
		if err := c.registry.GetBlogStore().PostUpdate(ctx, post); err != nil {
			c.registry.GetLogger().Error("BlogAi.PostUpdateV2.RegenerateImage.Save", "error", err.Error())
			c.FormErrorMessage = "Failed to save generated image"
			c.FormSuccessMessage = ""
			return nil
		}

		c.FormImageUrl = imageURL
		c.FormErrorMessage = ""
		c.FormSuccessMessage = "Image regenerated successfully"
	}

	return nil
}

func (c *postDetailsComponent) Render(ctx context.Context) hb.TagInterface {
	frm := form.NewForm(form.FormOptions{
		ID: "FormPostUpdateDetails",
	})

	fields := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Status",
			Name:  "post_status",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: c.FormStatus,
			Options: []form.FieldOption{
				{Key: "", Value: "- not selected -"},
				{Key: blogstore.POST_STATUS_DRAFT, Value: "Draft"},
				{Key: blogstore.POST_STATUS_PUBLISHED, Value: "Published"},
				{Key: blogstore.POST_STATUS_UNPUBLISHED, Value: "Unpublished"},
				{Key: blogstore.POST_STATUS_TRASH, Value: "In Trash Bin"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Image URL",
			Name:  "post_image_url",
			Type:  form.FORM_FIELD_TYPE_IMAGE,
			Value: c.FormImageUrl,
		}),
		form.NewField(form.FieldOptions{
			Label: "Featured",
			Name:  "post_featured",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: c.FormFeatured,
			Options: []form.FieldOption{
				{Key: "", Value: "- not selected -"},
				{Key: "no", Value: "No"},
				{Key: "yes", Value: "Yes"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Published At",
			Name:  "post_published_at",
			Type:  form.FORM_FIELD_TYPE_DATETIME,
			Value: c.FormPublishedAt,
		}),
		form.NewField(form.FieldOptions{
			Label: "Editor",
			Name:  "post_editor",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: c.FormEditor,
			Options: []form.FieldOption{
				{Key: "", Value: "- not selected -"},
				{Key: blogstore.POST_EDITOR_BLOCKAREA, Value: "BlockArea"},
				{Key: blogstore.POST_EDITOR_BLOCKEDITOR, Value: "BlockEditor"},
				{Key: blogstore.POST_EDITOR_MARKDOWN, Value: "Markdown"},
				{Key: PostEditorMarkdownEasyMDE, Value: "Markdown (EasyMDE)"},
				{Key: blogstore.POST_EDITOR_HTMLAREA, Value: "HTML Area (WYSIWYG)"},
				{Key: blogstore.POST_EDITOR_TEXTAREA, Value: "Text Area"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Admin Notes",
			Name:  "post_memo",
			Type:  form.FORM_FIELD_TYPE_TEXTAREA,
			Value: c.FormMemo,
		}),
		form.NewField(form.FieldOptions{
			Label:    "Post ID",
			Name:     "post_id",
			Type:     form.FORM_FIELD_TYPE_HIDDEN,
			Value:    c.PostID,
			Readonly: true,
		}),
	}

	frm.SetFields(fields)

	if c.FormErrorMessage != "" {
		frm.AddField(form.NewField(form.FieldOptions{
			Type: form.FORM_FIELD_TYPE_RAW,
			Value: hb.Swal(hb.SwalOptions{
				Icon:             "error",
				Text:             c.FormErrorMessage,
				Position:         "top-end",
				Timer:            3000,
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

	regenSpinner := hb.Span().
		Class("post-image-spinner spinner-border spinner-border-sm align-middle ms-2").
		Style("display: none;").
		Attr("role", "status").
		Child(hb.Span().Class("visually-hidden").Text("Loading"))

	// Visible button that shows confirmation dialog
	regenBtn := hb.Button().
		Type("button").
		Class("btn btn-outline-warning me-2 js-regenerate-image").
		Child(hb.I().Class("bi bi-magic me-1")).
		Child(hb.Span().Text("Regenerate Image")).
		Child(regenSpinner)

	// Hidden action button that actually triggers the Liveflux action
	regenActionBtn := hb.Button().
		Type("button").
		Class("d-none js-regenerate-image-action").
		Attr(liveflux.DataFluxAction, "regenerate_image").
		Attr(liveflux.DataFluxTargetKind, c.GetKind()).
		Attr(liveflux.DataFluxTargetID, c.GetID()).
		Attr(liveflux.DataFluxIndicator, "this, .post-image-spinner")

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

	// Attach SweetAlert confirmation for regenerate image
	confirmScript := hb.Script(`
document.addEventListener('click', function(e) {
	const btn = e.target.closest('.js-regenerate-image');
	if (!btn || typeof Swal === 'undefined') return;
	Swal.fire({
		title: 'Regenerate image?',
		text: 'This will replace the current image URL for this post.',
		icon: 'warning',
		showCancelButton: true,
		confirmButtonText: 'Yes, regenerate',
		cancelButtonText: 'Cancel'
	}).then((result) => {
		if (result.isConfirmed) {
			const actionBtn = document.querySelector('.js-regenerate-image-action');
			if (actionBtn) { actionBtn.click(); }
		}
	});
});

document.addEventListener('click', function(e) {
	const toggle = e.target.closest('.js-toggle-advanced-tools');
	if (!toggle) return;
	const panel = document.querySelector('.js-advanced-tools-panel');
	if (!panel) return;
	panel.classList.toggle('d-none');
});
	`).
		ToHTML()

	advancedToggleBtn := hb.Button().
		Type("button").
		Class("btn btn-link btn-sm js-toggle-advanced-tools").
		Child(hb.I().Class("bi bi-tools me-1")).
		Child(hb.Span().Text("Toggle Advanced Tools"))

	footerMain := hb.Div().
		Class("mt-3 d-flex justify-content-between align-items-center").
		Child(hb.Div().Child(advancedToggleBtn)).
		Child(hb.Div().Child(saveBtn))

	advancedPanel := hb.Div().
		Class("mt-3 js-advanced-tools-panel d-none").
		Child(
			hb.Div().
				Class("border rounded p-2 bg-dark bg-opacity-25").
				Child(
					hb.Div().
						Class("small text-muted mb-2").
						Child(hb.I().Class("bi bi-tools me-1")).
						Text("Advanced Tools"),
				).
				Child(
					hb.Div().
						Class("d-flex justify-content-start").
						Child(regenBtn).
						Child(regenActionBtn),
				),
		)

	content := hb.Div().
		Child(frm.Build()).
		Child(footerMain).
		Child(advancedPanel).
		Child(hb.Raw(confirmScript))

	return c.Root(content)
}

func init() {
	if err := liveflux.Register(&postDetailsComponent{}); err != nil {
		log.Printf("Failed to register postDetailsComponent: %v", err)
	}
}
