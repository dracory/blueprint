package partials

import (
	"project/internal/layouts"
	"project/internal/links"

	"github.com/dracory/hb"
	"github.com/dracory/uid"

	"github.com/samber/lo"
)

func PageHeader(iconName string, title string, breadcrumbs ...[]layouts.Breadcrumb) *hb.Tag {
	b := lo.FirstOr(breadcrumbs, []layouts.Breadcrumb{
		{
			Name: "Home",
			URL:  links.User().Home(),
		},
	})
	// if first is not home, add home
	if len(b) > 0 && b[0].Name != "Home" {
		b = append([]layouts.Breadcrumb{
			{
				Name: "Home",
				URL:  links.User().Home(),
			},
		}, b...)
	}
	breadcrumbsTag := layouts.Breadcrumbs(b)

	icon := hb.NewDiv().
		Class("d-flex align-items-center justify-content-center rounded-circle bg-white text-primary me-3").
		Style("width:56px;height:56px;").
		Children([]hb.TagInterface{
			hb.I().Class("bi " + iconName).Style("font-size:28px;"),
		})

	heading := hb.NewHeading1().
		Class("h3 mb-0 text-white").
		HTML(title)

	layout := hb.NewBorderLayout().
		AddLeft(icon, hb.BORDER_LAYOUT_ALIGN_CENTER, hb.BORDER_LAYOUT_ALIGN_MIDDLE).
		AddCenter(hb.NewDiv().
			Children([]hb.TagInterface{
				hb.NewDiv().Class("small mb-1 breadcrumbs-wrapper text-white-50").Child(breadcrumbsTag),
				heading,
			}),
			hb.BORDER_LAYOUT_ALIGN_LEFT,
			hb.BORDER_LAYOUT_ALIGN_MIDDLE)

	id := "PageHeader" + uid.HumanUid()

	style := hb.NewStyle(`
		#` + id + ` .breadcrumb{
			background:none;
			margin-bottom:0;
			padding:0;
		}
		#` + id + ` .breadcrumb-item + .breadcrumb-item::before{
			color: rgba(255,255,255,0.6);
		}
		#` + id + ` .breadcrumb a{
			color: rgba(255,255,255,0.75);
			text-decoration:none;
		}
		#` + id + ` .breadcrumb a:hover{
			color: #ffffff;
			text-decoration:underline;
		}
	`)

	pageHeader := hb.NewSection().
		ID(id).
		Child(layout).
		Class("py-3 px-4 mb-4 rounded-4 shadow-sm bg-secondary text-white")

	return hb.Wrap(style, pageHeader)
}
