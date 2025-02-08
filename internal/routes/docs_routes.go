package routes

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func setupDocsRoutes(app *fiber.App) {
	app.Get("/docs/*", fiberSwagger.WrapHandler)
	app.Get("/", redirectToDocs)
	app.Get("/api", redirectToDocs)
	app.Get("/docs", redirectToDocs)
}

func redirectToDocs(c *fiber.Ctx) error {
	return c.Redirect("/docs/index.html", fiber.StatusFound)
}