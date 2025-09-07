package routes

import (
	"file_project/controllers"
	"file_project/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, authCtrl *controllers.AuthController) {
	auth := app.Group("/api/auth")
	auth.Post("/register", authCtrl.Register)
	auth.Post("/login", authCtrl.Login)
	auth.Get("/me", middleware.JWTProtected, authCtrl.Me)
}
