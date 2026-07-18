package users

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

	users := rtr.NewRoute().
		SetName("Admin > Users").
		SetPath(links.ADMIN_USERS).
		SetHandler(NewUsersAdminController(app).Handler)

	return []rtr.RouteInterface{
		users,
	}, nil
}
