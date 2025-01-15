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
	userGroup.Get("/", handler.GetAllUsers)
	userGroup.Post("/", handler.CreateUser)
	userGroup.Put("/:id", handler.UpdateUser)
	userGroup.Delete("/:id", handler.DeleteUser)
}

func setupAuthRoutes(app *fiber.App, handler *handler.AuthHandler) {
	authGroup := app.Group("/api/auth")
	authGroup.Post("/register", handler.Register)
	authGroup.Post("/login", handler.Login)
}

func setupPostRoutes(app *fiber.App, handler *handler.PostHandler) {
	postGroup := app.Group("/api/posts")
	postGroup.Get("/", handler.GetAllPosts)
	postGroup.Get("/:id", handler.GetPostByID)
	postGroup.Post("/", handler.CreatePost)
	postGroup.Put("/:id", handler.UpdatePost)
	postGroup.Delete("/:id", handler.DeletePost)
}