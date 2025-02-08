package routes

import (
	"fiber-starter/internal/di"
	"fiber-starter/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func setupRESTRoutes(app *fiber.App, container di.Container, jwtSecret string) {
	setupUserRoutes(app, container)
	setupAuthRoutes(app, container, jwtSecret)
	setupPostRoutes(app, container)
	setupSearchRoutes(app, container)
	setupCommentRoutes(app, container)
	setupLikeRoutes(app, container)
}

func setupLikeRoutes(app *fiber.App, container di.Container) {
	likeGroup := app.Group("/api/v1/posts")
	likeGroup.Post("/:post_id/like", container.LikeHandler.LikePost)
	likeGroup.Post("/:post_id/unlike", container.LikeHandler.UnlikePost)
	likeGroup.Get("/:post_id/likes", container.LikeHandler.GetLikesByPostID)

	userLikeGroup := app.Group("/api/v1/users")
	userLikeGroup.Get("/:user_id/likes",container.LikeHandler.GetLikesByUserID)
}

func setupSearchRoutes(app *fiber.App, container di.Container) {
	searchGroup := app.Group("/api/v1/search")
	searchGroup.Get("/users", container.UserHandler.SearchUsers)
	searchGroup.Get("/posts", container.PostHandler.SearchPosts)
}

func setupUserRoutes(app *fiber.App, container di.Container) {
	userGroup := app.Group("/api/v1/users")
	userGroup.Patch("/", container.UserHandler.UpdateUser)
	userGroup.Get("/", container.UserHandler.GetAllUsers)
	userGroup.Get("/:id", container.UserHandler.GetUserByID)
}

func setupAuthRoutes(app *fiber.App, container di.Container, jwtSecret string) {
	authGroup := app.Group("/api/v1/auth")
	authGroup.Post("/register", container.AuthHandler.Register)
	authGroup.Post("/login", container.AuthHandler.Login)
	authGroup.Post("/refresh-token", container.AuthHandler.RefreshToken)
	authGroup.Get("/current-user", container.AuthHandler.GetUserInfo, middleware.TokenValidationMiddleware(jwtSecret))
	authGroup.Patch("/change-password", container.AuthHandler.ChangePassword, middleware.TokenValidationMiddleware(jwtSecret))
}

func setupPostRoutes(app *fiber.App, container di.Container) {
	postGroup := app.Group("/api/v1/posts")
	postGroup.Post("/", container.PostHandler.CreatePost)
	postGroup.Patch("/:id", container.PostHandler.UpdatePost)
	postGroup.Delete("/:id", container.PostHandler.DeletePost)
	postGroup.Get("/", container.PostHandler.GetAllPosts)
	postGroup.Get("/:id", container.PostHandler.GetPostByID)
	postGroup.Get("/user/:user_id", container.PostHandler.GetPostsByUserID)
}

func setupCommentRoutes(app *fiber.App, container di.Container) {
	commentGroup := app.Group("/api/v1/posts/:post_id/comments")
	commentGroup.Get("/", container.CommentHandler.GetCommentsByPostID)
	commentGroup.Post("/", container.CommentHandler.CreateComment)
	commentGroup.Delete("/:id", container.CommentHandler.DeleteComment)
}