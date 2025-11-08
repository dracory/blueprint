package liveflux

import (
	"context"
	"net/http"
	"net/http/httptest"
	"project/internal/types"

	"github.com/dracory/liveflux"
)

// livefluxController adapts liveflux.Handler to the rtr HTML handler signature.
type livefluxController struct {
	Engine http.Handler
	App    types.AppInterface
}

type contextKey string

const AppContextKey contextKey = "app"

func NewController(app types.AppInterface) *livefluxController {
	return &livefluxController{
		App:    app,
		Engine: liveflux.NewHandler(nil),
	}
}

// Handler returns the rendered HTML string for the component action/mount.
func (c *livefluxController) Handler(w http.ResponseWriter, r *http.Request) string {
	// add app to context
	ctx := context.WithValue(r.Context(), AppContextKey, c.App)
	r = r.WithContext(ctx)

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
