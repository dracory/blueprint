package post_update

import (
	"context"
	"log"
	"net/url"
	"strings"

	livefluxctl "project/internal/controllers/liveflux"
	"project/internal/types"

	"github.com/dracory/blogstore"
	"github.com/dracory/form"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
)

type postSEOComponent struct {
	liveflux.Base

	App types.AppInterface

	PostID string
	Post   *blogstore.Post

	FormCanonicalURL    string
	FormMetaDescription string
	FormMetaKeywords    string
	FormMetaRobots      string

	FormErrorMessage   string
	FormSuccessMessage string
}

func NewPostSEOComponent(app types.AppInterface) liveflux.ComponentInterface {
	inst, err := liveflux.New(&postSEOComponent{})
	if err != nil {
		log.Println(err)
		return nil
	}

	if c, ok := inst.(*postSEOComponent); ok {
		c.App = app
	}

	return inst
}

func (c *postSEOComponent) GetKind() string {
	return "admin_blog_post_seo"
}

func (c *postSEOComponent) Mount(ctx context.Context, params map[string]string) error {
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

	post, err := c.App.GetBlogStore().PostFindByID(ctx, c.PostID)
	if err != nil {
		c.App.GetLogger().Error("Error loading post for SEO component", "error", err.Error())
		c.FormErrorMessage = "Error loading post"
		return nil
	}

	if post == nil {
		c.FormErrorMessage = "Post not found"
		return nil
	}

	c.Post = post
	c.FormCanonicalURL = post.CanonicalURL()
	c.FormMetaDescription = post.MetaDescription()
	c.FormMetaKeywords = post.MetaKeywords()
	c.FormMetaRobots = post.MetaRobots()

	return nil
}

func (c *postSEOComponent) Handle(ctx context.Context, action string, data url.Values) error {
	switch action {
	case "save":
		if data == nil {
			data = url.Values{}
		}

		c.FormCanonicalURL = strings.TrimSpace(data.Get("post_canonical_url"))
		c.FormMetaDescription = strings.TrimSpace(data.Get("post_meta_description"))
		c.FormMetaKeywords = strings.TrimSpace(data.Get("post_meta_keywords"))
		c.FormMetaRobots = strings.TrimSpace(data.Get("post_meta_robots"))

		if c.App == nil || c.App.GetBlogStore() == nil {
			c.FormErrorMessage = "Blog store not available"
			c.FormSuccessMessage = ""
			return nil
		}

		post, err := c.App.GetBlogStore().PostFindByID(ctx, c.PostID)
		if err != nil || post == nil {
			c.FormErrorMessage = "Post not found"
			c.FormSuccessMessage = ""
			return nil
		}

		post.SetCanonicalURL(c.FormCanonicalURL)
		post.SetMetaDescription(c.FormMetaDescription)
		post.SetMetaKeywords(c.FormMetaKeywords)
		post.SetMetaRobots(c.FormMetaRobots)

		if err := c.App.GetBlogStore().PostUpdate(ctx, post); err != nil {
			c.App.GetLogger().Error("Error saving post SEO", "error", err.Error())
			c.FormErrorMessage = "System error. Saving post failed"
			c.FormSuccessMessage = ""
			return nil
		}

		c.FormErrorMessage = ""
		c.FormSuccessMessage = "Post saved successfully"
	}

	return nil
}

func (c *postSEOComponent) Render(ctx context.Context) hb.TagInterface {
	frm := form.NewForm(form.FormOptions{
		ID: "FormPostUpdateSEO",
	})

	fields := []form.FieldInterface{
		form.NewField(form.FieldOptions{
			Label: "Meta Description",
			Name:  "post_meta_description",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: c.FormMetaDescription,
		}),
		form.NewField(form.FieldOptions{
			Label: "Meta Keywords",
			Name:  "post_meta_keywords",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: c.FormMetaKeywords,
		}),
		form.NewField(form.FieldOptions{
			Label: "Meta Robots",
			Name:  "post_meta_robots",
			Type:  form.FORM_FIELD_TYPE_SELECT,
			Value: c.FormMetaRobots,
			Options: []form.FieldOption{
				{Key: "", Value: "- not selected -"},
				{Key: "INDEX, FOLLOW", Value: "INDEX, FOLLOW"},
				{Key: "NOINDEX, FOLLOW", Value: "NOINDEX, FOLLOW"},
				{Key: "INDEX, NOFOLLOW", Value: "INDEX, NOFOLLOW"},
				{Key: "NOINDEX, NOFOLLOW", Value: "NOINDEX, NOFOLLOW"},
			},
		}),
		form.NewField(form.FieldOptions{
			Label: "Canonical URL",
			Name:  "post_canonical_url",
			Type:  form.FORM_FIELD_TYPE_STRING,
			Value: c.FormCanonicalURL,
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
	if err := liveflux.Register(&postSEOComponent{}); err != nil {
		log.Printf("Failed to register postSEOComponent: %v", err)
	}
}
