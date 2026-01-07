package aipostcontentupdate

import (
	"net/http"
	"strings"

	"project/internal/controllers/admin/blog/shared"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/registry"

	"github.com/dracory/cdn"
	"github.com/dracory/liveflux"
	"github.com/dracory/req"
)

type Controller struct {
	registry registry.RegistryInterface
}

func NewController(registry registry.RegistryInterface) *Controller {
	return &Controller{registry: registry}
}

func (c *Controller) Handler(w http.ResponseWriter, r *http.Request) string {
	postID := req.GetStringTrimmed(r, "post_id")
	if strings.TrimSpace(postID) == "" {
		return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, "Post ID is required", shared.NewLinks().PostManager(), 10)
	}

	component := NewFormAiPostContentUpdate(c.registry)
	if component == nil {
		return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, "Failed to initialize AI content editor", shared.NewLinks().PostManager(), 10)
	}

	rendered := liveflux.SSR(component, map[string]string{
		"post_id": postID,
	})
	if rendered == nil {
		return helpers.ToFlashError(c.registry.GetCacheStore(), w, r, "Error rendering AI content editor", shared.NewLinks().PostManager(), 10)
	}

	return layouts.NewAdminLayout(c.registry, r, layouts.Options{
		Title:   "Edit Post Content",
		Content: rendered,
		ScriptURLs: []string{
			cdn.Sweetalert2_11(),
			"https://cdn.jsdelivr.net/npm/sortablejs@1.15.0/Sortable.min.js",
		},
		Scripts: []string{
			liveflux.Script().ToHTML(),
		},
	}).ToHTML()
}

// func (c *Controller) prepareData(r *http.Request) (pageData, string) {
// 	var data pageData

// 	postID := req.GetStringTrimmed(r, "post_id")
// 	if postID == "" {
// 		return data, "Post ID is required"
// 	}

// 	post, err := c.app.GetBlogStore().PostFindByID(postID)
// 	if err != nil {
// 		c.app.GetLogger().Error("BlogAi.ContentUpdate.PrepareData. Post find failed", slog.String("error", err.Error()))
// 		return data, "Post not found"
// 	}
// 	if post == nil {
// 		return data, "Post not found"
// 	}

// 	data.Post = post
// 	data.Request = r

// 	record := MarkdownToRecordPost(post.Content(), post.Title())
// 	record.ID = post.ID()
// 	record.Summary = post.Summary()
// 	record.MetaDescription = post.MetaDescription()
// 	record.MetaTitle = post.Title()
// 	record.MetaKeywords = keywordsToSlice(post.MetaKeywords())
// 	record.Keywords = record.MetaKeywords
// 	record.Image = post.ImageUrl()

// 	data.Record = record

// 	return data, ""
// }

// func (c *Controller) view(data pageData) *hb.Tag {
// 	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
// 		{Name: "Dashboard", URL: links.Admin().Home()},
// 		{Name: "Blog", URL: links.Admin().Blog()},
// 		{Name: "Post Manager", URL: shared.NewLinks().PostManager()},
// 		{Name: "Edit Post", URL: shared.NewLinks().PostUpdate(map[string]string{"post_id": data.Post.ID()})},
// 		{Name: "AI Content Editor"},
// 	})

// 	header := hb.Heading1().HTML("AI Post Content Editor")
// 	backButton := hb.A().
// 		Class("btn btn-secondary").
// 		Href(shared.NewLinks().PostUpdate(map[string]string{"post_id": data.Post.ID()})).
// 		HTML("← Back to Post")

// 	return hb.Div().
// 		Class("container py-4").
// 		Child(breadcrumbs).
// 		Child(
// 			hb.Div().
// 				Class("d-flex justify-content-between align-items-center mb-3").
// 				Child(header).
// 				Child(backButton),
// 		)
// }

// func keywordsToSlice(keywords string) []string {
// 	if strings.TrimSpace(keywords) == "" {
// 		return []string{}
// 	}
// 	parts := strings.Split(keywords, ",")
// 	return lo.FilterMap(parts, func(part string, _ int) (string, bool) {
// 		trimmed := strings.TrimSpace(part)
// 		return trimmed, trimmed != ""
// 	})
// }

// func sliceToKeywords(items []string) string {
// 	cleaned := lo.FilterMap(items, func(item string, _ int) (string, bool) {
// 		trimmed := strings.TrimSpace(item)
// 		return trimmed, trimmed != ""
// 	})
// 	return strings.Join(cleaned, ", ")
// }

// func (c *Controller) onLoadPost(data pageData) string {
// 	return api.SuccessWithData("Post loaded", map[string]any{
// 		"post":   data.Record,
// 		"postID": data.Post.ID(),
// 	}).ToString()
// }

// func (c *Controller) onRegenerateContent(data pageData) string {
// 	record, err := c.parseRecordPostFromRequest(data)
// 	if err != nil {
// 		return api.Error(err.Error()).ToString()
// 	}

// 	target := req.GetStringTrimmed(data.Request, "target")
// 	if target == "" {
// 		target = "section"
// 	}

// 	agent := blogai.NewBlogWriterAgent(c.app.GetLogger())
// 	if agent == nil {
// 		return api.Error("failed to initialize LLM engine").ToString()
// 	}

// 	llmEngine, err := shared.LlmEngine(c.app)
// 	if err != nil {
// 		return api.Error("failed to initialize LLM engine: " + err.Error()).ToString()
// 	}
// 	if llmEngine == nil {
// 		return api.Error("failed to initialize LLM engine").ToString()
// 	}

// 	switch target {
// 	case "section":
// 		section := req.GetStringTrimmed(data.Request, "section")
// 		if section == "" {
// 			return api.Error("Section parameter is missing").ToString()
// 		}

// 		updatedRecord, regenErr := agent.RegenerateSection(llmEngine, record, section)
// 		if regenErr != nil {
// 			c.app.GetLogger().Error("BlogAi.ContentUpdate.RegenerateSection", slog.String("error", regenErr.Error()), slog.String("section", section))
// 			return api.Error("Failed to regenerate section: " + regenErr.Error()).ToString()
// 		}

// 		updatedRecord = normalizeRecord(updatedRecord, data.Post.ID())
// 		sectionData, payloadErr := sectionPayload(updatedRecord, section)
// 		if payloadErr != nil {
// 			return api.Error(payloadErr.Error()).ToString()
// 		}

// 		return api.SuccessWithData("Section regenerated successfully", map[string]any{
// 			"post":    updatedRecord,
// 			"section": sectionData,
// 		}).ToString()

// 	case "paragraph":
// 		sectionType := req.GetStringTrimmed(data.Request, "section_type")
// 		if sectionType == "" {
// 			return api.Error("Section type is required").ToString()
// 		}

// 		sectionIndex := cast.ToInt(req.GetStringTrimmed(data.Request, "section_index"))
// 		paragraphIndex := cast.ToInt(req.GetStringTrimmed(data.Request, "paragraph_index"))

// 		newParagraph, regenErr := agent.RegenerateParagraph(llmEngine, record, sectionType, sectionIndex, paragraphIndex)
// 		if regenErr != nil {
// 			c.app.GetLogger().Error("BlogAi.ContentUpdate.RegenerateParagraph", slog.String("error", regenErr.Error()), slog.String("section_type", sectionType))
// 			return api.Error("Failed to regenerate paragraph: " + regenErr.Error()).ToString()
// 		}

// 		newParagraph = strings.TrimSpace(newParagraph)
// 		if newParagraph == "" {
// 			return api.Error("Generated paragraph is empty").ToString()
// 		}

// 		switch sectionType {
// 		case "introduction":
// 			if paragraphIndex < 0 || paragraphIndex > len(record.Introduction.Paragraphs) {
// 				return api.Error("Invalid introduction paragraph index").ToString()
// 			}
// 			if paragraphIndex == len(record.Introduction.Paragraphs) {
// 				record.Introduction.Paragraphs = append(record.Introduction.Paragraphs, newParagraph)
// 			} else {
// 				record.Introduction.Paragraphs[paragraphIndex] = newParagraph
// 			}
// 		case "conclusion":
// 			if paragraphIndex < 0 || paragraphIndex > len(record.Conclusion.Paragraphs) {
// 				return api.Error("Invalid conclusion paragraph index").ToString()
// 			}
// 			if paragraphIndex == len(record.Conclusion.Paragraphs) {
// 				record.Conclusion.Paragraphs = append(record.Conclusion.Paragraphs, newParagraph)
// 			} else {
// 				record.Conclusion.Paragraphs[paragraphIndex] = newParagraph
// 			}
// 		case "section":
// 			if sectionIndex < 0 || sectionIndex >= len(record.Sections) {
// 				return api.Error("Invalid section index").ToString()
// 			}
// 			if paragraphIndex < 0 || paragraphIndex > len(record.Sections[sectionIndex].Paragraphs) {
// 				return api.Error("Invalid section paragraph index").ToString()
// 			}
// 			if paragraphIndex == len(record.Sections[sectionIndex].Paragraphs) {
// 				record.Sections[sectionIndex].Paragraphs = append(record.Sections[sectionIndex].Paragraphs, newParagraph)
// 			} else {
// 				record.Sections[sectionIndex].Paragraphs[paragraphIndex] = newParagraph
// 			}
// 		default:
// 			return api.Error("Unsupported section type").ToString()
// 		}

// 		record = normalizeRecord(record, data.Post.ID())

// 		return api.SuccessWithData("Paragraph regenerated successfully", map[string]any{
// 			"post":      record,
// 			"paragraph": newParagraph,
// 		}).ToString()
// 	}

// 	return api.Error("Unsupported regeneration target").ToString()
// }

// func (c *Controller) onSaveDraft(data pageData) string {
// 	record, err := c.parseRecordPostFromRequest(data)
// 	if err != nil {
// 		return api.Error(err.Error()).ToString()
// 	}

// 	if record.Title == "" {
// 		return api.Error("Title is required").ToString()
// 	}

// 	record.Status = blogai.POST_STATUS_DRAFT

// 	if err := c.persistRecordToPost(data.Post, record, blogstore.POST_STATUS_DRAFT); err != nil {
// 		c.app.GetLogger().Error("BlogAi.ContentUpdate.SaveDraft", slog.String("error", err.Error()))
// 		return api.Error("Failed to save draft: " + err.Error()).ToString()
// 	}

// 	return api.SuccessWithData("Draft saved successfully", map[string]any{
// 		"post": record,
// 	}).ToString()
// }

// func (c *Controller) onSaveFinal(data pageData) string {
// 	record, err := c.parseRecordPostFromRequest(data)
// 	if err != nil {
// 		return api.Error(err.Error()).ToString()
// 	}

// 	if record.Title == "" {
// 		return api.Error("Title is required").ToString()
// 	}

// 	if record.Summary == "" {
// 		return api.Error("Summary is required").ToString()
// 	}

// 	record.Status = blogai.POST_STATUS_PUBLISHED

// 	if err := c.persistRecordToPost(data.Post, record, blogstore.POST_STATUS_PUBLISHED); err != nil {
// 		c.app.GetLogger().Error("BlogAi.ContentUpdate.SaveFinal", slog.String("error", err.Error()))
// 		return api.Error("Failed to save post: " + err.Error()).ToString()
// 	}

// 	return api.SuccessWithData("Post saved successfully", map[string]any{
// 		"post": record,
// 	}).ToString()
// }

// func (c *Controller) parseRecordPostFromRequest(data pageData) (blogai.RecordPost, error) {
// 	postJSON := req.GetStringTrimmed(data.Request, "post")
// 	if strings.TrimSpace(postJSON) == "" {
// 		return blogai.RecordPost{}, fmt.Errorf("Post data is missing")
// 	}

// 	var record blogai.RecordPost
// 	if err := json.Unmarshal([]byte(postJSON), &record); err != nil {
// 		return blogai.RecordPost{}, fmt.Errorf("Failed to parse post data: %w", err)
// 	}

// 	return normalizeRecord(record, data.Post.ID()), nil
// }

// func (c *Controller) persistRecordToPost(post *blogstore.Post, record blogai.RecordPost, status string) error {
// 	markdown := RecordPostToMarkdown(record)

// 	post.SetTitle(record.Title)
// 	if record.Summary != "" {
// 		post.SetSummary(record.Summary)
// 	}
// 	post.SetContent(markdown)
// 	post.SetMetaDescription(record.MetaDescription)
// 	post.SetMetaKeywords(sliceToKeywords(record.MetaKeywords))
// 	post.SetStatus(status)
// 	post.SetEditor(blogstore.POST_EDITOR_MARKDOWN)
// 	post.SetImageUrl(record.Image)

// 	if err := c.app.GetBlogStore().PostUpdate(post); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func normalizeRecord(record blogai.RecordPost, postID string) blogai.RecordPost {
// 	record.ID = postID
// 	record.Title = strings.TrimSpace(record.Title)
// 	record.Subtitle = strings.TrimSpace(record.Subtitle)
// 	record.Status = strings.TrimSpace(record.Status)
// 	record.Summary = strings.TrimSpace(record.Summary)
// 	record.MetaDescription = strings.TrimSpace(record.MetaDescription)
// 	record.MetaTitle = strings.TrimSpace(record.MetaTitle)
// 	if record.MetaTitle == "" {
// 		record.MetaTitle = record.Title
// 	}
// 	record.Image = strings.TrimSpace(record.Image)

// 	record.MetaKeywords = cleanKeywordSlice(record.MetaKeywords)
// 	if len(record.MetaKeywords) == 0 {
// 		record.MetaKeywords = cleanKeywordSlice(record.Keywords)
// 	}
// 	record.Keywords = record.MetaKeywords

// 	record.Introduction.Title = strings.TrimSpace(record.Introduction.Title)
// 	record.Introduction.Paragraphs = cleanParagraphs(record.Introduction.Paragraphs)

// 	for idx := range record.Sections {
// 		record.Sections[idx].Title = strings.TrimSpace(record.Sections[idx].Title)
// 		record.Sections[idx].Paragraphs = cleanParagraphs(record.Sections[idx].Paragraphs)
// 	}

// 	record.Conclusion.Title = strings.TrimSpace(record.Conclusion.Title)
// 	record.Conclusion.Paragraphs = cleanParagraphs(record.Conclusion.Paragraphs)

// 	return record
// }

// func cleanKeywordSlice(items []string) []string {
// 	cleaned := lo.FilterMap(items, func(item string, _ int) (string, bool) {
// 		trimmed := strings.TrimSpace(item)
// 		return trimmed, trimmed != ""
// 	})
// 	if len(cleaned) == 0 {
// 		return []string{}
// 	}
// 	return lo.Uniq(cleaned)
// }

// func cleanParagraphs(paragraphs []string) []string {
// 	return lo.FilterMap(paragraphs, func(p string, _ int) (string, bool) {
// 		trimmed := strings.TrimSpace(p)
// 		return trimmed, trimmed != ""
// 	})
// }

// func sectionPayload(record blogai.RecordPost, section string) (any, error) {
// 	switch section {
// 	case "introduction":
// 		return record.Introduction, nil
// 	case "conclusion":
// 		return record.Conclusion, nil
// 	default:
// 		var index int
// 		if _, err := fmt.Sscanf(section, "section_%d", &index); err != nil {
// 			return nil, fmt.Errorf("Invalid section index")
// 		}
// 		if index < 0 || index >= len(record.Sections) {
// 			return nil, fmt.Errorf("Invalid section index")
// 		}
// 		return record.Sections[index], nil
// 	}
// }
