package main

import (
	"log"

	"file_project/config"
	"file_project/controllers"
	"file_project/database"
	"file_project/repositories"
	"file_project/routes"
	"file_project/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config.Load()

	if err := database.Connect(); err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	app := fiber.New()

	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New())

	// Health
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendString("OK") })

	// Wire repositories and controllers
	userRepo := repositories.NewUserRepository(database.DB)
	authCtrl := &controllers.AuthController{Users: userRepo}

	fileRepo := repositories.NewFileRepository(database.DB)
	fileSvc := services.NewFileService(fileRepo)
	fileCtrl := &controllers.FileController{Files: fileSvc}

	// Register routes
	routes.AuthRoutes(app, authCtrl)
	routes.FileRoutes(app, fileCtrl)

	log.Printf("server running on :%s", config.C.AppPort)
	if err := app.Listen(":" + config.C.AppPort); err != nil {
		log.Fatal(err)
	}
}
