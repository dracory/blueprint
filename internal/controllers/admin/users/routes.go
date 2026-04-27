package users

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

	users := rtr.NewRoute().
		SetName("Admin > Users").
		SetPath(links.ADMIN_USERS).
		SetHandler(NewUsersAdminController(registry).Handler)

	return []rtr.RouteInterface{
		users,
	}, nil
}
