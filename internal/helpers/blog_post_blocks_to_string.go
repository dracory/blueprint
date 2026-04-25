package helpers

import (
	"encoding/json"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Block struct {
	ID         string         `json:"id"`
	Type       string         `json:"type"`
	Sequence   int            `json:"sequence"`
	ParentID   string         `json:"ParentId"`
	Text       string         `json:"content"`
	Attributes map[string]any `json:"attributes"`
}

func BlogPostBlocksToString(blocksString string) string {
	blocksAny := []map[string]any{}
	err := json.Unmarshal([]byte(blocksString), &blocksAny)

	if err != nil {
		return "Error parsing content. Please try again later."
	}

	blocksMap := blocksAny

	html := ""
	for _, blockMap := range blocksMap {
		blockType := blockMap["Type"].(string)
		blockID := blockMap["Id"].(string)
		parentID := blockMap["ParentId"].(string)
		attributes := blockMap["Attributes"].(map[string]any)
		sequence := blockMap["Sequence"].(float64)
		sequenceInt := cast.ToInt(sequence)

		block := Block{
			ID:         blockID,
			Type:       blockType,
			Sequence:   int(sequenceInt),
			ParentID:   parentID,
			Attributes: attributes,
		}

		html += processBlock(block)

	}

	return html
}

func processBlock(block Block) string {
	switch block.Type {
	case "code", "Code":
		return blockEditorBlockCodeToHtml(block)
	case "heading", "Heading":
		return blockEditorBlockHeadingToHtml(block)
	case "image", "Image":
		return blockEditorBlockImageToHtml(block)
	case "text", "Text":
		return blockEditorBlockTextToHtml(block)
	case "raw-html", "RawHtml":
		return blockEditorBlockRawHtmlToHtml(block)
	default:
		return "Block " + block.Type + " renderer does not exist"
	}
}

func blockEditorBlockCodeToHtml(block Block) string {
	code := lo.ValueOr(block.Attributes, "Code", "").(string)
	language := lo.ValueOr(block.Attributes, "Language", "").(string)

	html := ``
	html += `<div class="card" style="margin-bottom:20px;">`
	html += `  <div class="card-header">Language: ` + language + `</div>`
	html += `  <div class="card-body"><pre><code>` + code + `</code></pre></div>`
	html += `</div>`
	return html
}

func blockEditorBlockHeadingToHtml(block Block) string {
	level := lo.ValueOr(block.Attributes, "Level", "1").(string)
	text := lo.ValueOr(block.Attributes, "Text", "").(string)
	levelInt := cast.ToInt32(level)
	levelStr := cast.ToString(levelInt)

	return `<h` + levelStr + ` style="margin-bottom:20px;margin-top:20px;">` + text + `</h` + levelStr + `>`
}

func blockEditorBlockImageToHtml(block Block) string {
	url := lo.ValueOr(block.Attributes, "Url", "").(string)
	return `<img src="` + url + `" class="img img-responsive img-thumbnail" />`
}

func blockEditorBlockTextToHtml(block Block) string {
	text := lo.ValueOr(block.Attributes, "Text", "").(string)
	return `<p>` + text + `</p>`
}

func blockEditorBlockRawHtmlToHtml(block Block) string {
	text := lo.ValueOr(block.Attributes, "Text", "").(string)
	return text
}
