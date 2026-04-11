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
	mediaManager := rtr.NewRoute().
		SetName("Admin > Media Manager").
		SetPath(links.ADMIN_MEDIA).
		SetHTMLHandler(NewMediaManagerController(registry).AnyIndex)

	return []rtr.RouteInterface{
		mediaManager,
	}, nil
}
