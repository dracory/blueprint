package account

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"

	"project/internal/controllers/user/partials"
	"project/internal/ext"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/cdn"
	"github.com/dracory/geostore"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dracory/req"
	"github.com/dracory/userstore"
)

// == CONTROLLER ==============================================================

type profileController struct {
	app                                    types.AppInterface
	actionOnCountrySelectedTimezoneOptions string
	formCountry                            string
	formTimezone                           string
}

// == CONSTRUCTOR =============================================================

func NewProfileController(app types.AppInterface) *profileController {
	return &profileController{
		app:                                    app,
		actionOnCountrySelectedTimezoneOptions: "on-country-selected-timezone-options",
		formCountry:                            "country",
		formTimezone:                           "timezone",
	}
}

// == PUBLIC METHODS ==========================================================

func (controller *profileController) Handler(w http.ResponseWriter, r *http.Request) string {
	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, errorMessage, links.User().Home(), 10)
	}

	if data.action == controller.actionOnCountrySelectedTimezoneOptions {
		return controller.onCountrySelectedTimezoneOptions(data)
	}

	params := map[string]string{
		"user_id": data.authUser.ID(),
	}

	rendered := liveflux.SSR(
		NewFormProfileUpdate(controller.app),
		params,
	)

	if rendered == nil {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, "Error rendering profile form", links.User().Home(), 10)
	}

	pageHeader := partials.PageHeader("bi-person", "My Account", []layouts.Breadcrumb{
		{Name: "Dashboard", Icon: "bi-speedometer2", URL: links.User().Home()},
		{Name: "My Account", URL: links.User().Profile()},
	})

	page := hb.Section().
		Child(hb.Div().
			Class("container").
			Child(pageHeader)).
		Child(
			hb.Div().
				Class("container").
				Child(hb.Paragraph().Text("Please keep your details updated so that we can contact you if you need our help.").Style("margin-bottom:20px;")).
				Child(rendered).
				Child(hb.BR()).
				Child(hb.BR()),
		)

	return layouts.NewUserLayout(controller.app, r, layouts.Options{
		Title:   "My Account",
		Content: hb.NewDiv().Class("p-3").Child(page),
		ScriptURLs: []string{
			cdn.Sweetalert2_10(),
		},
		Scripts: []string{
			liveflux.Script().ToHTML(),
		},
	}).ToHTML()
}

func (controller *profileController) onCountrySelectedTimezoneOptions(data profileControllerData) string {
	component := NewFormProfileUpdate(controller.app)
	if component == nil {
		controller.app.GetLogger().Error("Error creating profile update form component for country change")
		return ""
	}

	form, ok := component.(*formProfileUpdate)
	if !ok {
		controller.app.GetLogger().Error("Unexpected component type for profile update form during country change")
		return ""
	}

	if err := form.Handle(context.Background(), "country_change", url.Values{
		"country":  {data.country},
		"timezone": {data.timezone},
	}); err != nil {
		controller.app.GetLogger().Error("Error handling country change", slog.String("error", err.Error()))
		return ""
	}

	return form.timezoneSelect().ToHTML()
}

func (controller *profileController) prepareData(r *http.Request) (data profileControllerData, errorMessage string) {
	authUser := helpers.GetAuthUser(r)

	if authUser == nil {
		return profileControllerData{}, "User not found"
	}

	countryList, err := controller.app.GetGeoStore().CountryList(geostore.CountryQueryOptions{
		SortOrder: "asc",
		OrderBy:   geostore.COLUMN_NAME,
	})

	if err != nil {
		controller.app.GetLogger().Error("Error listing countries", slog.String("error", err.Error()))
		return profileControllerData{}, "Error listing countries"
	}

	email, firstName, lastName, buinessName, phone, err := ext.UserUntokenizeTransparently(
		r.Context(),
		controller.app,
		authUser,
	)

	if err != nil {
		controller.app.GetLogger().Error("Error reading profile data", slog.String("error", err.Error()))
		return profileControllerData{}, "Error reading profile data"
	}

	data.request = r
	data.authUser = authUser
	data.countryList = countryList

	if r.Method == http.MethodGet {
		data.email = email
		data.firstName = firstName
		data.lastName = lastName
		data.buinessName = buinessName
		data.phone = phone
		data.timezone = authUser.Timezone()
		data.country = authUser.Country()
	}

	if r.Method == http.MethodPost {
		data.email = req.GetStringTrimmed(r, "email")
		data.firstName = req.GetStringTrimmed(r, "first_name")
		data.lastName = req.GetStringTrimmed(r, "last_name")
		data.buinessName = req.GetStringTrimmed(r, "business_name")
		data.phone = req.GetStringTrimmed(r, "phone")
		data.timezone = req.GetStringTrimmed(r, "timezone")
		data.country = req.GetStringTrimmed(r, "country")
	}

	return data, ""
}

// func (controller *profileController) untokenizeProfileData(ctx context.Context, user userstore.UserInterface) (email string, firstName string, lastName string, businessName string, phone string, err error) {
// 	if user == nil {
// 		return "", "", "", "", "", errors.New("user is nil")
// 	}

// 	// If vault is not used, treat fields as plaintext
// 	if controller.app.GetConfig().GetUserStoreVaultEnabled() {
// 		if controller.app.GetUserStore() == nil {
// 			return "", "", "", "", "", errors.New("UserStore is not initialized")
// 		}
// 		return ext.UserUntokenize(
// 			ctx,
// 			controller.app,
// 			controller.app.GetConfig().GetVaultStoreKey(),
// 			user)
// 	}

// 	email = user.Email()
// 	firstName = user.FirstName()
// 	lastName = user.LastName()
// 	businessName = user.BusinessName()
// 	phone = user.Phone()

// 	return email, firstName, lastName, businessName, phone, nil
// }

type profileControllerData struct {
	request            *http.Request
	action             string
	authUser           userstore.UserInterface
	email              string
	firstName          string
	lastName           string
	buinessName        string
	phone              string
	country            string
	countryList        []geostore.Country
	timezone           string
	formErrorMessage   string
	formSuccessMessage string
	formRedirectURL    string
}
