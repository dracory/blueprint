package webtheme

import (
	"github.com/dracory/hb"
	"github.com/dracory/ui"
)

func (t *theme) rowToHtml(_ ui.BlockInterface) *hb.Tag {
	return hb.Div().Class("row")
}
