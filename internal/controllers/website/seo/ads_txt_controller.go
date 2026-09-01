package seo

import (
	"net/http"
)

type adsTxtController struct{}

// NewAdsTxtController creates a new instance of the adsTxtController struct.
//
// Returns:
// - *adsTxtController: a pointer to the newly created adsTxtController.
func NewAdsTxtController() *adsTxtController {
	return &adsTxtController{}
}

func (c adsTxtController) Handler(w http.ResponseWriter, r *http.Request) string {
	w.Header().Set("Content-Type", "text/plain")
	return "google.com, pub-8821108004642146, DIRECT, f08c47fec0942fa0"
}
