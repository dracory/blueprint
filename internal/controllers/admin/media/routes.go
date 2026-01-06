package admin

import (
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/rtr"
)

func Routes(app types.RegistryInterface) []rtr.RouteInterface {
	mediaManager := rtr.NewRoute().
		SetName("Admin > Media Manager").
		SetPath(links.ADMIN_MEDIA).
		SetHTMLHandler(NewMediaManagerController(app).AnyIndex)

	return []rtr.RouteInterface{
		mediaManager,
	}
}
