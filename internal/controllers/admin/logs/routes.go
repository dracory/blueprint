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

	logs := rtr.NewRoute().
		SetName("Admin > Logs").
		SetPath(links.ADMIN_LOGS).
		SetHandler(NewLogsAdminController(app).Handler)

	logsCatchAll := rtr.NewRoute().
		SetName("Admin > Logs > Catchall").
		SetPath(links.ADMIN_LOGS + links.CATCHALL).
		SetHandler(NewLogsAdminController(app).Handler)

	return []rtr.RouteInterface{
		logs,
		logsCatchAll,
	}, nil
}
