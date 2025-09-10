package page_not_found

import (
	"net/http"
)

// == CONSTRUCTOR =============================================================

func PageNotFoundController() *pageNotFoundController {
	return &pageNotFoundController{}
}

// == CONTROLLER ==============================================================

type pageNotFoundController struct{}

// PUBLIC METHODS =============================================================

func (controller *pageNotFoundController) Handler(w http.ResponseWriter, r *http.Request) string {
	w.WriteHeader(http.StatusNotFound)
	return "Sorry, page not found."
}
