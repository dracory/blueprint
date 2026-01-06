package cms

import (
	"net/http"
	"project/internal/types"
	"project/internal/widgets"

	"github.com/dracory/req"
)

// == CONTROLLER ===============================================================

type widgetController struct {
	app types.RegistryInterface
}

// == CONSTRUCTOR ==============================================================

func NewWidgetController(app types.RegistryInterface) *widgetController {
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
