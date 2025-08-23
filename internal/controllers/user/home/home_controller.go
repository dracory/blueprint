package user

import (
	"log/slog"
	"net/http"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"
	"strings"

	"github.com/gouniverse/userstore"

	"github.com/gouniverse/hb"
)

// == CONTROLLER ==============================================================

type homeController struct {
	app types.AppInterface
}

// == CONSTRUCTOR =============================================================

func NewHomeController(app types.AppInterface) *homeController {
	return &homeController{app: app}
}

// == PUBLIC METHODS ==========================================================

func (controller *homeController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, errorMessage, links.User().Home(map[string]string{}), 10)
	}

	return layouts.NewUserLayout(controller.app, r, layouts.Options{
		Request:    r,
		Title:      "Home | Client",
		Content:    controller.view(data),
		StyleURLs:  []string{},
		ScriptURLs: []string{},
		Scripts:    []string{},
		Styles:     []string{},
		VaultStore: controller.app.GetVaultStore(),
		VaultKey:   controller.app.GetConfig().GetVaultKey(),
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

	if controller.app.GetConfig().GetVaultStoreUsed() {
		userFirstName, userLastName, userEmail, err = helpers.UserUntokenized(r.Context(), controller.app, controller.app.GetConfig().GetVaultKey(), authUser)

		if err != nil {
			controller.app.GetLogger().Error("Error: user > home > prepareData", slog.String("error", err.Error()))
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
