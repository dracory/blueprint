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

	mediaManager := rtr.NewRoute().
		SetName("Admin > Media Manager").
		SetPath(links.ADMIN_MEDIA).
		SetHandler(NewMediaManagerController(app).Handler)

	mediaManagerCatchAll := rtr.NewRoute().
		SetName("Admin > Media Manager > Catchall").
		SetPath(links.ADMIN_MEDIA + links.CATCHALL).
		SetHandler(NewMediaManagerController(app).Handler)

	return []rtr.RouteInterface{
		mediaManager,
		mediaManagerCatchAll,
	}, nil
}
