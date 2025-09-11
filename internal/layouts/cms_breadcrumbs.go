package layouts

import (
	"project/internal/links"

	"github.com/dracory/bs"
	"github.com/dracory/hb"
)

func websiteBreadcrumbs(path []bs.Breadcrumb) hb.TagInterface {
	breadcrumbsPath := []bs.Breadcrumb{
		{
			Name: "",
			URL:  links.Website().Home(),
			Icon: hb.I().Class("bi bi-house").Style("font-size: 16px; color: gray").ToHTML(),
		},
	}

	breadcrumbsPath = append(breadcrumbsPath, path...)

	breadcrumbs := bs.Breadcrumbs(breadcrumbsPath)

	return breadcrumbs
}

func NewWebsiteBreadcrumbsSection(path []bs.Breadcrumb) hb.TagInterface {
	return hb.Section().
		ID("SectionBreadcrumbs").
		Child(websiteBreadcrumbs(path)).
		Style("margin-bottom:10px;")
}

func NewWebsiteBreadcrumbsSectionWithContainer(path []bs.Breadcrumb) hb.TagInterface {
	return hb.Section().
		ID("SectionBreadcrumbs").
		Child(
			hb.Div().
				Class("container").
				Child(websiteBreadcrumbs(path)),
		)
}
