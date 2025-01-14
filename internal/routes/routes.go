package routes

import (
	"fiber-starter/internal/handler"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func SetupRoutes(
	app *fiber.App,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	postHandler *handler.PostHandler,
) {
	app.Get("/api/docs/*", fiberSwagger.WrapHandler, func(c *fiber.Ctx) error {
		return c.Redirect("/api/docs/index.html")
	})

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

	// Post Routes
	postGroup := app.Group("/api/posts")
	postGroup.Get("/", postHandler.GetAllPosts)
	postGroup.Get("/:id", postHandler.GetPostByID)
	postGroup.Post("/", postHandler.CreatePost)
	postGroup.Put("/:id", postHandler.UpdatePost)
	postGroup.Delete("/:id", postHandler.DeletePost)
}
