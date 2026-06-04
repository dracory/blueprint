package admin

import (
	"errors"
	"project/internal/links"
	"project/internal/app"

	"github.com/dracory/rtr"
)

func Routes(app app.AppInterface) ([]rtr.RouteInterface, error) {
	if app == nil {
		return nil, errors.New("app cannot be nil")
	}
	mediaManager := rtr.NewRoute().
		SetName("Admin > Media Manager").
		SetPath(links.ADMIN_MEDIA).
		SetHTMLHandler(NewMediaManagerController(app).AnyIndex)

	return []rtr.RouteInterface{
		mediaManager,
	}, nil
}
