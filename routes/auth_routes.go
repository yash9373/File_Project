package routes

import (
	"file_project/controllers"
	"file_project/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	auth := app.Group("/api/auth") // route group

	auth.Post("/register", controllers.Register)
	auth.Post("/login", controllers.Login)

	// Protected sub-group
	me := auth.Group("/me", middleware.JWTProtected)
	me.Get("/", controllers.Me)
}
