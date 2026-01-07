package admin

import (
	"net/http"
	"project/internal/controllers/admin/logs/log_manager"
	"project/internal/controllers/admin/logs/shared"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/req"
	"github.com/dracory/rtr"
)

func Routes(registry registry.RegistryInterface) []rtr.RouteInterface {
	handler := func(w http.ResponseWriter, r *http.Request) string {
		controller := req.GetStringTrimmed(r, "controller")

		switch controller {
		case shared.CONTROLLER_LOG_MANAGER:
			return log_manager.NewLogManagerController(registry).Handler(w, r)
		}

		// Default to post manager
		return log_manager.NewLogManagerController(registry).Handler(w, r)
	}

	blog := rtr.NewRoute().
		SetName("Admin > Logs").
		SetPath(links.ADMIN_LOGS).
		SetHTMLHandler(handler)

	blogCatchAll := rtr.NewRoute().
		SetName("Admin > Logs > Catchall").
		SetPath(links.ADMIN_LOGS + links.CATCHALL).
		SetHTMLHandler(handler)

	return []rtr.RouteInterface{
		blog,
		blogCatchAll,
	}
}
