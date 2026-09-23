package register

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"html"
	"log/slog"
	"net/http"
	"project/internal/app"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	authrules "project/internal/rules/auth"
	"strings"

	baselayouts "github.com/dracory/base/layouts"
	basesession "github.com/dracory/base/session"
	"github.com/dracory/cdn"
	"github.com/dracory/geostore"
	"github.com/dracory/hb"
	"github.com/dracory/neat"
	"github.com/dracory/req"
	"github.com/dracory/userstore"
)

//go:embed app.html
var templateHTML string

//go:embed app.js
var appJS string

//go:embed app.css
var appCSS string

// == CONTROLLER ==============================================================

type registerController struct {
	app app.AppInterface
}

// == CONSTRUCTOR =============================================================

func NewRegisterController(app app.AppInterface) *registerController {
	return &registerController{app: app}
}

// == PUBLIC METHODS ==========================================================

// PageHandler renders the registration page (GET). Registered via rtr.GetHTML
// in routes.go. Guard failures redirect through the flash page as before.
func (controller *registerController) PageHandler(w http.ResponseWriter, r *http.Request) string {
	if flash := controller.guardFlash(w, r); flash != "" {
		return flash
	}

	authUser := basesession.GetAuthUser(r)

	countries, err := controller.countryList(r.Context())
	if err != nil {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, "Error listing countries", links.Website().Home(), 10)
	}

	email, firstName, lastName, businessName, phone, err := controller.getUserData(r.Context(), authUser)
	if err != nil {
		controller.app.GetLogger().Error("Error reading user data", slog.String("error", err.Error()))
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, "Error reading user data", links.Website().Home(), 10)
	}

	// Registration completion only makes sense for users created through an
	// email-based login. An empty email means the user record is broken.
	if email == "" {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, "Your account is missing an email address", links.Website().Home(), 10)
	}

	initialData, _ := json.Marshal(map[string]string{
		"email":         email,
		"first_name":    firstName,
		"last_name":     lastName,
		"business_name": businessName,
		"phone":         phone,
		"country":       authUser.GetCountry(),
		"timezone":      authUser.GetTimezone(),
	})

	type countryOption struct {
		Code string `json:"code"`
		Name string `json:"name"`
	}
	countryOptions := make([]countryOption, 0, len(countries))
	for _, c := range countries {
		countryOptions = append(countryOptions, countryOption{Code: c.IsoCode2(), Name: c.Name()})
	}
	countriesJSON, _ := json.Marshal(countryOptions)

	appName := "App"
	if controller.app.GetConfig() != nil && controller.app.GetConfig().GetAppName() != "" {
		appName = controller.app.GetConfig().GetAppName()
	}

	htmlContent := strings.ReplaceAll(templateHTML, "{{ appName }}", html.EscapeString(appName))

	script := `
	const APP_NAME = ` + string(mustJSON(appName)) + `;
	const REGISTER_AJAX_URL = ` + string(mustJSON(links.Auth().Register(map[string]string{"action": "save"}))) + `;
	const TIMEZONES_AJAX_URL = ` + string(mustJSON(links.Auth().Register(map[string]string{"action": "timezones"}))) + `;
	const INITIAL_DATA = ` + string(initialData) + `;
	const COUNTRIES = ` + string(countriesJSON) + `;
	` + appJS

	return layouts.NewBlankLayout(controller.app, r, baselayouts.Options{
		AppName: appName,
		Title:   "Complete Registration",
		Content: hb.Div().HTML(htmlContent),
		StyleURLs: []string{
			cdn.BootstrapIconsCss_1_11_3(),
			cdn.Notiflix_3_2_8_CSS(),
		},
		ScriptURLs: []string{
			cdn.VueJs_3_5_32(),
			cdn.Notiflix_3_2_8(),
		},
		Scripts: []string{script},
		Styles:  []string{appCSS},
	}).ToHTML()
}

// AjaxHandler dispatches POST actions on the register path. Registered via
// rtr.PostJSON — returns a JSON string written with the application/json
// content type.
//
// Actions:
//   - "save"      → validates and persists the profile, returns redirect
//   - "timezones" → returns the timezone list for a country code
func (controller *registerController) AjaxHandler(w http.ResponseWriter, r *http.Request) string {
	if err := controller.guardError(w); err != "" {
		controller.sendErrorResponse(w, err, http.StatusForbidden)
		return ""
	}

	if basesession.GetAuthUser(r) == nil {
		controller.sendErrorResponse(w, "You must be logged in to access this page", http.StatusUnauthorized)
		return ""
	}

	action := r.URL.Query().Get("action")
	switch action {
	case "save":
		controller.handleSave(w, r)
	case "timezones":
		controller.handleTimezones(w, r)
	default:
		controller.sendErrorResponse(w, "Invalid action", http.StatusBadRequest)
	}
	return ""
}

// == PRIVATE METHODS =========================================================

// guardFlash runs the shared preconditions and returns a flash-redirect HTML
// string on failure ("" when everything is OK). Used by PageHandler.
func (controller *registerController) guardFlash(w http.ResponseWriter, r *http.Request) string {
	if msg := controller.guardError(w); msg != "" {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, msg, links.Website().Home(), 10)
	}

	if basesession.GetAuthUser(r) == nil {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r,
			"You must be logged in to access this page", links.Website().Home(), 10)
	}

	if controller.app.IsDisabledGeoStore() {
		return helpers.ToFlashError(controller.app.GetCacheStore(), w, r, "Geo store is nil", links.Website().Home(), 10)
	}

	return ""
}

// guardError returns a human-readable error message when a precondition fails
// ("" when everything is OK). Shared by the page and AJAX paths.
func (controller *registerController) guardError(w http.ResponseWriter) string {
	canRegister := authrules.NewCanRegisterRule(controller.app, "")
	if canRegister.Fails() {
		return canRegister.FailMessageFirst()
	}

	if controller.app.IsDisabledUserStore() {
		return "user store is required"
	}

	if controller.app.GetConfig().GetUserStoreVaultEnabled() && controller.app.IsDisabledVaultStore() {
		return "vault store is required"
	}

	return ""
}

// handleSave validates the submitted profile and persists it.
func (controller *registerController) handleSave(w http.ResponseWriter, r *http.Request) {
	if controller.app.IsDisabledUserStore() {
		controller.sendErrorResponse(w, "We are very sorry user store is not configured. Saving the details not possible.", http.StatusInternalServerError)
		return
	}

	authUser := basesession.GetAuthUser(r)

	data := registerFormData{
		firstName:    strings.TrimSpace(req.GetStringTrimmed(r, "first_name")),
		lastName:     strings.TrimSpace(req.GetStringTrimmed(r, "last_name")),
		businessName: strings.TrimSpace(req.GetStringTrimmed(r, "business_name")),
		phone:        strings.TrimSpace(req.GetStringTrimmed(r, "phone")),
		country:      strings.TrimSpace(req.GetStringTrimmed(r, "country")),
		timezone:     strings.TrimSpace(req.GetStringTrimmed(r, "timezone")),
	}

	formValidation := authrules.NewRegisterFormValidationRule(authrules.RegisterFormData{
		FirstName: data.firstName,
		LastName:  data.lastName,
		Email:     req.GetStringTrimmed(r, "email"),
		Country:   data.country,
		Timezone:  data.timezone,
	})
	if formValidation.Fails() {
		controller.sendErrorResponse(w, formValidation.Message(), http.StatusBadRequest)
		return
	}

	if err := controller.applyProfile(r.Context(), authUser, data); err != nil {
		controller.app.GetLogger().Error("Error saving registration", slog.String("error", err.Error()))
		controller.sendErrorResponse(w, "We are very sorry. Saving the details failed. Please try again later.", http.StatusInternalServerError)
		return
	}

	if err := controller.app.GetUserStore().UserUpdate(r.Context(), authUser); err != nil {
		controller.app.GetLogger().Error("Error updating user profile", slog.String("error", err.Error()))
		controller.sendErrorResponse(w, "We are very sorry. Saving the details failed. Please try again later.", http.StatusInternalServerError)
		return
	}

	controller.sendJSONResponse(w, map[string]interface{}{
		"status":   "success",
		"message":  "Your registration completed successfully.",
		"redirect": links.User().Home(),
	}, http.StatusOK)
}

// handleTimezones returns the timezone list for the given country code.
func (controller *registerController) handleTimezones(w http.ResponseWriter, r *http.Request) {
	if controller.app.IsDisabledGeoStore() {
		controller.sendErrorResponse(w, "Geo store is nil", http.StatusInternalServerError)
		return
	}

	country := strings.TrimSpace(req.GetStringTrimmed(r, "country"))

	query := geostore.TimezoneQueryOptions{
		SortOrder: neat.SortAsc,
		OrderBy:   geostore.COLUMN_TIMEZONE,
	}
	if country != "" {
		query.CountryCode = country
	}

	timezones, err := controller.app.GetGeoStore().TimezoneList(r.Context(), query)
	if err != nil {
		controller.app.GetLogger().Error("Error listing timezones", slog.String("error", err.Error()))
		controller.sendErrorResponse(w, "Error listing timezones", http.StatusInternalServerError)
		return
	}

	list := make([]string, 0, len(timezones))
	for _, tz := range timezones {
		list = append(list, tz.Timezone())
	}

	controller.sendJSONResponse(w, map[string]interface{}{
		"status":    "success",
		"timezones": list,
	}, http.StatusOK)
}

// applyProfile sets the profile fields on the user, vault-tokenizing the
// sensitive fields when the vault store is enabled.
func (controller *registerController) applyProfile(ctx context.Context, user userstore.UserInterface, data registerFormData) error {
	if !controller.app.GetConfig().GetUserStoreVaultEnabled() {
		user.SetFirstName(data.firstName)
		user.SetLastName(data.lastName)
		user.SetBusinessName(data.businessName)
		user.SetPhone(data.phone)
		user.SetCountry(data.country)
		user.SetTimezone(data.timezone)
		return nil
	}

	if controller.app.IsDisabledVaultStore() {
		return errors.New("vault store is not configured")
	}

	tokenize := func(value string) (string, error) {
		return controller.app.GetVaultStore().TokenCreate(ctx, value, controller.app.GetConfig().GetVaultStoreKey(), 20)
	}

	firstNameToken, err := tokenize(data.firstName)
	if err != nil {
		return err
	}
	lastNameToken, err := tokenize(data.lastName)
	if err != nil {
		return err
	}
	businessNameToken, err := tokenize(data.businessName)
	if err != nil {
		return err
	}
	phoneToken, err := tokenize(data.phone)
	if err != nil {
		return err
	}

	user.SetFirstName(firstNameToken)
	user.SetLastName(lastNameToken)
	user.SetBusinessName(businessNameToken)
	user.SetPhone(phoneToken)
	user.SetCountry(data.country)
	user.SetTimezone(data.timezone)
	return nil
}

// countryList returns all countries from the geo store.
func (controller *registerController) countryList(ctx context.Context) ([]geostore.Country, error) {
	return controller.app.GetGeoStore().CountryList(ctx, geostore.CountryQueryOptions{
		SortOrder: "asc",
		OrderBy:   geostore.COLUMN_NAME,
	})
}

// getUserData reads the user profile, decrypting vault-tokenized fields when
// the vault store is enabled.
func (controller *registerController) getUserData(ctx context.Context, user userstore.UserInterface) (email string, firstName string, lastName string, businessName string, phone string, err error) {
	if user == nil {
		return "", "", "", "", "", errors.New("user is nil")
	}

	email = user.GetEmail()
	firstName = user.GetFirstName()
	lastName = user.GetLastName()
	businessName = user.GetBusinessName()
	phone = user.GetPhone()

	if !controller.app.GetConfig().GetUserStoreVaultEnabled() {
		return email, firstName, lastName, businessName, phone, nil
	}

	if controller.app.IsDisabledVaultStore() {
		return "", "", "", "", "", errors.New("vault store is nil")
	}

	// assign tokenized values
	emailToken := email
	firstNameToken := firstName
	lastNameToken := lastName
	businessNameToken := businessName
	phoneToken := phone

	if emailToken != "" {
		email, err = controller.app.GetVaultStore().TokenRead(ctx, emailToken, controller.app.GetConfig().GetVaultStoreKey())

		if err != nil {
			controller.app.GetLogger().Error("Error reading email", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	if firstNameToken != "" {
		firstName, err = controller.app.GetVaultStore().TokenRead(ctx, firstNameToken, controller.app.GetConfig().GetVaultStoreKey())

		if err != nil {
			controller.app.GetLogger().Error("Error reading first name", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	if lastNameToken != "" {
		lastName, err = controller.app.GetVaultStore().TokenRead(ctx, lastNameToken, controller.app.GetConfig().GetVaultStoreKey())

		if err != nil {
			controller.app.GetLogger().Error("Error reading last name", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	if businessNameToken != "" {
		businessName, err = controller.app.GetVaultStore().TokenRead(ctx, businessNameToken, controller.app.GetConfig().GetVaultStoreKey())

		if err != nil {
			controller.app.GetLogger().Error("Error reading business name", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	if phoneToken != "" {
		phone, err = controller.app.GetVaultStore().TokenRead(ctx, phoneToken, controller.app.GetConfig().GetVaultStoreKey())

		if err != nil {
			controller.app.GetLogger().Error("Error reading phone", slog.String("error", err.Error()))
			return "", "", "", "", "", err
		}
	}

	return email, firstName, lastName, businessName, phone, nil
}

// sendJSONResponse sends a JSON response
func (controller *registerController) sendJSONResponse(w http.ResponseWriter, data map[string]interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// sendErrorResponse sends a JSON error response
func (controller *registerController) sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	controller.sendJSONResponse(w, map[string]interface{}{
		"status":  "error",
		"message": message,
	}, statusCode)
}

// mustJSON marshals a value to JSON, falling back to "null" on error.
// encoding/json escapes <, > and & by default so the output is safe to embed
// inside a <script> block.
func mustJSON(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		return []byte("null")
	}
	return b
}
