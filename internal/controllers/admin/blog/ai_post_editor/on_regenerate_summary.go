package aiposteditor

import (
	"project/internal/controllers/admin/blog/shared"
	"project/pkg/blogai"

	"github.com/dracory/api"
)

func (c *AiPostEditorController) onRegenerateSummary(data pageData) string {
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

	summary, err := agent.GenerateSummary(llmEngine, data.BlogAiPost)
	if err != nil {
		return api.Error("Failed to regenerate summary: " + err.Error()).ToString()
	}

	data.BlogAiPost.Summary = summary
	data.Record.SetPayload(data.BlogAiPost.ToJSON())
	if err := c.app.GetCustomStore().RecordUpdate(data.Record); err != nil {
		return api.Error("Failed to save updated blog post: " + err.Error()).ToString()
	}

	return api.SuccessWithData("Summary regenerated successfully", map[string]any{"summary": summary}).ToString()
}
