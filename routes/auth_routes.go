package routes

import (
	"file_project/controllers"
	"file_project/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, authCtrl *controllers.AuthController) {
	auth := app.Group("/api/auth") // route group

	auth.Post("/register", authCtrl.Register)
	auth.Post("/login", authCtrl.Login)

	// Protected sub-group
	me := auth.Group("/me", middleware.JWTProtected)
	me.Get("/", authCtrl.Me)
}
