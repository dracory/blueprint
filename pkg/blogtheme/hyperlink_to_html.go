package blogtheme

import (
	"github.com/dracory/hb"
	"github.com/dracory/ui"
)

func (t *theme) hyperlinkToHtml(block ui.BlockInterface) *hb.Tag {
	url := block.Parameter("url")
	text := block.Parameter("content")
	return hb.A().Href(url).HTML(text)
}
