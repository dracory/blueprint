package user

import (
	"net/http"
	"project/app/layouts"
	"project/app/links"
	"project/config"
	"project/internal/helpers"
	"strings"

	"github.com/gouniverse/userstore"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/router"
	"github.com/samber/lo"
)

// == CONTROLLER ==============================================================

type homeController struct{}

// == CONSTRUCTOR =============================================================

func NewHomeController() router.HTMLControllerInterface {
	return &homeController{}
}

// == PUBLIC METHODS ==========================================================

func (controller *homeController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return helpers.ToFlashError(w, r, errorMessage, links.NewUserLinks().Home(map[string]string{}), 10)
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
	authUser := helpers.GetAuthUser(r)

	if authUser == nil {
		return data, "User not found"
	}

	untokenized, err := helpers.Untokenize(r.Context(), map[string]string{
		"first_name": authUser.FirstName(),
		"last_name":  authUser.LastName(),
		"email":      authUser.Email(),
	})

	if err != nil {
		config.LogStore.ErrorWithContext("At orderListController > prepareData", err.Error())
		return data, "User data failed to be fetched"
	}

	userFirstName := lo.ValueOr(untokenized, "first_name", "")
	userLastName := lo.ValueOr(untokenized, "last_name", "")
	userEmail := lo.ValueOr(untokenized, "email", "n/a")

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
