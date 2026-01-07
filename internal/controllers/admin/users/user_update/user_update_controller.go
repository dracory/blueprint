package admin

import (
	"log/slog"
	"net/http"
	"strings"

	"project/internal/ext"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"

	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dracory/req"
	"github.com/dracory/userstore"
)

type userUpdateController struct {
	registry registry.RegistryInterface
}

func NewUserUpdateController(registry registry.RegistryInterface) *userUpdateController {
	return &userUpdateController{registry: registry}
}

func (controller *userUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, errorMessage, links.Admin().UsersUserManager(), 10)
	}

	rendered := liveflux.SSR(NewFormUserUpdate(controller.registry), map[string]string{
		"user_id":    data.userID,
		"return_url": data.returnURL,
	})

	if rendered == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "Error rendering user form", links.Admin().UsersUserManager(), 10)
	}

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Edit User | Users",
		Content: controller.page(data, rendered),
		ScriptURLs: []string{
			cdn.Sweetalert2_10(),
		},
		Scripts: []string{
			liveflux.Script().ToHTML(),
		},
	}).ToHTML()
}

func (controller *userUpdateController) page(data userUpdateControllerData, component hb.TagInterface) hb.TagInterface {
	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{
			Name: "Home",
			URL:  links.Admin().Home(),
		},
		{
			Name: "User Manager",
			URL:  links.Admin().UsersUserManager(),
		},
		{
			Name: "Edit User",
			URL:  links.Admin().UsersUserUpdate(map[string]string{"user_id": data.userID}),
		},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(links.Admin().UsersUserManager())

	heading := hb.Heading1().
		HTML("Edit User").
		Child(buttonCancel)

	card := hb.Div().
		Class("card").
		Child(
			hb.Div().
				Class("card-header").
				Style(`display:flex;justify-content:space-between;align-items:center;`).
				Child(hb.Heading4().
					HTML("User Details").
					Style("margin-bottom:0;display:inline-block;")),
		).
		Child(
			hb.Div().
				Class("card-body").
				Child(component),
		)

	userTitle := hb.Heading2().
		Class("mb-3").
		Text("User: ").
		Text(data.userDisplayName)

	return layouts.AdminPage(
		breadcrumbs,
		hb.HR(),
		heading,
		userTitle,
		card,
	)
}

func (controller *userUpdateController) prepareData(r *http.Request) (data userUpdateControllerData, errorMessage string) {
	if controller.registry.GetUserStore() == nil {
		return data, "User store is not configured"
	}

	data.userID = req.GetStringTrimmed(r, "user_id")

	if data.userID == "" {
		return data, "User ID is required"
	}

	user, err := controller.registry.GetUserStore().UserFindByID(r.Context(), data.userID)

	if err != nil {
		if controller.registry.GetLogger() != nil {
			controller.registry.GetLogger().Error("At userUpdateController > prepareData", slog.String("error", err.Error()))
		}
		return data, "User not found"
	}

	if user == nil {
		return data, "User not found"
	}

	firstName := user.FirstName()
	lastName := user.LastName()
	email := user.Email()

	if controller.registry.GetConfig().GetVaultStoreUsed() && controller.registry.GetVaultStore() != nil {
		firstName, lastName, email, _, _, err = ext.UserUntokenize(r.Context(), controller.registry, controller.registry.GetConfig().GetVaultStoreKey(), user)

		if err != nil {
			if controller.registry.GetLogger() != nil {
				controller.registry.GetLogger().Error("At userUpdateController > prepareData", slog.String("error", err.Error()))
			}
			return data, "Tokens failed to be read"
		}
	}

	data.user = user
	data.userDisplayName = strings.TrimSpace(firstName + " " + lastName)
	if data.userDisplayName == "" {
		data.userDisplayName = user.ID()
	}
	data.returnURL = links.Admin().UsersUserManager()

	_ = email

	return data, ""
}

type userUpdateControllerData struct {
	userID          string
	user            userstore.UserInterface
	userDisplayName string
	returnURL       string
}
