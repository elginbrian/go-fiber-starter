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
	app.Get("/api", redirectToDocs)
	app.Get("/docs", redirectToDocs)
	app.Get("/docs/*", fiberSwagger.WrapHandler)

	setupUserRoutes(app, userHandler, jwtSecret)
	setupAuthRoutes(app, authHandler, jwtSecret)
	setupPostRoutes(app, postHandler, jwtSecret)
}

func redirectToDocs(c *fiber.Ctx) error {
	return c.Redirect("/docs/index.html")
}

func setupUserRoutes(app *fiber.App, handler *handler.UserHandler, jwtSecret string) { 
    userGroup := app.Group("/api/users")
    userGroup.Put("/:id", middleware.TokenValidationMiddleware(jwtSecret), handler.UpdateUser)
    userGroup.Get("/", handler.GetAllUsers)
    userGroup.Get("/:id", handler.GetUserByID)
    userGroup.Get("/search", handler.SearchUsers)
}

func setupAuthRoutes(app *fiber.App, handler *handler.AuthHandler, jwtSecret string) {
	authGroup := app.Group("/api/auth")
	authGroup.Post("/register", handler.Register)
	authGroup.Post("/login", handler.Login)
	authGroup.Post("/change-password", handler.ChangePassword, middleware.TokenValidationMiddleware(jwtSecret))
}

func setupPostRoutes(app *fiber.App, handler *handler.PostHandler, jwtSecret string) {
	postGroup := app.Group("/api/posts")
	postGroup.Post("/", middleware.TokenValidationMiddleware(jwtSecret), handler.CreatePost)
	postGroup.Put("/:id", middleware.TokenValidationMiddleware(jwtSecret), handler.UpdatePost)
	postGroup.Delete("/:id", middleware.TokenValidationMiddleware(jwtSecret), handler.DeletePost)
	postGroup.Get("/", handler.GetAllPosts)
	postGroup.Get("/:id", handler.GetPostByID)
	postGroup.Get("/user/:user_id", handler.GetPostsByUserID)
}