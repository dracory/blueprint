package aiposteditor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"project/internal/controllers/admin/blog/ai_post_editor/templates"
	"project/internal/controllers/admin/blog/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"
	"project/pkg/blogai"

	"github.com/dracory/base/req"
	"github.com/dracory/cdn"
	"github.com/dracory/customstore"
	"github.com/dracory/hb"
)

const (
	ACTION_REGENERATE_SECTION   = "regenerate_section"
	ACTION_REGENERATE_IMAGE     = "regenerate_image"
	ACTION_CREATE_FINAL_POST    = "create_final_post"
	ACTION_SAVE_DRAFT           = "save_draft"
	ACTION_REGENERATE_PARAGRAPH = "regenerate_paragraph"
	ACTION_LOAD_POST            = "load_post"
	ACTION_REGENERATE_SUMMARY   = "regenerate_summary"
	ACTION_REGENERATE_METAS     = "regenerate_metas"
)

type AiPostEditorController struct {
	app types.AppInterface
}

type pageData struct {
	Request    *http.Request
	BlogAiPost blogai.RecordPost
	Record     customstore.RecordInterface
}

func NewAiPostEditorController(app types.AppInterface) *AiPostEditorController {
	return &AiPostEditorController{app: app}
}

func (c *AiPostEditorController) Handler(w http.ResponseWriter, r *http.Request) string {
	c.app.GetLogger().Info("Post Editor Handler called")

	data, errorMessage := c.prepareDataAndValidate(r)
	if errorMessage != "" {
		return helpers.ToFlashError(c.app.GetCacheStore(), w, r, errorMessage, shared.NewLinks().Home(), 10)
	}

	action := req.Value(r, "action")
	switch {
	case r.Method == http.MethodPost && action == ACTION_REGENERATE_SECTION:
		return c.onRegenerateSection(data)
	case r.Method == http.MethodPost && action == ACTION_REGENERATE_IMAGE:
		return c.onRegenerateImage(data)
	case r.Method == http.MethodPost && action == ACTION_REGENERATE_PARAGRAPH:
		return c.onRegenerateParagraph(data)
	case r.Method == http.MethodPost && action == ACTION_CREATE_FINAL_POST:
		return c.onCreateFinalPost(data)
	case r.Method == http.MethodPost && action == ACTION_SAVE_DRAFT:
		return c.onSaveDraft(data)
	case r.Method == http.MethodPost && action == ACTION_LOAD_POST:
		return c.onLoadPost(data)
	case r.Method == http.MethodPost && action == ACTION_REGENERATE_SUMMARY:
		return c.onRegenerateSummary(data)
	case r.Method == http.MethodPost && action == ACTION_REGENERATE_METAS:
		return c.onRegenerateMetas(data)
	}

	return layouts.NewAdminLayout(c.app, r, layouts.Options{
		Title:   "Edit & Save Blog Post",
		Content: c.view(data),
		ScriptURLs: []string{
			cdn.Sweetalert2_11(),
			cdn.VueJs_3(),
		},
		StyleURLs: []string{},
	}).ToHTML()
}

func (c *AiPostEditorController) buildPostMarkdownContent(_ *http.Request, record *blogai.RecordPost) string {
	content := "# " + record.Title + "\n\n"

	content += "## " + record.Introduction.Title + "\n\n"
	for _, paragraph := range record.Introduction.Paragraphs {
		content += paragraph + "\n\n"
	}

	for _, section := range record.Sections {
		content += "## " + section.Title + "\n\n"
		for _, paragraph := range section.Paragraphs {
			content += paragraph + "\n\n"
		}
	}

	content += "## " + record.Conclusion.Title + "\n\n"
	for _, paragraph := range record.Conclusion.Paragraphs {
		content += paragraph + "\n\n"
	}

	return content
}

func (c *AiPostEditorController) view(data pageData) *hb.Tag {
	header := hb.Heading1().HTML("Edit Blog Post")

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{
			Name: "Dashboard",
			URL:  links.Admin().Home(),
		},
		{
			Name: "Blog",
			URL:  links.Admin().Blog(),
		},
		{
			Name: "AI Tools",
			URL:  shared.NewLinks().AiTools(),
		},
		{
			Name: "Post Editor",
			URL:  shared.NewLinks().AiPostEditor(),
		},
		{
			Name: "Post Editor",
			URL:  shared.NewLinks().AiPostEditor(map[string]string{"id": data.BlogAiPost.ID}),
		},
	})

	backButton := hb.A().
		Class("btn btn-secondary me-3").
		Href(shared.NewLinks().AiPostGenerator(map[string]string{})).
		HTML("← Back to Post Generator")

	vueApp := hb.Raw(templates.Tpl("app.html", map[string]any{}))
	vueScript := hb.Script(templates.Tpl("app.js", map[string]any{
		"postJSON": data.BlogAiPost.ToJSON(),
		"id":       data.BlogAiPost.ID,
		"url":      shared.NewLinks().AiPostEditor(map[string]string{}),
	}))
	vueStyles := hb.Style(templates.Tpl("app.css", nil))

	return hb.Div().
		Class("container min-vh-100 py-4 bg-light").
		Child(breadcrumbs).
		Child(hb.Div().
			Class("container").
			Child(header).
			Child(hb.Div().
				Class("d-flex justify-content-between mb-3").
				Child(backButton),
			).
			Child(vueApp).
			Child(vueScript).
			Child(vueStyles),
		)
}

func (c *AiPostEditorController) prepareDataAndValidate(r *http.Request) (pageData, string) {
	var (
		data pageData
		err  error
	)

	data.Request = r
	recordPostID := req.Value(r, "id")
	if recordPostID == "" {
		return data, "Record Post ID is missing"
	}

	data.Record, err = c.app.GetCustomStore().RecordFindByID(recordPostID)
	if err != nil {
		c.app.GetLogger().Error("BlogAi. Post Editor. Prepare Data. Error finding record post", slog.String("error", err.Error()))
		return data, fmt.Sprintf("Failed to find record post: %s", err)
	}

	if data.Record == nil {
		c.app.GetLogger().Error("BlogAi. Post Editor. Prepare Data. Post record not found", slog.String("record_id", recordPostID))
		return data, "Post record not found"
	}

	if data.Record.Type() != blogai.POST_RECORD_TYPE {
		c.app.GetLogger().Error("BlogAi. Post Editor. Prepare Data. Invalid record type", slog.String("record_type", data.Record.Type()), slog.String("record_id", recordPostID))
		return data, "Invalid record type"
	}

	data.BlogAiPost, err = blogai.NewRecordPostFromCustomRecord(data.Record)
	if err != nil {
		c.app.GetLogger().Error("BlogAi. Post Editor. Prepare Data. Failed to parse blog record", slog.String("error", err.Error()))
		return data, fmt.Sprintf("Failed to parse blog record: %s", err)
	}

	return data, ""
}

func RecordFromJSON(jsonStr string) (*blogai.RecordPost, error) {
	type postData struct {
		Title           string                         `json:"title"`
		Subtitle        string                         `json:"subtitle,omitempty"`
		Summary         string                         `json:"summary,omitempty"`
		Introduction    blogai.PostContentIntroduction `json:"introduction"`
		Sections        []blogai.PostContentSection    `json:"sections"`
		Conclusion      blogai.PostContentConclusion   `json:"conclusion"`
		Keywords        []string                       `json:"keywords,omitempty"`
		MetaDescription string                         `json:"metaDescription,omitempty"`
		MetaKeywords    []string                       `json:"metaKeywords,omitempty"`
		MetaTitle       string                         `json:"metaTitle,omitempty"`
		Image           string                         `json:"image,omitempty"`
	}

	var record postData
	if err := json.Unmarshal([]byte(jsonStr), &record); err != nil {
		return nil, err
	}

	recordPost := blogai.RecordPost{}
	recordPost.Title = record.Title
	recordPost.Subtitle = record.Subtitle
	recordPost.Summary = record.Summary
	recordPost.Introduction = record.Introduction
	recordPost.Sections = record.Sections
	recordPost.Conclusion = record.Conclusion
	recordPost.Keywords = record.Keywords
	recordPost.MetaDescription = record.MetaDescription
	recordPost.MetaKeywords = record.MetaKeywords
	recordPost.MetaTitle = record.MetaTitle
	recordPost.Image = record.Image

	return &recordPost, nil
}
