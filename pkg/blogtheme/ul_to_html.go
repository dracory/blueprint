package blogtheme

import (
	"github.com/dracory/hb"
	"github.com/dracory/ui"
)

func (t *theme) ulToHtml(block ui.BlockInterface) *hb.Tag {
	text := block.Parameter("content")
	return hb.UL().HTML(text)
}

func (t *theme) olToHtml(block ui.BlockInterface) *hb.Tag {
	text := block.Parameter("content")
	return hb.OL().HTML(text)
}

func (t *theme) liToHtml(block ui.BlockInterface) *hb.Tag {
	text := block.Parameter("content")
	return hb.LI().HTML(text)
}
