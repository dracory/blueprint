package user_manager

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/dracory/api"
	"github.com/dracory/userstore"
)

func (controller *userManagerController) handleUserCreateAjax(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		api.Respond(w, r, api.Error("Method not allowed"))
		return ""
	}

	if controller.app.GetUserStore() == nil {
		api.Respond(w, r, api.Error("User store not configured"))
		return ""
	}

	var reqBody struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		api.Respond(w, r, api.Error("Invalid request body"))
		return ""
	}

	if strings.TrimSpace(reqBody.FirstName) == "" {
		api.Respond(w, r, api.Error("First name is required"))
		return ""
	}
	if strings.TrimSpace(reqBody.LastName) == "" {
		api.Respond(w, r, api.Error("Last name is required"))
		return ""
	}
	if strings.TrimSpace(reqBody.Email) == "" {
		api.Respond(w, r, api.Error("Email is required"))
		return ""
	}

	user := userstore.NewUser()
	user.SetFirstName(strings.TrimSpace(reqBody.FirstName))
	user.SetLastName(strings.TrimSpace(reqBody.LastName))
	user.SetEmail(strings.TrimSpace(reqBody.Email))

	if err := controller.app.GetUserStore().UserCreate(r.Context(), user); err != nil {
		controller.app.GetLogger().Error("userManagerController.handleUserCreateAjax", slog.String("error", err.Error()))
		api.Respond(w, r, api.Error("Failed to create user"))
		return ""
	}

	api.Respond(w, r, api.SuccessWithData("User created successfully", map[string]interface{}{FieldUserID: user.GetID()}))
	return ""
}
