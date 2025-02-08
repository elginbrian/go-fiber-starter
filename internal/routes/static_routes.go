package routes

import "github.com/gofiber/fiber/v2"

func setupStaticRoutes(app *fiber.App) {
	app.Static("/uploads/posts", "./public/uploads/posts")
	app.Static("/uploads/profile", "./public/uploads/profile")
}