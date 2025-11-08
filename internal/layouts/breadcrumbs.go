package layouts

import (
	"strings"

	"github.com/dracory/hb"
	"github.com/samber/lo"
)

func Breadcrumbs(breadcrumbs []Breadcrumb) *hb.Tag {
	nav := hb.Nav().Aria("label", "breadcrumb")
	ol := hb.OL().Class("breadcrumb").Style("border:none;")

	for _, breadcrumb := range breadcrumbs {
		// Icon
		icon := lo.IfF(breadcrumb.Icon != "", func() hb.TagInterface {
			if strings.HasPrefix(breadcrumb.Icon, "bi") {
				return hb.I().
					Class("bi " + breadcrumb.Icon).
					Style("margin-right: 8px;")
			}
			return hb.NewHTML(breadcrumb.Icon)
		}).ElseF(func() hb.TagInterface {
			return nil
		})

		// Link
		link := hb.Hyperlink().
			HTML(breadcrumb.Name).
			Href(breadcrumb.URL)

		// List Item
		li := hb.LI().
			Class("breadcrumb-item").
			Child(icon).
			Child(link)

		// Add List Item to Order List
		ol.Child(li)
	}

	nav.Child(ol)

	return nav
}
