package webtheme

import (
	"github.com/dracory/hb"
	"github.com/dracory/ui"
)

func (t *theme) sectionToHtml(block ui.BlockInterface) *hb.Tag {
	return hb.Section()
}
