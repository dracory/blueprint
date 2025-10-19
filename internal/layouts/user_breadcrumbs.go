package layouts

import (
	"project/internal/links"

	"github.com/dracory/hb"
)

// userBreadcrumbs generates the user breadcrumbs
// the first breadcrumb is always the dashboard
func userBreadcrumbs(path []Breadcrumb) hb.TagInterface {
	breadcrumbsPath := []Breadcrumb{
		{
			Name: "Dashboard",
			URL:  links.User().Home(),
			Icon: hb.I().Class("bi bi-speedometer").Style("font-size: 16px; color: gray").ToHTML(),
		},
	}

	breadcrumbsPath = append(breadcrumbsPath, path...)

	breadcrumbs := Breadcrumbs(breadcrumbsPath)

	return breadcrumbs
}

func NewUserBreadcrumbsSection(path []Breadcrumb) hb.TagInterface {
	return hb.Section().
		ID("SectionBreadcrumbs").
		Child(userBreadcrumbs(path)).
		Style("margin-bottom:10px;")
}

func NewUserBreadcrumbsSectionWithContainer(path []Breadcrumb) hb.TagInterface {
	return hb.Section().
		ID("SectionBreadcrumbs").
		Child(
			hb.Div().
				Class("container").
				Child(userBreadcrumbs(path)),
		)
}
