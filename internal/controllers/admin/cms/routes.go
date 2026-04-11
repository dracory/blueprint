package admin

import (
	"errors"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/rtr"
)

func Routes(registry registry.RegistryInterface) ([]rtr.RouteInterface, error) {
	if registry == nil {
		return nil, errors.New("registry cannot be nil")
	}
	cmsManager := rtr.NewRoute().
		SetName("Admin > Cms Manager").
		SetPath(links.ADMIN_CMS).
		SetHandler(NewCmsNewController(registry).Handler)

	return []rtr.RouteInterface{
		cmsManager,
	}, nil
}
