package blogai

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/dracory/llm"
)

// This file defines the BlogWriterAgent, which is responsible for
// generating blog post content based on a given title.

// BlogWriterAgent is responsible for generating blog post content.
type BlogWriterAgent struct {
	logger *slog.Logger
}

// NewBlogWriterAgent creates a new instance of the BlogWriterAgent.
func NewBlogWriterAgent(logger *slog.Logger) *BlogWriterAgent {
	return &BlogWriterAgent{
		logger: logger,
	}
}

// WriteBlogPost takes a title and generates a detailed blog post.
func (b *BlogWriterAgent) WriteBlogPost(ai llm.LlmInterface, title string) (RecordPost, error) {
	userPrompt := `You are an experienced blog post writer. Your task is create a blog post with the title: %s`

	systemPrompt := `Write a professional, engaging, and well-researched blog post.

Requirements:
- Write in a professional yet engaging tone
- Include relevant examples and data points
- Use proper formatting with headers (use markdown)
- Ensure smooth transitions between sections
- Add a compelling introduction and conclusion
- Aim for clarity and actionable insights
- Include citations or references where appropriate
- Where it would genuinely benefit the reader, and only if it fits naturally, subtly include a backlink to https://roastmycontract.com (for example, as a further resource, tool, or actionable suggestion). Do NOT force a backlink in every paragraph or if it would feel unnatural or promotional.

IMPORTANT: Your response must be a valid JSON object with this exact structure:
{
	"title": "Post Title",
	"subtitle": "Post Subtitle",
	"summary": "Brief summary of the post",
	"introduction": {
		"title": "Introduction Paragraph Title",
		"paragraphs": ["Introduction paragraph content"]
	},
	"sections": [
		{
			"title": "Section Paragraph Title",
			"paragraphs": ["Section paragraph content"]
		}
	],
	"conclusion": {
		"title": "Conclusion Paragraph Title",
		"paragraphs": ["Conclusion paragraph content"]
	}
}

Do not include any text outside of this JSON structure. The response must be parseable as JSON.
`

	// Format prompts
	userPrompt = fmt.Sprintf(userPrompt, title)

	// Generate the blog post
	content, err := ai.Generate(systemPrompt, userPrompt, llm.LlmOptions{
		MaxTokens:    8192, // Increased token limit for longer responses
		Temperature:  0.7,
		OutputFormat: llm.OutputFormatJSON,
	})
	if err != nil {
		return RecordPost{}, fmt.Errorf("failed to generate blog post: %w", err)
	}

	fmt.Printf("Raw LLM response:\n%s\n", content)

	cleanContent := sanitizeJSONContent(content)

	// Parse the response into a RecordPost
	var post RecordPost
	err = json.Unmarshal([]byte(cleanContent), &post)
	if err != nil {
		fmt.Printf("JSON parsing error: %v\n", err)
		fmt.Printf("Attempted to parse (sanitized):\n%s\n", cleanContent)
		return RecordPost{}, fmt.Errorf("failed to parse blog post JSON: %w", err)
	}

	// Validate required fields
	if post.Title == "" {
		return RecordPost{}, fmt.Errorf("missing title in response")
	}
	if post.Subtitle == "" {
		return RecordPost{}, fmt.Errorf("missing subtitle in response")
	}
	if post.Summary == "" {
		return RecordPost{}, fmt.Errorf("missing summary in response")
	}
	if post.Introduction.Title == "" || len(post.Introduction.Paragraphs) == 0 {
		return RecordPost{}, fmt.Errorf("missing introduction fields in response")
	}
	if len(post.Sections) == 0 {
		return RecordPost{}, fmt.Errorf("missing sections in response")
	}
	if post.Conclusion.Title == "" || len(post.Conclusion.Paragraphs) == 0 {
		return RecordPost{}, fmt.Errorf("missing conclusion fields in response")
	}

	fmt.Printf("Successfully parsed post: %+v\n", post)

	return post, nil
}

// buildSectionRegenerationPrompt constructs the system and user prompts for regenerating a section
func buildSectionRegenerationPrompt(post RecordPost, section string) (systemPrompt, userPrompt string) {
	if section == "introduction" {
		systemPrompt = `You are an experienced blog post writer. Your task is to regenerate the introduction section of a blog post.

Requirements:
- Write in a professional yet engaging tone
- Use as many paragraphs as are appropriate for the introduction (can be one or several)
- Use proper formatting (use markdown, allowed bold, italic, links)
- Hook the reader's attention immediately
- Provide context for the topic
- Preview what the post will cover
- Set the tone for the rest of the article
- Where it would genuinely benefit the reader, and only if it fits naturally, subtly include a backlink to https://roastmycontract.com (for example, as a further resource, tool, or actionable suggestion). Do NOT force a backlink if it would feel unnatural or promotional.
- Your response must be a valid JSON object with this exact structure:
{
	"title": "The section title",
	"paragraphs": ["First paragraph.", "Second paragraph.", "..."]
}`
		userPrompt = fmt.Sprintf(`Write an engaging introduction for a blog post with the title: %s`, post.Title)
		return
	}

	if section == "conclusion" {
		systemPrompt = `You are an experienced blog post writer. Your task is to regenerate the conclusion section of a blog post.

Requirements:
- Write in a professional yet engaging tone
- Use as many paragraphs as are appropriate for the conclusion (can be one or several)
- Summarize the key points
- Reinforce the main message
- Provide a clear call to action or next steps
- Leave the reader with something to think about
- Use proper formatting (use markdown, allowed bold, italic, links)
- Where it would genuinely benefit the reader, and only if it fits naturally, subtly include a backlink to https://roastmycontract.com (for example, as a further resource, tool, or actionable suggestion). Do NOT force a backlink if it would feel unnatural or promotional.
- Your response must be a valid JSON object with this exact structure:
{
	"title": "The section title",
	"paragraphs": ["First paragraph.", "Second paragraph.", "..."]
}`
		userPrompt = fmt.Sprintf(`Write a strong conclusion for a blog post with the title: %s`, post.Title)
		return
	}

	// Handle numbered sections (section_0, section_1, etc.)
	var sectionIndex int
	_, err := fmt.Sscanf(section, "section_%d", &sectionIndex)
	if err != nil || sectionIndex < 0 || sectionIndex >= len(post.Sections) {
		return "", ""
	}
	systemPrompt = `You are an experienced blog post writer. Your task is to regenerate a specific section of a blog post.

Requirements:
- Write in a professional yet engaging tone
- Use as many paragraphs as are appropriate for the section (can be one or several)
- Be informative, engaging, and relevant
- Use proper formatting (use markdown, allowed bold, italic, links)
- Include examples or data if possible
- Maintain consistent tone with the rest of the post
- Ensure smooth transitions with surrounding content
- Where it would genuinely benefit the reader, and only if it fits naturally, subtly include a backlink to https://roastmycontract.com (for example, as a further resource, tool, or actionable suggestion). Do NOT force a backlink if it would feel unnatural or promotional.
- Your response must be a valid JSON object with this exact structure:
{
	"title": "The section title",
	"paragraphs": ["First paragraph.", "Second paragraph.", "..."]
}`
	userPrompt = fmt.Sprintf(`Write a section for a blog post titled: %s\nSection title: %s`, post.Title, post.Sections[sectionIndex].Title)
	return
}

// RegenerateSection regenerates a specific section of the blog post
//
// Business logic:
// 1. Validate input
// 2. Prepare prompt
// 3. Generate new section content
// 4. Update post
// 5. Return updated post
func (b *BlogWriterAgent) RegenerateSection(ai llm.LlmInterface, post RecordPost, section string) (RecordPost, error) {
	if ai == nil {
		return RecordPost{}, fmt.Errorf("LLM interface not initialized")
	}

	systemPrompt, userPrompt := buildSectionRegenerationPrompt(post, section)

	b.logger.Info("BlogAi. BlogWriterAgent. Regenerate Section. Prompt:",
		slog.String("systemPrompt", systemPrompt),
		slog.String("userPrompt", userPrompt))

	// Generate the content
	content, err := ai.Generate(systemPrompt, userPrompt, llm.LlmOptions{
		MaxTokens:    4096,
		Temperature:  0.7,
		OutputFormat: llm.OutputFormatJSON,
	})
	if err != nil {
		b.logger.Error("BlogAi. BlogWriterAgent. Regenerate Section. Failed to generate section content",
			slog.String("error", err.Error()))
		return post, fmt.Errorf("failed to generate section content: %w", err)
	}

	b.logger.Info("BlogAi. BlogWriterAgent. Regenerate Section. Raw Response:",
		slog.String("content", content))

	cleanContent := sanitizeJSONContent(content)

	// Parse the response
	var response struct {
		Title      string   `json:"title"`
		Paragraphs []string `json:"paragraphs"`
	}
	err = json.Unmarshal([]byte(cleanContent), &response)
	if err != nil {
		b.logger.Error("BlogAi. BlogWriterAgent. Regenerate Section. Failed to parse section content JSON",
			slog.String("error", err.Error()),
			slog.String("content", cleanContent))
		return post, fmt.Errorf("failed to parse section content JSON: %w raw content: %s", err, content)
	}

	if response.Title == "" || len(response.Paragraphs) == 0 {
		b.logger.Error("BlogAi. BlogWriterAgent. Regenerate Section. Generated content or title is empty")
		return post, fmt.Errorf("generated content or title is empty")
	}

	// Update the appropriate section
	switch section {
	case "introduction":
		post.Introduction.Title = response.Title
		post.Introduction.Paragraphs = response.Paragraphs
	case "conclusion":
		post.Conclusion.Title = response.Title
		post.Conclusion.Paragraphs = response.Paragraphs
	default:
		var sectionIndex int
		fmt.Sscanf(section, "section_%d", &sectionIndex)
		post.Sections[sectionIndex].Title = response.Title
		post.Sections[sectionIndex].Paragraphs = response.Paragraphs
	}

	return post, nil
}

// RegenerateParagraph regenerates a single paragraph in the post based on full context.
func (b *BlogWriterAgent) RegenerateParagraph(ai llm.LlmInterface, post RecordPost, sectionType string, sectionIndex, paragraphIndex int) (string, error) {
	if ai == nil {
		return "", fmt.Errorf("LLM interface not initialized")
	}

	const PLACEHOLDER_TEXT = "<<<REGENERATE_THIS_PARAGRAPH>>>"

	// Clone post and insert placeholder
	postWithPlaceholder := post
	switch sectionType {
	case "introduction":
		if paragraphIndex < 0 || paragraphIndex > len(postWithPlaceholder.Introduction.Paragraphs) {
			return "", fmt.Errorf("invalid introduction paragraph index")
		}
		if paragraphIndex == len(postWithPlaceholder.Introduction.Paragraphs) {
			postWithPlaceholder.Introduction.Paragraphs = append(postWithPlaceholder.Introduction.Paragraphs, PLACEHOLDER_TEXT)
		} else {
			postWithPlaceholder.Introduction.Paragraphs[paragraphIndex] = PLACEHOLDER_TEXT
		}
	case "conclusion":
		if paragraphIndex < 0 || paragraphIndex > len(postWithPlaceholder.Conclusion.Paragraphs) {
			return "", fmt.Errorf("invalid conclusion paragraph index")
		}
		if paragraphIndex == len(postWithPlaceholder.Conclusion.Paragraphs) {
			postWithPlaceholder.Conclusion.Paragraphs = append(postWithPlaceholder.Conclusion.Paragraphs, PLACEHOLDER_TEXT)
		} else {
			postWithPlaceholder.Conclusion.Paragraphs[paragraphIndex] = PLACEHOLDER_TEXT
		}
	case "section":
		if sectionIndex < 0 || sectionIndex >= len(postWithPlaceholder.Sections) {
			return "", fmt.Errorf("invalid section index")
		}
		if paragraphIndex < 0 || paragraphIndex > len(postWithPlaceholder.Sections[sectionIndex].Paragraphs) {
			return "", fmt.Errorf("invalid section paragraph index")
		}
		if paragraphIndex == len(postWithPlaceholder.Sections[sectionIndex].Paragraphs) {
			postWithPlaceholder.Sections[sectionIndex].Paragraphs = append(postWithPlaceholder.Sections[sectionIndex].Paragraphs, PLACEHOLDER_TEXT)
		} else {
			postWithPlaceholder.Sections[sectionIndex].Paragraphs[paragraphIndex] = PLACEHOLDER_TEXT
		}
	default:
		return "", fmt.Errorf("unknown section type: %s", sectionType)
	}

	postJSON, err := json.MarshalIndent(postWithPlaceholder, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal post JSON: %w", err)
	}

	systemPrompt := `
You are an expert blog writer. You are given the full JSON of a blog post.
One paragraph contains the text <<<REGENERATE_THIS_PARAGRAPH>>>.
Your job is to generate a single paragraph that fits perfectly in that spot,
matching the style, tone, and context of the rest of the post. 
Return ONLY the paragraph as a plain string, no JSON, no markdown, no explanations.
`
	userPrompt := "POST JSON:\n" + string(postJSON)

	b.logger.Info("BlogAi. BlogWriterAgent. Regenerate Paragraph. Prompt:",
		slog.String("systemPrompt", systemPrompt),
		slog.String("userPrompt", userPrompt))

	paragraph, err := ai.Generate(systemPrompt, userPrompt, llm.LlmOptions{
		MaxTokens:    512,
		Temperature:  0.7,
		OutputFormat: llm.OutputFormatText,
	})
	if err != nil {
		b.logger.Error("BlogAi. BlogWriterAgent. Regenerate Paragraph. Failed to generate paragraph",
			slog.String("error", err.Error()))
		return "", fmt.Errorf("failed to generate paragraph: %w", err)
	}

	return paragraph, nil
}

// GenerateImage generates an image for the blog post using the title and summary
func (b *BlogWriterAgent) GenerateImage(ai llm.LlmInterface, title string, summary string) (string, error) {
	systemPrompt := `
You are an AI image generator. Generate an image that represents the given blog post title and summary.
The image should be visually appealing and relevant to the content.
`

	userPrompt := fmt.Sprintf("Title: %s\nSummary: %s", title, summary)

	prompt := systemPrompt + "\n\n" + userPrompt

	// Generate the image URL
	imageURL, err := ai.GenerateImage(prompt, llm.LlmOptions{
		Verbose:      true,
		MaxTokens:    100,
		Temperature:  0.7,
		OutputFormat: llm.OutputFormatImagePNG,
		Model:        llm.OPENROUTER_MODEL_GEMINI_2_5_FLASH_IMAGE,
	})
	if err != nil {
		b.logger.Error("BlogAi. BlogWriterAgent. GenerateImage. Failed to generate image",
			slog.String("error", err.Error()))
		return "", fmt.Errorf("failed to generate image: %w", err)
	}

	// convert to data URL
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageURL), nil
}

// GenerateSummary generates a summary for the given blog post.
func (b *BlogWriterAgent) GenerateSummary(ai llm.LlmInterface, post RecordPost) (string, error) {
	if ai == nil {
		return "", fmt.Errorf("LLM interface not initialized")
	}

	systemPrompt := `
You are an experienced blog post writer. Your task is to generate a concise,
informative summary for the following blog post. The summary should capture
the main points and be suitable for use as an excerpt or meta description.
`

	// We'll serialize the post to JSON and include it as context
	postJSON, err := json.Marshal(post)
	if err != nil {
		return "", fmt.Errorf("failed to serialize post: %w", err)
	}

	userPrompt := "POST JSON:\n" + string(postJSON)

	b.logger.Info("BlogAi. BlogWriterAgent. GenerateSummary. Prompt:",
		slog.String("systemPrompt", systemPrompt),
		slog.String("userPrompt", userPrompt))

	summary, err := ai.Generate(systemPrompt, userPrompt, llm.LlmOptions{
		MaxTokens:    200,
		Temperature:  0.6,
		OutputFormat: llm.OutputFormatText,
	})
	if err != nil {
		b.logger.Error("BlogAi. BlogWriterAgent. GenerateSummary. Failed to generate summary",
			slog.String("error", err.Error()))
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	return summary, nil
}

// GenerateMetas generates meta title, meta description, and keywords for the given blog post.
func (b *BlogWriterAgent) GenerateMetas(ai llm.LlmInterface, post RecordPost) (metaTitle, metaDesc, metaKeywords string, err error) {
	if ai == nil {
		return "", "", "", fmt.Errorf("LLM interface not initialized")
	}

	systemPrompt := `
You are an experienced SEO specialist. Generate the following for the given blog post:
1. Meta title (max 60 chars)
2. Meta description (max 160 chars)
3. Meta keywords (comma separated)
`

	postJSON, err := json.Marshal(post)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to serialize post: %w", err)
	}

	userPrompt := "POST JSON:\n" + string(postJSON)

	b.logger.Info("BlogAi. BlogWriterAgent. GenerateMetas. Prompt:",
		slog.String("systemPrompt", systemPrompt),
		slog.String("userPrompt", userPrompt))

	response, err := ai.Generate(systemPrompt, userPrompt, llm.LlmOptions{
		MaxTokens:    300,
		Temperature:  0.6,
		OutputFormat: llm.OutputFormatJSON,
	})
	if err != nil {
		b.logger.Error("BlogAi. BlogWriterAgent. GenerateMetas. Failed to generate metas",
			slog.String("error", err.Error()))
		return "", "", "", fmt.Errorf("failed to generate metas: %w", err)
	}

	cleanResponse := sanitizeJSONContent(response)

	var resp struct {
		MetaTitle       string `json:"meta_title"`
		MetaDescription string `json:"meta_description"`
		MetaKeywords    string `json:"meta_keywords"`
	}
	err = json.Unmarshal([]byte(cleanResponse), &resp)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse metas response: %w", err)
	}

	return resp.MetaTitle, resp.MetaDescription, resp.MetaKeywords, nil
}

func sanitizeJSONContent(content string) string {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "```") {
		return trimmed
	}

	trimmed = strings.TrimPrefix(trimmed, "```json")
	trimmed = strings.TrimPrefix(trimmed, "```")
	trimmed = strings.TrimSpace(trimmed)
	trimmed = strings.TrimSuffix(trimmed, "```")
	trimmed = strings.TrimSpace(trimmed)

	if trimmed == "" {
		return trimmed
	}

	const jsonStartChars = "{"

	if idx := strings.IndexAny(trimmed, jsonStartChars); idx > 0 {
		trimmed = trimmed[idx:]
	}

	return strings.TrimSpace(trimmed)
}
