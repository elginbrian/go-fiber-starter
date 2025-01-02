package routes

import (
	"fiber-starter/internal/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler, authHandler *handler.AuthHandler) {
	// User Routes
	userGroup := app.Group("/api/users")
	userGroup.Get("/", userHandler.GetAllUsers)
	userGroup.Post("/", userHandler.CreateUser)
	userGroup.Put("/:id", userHandler.UpdateUser)
	userGroup.Delete("/:id", userHandler.DeleteUser)

	// Auth Routes
	authGroup := app.Group("/api/auth")
	authGroup.Post("/register", authHandler.Register)
	authGroup.Post("/login", authHandler.Login)
}
