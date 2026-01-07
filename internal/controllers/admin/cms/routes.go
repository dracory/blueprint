package admin

import (
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/rtr"
)

func Routes(registry registry.RegistryInterface) []rtr.RouteInterface {
	cmsManager := rtr.NewRoute().
		SetName("Admin > Cms Manager").
		SetPath(links.ADMIN_CMS).
		SetHandler(NewCmsNewController(registry).Handler)

	cmsOldManager := rtr.NewRoute().
		SetName("Admin > CMS Old Manager").
		SetPath(links.ADMIN_CMS_OLD).
		SetHTMLHandler(NewCmsOldController(registry).Handler)

	return []rtr.RouteInterface{
		cmsManager,
		cmsOldManager,
	}
}
