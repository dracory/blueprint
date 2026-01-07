package admin

import (
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/rtr"
)

func Routes(registry registry.RegistryInterface) []rtr.RouteInterface {
	mediaManager := rtr.NewRoute().
		SetName("Admin > Media Manager").
		SetPath(links.ADMIN_MEDIA).
		SetHTMLHandler(NewMediaManagerController(registry).AnyIndex)

	return []rtr.RouteInterface{
		mediaManager,
	}
}
