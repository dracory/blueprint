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

	return []rtr.RouteInterface{
		cmsManager,
	}
}
