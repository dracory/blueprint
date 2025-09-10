package account

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"project/internal/controllers/user/partials"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/types"

	"github.com/dracory/bs"
	"github.com/dracory/cdn"
	"github.com/dracory/geostore"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/dracory/userstore"
	"github.com/samber/lo"
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
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, errorMessage, links.NewUserLinks().Home(map[string]string{}), 10)
	}

	if data.action == controller.actionOnCountrySelectedTimezoneOptions {
		return controller.onCountrySelectedTimezoneOptions(data)
	}

	if r.Method == http.MethodPost {
		return controller.postUpdate(data)
	}

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.User().Home(map[string]string{})},
		{Name: "My Profile", URL: ""},
	})

	title := hb.Heading1().
		Text("My Account").
		Style("margin:30px 0px 30px 0px;")

	paragraph1 := hb.Paragraph().
		Text("Please keep your details updated so that we can contact you if you need our help.").
		Style("margin-bottom:20px;")

	formProfile := controller.formProfile(data)

	page := hb.Section().
		Child(hb.BR()).
		Child(breadcrumbs).
		Child(hb.HR()).
		Child(partials.UserQuickLinks(data.request)).
		Child(hb.HR()).
		Child(
			hb.Div().
				Class("container").
				Child(title).
				Child(paragraph1).
				Child(formProfile).
				Child(hb.BR()).
				Child(hb.BR()),
		)

	return layouts.NewUserLayout(controller.app, r, layouts.Options{
		Title:      "My Profile",
		Content:    page,
		ScriptURLs: []string{cdn.Sweetalert2_10()},
	}).ToHTML()
}

func (controller *profileController) postUpdate(data profileControllerData) string {
	if controller.app.GetConfig().GetVaultStoreUsed() && controller.app.GetVaultStore() == nil {
		data.formErrorMessage = "We are very sorry vault store is not configured. Saving the details not possible."
		return controller.formProfile(data).ToHTML()
	}

	// Basic validations for both modes
	if data.firstName == "" {
		data.formErrorMessage = "First name is required field"
		return controller.formProfile(data).ToHTML()
	}

	if data.lastName == "" {
		data.formErrorMessage = "Last name is required field"
		return controller.formProfile(data).ToHTML()
	}

	if data.email == "" {
		data.formErrorMessage = "Email is required field"
		return controller.formProfile(data).ToHTML()
	}

	if data.country == "" {
		data.formErrorMessage = "Country is required field"
		return controller.formProfile(data).ToHTML()
	}

	if data.timezone == "" {
		data.formErrorMessage = "Timezone is required field"
		return controller.formProfile(data).ToHTML()
	}

	if !controller.app.GetConfig().GetVaultStoreUsed() {
		// Direct write without tokenization
		data.authUser.SetFirstName(data.firstName)
		data.authUser.SetLastName(data.lastName)
		data.authUser.SetBusinessName(data.buinessName)
		data.authUser.SetPhone(data.phone)
	} else {
		// First name
		if err := controller.app.GetVaultStore().TokenUpdate(data.request.Context(), data.authUser.FirstName(), data.firstName, controller.app.GetConfig().GetVaultKey()); err != nil {
			controller.app.GetLogger().Error("Error saving first name", slog.String("error", err.Error()))
			data.formErrorMessage = "Saving profile failed. Please try again later."
			return controller.formProfile(data).ToHTML()
		}

		// Last name
		if err := controller.app.GetVaultStore().TokenUpdate(data.request.Context(), data.authUser.LastName(), data.lastName, controller.app.GetConfig().GetVaultKey()); err != nil {
			controller.app.GetLogger().Error("Error saving last name", slog.String("error", err.Error()))
			data.formErrorMessage = "Saving profile failed. Please try again later."
			return controller.formProfile(data).ToHTML()
		}

		// Business name
		if err := controller.app.GetVaultStore().TokenUpdate(data.request.Context(), data.authUser.BusinessName(), data.buinessName, controller.app.GetConfig().GetVaultKey()); err != nil {
			controller.app.GetLogger().Error("Error saving business name", slog.String("error", err.Error()))
			data.formErrorMessage = "Saving profile failed. Please try again later."
			return controller.formProfile(data).ToHTML()
		}

		// Phone
		if err := controller.app.GetVaultStore().TokenUpdate(data.request.Context(), data.authUser.Phone(), data.phone, controller.app.GetConfig().GetVaultKey()); err != nil {
			controller.app.GetLogger().Error("Error saving phone", slog.String("error", err.Error()))
			data.formErrorMessage = "Saving profile failed. Please try again later."
			return controller.formProfile(data).ToHTML()
		}
	}

	// Common updates
	data.authUser.SetCountry(data.country)
	data.authUser.SetTimezone(data.timezone)

	if controller.app.GetUserStore() == nil {
		controller.app.GetLogger().Warn("At profileController > post update. UserStore is nil.")
		data.formErrorMessage = "Saving profile failed. Please try again later."
		return controller.formProfile(data).ToHTML()
	}

	if err := controller.app.GetUserStore().UserUpdate(context.Background(), data.authUser); err != nil {
		controller.app.GetLogger().Error("Error updating user profile", slog.String("error", err.Error()))
		data.formErrorMessage = "Saving profile failed. Please try again later."
		return controller.formProfile(data).ToHTML()
	}

	data.formSuccessMessage = "Profile updated successfully"
	data.formRedirectURL = helpers.ToFlashSuccessURL(controller.app.GetCacheStore(), data.formSuccessMessage, links.User().Home(), 5)
	return controller.formProfile(data).ToHTML()
}

func (controller *profileController) formProfile(data profileControllerData) *hb.Tag {
	required := hb.Sup().
		Text("required").
		Style("margin-left:5px;color:lightcoral;")

	groupFirstName := bs.FormGroup().
		Child(bs.FormLabel("First name").
			Child(required)).
		Child(bs.FormInput().
			Name("first_name").
			Value(data.firstName))

	groupLastName := bs.FormGroup().
		Child(bs.FormLabel("Last name").
			Child(required)).
		Child(bs.FormInput().
			Name("last_name").
			Value(data.lastName))

	groupEmail := bs.FormGroup().
		Child(bs.FormLabel("Email").
			Child(required)).
		Child(bs.FormInput().
			Name("email").
			Value(data.email).
			Attr("readonly", "readonly").
			Style("background-color:#F8F8F8;"))

	groupBuinessName := bs.FormGroup().
		Child(bs.FormLabel("Company / buiness name")).
		Child(bs.FormInput().
			Name("business_name").
			Value(data.buinessName))

	groupPhone := bs.FormGroup().
		Child(bs.FormLabel("Phone")).
		Child(bs.FormInput().
			Name("phone").
			Value(data.phone))

	selectCountries := bs.FormSelect().
		ID("SelectCountries").
		Name(controller.formCountry).
		Child(bs.FormSelectOption("", "")).
		Children(lo.Map(data.countryList, func(country geostore.Country, _ int) hb.TagInterface {
			return bs.FormSelectOption(country.IsoCode2(), country.Name()).
				AttrIf(data.country == country.IsoCode2(), "selected", "selected")
		})).
		HxPost(links.NewAuthLinks().Register(map[string]string{
			"action": controller.actionOnCountrySelectedTimezoneOptions,
		})).
		HxTarget("#SelectTimezones").
		HxSwap("outerHTML")

	countryGroup := hb.Div().
		Class("form-group").
		Children([]hb.TagInterface{
			bs.FormLabel("Country").
				Child(required),
			selectCountries,
		})

	timezoneGroup := hb.Div().
		Class("form-group").
		Children([]hb.TagInterface{
			bs.FormLabel("Timezone").
				Child(required),
			controller.selectTimezoneByCountry(data.country, data.timezone),
		})

	buttonSave := bs.Button().
		Class("btn-primary mb-0").
		Attr("type", "button").
		Text("Save changes").
		HxInclude("#FormProfile").
		HxTarget("#CardUserProfile").
		HxTrigger("click").
		HxSwap("outerHTML").
		HxPost(links.NewUserLinks().Profile(map[string]string{}))

	formProfile := hb.Div().ID("FormProfile").Children([]hb.TagInterface{
		bs.Row().
			Class("g-4").
			Children([]hb.TagInterface{
				bs.Column(12).Child(groupEmail),
				bs.Column(6).Child(groupFirstName),
				bs.Column(6).Child(groupLastName),
				bs.Column(6).Child(groupBuinessName),
				bs.Column(6).Child(groupPhone),
				bs.Column(6).Child(countryGroup),
				bs.Column(6).Child(timezoneGroup),
			}),
		bs.Row().
			Class("mt-3").
			Child(
				bs.Column(12).
					Class("d-sm-flex justify-content-end").
					Child(buttonSave),
			),
	})

	return hb.Div().ID("CardUserProfile").
		Class("card bg-transparent border rounded-3").
		Style("text-align:left;").
		Children([]hb.TagInterface{
			hb.Div().Class("card-header  bg-transparent").Children([]hb.TagInterface{
				hb.Heading3().
					Text("Your Details").
					Style("text-align:left;font-size:23px;color:#333;"),
			}),
			hb.Div().Class("card-body").Children([]hb.TagInterface{
				formProfile,
			}),
		}).
		ChildIf(data.formErrorMessage != "", hb.Swal(hb.SwalOptions{
			Icon:              "error",
			Title:             "Error",
			Text:              data.formErrorMessage,
			ShowCancelButton:  false,
			ConfirmButtonText: "OK",
		})).
		ChildIf(data.formSuccessMessage != "", hb.Swal(hb.SwalOptions{
			Icon:              "success",
			Title:             "Saved",
			Text:              data.formSuccessMessage,
			ShowCancelButton:  false,
			ConfirmButtonText: "OK",
			ConfirmCallback:   "window.location.href = window.location.href",
		})).
		ChildIf(data.formRedirectURL != "", hb.Script(`window.location.href = '`+data.formRedirectURL+`'`))

}

func (controller *profileController) onCountrySelectedTimezoneOptions(data profileControllerData) string {
	return controller.selectTimezoneByCountry(data.country, data.timezone).ToHTML()
}

func (controller *profileController) selectTimezoneByCountry(country string, selectedTimezone string) *hb.Tag {
	query := geostore.TimezoneQueryOptions{
		SortOrder: sb.ASC,
		OrderBy:   geostore.COLUMN_TIMEZONE,
	}

	if country != "" {
		query.CountryCode = country
	}

	timezones, errZones := controller.app.GetGeoStore().TimezoneList(query)

	if errZones != nil {
		controller.app.GetLogger().Error("Error listing timezones", slog.String("error", errZones.Error()))
		return hb.Text("Error listing timezones")
	}

	selectTimezones := bs.FormSelect().
		ID("SelectTimezones").
		Name(controller.formTimezone).
		Child(bs.FormSelectOption("", "")).
		Children(lo.Map(timezones, func(timezone geostore.Timezone, _ int) hb.TagInterface {
			return bs.FormSelectOption(timezone.Timezone(), timezone.Timezone()).
				AttrIf(selectedTimezone == timezone.Timezone(), "selected", "selected")
		}))

	return selectTimezones
}

func (controller *profileController) prepareData(r *http.Request) (data profileControllerData, errorMessage string) {
	authUser := helpers.GetAuthUser(r)

	if authUser == nil {
		return profileControllerData{}, "User not found"
	}

	countryList, errCountries := controller.app.GetGeoStore().CountryList(geostore.CountryQueryOptions{
		SortOrder: "asc",
		OrderBy:   geostore.COLUMN_NAME,
	})

	if errCountries != nil {
		controller.app.GetLogger().Error("Error listing countries", slog.String("error", errCountries.Error()))
		return profileControllerData{}, "Error listing countries"
	}

	email, firstName, lastName, buinessName, phone, err := controller.untokenizeProfileData(r.Context(), authUser)

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

func (controller *profileController) untokenizeProfileData(ctx context.Context, user userstore.UserInterface) (email string, firstName string, lastName string, businessName string, phone string, err error) {
	if user == nil {
		return "", "", "", "", "", errors.New("user is nil")
	}

	// Start with whatever is stored on the user
	email = user.Email()
	firstName = user.FirstName()
	lastName = user.LastName()
	businessName = user.BusinessName()
	phone = user.Phone()

	// If vault is not used, treat fields as plaintext
	if !controller.app.GetConfig().GetVaultStoreUsed() {
		return email, firstName, lastName, businessName, phone, nil
	}

	if controller.app.GetVaultStore() == nil {
		return "", "", "", "", "", errors.New("VaultStore is not initialized")
	}

	// Interpret existing values as tokens and untokenize
	emailToken := email
	firstNameToken := firstName
	lastNameToken := lastName
	businessNameToken := businessName
	phoneToken := phone

	if emailToken != "" {
		if email, err = controller.app.GetVaultStore().TokenRead(ctx, emailToken, controller.app.GetConfig().GetVaultKey()); err != nil {
			controller.app.GetLogger().Error("Error reading email", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	if firstNameToken != "" {
		if firstName, err = controller.app.GetVaultStore().TokenRead(ctx, firstNameToken, controller.app.GetConfig().GetVaultKey()); err != nil {
			controller.app.GetLogger().Error("Error reading first name", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	if lastNameToken != "" {
		if lastName, err = controller.app.GetVaultStore().TokenRead(ctx, lastNameToken, controller.app.GetConfig().GetVaultKey()); err != nil {
			controller.app.GetLogger().Error("Error reading last name", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	if businessNameToken != "" {
		if businessName, err = controller.app.GetVaultStore().TokenRead(ctx, businessNameToken, controller.app.GetConfig().GetVaultKey()); err != nil {
			controller.app.GetLogger().Error("Error reading business name", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	if phoneToken != "" {
		if phone, err = controller.app.GetVaultStore().TokenRead(ctx, phoneToken, controller.app.GetConfig().GetVaultKey()); err != nil {
			controller.app.GetLogger().Error("Error reading phone", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	return email, firstName, lastName, businessName, phone, nil
}

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
