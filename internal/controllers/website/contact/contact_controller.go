package contact

import (
	"net/http"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/registry"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
)

// == CONTROLLER ==============================================================

type contactController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewContactController(registry registry.RegistryInterface) *contactController {
	return &contactController{registry: registry}
}

func (controller *contactController) AnyIndex(w http.ResponseWriter, r *http.Request) string {
	authUser := helpers.GetAuthUser(r)
	userID := ""
	if authUser != nil {
		userID = authUser.ID()
	}

	component := NewFormContact(controller.registry)
	rendered := liveflux.SSR(component, map[string]string{
		"user_id": userID,
	})

	title := hb.Heading1().HTML("Contact").Style("margin:30px 0px 30px 0px;")

	paragraph1 := hb.Paragraph().
		HTML("Please add and check your details below are correct so that we can respond to you as requested.").
		Style("margin-bottom:20px;")

	page := hb.Section().
		Child(
			hb.Div().
				Class("container").
				Child(title).
				Child(paragraph1).
				Child(rendered),
		)

	return layouts.NewUserLayout(controller.registry, r, layouts.Options{
		Title:   "Contact",
		Content: page,
		ScriptURLs: []string{
			cdn.Sweetalert2_11(),
			"/liveflux",
		},
		Scripts: []string{
			liveflux.Script().ToHTML(),
		},
	}).ToHTML()
}
