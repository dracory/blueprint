package user_update

import (
	"context"
	"log"
	"net/url"
	"strings"

	"project/internal/ext"
	"project/internal/registry"
	"project/pkg/useradmin/shared"

	"github.com/asaskevich/govalidator"
	"github.com/dracory/geostore"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dracory/sb"
	"github.com/dracory/taskstore"
	"github.com/dracory/userstore"
)

type formUserUpdate struct {
	liveflux.Base
	registry                registry.RegistryInterface
	UserID                  string
	ReturnURL               string
	FormStatus              string
	FormFirstName           string
	FormLastName            string
	FormEmail               string
	FormMemo                string
	FormBusiness            string
	FormPhone               string
	FormCountry             string
	FormTimezone            string
	FormError               string
	FormSuccess             string
	FormRedirectTo          string
	FieldStatusFirstName    bool
	FieldStatusLastName     bool
	FieldStatusEmail        bool
	FieldStatusBusinessName bool
	FieldStatusPhone        bool

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

func NewFormUserUpdate(registry registry.RegistryInterface) liveflux.ComponentInterface {
	inst, err := liveflux.New(&formUserUpdate{})
	if err != nil {
		log.Println(err)
		return nil
	}
	if c, ok := inst.(*formUserUpdate); ok {
		c.registry = registry
	}
	return inst
}

func (c *formUserUpdate) GetKind() string {
	return "admin_user_update_form"
}

func (c *formUserUpdate) Mount(ctx context.Context, params map[string]string) error {
	if c.registry == nil {
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
		c.ReturnURL = shared.NewLinks("/admin/users").UserManager()
	}

	c.FieldStatusFirstName = params["field_status_first_name"] == "true"
	c.FieldStatusLastName = params["field_status_last_name"] == "true"
	c.FieldStatusEmail = params["field_status_email"] == "true"
	c.FieldStatusBusinessName = params["field_status_business_name"] == "true"
	c.FieldStatusPhone = params["field_status_phone"] == "true"

	if c.registry.GetUserStore() == nil {
		c.FormError = "User store is not configured"
		return nil
	}

	user, err := c.registry.GetUserStore().UserFindByID(ctx, c.UserID)
	if err != nil {
		if c.registry.GetLogger() != nil {
			c.registry.GetLogger().Error("Error loading user", "error", err.Error())
		}
		c.FormError = "Error loading user"
		return nil
	}

	if user == nil {
		c.FormError = "User not found"
		return nil
	}

	firstName := user.GetFirstName()
	lastName := user.GetLastName()
	email := user.GetEmail()
	businessName := user.GetBusinessName()
	phone := user.GetPhone()

	// Use untokenize to decrypt user fields
	if c.registry.GetConfig().GetVaultStoreUsed() && c.registry.GetVaultStore() != nil {
		firstName, lastName, email, businessName, phone, _ = ext.UserUntokenize(ctx, c.registry, c.registry.GetConfig().GetVaultStoreKey(), user)
		// Field status is already passed from controller, no need to re-check here
	}

	c.FormStatus = user.GetStatus()
	c.FormFirstName = firstName
	c.FormLastName = lastName
	c.FormEmail = email
	c.FormMemo = user.GetMemo()
	c.FormBusiness = businessName
	c.FormPhone = phone
	c.FormCountry = user.GetCountry()
	c.FormTimezone = user.GetTimezone()
	c.DisplayName = strings.TrimSpace(strings.Join([]string{firstName, lastName}, " "))
	c.StatusOptions = newUserStatusOptions()

	if c.registry.GetGeoStore() == nil {
		c.FormError = "Geo store is not configured"
		return nil
	}

	countries, err := c.registry.GetGeoStore().CountryList(ctx, geostore.CountryQueryOptions{
		SortOrder: sb.ASC,
		OrderBy:   geostore.COLUMN_NAME,
	})
	if err != nil {
		if c.registry.GetLogger() != nil {
			c.registry.GetLogger().Error("Error listing countries", "error", err.Error())
		}
		c.FormError = "Error listing countries"
		return nil
	}

	c.Countries = countries
	c.refreshTimezones(ctx)

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
		c.refreshTimezones(ctx)
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

	// Prevent updates when any field tokens are unreadable
	hasUnreadableFields := !c.FieldStatusFirstName || !c.FieldStatusLastName || !c.FieldStatusEmail || !c.FieldStatusBusinessName || !c.FieldStatusPhone
	if hasUnreadableFields {
		c.FormError = "Cannot update user: some field tokens are unreadable. Please resolve the vault key issue first."
		c.FormSuccess = ""
		return nil
	}

	userID := strings.TrimSpace(firstNonEmpty(data.Get("user_id"), c.UserID))
	if userID == "" {
		c.FormError = "User ID is required"
		c.FormSuccess = ""
		return nil
	}

	if c.registry == nil || c.registry.GetUserStore() == nil {
		c.FormError = "User store is not configured"
		c.FormSuccess = ""
		return nil
	}

	user, err := c.registry.GetUserStore().UserFindByID(ctx, userID)
	if err != nil {
		if c.registry.GetLogger() != nil {
			c.registry.GetLogger().Error("Error loading user", "error", err.Error())
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

	// Store original email for comparison (plaintext before tokenizing)
	originalEmail := user.GetEmail()

	c.FormStatus = strings.TrimSpace(data.Get("user_status"))
	c.FormFirstName = strings.TrimSpace(data.Get("user_first_name"))
	c.FormLastName = strings.TrimSpace(data.Get("user_last_name"))
	c.FormEmail = strings.TrimSpace(data.Get("user_email"))
	c.FormBusiness = strings.TrimSpace(data.Get("user_business_name"))
	c.FormPhone = strings.TrimSpace(data.Get("user_phone"))
	c.FormMemo = strings.TrimSpace(data.Get("user_memo"))
	c.FormCountry = strings.TrimSpace(data.Get("user_country"))
	c.FormTimezone = strings.TrimSpace(data.Get("user_timezone"))

	c.refreshTimezones(ctx)

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

	if c.registry.GetConfig().GetVaultStoreUsed() && c.registry.GetVaultStore() != nil {
		firstToken, lastToken, emailToken, phoneToken, businessToken, err := ext.UserTokenize(
			ctx,
			c.registry.GetVaultStore(),
			c.registry.GetConfig().GetVaultStoreKey(),
			user,
			c.FormFirstName,
			c.FormLastName,
			c.FormEmail,
			c.FormPhone,
			c.FormBusiness,
		)
		if err != nil {
			if c.registry.GetLogger() != nil {
				c.registry.GetLogger().Error("Error tokenizing user", "error", err.Error())
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

	if err := c.registry.GetUserStore().UserUpdate(ctx, user); err != nil {
		if c.registry.GetLogger() != nil {
			c.registry.GetLogger().Error("Error updating user", "error", err.Error())
		}
		c.FormError = "System error. Saving user failed"
		c.FormSuccess = ""
		return nil
	}

	// If email changed and vault is enabled, enqueue blind index rebuild for email
	if c.registry.GetConfig().GetVaultStoreUsed() && c.registry.GetVaultStore() != nil {
		if originalEmail != c.FormEmail {
			if c.registry.GetTaskStore() != nil {
				_, err := c.registry.GetTaskStore().TaskDefinitionEnqueueByAlias(
					ctx,
					taskstore.DefaultQueueName,
					"BlindIndexUpdate",
					map[string]any{
						"index":    "email",
						"truncate": "no",
					},
				)
				if err != nil {
					if c.registry.GetLogger() != nil {
						c.registry.GetLogger().Error("Error enqueuing blind index rebuild", "error", err.Error())
					}
				}
			}
		}
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

	// Warning when specific field tokens are unreadable
	var unreadableFields []string
	if !c.FieldStatusFirstName {
		unreadableFields = append(unreadableFields, "first name")
	}
	if !c.FieldStatusLastName {
		unreadableFields = append(unreadableFields, "last name")
	}
	if !c.FieldStatusEmail {
		unreadableFields = append(unreadableFields, "email")
	}
	if !c.FieldStatusBusinessName {
		unreadableFields = append(unreadableFields, "business name")
	}
	if !c.FieldStatusPhone {
		unreadableFields = append(unreadableFields, "phone")
	}

	if len(unreadableFields) > 0 {
		fieldsList := strings.Join(unreadableFields, ", ")
		alerts = alerts.Child(hb.Div().
			Class("alert alert-warning").
			Child(hb.I().Class("bi bi-exclamation-triangle-fill me-2")).
			Text("The following field tokens could not be decrypted: " + fieldsList + ". This may be due to a vault key mismatch. These fields show encrypted data and are disabled until the tokens can be read."))
	}

	statusSelect := hb.Select().
		Class("form-select").
		Name("user_status").
		AttrIf(!c.FieldStatusFirstName || !c.FieldStatusLastName || !c.FieldStatusEmail, "disabled", "disabled")
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
		AttrIf(!c.FieldStatusFirstName || !c.FieldStatusLastName || !c.FieldStatusEmail || !c.FieldStatusBusinessName || !c.FieldStatusPhone, "disabled", "disabled").
		Child(hb.I().Class("bi bi-check2 me-2")).
		Text("Apply")

	buttonSave := hb.Button().
		Type("submit").
		Class("btn btn-success").
		Attr("data-flux-action", "save").
		AttrIf(!c.FieldStatusFirstName || !c.FieldStatusLastName || !c.FieldStatusEmail || !c.FieldStatusBusinessName || !c.FieldStatusPhone, "disabled", "disabled").
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
		Attr("onchange", `document.getElementById('btnCountryChange').click();`).
		AttrIf(!c.FieldStatusFirstName || !c.FieldStatusLastName || !c.FieldStatusEmail, "disabled", "disabled")
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
		Name("user_timezone").
		AttrIf(!c.FieldStatusFirstName || !c.FieldStatusLastName || !c.FieldStatusEmail, "disabled", "disabled")
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
			Child(hb.Input().Class("form-control").Name("user_first_name").Value(c.FormFirstName).AttrIf(!c.FieldStatusFirstName, "disabled", "disabled"))).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Last Name")).
			Child(hb.Input().Class("form-control").Name("user_last_name").Value(c.FormLastName).AttrIf(!c.FieldStatusLastName, "disabled", "disabled"))).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Email")).
			Child(hb.Input().Class("form-control").Name("user_email").Value(c.FormEmail).AttrIf(!c.FieldStatusEmail, "disabled", "disabled"))).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Business Name")).
			Child(hb.Input().Class("form-control").Name("user_business_name").Value(c.FormBusiness).AttrIf(!c.FieldStatusBusinessName, "disabled", "disabled"))).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Phone")).
			Child(hb.Input().Class("form-control").Name("user_phone").Value(c.FormPhone).AttrIf(!c.FieldStatusPhone, "disabled", "disabled"))).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Country")).
			Child(countrySelect)).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Timezone")).
			Child(timezoneSelect)).
		Child(hb.Div().Class("mb-3").
			Child(hb.Label().Class("form-label").Text("Admin Notes")).
			Child(hb.Textarea().Class("form-control").Name("user_memo").Text(c.FormMemo).AttrIf(!c.FieldStatusFirstName || !c.FieldStatusLastName || !c.FieldStatusEmail || !c.FieldStatusBusinessName || !c.FieldStatusPhone, "disabled", "disabled"))).
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

func (c *formUserUpdate) refreshTimezones(ctx context.Context) {
	if c.registry == nil || c.registry.GetGeoStore() == nil {
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

	timezones, err := c.registry.GetGeoStore().TimezoneList(ctx, query)
	if err != nil {
		if c.registry.GetLogger() != nil {
			c.registry.GetLogger().Error("Error listing timezones", "error", err.Error())
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
