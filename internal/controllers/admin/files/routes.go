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
	fileManager := rtr.NewRoute().
		SetName("Admin > File Manager").
		SetPath(links.ADMIN_FILE_MANAGER).
		SetHandler(NewFileManagerController(registry).Handler)

	return []rtr.RouteInterface{
		fileManager,
	}, nil
}
