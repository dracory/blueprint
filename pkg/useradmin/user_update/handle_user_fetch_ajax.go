package user_update

import (
	"log/slog"
	"net/http"

	"project/internal/ext"

	"github.com/dracory/api"
	"github.com/dracory/geostore"
	"github.com/dracory/neat"
	"github.com/dracory/req"
)

func (controller *userUpdateController) handleUserFetchAjax(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return
	}
	userID := req.GetStringTrimmed(r, "user_id")
	if userID == "" {
		api.Respond(w, r, api.Error("User ID is required"))
		return
	}

	user, err := controller.app.GetUserStore().UserFindByID(r.Context(), userID)
	if err != nil {
		if controller.app.GetLogger() != nil {
			controller.app.GetLogger().Error("handleUserFetchAjax UserFindByID", slog.String("user_id", userID), slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Error loading user"))
		return
	}
	if user == nil {
		if controller.app.GetLogger() != nil {
			controller.app.GetLogger().Error("handleUserFetchAjax user not found", slog.String("user_id", userID))
		}
		api.Respond(w, r, api.Error("User not found"))
		return
	}

	firstName := user.GetFirstName()
	lastName := user.GetLastName()
	email := user.GetEmail()
	phone := user.GetPhone()
	business := user.GetBusinessName()
	memo := user.GetMemo()
	status := user.GetStatus()
	role := user.GetRole()
	country := user.GetCountry()
	timezone := user.GetTimezone()

	fieldStatus := map[string]bool{
		"first_name":    true,
		"last_name":     true,
		"email":         true,
		"business_name": true,
		"phone":         true,
		"role":          true,
	}

	if controller.app.GetConfig().GetVaultStoreUsed() && controller.app.GetVaultStore() != nil {
		firstName, lastName, email, phone, business, err = ext.UserUntokenize(r.Context(), controller.app, controller.app.GetConfig().GetVaultStoreKey(), user)
		if err != nil {
			if controller.app.GetLogger() != nil {
				controller.app.GetLogger().Error("userUpdateController.handleUserFetchAjax UserUntokenize", slog.String("error", err.Error()))
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

	if controller.app.GetGeoStore() == nil {
		if controller.app.GetLogger() != nil {
			controller.app.GetLogger().Error("userUpdateController.handleUserFetchAjax GeoStore not configured")
		}
		api.Respond(w, r, api.Error("GeoStore is not configured"))
		return
	}

	countryList, err := controller.app.GetGeoStore().CountryList(r.Context(), geostore.CountryQueryOptions{
		SortOrder: neat.SortAsc,
		OrderBy:   geostore.COLUMN_NAME,
	})
	if err != nil {
		if controller.app.GetLogger() != nil {
			controller.app.GetLogger().Error("userUpdateController.handleUserFetchAjax CountryList", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to load countries"))
		return
	}
	countries := make([]map[string]string, 0, len(countryList))
	for _, c := range countryList {
		countries = append(countries, map[string]string{
			FieldIsoCode2: c.IsoCode2(),
			FieldName:     c.Name(),
		})
	}

	timezoneList, err := controller.app.GetGeoStore().TimezoneList(r.Context(), geostore.TimezoneQueryOptions{
		SortOrder:   neat.SortAsc,
		OrderBy:     geostore.COLUMN_TIMEZONE,
		CountryCode: country,
	})
	if err != nil {
		if controller.app.GetLogger() != nil {
			controller.app.GetLogger().Error("userUpdateController.handleUserFetchAjax TimezoneList", slog.String("error", err.Error()))
		}
		api.Respond(w, r, api.Error("Failed to load timezones"))
		return
	}
	timezones := make([]map[string]string, 0, len(timezoneList))
	for _, tz := range timezoneList {
		timezones = append(timezones, map[string]string{
			FieldTimezone: tz.Timezone(),
		})
	}

	api.Respond(w, r, api.SuccessWithData("", map[string]any{
		FieldStatus:       status,
		FieldRole:         role,
		FieldFirstName:    firstName,
		FieldLastName:     lastName,
		FieldEmail:        email,
		FieldBusinessName: business,
		FieldPhone:        phone,
		FieldCountry:      country,
		FieldTimezone:     timezone,
		FieldMemo:         memo,
		FieldStatusField:  fieldStatus,
		FieldCountries:    countries,
		FieldTimezones:    timezones,
	}))
}
