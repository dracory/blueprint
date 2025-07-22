package user

import (
	"log/slog"
	"net/http"
	"project/app/layouts"
	"project/app/links"
	"project/config"
	"project/internal/helpers"
	"strings"

	"github.com/gouniverse/userstore"

	"github.com/gouniverse/hb"
)

// == CONTROLLER ==============================================================

type homeController struct{}

// == CONSTRUCTOR =============================================================

func NewHomeController() *homeController {
	return &homeController{}
}

// == PUBLIC METHODS ==========================================================

func (controller *homeController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return helpers.ToFlashError(w, r, errorMessage, links.User().Home(map[string]string{}), 10)
	}

	return layouts.NewUserLayout(r, layouts.Options{
		Request:    r,
		Title:      "Home | Client",
		Content:    controller.view(data),
		StyleURLs:  []string{},
		ScriptURLs: []string{},
		Scripts:    []string{},
		Styles:     []string{},
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

	if config.VaultStoreUsed {
		userFirstName, userLastName, userEmail, err = helpers.UserUntokenized(r.Context(), authUser)

		if err != nil {
			config.Logger.Error("Error: user > home > prepareData", slog.String("error", err.Error()))
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
