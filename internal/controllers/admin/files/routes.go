package admin

import (
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/rtr"
)

func Routes(registry registry.RegistryInterface) []rtr.RouteInterface {
	fileManager := rtr.NewRoute().
		SetName("Admin > File Manager").
		SetPath(links.ADMIN_FILE_MANAGER).
		SetHTMLHandler(NewFileManagerController(registry).Handler)

	return []rtr.RouteInterface{
		fileManager,
	}
}
