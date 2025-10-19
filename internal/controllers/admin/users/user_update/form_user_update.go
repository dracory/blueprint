package admin

import (
	"context"
	"log"
	"net/url"
	"strings"

	"project/internal/ext"
	"project/internal/links"
	"project/internal/types"

	"github.com/asaskevich/govalidator"
	"github.com/dracory/geostore"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dracory/sb"
	"github.com/dracory/userstore"
)

type formUserUpdate struct {
	liveflux.Base
	App            types.AppInterface
	UserID         string
	ReturnURL      string
	FormStatus     string
	FormFirstName  string
	FormLastName   string
	FormEmail      string
	FormMemo       string
	FormBusiness   string
	FormPhone      string
	FormCountry    string
	FormTimezone   string
	FormError      string
	FormSuccess    string
	FormRedirectTo string

	StatusOptions []userStatusOption
	DisplayName   string
	Countries     []geostore.Country
	Timezones     []geostore.Timezone
}

type userStatusOption struct {
	Key   string
	Label string
}

func newUserStatusOptions() []userStatusOption {
	return []userStatusOption{
		{Key: "", Label: "- not selected -"},
		{Key: userstore.USER_STATUS_ACTIVE, Label: "Active"},
		{Key: userstore.USER_STATUS_UNVERIFIED, Label: "Unverified"},
		{Key: userstore.USER_STATUS_INACTIVE, Label: "Inactive"},
		{Key: userstore.USER_STATUS_DELETED, Label: "In Trash Bin"},
	}
}

func NewFormUserUpdate(app types.AppInterface) liveflux.ComponentInterface {
	inst, err := liveflux.New(&formUserUpdate{})
	if err != nil {
		log.Println(err)
		return nil
	}
	if c, ok := inst.(*formUserUpdate); ok {
		c.App = app
	}
	return inst
}

func (c *formUserUpdate) GetAlias() string {
	return "admin_user_update_form"
}

func (c *formUserUpdate) Mount(ctx context.Context, params map[string]string) error {
	if c.App == nil {
		c.FormError = "Application not initialized"
		return nil
	}

	c.UserID = strings.TrimSpace(params["user_id"])
	if c.UserID == "" {
		c.FormError = "User ID is required"
		return nil
	}

	c.ReturnURL = strings.TrimSpace(params["return_url"])
	if c.ReturnURL == "" {
		c.ReturnURL = links.Admin().UsersUserManager()
	}

	if c.App.GetUserStore() == nil {
		c.FormError = "User store is not configured"
		return nil
	}

	user, err := c.App.GetUserStore().UserFindByID(ctx, c.UserID)
	if err != nil {
		if c.App.GetLogger() != nil {
			c.App.GetLogger().Error("Error loading user", "error", err.Error())
		}
		c.FormError = "Error loading user"
		return nil
	}

	if user == nil {
		c.FormError = "User not found"
		return nil
	}

	firstName := user.FirstName()
	lastName := user.LastName()
	email := user.Email()
	businessName := user.BusinessName()
	phone := user.Phone()

	if c.App.GetConfig().GetVaultStoreUsed() && c.App.GetVaultStore() != nil {
		firstName, lastName, email, businessName, phone, err = ext.UserUntokenize(ctx, c.App, c.App.GetConfig().GetVaultStoreKey(), user)
		if err != nil {
			if c.App.GetLogger() != nil {
				c.App.GetLogger().Error("Error untokenizing user", "error", err.Error())
			}
			c.FormError = "Error reading user details"
			return nil
		}
	}

	c.FormStatus = user.Status()
	c.FormFirstName = firstName
	c.FormLastName = lastName
	c.FormEmail = email
	c.FormMemo = user.Memo()
	c.FormBusiness = businessName
	c.FormPhone = phone
	c.FormCountry = user.Country()
	c.FormTimezone = user.Timezone()
	c.DisplayName = strings.TrimSpace(strings.Join([]string{firstName, lastName}, " "))
	c.StatusOptions = newUserStatusOptions()

	if c.App.GetGeoStore() == nil {
		c.FormError = "Geo store is not configured"
		return nil
	}

	countries, err := c.App.GetGeoStore().CountryList(geostore.CountryQueryOptions{
		SortOrder: sb.ASC,
		OrderBy:   geostore.COLUMN_NAME,
	})
	if err != nil {
		if c.App.GetLogger() != nil {
			c.App.GetLogger().Error("Error listing countries", "error", err.Error())
		}
		c.FormError = "Error listing countries"
		return nil
	}

	c.Countries = countries
	c.refreshTimezones()

	return nil
}

func (c *formUserUpdate) Handle(ctx context.Context, action string, data url.Values) error {
	switch action {
	case "apply", "save":
		return c.handleUpdate(ctx, action, data)
	case "country_change":
		if data == nil {
			data = url.Values{}
		}
		c.FormCountry = strings.TrimSpace(data.Get("user_country"))
		c.FormTimezone = ""
		c.refreshTimezones()
		c.FormSuccess = ""
		return nil
	default:
		return nil
	}
}

func (c *formUserUpdate) handleUpdate(ctx context.Context, action string, data url.Values) error {
	if data == nil {
		data = url.Values{}
	}

	userID := strings.TrimSpace(firstNonEmpty(data.Get("user_id"), c.UserID))
	if userID == "" {
		c.FormError = "User ID is required"
		c.FormSuccess = ""
		return nil
	}

	if c.App == nil || c.App.GetUserStore() == nil {
		c.FormError = "User store is not configured"
		c.FormSuccess = ""
		return nil
	}

	user, err := c.App.GetUserStore().UserFindByID(ctx, userID)
	if err != nil {
		if c.App.GetLogger() != nil {
			c.App.GetLogger().Error("Error loading user", "error", err.Error())
		}
		c.FormError = "Error loading user"
		c.FormSuccess = ""
		return nil
	}

	if user == nil {
		c.FormError = "User not found"
		c.FormSuccess = ""
		return nil
	}

	c.FormStatus = strings.TrimSpace(data.Get("user_status"))
	c.FormFirstName = strings.TrimSpace(data.Get("user_first_name"))
	c.FormLastName = strings.TrimSpace(data.Get("user_last_name"))
	c.FormEmail = strings.TrimSpace(data.Get("user_email"))
	c.FormBusiness = strings.TrimSpace(data.Get("user_business_name"))
	c.FormPhone = strings.TrimSpace(data.Get("user_phone"))
	c.FormMemo = strings.TrimSpace(data.Get("user_memo"))
	c.FormCountry = strings.TrimSpace(data.Get("user_country"))
	c.FormTimezone = strings.TrimSpace(data.Get("user_timezone"))

	c.refreshTimezones()

	if c.FormStatus == "" {
		c.FormError = "Status is required"
		c.FormSuccess = ""
		return nil
	}

	if c.FormFirstName == "" {
		c.FormError = "First name is required"
		c.FormSuccess = ""
		return nil
	}

	if c.FormLastName == "" {
		c.FormError = "Last name is required"
		c.FormSuccess = ""
		return nil
	}

	if c.FormEmail == "" {
		c.FormError = "Email is required"
		c.FormSuccess = ""
		return nil
	}

	if !govalidator.IsEmail(c.FormEmail) {
		c.FormError = "Invalid email address"
		c.FormSuccess = ""
		return nil
	}

	if c.FormCountry == "" {
		c.FormError = "Country is required"
		c.FormSuccess = ""
		return nil
	}

	if c.FormTimezone == "" {
		c.FormError = "Timezone is required"
		c.FormSuccess = ""
		return nil
	}

	user.SetMemo(c.FormMemo)
	user.SetStatus(c.FormStatus)
	user.SetCountry(c.FormCountry)
	user.SetTimezone(c.FormTimezone)

	if c.App.GetConfig().GetUserStoreVaultEnabled() && c.App.GetVaultStore() != nil {
		firstToken, lastToken, emailToken, phoneToken, businessToken, err := ext.UserTokenize(
			ctx,
			c.App.GetVaultStore(),
			c.App.GetConfig().GetVaultStoreKey(),
			user,
			c.FormFirstName,
			c.FormLastName,
			c.FormEmail,
			c.FormPhone,
			c.FormBusiness,
		)
		if err != nil {
			if c.App.GetLogger() != nil {
				c.App.GetLogger().Error("Error tokenizing user", "error", err.Error())
			}
			c.FormError = "System error. Saving user failed"
			c.FormSuccess = ""
			return nil
		}
		user.SetFirstName(firstToken)
		user.SetLastName(lastToken)
		user.SetEmail(emailToken)
		user.SetPhone(phoneToken)
		user.SetBusinessName(businessToken)
	} else {
		user.SetFirstName(c.FormFirstName)
		user.SetLastName(c.FormLastName)
		user.SetEmail(c.FormEmail)
		user.SetPhone(c.FormPhone)
		user.SetBusinessName(c.FormBusiness)
	}

	if err := c.App.GetUserStore().UserUpdate(ctx, user); err != nil {
		if c.App.GetLogger() != nil {
			c.App.GetLogger().Error("Error updating user", "error", err.Error())
		}
		c.FormError = "System error. Saving user failed"
		c.FormSuccess = ""
		return nil
	}

	c.FormError = ""
	c.FormSuccess = "User saved successfully"
	c.DisplayName = strings.TrimSpace(strings.Join([]string{c.FormFirstName, c.FormLastName}, " "))

	if action == "save" {
		c.FormRedirectTo = c.ReturnURL
	} else {
		c.FormRedirectTo = ""
	}

	return nil
}

func (c *formUserUpdate) Render(ctx context.Context) hb.TagInterface {
	alerts := hb.Div()
	if c.FormError != "" {
		alerts = alerts.Child(hb.SwalError(hb.SwalOptions{
			Text:             c.FormError,
			Timer:            5000,
			TimerProgressBar: true,
			Position:         "top-end",
		}))
	}
	if c.FormSuccess != "" {
		if c.FormRedirectTo != "" {
			alerts = alerts.Child(hb.SwalSuccess(hb.SwalOptions{
				Text:             c.FormSuccess,
				RedirectURL:      c.FormRedirectTo,
				RedirectSeconds:  5,
				Timer:            5000,
				TimerProgressBar: true,
				Position:         "top-end",
			}))
		} else {
			alerts = alerts.Child(hb.SwalSuccess(hb.SwalOptions{
				Text:             c.FormSuccess,
				Timer:            5000,
				TimerProgressBar: true,
				Position:         "top-end",
			}))
		}
	}

	statusSelect := hb.Select().
		Class("form-select").
		Name("user_status")
	for _, option := range c.StatusOptions {
		statusSelect = statusSelect.Child(
			hb.Option().
				Value(option.Key).
				Text(option.Label).
				AttrIf(option.Key == c.FormStatus, "selected", "selected"),
		)
	}

	buttonApply := hb.Button().
		Type("submit").
		Class("btn btn-primary").
		Attr("data-flux-action", "apply").
		Child(hb.I().Class("bi bi-check2 me-2")).
		Text("Apply")

	buttonSave := hb.Button().
		Type("submit").
		Class("btn btn-success").
		Attr("data-flux-action", "save").
		Child(hb.I().Class("bi bi-check2-all me-2")).
		Text("Save & Close")

	buttonCancel := hb.A().
		Href(c.ReturnURL).
		Class("btn btn-secondary").
		Child(hb.I().Class("bi bi-chevron-left me-2")).
		Text("Back")

	buildActions := func(marginClass string) hb.TagInterface {
		return hb.Div().
			Class(marginClass + " d-flex justify-content-between align-items-center flex-wrap gap-2").
			Child(buttonCancel).
			Child(hb.Div().Class("d-flex gap-2").
				Child(buttonApply).
				Child(buttonSave))
	}

	countrySelect := hb.Select().
		Class("form-select").
		Name("user_country").
		Attr("onchange", `document.getElementById('btnCountryChange').click();`)
	countrySelect = countrySelect.Child(
		hb.Option().
			Value("").
			Text("Select country").
			AttrIf(c.FormCountry == "", "selected", "selected"),
	)
	for _, country := range c.Countries {
		countrySelect = countrySelect.Child(
			hb.Option().
				Value(country.IsoCode2()).
				Text(country.Name()).
				AttrIf(country.IsoCode2() == c.FormCountry, "selected", "selected"),
		)
	}

	timezoneSelect := hb.Select().
		Class("form-select").
		Name("user_timezone")
	timezoneSelect = timezoneSelect.Child(
		hb.Option().
			Value("").
			Text("Select timezone").
			AttrIf(c.FormTimezone == "", "selected", "selected"),
	)
	for _, tz := range c.Timezones {
		timezoneSelect = timezoneSelect.Child(
			hb.Option().
				Value(tz.Timezone()).
				Text(tz.Timezone()).
				AttrIf(tz.Timezone() == c.FormTimezone, "selected", "selected"),
		)
	}

	buttonCountryChange := hb.Button().
		ID("btnCountryChange").
		Type("submit").
		Attr("data-flux-action", "country_change").
		Style("display:none;")

	form := hb.Form().
		ID("FormUserUpdate").
		Child(alerts).
		Child(hb.Input().Type("hidden").Name("user_id").Value(c.UserID)).
		Child(buttonCountryChange).
		Child(buildActions("mb-3")).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Status")).
			Child(statusSelect)).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("First Name")).
			Child(hb.Input().Class("form-control").Name("user_first_name").Value(c.FormFirstName))).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Last Name")).
			Child(hb.Input().Class("form-control").Name("user_last_name").Value(c.FormLastName))).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Email")).
			Child(hb.Input().Class("form-control").Name("user_email").Value(c.FormEmail))).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Business Name")).
			Child(hb.Input().Class("form-control").Name("user_business_name").Value(c.FormBusiness))).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Phone")).
			Child(hb.Input().Class("form-control").Name("user_phone").Value(c.FormPhone))).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Country")).
			Child(countrySelect)).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Timezone")).
			Child(timezoneSelect)).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Admin Notes")).
			Child(hb.Textarea().Class("form-control").Name("user_memo").Text(c.FormMemo))).
		Child(buildActions("mt-4"))

	content := hb.Div().
		Child(form)

	return c.Root(content)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func (c *formUserUpdate) refreshTimezones() {
	if c.App == nil || c.App.GetGeoStore() == nil {
		c.Timezones = nil
		if c.FormCountry == "" {
			c.FormTimezone = ""
		}
		return
	}

	query := geostore.TimezoneQueryOptions{
		SortOrder: sb.ASC,
		OrderBy:   geostore.COLUMN_TIMEZONE,
	}

	if c.FormCountry != "" {
		query.CountryCode = c.FormCountry
	}

	timezones, err := c.App.GetGeoStore().TimezoneList(query)
	if err != nil {
		if c.App.GetLogger() != nil {
			c.App.GetLogger().Error("Error listing timezones", "error", err.Error())
		}
		c.FormError = "Error listing timezones"
		return
	}

	c.Timezones = timezones

	if c.FormCountry == "" {
		c.FormTimezone = ""
		return
	}

	if c.FormTimezone == "" {
		return
	}

	found := false
	for _, tz := range timezones {
		if tz.Timezone() == c.FormTimezone {
			found = true
			break
		}
	}

	if !found {
		c.FormTimezone = ""
	}
}

func init() {
	if err := liveflux.Register(&formUserUpdate{}); err != nil {
		log.Printf("Failed to register formUserUpdate component: %v", err)
	}
}
