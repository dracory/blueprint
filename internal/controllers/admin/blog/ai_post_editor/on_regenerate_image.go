package aiposteditor

import (
	"log/slog"
	"project/internal/controllers/admin/blog/shared"
	"project/pkg/blogai"

	"github.com/dracory/api"
)

func (c *AiPostEditorController) onRegenerateImage(data pageData) string {
	agent := blogai.NewBlogWriterAgent(c.registry.GetLogger())
	if agent == nil {
		return api.Error("failed to initialize LLM engine").ToString()
	}

	llmEngine, err := shared.LlmEngine(c.registry)
	if err != nil {
		return api.Error("failed to initialize LLM engine").ToString()
	}

	imageDataURL, err := agent.GenerateImage(llmEngine, data.BlogAiPost.Title, data.BlogAiPost.Summary)
	if err != nil {
		c.registry.GetLogger().Error("BlogAi. Post Editor. Generate Image. Failed to generate image", slog.String("error", err.Error()))
		return api.Error("Failed to generate image").ToString()
	}

	data.BlogAiPost.Image = imageDataURL
	data.Record.SetPayload(data.BlogAiPost.ToJSON())
	if err := c.registry.GetCustomStore().RecordUpdate(data.Record); err != nil {
		return api.Error("Failed to save updated blog post with image: " + err.Error()).ToString()
	}

	return api.SuccessWithData("Image generated successfully", map[string]any{"image": imageDataURL}).ToString()
}
