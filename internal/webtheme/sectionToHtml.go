package webtheme

import (
	"github.com/gouniverse/hb"
	"github.com/gouniverse/ui"
)

func (t *theme) sectionToHtml(block ui.BlockInterface) *hb.Tag {
	return hb.Section()
}
