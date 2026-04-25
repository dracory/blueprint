package user_manager

import (
	_ "embed"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"project/internal/ext"
	"project/internal/layouts"
	"project/internal/links"
	"project/internal/registry"
	"project/pkg/useradmin/shared"

	"github.com/dracory/blindindexstore"
	"github.com/dracory/cdn"
	"github.com/dracory/hb"
	"github.com/dracory/req"
	"github.com/dracory/sb"
	"github.com/dracory/userstore"
)

var (
	//go:embed users.html
	usersHTML string

	//go:embed users.js
	usersJS string
)

type userManagerController struct{ registry registry.RegistryInterface }

type JSONResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewUserManagerController(registry registry.RegistryInterface) *userManagerController {
	return &userManagerController{registry: registry}
}

func (controller *userManagerController) Handler(w http.ResponseWriter, r *http.Request) string {
	action := req.GetStringTrimmed(r, "action")

	switch action {
	case actionLoadUsers:
		return controller.handleLoadUsers(w, r)
	case actionDeleteUser:
		return controller.handleDeleteUser(w, r)
	case actionCreateUser:
		return controller.handleCreateUser(w, r)
	default:
		return controller.renderPage(w, r)
	}
}

func (controller *userManagerController) renderPage(w http.ResponseWriter, r *http.Request) string {
	if controller.registry == nil {
		http.Error(w, "Registry not initialized", http.StatusInternalServerError)
		return ""
	}

	urlUsersLoad := shared.NewLinks("/admin/users").UserManager(map[string]string{"action": actionLoadUsers})
	urlUserDelete := shared.NewLinks("/admin/users").UserManager(map[string]string{"action": actionDeleteUser})
	urlUserCreate := shared.NewLinks("/admin/users").UserManager(map[string]string{"action": actionCreateUser})
	urlUserUpdate := shared.NewLinks("/admin/users").UserUpdate(map[string]string{"user_id": "USER_ID_PLACEHOLDER"})
	urlUserImpersonate := shared.NewLinks("/admin/users").UserImpersonate(map[string]string{"user_id": "USER_ID_PLACEHOLDER"})

	html := strings.ReplaceAll(usersHTML, "urlUsersLoad", "'"+urlUsersLoad+"'")
	html = strings.ReplaceAll(html, "urlUserUpdate", "'"+urlUserUpdate+"'")
	html = strings.ReplaceAll(html, "urlUserImpersonate", "'"+urlUserImpersonate+"'")
	js := strings.ReplaceAll(usersJS, "urlUsersLoad", "'"+urlUsersLoad+"'")
	js = strings.ReplaceAll(js, "urlUserDelete", "'"+urlUserDelete+"'")
	js = strings.ReplaceAll(js, "urlUserCreate", "'"+urlUserCreate+"'")
	js = strings.ReplaceAll(js, "urlUserUpdate", "'"+urlUserUpdate+"'")
	js = strings.ReplaceAll(js, "urlUserImpersonate", "'"+urlUserImpersonate+"'")

	breadcrumbs := layouts.Breadcrumbs([]layouts.Breadcrumb{
		{Name: "Home", URL: links.Admin().Home(map[string]string{})},
		{Name: "User Manager", URL: shared.NewLinks("/admin/users").UserManager()},
	})

	vueCDN := hb.Script("").Src("https://unpkg.com/vue@3/dist/vue.global.js")

	content := hb.Div().
		Child(vueCDN).
		Child(hb.Raw(html)).
		Child(hb.Script(js))

	return layouts.NewAdminLayout(controller.registry, r, layouts.Options{
		Title:   "Users | User Manager",
		Content: layouts.AdminPage(breadcrumbs, content),
		ScriptURLs: []string{
			cdn.Sweetalert2_10(),
		},
		Styles: []string{},
	}).ToHTML()
}

func (controller *userManagerController) handleLoadUsers(w http.ResponseWriter, r *http.Request) string {
	if controller.registry.GetUserStore() == nil {
		return controller.jsonResponse(w, false, "User store not configured", nil)
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
		return controller.jsonResponse(w, false, "Invalid request body", nil)
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
	sortOrder := getStringTrimmed(reqBody.SortOrder, sb.DESC)
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
		ids, err := controller.registry.GetBlindIndexStoreFirstName().Search(r.Context(), firstName, blindindexstore.SEARCH_TYPE_CONTAINS)
		if err != nil {
			controller.registry.GetLogger().Error("userManagerController.handleLoadUsers blind index first_name", slog.String("error", err.Error()))
		}
		if len(ids) == 0 {
			return controller.jsonResponse(w, true, "", map[string]interface{}{"users": []interface{}{}, "total": 0})
		}
		userIDs = append(userIDs, ids...)
	}

	if lastName != "" {
		ids, err := controller.registry.GetBlindIndexStoreLastName().Search(r.Context(), lastName, blindindexstore.SEARCH_TYPE_CONTAINS)
		if err != nil {
			controller.registry.GetLogger().Error("userManagerController.handleLoadUsers blind index last_name", slog.String("error", err.Error()))
		}
		if len(ids) == 0 {
			return controller.jsonResponse(w, true, "", map[string]interface{}{"users": []interface{}{}, "total": 0})
		}
		userIDs = append(userIDs, ids...)
	}

	if email != "" {
		ids, err := controller.registry.GetBlindIndexStoreEmail().Search(r.Context(), email, blindindexstore.SEARCH_TYPE_CONTAINS)
		if err != nil {
			controller.registry.GetLogger().Error("userManagerController.handleLoadUsers blind index email", slog.String("error", err.Error()))
		}
		if len(ids) == 0 {
			return controller.jsonResponse(w, true, "", map[string]interface{}{"users": []interface{}{}, "total": 0})
		}
		userIDs = append(userIDs, ids...)
	}

	if len(userIDs) > 0 {
		query.SetIDIn(userIDs)
	}

	userList, err := controller.registry.GetUserStore().UserList(r.Context(), query)
	if err != nil {
		controller.registry.GetLogger().Error("userManagerController.handleLoadUsers UserList", slog.String("error", err.Error()))
		return controller.jsonResponse(w, false, "Failed to load users", nil)
	}

	userCount, err := controller.registry.GetUserStore().UserCount(r.Context(), query)
	if err != nil {
		controller.registry.GetLogger().Error("userManagerController.handleLoadUsers UserCount", slog.String("error", err.Error()))
		return controller.jsonResponse(w, false, "Failed to count users", nil)
	}

	users := make([]map[string]interface{}, 0, len(userList))
	for _, user := range userList {
		firstNameVal := user.GetFirstName()
		lastNameVal := user.GetLastName()
		emailVal := user.GetEmail()

		if controller.registry.GetConfig().GetVaultStoreUsed() && controller.registry.GetVaultStore() != nil {
			var err error
			firstNameVal, lastNameVal, emailVal, _, _, err = ext.UserUntokenize(r.Context(), controller.registry, controller.registry.GetConfig().GetVaultStoreKey(), user)
			if err != nil {
				controller.registry.GetLogger().Error("userManagerController.handleLoadUsers UserUntokenize", slog.String("error", err.Error()))
				firstNameVal = "n/a"
				lastNameVal = "n/a"
				emailVal = "n/a"
			}
		}

		users = append(users, map[string]interface{}{
			"id":         user.GetID(),
			"first_name": firstNameVal,
			"last_name":  lastNameVal,
			"email":      emailVal,
			"status":     user.GetStatus(),
			"created_at": user.GetCreatedAtCarbon().Format("d M Y"),
			"updated_at": user.GetUpdatedAtCarbon().Format("d M Y"),
		})
	}

	return controller.jsonResponse(w, true, "", map[string]interface{}{
		"users": users,
		"total": userCount,
	})
}

func (controller *userManagerController) handleDeleteUser(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		return controller.jsonResponse(w, false, "Method not allowed", nil)
	}

	if controller.registry.GetUserStore() == nil {
		return controller.jsonResponse(w, false, "User store not configured", nil)
	}

	var reqBody struct {
		UserID string `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return controller.jsonResponse(w, false, "Invalid request body", nil)
	}

	if reqBody.UserID == "" {
		return controller.jsonResponse(w, false, "User ID is required", nil)
	}

	user, err := controller.registry.GetUserStore().UserFindByID(r.Context(), reqBody.UserID)
	if err != nil {
		controller.registry.GetLogger().Error("userManagerController.handleDeleteUser", slog.String("error", err.Error()))
		return controller.jsonResponse(w, false, "User not found", nil)
	}
	if user == nil {
		return controller.jsonResponse(w, false, "User not found", nil)
	}

	if err := controller.registry.GetUserStore().UserSoftDelete(r.Context(), user); err != nil {
		controller.registry.GetLogger().Error("userManagerController.handleDeleteUser", slog.String("error", err.Error()))
		return controller.jsonResponse(w, false, "Failed to delete user", nil)
	}

	return controller.jsonResponse(w, true, "User deleted successfully", nil)
}

func (controller *userManagerController) handleCreateUser(w http.ResponseWriter, r *http.Request) string {
	if r.Method != http.MethodPost {
		return controller.jsonResponse(w, false, "Method not allowed", nil)
	}

	if controller.registry.GetUserStore() == nil {
		return controller.jsonResponse(w, false, "User store not configured", nil)
	}

	var reqBody struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		return controller.jsonResponse(w, false, "Invalid request body", nil)
	}

	if strings.TrimSpace(reqBody.FirstName) == "" {
		return controller.jsonResponse(w, false, "First name is required", nil)
	}
	if strings.TrimSpace(reqBody.LastName) == "" {
		return controller.jsonResponse(w, false, "Last name is required", nil)
	}
	if strings.TrimSpace(reqBody.Email) == "" {
		return controller.jsonResponse(w, false, "Email is required", nil)
	}

	user := userstore.NewUser()
	user.SetFirstName(strings.TrimSpace(reqBody.FirstName))
	user.SetLastName(strings.TrimSpace(reqBody.LastName))
	user.SetEmail(strings.TrimSpace(reqBody.Email))

	if err := controller.registry.GetUserStore().UserCreate(r.Context(), user); err != nil {
		controller.registry.GetLogger().Error("userManagerController.handleCreateUser", slog.String("error", err.Error()))
		return controller.jsonResponse(w, false, "Failed to create user", nil)
	}

	return controller.jsonResponse(w, true, "User created successfully", map[string]interface{}{"user_id": user.GetID()})
}

func (controller *userManagerController) jsonResponse(w http.ResponseWriter, success bool, message string, data interface{}) string {
	w.Header().Set("Content-Type", "application/json")
	resp := JSONResponse{Success: success, Message: message, Data: data}
	json.NewEncoder(w).Encode(resp)
	return ""
}

const (
	actionLoadUsers  = "load-users"
	actionDeleteUser = "delete-user"
	actionCreateUser = "create-user"
)
