package webtheme

import (
	"github.com/dracory/hb"
	"github.com/dracory/ui"
)

func (t *theme) iconToHtml(block ui.BlockInterface) *hb.Tag {
	icon := block.Parameter("icon")
	return hb.I().Class(icon)
}
