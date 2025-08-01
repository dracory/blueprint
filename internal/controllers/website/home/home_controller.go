package home

import (
	"net/http"
)

// == CONSTRUCTOR ==============================================================

func NewHomeController() *homeController {
	return &homeController{}
}

// == CONTROLLER ===============================================================

type homeController struct{}

// == PUBLIC METHODS ===========================================================

func (controller *homeController) Handler(w http.ResponseWriter, r *http.Request) string {
	return "You are at the website home page"
}
