package routes

import (
	"file_project/controllers"
	"file_project/middleware"

	"github.com/gofiber/fiber/v2"
)

func FileRoutes(app *fiber.App, fc *controllers.FileController) {
	g := app.Group("/api/files", middleware.JWTProtected)
	g.Post("/upload", fc.Upload)
	g.Get("/:id/download", fc.Download) // password in query param ?password=...
	g.Patch("/:id/password", fc.ChangePassword)
	g.Delete("/:id", fc.Delete)
	g.Get("/", fc.List)
}
