package aiposteditor

import (
	"log/slog"
	"project/pkg/blogai"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/base/req"
)

func (c *AiPostEditorController) onSaveDraft(data pageData) string {
	postJSON := req.Value(data.Request, "post")
	if postJSON == "" {
		return api.Error("Post data is missing").ToString()
	}
	if !strings.HasPrefix(postJSON, "{") {
		return api.Error("Invalid post data format").ToString()
	}
	record, err := RecordFromJSON(postJSON)
	if err != nil {
		return api.Error("Failed to parse post data: " + err.Error()).ToString()
	}

	data.BlogAiPost.Status = blogai.POST_STATUS_DRAFT
	data.BlogAiPost.Conclusion = record.Conclusion
	data.BlogAiPost.Image = record.Image
	data.BlogAiPost.Introduction = record.Introduction
	data.BlogAiPost.MetaDescription = record.MetaDescription
	data.BlogAiPost.MetaKeywords = record.MetaKeywords
	data.BlogAiPost.MetaTitle = record.MetaTitle
	data.BlogAiPost.Sections = record.Sections
	data.BlogAiPost.Summary = record.Summary
	data.BlogAiPost.Title = record.Title

	data.Record.SetPayload(data.BlogAiPost.ToJSON())
	if err := c.app.GetCustomStore().RecordUpdate(data.Record); err != nil {
		c.app.GetLogger().Error("failed to update blog post draft", slog.String("error", err.Error()))
		return api.Error("Failed to save draft: " + err.Error()).ToString()
	}

	return api.SuccessWithData("Draft saved successfully", map[string]any{}).ToString()
}
