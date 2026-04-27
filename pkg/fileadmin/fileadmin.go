// Package fileadmin provides a file admin interface following the folder-per-controller pattern.
// Each controller is in its own subfolder and handles its own views and AJAX data.
// This structure allows for future migration to a standalone external package.
package fileadmin

import (
	"net/http"

	"project/internal/registry"
	"project/pkg/fileadmin/file_manager"
)

// AdminOptions contains all dependencies and configuration for the file admin
type AdminOptions struct {
	// Registry provides access to all stores and services
	Registry registry.RegistryInterface

	// AdminHomeURL is the URL for the admin home page
	AdminHomeURL string

	// FileAdminURL is the base URL for the file admin (e.g., "/admin/file-manager")
	FileAdminURL string

	// AuthUserID returns the authenticated user ID from the request
	AuthUserID func(r *http.Request) string
}

// AdminInterface defines the interface for the file admin
type AdminInterface interface {
	Handle(w http.ResponseWriter, r *http.Request) string
}

// admin implements AdminInterface
type admin struct {
	opts AdminOptions
}

// New creates a new file admin instance
func New(opts AdminOptions) (AdminInterface, error) {
	if opts.Registry == nil {
		return nil, ErrRegistryRequired
	}

	// Set defaults
	if opts.FileAdminURL == "" {
		opts.FileAdminURL = "/admin/file-manager"
	}

	return &admin{opts: opts}, nil
}

// Handle processes all file admin requests
func (a *admin) Handle(w http.ResponseWriter, r *http.Request) string {
	// Check authentication
	if a.opts.AuthUserID != nil && a.opts.AuthUserID(r) == "" {
		http.Redirect(w, r, a.opts.AdminHomeURL, http.StatusSeeOther)
		return ""
	}

	if a.opts.Registry == nil {
		http.Error(w, "Registry not configured", http.StatusInternalServerError)
		return ""
	}

	// Delegate to the file manager controller
	controller := file_manager.NewFileManagerController(a.opts.Registry)
	return controller.Handler(w, r)
}
