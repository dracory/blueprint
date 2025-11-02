package aiposteditor

import (
	"log/slog"
	"project/internal/controllers/admin/blog/shared"
	"project/pkg/blogai"

	"github.com/dracory/api"
	"github.com/dracory/base/req"
	"github.com/spf13/cast"
)

func (c *AiPostEditorController) onRegenerateParagraph(data pageData) string {
	sectionType := req.Value(data.Request, "section_type")
	sectionIndexStr := req.Value(data.Request, "section_index")
	paragraphIndexStr := req.Value(data.Request, "paragraph_index")

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

	sectionIndex := cast.ToInt(sectionIndexStr)
	paragraphIndex := cast.ToInt(paragraphIndexStr)

	newParagraph, err := agent.RegenerateParagraph(llmEngine, data.BlogAiPost, sectionType, sectionIndex, paragraphIndex)
	if err != nil {
		c.app.GetLogger().Error("BlogAi. Post Editor. Regenerate Paragraph. Failed", slog.String("error", err.Error()))
		return api.Error("Failed to regenerate paragraph: " + err.Error()).ToString()
	}

	switch sectionType {
	case "introduction":
		if paragraphIndex == len(data.BlogAiPost.Introduction.Paragraphs) {
			data.BlogAiPost.Introduction.Paragraphs = append(data.BlogAiPost.Introduction.Paragraphs, newParagraph)
		} else {
			data.BlogAiPost.Introduction.Paragraphs[paragraphIndex] = newParagraph
		}
	case "conclusion":
		if paragraphIndex == len(data.BlogAiPost.Conclusion.Paragraphs) {
			data.BlogAiPost.Conclusion.Paragraphs = append(data.BlogAiPost.Conclusion.Paragraphs, newParagraph)
		} else {
			data.BlogAiPost.Conclusion.Paragraphs[paragraphIndex] = newParagraph
		}
	case "section":
		if paragraphIndex == len(data.BlogAiPost.Sections[sectionIndex].Paragraphs) {
			data.BlogAiPost.Sections[sectionIndex].Paragraphs = append(data.BlogAiPost.Sections[sectionIndex].Paragraphs, newParagraph)
		} else {
			data.BlogAiPost.Sections[sectionIndex].Paragraphs[paragraphIndex] = newParagraph
		}
	}

	data.Record.SetPayload(data.BlogAiPost.ToJSON())
	if err := c.app.GetCustomStore().RecordUpdate(data.Record); err != nil {
		return api.Error("Failed to save updated blog post: " + err.Error()).ToString()
	}

	return api.SuccessWithData("Paragraph regenerated successfully", map[string]any{"paragraph": newParagraph}).ToString()
}
