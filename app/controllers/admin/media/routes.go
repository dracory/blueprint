package admin

import (
	"project/internal/links"

	"github.com/dracory/rtr"
)

func Routes() []rtr.RouteInterface {
	mediaManager := rtr.NewRoute().
		SetName("Admin > Media Manager").
		SetPath(links.ADMIN_MEDIA).
		SetHTMLHandler(NewMediaManagerController().AnyIndex)

	return []rtr.RouteInterface{
		mediaManager,
	}
}
