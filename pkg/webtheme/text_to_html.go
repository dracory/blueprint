package webtheme

import (
	"github.com/dracory/hb"
	"github.com/dracory/ui"
)

func (t *theme) textToHtml(block ui.BlockInterface) *hb.Tag {
	text := block.Parameter("content")
	return hb.Span().HTML(text)
}
