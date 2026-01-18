package user

import (
	"log/slog"
	"net/http"
	"project/internal/ext"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"
	"strings"

	"github.com/dracory/userstore"

	"github.com/dracory/hb"
)

// == CONTROLLER ==============================================================

type homeController struct {
	registry registry.RegistryInterface
}

// == CONSTRUCTOR =============================================================

func NewHomeController(registry registry.RegistryInterface) *homeController {
	return &homeController{registry: registry}
}

// == PUBLIC METHODS ==========================================================

func (controller *homeController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, errorMessage, links.User().Home(map[string]string{}), 10)
	}

	return layouts.NewUserLayout(controller.registry, r, layouts.Options{
		Title:   "Home",
		Content: controller.view(data),
	}).ToHTML()
}

func (controller *homeController) view(data homeControllerData) hb.TagInterface {
	userName := data.userFirstName + " " + data.userLastName

	if strings.TrimSpace(userName) == "" {
		userName = data.userEmail
	}

	return hb.Wrap().HTML("Hi, " + userName + ". You are in user dashboard")
}

func (controller *homeController) prepareData(r *http.Request) (data homeControllerData, errorMessage string) {
	var err error
	authUser := helpers.GetAuthUser(r)

	if authUser == nil {
		return data, "User not found"
	}

	userFirstName := authUser.FirstName()
	userLastName := authUser.LastName()
	userEmail := authUser.Email()

	if controller.registry.GetConfig().GetUserStoreVaultEnabled() {
		userFirstName, userLastName, userEmail, _, _, err = ext.UserUntokenize(r.Context(), controller.registry, controller.registry.GetConfig().GetVaultStoreKey(), authUser)

		if err != nil {
			controller.registry.GetLogger().Error("Error: user > home > prepareData", slog.String("error", err.Error()))
			return data, "User data failed to be fetched"
		}
	}

	return homeControllerData{
		request:       r,
		user:          authUser,
		userFirstName: userFirstName,
		userLastName:  userLastName,
		userEmail:     userEmail,
	}, ""
}

type homeControllerData struct {
	request       *http.Request
	user          userstore.UserInterface
	userFirstName string
	userLastName  string
	userEmail     string
}
