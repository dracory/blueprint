package home

import (
	"net/http"

	"github.com/gouniverse/router"
)

// == CONSTRUCTOR ==============================================================

func NewHomeController() router.HTMLControllerInterface {
	return &homeController{}
}

// == CONTROLLER ===============================================================

type homeController struct{}

// == PUBLIC METHODS ===========================================================

func (controller *homeController) Handler(w http.ResponseWriter, r *http.Request) string {
	return "You are at the website home page"
}
