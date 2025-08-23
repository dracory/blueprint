package admin

import (
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/rtr"
)

func Routes(app types.AppInterface) []rtr.RouteInterface {
	cmsManager := rtr.NewRoute().
		SetName("Admin > Cms Manager").
		SetPath(links.ADMIN_CMS_NEW).
		SetHandler(NewCmsNewController(app).Handler)

	cmsOldManager := rtr.NewRoute().
		SetName("Admin > CMS Old Manager").
		SetPath(links.ADMIN_CMS_NEW).
		SetHTMLHandler(NewCmsOldController(app).Handler)

	return []rtr.RouteInterface{
		cmsManager,
		cmsOldManager,
	}
}
