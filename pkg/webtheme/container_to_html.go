package webtheme

import (
	"github.com/dracory/hb"
	"github.com/dracory/ui"
)

func (t *theme) containerToHtml(block ui.BlockInterface) *hb.Tag {
	return hb.Div().Class("container")
}
