package user_update

import (
	_ "embed"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"project/internal/ext"
	"project/internal/helpers"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"
	"project/pkg/useradmin/shared"

	"github.com/asaskevich/govalidator"
	"github.com/dracory/cdn"
	"github.com/dracory/geostore"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/dracory/taskstore"
)

var (
	//go:embed form.html
	formHTML string

	//go:embed form.js
	formJS string
)

type userUpdateController struct {
	registry registry.RegistryInterface
}

func NewUserUpdateController(registry registry.RegistryInterface) *userUpdateController {
	return &userUpdateController{registry: registry}
}

func (controller *userUpdateController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case "get-user":
		controller.handleGetUser(w, r)
		return ""
	case "get-timezones":
		controller.handleGetTimezones(w, r)
		return ""
	case "update-user":
		controller.handleUpdateUser(w, r)
		return ""
	default:
		return controller.renderPage(w, r)
	}
}

func (controller *userUpdateController) renderPage(w http.ResponseWriter, r *http.Request) string {
	userID := req.GetStringTrimmed(r, "user_id")

	if controller.registry.GetUserStore() == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "User store is not configured", shared.NewLinks("/admin/users").UserManager(), 10)
	}

	if userID == "" {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "User ID is required", shared.NewLinks("/admin/users").UserManager(), 10)
	}

	user, err := controller.registry.GetUserStore().UserFindByID(r.Context(), userID)
	if err != nil || user == nil {
		return helpers.ToFlashError(controller.registry.GetCacheStore(), w, r, "User not found", shared.NewLinks("/admin/users").UserManager(), 10)
	}

	firstName := user.GetFirstName()
	lastName := user.GetLastName()
	if controller.registry.GetConfig().GetVaultStoreUsed() && controller.registry.GetVaultStore() != nil {
		firstName, lastName, _, _, _, err = ext.UserUntokenize(r.Context(), controller.registry, controller.registry.GetConfig().GetVaultStoreKey(), user)
		if err != nil {
			if controller.registry.GetLogger() != nil {
				controller.registry.GetLogger().Error("At userUpdateController > renderPage", slog.String("error", err.Error()))
			}
		}
	}

	displayName := strings.TrimSpace(firstName + " " + lastName)
	if displayName == "" {
		displayName = user.GetID()
	}

	returnURL := shared.NewLinks("/admin/users").UserManager()
	urlGetUser := shared.NewLinks("/admin/users").UserUpdate(map[string]string{"action": "get-user", "user_id": userID})
	urlGetTimezones := shared.NewLinks("/admin/users").UserUpdate(map[string]string{"action": "get-timezones"})
	urlUpdateUser := shared.NewLinks("/admin/users").UserUpdate(map[string]string{"action": "update-user"})

	html := strings.ReplaceAll(formHTML, "USER_ID_PLACEHOLDER", "'"+userID+"'")
	html = strings.ReplaceAll(html, "RETURN_URL_PLACEHOLDER", "'"+returnURL+"'")
	js := strings.ReplaceAll(formJS, "USER_ID_PLACEHOLDER", "'"+userID+"'")
	js = strings.ReplaceAll(js, "RETURN_URL_PLACEHOLDER", "'"+returnURL+"'")
	js = strings.ReplaceAll(js, "urlGetUser", "'"+urlGetUser+"'")
	js = strings.ReplaceAll(js, "urlGetTimezones", "'"+urlGetTimezones+"'")
	js = strings.ReplaceAll(js, "urlUpdateUser", "'"+urlUpdateUser+"'")

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")
	appDiv := hb.Div().ID("app-user-update").Class("mt-3").HTML(html)

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.Admin().Home()},
		{Name: "User Manager", URL: shared.NewLinks("/admin/users").UserManager()},
		{Name: "Edit User", URL: shared.NewLinks("/admin/users").UserUpdate(map[string]string{"user_id": userID})},
	})

	buttonCancel := hb.Hyperlink().
		Class("btn btn-secondary ms-2 float-end").
		Child(hb.I().Class("bi bi-chevron-left").Style("margin-top:-4px;margin-right:8px;font-size:16px;")).
		HTML("Back").
		Href(shared.NewLinks("/admin/users").UserManager())

	heading := hb.Heading1().HTML("Edit User").Child(buttonCancel)

	userTitle := hb.Heading2().Class("mb-3").Text("User: ").Text(displayName)

	card := hb.Div().Class("card").Child(
		hb.Div().Class("card-header").Style("display:flex;justify-content:space-between;align-items:center;").
			Child(hb.Heading4().HTML("User Details").Style("margin-bottom:0;display:inline-block;")),
	).Child(
		hb.Div().Class("card-body").Child(vueCDN).Child(appDiv),
	)

	content := layouts.AdminPage(
		breadcrumbs,
		hb.HR(),
		heading,
		userTitle,
		card,
	)

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:      "Edit User | Users",
		Content:    content,
		ScriptURLs: []string{cdn.Sweetalert2_10()},
		Scripts:    []string{js},
	}).ToHTML()
}

func (controller *userUpdateController) jsonResponse(w http.ResponseWriter, success bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": success,
		"message": message,
		"data":    data,
	})
}

func (controller *userUpdateController) handleGetUser(w http.ResponseWriter, r *http.Request) {
	userID := req.GetStringTrimmed(r, "user_id")
	if userID == "" {
		controller.jsonResponse(w, false, "User ID is required", nil)
		return
	}

	user, err := controller.registry.GetUserStore().UserFindByID(r.Context(), userID)
	if err != nil || user == nil {
		controller.jsonResponse(w, false, "User not found", nil)
		return
	}

	firstName := user.GetFirstName()
	lastName := user.GetLastName()
	email := user.GetEmail()
	phone := user.GetPhone()
	business := user.GetBusinessName()
	memo := user.GetMemo()
	status := user.GetStatus()
	country := user.GetCountry()
	timezone := user.GetTimezone()

	fieldStatus := map[string]bool{
		"first_name":    true,
		"last_name":     true,
		"email":         true,
		"business_name": true,
		"phone":         true,
	}

	if controller.registry.GetConfig().GetVaultStoreUsed() && controller.registry.GetVaultStore() != nil {
		firstName, lastName, email, phone, business, err = ext.UserUntokenize(r.Context(), controller.registry, controller.registry.GetConfig().GetVaultStoreKey(), user)
		if err != nil {
			if controller.registry.GetLogger() != nil {
				controller.registry.GetLogger().Error("userUpdateController.handleGetUser UserUntokenize", slog.String("error", err.Error()))
			}
			fieldStatus["first_name"] = false
			fieldStatus["last_name"] = false
			fieldStatus["email"] = false
			fieldStatus["business_name"] = false
			fieldStatus["phone"] = false
			firstName = "n/a"
			lastName = "n/a"
			email = "n/a"
			phone = "n/a"
			business = "n/a"
		}
	}

	countryList, _ := controller.registry.GetGeoStore().CountryList(r.Context(), geostore.CountryQueryOptions{
		SortOrder: sb.ASC,
		OrderBy:   geostore.COLUMN_NAME,
	})
	countries := make([]map[string]string, 0, len(countryList))
	for _, c := range countryList {
		countries = append(countries, map[string]string{
			"iso_code_2": c.IsoCode2(),
			"name":       c.Name(),
		})
	}

	timezoneList, _ := controller.registry.GetGeoStore().TimezoneList(r.Context(), geostore.TimezoneQueryOptions{
		SortOrder:   sb.ASC,
		OrderBy:     geostore.COLUMN_TIMEZONE,
		CountryCode: country,
	})
	timezones := make([]map[string]string, 0, len(timezoneList))
	for _, tz := range timezoneList {
		timezones = append(timezones, map[string]string{
			"timezone": tz.Timezone(),
		})
	}

	controller.jsonResponse(w, true, "", map[string]interface{}{
		"status":        status,
		"first_name":    firstName,
		"last_name":     lastName,
		"email":         email,
		"business_name": business,
		"phone":         phone,
		"country":       country,
		"timezone":      timezone,
		"memo":          memo,
		"field_status":  fieldStatus,
		"countries":     countries,
		"timezones":     timezones,
	})
}

func (controller *userUpdateController) handleGetTimezones(w http.ResponseWriter, r *http.Request) {
	countryCode := req.GetStringTrimmed(r, "country_code")
	if countryCode == "" {
		controller.jsonResponse(w, false, "Country code is required", nil)
		return
	}

	timezoneList, err := controller.registry.GetGeoStore().TimezoneList(r.Context(), geostore.TimezoneQueryOptions{
		SortOrder:   sb.ASC,
		OrderBy:     geostore.COLUMN_TIMEZONE,
		CountryCode: countryCode,
	})
	if err != nil {
		controller.jsonResponse(w, false, "Failed to load timezones", nil)
		return
	}

	timezones := make([]map[string]string, 0, len(timezoneList))
	for _, tz := range timezoneList {
		timezones = append(timezones, map[string]string{
			"timezone": tz.Timezone(),
		})
	}

	controller.jsonResponse(w, true, "", map[string]interface{}{
		"timezones": timezones,
	})
}

func (controller *userUpdateController) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		UserID       string `json:"user_id"`
		Status       string `json:"status"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Email        string `json:"email"`
		BusinessName string `json:"business_name"`
		Phone        string `json:"phone"`
		Country      string `json:"country"`
		Timezone     string `json:"timezone"`
		Memo         string `json:"memo"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		controller.jsonResponse(w, false, "Invalid request body", nil)
		return
	}

	if payload.UserID == "" {
		controller.jsonResponse(w, false, "User ID is required", nil)
		return
	}

	user, err := controller.registry.GetUserStore().UserFindByID(r.Context(), payload.UserID)
	if err != nil || user == nil {
		controller.jsonResponse(w, false, "User not found", nil)
		return
	}

	if payload.Status == "" {
		controller.jsonResponse(w, false, "Status is required", nil)
		return
	}
	if strings.TrimSpace(payload.FirstName) == "" {
		controller.jsonResponse(w, false, "First name is required", nil)
		return
	}
	if strings.TrimSpace(payload.LastName) == "" {
		controller.jsonResponse(w, false, "Last name is required", nil)
		return
	}
	if strings.TrimSpace(payload.Email) == "" {
		controller.jsonResponse(w, false, "Email is required", nil)
		return
	}
	if !govalidator.IsEmail(strings.TrimSpace(payload.Email)) {
		controller.jsonResponse(w, false, "Invalid email address", nil)
		return
	}
	if payload.Country == "" {
		controller.jsonResponse(w, false, "Country is required", nil)
		return
	}
	if payload.Timezone == "" {
		controller.jsonResponse(w, false, "Timezone is required", nil)
		return
	}

	originalEmail := user.GetEmail()
	if controller.registry.GetConfig().GetVaultStoreUsed() && controller.registry.GetVaultStore() != nil {
		_, _, originalEmail, _, _, _ = ext.UserUntokenize(r.Context(), controller.registry, controller.registry.GetConfig().GetVaultStoreKey(), user)
	}

	user.SetMemo(strings.TrimSpace(payload.Memo))
	user.SetStatus(payload.Status)
	user.SetCountry(payload.Country)
	user.SetTimezone(payload.Timezone)

	if controller.registry.GetConfig().GetVaultStoreUsed() && controller.registry.GetVaultStore() != nil {
		firstToken, lastToken, emailToken, phoneToken, businessToken, err := ext.UserTokenize(
			r.Context(),
			controller.registry.GetVaultStore(),
			controller.registry.GetConfig().GetVaultStoreKey(),
			user,
			strings.TrimSpace(payload.FirstName),
			strings.TrimSpace(payload.LastName),
			strings.TrimSpace(payload.Email),
			strings.TrimSpace(payload.Phone),
			strings.TrimSpace(payload.BusinessName),
		)
		if err != nil {
			if controller.registry.GetLogger() != nil {
				controller.registry.GetLogger().Error("Error tokenizing user", slog.String("error", err.Error()))
			}
			controller.jsonResponse(w, false, "System error. Saving user failed", nil)
			return
		}
		user.SetFirstName(firstToken)
		user.SetLastName(lastToken)
		user.SetEmail(emailToken)
		user.SetPhone(phoneToken)
		user.SetBusinessName(businessToken)
	} else {
		user.SetFirstName(strings.TrimSpace(payload.FirstName))
		user.SetLastName(strings.TrimSpace(payload.LastName))
		user.SetEmail(strings.TrimSpace(payload.Email))
		user.SetPhone(strings.TrimSpace(payload.Phone))
		user.SetBusinessName(strings.TrimSpace(payload.BusinessName))
	}

	if err := controller.registry.GetUserStore().UserUpdate(r.Context(), user); err != nil {
		if controller.registry.GetLogger() != nil {
			controller.registry.GetLogger().Error("Error updating user", slog.String("error", err.Error()))
		}
		controller.jsonResponse(w, false, "System error. Saving user failed", nil)
		return
	}

	if controller.registry.GetConfig().GetVaultStoreUsed() && controller.registry.GetVaultStore() != nil {
		if originalEmail != strings.TrimSpace(payload.Email) {
			if controller.registry.GetTaskStore() != nil {
				_, err := controller.registry.GetTaskStore().TaskDefinitionEnqueueByAlias(
					r.Context(),
					taskstore.DefaultQueueName,
					"BlindIndexUpdate",
					map[string]any{
						"index":    "email",
						"truncate": "no",
					},
				)
				if err != nil {
					if controller.registry.GetLogger() != nil {
						controller.registry.GetLogger().Error("Error enqueuing blind index rebuild", slog.String("error", err.Error()))
					}
				}
			}
		}
	}

	controller.jsonResponse(w, true, "User saved successfully", nil)
}
