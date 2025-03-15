package auth

import (
	"log"
	"net/http"
	"time"

	"project/app/config"
	"project/app/controllers/base"
	"project/app/models"
	"project/internal/platform/database"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/uid"
)

// AuthController handles authentication requests
type AuthController struct {
	base.Controller
	userStore *models.UserStore
}

// NewAuthController creates a new AuthController instance
func NewAuthController(cfg *config.Config, db *database.Database) *AuthController {
	return &AuthController{
		Controller: *base.NewController(cfg, db),
		userStore:  models.NewUserStore(db),
	}
}

// Login handles the login page request
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	// Check if the request is a POST
	if r.Method == http.MethodPost {
		// Parse the form
		err := r.ParseForm()
		if err != nil {
			c.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		// Get the form values
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Authenticate the user
		user, err := c.userStore.Authenticate(r.Context(), email, password)
		if err != nil {
			// Create a view with an error message
			content := c.loginForm("Invalid email or password")
			c.Render(w, c.Layout("Dracory - Login", content))
			return
		}

		log.Println(user)

		// Create a session token
		token := uid.MicroUid()
		expires := time.Now().Add(24 * time.Hour)

		// Set the session cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    token,
			Expires:  expires,
			HttpOnly: true,
			Path:     "/",
			SameSite: http.SameSiteStrictMode,
			Secure:   c.Config.IsProduction(),
		})

		// In a real application, you would store the session in a database
		// For simplicity, we're just redirecting to the home page
		c.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Render the login form
	content := c.loginForm("")
	c.Render(w, c.Layout("Dracory - Login", content))
}

// Register handles the registration page request
func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	// Check if the request is a POST
	if r.Method == http.MethodPost {
		// Parse the form
		err := r.ParseForm()
		if err != nil {
			c.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		// Get the form values
		email := r.FormValue("email")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirm_password")
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")

		// Validate the form values
		if email == "" || password == "" || confirmPassword == "" || firstName == "" || lastName == "" {
			content := c.registerForm("All fields are required")
			c.Render(w, c.Layout("Dracory - Register", content))
			return
		}

		// Check if the passwords match
		if password != confirmPassword {
			content := c.registerForm("Passwords do not match")
			c.Render(w, c.Layout("Dracory - Register", content))
			return
		}

		// Check if the user already exists
		existingUser, err := c.userStore.FindByEmail(r.Context(), email)
		if err != nil {
			c.Error(w, "Error checking user", http.StatusInternalServerError)
			return
		}

		if existingUser != nil {
			content := c.registerForm("Email already exists")
			c.Render(w, c.Layout("Dracory - Register", content))
			return
		}

		// Create the user
		user := &models.User{
			Entity: models.Entity{
				ID: uid.MicroUid(),
			},
			Email:     email,
			Password:  password,
			FirstName: firstName,
			LastName:  lastName,
			IsActive:  true,
		}

		err = c.userStore.Create(r.Context(), user)
		if err != nil {
			c.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		// Redirect to the login page
		c.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Render the registration form
	content := c.registerForm("")
	c.Render(w, c.Layout("Dracory - Register", content))
}

// Logout handles the logout request
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Secure:   c.Config.IsProduction(),
	})

	// Redirect to the home page
	c.Redirect(w, r, "/", http.StatusSeeOther)
}

// loginForm creates a login form with an optional error message
func (c *AuthController) loginForm(errorMessage string) *hb.Tag {
	// Create an error alert if there's an error message
	var errorAlert *hb.Tag
	if errorMessage != "" {
		errorAlert = hb.NewDiv().Class("alert alert-danger").Text(errorMessage)
	} else {
		errorAlert = hb.NewDiv()
	}

	// Create the login form
	return hb.NewDiv().Class("row justify-content-center").AddChildren([]hb.TagInterface{
		hb.NewDiv().Class("col-md-6").AddChildren([]hb.TagInterface{
			hb.NewDiv().Class("card").AddChildren([]hb.TagInterface{
				hb.NewDiv().Class("card-header").Text("Login"),
				hb.NewDiv().Class("card-body").AddChildren([]hb.TagInterface{
					errorAlert,
					hb.NewForm().Attr("method", "post").Attr("action", "/login").AddChildren([]hb.TagInterface{
						hb.NewDiv().Class("mb-3").AddChildren([]hb.TagInterface{
							hb.NewLabel().Class("form-label").Attr("for", "email").Text("Email"),
							hb.NewInput().Class("form-control").Attr("type", "email").Attr("id", "email").Attr("name", "email").Attr("required", "required"),
						}),
						hb.NewDiv().Class("mb-3").AddChildren([]hb.TagInterface{
							hb.NewLabel().Class("form-label").Attr("for", "password").Text("Password"),
							hb.NewInput().Class("form-control").Attr("type", "password").Attr("id", "password").Attr("name", "password").Attr("required", "required"),
						}),
						hb.NewDiv().Class("d-grid gap-2").AddChildren([]hb.TagInterface{
							hb.NewButton().Class("btn btn-primary").Attr("type", "submit").Text("Login"),
						}),
					}),
				}),
				hb.NewDiv().Class("card-footer text-center").AddChildren([]hb.TagInterface{
					hb.NewP().Text("Don't have an account? ").AddChild(
						hb.NewA().Attr("href", "/register").Text("Register"),
					),
				}),
			}),
		}),
	})
}

// registerForm creates a registration form with an optional error message
func (c *AuthController) registerForm(errorMessage string) *hb.Tag {
	// Create an error alert if there's an error message
	var errorAlert *hb.Tag
	if errorMessage != "" {
		errorAlert = hb.NewDiv().Class("alert alert-danger").Text(errorMessage)
	} else {
		errorAlert = hb.NewDiv()
	}

	// Create the registration form
	return hb.NewDiv().Class("row justify-content-center").AddChildren([]hb.TagInterface{
		hb.NewDiv().Class("col-md-6").AddChildren([]hb.TagInterface{
			hb.NewDiv().Class("card").AddChildren([]hb.TagInterface{
				hb.NewDiv().Class("card-header").Text("Register"),
				hb.NewDiv().Class("card-body").AddChildren([]hb.TagInterface{
					errorAlert,
					hb.NewForm().Attr("method", "post").Attr("action", "/register").AddChildren([]hb.TagInterface{
						hb.NewDiv().Class("mb-3").AddChildren([]hb.TagInterface{
							hb.NewLabel().Class("form-label").Attr("for", "email").Text("Email"),
							hb.NewInput().Class("form-control").Attr("type", "email").Attr("id", "email").Attr("name", "email").Attr("required", "required"),
						}),
						hb.NewDiv().Class("row mb-3").AddChildren([]hb.TagInterface{
							hb.NewDiv().Class("col").AddChildren([]hb.TagInterface{
								hb.NewLabel().Class("form-label").Attr("for", "first_name").Text("First Name"),
								hb.NewInput().Class("form-control").Attr("type", "text").Attr("id", "first_name").Attr("name", "first_name").Attr("required", "required"),
							}),
							hb.NewDiv().Class("col").AddChildren([]hb.TagInterface{
								hb.NewLabel().Class("form-label").Attr("for", "last_name").Text("Last Name"),
								hb.NewInput().Class("form-control").Attr("type", "text").Attr("id", "last_name").Attr("name", "last_name").Attr("required", "required"),
							}),
						}),
						hb.NewDiv().Class("mb-3").AddChildren([]hb.TagInterface{
							hb.NewLabel().Class("form-label").Attr("for", "password").Text("Password"),
							hb.NewInput().Class("form-control").Attr("type", "password").Attr("id", "password").Attr("name", "password").Attr("required", "required"),
						}),
						hb.NewDiv().Class("mb-3").AddChildren([]hb.TagInterface{
							hb.NewLabel().Class("form-label").Attr("for", "confirm_password").Text("Confirm Password"),
							hb.NewInput().Class("form-control").Attr("type", "password").Attr("id", "confirm_password").Attr("name", "confirm_password").Attr("required", "required"),
						}),
						hb.NewDiv().Class("d-grid gap-2").AddChildren([]hb.TagInterface{
							hb.NewButton().Class("btn btn-primary").Attr("type", "submit").Text("Register"),
						}),
					}),
				}),
				hb.NewDiv().Class("card-footer text-center").AddChildren([]hb.TagInterface{
					hb.NewP().Text("Already have an account? ").AddChild(
						hb.NewA().Attr("href", "/login").Text("Login"),
					),
				}),
			}),
		}),
	})
}
