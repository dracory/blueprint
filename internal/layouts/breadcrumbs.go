package layouts

import (
	"strings"

	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func Breadcrumbs(breadcrumbs []Breadcrumb) *hb.Tag {
	nav := hb.Nav().Aria("label", "breadcrumb")
	ol := hb.OL().Class("breadcrumb")

	for _, breadcrumb := range breadcrumbs {
		icon := lo.IfF(breadcrumb.Icon != "", func() hb.TagInterface {
			if strings.HasSuffix(breadcrumb.Icon, "bi") {
				return hb.I().Class("bi " + breadcrumb.Icon)
			}
			return hb.NewHTML(breadcrumb.Icon)
		}).ElseF(func() hb.TagInterface {
			return nil
		})

		link := hb.Hyperlink().
			HTML(breadcrumb.Name).
			Href(breadcrumb.URL)

		li := hb.LI().
			Class("breadcrumb-item").
			Child(icon).
			Child(link)

		ol.Child(li)
	}

	nav.Child(ol)

	return nav
}
