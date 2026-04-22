// Package logadmin provides a log admin interface following the folder-per-controller pattern.
// Each controller is in its own subfolder and handles its own views and AJAX data.
// This structure allows for future migration to a standalone external package.
package logadmin

import (
	"net/http"
	"strings"

	"project/internal/registry"
)

// AdminOptions contains all dependencies and configuration for the log admin
type AdminOptions struct {
	// Registry provides access to all stores and services
	Registry registry.RegistryInterface

	// AdminHomeURL is the URL for the admin home page
	AdminHomeURL string

	// LogAdminURL is the base URL for the log admin (e.g., "/admin/logs")
	LogAdminURL string

	// AuthUserID returns the authenticated user ID from the request
	AuthUserID func(r *http.Request) string
}

// AdminInterface defines the interface for the log admin
type AdminInterface interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

// admin implements AdminInterface
type admin struct {
	opts AdminOptions
}

// New creates a new log admin instance
func New(opts AdminOptions) (AdminInterface, error) {
	if opts.Registry == nil {
		return nil, ErrRegistryRequired
	}

	// Set defaults
	if opts.LogAdminURL == "" {
		opts.LogAdminURL = "/admin/logs"
	}

	return &admin{opts: opts}, nil
}

// Handle processes all log admin requests
func (a *admin) Handle(w http.ResponseWriter, r *http.Request) {
	// Check authentication
	if a.opts.AuthUserID != nil && a.opts.AuthUserID(r) == "" {
		http.Redirect(w, r, a.opts.AdminHomeURL, http.StatusSeeOther)
		return
	}

	// Use Routes() with registry to handle the request
	if a.opts.Registry == nil {
		http.Error(w, "Registry not configured", http.StatusInternalServerError)
		return
	}

	// Get routes and find matching one
	routes := Routes(a.opts.Registry)

	// Find matching route by path
	for _, route := range routes {
		if strings.HasPrefix(r.URL.Path, route.GetPath()) {
			// Execute the route's handler
			if handler := route.GetHandler(); handler != nil {
				handler(w, r)
				return
			}
			if htmlHandler := route.GetHTMLHandler(); htmlHandler != nil {
				htmlHandler(w, r)
				return
			}
		}
	}

	// No route matched - redirect to log manager
	http.Redirect(w, r, a.opts.LogAdminURL+"?controller=log-manager", http.StatusSeeOther)
}
