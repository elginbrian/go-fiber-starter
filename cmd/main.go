package main

import (
	"fiber-starter/config"
	"fiber-starter/internal/handler"
	"fiber-starter/internal/repository"
	"fiber-starter/internal/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load configuration
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
	userRepo := repository.NewUserRepository(config.DB)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Routes
	app.Get("/api/users", userHandler.GetAllUsers)

	// Start server
	port := config.GetEnv("PORT", "3000")
	log.Printf("Fiber Starter is running on http://localhost:%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start Fiber Starter: %v", err)
	}
}
