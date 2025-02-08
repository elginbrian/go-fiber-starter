package routes

import (
	"fiber-starter/internal/di"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, container di.Container, jwtSecret string) {
	setupStaticRoutes(app)
	setupDocsRoutes(app)
	setupRESTRoutes(app, container, jwtSecret)
	SetupGraphQLRoute(app, container)
	setupErrorRoutes(app)
}