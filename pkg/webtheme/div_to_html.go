package webtheme

import (
	"github.com/dracory/hb"
	"github.com/dracory/ui"
)

func (t *theme) divToHtml(_ ui.BlockInterface) *hb.Tag {
	return hb.Div()
}
