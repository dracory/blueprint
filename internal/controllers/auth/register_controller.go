package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"
	"strings"

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

type registerController struct {
	registry                               registry.RegistryInterface
	actionOnCountrySelectedTimezoneOptions string
	formFirstName                          string
	formLastName                           string
	formBusinessName                       string
	formEmail                              string
	formPhone                              string
	formCountry                            string
	formTimezone                           string
}

type registerControllerData struct {
	action             string
	authUser           userstore.UserInterface
	email              string
	firstName          string
	lastName           string
	buinessName        string
	phone              string
	country            string
	timezone           string
	countryList        []geostore.Country
	formErrorMessage   string
	formSuccessMessage string
	formRedirectURL    string
}

// == CONSTRUCTOR =============================================================

func NewRegisterController(registry registry.RegistryInterface) *registerController {
	return &registerController{
		registry:                               registry,
		actionOnCountrySelectedTimezoneOptions: "on-country-selected-timezone-options",
		formCountry:                            "country",
		formTimezone:                           "timezone",
		formPhone:                              "phone",
		formEmail:                              "email",
		formFirstName:                          "first_name",
		formLastName:                           "last_name",
		formBusinessName:                       "buiness_name",
	}
}

// == PUBLIC METHODS ==========================================================

func (controller *registerController) Handler(w http.ResponseWriter, r *http.Request) string {
	if !controller.registry.GetConfig().GetRegistrationEnabled() {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, `Registrations are currently disabled`, links.Website().Home(), 10)
	}

	if controller.registry.GetUserStore() == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, `user store is required`, links.Website().Home(), 5)
	}

	if controller.registry.GetConfig().GetUserStoreVaultEnabled() && controller.registry.GetVaultStore() == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, `vault store is required`, links.Website().Home(), 5)
	}

	data, errorMessage := controller.prepareData(r)

	if errorMessage != "" {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, errorMessage, links.Website().Home(), 10)
	}

	if data.action == controller.actionOnCountrySelectedTimezoneOptions {
		return controller.selectTimezoneByCountry(r.Context(), data.country, data.timezone).ToHTML()
	}

	if r.Method == http.MethodPost {
		return controller.postUpdate(r.Context(), data)
	}

	scripts := []string{}
	scriptURLs := []string{
		cdn.BootstrapJs_5_3_3(),
		cdn.Htmx_2_0_0(),
		cdn.Sweetalert2_11(),
	}

	if controller.registry.GetConfig().IsEnvProduction() {
		//scriptURLs = append([]string{helpers.GoogleTagScriptURL()}, scriptURLs...)
		//scripts = append(scripts, helpers.GoogleTagInitScript(), helpers.GoogleConversionScript())
	}

	return layouts.NewBlankLayout(
		controller.registry,
		r,
		layouts.Options{
			Title: "Register",
			// CanonicalURL: links.NewWebsiteLinks().Flash(map[string]string{}),
			Content:    controller.pageHTML(r.Context(), data),
			ScriptURLs: scriptURLs,
			Scripts:    scripts,
			StyleURLs:  []string{cdn.BootstrapIconsCss_1_11_3()},
			Styles: []string{`.Center > div{padding:0px !important;margin:0px !important;}
		@media (min-width: 576px) {.container.container-xs {max-width: 520px;}}
		body{background:rgba(128,0,128,0.05);}`,
				`#CardRegister{border-radius:24px;box-shadow:0 20px 60px rgba(33,37,41,0.08);overflow:hidden;}
		#CardRegister .card-header{padding:24px;border-bottom:1px solid rgba(0,0,0,0.05);background:#f8f9ff;}
		#CardRegister .card-header h3{font-size:14px;font-weight:600;letter-spacing:0.08em;color:#4b4b63;text-transform:uppercase;}
		#CardRegister .card-body{padding:32px;}
		#CardRegister .form-group{margin-bottom:18px;}
		#CardRegister .form-group label{display:flex;justify-content:space-between;align-items:center;font-size:13px;font-weight:600!important;color:#2b2b3f;text-transform:none;letter-spacing:0.02em;margin-bottom:6px;}
		#CardRegister .form-group label sup{font-size:12px;font-weight:500;color:#e26d78;margin-left:8px;}
		#CardRegister .form-control,#CardRegister .form-select{border-radius:14px;border-color:rgba(111,108,212,0.4);padding:12px 15px;transition:box-shadow 0.2s ease,border-color 0.2s ease;}
		#CardRegister .form-select{background-image:url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='14' height='14' viewBox='0 0 16 16'%3E%3Cpath fill='%234e73df' d='M4.646 6.146a.5.5 0 0 1 .708 0L8 8.793l2.646-2.647a.5.5 0 0 1 .708.708l-3 3a.5.5 0 0 1-.708 0l-3-3a.5.5 0 0 1 0-.708'/%3E%3C/svg%3E");background-repeat:no-repeat;background-position:right 1rem center;background-size:14px;}
		#CardRegister .form-control:focus{box-shadow:0 0 0 0.25rem rgba(78,115,223,0.2);border-color:#4e73df;}`},
		}).ToHTML()
}

// == PRIVATE METHODS =========================================================

func (controller *registerController) postUpdate(ctx context.Context, data registerControllerData) string {
	if controller.registry.GetUserStore() == nil {
		data.formErrorMessage = "We are very sorry user store is not configured. Saving the details not possible."
		return controller.formRegister(ctx, data).ToHTML()
	}

	if data.firstName == "" {
		data.formErrorMessage = "First name is required field"
		return controller.formRegister(ctx, data).ToHTML()
	}

	if data.lastName == "" {
		data.formErrorMessage = "Last name is required field"
		return controller.formRegister(ctx, data).ToHTML()
	}

	if data.country == "" {
		data.formErrorMessage = "Country is required field"
		return controller.formRegister(ctx, data).ToHTML()
	}

	if data.timezone == "" {
		data.formErrorMessage = "Timezone is required field"
		return controller.formRegister(ctx, data).ToHTML()
	}

	if controller.registry.GetConfig().GetUserStoreVaultEnabled() {
		if controller.registry.GetVaultStore() == nil {
			data.formErrorMessage = "We are very sorry vault store is not configured. Saving the details not possible."
			return controller.formRegister(ctx, data).ToHTML()
		}

		firstNameToken, err := controller.registry.GetVaultStore().TokenCreate(ctx, data.firstName, controller.registry.GetConfig().GetVaultStoreKey(), 20)

		if err != nil {
			data.formErrorMessage = "We are very sorry. Saving the details failed. Please try again later."
			return controller.formRegister(ctx, data).ToHTML()
		}

		lastNameToken, err := controller.registry.GetVaultStore().TokenCreate(ctx, data.lastName, controller.registry.GetConfig().GetVaultStoreKey(), 20)

		if err != nil {
			controller.registry.GetLogger().Error("Error creating last name token", slog.String("error", err.Error()))
			data.formErrorMessage = "We are very sorry. Saving the details failed. Please try again later."
			return controller.formRegister(ctx, data).ToHTML()
		}

		businessNameToken, err := controller.registry.GetVaultStore().TokenCreate(ctx, data.buinessName, controller.registry.GetConfig().GetVaultStoreKey(), 20)

		if err != nil {
			controller.registry.GetLogger().Error("Error creating business name token", slog.String("error", err.Error()))
			data.formErrorMessage = "We are very sorry. Saving the details failed. Please try again later."
			return controller.formRegister(ctx, data).ToHTML()
		}

		phoneToken, err := controller.registry.GetVaultStore().TokenCreate(ctx, data.phone, controller.registry.GetConfig().GetVaultStoreKey(), 20)

		if err != nil {
			controller.registry.GetLogger().Error("Error creating phone token", slog.String("error", err.Error()))
			data.formErrorMessage = "We are very sorry. Saving the details failed. Please try again later."
			return controller.formRegister(ctx, data).ToHTML()
		}

		data.authUser.SetFirstName(firstNameToken)
		data.authUser.SetLastName(lastNameToken)
		data.authUser.SetBusinessName(businessNameToken)
		data.authUser.SetPhone(phoneToken)
		data.authUser.SetCountry(data.country)
		data.authUser.SetTimezone(data.timezone)
	} else {
		data.authUser.SetFirstName(data.firstName)
		data.authUser.SetLastName(data.lastName)
		data.authUser.SetBusinessName(data.buinessName)
		data.authUser.SetPhone(data.phone)
		data.authUser.SetCountry(data.country)
		data.authUser.SetTimezone(data.timezone)
	}

	err := controller.registry.GetUserStore().UserUpdate(ctx, data.authUser)

	if err != nil {
		controller.registry.GetLogger().Error("Error updating user profile", slog.String("error", err.Error()))
		data.formErrorMessage = "We are very sorry. Saving the details failed. Please try again later."
		return controller.formRegister(ctx, data).ToHTML()
	}

	data.formSuccessMessage = "Your registration completed successfully. You can now continue browsing the website."
	data.formRedirectURL = links.User().Home()
	return controller.formRegister(ctx, data).ToHTML()
}

func (controller *registerController) pageHTML(ctx context.Context, data registerControllerData) hb.TagInterface {
	form := controller.formRegister(ctx, data)
	return hb.Div().
		Class(`container container-xs text-center`).
		Child(hb.BR()).
		Child(hb.BR()).
		Child(hb.Raw(layouts.LogoHTML())).
		Child(hb.BR()).
		Child(hb.BR()).
		Child(hb.Heading1().Text("Complete registration").Style(`font-size:24px;`)).
		Child(hb.BR()).
		Child(form).
		Child(hb.BR()).
		Child(hb.BR())
}

func (controller *registerController) formRegister(ctx context.Context, data registerControllerData) hb.TagInterface {
	required := hb.Sup().
		Text("required").
		Style("margin-left:5px;color:lightcoral;")

	buttonSave := bs.Button().
		Class("btn-primary mb-0 w-100 py-3 fs-5").
		Attr("type", "button").
		Child(hb.I().Class("bi bi-check-circle me-2")).
		Text("Save changes").
		HxInclude("#FormRegister").
		HxTarget("#CardRegister").
		HxTrigger("click").
		HxSwap("outerHTML").
		HxPost(links.Auth().Register())

	firstNameGroup := hb.Div().
		Class("form-group").
		Children([]hb.TagInterface{
			bs.FormLabel("First name").
				Child(required),
			bs.FormInput().
				Name(controller.formFirstName).
				Value(data.firstName),
		})

	lastNameGroup := hb.Div().
		Class("form-group").
		Children([]hb.TagInterface{
			bs.FormLabel("Last name").
				Child(required),
			bs.FormInput().
				Name(controller.formLastName).
				Value(data.lastName),
		})

	businessNameGroup := hb.Div().
		Class("form-group").
		Children([]hb.TagInterface{
			bs.FormLabel("Company / buiness name"),
			bs.FormInput().
				Name("business_name").
				Value(data.buinessName),
		})

	phoneGroup := hb.Div().
		Class("form-group").
		Children([]hb.TagInterface{
			bs.FormLabel("Phone"),
			bs.FormInput().
				Name("phone").
				Value(data.phone),
		})

	emailGroup := hb.Div().
		Class("form-group").
		Children([]hb.TagInterface{
			bs.FormLabel("Email").
				Child(required),
			bs.FormInput().
				Name("email").
				Value(data.email).
				Attr("readonly", "readonly").
				Style("background-color:#F8F8F8;"),
		})

	selectCountries := bs.FormSelect().
		ID("SelectCountries").
		Name(controller.formCountry).
		Child(bs.FormSelectOption("", "")).
		Children(lo.Map(data.countryList, func(country geostore.Country, _ int) hb.TagInterface {
			return bs.FormSelectOption(country.IsoCode2(), country.Name()).
				AttrIf(data.country == country.IsoCode2(), "selected", "selected")
		})).
		Hx("post", links.Auth().Register(map[string]string{
			"action": "on-country-selected-timezone-options",
		})).
		Hx("target", "#SelectTimezones").
		Hx("swap", "outerHTML")

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
			controller.selectTimezoneByCountry(ctx, data.country, data.timezone),
		})

	formProfile := hb.Div().
		ID("FormRegister").
		Child(
			bs.Row().
				Class("g-3").
				Children([]hb.TagInterface{
					bs.Column(12).
						Child(emailGroup),
					bs.Column(6).
						Child(firstNameGroup),
					bs.Column(6).
						Child(lastNameGroup),
					bs.Column(6).
						Child(businessNameGroup),
					bs.Column(6).
						Child(phoneGroup),
					bs.Column(6).
						Child(countryGroup),
					bs.Column(6).
						Child(timezoneGroup),
				}),
		).
		Child(
			bs.Row().Class("mt-4").Children([]hb.TagInterface{
				bs.Column(12).Class("d-sm-flex justify-content-end").
					Children([]hb.TagInterface{
						buttonSave,
					}),
			}),
		)

	return hb.Div().ID("CardRegister").
		Class("card bg-white border rounded-3").
		Style("text-align:left;").
		Children([]hb.TagInterface{
			hb.Div().Class("card-header  bg-transparent").Children([]hb.TagInterface{
				hb.Heading3().
					Text("Your Details").
					Style("text-align:left;font-size:12px;color:#333;margin:0px;"),
			}),
			hb.Div().Class("card-body").Children([]hb.TagInterface{
				formProfile,
			}),
		}).
		ChildIf(data.formErrorMessage != "", hb.Swal(hb.SwalOptions{
			Icon: "error",
			// Title:             "Oops...",
			Text:              data.formErrorMessage,
			ShowCancelButton:  false,
			ShowConfirmButton: false,
			Timer:             5000,
			TimerProgressBar:  true,
			Position:          "top-end",
		})).
		ChildIf(data.formSuccessMessage != "", hb.Swal(hb.SwalOptions{
			Icon:              "success",
			Title:             "Saved",
			Text:              data.formSuccessMessage,
			ShowCancelButton:  false,
			ShowConfirmButton: false,
			ConfirmCallback:   "window.location.href = window.location.href",
			Timer:             5000,
			TimerProgressBar:  true,
			Position:          "top-end",
		})).
		ChildIf(data.formRedirectURL != "", hb.Script(`window.location.href = '`+data.formRedirectURL+`'`))
}

func (controller *registerController) getUserData(ctx context.Context, user userstore.UserInterface) (email string, firstName string, lastName string, businessName string, phone string, err error) {
	if user == nil {
		return "", "", "", "", "", errors.New("user is nil")
	}

	email = user.Email()
	firstName = user.FirstName()
	lastName = user.LastName()
	businessName = user.BusinessName()
	phone = user.Phone()

	if !controller.registry.GetConfig().GetUserStoreVaultEnabled() {
		return email, firstName, lastName, businessName, phone, nil
	}

	if controller.registry.GetVaultStore() == nil {
		return "", "", "", "", "", errors.New("vault store is nil")
	}

	// assign tokenized values
	emailToken := email
	firstNameToken := firstName
	lastNameToken := lastName
	businessNameToken := businessName
	phoneToken := phone

	if emailToken != "" {
		email, err = controller.registry.GetVaultStore().TokenRead(ctx, emailToken, controller.registry.GetConfig().GetVaultStoreKey())

		if err != nil {
			controller.registry.GetLogger().Error("Error reading email", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	if firstNameToken != "" {
		firstName, err = controller.registry.GetVaultStore().TokenRead(ctx, firstNameToken, controller.registry.GetConfig().GetVaultStoreKey())

		if err != nil {
			controller.registry.GetLogger().Error("Error reading first name", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	if lastNameToken != "" {
		lastName, err = controller.registry.GetVaultStore().TokenRead(ctx, lastNameToken, controller.registry.GetConfig().GetVaultStoreKey())

		if err != nil {
			controller.registry.GetLogger().Error("Error reading last name", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	if businessNameToken != "" {
		businessName, err = controller.registry.GetVaultStore().TokenRead(ctx, businessNameToken, controller.registry.GetConfig().GetVaultStoreKey())

		if err != nil {
			controller.registry.GetLogger().Error("Error reading business name", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	if phoneToken != "" {
		phone, err = controller.registry.GetVaultStore().TokenRead(ctx, phoneToken, controller.registry.GetConfig().GetVaultStoreKey())

		if err != nil {
			controller.registry.GetLogger().Error("Error reading phone", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	return email, firstName, lastName, businessName, phone, nil
}

func (controller *registerController) prepareData(r *http.Request) (data registerControllerData, errorMessage string) {
	if controller.registry.GetUserStore() == nil {
		return registerControllerData{}, "User store is nil"
	}

	action := req.GetStringTrimmed(r, "action")
	authUser := helpers.GetAuthUser(r)

	if authUser == nil {
		return registerControllerData{}, "You must be logged in to access this page"
	}

	if controller.registry.GetGeoStore() == nil {
		return registerControllerData{}, "Geo store is nil"
	}

	countries, errCountries := controller.registry.GetGeoStore().CountryList(r.Context(), geostore.CountryQueryOptions{
		SortOrder: "asc",
		OrderBy:   geostore.COLUMN_NAME,
	})

	if errCountries != nil {
		controller.registry.GetLogger().Error("Error listing countries", slog.String("error", errCountries.Error()))
		return registerControllerData{}, "Error listing countries"
	}

	email, firstName, lastName, businessName, phone, err := controller.getUserData(r.Context(), authUser)

	if r.Method == http.MethodGet {
		if err != nil {
			controller.registry.GetLogger().Error("Error reading email", slog.String("error", err.Error()))
			return registerControllerData{}, "Error reading email"
		}

		data = registerControllerData{
			action:      action,
			authUser:    authUser,
			email:       email,
			firstName:   firstName,
			lastName:    lastName,
			buinessName: businessName,
			phone:       phone,
			timezone:    authUser.Timezone(),
			country:     authUser.Country(),
			countryList: countries,
		}
	}

	if r.Method == http.MethodPost {
		data = registerControllerData{
			action:      action,
			authUser:    authUser,
			email:       email,
			firstName:   strings.TrimSpace(req.GetStringTrimmed(r, controller.formFirstName)),
			lastName:    strings.TrimSpace(req.GetStringTrimmed(r, controller.formLastName)),
			buinessName: strings.TrimSpace(req.GetStringTrimmed(r, controller.formBusinessName)),
			phone:       strings.TrimSpace(req.GetStringTrimmed(r, controller.formPhone)),
			timezone:    strings.TrimSpace(req.GetStringTrimmed(r, controller.formTimezone)),
			country:     strings.TrimSpace(req.GetStringTrimmed(r, controller.formCountry)),
			countryList: countries,
		}
	}

	return data, ""
}

func (controller *registerController) selectTimezoneByCountry(ctx context.Context, country string, selectedTimezone string) hb.TagInterface {
	query := geostore.TimezoneQueryOptions{
		SortOrder: sb.ASC,
		OrderBy:   geostore.COLUMN_TIMEZONE,
	}

	if country != "" {
		query.CountryCode = country
	}

	timezones, errZones := controller.registry.GetGeoStore().TimezoneList(ctx, query)

	if errZones != nil {
		controller.registry.GetLogger().Error("Error listing timezones", slog.String("error", errZones.Error()))
		return hb.Text("Error listing timezones")
	}

	selectTimezones := bs.FormSelect().
		ID("SelectTimezones").
		Name("timezone").
		Child(bs.FormSelectOption("", "")).
		Children(lo.Map(timezones, func(timezone geostore.Timezone, _ int) hb.TagInterface {
			return bs.FormSelectOption(timezone.Timezone(), timezone.Timezone()).
				AttrIf(selectedTimezone == timezone.Timezone(), "selected", "selected")
		}))

	return selectTimezones
}
