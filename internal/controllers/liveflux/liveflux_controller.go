package liveflux

import (
	"context"
	"net/http"
	"net/http/httptest"
	"project/internal/controllers/user/account"
	"project/internal/controllers/user/discuss/chat"
	"project/internal/types"

	"github.com/dracory/liveflux"
)

// livefluxController adapts liveflux.Handler to the rtr HTML handler signature.
type livefluxController struct {
	Engine http.Handler
	App    types.AppInterface
}

func NewController(app types.AppInterface) *livefluxController {
	// Register all Liveflux components
	// Components are registered via liveflux.New() in their constructors
	account.NewFormProfileUpdate(app)
	chat.NewChatComponent(app)

	return &livefluxController{
		App:    app,
		Engine: liveflux.NewHandler(nil),
	}
}

// Handler returns the rendered HTML string for the component action/mount.
func (c *livefluxController) Handler(w http.ResponseWriter, r *http.Request) string {
	// add app to context
	ctx := context.WithValue(r.Context(), "app", c.App)
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
