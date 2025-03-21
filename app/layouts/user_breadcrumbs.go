package layouts

import (
	"project/app/links"

	"github.com/gouniverse/bs"
	"github.com/gouniverse/hb"
	"github.com/gouniverse/icons"
)

// userBreadcrumbs generates the user breadcrumbs
// the first breadcrumb is always the dashboard
func userBreadcrumbs(path []bs.Breadcrumb) hb.TagInterface {
	breadcrumbsPath := []bs.Breadcrumb{
		{
			Name: "Dashboard",
			URL:  links.NewUserLinks().Home(map[string]string{}),
			Icon: icons.Icon("bi-speedometer", 16, 16, "gray").ToHTML(),
		},
	}

	breadcrumbsPath = append(breadcrumbsPath, path...)

	breadcrumbs := bs.Breadcrumbs(breadcrumbsPath)

	return breadcrumbs
}

func NewUserBreadcrumbsSection(path []bs.Breadcrumb) hb.TagInterface {
	return hb.Section().
		ID("SectionBreadcrumbs").
		Child(userBreadcrumbs(path)).
		Style("margin-bottom:10px;")
}

func NewUserBreadcrumbsSectionWithContainer(path []bs.Breadcrumb) hb.TagInterface {
	return hb.Section().
		ID("SectionBreadcrumbs").
		Child(
			hb.Div().
				Class("container").
				Child(userBreadcrumbs(path)),
		)
}
