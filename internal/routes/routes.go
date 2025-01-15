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
    jwtSecret string,
) {
    app.Get("/", redirectToDocs)  
    app.Get("/api", redirectToDocs) 
    app.Get("/docs", redirectToDocs)
    app.Get("/docs/*", fiberSwagger.WrapHandler)

    app.Static("/docs", "./docs")

    setupUserRoutes(app, userHandler, jwtSecret)
    setupAuthRoutes(app, authHandler)
    setupPostRoutes(app, postHandler, jwtSecret)
}

func redirectToDocs(c *fiber.Ctx) error {
    return c.Redirect("/docs")
}


func setupUserRoutes(app *fiber.App, handler *handler.UserHandler, jwtSecret string) {
	userGroup := app.Group("/api/users")
	userGroup.Post("/", middleware.TokenValidationMiddleware(jwtSecret), handler.CreateUser)
	userGroup.Put("/:id", middleware.TokenValidationMiddleware(jwtSecret), handler.UpdateUser)
	userGroup.Delete("/:id", middleware.TokenValidationMiddleware(jwtSecret), handler.DeleteUser)
	userGroup.Get("/", handler.GetAllUsers)
}

func setupAuthRoutes(app *fiber.App, handler *handler.AuthHandler) {
	authGroup := app.Group("/api/auth")
	authGroup.Post("/register", handler.Register)
	authGroup.Post("/login", handler.Login)
}

func setupPostRoutes(app *fiber.App, handler *handler.PostHandler, jwtSecret string) {
	postGroup := app.Group("/api/posts")
	postGroup.Post("/", middleware.TokenValidationMiddleware(jwtSecret), func(c *fiber.Ctx) error {
		if c.Accepts("multipart/form-data") == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Expected multipart/form-data",
			})
		}
		return handler.CreatePost(c)
	})
	postGroup.Put("/:id", middleware.TokenValidationMiddleware(jwtSecret), handler.UpdatePost)
	postGroup.Delete("/:id", middleware.TokenValidationMiddleware(jwtSecret), handler.DeletePost)
	postGroup.Get("/", handler.GetAllPosts)
	postGroup.Get("/:id", handler.GetPostByID)
}