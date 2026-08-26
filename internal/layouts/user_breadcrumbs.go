package layouts

import (
	"project/internal/links"

	baselayouts "github.com/dracory/base/layouts"
	"github.com/dracory/hb"
)

// userBreadcrumbs generates the user breadcrumbs
// the first breadcrumb is always the dashboard
func userBreadcrumbs(path []baselayouts.Breadcrumb) hb.TagInterface {
	breadcrumbsPath := []baselayouts.Breadcrumb{
		{
			Name: "Dashboard",
			URL:  links.User().Home(),
			Icon: hb.I().Class("bi bi-speedometer").Style("font-size: 16px; color: gray").ToHTML(),
		},
	}

	breadcrumbsPath = append(breadcrumbsPath, path...)

	breadcrumbs := baselayouts.Breadcrumbs(breadcrumbsPath)

	return breadcrumbs
}

func NewUserBreadcrumbsSection(path []baselayouts.Breadcrumb) hb.TagInterface {
	return hb.Section().
		ID("SectionBreadcrumbs").
		Child(userBreadcrumbs(path)).
		Style("margin-bottom:10px;")
}

func NewUserBreadcrumbsSectionWithContainer(path []baselayouts.Breadcrumb) hb.TagInterface {
	return hb.Section().
		ID("SectionBreadcrumbs").
		Child(
			hb.Div().
				Class("container").
				Child(userBreadcrumbs(path)),
		)
}
