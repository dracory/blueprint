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

	logs := rtr.NewRoute().
		SetName("Admin > Logs").
		SetPath(links.ADMIN_LOGS).
		SetHandler(NewLogsAdminController(registry).Handler)

	return []rtr.RouteInterface{
		logs,
	}, nil
}
