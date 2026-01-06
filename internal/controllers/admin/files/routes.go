package admin

import (
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/rtr"
)

func Routes(app types.RegistryInterface) []rtr.RouteInterface {
	fileManager := rtr.NewRoute().
		SetName("Admin > File Manager").
		SetPath(links.ADMIN_FILE_MANAGER).
		SetHTMLHandler(NewFileManagerController(app).Handler)

	return []rtr.RouteInterface{
		fileManager,
	}
}
