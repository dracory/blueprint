package auth

import (
	"project/app/config"
	"project/internal/platform/database"

	"github.com/dracory/base/router"
)

func Routes(cfg *config.Config, db *database.Database) router.GroupInterface {
	// Create controllers
	authController := NewAuthController(cfg, db)

	// Create auth group
	authGroup := router.NewGroup().SetPrefix("/auth")

	// Add routes to auth group
	loginGetRoute := router.NewRoute().
		SetMethod("GET").
		SetPath("/login").
		SetHandler(authController.Login)
	authGroup.AddRoute(loginGetRoute)

	loginPostRoute := router.NewRoute().
		SetMethod("POST").
		SetPath("/login").
		SetHandler(authController.Login)
	authGroup.AddRoute(loginPostRoute)

	registerGetRoute := router.NewRoute().
		SetMethod("GET").
		SetPath("/register").
		SetHandler(authController.Register)
	authGroup.AddRoute(registerGetRoute)

	registerPostRoute := router.NewRoute().
		SetMethod("POST").
		SetPath("/register").
		SetHandler(authController.Register)
	authGroup.AddRoute(registerPostRoute)

	logoutRoute := router.NewRoute().
		SetMethod("GET").
		SetPath("/logout").
		SetHandler(authController.Logout)
	authGroup.AddRoute(logoutRoute)

	return authGroup
}
