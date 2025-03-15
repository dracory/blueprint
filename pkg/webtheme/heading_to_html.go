package webtheme

import (
	"github.com/gouniverse/hb"
	"github.com/gouniverse/ui"
	"github.com/spf13/cast"
)

func (t *theme) headingToHtml(block ui.BlockInterface) *hb.Tag {
	level := block.Parameter("level")
	if level == "" {
		level = "1"
	}

	text := block.Parameter("content")

	levelInt := cast.ToInt(level)
	levelStr := cast.ToString(levelInt)

	return hb.NewTag(`h` + levelStr).
		Style("margin-bottom:20px;margin-top:20px;").
		HTML(text)
}
