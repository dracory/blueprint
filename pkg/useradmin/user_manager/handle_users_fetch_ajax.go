package user_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"project/internal/ext"

	"github.com/dracory/api"
	"github.com/dracory/blindindexstore"
	"github.com/dracory/neat"
	"github.com/dracory/userstore"
)

func (controller *userManagerController) handleUsersFetchAjax(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}
	if controller.app.GetUserStore() == nil {
		api.Respond(w, r, api.Error("User store not configured"))
		return ""
	}

	// Parse request body
	var reqBody struct {
		Page        int    `json:"page"`
		PerPage     int    `json:"per_page"`
		SortOrder   string `json:"sort_order"`
		SortBy      string `json:"sort_by"`
		Status      string `json:"status"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		Email       string `json:"email"`
		UserID      string `json:"user_id"`
		CreatedFrom string `json:"created_from"`
		CreatedTo   string `json:"created_to"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return ""
	}

	// Helper functions for trimmed values with defaults (similar to req.GetStringTrimmedOr)
	getInt := func(val int, defaultVal int) int {
		if val == 0 {
			return defaultVal
		}
		return val
	}

	getStringTrimmed := func(val string, defaultVal string) string {
		val = strings.TrimSpace(val)
		if val == "" {
			return defaultVal
		}
		return val
	}

	page := getInt(reqBody.Page, 0)
	perPage := getInt(reqBody.PerPage, 10)
	sortOrder := getStringTrimmed(reqBody.SortOrder, neat.SortDesc)
	sortBy := getStringTrimmed(reqBody.SortBy, userstore.COLUMN_CREATED_AT)
	status := getStringTrimmed(reqBody.Status, "")
	firstName := getStringTrimmed(reqBody.FirstName, "")
	lastName := getStringTrimmed(reqBody.LastName, "")
	email := getStringTrimmed(reqBody.Email, "")
	userID := getStringTrimmed(reqBody.UserID, "")
	createdFrom := getStringTrimmed(reqBody.CreatedFrom, "")
	createdTo := getStringTrimmed(reqBody.CreatedTo, "")

	query := userstore.NewUserQuery().
		SetSortDirection(sortOrder).
		SetOrderBy(sortBy).
		SetOffset(page * perPage).
		SetLimit(perPage)

	if status != "" {
		query.SetStatus(status)
	}

	if userID != "" {
		query.SetID(userID)
	}

	if createdFrom != "" {
		query.SetCreatedAtGte(createdFrom + " 00:00:00")
	}

	if createdTo != "" {
		query.SetCreatedAtLte(createdTo + " 23:59:59")
	}

	userIDs := []string{}

	if firstName != "" {
		ids, err := controller.app.GetBlindIndexStoreFirstName().Search(r.Context(), firstName, blindindexstore.SEARCH_TYPE_CONTAINS)
		if err != nil {
			controller.app.GetLogger().Error("userManagerController.handleUsersFetchAjax blind index first_name", slog.String("error", err.Error()))
		}
		if len(ids) == 0 {
			api.Respond(w, r, api.SuccessWithData("", map[string]interface{}{FieldUsers: []interface{}{}, FieldTotal: 0}))
			return ""
		}
		userIDs = append(userIDs, ids...)
	}

	if lastName != "" {
		ids, err := controller.app.GetBlindIndexStoreLastName().Search(r.Context(), lastName, blindindexstore.SEARCH_TYPE_CONTAINS)
		if err != nil {
			controller.app.GetLogger().Error("userManagerController.handleUsersFetchAjax blind index last_name", slog.String("error", err.Error()))
		}
		if len(ids) == 0 {
			api.Respond(w, r, api.SuccessWithData("", map[string]interface{}{FieldUsers: []interface{}{}, FieldTotal: 0}))
			return ""
		}
		userIDs = append(userIDs, ids...)
	}

	if email != "" {
		ids, err := controller.app.GetBlindIndexStoreEmail().Search(r.Context(), email, blindindexstore.SEARCH_TYPE_CONTAINS)
		if err != nil {
			controller.app.GetLogger().Error("userManagerController.handleUsersFetchAjax blind index email", slog.String("error", err.Error()))
		}
		if len(ids) == 0 {
			api.Respond(w, r, api.SuccessWithData("", map[string]interface{}{FieldUsers: []interface{}{}, FieldTotal: 0}))
			return ""
		}
		userIDs = append(userIDs, ids...)
	}

	if len(userIDs) > 0 {
		query.SetIDIn(userIDs)
	}

	userList, err := controller.app.GetUserStore().UserList(r.Context(), query)
	if err != nil {
		controller.app.GetLogger().Error("userManagerController.handleUsersFetchAjax UserList", slog.String("error", err.Error()))
		api.Respond(w, r, api.Error("Failed to load users"))
		return ""
	}

	userCount, err := controller.app.GetUserStore().UserCount(r.Context(), query)
	if err != nil {
		controller.app.GetLogger().Error("userManagerController.handleUsersFetchAjax UserCount", slog.String("error", err.Error()))
		api.Respond(w, r, api.Error("Failed to count users"))
		return ""
	}

	users := make([]map[string]interface{}, 0, len(userList))
	for _, user := range userList {
		firstNameVal := user.GetFirstName()
		lastNameVal := user.GetLastName()
		emailVal := user.GetEmail()

		if controller.app.GetConfig().GetVaultStoreUsed() && controller.app.GetVaultStore() != nil {
			var err error
			firstNameVal, lastNameVal, emailVal, _, _, err = ext.UserUntokenize(r.Context(), controller.app, controller.app.GetConfig().GetVaultStoreKey(), user)
			if err != nil {
				controller.app.GetLogger().Error("userManagerController.handleUsersFetchAjax UserUntokenize", slog.String("error", err.Error()))
				firstNameVal = "n/a"
				lastNameVal = "n/a"
				emailVal = "n/a"
			}
		}

		users = append(users, map[string]interface{}{
			FieldID:        user.GetID(),
			FieldFirstName: firstNameVal,
			FieldLastName:  lastNameVal,
			FieldEmail:     emailVal,
			FieldStatus:    user.GetStatus(),
			FieldCreatedAt: user.GetCreatedAtCarbon().Format("d M Y"),
			FieldUpdatedAt: user.GetUpdatedAtCarbon().Format("d M Y"),
		})
	}

	api.Respond(w, r, api.SuccessWithData("", map[string]interface{}{
		FieldUsers: users,
		FieldTotal: userCount,
	}))
	return ""
}
