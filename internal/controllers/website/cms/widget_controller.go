package cms

import (
	"net/http"
	"project/internal/app"
	"project/internal/widgets"

	"github.com/dracory/req"
)

// == CONTROLLER ===============================================================

type widgetController struct {
	app app.AppInterface
}

// == CONSTRUCTOR ==============================================================

func NewWidgetController(app app.AppInterface) *widgetController {
	return &widgetController{app: app}
}

// == PUBLIC METHODS ==========================================================

func (controller *widgetController) Handler(w http.ResponseWriter, r *http.Request) string {
	alias := req.GetStringTrimmed(r, "alias")

	if alias == "" {
		return "Widget type not specified"
	}

	widgetList := widgets.WidgetRegistry(controller.app)

	for _, widget := range widgetList {
		if widget.Alias() == alias {
			return widget.Render(r, "", nil)
		}
	}

	return alias
}
