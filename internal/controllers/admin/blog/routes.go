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
	blog := rtr.NewRoute().
		SetName("Admin > Blog").
		SetPath(links.ADMIN_BLOG).
		SetHandler(NewBlogAdminController(app).Handler)

	return []rtr.RouteInterface{
		blog,
	}, nil
}
