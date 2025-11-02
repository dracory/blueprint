package aiposteditor

import (
	"project/internal/controllers/admin/blog/shared"
	"project/pkg/blogai"
	"strings"

	"github.com/dracory/api"
)

func (c *AiPostEditorController) onRegenerateMetas(data pageData) string {
	agent := blogai.NewBlogWriterAgent(c.app.GetLogger())
	if agent == nil {
		return api.Error("failed to initialize LLM engine").ToString()
	}

	llmEngine, err := shared.LlmEngine(c.app)
	if err != nil {
		return api.Error("failed to initialize LLM engine: " + err.Error()).ToString()
	}
	if llmEngine == nil {
		return api.Error("failed to initialize LLM engine").ToString()
	}

	metaTitle, metaDescription, metaKeywords, err := agent.GenerateMetas(llmEngine, data.BlogAiPost)
	if err != nil {
		return api.Error("Failed to regenerate meta information: " + err.Error()).ToString()
	}

	data.BlogAiPost.MetaTitle = metaTitle
	data.BlogAiPost.MetaDescription = metaDescription
	data.BlogAiPost.Keywords = strings.Split(metaKeywords, ",")
	data.Record.SetPayload(data.BlogAiPost.ToJSON())
	if err := c.app.GetCustomStore().RecordUpdate(data.Record); err != nil {
		return api.Error("Failed to save updated blog post: " + err.Error()).ToString()
	}

	return api.SuccessWithData("Meta information regenerated successfully", map[string]any{
		"metaTitle":       metaTitle,
		"metaDescription": metaDescription,
		"metaKeywords":    metaKeywords,
	}).ToString()
}
