package cms

import (
	"net/http"
	"project/internal/registry"
	"project/internal/widgets"

	"github.com/dracory/req"
)

// == CONTROLLER ===============================================================

type widgetController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR ==============================================================

func NewWidgetController(registry registry.RegistryInterface) *widgetController {
	return &widgetController{registry: registry}
}

// == PUBLIC METHODS ==========================================================

func (controller *widgetController) Handler(w http.ResponseWriter, r *http.Request) string {
	alias := req.GetStringTrimmed(r, "alias")

	if alias == "" {
		return "Widget type not specified"
	}

	widgetList := widgets.WidgetRegistry(controller.registry)

	for _, widget := range widgetList {
		if widget.Alias() == alias {
			return widget.Render(r, "", nil)
		}
	}

	return alias
}
