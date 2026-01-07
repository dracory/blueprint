package aiposteditor

import (
	"fmt"
	"log/slog"
	"project/internal/controllers/admin/blog/shared"
	"project/pkg/blogai"

	"github.com/dracory/api"
	"github.com/dracory/base/req"
)

func (c *AiPostEditorController) onRegenerateSection(data pageData) string {
	section := req.Value(data.Request, "section")
	if section == "" {
		return api.Error("Section parameter is missing").ToString()
	}

	agent := blogai.NewBlogWriterAgent(c.registry.GetLogger())
	if agent == nil {
		return api.Error("failed to initialize LLM engine").ToString()
	}

	llmEngine, err := shared.LlmEngine(c.registry)
	if err != nil {
		return api.Error("failed to initialize LLM engine: " + err.Error()).ToString()
	}
	if llmEngine == nil {
		return api.Error("failed to initialize LLM engine").ToString()
	}

	data.BlogAiPost, err = agent.RegenerateSection(llmEngine, data.BlogAiPost, section)
	if err != nil {
		c.registry.GetLogger().Error("BlogAi. Post Editor. Regenerate Section. Failed to regenerate section", slog.String("error", err.Error()))
		return api.Error("Failed to regenerate section: " + section).ToString()
	}

	data.Record.SetPayload(data.BlogAiPost.ToJSON())
	if err := c.registry.GetCustomStore().RecordUpdate(data.Record); err != nil {
		return api.Error("Failed to save updated blog post: " + err.Error()).ToString()
	}

	var sectionData any
	switch section {
	case "introduction":
		sectionData = data.BlogAiPost.Introduction
	case "conclusion":
		sectionData = data.BlogAiPost.Conclusion
	default:
		var sectionIndex int
		if _, err := fmt.Sscanf(section, "section_%d", &sectionIndex); err != nil {
			return api.Error("Invalid section index format").ToString()
		}
		if sectionIndex >= 0 && sectionIndex < len(data.BlogAiPost.Sections) {
			sectionData = data.BlogAiPost.Sections[sectionIndex]
		} else {
			return api.Error("Invalid section index").ToString()
		}
	}

	return api.SuccessWithData("Section regenerated successfully", map[string]any{"section": sectionData}).ToString()
}
