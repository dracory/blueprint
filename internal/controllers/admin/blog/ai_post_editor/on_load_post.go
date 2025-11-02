package aiposteditor

import "github.com/dracory/api"

func (c *AiPostEditorController) onLoadPost(data pageData) string {
	if data.BlogAiPost.ID == "" {
		return api.Error("Post not found").ToString()
	}

	return api.SuccessWithData("Post loaded successfully", map[string]any{
		"id":              data.BlogAiPost.ID,
		"title":           data.BlogAiPost.Title,
		"summary":         data.BlogAiPost.Summary,
		"metaTitle":       data.BlogAiPost.MetaTitle,
		"metaDescription": data.BlogAiPost.MetaDescription,
		"metaKeywords":    data.BlogAiPost.MetaKeywords,
		"keywords":        data.BlogAiPost.Keywords,
		"image":           data.BlogAiPost.Image,
		"introduction":    data.BlogAiPost.Introduction,
		"sections":        data.BlogAiPost.Sections,
		"conclusion":      data.BlogAiPost.Conclusion,
		"status":          data.BlogAiPost.Status,
	}).ToString()
}
