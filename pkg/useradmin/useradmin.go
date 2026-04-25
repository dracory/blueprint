// Package useradmin provides a user admin interface following the folder-per-controller pattern.
// Each controller is in its own subfolder and handles its own views and AJAX data.
// This structure allows for future migration to a standalone external package.
package useradmin

import (
	"net/http"

	"project/internal/registry"
	"project/pkg/useradmin/shared"
	"project/pkg/useradmin/user_create"
	"project/pkg/useradmin/user_delete"
	"project/pkg/useradmin/user_impersonate"
	"project/pkg/useradmin/user_manager"
	"project/pkg/useradmin/user_update"

	"github.com/dracory/req"
)

// AdminOptions contains all dependencies and configuration for the user admin
type AdminOptions struct {
	// Registry provides access to all stores and services
	Registry registry.RegistryInterface

	// AdminHomeURL is the URL for the admin home page
	AdminHomeURL string

	// UserAdminURL is the base URL for the user admin (e.g., "/admin/users")
	UserAdminURL string

	// AuthUserID returns the authenticated user ID from the request
	AuthUserID func(r *http.Request) string
}

// AdminInterface defines the interface for the user admin
type AdminInterface interface {
	Handle(w http.ResponseWriter, r *http.Request) string
}

// admin implements AdminInterface
type admin struct {
	opts AdminOptions
}

// New creates a new user admin instance
func New(opts AdminOptions) (AdminInterface, error) {
	if opts.Registry == nil {
		return nil, ErrRegistryRequired
	}

	// Set defaults
	if opts.UserAdminURL == "" {
		opts.UserAdminURL = "/admin/users"
	}

	return &admin{opts: opts}, nil
}

// Handle processes all user admin requests
func (a *admin) Handle(w http.ResponseWriter, r *http.Request) string {
	// Check authentication
	if a.opts.AuthUserID != nil && a.opts.AuthUserID(r) == "" {
		http.Redirect(w, r, a.opts.AdminHomeURL, http.StatusSeeOther)
		return ""
	}

	controller := req.GetStringTrimmed(r, "controller")

	switch controller {
	case shared.CONTROLLER_USER_MANAGER:
		return user_manager.NewUserManagerController(a.opts.Registry).Handler(w, r)
	case shared.CONTROLLER_USER_CREATE:
		return user_create.NewUserCreateController(a.opts.Registry).Handler(w, r)
	case shared.CONTROLLER_USER_DELETE:
		return user_delete.NewUserDeleteController(a.opts.Registry).Handler(w, r)
	case shared.CONTROLLER_USER_UPDATE:
		return user_update.NewUserUpdateController(a.opts.Registry).Handler(w, r)
	case shared.CONTROLLER_USER_IMPERSONATE:
		return user_impersonate.NewUserImpersonateController(a.opts.Registry).Handler(w, r)
	}

	// Default to user manager
	return user_manager.NewUserManagerController(a.opts.Registry).Handler(w, r)
}
