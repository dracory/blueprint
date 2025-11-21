package aiposteditor

import (
	"context"
	"log/slog"
	"project/internal/controllers/admin/blog/shared"
	"project/pkg/blogai"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/base/req"
	"github.com/dracory/blogstore"
)

func (c *AiPostEditorController) onCreateFinalPost(data pageData) string {
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
	if record.Title == "" {
		return api.Error("Title is required").ToString()
	}
	if record.Summary == "" {
		return api.Error("Summary is required").ToString()
	}

	content := c.buildPostMarkdownContent(data.Request, record)
	post := blogstore.NewPost().
		SetID(data.BlogAiPost.ID).
		SetStatus(blogstore.POST_STATUS_PUBLISHED).
		SetTitle(record.Title).
		SetSummary(record.Summary).
		SetMetaKeywords(strings.Join(record.MetaKeywords, ", ")).
		SetMetaDescription(record.MetaDescription).
		SetContent(content).
		SetEditor(blogstore.POST_EDITOR_MARKDOWN)

	if record.Image != "" {
		post.SetImageUrl(record.Image)
	}

	data.BlogAiPost.Title = record.Title
	data.BlogAiPost.Summary = record.Summary
	data.BlogAiPost.Introduction = record.Introduction
	data.BlogAiPost.Sections = record.Sections
	data.BlogAiPost.Conclusion = record.Conclusion
	data.BlogAiPost.MetaTitle = record.MetaTitle
	data.BlogAiPost.MetaDescription = record.MetaDescription
	data.BlogAiPost.MetaKeywords = record.MetaKeywords
	data.BlogAiPost.Image = record.Image
	data.BlogAiPost.Status = blogai.POST_STATUS_PUBLISHED

	if err := c.app.GetBlogStore().PostCreate(context.Background(), post); err != nil {
		c.app.GetLogger().Error("failed to create blog post", slog.String("error", err.Error()))
		return api.Error("Failed to save blog post: " + err.Error()).ToString()
	}

	data.Record.SetPayload(data.BlogAiPost.ToJSON())
	if err := c.app.GetCustomStore().RecordUpdate(data.Record); err != nil {
		c.app.GetLogger().Error("failed to update blog post", slog.String("error", err.Error()))
		return api.Error("Failed to update blog post record: " + err.Error()).ToString()
	}

	return api.SuccessWithData("Blog post created successfully", map[string]any{
		"redirect": shared.NewLinks().PostManager(),
	}).ToString()
}
