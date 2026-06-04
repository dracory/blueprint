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

	logs := rtr.NewRoute().
		SetName("Admin > Logs").
		SetPath(links.ADMIN_LOGS).
		SetHandler(NewLogsAdminController(app).Handler)

	return []rtr.RouteInterface{
		logs,
	}, nil
}
