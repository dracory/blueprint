package cms

import (
	"net/http"
	"project/internal/types"
	"project/internal/widgets"

	"github.com/dracory/base/req"
)

// == CONTROLLER ===============================================================

type widgetController struct {
	app types.AppInterface
}

// == CONSTRUCTOR ==============================================================

func NewWidgetController(app types.AppInterface) *widgetController {
	return &widgetController{app: app}
}

// == PUBLIC METHODS ==========================================================

func (controller *widgetController) Handler(w http.ResponseWriter, r *http.Request) string {
	alias := req.Value(r, "alias")

	if alias == "" {
		return "Widget type not specified"
	}

	widgetList := widgets.WidgetRegistry(controller.app.GetConfig())

	for _, widget := range widgetList {
		if widget.Alias() == alias {
			return widget.Render(r, "", nil)
		}
	}

	return alias
}
