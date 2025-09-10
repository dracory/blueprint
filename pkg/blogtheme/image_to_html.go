package blogtheme

import (
	"github.com/dracory/hb"
	"github.com/dracory/ui"
)

func (t *theme) imageToHtml(block ui.BlockInterface) *hb.Tag {
	imageUrl := block.Parameter("image_url")
	return hb.Img(imageUrl)
}
