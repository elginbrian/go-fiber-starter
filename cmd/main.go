package main

import (
	"fiber-starter/config"
	"fiber-starter/internal/di"
	"fiber-starter/internal/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to the database
	config.ConnectDatabase()
	defer config.DB.Close()

	// Create a new Fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Dependency Injection
	container := di.NewContainer()

	// Setup Routes
	routes.SetupRoutes(app, container.UserHandler, container.AuthHandler)

	// Start the server
	port := config.GetEnv("PORT", "3000")
	log.Printf("Fiber Starter is running on http://localhost:%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start Fiber Starter: %v", err)
	}
}
