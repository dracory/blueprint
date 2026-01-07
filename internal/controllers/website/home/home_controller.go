package home

import (
	"net/http"
	"project/internal/registry"
)

// == CONSTRUCTOR ==============================================================

func NewHomeController(registry registry.RegistryInterface) *homeController {
	return &homeController{
		registry: registry,
	}
}

// == CONTROLLER ===============================================================

type homeController struct {
	registry registry.RegistryInterface
}

// == PUBLIC METHODS ===========================================================

func (controller *homeController) Handler(w http.ResponseWriter, r *http.Request) string {
	return "You are at the website home page"
}
