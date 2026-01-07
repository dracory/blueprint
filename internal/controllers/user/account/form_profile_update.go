package account

import (
	"context"
	"log"
	"net/url"
	"strings"

	"project/internal/ext"
	"project/internal/links"
	"project/internal/registry"
	"project/internal/types"

	"github.com/dracory/geostore"
	"github.com/dracory/hb"
	"github.com/dracory/liveflux"
	"github.com/dracory/sb"
	"github.com/samber/lo"
)

// == COMPONENT ===============================================================

type formProfileUpdate struct {
	liveflux.Base
	App                    types.RegistryInterface
	UserID                 string
	ReturnURL              string
	FormEmail              string
	FormFirstName          string
	FormLastName           string
	FormBusinessName       string
	FormPhone              string
	FormCountry            string
	FormTimezone           string
	FormError              string
	FormSuccess            string
	FormRedirectTo         string
	Countries              []geostore.Country
	Timezones              []geostore.Timezone
	CountrySelectActionURL string
}

// == CONSTRUCTOR =============================================================

func NewFormProfileUpdate(registry registry.RegistryInterface) liveflux.ComponentInterface {
	inst, err := liveflux.New(&formProfileUpdate{})
	if err != nil {
		log.Println(err)
		return nil
	}
	if c, ok := inst.(*formProfileUpdate); ok {
		c.App = registry
	}
	return inst
}

// == PUBLIC METHODS ==========================================================

func (c *formProfileUpdate) GetKind() string {
	return "profile_update_form"
}

func (c *formProfileUpdate) Mount(ctx context.Context, params map[string]string) error {
	c.UserID = strings.TrimSpace(params["user_id"])
	if c.UserID == "" {
		c.FormError = "User ID is required"
		return nil
	}

	c.ReturnURL = strings.TrimSpace(params["return_url"])
	if c.ReturnURL == "" {
		c.ReturnURL = links.User().Profile()
	}

	if c.App.GetUserStore() == nil {
		c.App.GetLogger().Error("User store not initialized")
		c.FormError = "Error getting user"
		return nil
	}

	user, err := c.App.GetUserStore().UserFindByID(ctx, c.UserID)
	if err != nil {
		c.App.GetLogger().Error("Error getting user", "error", err.Error())
		c.FormError = "Error getting user"
		return nil
	}

	if user == nil {
		c.App.GetLogger().Error("Error getting user", "error", "user not found")
		c.FormError = "Error getting user"
		return nil
	}

	if c.App.GetGeoStore() == nil {
		c.App.GetLogger().Error("Geo store not initialized")
		c.FormError = "Error listing countries"
		return nil
	}

	countryList, err := c.App.GetGeoStore().CountryList(ctx, geostore.CountryQueryOptions{
		SortOrder: sb.ASC,
		OrderBy:   geostore.COLUMN_NAME,
	})
	if err != nil {
		c.App.GetLogger().Error("Error listing countries", "error", err.Error())
		c.FormError = "Error listing countries"
		return nil
	}
	c.Countries = countryList

	if c.App.GetConfig().GetUserStoreVaultEnabled() && c.App.GetVaultStore() != nil {
		firstName, lastName, email, businessName, phone, err := ext.UserUntokenize(
			ctx,
			c.App,
			c.App.GetConfig().GetVaultStoreKey(),
			user,
		)

		if err != nil {
			c.App.GetLogger().Error("Error reading profile data", "error", err.Error())
			c.FormError = "Error reading profile data"
			return nil
		}

		c.FormEmail = email
		c.FormFirstName = firstName
		c.FormLastName = lastName
		c.FormBusinessName = businessName
		c.FormPhone = phone
	} else {
		c.FormEmail = user.Email()
		c.FormFirstName = user.FirstName()
		c.FormLastName = user.LastName()
		c.FormBusinessName = user.BusinessName()
		c.FormPhone = user.Phone()
	}
	c.FormCountry = user.Country()
	c.FormTimezone = user.Timezone()
	c.refreshTimezones(ctx)

	return nil
}

func (c *formProfileUpdate) Handle(ctx context.Context, action string, data url.Values) error {
	if data == nil {
		data = url.Values{}
	}
	switch action {
	case "apply", "save":
		return c.handleUpdate(ctx, action, data)
	case "country_change":
		c.FormCountry = strings.TrimSpace(data.Get("country"))
		c.FormTimezone = strings.TrimSpace(data.Get("timezone"))
		c.refreshTimezones(ctx)
		return nil
	default:
		return nil
	}
}

// == PRIVATE METHODS =========================================================

func (c *formProfileUpdate) handleUpdate(ctx context.Context, action string, data url.Values) error {
	if data == nil {
		data = url.Values{}
	}
	userID := strings.TrimSpace(data.Get("user_id"))
	if userID == "" {
		c.FormError = "User ID is required"
		c.FormSuccess = ""
		return nil
	}

	user, err := c.App.GetUserStore().UserFindByID(ctx, userID)
	if err != nil {
		c.App.GetLogger().Error("Error getting user", "error", err.Error())
		c.FormError = "Error getting user"
		c.FormSuccess = ""
		return nil
	}

	if user == nil {
		c.App.GetLogger().Error("Error getting user", "error", "user not found")
		c.FormError = "Error getting user"
		c.FormSuccess = ""
		return nil
	}

	c.FormEmail = strings.TrimSpace(data.Get("email"))
	c.FormFirstName = strings.TrimSpace(data.Get("first_name"))
	c.FormLastName = strings.TrimSpace(data.Get("last_name"))
	c.FormBusinessName = strings.TrimSpace(data.Get("business_name"))
	c.FormPhone = strings.TrimSpace(data.Get("phone"))
	c.FormCountry = strings.TrimSpace(data.Get("country"))
	c.FormTimezone = strings.TrimSpace(data.Get("timezone"))

	if c.FormFirstName == "" {
		c.FormError = "First name is required field"
		c.FormSuccess = ""
		return nil
	}

	if c.FormLastName == "" {
		c.FormError = "Last name is required field"
		c.FormSuccess = ""
		return nil
	}

	if c.FormEmail == "" {
		c.FormError = "Email is required field"
		c.FormSuccess = ""
		return nil
	}

	if c.FormCountry == "" {
		c.FormError = "Country is required field"
		c.FormSuccess = ""
		return nil
	}

	if c.FormTimezone == "" {
		c.FormError = "Timezone is required field"
		c.FormSuccess = ""
		return nil
	}

	if c.App.GetConfig().GetUserStoreVaultEnabled() && c.App.GetVaultStore() == nil {
		c.FormError = "We are very sorry vault store is not configured. Saving the details not possible."
		c.FormSuccess = ""
		return nil
	}

	if !c.App.GetConfig().GetUserStoreVaultEnabled() {
		user.SetFirstName(c.FormFirstName)
		user.SetLastName(c.FormLastName)
		user.SetBusinessName(c.FormBusinessName)
		user.SetPhone(c.FormPhone)
	} else {
		if err := c.App.GetVaultStore().TokenUpdate(ctx, user.FirstName(), c.FormFirstName, c.App.GetConfig().GetVaultStoreKey()); err != nil {
			c.App.GetLogger().Error("Error saving first name", "error", err.Error())
			c.FormError = "Saving profile failed. Please try again later."
			c.FormSuccess = ""
			return nil
		}
		if err := c.App.GetVaultStore().TokenUpdate(ctx, user.LastName(), c.FormLastName, c.App.GetConfig().GetVaultStoreKey()); err != nil {
			c.App.GetLogger().Error("Error saving last name", "error", err.Error())
			c.FormError = "Saving profile failed. Please try again later."
			c.FormSuccess = ""
			return nil
		}
		if err := c.App.GetVaultStore().TokenUpdate(ctx, user.BusinessName(), c.FormBusinessName, c.App.GetConfig().GetVaultStoreKey()); err != nil {
			c.App.GetLogger().Error("Error saving business name", "error", err.Error())
			c.FormError = "Saving profile failed. Please try again later."
			c.FormSuccess = ""
			return nil
		}
		if err := c.App.GetVaultStore().TokenUpdate(ctx, user.Phone(), c.FormPhone, c.App.GetConfig().GetVaultStoreKey()); err != nil {
			c.App.GetLogger().Error("Error saving phone", "error", err.Error())
			c.FormError = "Saving profile failed. Please try again later."
			c.FormSuccess = ""
			return nil
		}
	}

	user.SetCountry(c.FormCountry)
	user.SetTimezone(c.FormTimezone)

	if c.App.GetUserStore() == nil {
		c.App.GetLogger().Warn("At formProfileUpdate > handleUpdate. UserStore is nil.")
		c.FormError = "Saving profile failed. Please try again later."
		c.FormSuccess = ""
		return nil
	}

	if err := c.App.GetUserStore().UserUpdate(context.Background(), user); err != nil {
		c.App.GetLogger().Error("Error updating user profile", "error", err.Error())
		c.FormError = "Saving profile failed. Please try again later."
		c.FormSuccess = ""
		return nil
	}

	c.FormError = ""

	if action == "save" {
		c.FormSuccess = "Profile updated successfully"
		c.FormRedirectTo = links.User().Home()
	} else {
		c.FormSuccess = "Profile updated successfully"
		c.FormRedirectTo = ""
	}

	return nil
}

func (c *formProfileUpdate) refreshTimezones(ctx context.Context) {
	query := geostore.TimezoneQueryOptions{
		SortOrder: sb.ASC,
		OrderBy:   geostore.COLUMN_TIMEZONE,
	}

	if c.FormCountry != "" {
		query.CountryCode = c.FormCountry
	}

	timezones, err := c.App.GetGeoStore().TimezoneList(ctx, query)
	if err != nil {
		c.App.GetLogger().Error("Error listing timezones", "error", err.Error())
		c.FormError = "Error listing timezones"
		return
	}

	c.Timezones = timezones
}

func (c *formProfileUpdate) timezoneOptions() []hb.TagInterface {
	return lo.Map(c.Timezones, func(tz geostore.Timezone, _ int) hb.TagInterface {
		return hb.Option().Text(tz.Timezone()).Value(tz.Timezone()).
			AttrIf(c.FormTimezone == tz.Timezone(), "selected", "selected")
	})
}

func (c *formProfileUpdate) timezoneSelect() hb.TagInterface {
	selectTag := hb.Select().
		ID("SelectTimezones").
		Class("form-select").
		Name("timezone")
	selectTag.Child(hb.Option().Text(""))
	selectTag.Children(c.timezoneOptions())
	return selectTag
}

func (c *formProfileUpdate) Render(ctx context.Context) hb.TagInterface {
	required := hb.Sup().
		Text("required").
		Style("margin-left:5px;color:lightcoral;")

	emailGroup := hb.Div().
		Class("mb-3 form-group").
		Child(hb.Label().
			Class("form-label").
			Text("Email").
			Child(required)).
		Child(hb.Input().
			Class("form-control").
			Name("email").
			Value(c.FormEmail).
			Attr("readonly", "readonly"))

	firstNameGroup := hb.Div().
		Class("mb-3 form-group").
		Child(hb.Label().
			Class("form-label").
			Text("First name").
			Child(required)).
		Child(hb.Input().
			Class("form-control").
			Name("first_name").
			Value(c.FormFirstName))

	lastNameGroup := hb.Div().
		Class("mb-3 form-group").
		Child(hb.Label().
			Class("form-label").
			Text("Last name").
			Child(required)).
		Child(hb.Input().
			Class("form-control").
			Name("last_name").
			Value(c.FormLastName))

	businessNameGroup := hb.Div().
		Class("mb-3 form-group").
		Child(hb.Label().
			Class("form-label").
			Text("Company / business name")).
		Child(hb.Input().
			Class("form-control").
			Name("business_name").
			Value(c.FormBusinessName))

	phoneGroup := hb.Div().
		Class("mb-3 form-group").
		Child(hb.Label().
			Class("form-label").
			Text("Phone")).
		Child(hb.Input().
			Class("form-control").
			Name("phone").
			Value(c.FormPhone))

	countrySelect := hb.Select().
		ID("SelectCountries").
		Class("form-select").
		Name("country").
		Attr("onchange", `
			const btn = document.getElementById('btnCountryChange');
			if (btn) btn.click();
		`)
	countrySelect.Child(hb.Option().Text("").Value(""))
	for _, country := range c.Countries {
		countrySelect.Child(hb.Option().Text(country.Name()).Value(country.IsoCode2()).
			AttrIf(c.FormCountry == country.IsoCode2(), "selected", "selected"))
	}

	timezoneSelect := c.timezoneSelect()

	countryGroup := hb.Div().
		Class("mb-3 form-group").
		Children([]hb.TagInterface{
			hb.Label().
				Text("Country").
				Class("form-label").
				Child(required),
			countrySelect,
		})

	timezoneGroup := hb.Div().
		Class("mb-3").
		Children([]hb.TagInterface{
			hb.Label().
				Text("Timezone").
				Class("form-label").
				Child(required),
			timezoneSelect,
		})

	buttonApply := hb.Button().
		Type("submit").
		Class("btn btn-primary").
		Attr("data-flux-action", "apply").
		Child(hb.I().Class("bi bi-check2 me-2")).
		Text("Apply")

	buttonSaveAndClose := hb.Button().
		Type("submit").
		Class("btn btn-success").
		Attr("data-flux-action", "save").
		Child(hb.I().Class("bi bi-check2-all me-2")).
		Text("Save & Close")

	buildActions := func(marginClass string) hb.TagInterface {
		return hb.Div().Class(marginClass + " d-flex justify-content-between align-items-center").
			Child(hb.A().
				Href(c.ReturnURL).
				Class("btn btn-secondary").
				Child(hb.I().Class("bi bi-chevron-left me-2")).
				Text("Cancel")).
			Child(hb.Div().Class("d-flex gap-2").
				Child(buttonApply).
				Child(buttonSaveAndClose))
	}

	alerts := hb.Div()
	if c.FormError != "" {
		swal := hb.SwalError(hb.SwalOptions{
			Text:             c.FormError,
			Timer:            5000,
			TimerProgressBar: true,
			Position:         "top-end",
		})
		alerts = alerts.Child(swal)
	}
	if c.FormSuccess != "" && c.FormRedirectTo != "" {
		swal := hb.SwalSuccess(hb.SwalOptions{
			Text:             c.FormSuccess,
			RedirectURL:      c.FormRedirectTo,
			RedirectSeconds:  5,
			Timer:            5000,
			TimerProgressBar: true,
			Position:         "top-end",
		})
		alerts = alerts.Child(swal)
	} else if c.FormSuccess != "" {
		swal := hb.SwalSuccess(hb.SwalOptions{
			Text:             c.FormSuccess,
			Timer:            5000,
			TimerProgressBar: true,
			Position:         "top-end",
		})
		alerts = alerts.Child(swal)
	}

	buttonCountryChangeHidden := hb.Button().
		ID("btnCountryChange").
		Type("submit").
		Attr("data-flux-action", "country_change").
		Style("display:none;")

	userIdHidden := hb.Input().
		Type("hidden").
		Name("user_id").
		Value(c.UserID)

	form := hb.Form().
		Child(alerts).
		Child(userIdHidden).
		Child(buttonCountryChangeHidden).
		Child(buildActions("mb-3")).
		Child(emailGroup).
		Child(firstNameGroup).
		Child(lastNameGroup).
		Child(businessNameGroup).
		Child(phoneGroup).
		Child(countryGroup).
		Child(timezoneGroup).
		Child(buildActions("mt-4"))

	card := hb.Div().
		Class("card").
		Child(hb.Div().
			Class("card-header").
			Child(hb.H4().
				Class("card-title").
				Text("Your Details"))).
		Child(hb.Div().
			Class("card-body").
			Child(form))

	return c.Root(card)
}

func init() {
	if err := liveflux.Register(&formProfileUpdate{}); err != nil {
		log.Printf("Failed to register formProfileUpdate component: %v", err)
	}
}
