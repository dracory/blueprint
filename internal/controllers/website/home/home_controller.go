package home

import (
	"net/http"
	"project/internal/types"

	"github.com/gouniverse/userstore"
)

// == CONSTRUCTOR ==============================================================

func NewHomeController(app types.AppInterface) *homeController {
	return &homeController{
		app: app,
	}
}

// == CONTROLLER ===============================================================

type homeController struct {
	app types.AppInterface
}

type homeControllerData struct {
	AuthUser userstore.UserInterface
}

// == PUBLIC METHODS ===========================================================

func (controller *homeController) Handler(w http.ResponseWriter, r *http.Request) string {
	return "You are at the website home page"
}
