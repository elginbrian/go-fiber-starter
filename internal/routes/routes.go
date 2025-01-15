package routes

import (
	"fiber-starter/internal/handler"
	"fiber-starter/internal/middleware"

	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func SetupRoutes(
	app *fiber.App,
	userHandler *handler.UserHandler,
	authHandler *handler.AuthHandler,
	postHandler *handler.PostHandler,
) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/api/docs/index.html")
	})
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.Redirect("/api/docs/index.html")
	})
	app.Get("/api/docs", func(c *fiber.Ctx) error {
		return c.Redirect("/api/docs/index.html")
	})

	app.Get("/api/docs/*", fiberSwagger.WrapHandler)

	setupUserRoutes(app, userHandler)
	setupAuthRoutes(app, authHandler)
	setupPostRoutes(app, postHandler)
}

func setupUserRoutes(app *fiber.App, handler *handler.UserHandler) {
	userGroup := app.Group("/api/users")
	userGroup.Post("/", middleware.TokenValidationMiddleware, handler.CreateUser)
	userGroup.Put("/:id", middleware.TokenValidationMiddleware, handler.UpdateUser)
	userGroup.Delete("/:id", middleware.TokenValidationMiddleware, handler.DeleteUser)
	userGroup.Get("/", handler.GetAllUsers)
}

func setupAuthRoutes(app *fiber.App, handler *handler.AuthHandler) {
	authGroup := app.Group("/api/auth")
	authGroup.Post("/register", handler.Register)
	authGroup.Post("/login", handler.Login)
}

func setupPostRoutes(app *fiber.App, handler *handler.PostHandler) {
	postGroup := app.Group("/api/posts")
	postGroup.Post("/", middleware.TokenValidationMiddleware, func(c *fiber.Ctx) error {
		if c.Accepts("multipart/form-data") == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Expected multipart/form-data",
			})
		}
		return handler.CreatePost(c)
	})
	postGroup.Put("/:id", middleware.TokenValidationMiddleware, handler.UpdatePost)
	postGroup.Delete("/:id", middleware.TokenValidationMiddleware, handler.DeletePost)
	postGroup.Get("/", handler.GetAllPosts)
	postGroup.Get("/:id", handler.GetPostByID)
}