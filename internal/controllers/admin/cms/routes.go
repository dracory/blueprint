package admin

import (
	"errors"
	"project/internal/app"
	"project/internal/links"

	"github.com/dracory/rtr"
)

func Routes(app app.AppInterface) ([]rtr.RouteInterface, error) {
	if app == nil {
		return nil, errors.New("app cannot be nil")
	}
	cmsManager := rtr.NewRoute().
		SetName("Admin > Cms Manager").
		SetPath(links.ADMIN_CMS).
		SetHandler(NewCmsNewController(app).Handler)

	return []rtr.RouteInterface{
		cmsManager,
	}, nil
}
