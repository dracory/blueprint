package liveflux

import (
	"net/http"
	"net/http/httptest"

	"github.com/dracory/liveflux"
)

// Controller adapts liveflux.Handler to the rtr HTML handler signature.
type Controller struct {
	Engine http.Handler
}

func NewController() *Controller {
	return &Controller{Engine: liveflux.NewHandler(nil)}
}

// Handler returns the rendered HTML string for the component action/mount.
func (c *Controller) Handler(w http.ResponseWriter, r *http.Request) string {
	rec := httptest.NewRecorder()
	c.Engine.ServeHTTP(rec, r)
	// Propagate headers (e.g., redirect headers) to the real response
	for k, vv := range rec.Header() {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	// Prefer body regardless of status; caller can still show errors if needed.
	return rec.Body.String()
}
