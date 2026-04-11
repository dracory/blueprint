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
	blog := rtr.NewRoute().
		SetName("Admin > Blog").
		SetPath(links.ADMIN_BLOG).
		SetHandler(NewBlogAdminController(registry).Handler)

	return []rtr.RouteInterface{
		blog,
	}, nil
}
