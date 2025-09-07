package routes

import (
	"file_project/controllers"
	"file_project/middleware"

	"github.com/gofiber/fiber/v2"
)

// Share routes
func ShareRoutes(app *fiber.App, sc *controllers.ShareController) {
	// Owner creates/deletes share links (protected)
	g := app.Group("/api/share", middleware.JWTProtected)
	g.Post("/", sc.Create)
	g.Delete("/:token", sc.Delete)

	// Public download (no JWT), but still needs password
	app.Get("/share/:token/download", sc.PublicDownload)
}
