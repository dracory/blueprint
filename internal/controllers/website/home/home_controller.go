package home

import (
	"net/http"
	"project/internal/types"
)

// == CONSTRUCTOR ==============================================================

func NewHomeController(app types.RegistryInterface) *homeController {
	return &homeController{
		app: app,
	}
}

// == CONTROLLER ===============================================================

type homeController struct {
	app types.RegistryInterface
}

// == PUBLIC METHODS ===========================================================

func (controller *homeController) Handler(w http.ResponseWriter, r *http.Request) string {
	return "You are at the website home page"
}
