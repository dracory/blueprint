package aipostcontentupdate

import (
	"encoding/json"
	"net/http"
	"strings"

	"project/pkg/blogadmin/shared"

	"github.com/aws/smithy-go/ptr"
	"github.com/dracory/api"
	"github.com/dracory/llm"
	"github.com/samber/lo"
)

func (c *Controller) handleRegenerateBlock(r *http.Request) string {
	var reqBody struct {
		BlockID string              `json:"block_id"`
		Blocks  []map[string]string `json:"blocks"`
		Title   string              `json:"title"`
		Summary string              `json:"summary"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return api.Error("Invalid request body").ToString()
	}

	if strings.TrimSpace(reqBody.BlockID) == "" {
		return api.Error("Block ID is required").ToString()
	}

	// Reconstruct blocks from request
	blocks := make([]Block, 0, len(reqBody.Blocks))
	for _, b := range reqBody.Blocks {
		blocks = append(blocks, Block{
			ID:   b["id"],
			Type: BlockType(b["type"]),
			Text: b["text"],
		})
	}

	_, idx, found := lo.FindIndexOf(blocks, func(b Block) bool {
		return b.ID == reqBody.BlockID
	})
	if !found || idx < 0 || idx >= len(blocks) {
		return api.Error("Block not found").ToString()
	}

	block := blocks[idx]
	if strings.TrimSpace(block.Text) == "" {
		return api.Error("Block text is empty").ToString()
	}

	// Build full-post context where this block is replaced by a marker
	contextBlocks := make([]Block, len(blocks))
	copy(contextBlocks, blocks)
	contextBlocks[idx].Text = "== BLOCK TO REPLACE HERE =="
	contextMarkdown := BlocksToMarkdown(contextBlocks)

	systemPrompt := `You are an expert blog editor.
You will receive the FULL blog post as markdown, where exactly ONE block is replaced
with the marker "== BLOCK TO REPLACE HERE ==".

Your task is to REWRITE that missing block to improve clarity, style, and readability,
while preserving the original meaning and fitting naturally into the surrounding content.

IMPORTANT:
- You MUST significantly rephrase the original block. Do NOT return the same text.
- Do NOT copy long spans verbatim from the original block content.
- If the original block is a heading, return a short, strong heading.
- If it is a paragraph, return one or more paragraphs of body text.

Return ONLY the rewritten text for that block as markdown, with no additional explanations.`

	userPrompt := []string{}
	userPrompt = append(userPrompt, "Post title: "+reqBody.Title)
	userPrompt = append(userPrompt, "Block type: "+string(block.Type))
	userPrompt = append(userPrompt, "Full post markdown with marker:\n"+contextMarkdown)
	userPrompt = append(userPrompt, "Original block content to be regenerated (for meaning only; MUST BE REPHRASED SIGNIFICANTLY; DO NOT copy it verbatim):\n"+block.Text)

	engine, err := shared.LlmEngine(c.app)
	if err != nil || engine == nil {
		return api.Error("Failed to initialize LLM engine. Please try again later.").ToString()
	}
	resp, err := engine.Generate(systemPrompt, strings.Join(userPrompt, "\n\n"), llm.LlmOptions{
		MaxTokens:    512,
		Temperature:  ptr.Float64(0.7),
		OutputFormat: llm.OutputFormatText,
	})
	if err != nil {
		return api.Error("Failed to regenerate block content. Please try again later.").ToString()
	}

	newText := strings.TrimSpace(resp)
	if newText == "" || newText == block.Text {
		return api.Error("LLM did not provide a meaningful rewrite for this block.").ToString()
	}

	return api.SuccessWithData("Block regenerated", map[string]any{
		"text": newText,
	}).ToString()
}
