package logadmin

import (
	"errors"
	"net/http"
	"project/internal/app"
	"project/internal/links"

	"github.com/dracory/req"
	"github.com/dracory/rtr"

	"project/pkg/logadmin/log_manager"
	"project/pkg/logadmin/shared"
)

func Routes(app app.AppInterface) ([]rtr.RouteInterface, error) {
	if app == nil {
		return nil, errors.New("app cannot be nil")
	}
	handler := func(w http.ResponseWriter, r *http.Request) string {
		controller := req.GetStringTrimmed(r, "controller")

		switch controller {
		case shared.CONTROLLER_LOG_MANAGER:
			return log_manager.NewLogManagerController(app).Handler(w, r)
		}

		// Default to log manager
		return log_manager.NewLogManagerController(app).Handler(w, r)
	}

	logs := rtr.NewRoute().
		SetName("Admin > Logs").
		SetPath(links.ADMIN_LOGS).
		SetHTMLHandler(handler)

	logsCatchAll := rtr.NewRoute().
		SetName("Admin > Logs > Catchall").
		SetPath(links.ADMIN_LOGS + links.CATCHALL).
		SetHTMLHandler(handler)

	return []rtr.RouteInterface{
		logs,
		logsCatchAll,
	}, nil
}
