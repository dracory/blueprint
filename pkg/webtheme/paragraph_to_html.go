package webtheme

import (
	"github.com/dracory/hb"
	"github.com/dracory/ui"
)

func (t *theme) paragraphToHtml(block ui.BlockInterface) *hb.Tag {
	text := block.Parameter("content")
	return hb.P().HTML(text)
}
