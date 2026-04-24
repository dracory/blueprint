package logadmin

import (
	"errors"
	"net/http"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/req"
	"github.com/dracory/rtr"

	"project/pkg/logadmin/log_manager"
	"project/pkg/logadmin/shared"
)

func Routes(registry registry.RegistryInterface) ([]rtr.RouteInterface, error) {
	if registry == nil {
		return nil, errors.New("registry cannot be nil")
	}
	handler := func(w http.ResponseWriter, r *http.Request) string {
		controller := req.GetStringTrimmed(r, "controller")

		switch controller {
		case shared.CONTROLLER_LOG_MANAGER:
			return log_manager.NewLogManagerController(registry).Handler(w, r)
		}

		// Default to log manager
		return log_manager.NewLogManagerController(registry).Handler(w, r)
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
