package admin

import (
	"project/internal/links"

	"github.com/dracory/rtr"
)

func Routes() []rtr.RouteInterface {
	cmsManager := rtr.NewRoute().
		SetName("Admin > CMS Manager").
		SetPath(links.ADMIN_CMS).
		SetHTMLHandler(NewCmsController().Handler)

	cmsNewManager := rtr.NewRoute().
		SetName("Admin > CMS New Manager").
		SetPath(links.ADMIN_CMS_NEW).
		SetHandler(NewCmsNewController().Handler)

	return []rtr.RouteInterface{
		cmsManager,
		cmsNewManager,
	}
}
